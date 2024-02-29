package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/cache"
)

var cacheManager = cache.NewCache(7)

type PokeApi interface {
	GetLocationArea(url string) LocationAreaRS
	GetPokemon(area string) PokemonRS
}

type Area struct {
	name string
	url  string
}

type PokeApiImpl struct{}
type LocationAreaRS struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		URL  string
	}
}

type PokemonRS struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (api PokeApiImpl) GetLocationArea(url string) LocationAreaRS {
	areas := LocationAreaRS{}
	err := getJson(url, &areas)
	if err != nil {
		panic(err)
	}
	return areas
}

func (api PokeApiImpl) GetPokemon(area string) PokemonRS {
	pokemons := PokemonRS{}
	err := getJson(area, &pokemons)
	if err != nil {
		panic(err)
	}
	return pokemons
}

func getJson(url string, target interface{}) error {
	val, found := cacheManager.Get(url)
	if found {
		return json.Unmarshal(val, target)
	}
	r, err := http.Get(url)
	fmt.Println(r.Status)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	go cacheManager.Add(url, body)
	return json.Unmarshal(body, target)
}
