package timeout

import (
	"fmt"
	"time"
)

func main() {
	// Create an unbuffered channel for string communication
	// Unbuffered means sender blocks until receiver is ready
	ch := make(chan string)

	// Start a goroutine that simulates slow work
	go func() {
		// Simulate a long-running operation (e.g., network request, database query)
		// This goroutine sleeps for 2 seconds before sending result
		time.Sleep(2 * time.Second)

		// Send result to the channel after 2 seconds
		// If main is still waiting, this will unblock the select statement
		ch <- "Result"
	}()

	// select statement allows waiting on multiple channel operations
	// It executes the first case that becomes ready
	select {
	// Case 1: Try to receive from ch channel
	// This case will execute if goroutine sends data within timeout
	case res := <-ch:
		fmt.Println("Received:", res)

	// Case 2: Timeout case using time.After()
	// time.After(1 * time.Second) returns a channel that sends current time after 1 second
	// We ignore the time value using <-time.After() (no variable assignment)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout")
	}

	// Execution analysis:
	// - Goroutine takes 2 seconds to send result
	// - Timeout is set to 1 second
	// - Since 1 second < 2 seconds, timeout case will execute first
	// - Program will print "Timeout" and exit
	// - The goroutine continues running and sends to ch, but no one is listening
}
