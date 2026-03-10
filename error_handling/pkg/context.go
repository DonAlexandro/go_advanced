package main

import (
	"context"
	"fmt"
	"time"
)

// worker listens to the context signal to stop execution
func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// The channel is closed when cancel() is called
			fmt.Println("Worker received cancellation signal and is stopping...")
			return
		default:
			// Simulate some background work
			fmt.Println("Worker is performing task...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func contextExample() {
	// Create a cancelable context based on the background context
	ctx, cancel := context.WithCancel(context.Background())

	// Ensuring resources are cleaned up
	defer cancel()

	// Start the worker in a separate goroutine
	go worker(ctx)

	// Let the worker run for 2 seconds
	time.Sleep(2 * time.Second)

	// Send the cancellation signal to the worker
	fmt.Println("Main: Sending cancellation signal...")
	cancel()

	// Give the worker a moment to print its exit message
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main: Program finished")
}
