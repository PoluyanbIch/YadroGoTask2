package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func LoadStopWords(path string) (map[string]struct{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %s", err.Error())
		}
	}()

	stopwords := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			stopwords[word] = struct{}{}
		}
	}
	return stopwords, scanner.Err()
}
