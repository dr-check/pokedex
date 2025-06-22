package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	pokecache "github.com/dr-check/pokedex/internal"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type Client struct {
	cache      *pokecache.Cache
	httpClient http.Client
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache(cacheInterval),
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

type Pokemon struct {
	Id             int            `json:"id"`
	Name           string         `json:"name"`
	BaseExperience int            `json:"base_experience"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	Stats          []PokemonStats `json:"stats"`
	Types          []PokemonTypes `json:"types"`
}

type PokemonStats struct {
	Stat     LocationAreaResource `json:"stat"`
	Effort   int                  `json:"effort"`
	BaseStat int                  `json:"base_stat"`
}

type PokemonTypes struct {
	Slot int                  `json:"slot"`
	Type LocationAreaResource `json:"type"`
}

func (c *Client) ListLocations(pageURL *string) (LocationArea, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	if val, ok := c.cache.Get(url); ok {
		locationsResp := LocationArea{}
		err := json.Unmarshal(val, &locationsResp)
		if err != nil {
			return LocationArea{}, err
		}

		return locationsResp, nil
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
	c.cache.Add(url, data)
	return locationAreas, nil
}

func (c *Client) GetLocationInfo(locationName string) (LocationArea, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", locationName)

	if val, ok := c.cache.Get(url); ok {
		locationsResp := LocationArea{}
		err := json.Unmarshal(val, &locationsResp)
		if err != nil {
			return LocationArea{}, err
		}

		return locationsResp, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error fetching location info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LocationArea{}, fmt.Errorf("location not found: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error reading response body: %w", err)
	}

	var location LocationArea
	err = json.Unmarshal(body, &location)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error unmarshalling location data: %w", err)
	}

	c.cache.Add(url, body)
	return location, nil
}

func (c *Client) GetPokemonInfo(pokemonName string) (Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)

	if val, ok := c.cache.Get(url); ok {
		locationsResp := Pokemon{}
		err := json.Unmarshal(val, &locationsResp)
		if err != nil {
			return Pokemon{}, err
		}

		return locationsResp, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error fetching creature info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Pokemon{}, fmt.Errorf("creature not found: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error reading response body: %w", err)
	}

	var foundPokemon Pokemon
	err = json.Unmarshal(body, &foundPokemon)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error unmarshalling creature data: %w", err)
	}

	c.cache.Add(url, body)
	return foundPokemon, nil
}
