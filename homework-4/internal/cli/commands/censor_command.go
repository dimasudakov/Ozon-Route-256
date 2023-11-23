package commands

import (
	"fmt"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-4/internal/cli"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-4/internal/set"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-4/internal/utils"
	"strings"
	"sync"
	"unicode"
)

const (
	censorDescription = "Censors banned words in the input text."
)

type censorCommand struct {
	next        cli.FileFormatter
	banWords    set.Set[string]
	description string
}

func NewCensorCommand(banWords set.Set[string]) *censorCommand {
	return &censorCommand{
		banWords:    banWords,
		description: censorDescription,
	}
}

func (c *censorCommand) Execute(args ...string) {
	if len(args) != 2 {
		fmt.Println("Usage: <input_file_name> <output_file_name>")
		return
	}

	text, err := utils.ReadFileContent(args[0])
	if err != nil {
		fmt.Println("Error: ", err)
	}

	text = c.processText(text)

	if err := utils.WriteFileContent(text, args[1]); err != nil {
		fmt.Println("Error: ", err)
	}

	if c.next != nil {
		c.next.Execute(args[1], args[1])
	}
}

func (c *censorCommand) GetDescription() string {
	if c.next != nil {
		return c.description + " " + c.next.GetDescription()
	}
	return c.description
}

func (c *censorCommand) SetNext(next cli.FileFormatter) {
	c.next = next
}

func (c *censorCommand) processText(text *string) *string {
	var wg sync.WaitGroup
	paragraphs := strings.Split(*text, "\n\n")

	wg.Add(len(paragraphs))
	for i := range paragraphs {
		i := i

		go func() {
			defer wg.Done()
			words := strings.Fields(paragraphs[i])
			for j := range words {
				word := words[j]
				punctuation := ""
				for unicode.IsPunct(rune(word[len(word)-1])) {
					punctuation = string(word[len(word)-1]) + punctuation
					word = word[:len(word)-1]
				}

				if c.banWords.Contains(strings.ToLower(word)) {
					// Если слово является запрещенным, заменяем его на звездочки
					words[j] = strings.Repeat("*", len(word)) + punctuation
				}
			}
			paragraphs[i] = strings.Join(words, " ")
		}()

	}

	wg.Wait()

	*text = strings.Join(paragraphs, "\n\n")

	return text
}
