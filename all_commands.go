package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	cache "github.com/josephus-git/pokedex/pokecache"
)

func allcommands(conf *config, myCache *cache.Cache, location string, pokedex *map[string]pokeStruct) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for key := range getcommands(conf, myCache, location, pokedex) {
		fmt.Printf("%s: %s\n", getcommands(conf, myCache, location, pokedex)[key].name, getcommands(conf, myCache, location, pokedex)[key].description)
	}
	return nil
}

func checkCache(myCache *cache.Cache, fullUrl string) ([]byte, error) {
	// check in cache if next url already available
	// if available, move to unmarshal, otherwise, get url

	cacheEn, ok := myCache.MapCache[fullUrl]
	var body []byte
	if ok {
		body = cacheEn.Val
	} else {
		res, err := http.Get(fullUrl)
		if err != nil {
			return nil, fmt.Errorf("location not found: %s", err)
		}
		body, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return nil, fmt.Errorf("status code: %d", res.StatusCode)
		}
		if err != nil {
			return nil, err
		}
		myCache.Add(fullUrl, body)
	}
	return body, nil
}

func listPokemon(pokedex *map[string]pokeStruct) error {
	if len(*pokedex) > 0 {
		fmt.Println("Your Pokedex:")
		for key := range *pokedex {
			fmt.Printf("- %s\n", key)
		}
		return nil
	}
	fmt.Println("You havent caught any pokemon")
	return nil
}

func inspectPokemon(pokemon string, pokedex *map[string]pokeStruct) error {
	pokemonStats, ok := (*pokedex)[pokemon]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", pokemonStats.Name)
	fmt.Printf("Height: %d\n", pokemonStats.Height)
	fmt.Printf("Weight: %d\n", pokemonStats.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemonStats.Stats {
		fmt.Printf("	-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typ := range pokemonStats.Types {
		fmt.Printf("	-%s\n", typ.Type.Name)
	}

	return nil
}

func catchPokemon(myCache *cache.Cache, pokemon string, pokedex *map[string]pokeStruct) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
	pokeDetails := pokeStruct{}
	url := "https://pokeapi.co/api/v2/pokemon/"
	fullUrl := url + pokemon

	body, err := checkCache(myCache, fullUrl)
	if err != nil {
		return err
	}

	// conf is just a temporary struct which stores details of the maps obtained
	err = json.Unmarshal(body, &pokeDetails)
	if err != nil {
		fmt.Println(err)
	}

	prob := rand.Intn(6) + 1

	chance := pokeDetails.BaseExperience
	if chance > 100 {
		if prob > 2 {
			addPokemonToPokedex(pokemon, pokeDetails, pokedex)
			fmt.Printf("%s was caught!\n", pokemon)
		} else {
			fmt.Printf("%s escaped!\n", pokemon)
		}
	} else if chance > 50 {
		if prob > 1 {
			addPokemonToPokedex(pokemon, pokeDetails, pokedex)
			fmt.Printf("%s was caught!\n", pokemon)
		} else {
			fmt.Printf("%s escaped!\n", pokemon)
		}
	} else if chance > 0 {
		if prob > 0 {
			addPokemonToPokedex(pokemon, pokeDetails, pokedex)
			fmt.Printf("%s was caught!\n", pokemon)
		} else {
			fmt.Printf("%s escaped!\n", pokemon)
		}
	} else {
		return fmt.Errorf("base experience is not a valid integer")
	}
	return nil
}
func addPokemonToPokedex(pokemon string, pokeDetails pokeStruct, pokedex *map[string]pokeStruct) {
	(*pokedex)[pokemon] = pokeDetails
}

func listOfPokemons(myCache *cache.Cache, location string) error {
	fmt.Printf("Exploring %s...\n", location)
	myPoke := pokelocation{}
	url := "https://pokeapi.co/api/v2/location-area/"
	fullurl := url + location

	body, err := checkCache(myCache, fullurl)
	if err != nil {
		return err
	}

	// conf is just a temporary struct which stores details of the maps obtained
	err = json.Unmarshal(body, &myPoke)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Found Pokemon:")
	for _, pokemon := range myPoke.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func mapsbackwards(conf *config, myCache *cache.Cache) error {
	//take note this logic may affect later on!
	if len(conf.Result) < 1 {
		fmt.Println("you havent seached any map")
		return nil
	}
	if conf.Result[0].Name == "canalave-city-area" {
		fmt.Println("you're on the first page")
		return nil
	}
	conf.Next = conf.Previous
	mapsTwenty(conf, myCache)
	return nil
}

func mapsTwenty(conf *config, myCache *cache.Cache) error {
	url := ""
	if conf.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = conf.Next
	}

	body, err := checkCache(myCache, url)
	if err != nil {
		return err
	}

	// conf is just a temporary struct which stores details of the maps obtained
	err = json.Unmarshal(body, &conf)
	if err != nil {
		fmt.Println(err)
	}

	for _, result := range conf.Result {
		fmt.Println(result.Name)
	}

	return nil
}

func getcommands(conf *config, myCache *cache.Cache, argument2 string, pokedex *map[string]pokeStruct) map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback: func() error {
				return commandExit()
			},
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback: func() error {
				return allcommands(conf, myCache, argument2, pokedex)
			},
		},
		"map": {
			name:        "map",
			description: "get map list of 20 locations in pokemon",
			callback: func() error {
				return mapsTwenty(conf, myCache)
			},
		},
		"mapb": {
			name:        "mapb",
			description: "go backwards to get previous 20 location in pokemon",
			callback: func() error {
				return mapsbackwards(conf, myCache)
			},
		},
		"explore": {
			name:        "explore",
			description: "returns all pokemons in the given location",
			callback: func() error {
				return listOfPokemons(myCache, argument2)
			},
		},
		"catch": {
			name:        "catch",
			description: "attempt to catch pikachu",
			callback: func() error {
				return catchPokemon(myCache, argument2, pokedex)
			},
		},
		"inspect": {
			name:        "inspect",
			description: "inspect stats of pokemon",
			callback: func() error {
				return inspectPokemon(argument2, pokedex)
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "see all captured pokemon",
			callback: func() error {
				return listPokemon(pokedex)
			},
		},
	}
	return commands
}
