package context

import (
	"context"
	"fmt"
	"time"
)

// «The `context` package in Go provides a way to signal cancellation across multiple goroutines, which is essential for robust and responsive applications.»

// workerWithContext simulates a worker that can be cancelled via context
func workerWithContext(ctx context.Context, id int) {
	// Infinite loop - worker keeps running until told to stop
	for {
		// select statement allows non-blocking channel operations
		select {
		// Case 1: Check if context was cancelled/timed out
		case <-ctx.Done():
			// ctx.Done() returns a channel that closes when context expires
			fmt.Printf("Worker %d stopping\n", id)
			return // Exit the goroutine completely
		// Case 2: Default case executes when no other case is ready
		default:
			// Simulate work being done
			fmt.Printf("Worker %d working\n", id)
			// Sleep for 500ms to simulate processing time
			time.Sleep(500 * time.Millisecond)
			// After sleep, loop continues and checks ctx.Done() again
		}
	}
}

func main() {
	// Create context that automatically cancels after 2 seconds
	// ctx will signal cancellation, cancel() function allows manual cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	// defer ensures cancel() is called when main() exits
	// This prevents context leak and releases resources
	defer cancel()

	// Start 3 worker goroutines concurrently
	for i := 1; i <= 3; i++ {
		// Each goroutine runs workerWithContext with the same context
		// All workers will stop when ctx times out after 2 seconds
		go workerWithContext(ctx, i)
	}

	// Main goroutine sleeps for 3 seconds
	// This is longer than the 2-second context timeout
	// Workers will stop at 2 seconds, but main continues until 3 seconds
	time.Sleep(3 * time.Second)

	// This prints after workers have already stopped
	fmt.Println("Main function completed")
}
