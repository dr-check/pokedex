package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dr-check/pokedex/pokeapi"
)

type config struct {
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
			}
			break
		}
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		command, ok := commands[cleanedInput[0]]
		if !ok {
			fmt.Println("Unknown command:", cleanedInput[0])
			continue
		}
		switch cleanedInput[0] {
		case "explore":
			if len(cleanedInput) < 2 {
				fmt.Println("Please provide a location name to explore.")
				continue
			}
			command.callback(cleanedInput[1:], cfg)
		case "catch":
			if len(cleanedInput) < 2 {
				fmt.Println("Please find a Pokemon to capture!")
			}
			command.callback(cleanedInput[1:], cfg)
		default:
			if err := command.callback(cleanedInput, cfg); err != nil {
				fmt.Fprintln(os.Stderr, "Error executing command:", err)
			}
		}
	}
}
