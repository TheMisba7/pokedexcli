package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
	r, err := http.Get(url)
	fmt.Println(r.Status)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	return json.Unmarshal(body, target)
}
