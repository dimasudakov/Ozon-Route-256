package cli

import "fmt"

type Command interface {
	Execute(args ...string)
	GetDescription() string
}

type FileFormatter interface {
	Command
	SetNext(FileFormatter)
}

type CommandHandler struct {
	commands map[string]Command
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		commands: make(map[string]Command),
	}
}

func (ch *CommandHandler) RegisterCommand(name string, command Command) {
	ch.commands[name] = command
}

func (ch *CommandHandler) ExecuteCommand(name string, args ...string) {
	command, found := ch.commands[name]
	if found {
		command.Execute(args...)
	} else {
		fmt.Println("Unknown command: ", name)
	}
}

func (ch *CommandHandler) GetCommands() map[string]Command {
	return ch.commands
}
