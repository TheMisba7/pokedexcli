package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/cache"
	"pokedexcli/model"
)

var cacheManager = cache.NewCache(60)

type PokeApiImpl struct{}

func (api PokeApiImpl) GetLocationArea(url string) (model.LocationAreaRS, error) {
	areas := model.LocationAreaRS{}
	err := httpGet(url, &areas)
	if err != nil {
		return areas, err
	}
	return areas, nil
}

func (api PokeApiImpl) GetListOfPokemon(area string) (model.PokemonRS, error) {
	pokemonList := model.PokemonRS{}
	err := httpGet(area, &pokemonList)
	if err != nil {
		return pokemonList, err
	}
	return pokemonList, nil
}

func (api PokeApiImpl) GetPokemon(url string) (model.Pokemon, error) {
	pke := model.Pokemon{}
	err := httpGet(url, &pke)
	if err != nil {
		return pke, err
	}
	return pke, nil
}

func httpGet(url string, target interface{}) error {
	val, found := cacheManager.Get(url)
	if found {
		return json.Unmarshal(val, target)
	}
	r, err := http.Get(url)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	if r.StatusCode > 399 {
		return fmt.Errorf("bad status code: %v", r.StatusCode)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	go cacheManager.Add(url, body)
	return json.Unmarshal(body, target)
}
