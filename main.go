package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/command"
)

func main() {
	command.Commands = command.InitCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("pokedex> ")
		scanner.Scan()
		input := scanner.Text()
		if "exit" == input {
			break
		}
		cliCmd, ok := command.Commands[input]
		if ok {
			err := cliCmd.Callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("unkown command")
		}
	}
}
