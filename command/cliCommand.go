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
	Callback    func() error
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
	}
}

func mapBackCallBack() error {
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

func mapCallBack() error {
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

func exitCallBack() error {
	return nil
}

func helpCallback() error {
	for _, command := range Commands {
		fmt.Println("name: " + command.name)
		fmt.Println("description: " + command.description)
	}
	return nil
}
