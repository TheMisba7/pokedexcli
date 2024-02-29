package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/api"
)

const locationAreaURL = "https://pokeapi.co/api/v2/location-area"

var pokeapi = api.PokeApiImpl{}
var commands map[string]command

type command struct {
	name        string
	description string
	callback    func() error
	config      *config
}

type config struct {
	next     string
	previous string
}

func initCommands() map[string]command {
	config := &config{next: locationAreaURL}
	return map[string]command{
		"help": {
			name:        "help",
			description: "this command will help you to learn about pokedex",
			callback:    helpCallback,
		},
		"exit": {
			name:        "exit",
			description: "exist the pokedex program",
			callback:    exitCallBack,
		},
		"map": {
			name:        "map",
			description: "get next 20 location areas",
			callback:    mapCallBack,
			config:      config,
		},
		"mapb": {
			name:        "mapb",
			description: "map back to previous page",
			callback:    mapBackCallBack,
			config:      config,
		},
	}
}

func mapBackCallBack() error {
	mapCommand := commands["mapb"]
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
	mapCommand := commands["map"]
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
	for _, command := range commands {
		fmt.Println("name: " + command.name)
		fmt.Println("description: " + command.description)
	}
	return nil
}

func main() {
	commands = initCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("pokedex> ")
		scanner.Scan()
		input := scanner.Text()
		if "exist" == input {
			break
		}
		command, ok := commands[input]
		if ok {
			err := command.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("unkown command")
		}
	}
}
