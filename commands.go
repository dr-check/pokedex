package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/dr-check/pokedex/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string, cfg *config) error
}

var Pokedex = make(map[string]pokeapi.Pokemon)

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Navigates through the maps of the Pokedex",
		callback:    commandMapf,
	},
	"mapb": {
		name:        "mapb",
		description: "Get the previous page of locations",
		callback:    commandMapb,
	},
	"explore": {
		name:        "explore",
		description: "Explore a specific location",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Attempt to capture a wild Pokemon",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "Inspect a Pokemon that you own",
		callback:    commandInspect,
	},
	"pokedex": {
		name:        "pokedex",
		description: "Check the list of Pokemon you caught",
		callback:    commandPokedex,
	},
}

func commandExit(args []string, cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(args []string, cfg *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func commandMapf(args []string, cfg *config) error {
	locationAreas, err := cfg.pokeapiClient.ListLocations(cfg.nextLocationsURL)
	if err != nil {
		return err
	}
	cfg.nextLocationsURL = locationAreas.Next
	cfg.prevLocationsURL = locationAreas.Previous

	for _, loc := range locationAreas.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(args []string, cfg *config) error {
	if cfg.prevLocationsURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locationAreas, err := cfg.pokeapiClient.ListLocations(cfg.prevLocationsURL)
	if err != nil {
		return err
	}
	cfg.nextLocationsURL = locationAreas.Next
	cfg.prevLocationsURL = locationAreas.Previous

	for _, loc := range locationAreas.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExplore(locationName []string, cfg *config) error {
	if len(locationName) < 1 {
		fmt.Println("Please provide a location name to explore.")
		return nil
	}
	location, err := cfg.pokeapiClient.GetLocationInfo(locationName[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching location info: %v\n", err)
	}
	fmt.Printf("Exploring %s...\n", locationName[0])
	fmt.Println("Found Pokemon")
	for _, encounter := range location.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func CaptureChance(baseXP int) int {
	capture := ((1000 - baseXP) / 10) - 10
	return capture
}

func commandCatch(foundPokemon []string, cfg *config) error {
	if len(foundPokemon) < 1 {
		fmt.Println("Pick a pokemon to capture!")
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", foundPokemon[0])

	targetPokemon, err := cfg.pokeapiClient.GetPokemonInfo(foundPokemon[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding that creature: %v\n", err)
		return nil
	}
	randomNumber := rand.Intn(100) + 1
	catchRate := CaptureChance(targetPokemon.BaseExperience)

	if randomNumber <= catchRate {
		fmt.Printf("%s was caught!\n", targetPokemon.Name)
		if _, ok := Pokedex[targetPokemon.Name]; !ok {
			Pokedex[targetPokemon.Name] = targetPokemon
		} else {
			fmt.Printf("%s was caught!\n", targetPokemon.Name)
		}
	} else {
		fmt.Printf("%s escaped!\n", targetPokemon.Name)
	}
	return nil
}

func commandInspect(ownedPokemon []string, cfg *config) error {
	_, ok := Pokedex[ownedPokemon[0]]
	if ok {
		partner := Pokedex[ownedPokemon[0]]
		fmt.Printf("Name: %s\n", partner.Name)
		fmt.Printf("Height: %v\n", partner.Height)
		fmt.Printf("Weight: %v\n", partner.Weight)
		fmt.Println("Stats:")
		for _, stat := range partner.Stats {
			fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Type:")
		for _, t := range partner.Types {
			fmt.Printf("  -%s\n", t.Type.Name)
		}
	}
	return nil
}

func commandPokedex(ownedPokemon []string, cfg *config) error {
	if len(Pokedex) == 0 {
		fmt.Println("You haven't caught any Pokemon yet!")
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range Pokedex {
		fmt.Printf("-%s\n", pokemon.Name)
	}
	return nil
}
