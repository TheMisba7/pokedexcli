package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/cache"
	"pokedexcli/model"
)

var cacheManager = cache.NewCache(7)

type PokeApiImpl struct{}

func (api PokeApiImpl) GetLocationArea(url string) model.LocationAreaRS {
	areas := model.LocationAreaRS{}
	err := getJson(url, &areas)
	if err != nil {
		panic(err)
	}
	return areas
}

func (api PokeApiImpl) GetPokemons(area string) model.PokemonRS {
	pokemons := model.PokemonRS{}
	err := getJson(area, &pokemons)
	if err != nil {
		panic(err)
	}
	return pokemons
}

func (api PokeApiImpl) GetPokemon(url string) model.Pokemon {
	pke := model.Pokemon{}
	err := getJson(url, &pke)
	if err != nil {
		panic(err)
	}
	return pke
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
