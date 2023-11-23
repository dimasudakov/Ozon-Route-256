package utils

import (
	"bufio"
	"os"
)

func ReadFileContent(filename string) (*string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		content += line + "\n"
	}

	if err2 := scanner.Err(); err2 != nil {
		return nil, err2
	}

	return &content, nil
}

func WriteFileContent(content *string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(*content)
	if err != nil {
		return err
	}

	writer.Flush()

	return nil
}
