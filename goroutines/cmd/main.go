package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	internal "github.com/DonAlexandro/go_advanced/internal"
)

func worker(jobs <-chan string, results chan<- internal.FileWordFrequency, errChan chan<- error) {
	for filePath := range jobs {
		// Count word frequencies in the file
		words, err := internal.CountWordFrequency(filePath)

		if err != nil {
			errChan <- err
			continue
		}

		// Create result using the struct
		result := internal.FileWordFrequency{
			FileName: filepath.Base(filePath),
			Words:    words,
		}

		results <- result
	}
}

func main() {
	// Command line flags
	workers := flag.Int("w", 4, "Number of workers")
	flag.Parse()

	// Get positional arguments (non-flag arguments)
	args := flag.Args()

	// Check if directory path is provided
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: directory path is required\n")
		fmt.Fprintf(os.Stderr, "Usage: %s <directory_path> [-w <num_workers>]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s /path/to/files -w 8\n", os.Args[0])
		os.Exit(1)
	}

	directoryPath := args[0]

	// Validate worker count
	if *workers < 1 {
		fmt.Fprintf(os.Stderr, "Error: number of workers must be positive, got: %d\n", *workers)
		os.Exit(1)
	}

	// Get all .txt files from the directory
	txtFiles, err := internal.GetTxtFiles(directoryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
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
			worker(jobs, results, errChan)
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
			fmt.Fprintf(os.Stderr, "Error processing file: %v\n", err)
		}
	}
}
