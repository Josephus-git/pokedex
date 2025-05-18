package main

import (
	"fmt"
	"os"
)

func allcommands() error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for key := range getcommands() {
		fmt.Printf("%s: %s\n", getcommands()[key].name, getcommands()[key].description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func getcommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    allcommands,
		},
	}
	return commands
}
