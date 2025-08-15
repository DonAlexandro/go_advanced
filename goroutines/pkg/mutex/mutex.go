package mutex

import (
	"fmt"
	"sync"
)

// Counter struct provides thread-safe counting operations
type Counter struct {
	// mu (mutex) protects the value field from concurrent access
	// Mutex ensures only one goroutine can modify value at a time
	// Without mutex, concurrent increments would cause race conditions
	mu sync.Mutex

	// value holds the actual counter data
	// This field is protected by the mutex above
	// Direct access to this field would be unsafe in concurrent code
	value int
}

// Increment safely increases counter by 1 in concurrent environment
func (c *Counter) Increment() {
	// Lock the mutex before accessing shared state
	// This blocks other goroutines from entering critical section
	// Only one goroutine can hold the lock at a time
	c.mu.Lock()

	// defer ensures unlock happens when function exits
	// Critical: prevents deadlock if function panics or returns early
	// Unlock MUST happen even if increment operation fails
	defer c.mu.Unlock()

	// Perform the actual increment operation
	// This is the "critical section" - code that must be atomic
	// Only one goroutine can execute this line at a time due to mutex
	c.value++
}

// Value safely reads current counter value in concurrent environment
func (c *Counter) Value() int {
	// Lock mutex before reading shared state
	// Prevents reading partial/inconsistent data during concurrent writes
	// Even reads need protection in concurrent programs
	c.mu.Lock()

	// defer ensures unlock happens when function exits
	// Releases lock so other goroutines can access counter
	defer c.mu.Unlock()

	// Return current value safely
	// This read is atomic due to mutex protection
	// No other goroutine can modify value during this read
	return c.value
}

func main() {
	// WaitGroup synchronizes main with worker goroutines
	// Ensures main doesn't exit before all increments complete
	var wg sync.WaitGroup

	// Create counter instance with zero value
	// All goroutines will share this single counter instance
	// Without synchronization, this would cause race conditions
	counter := Counter{}

	// Start 1000 goroutines to increment counter concurrently
	for range 1000 {
		// Increment WaitGroup counter BEFORE starting goroutine
		// This prevents race where main calls wg.Wait() before goroutine calls wg.Add()
		wg.Add(1)

		// Start goroutine to perform concurrent increment
		go func() {
			// defer ensures wg.Done() executes when goroutine exits
			// This decrements WaitGroup counter, signaling completion
			defer wg.Done()

			// Call thread-safe increment method
			// Multiple goroutines call this simultaneously
			// Mutex ensures each increment is atomic and safe
			counter.Increment()
		}()
	}

	// Wait for all 1000 goroutines to complete their increments
	// Blocks until all goroutines call wg.Done() (counter reaches 0)
	// This ensures we read final value only after all work is done
	wg.Wait()

	// Read and print final counter value safely
	// At this point, all increments are guaranteed to be complete
	// Expected output: "Final counter value: 1000"
	fmt.Println("Final counter value:", counter.Value())

	// RACE CONDITION ANALYSIS:
	// Without mutex protection, multiple goroutines might:
	// 1. Read same value (e.g., 42) simultaneously
	// 2. Increment to 43 simultaneously
	// 3. Write 43 back simultaneously
	// Result: only 1 increment instead of 2 (lost update)
	//
	// With mutex protection:
	// 1. Goroutine A locks, reads 42, increments to 43, unlocks
	// 2. Goroutine B waits for lock, then reads 43, increments to 44, unlocks
	// Result: both increments are preserved (correct behavior)
	//
	// PERFORMANCE CONSIDERATIONS:
	// - Mutex adds overhead but ensures correctness
	// - Goroutines may wait (block) for lock availability
	// - Total time depends on contention level
	// - Alternative: atomic operations for simple counters
}
