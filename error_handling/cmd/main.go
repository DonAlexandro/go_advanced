package main

import (
	"fmt"
	"sync"
)

// worker sends its result (error or nil) to the channel
func worker(ch chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	// Simulate work that might fail
	err := fmt.Errorf("error from worker")

	// Send error to the channel
	ch <- err
}

func main() {
	numWorkers := 5
	// Use a buffered channel to prevent goroutines from blocking
	// if main returns early after the first error
	errCh := make(chan error, numWorkers)
	var wg sync.WaitGroup

	for range make([]struct{}, numWorkers) {
		wg.Add(1)
		go worker(errCh, &wg)
	}

	// Wait for all workers and close the channel in a separate goroutine
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Flag to track if we encountered any errors
	hasError := false

	// Iterate through results as they arrive
	for err := range errCh {
		if err != nil {
			fmt.Println("Caught Error:", err)
			hasError = true
			// If you want to stop on the first error, you can 'break' here.
			// Thanks to the buffer, other workers won't leak.
			break
		}
	}

	if !hasError {
		fmt.Println("All work completed successfully")
	} else {
		fmt.Println("Processing finished with errors")
	}
}
