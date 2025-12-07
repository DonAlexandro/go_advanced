package internal

import (
	"strings"
	"unicode"
)

// TextPreprocessor handles text preprocessing using a pipeline pattern
type TextPreprocessor struct{}

// ToLower creates a pipeline stage that converts text to lowercase
func (tp *TextPreprocessor) ToLower(in <-chan string) <-chan string {
	// Create unbuffered output channel
	out := make(chan string)

	// Start goroutine to process input stream asynchronously
	go func() {
		// defer ensures output channel is properly closed
		// Critical for pipeline termination propagation
		defer close(out)

		// Process each string from input channel until it's closed
		// Range automatically handles channel closing - exits when in is closed and drained
		for text := range in {
			// Convert to lowercase and send to output
			out <- strings.ToLower(text)
		}
		// When input channel closes and all text processed,
		// defer close(out) signals next stage that transformation is complete
	}()

	// Return receive-only channel immediately (non-blocking)
	// Allows caller to start consuming transformed data
	return out
}

// RemovePunctuation creates a pipeline stage that removes non-alphanumeric characters
func (tp *TextPreprocessor) RemovePunctuation(in <-chan string) <-chan string {
	// Create unbuffered output channel
	out := make(chan string)

	// Start goroutine to process input stream asynchronously
	go func() {
		// defer ensures output channel is properly closed
		defer close(out)

		// Process each string from input channel until it's closed
		for text := range in {
			// Remove all non-letter and non-digit characters except spaces
			// This preserves word boundaries for the next stage
			cleaned := strings.Map(func(r rune) rune {
				if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
					return r
				}
				// Replace punctuation with space to maintain word boundaries
				return ' '
			}, text)

			// Send cleaned text to output
			out <- cleaned
		}
		// When input channel closes, defer close(out) signals next stage
	}()

	// Return receive-only channel immediately (non-blocking)
	return out
}

// SplitIntoWords creates a pipeline stage that splits text into individual words
func (tp *TextPreprocessor) SplitIntoWords(in <-chan string) <-chan string {
	// Create unbuffered output channel
	out := make(chan string)

	// Start goroutine to process input stream asynchronously
	go func() {
		// defer ensures output channel is properly closed
		defer close(out)

		// Process each string from input channel until it's closed
		for text := range in {
			// Split text by whitespace into individual words
			words := strings.FieldsSeq(text)

			// Emit each non-empty word to output channel
			for word := range words {
				if len(word) > 0 { // Filter out empty strings
					out <- word
				}
			}
		}
		// When input channel closes and all words emitted,
		// defer close(out) signals consumer that no more words coming
	}()

	// Return receive-only channel immediately (non-blocking)
	return out
}

// PreprocessText orchestrates the 3-stage pipeline for text preprocessing
func (tp *TextPreprocessor) PreprocessText(text string) <-chan string {
	// Create unbuffered initial channel to feed raw text
	input := make(chan string)

	// Start goroutine to feed the entire text chunk into the pipeline
	go func() {
		// defer ensures channel is closed after sending
		defer close(input)

		// Send the entire text chunk as a single value
		input <- text
		// After sending, defer close(input) executes
		// This signals ToLower stage that no more text is coming
	}()

	// PIPELINE CONSTRUCTION: Chain processing stages together
	// Stage 1: Convert to lowercase
	lowercased := tp.ToLower(input)

	// Stage 2: Remove punctuation from lowercased text
	cleaned := tp.RemovePunctuation(lowercased)

	// Stage 3: Split cleaned text into individual words
	words := tp.SplitIntoWords(cleaned)

	// Return final word stream channel
	// Consumer will receive individual cleaned words one by one
	return words
}
