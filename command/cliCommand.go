package command

import (
	"fmt"
	"pokedexcli/api"
)

const locationAreaURL = "https://pokeapi.co/api/v2/location-area"

var pokeapi = api.PokeApiImpl{}
var Commands map[string]command

type command struct {
	name        string
	description string
	Callback    func([]string) error
	config      *config
}

type config struct {
	next     string
	previous string
}

func InitCommands() map[string]command {
	config := &config{next: locationAreaURL}
	return map[string]command{
		"help": {
			name:        "help",
			description: "this command will help you to learn about pokedex",
			Callback:    helpCallback,
		},
		"exit": {
			name:        "exit",
			description: "exist the pokedex program",
			Callback:    exitCallBack,
		},
		"map": {
			name:        "map",
			description: "get next 20 location areas",
			Callback:    mapCallBack,
			config:      config,
		},
		"mapb": {
			name:        "mapb",
			description: "map back to previous page",
			Callback:    mapBackCallBack,
			config:      config,
		},
		"explore": {
			name:        "explore",
			description: "list pokemons within a location area. \n usage: explore <area_name>.",
			Callback:    exploreCallBack,
		},
	}
}

func exploreCallBack(params []string) error {
	area := locationAreaURL + "/" + params[0]
	pokemon := pokeapi.GetPokemon(area)
	for _, pke := range pokemon.PokemonEncounters {
		fmt.Println(pke.Pokemon.Name)
	}
	return nil
}

func mapBackCallBack(params []string) error {
	mapCommand := Commands["mapb"]
	if mapCommand.config.previous == "" {
		return fmt.Errorf("no previous page to show")
	}
	area := pokeapi.GetLocationArea(mapCommand.config.previous)
	for _, result := range area.Results {
		fmt.Println(result.Name)
	}
	mapCommand.config.next = area.Next
	mapCommand.config.previous = area.Previous
	return nil
}

func mapCallBack(params []string) error {
	mapCommand := Commands["map"]
	fmt.Println("next: ", mapCommand.config.next)
	area := pokeapi.GetLocationArea(mapCommand.config.next)
	for _, result := range area.Results {
		fmt.Println(result.Name)
	}
	mapCommand.config.next = area.Next
	mapCommand.config.previous = area.Previous
	return nil
}

func exitCallBack(params []string) error {
	return nil
}

func helpCallback(params []string) error {
	for _, command := range Commands {
		fmt.Println("name: " + command.name)
		fmt.Println("description: " + command.description)
	}
	return nil
}
