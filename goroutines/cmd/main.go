package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	internal "github.com/DonAlexandro/go_advanced/internal"
	slogjson "github.com/veqryn/slog-json"
)

type Worker struct {
	jobs     <-chan string
	results  chan<- internal.FileWordFrequency
	errChan  chan<- error
	counters *int
}

func (w Worker) work() {
	for filePath := range w.jobs {
		// Count word frequencies in the file
		words, err := internal.CountWordFrequency(filePath, w.counters)

		if err != nil {
			w.errChan <- err
			continue
		}

		// Create result using the struct
		result := internal.FileWordFrequency{
			FileName: filepath.Base(filePath),
			Words:    words,
		}

		w.results <- result
	}
}

func main() {
	h := slogjson.NewHandler(os.Stdout, &slogjson.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelInfo,
		ReplaceAttr: nil, // Same signature and behavior as stdlib JSONHandler
	})
	// Default global logger
	slog.SetDefault(slog.New(h))

	// Command line flags
	workers := flag.Int("w", 4, "Number of workers to process files concurrently")
	counters := flag.Int("c", 2, "Number of goroutines counting the words in files")

	flag.Parse()

	// Get positional arguments (non-flag arguments)
	args := flag.Args()

	// Check if directory path is provided
	if len(args) == 0 {
		slog.Error("directory path is required")
		slog.Error(fmt.Sprintf("usage: %s <directory_path> [-w <num_workers>]\n", os.Args[0]))
		slog.Error(fmt.Sprintf("example: %s /path/to/files -w 8\n", os.Args[0]))
		os.Exit(1)
	}

	directoryPath := args[0]

	// Validate worker count
	if *workers < 1 {
		slog.Error(fmt.Sprintf("number of workers must be positive, got: %d", *workers))
		os.Exit(1)
	}

	// Validate counter goroutine count
	if *counters < 1 {
		slog.Error(fmt.Sprintf("number of counters must be positive, got: %d", *counters))
		os.Exit(1)
	}

	// Get all .txt files from the directory
	txtFiles, err := internal.GetTxtFiles(directoryPath)
	if err != nil {
		slog.Error("error reading directory", slog.Any("error", err))
		os.Exit(1)
	}

	jobsNum := len(txtFiles)

	jobs := make(chan string, jobsNum)
	results := make(chan internal.FileWordFrequency, jobsNum)
	errChan := make(chan error, jobsNum)

	var wg sync.WaitGroup

	// Create and start the worker pool
	for w := 1; w <= *workers; w++ {
		wg.Go(func() {
			worker := Worker{
				jobs:     jobs,
				results:  results,
				errChan:  errChan,
				counters: counters,
			}

			worker.work()
		})
	}

	// Send all file paths to jobs channel
	for _, filePath := range txtFiles {
		jobs <- filePath
	}

	// Close the jobs to indicate no more filePaths will be provided after the loop
	close(jobs)

	// Wait for all workers to complete their tasks
	wg.Wait()

	// Close results and error channels after all workers finished
	close(results)
	close(errChan)

	// Collect and print results
	for result := range results {
		fmt.Print(result.ToHumanReadable())
	}

	// Check for any errors
	for err := range errChan {
		if err != nil {
			slog.Error("error processing file", slog.Any("error", err))
		}
	}
}
