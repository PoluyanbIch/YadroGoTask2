package app

import (
	"log/slog"
	"regexp"
	"slices"
	"strings"

	"github.com/kljensen/snowball"
)

type Normalizer struct {
	stopwords map[string]struct{}
}

func NewNormalizer(stopwords map[string]struct{}) *Normalizer {
	return &Normalizer{stopwords: stopwords}
}

func (n *Normalizer) deleteStopWords(words []string) []string {
	words = slices.DeleteFunc(words, func(w string) bool {
		_, banned := n.stopwords[strings.ToLower(strings.TrimSpace(w))]
		return banned
	})
	return words
}

func norm(words []string) []string {
	var result []string
	for _, word := range words {
		newWord, err := snowball.Stem(word, "english", false)
		if err != nil {
			slog.Error(err.Error())
			return nil
		}
		result = append(result, newWord)
	}
	return result
}

func deleteDuplicates(words []string) []string {
	var result []string
	seen := make(map[string]struct{})
	for _, word := range words {
		if _, ok := seen[word]; ok {
			continue
		}
		result = append(result, word)
		seen[word] = struct{}{}
	}
	return result
}

func string2slice(s string) []string {
	re := regexp.MustCompile(`[^\w\s]`)
	s = re.ReplaceAllString(s, " ")
	words := strings.Fields(s)
	return words
}

func (n *Normalizer) Normalize(s string) []string {
	words := string2slice(s)
	words = n.deleteStopWords(words)
	words = norm(words)
	words = deleteDuplicates(words)
	return words
}
