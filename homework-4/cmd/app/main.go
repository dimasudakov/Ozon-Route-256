package main

import (
	"fmt"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-4/internal/cli"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-4/internal/cli/commands"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-4/internal/set"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-4/internal/utils"
	"os"
	"strings"
)

func main() {
	banWords := set.NewHashSet[string]()
	words, err := utils.ReadFileContent("banwords.txt")
	if err != nil {
		fmt.Println("Error occurred during reading banwords:", err)
		return
	}
	for _, banword := range strings.Fields(*words) {
		banWords.Add(strings.ToLower(banword))
	}

	c1 := commands.NewCensorCommand(banWords)
	c2 := commands.NewPointTabAddCommand()

	c1.SetNext(c2)

	commandHandler := cli.NewCommandHandler()

	commandHandler.RegisterCommand("help", commands.NewHelpCommand(commandHandler.GetCommands()))
	commandHandler.RegisterCommand("spell", commands.NewSpellCommand())
	commandHandler.RegisterCommand("format", c1)

	if len(os.Args) != 1 {
		commandHandler.ExecuteCommand(os.Args[1], os.Args[2:]...)
	}

}
