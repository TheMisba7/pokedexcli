package command

import (
	"fmt"
	"math/rand"
	"pokedexcli/api"
	"pokedexcli/model"
)

const locationAreaURL = "https://pokeapi.co/api/v2/location-area"
const pokemonUrl = "https://pokeapi.co/api/v2/pokemon/"

var catchedPokemon = make(map[string]model.Pokemon)
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
		"catch": {
			name:        "catch",
			description: "Used to catch a pokemon. \n usage: catch <pokemon_name>",
			Callback:    catchCallBack,
		},
		"inspect": {
			name:        "inspect",
			description: "It takes the name of a Pokemon as an argument. It should print the name, height, weight, stats and type(s) of the Pokemon",
			Callback:    inspectCallBack,
		},
	}
}

func inspectCallBack(params []string) error {
	pokeName := params[0]
	pokemon, ok := catchedPokemon[pokeName]
	if !ok {
		fmt.Println("you have not caught ", pokeName)
		return nil
	}
	fmt.Println("Name: ", pokeName)
	fmt.Println("Base experience: ", pokemon.BaseExperience)

	return nil
}

func catchCallBack(param []string) error {
	pokeUrl := pokemonUrl + param[0]
	fmt.Println(fmt.Sprintf("Throwing a Pokeball at %s ...", param[0]))
	pokemon := pokeapi.GetPokemon(pokeUrl)
	if pokemon.BaseExperience < rand.Intn(100) {
		fmt.Println(fmt.Sprintf("%s was caught!", param[0]))
		catchedPokemon[param[0]] = pokemon
	} else {
		fmt.Println(fmt.Sprintf("%s escaped!", param[0]))
	}
	fmt.Println(pokemon.BaseExperience)
	return nil
}

func exploreCallBack(params []string) error {
	area := locationAreaURL + "/" + params[0]
	pokemon := pokeapi.GetPokemons(area)
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
