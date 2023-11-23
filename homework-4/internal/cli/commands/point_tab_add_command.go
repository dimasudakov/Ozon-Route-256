package commands

import (
	"fmt"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-4/internal/cli"
	"gitlab.ozon.dev/sudakov.dima.2014/homework-4/internal/utils"
	"strings"
	"sync"
	"unicode"
)

const (
	pointTabAddDescription = "Adds dots at the end of sentences and tabs at the beginning of paragraphs."
)

type pointTabAddCommand struct {
	next        cli.FileFormatter
	description string
}

func NewPointTabAddCommand() *pointTabAddCommand {
	return &pointTabAddCommand{
		description: pointTabAddDescription,
	}
}

func (f *pointTabAddCommand) Execute(args ...string) {
	if len(args) != 2 {
		fmt.Println("Usage: <input_file_name> <output_file_name>")
		return
	}
	text, err := utils.ReadFileContent(args[0])
	if err != nil {
		fmt.Println("Error: ", err)
	}

	text = f.processText(text)

	if err := utils.WriteFileContent(text, args[1]); err != nil {
		fmt.Println("Error: ", err)
	}

	if f.next != nil {
		f.next.Execute(args[1], args[1])
	}
}

func (f *pointTabAddCommand) SetNext(next cli.FileFormatter) {
	f.next = next
}

func (f *pointTabAddCommand) GetDescription() string {
	if f.next != nil {
		return f.description + " " + f.next.GetDescription()
	}
	return f.description
}

func (f *pointTabAddCommand) processText(text *string) *string {
	var wg sync.WaitGroup
	paragraphs := strings.Split(*text, "\n\n")

	wg.Add(len(paragraphs))
	for i := range paragraphs {
		i := i

		go func() {
			defer wg.Done()
			words := strings.Fields(paragraphs[i])
			for j := range words {
				if j == 0 {
					words[j] = "\t" + words[j]
				}
				// Если слово еще не заканчивается на знак пунктуации, а следующее за ним слово с заглавной буквы, то ставим точку
				if !unicode.IsPunct(rune(words[j][len(words[j])-1])) && (j == len(words)-1 || unicode.IsUpper(rune(words[j+1][0]))) {
					words[j] += "."
				}
			}
			paragraphs[i] = strings.Join(words, " ")
		}()

	}

	wg.Wait()
	*text = strings.Join(paragraphs, "\n\n")

	return text
}
