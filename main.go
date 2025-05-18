package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for i := 0; ; i++ {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		fmt.Printf("Your command was: %s\n", input[0])
	}
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
