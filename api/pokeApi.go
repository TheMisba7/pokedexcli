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

func (api PokeApiImpl) GetLocationArea(url string) LocationAreaRS {
	areas := LocationAreaRS{}
	err := getJson(url, &areas)
	if err != nil {
		panic(err)
	}
	return areas
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
