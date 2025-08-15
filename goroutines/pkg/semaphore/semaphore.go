package semaphore

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	tokens chan struct{} // Buffered channel acts as token bucket - limits concurrent access
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		tokens: make(chan struct{}, n), // Create buffered channel with capacity n
	}
}

func (s *Semaphore) Acquire() {
	s.tokens <- struct{}{} // Send empty struct to claim a token - blocks if channel is full
}

func (s *Semaphore) Release() {
	<-s.tokens // Receive from channel to release a token - frees up space for others
}

func worker(id int, sem *Semaphore, wg *sync.WaitGroup) {
	defer wg.Done() // Signal completion when function exits

	sem.Acquire()                          // Block until semaphore token is available
	fmt.Printf("Worker %d starting\n", id) // Only 2 workers can reach this point simultaneously
	time.Sleep(1 * time.Second)            // Simulate work
	fmt.Printf("Worker %d done\n", id)
	sem.Release() // Return token, allowing another goroutine to proceed
}

func main() {
	sem := NewSemaphore(2) // Create semaphore allowing max 2 concurrent workers
	var wg sync.WaitGroup  // Coordinate waiting for all goroutines to complete

	for i := 1; i <= 10; i++ {
		wg.Add(1)              // Increment counter before starting goroutine
		go worker(i, sem, &wg) // Launch worker goroutine with pointer to WaitGroup
	}

	wg.Wait() // Block until all 10 workers complete
}
