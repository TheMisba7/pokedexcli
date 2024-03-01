package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/command"
	"strings"
)

func main() {
	command.Commands = command.InitCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("pokedex> ")
		scanner.Scan()
		input := scanner.Text()
		cmdParams := strings.Fields(input)
		if len(cmdParams) == 0 {
			continue
		}
		cliCmd, ok := command.Commands[cmdParams[0]]
		if ok {
			err := cliCmd.Callback(cmdParams[1:])
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("unknown command")
		}
	}
}
