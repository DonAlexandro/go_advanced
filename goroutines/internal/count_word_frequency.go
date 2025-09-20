package internal

import (
	"os"
	"strings"

	"github.com/mdobak/go-xerrors"
)

type Word struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// CountWordFrequency reads a file and counts the frequency of each word
func CountWordFrequency(filePath string) ([]Word, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, xerrors.Newf("failed to read a file %q: %w", filePath, err)
	}

	// Convert to lowercase and split into words
	text := strings.ToLower(string(content))
	// Split by whitespace and remove punctuation
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'))
	})

	// Count word frequencies using a map
	frequency := make(map[string]int)
	for _, word := range words {
		if len(word) > 0 { // Skip empty strings
			frequency[word]++
		}
	}

	// Convert map to slice of Word structs
	result := make([]Word, 0, len(frequency))
	for word, count := range frequency {
		result = append(result, Word{
			Word:  word,
			Count: count,
		})
	}

	return result, nil
}
