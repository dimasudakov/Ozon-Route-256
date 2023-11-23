package commands

import (
	"fmt"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-4/internal/cli"
)

const (
	helpDescription = "Displays brief information about commands."
)

type helpCommand struct {
	commandMap  map[string]cli.Command
	description string
}

func NewHelpCommand(commandMap map[string]cli.Command) *helpCommand {
	return &helpCommand{
		commandMap:  commandMap,
		description: helpDescription,
	}
}

func (hc *helpCommand) Execute(args ...string) {
	if len(args) != 0 {
		if len(args) > 1 {
			fmt.Println("Usage: help <command name>")
			return
		}
		if _, found := hc.commandMap[args[0]]; !found {
			fmt.Println("Unknown command:", args[0])
			return
		}
		fmt.Printf("\t%-15s\t%s\n", args[0], hc.commandMap[args[0]].GetDescription())
		return
	}

	fmt.Println("Available commands:")
	for cmd := range hc.commandMap {
		fmt.Printf("\t%-15s\t%s\n", cmd, hc.commandMap[cmd].GetDescription())
	}
}

func (hc *helpCommand) GetDescription() string {
	return hc.description
}
