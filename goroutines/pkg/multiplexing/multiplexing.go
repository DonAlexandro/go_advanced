package multiplexing

import (
	"fmt"
	"math/rand"
	"time"
)

// worker function simulates a worker that processes multiple jobs
func worker(id int, ch chan<- string) {
	// Process exactly 5 jobs per worker
	for i := range 5 {
		// Simulate variable processing time (0-999 milliseconds)
		// Random delay makes workers complete jobs at different times
		// This creates realistic concurrency where work takes unpredictable time
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		// Send formatted job completion message to worker's dedicated channel
		// Each worker has its own channel (ch1 for worker 1, ch2 for worker 2)
		// Message includes worker ID and job number for tracking
		ch <- fmt.Sprintf("Worker %d: Job %d", id, i)
	}
	// Worker exits after completing all 5 jobs
	// Note: channels are NOT closed, they remain open
}

func main() {
	// Create two unbuffered channels for worker communication
	// Each worker gets its own dedicated channel for sending results
	// Unbuffered channels provide synchronization between workers and main
	ch1 := make(chan string) // Channel for worker 1 messages
	ch2 := make(chan string) // Channel for worker 2 messages

	// Start two workers concurrently
	// Each worker processes 5 jobs independently
	// Workers run in parallel, not sequentially
	go worker(1, ch1) // Worker 1 sends to ch1
	go worker(2, ch2) // Worker 2 sends to ch2

	// Receive exactly 10 messages total (5 from each worker)
	// Loop runs 10 times to collect all job completion messages
	for range 10 {
		// select statement waits for message from EITHER channel
		// Executes whichever case becomes ready first
		// This creates fair, non-blocking message collection
		select {
		// Case 1: Receive message from worker 1
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)

		// Case 2: Receive message from worker 2
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		}
		// No default case - select blocks until one channel has data
		// This ensures we wait for actual worker messages
	}
}
