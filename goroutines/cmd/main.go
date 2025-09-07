package main

import (
	"flag"
	"fmt"
	"os"
)

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

	// Print the values
	fmt.Printf("Directory path: %s\n", directoryPath)
	fmt.Printf("Number of workers: %d\n", *workers)
}
