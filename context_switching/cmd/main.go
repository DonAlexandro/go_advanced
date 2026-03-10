package main

import (
	"fmt"
	"sync"
	"sync/atomic" // Import the atomic package for low-level synchronization
)

// Use int64 for compatibility with atomic functions
var atomicCounter int64

func main() {
	var wg sync.WaitGroup

	for range 1000 {
		wg.Go(func() {
			// Atomic addition is performed in one CPU cycle without blocking
			atomic.AddInt64(&atomicCounter, 1)
		})
	}

	wg.Wait()

	// Atomically load the final value to ensure consistency
	finalValue := atomic.LoadInt64(&atomicCounter)
	fmt.Printf("Final Atomic Counter Value: %d\n", finalValue)
}
