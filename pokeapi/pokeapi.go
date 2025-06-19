package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

type LocationArea struct {
	Id                  int                   `json:"id"`
	Name                string                `json:"name"`
	GameIndex           int                   `json:"game_index"`
	EncounterMethodRate []EncounterMethodRate `json:"encounter_method_rates"`
	Location            LocationAreaResource  `json:"location"`
	PokemonEncounters   []PokemonEncounter    `json:"pokemon_encounters"`
	Next                *string               `json:"next"`
	Previous            *string               `json:"previous"`
	Results             []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
}

type LocationAreaResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type EncounterMethodRate struct {
	EncounterMethod LocationAreaResource `json:"encounter_method"`
	VersionDetails  []VersionDetail      `json:"version_details"`
}

type VersionDetail struct {
	Rate    int                  `json:"rate"`
	Version LocationAreaResource `json:"version"`
}

type PokemonEncounter struct {
	Pokemon        LocationAreaResource `json:"pokemon"`
	VersionDetails []VersionDetail      `json:"version_details"`
}

type VersionEncounter struct {
	Version          LocationAreaResource `json:"version"`
	MaxChance        int                  `json:"max_chance"`
	EncounterDetails []EncounterDetail    `json:"encounter_details"`
}

type EncounterDetail struct {
	MinLevel        int                    `json:"min_level"`
	MaxLevel        int                    `json:"max_level"`
	ConditionValues []LocationAreaResource `json:"condition_values"`
	Chance          int                    `json:"chance"`
	Method          LocationAreaResource   `json:"method"`
}

func (c *Client) ListLocations(pageURL *string) (LocationArea, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching location areas: %v\n", err)
		return LocationArea{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return LocationArea{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
		return LocationArea{}, err
	}

	var locationAreas LocationArea
	err = json.Unmarshal(data, &locationAreas)
	if err != nil {
		return LocationArea{}, err
	}

	return locationAreas, nil
}
