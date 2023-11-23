package commands

import "fmt"

const (
	SpellDescription = "Prints all letters of input word to the console separated by a space."
)

type spellCommand struct {
	description string
}

func NewSpellCommand() *spellCommand {
	return &spellCommand{
		description: SpellDescription,
	}
}

func (c *spellCommand) Execute(args ...string) {
	if len(args) == 0 {
		fmt.Println("Usage: spell <words>")
		return
	}

	word := args[0]
	for _, letter := range word {
		fmt.Print(string(letter), " ")
	}
}

func (c *spellCommand) GetDescription() string {
	return c.description
}
