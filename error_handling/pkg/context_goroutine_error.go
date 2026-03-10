package main

import (
	"context"
	"fmt"
	"time"
)

// performTask performs a task and reports any errors or context cancellations
func performTask(ctx context.Context, ch chan<- error) {
	// Simulate work that takes 2 seconds
	// We use a timer to allow the worker to be interrupted immediately
	select {
	case <-time.After(2 * time.Second):
		// Work finished successfully
		ch <- nil
	case <-ctx.Done():
		// Context was cancelled or timed out before work finished
		ch <- ctx.Err()
	}
}

func contextGoroutineErrorExample() {
	// Create a context with a 1-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Channel to receive errors or success signals from the worker
	errCh := make(chan error)

	// Start the worker goroutine
	go performTask(ctx, errCh)

	// Wait for the worker to report back
	if err := <-errCh; err != nil {
		fmt.Printf("Worker failed: %v\n", err)
		return
	}

	fmt.Println("Work completed successfully")
}
