package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func allcommands(conf *config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for key := range getcommands(conf) {
		fmt.Printf("%s: %s\n", getcommands(conf)[key].name, getcommands(conf)[key].description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func mapsbackwards(conf *config) error {
	//take note this logic may affect later on!
	if conf.Previous == "" {
		return fmt.Errorf("you're on the first page")
	}
	conf.Next = conf.Previous
	mapsTwenty(conf)
	return nil
}

func mapsTwenty(conf *config) error {
	url := ""
	if conf.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = conf.Next
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &conf)
	if err != nil {
		fmt.Println(err)
	}

	for _, result := range conf.Result {
		fmt.Println(result.Name)
	}

	return nil
}

func getcommands(conf *config) map[string]cliCommand {
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
				return allcommands(conf)
			},
		},
		"map": {
			name:        "map",
			description: "get map list of 20 locations in pokemon",
			callback: func() error {
				return mapsTwenty(conf)
			},
		},
		"mapb": {
			name:        "mapb",
			description: "go backwards to get previous 20 location in pokemon",
			callback: func() error {
				return mapsbackwards(conf)
			},
		},
	}
	return commands
}
