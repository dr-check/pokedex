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
		if err := command.callback(cfg); err != nil {
			fmt.Fprintln(os.Stderr, "Error executing command:", err)
		}
		if cleanedInput[0] == "exit" {
			command.callback(cfg)
		}
	}
}
