package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	cache "github.com/josephus-git/pokedex/pokecache"
)

func startR() {
	scanner := bufio.NewScanner(os.Stdin)
	configptr := &config{}
	pokedexMade := make(map[string]pokeStruct)
	pokedex := &pokedexMade
	const interval = 1 * time.Minute
	myCache := cache.NewCache(interval)

	for i := 0; ; i++ {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		argument2 := ""
		if len(input) > 1 {
			argument2 = input[1]
		}

		cmd, ok := getcommands(configptr, myCache, argument2, pokedex)[input[0]]
		if !ok {
			fmt.Print("Unknown command\n")
			continue
		}
		err := cmd.callback()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		println("___________")
		println("")
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func cleanInput(text string) []string {
	words := []string{}
	sText := strings.Split(strings.TrimSpace(text), " ")
	for _, word := range sText {
		if len(word) != 0 {
			words = append(words, strings.ToLower(word))
		}
	}
	return words
}
