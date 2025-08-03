package unbufferedchannel

import (
	"fmt"
	"sync"
)

// «In this example, the send operation blocks until the main function is ready to receive the message.»

func main() {
	// WaitGroup is a counter that tracks running goroutines
	// It helps main goroutine wait for other goroutines to complete
	var wg sync.WaitGroup

	// Add(1) increments the counter by 1
	// This tells WaitGroup to expect 1 goroutine to complete
	wg.Add(1)

	// Create an unbuffered channel that can pass string messages
	// Unbuffered means sender blocks until receiver is ready
	messages := make(chan string)

	// Start a new goroutine (lightweight thread)
	go func() {
		// defer schedules wg.Done() to execute when this function exits
		// wg.Done() decrements the WaitGroup counter by 1
		// This signals that this goroutine has completed its work
		defer wg.Done()

		// Send a string message to the channel
		// This operation blocks until main goroutine receives from channel
		messages <- "Hello from Goroutine"
	}()

	// Receive a message from the channel
	// This operation blocks until the goroutine sends a message
	// Channel communication synchronizes the two goroutines
	msg := <-messages

	// Print the received message
	fmt.Println(msg)

	// Wait for all goroutines to call wg.Done()
	// Since we added 1 to the counter, this waits for 1 goroutine to finish
	// This ensures main doesn't exit before the goroutine completes
	wg.Wait()
}
