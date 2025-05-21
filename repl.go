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
	currentConfig := config{}
	configptr := &currentConfig
	const interval = 1 * time.Minute
	myCache := cache.NewCache(interval)

	for i := 0; ; i++ {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		cmd, ok := getcommands(configptr, myCache)[input[0]]
		if !ok {
			fmt.Print("Unknown command\n")
			continue
		}
		err := cmd.callback()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
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
