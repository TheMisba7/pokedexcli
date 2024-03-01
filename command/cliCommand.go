package command

import (
	"fmt"
	"math/rand"
	"os"
	"pokedexcli/api"
	"pokedexcli/model"
)

const locationAreaURL = "https://pokeapi.co/api/v2/location-area"
const pokemonUrl = "https://pokeapi.co/api/v2/pokemon/"

var caughtPokemon = make(map[string]model.Pokemon)
var pokeApi model.PokeApi = api.PokeApiImpl{}
var Commands map[string]model.Command

func InitCommands() map[string]model.Command {
	config := &model.Config{Next: locationAreaURL}
	return map[string]model.Command{
		"help": {
			Name:        "help",
			Description: "this command will help you to learn about pokedex",
			Callback:    helpCallback,
		},
		"exit": {
			Name:        "exit",
			Description: "exist the pokedex program",
			Callback:    exitCallBack,
		},
		"map": {
			Name:        "map",
			Description: "get next 20 location areas",
			Callback:    mapCallBack,
			Config:      config,
		},
		"mapb": {
			Name:        "mapb",
			Description: "map back to previous page",
			Callback:    mapBackCallBack,
			Config:      config,
		},
		"explore": {
			Name:        "explore",
			Description: "list pokemons within a location area. \n usage: explore <area_name>.",
			Callback:    exploreCallBack,
		},
		"catch": {
			Name:        "catch",
			Description: "Used to catch a pokemon. \n usage: catch <pokemon_name>",
			Callback:    catchCallBack,
		},
		"inspect": {
			Name:        "inspect",
			Description: "It takes the name of a Pokemon as an argument. It should print the name, height, weight, stats and type(s) of the Pokemon",
			Callback:    inspectCallBack,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "list the names of all caught Pokemon.",
			Callback:    pokedexCallBack,
		},
	}
}

func pokedexCallBack(_ []string) error {
	if len(caughtPokemon) == 0 {
		fmt.Println("Your Pokedex is empty.")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, poke := range caughtPokemon {
		fmt.Println("- ", poke.Name)
	}
	return nil
}

func inspectCallBack(params []string) error {
	pokeName := params[0]
	pokemon, ok := caughtPokemon[pokeName]
	if !ok {
		fmt.Println("you have not caught ", pokeName)
		return nil
	}
	fmt.Println("Name: ", pokeName)
	fmt.Println("Base experience: ", pokemon.BaseExperience)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Println(fmt.Sprintf("- %s: %d", stat.Stat.Name, stat.BaseStat))
	}

	fmt.Println("Types:")
	for _, pkeType := range pokemon.Types {
		fmt.Println("- ", pkeType.Type.Name)

	}
	return nil
}

func catchCallBack(param []string) error {
	pokeUrl := pokemonUrl + param[0]
	fmt.Println(fmt.Sprintf("Throwing a Pokeball at %s ...", param[0]))
	pokemon, err := pokeApi.GetPokemon(pokeUrl)
	if err != nil {
		return err
	}
	if pokemon.BaseExperience < rand.Intn(pokemon.BaseExperience+rand.Intn(20)) {
		fmt.Println(fmt.Sprintf("%s was caught!", param[0]))
		caughtPokemon[param[0]] = pokemon
	} else {
		fmt.Println(fmt.Sprintf("%s escaped!", param[0]))
	}
	fmt.Println(pokemon.BaseExperience)
	return nil
}

func exploreCallBack(params []string) error {
	area := locationAreaURL + "/" + params[0]
	pokemon, err := pokeApi.GetListOfPokemon(area)
	if err != nil {
		return err
	}
	for _, pke := range pokemon.PokemonEncounters {
		fmt.Println(pke.Pokemon.Name)
	}
	return nil
}

func mapBackCallBack(_ []string) error {
	mapCommand := Commands["mapb"]
	if mapCommand.Config.Previous == "" {
		return fmt.Errorf("no previous page to show")
	}
	area, err := pokeApi.GetLocationArea(mapCommand.Config.Previous)
	if err != nil {
		return err
	}
	for _, result := range area.Results {
		fmt.Println(result.Name)
	}
	mapCommand.Config.Next = area.Next
	mapCommand.Config.Previous = area.Previous
	return nil
}

func mapCallBack(_ []string) error {
	mapCommand := Commands["map"]
	area, err := pokeApi.GetLocationArea(mapCommand.Config.Next)
	if err != nil {
		return err
	}
	for _, result := range area.Results {
		fmt.Println(result.Name)
	}
	mapCommand.Config.Next = area.Next
	mapCommand.Config.Previous = area.Previous
	return nil
}

func exitCallBack(_ []string) error {
	os.Exit(0)
	return nil
}

func helpCallback(_ []string) error {
	for _, command := range Commands {
		fmt.Println("name: " + command.Name)
		fmt.Println("description: " + command.Description)
	}
	return nil
}
