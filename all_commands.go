package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	cache "github.com/josephus-git/pokedex/pokecache"
)

func allcommands(conf *config, myCache *cache.Cache) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for key := range getcommands(conf, myCache) {
		fmt.Printf("%s: %s\n", getcommands(conf, myCache)[key].name, getcommands(conf, myCache)[key].description)
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
	if conf.Result[0].Name == "canalave-city-area" {
		return fmt.Errorf("you're on the first page")
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

	// check in cache if next url already available
	// if available, move to unmarshal, otherwise, get url

	cacheEn, ok := myCache.MapCache[url]
	var body []byte
	if ok {
		body = cacheEn.Val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		body, err = io.ReadAll(res.Body)
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return fmt.Errorf("status code: %d", res.StatusCode)
		}
		if err != nil {
			return err
		}
		myCache.Add(url, body)
	}

	// conf is just a temporary struct which stores details of the maps obtained
	err := json.Unmarshal(body, &conf)
	if err != nil {
		fmt.Println(err)
	}

	for _, result := range conf.Result {
		fmt.Println(result.Name)
	}

	return nil
}

func getcommands(conf *config, myCache *cache.Cache) map[string]cliCommand {
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
				return allcommands(conf, myCache)
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
	}
	return commands
}
