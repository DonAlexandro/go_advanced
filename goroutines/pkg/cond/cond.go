package cond

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex         // Regular mutex for protecting shared data
	cond := sync.NewCond(&mu) // Condition variable - like event emitter but for goroutines
	ready := false            // Shared flag (like a global variable in Node.js)

	// Start a worker goroutine (like an async function)
	go func() {
		cond.L.Lock() // Lock the mutex (cond.L is the underlying mutex)
		for !ready {  // Keep checking if we're ready (like polling)
			cond.Wait() // Sleep until someone calls Signal() - releases lock temporarily
		}
		fmt.Println("Goroutine proceeding")
		cond.L.Unlock() // Release the lock
	}()

	time.Sleep(1 * time.Second) // Simulate work (like setTimeout in Node.js)

	// Tell the waiting goroutine to wake up
	cond.L.Lock()   // Must lock before modifying shared state
	ready = true    // Set the flag
	cond.Signal()   // Wake up ONE waiting goroutine (like emit('ready') in Node.js)
	cond.L.Unlock() // Release the lock

	time.Sleep(1 * time.Second) // Give goroutine time to print before main exits
}
