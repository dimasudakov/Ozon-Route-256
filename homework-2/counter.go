package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Counter interface {
	Count(n int, textFilePath, keywordFilePath string) (map[string]int, error)
}

type KeywordCounter struct {
	result map[string]int
	mutex  sync.Mutex
}

func NewKeywordCounter() *KeywordCounter {
	return &KeywordCounter{result: make(map[string]int)}
}

func (kc *KeywordCounter) Count(n int, textFilePath, keywordFilePath string) (map[string]int, error) {
	file, err := os.Open(keywordFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	keywords, err := kc.readKeywords(scanner)
	if err != nil {
		return nil, err
	}

	lines := kc.readText(textFilePath)

	kc.countKeywords(n, lines, keywords)

	return kc.result, nil
}

func (kc *KeywordCounter) readKeywords(scanner *bufio.Scanner) ([]string, error) {
	scanner.Scan()
	cntStr := scanner.Text()

	n, err := strconv.Atoi(cntStr)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		result = append(result, scanner.Text())
	}

	return result, nil
}

func (kc *KeywordCounter) readText(fileName string) <-chan string {
	out := make(chan string)

	file, _ := os.Open(fileName)

	scanner := bufio.NewScanner(file)
	go func() {
		for scanner.Scan() {
			out <- scanner.Text()
		}
		close(out)
	}()

	return out
}

func (kc *KeywordCounter) countKeywords(n int, in <-chan string, keywords []string) {
	var wg sync.WaitGroup

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for line := range in {
				line = strings.ToLower(line)
				for _, keyword := range keywords {
					keywordCnt := strings.Count(line, keyword)
					kc.mutex.Lock()
					kc.result[keyword] += keywordCnt
					kc.result["Всего"] += keywordCnt
					kc.mutex.Unlock()
				}
			}
		}()
	}

	wg.Wait()
}
