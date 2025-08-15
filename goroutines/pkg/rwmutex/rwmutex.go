package rwmutex

import (
	"fmt"
	"sync"
)

// SafeMap is like a JavaScript Map but thread-safe for concurrent access
type SafeMap struct {
	mu    sync.RWMutex      // Read-Write lock (like a shared/exclusive lock)
	store map[string]string // The actual map data
}

// NewSafeMap creates a new thread-safe map (like `new Map()` in JS)
func NewSafeMap() *SafeMap {
	return &SafeMap{
		store: make(map[string]string), // Initialize empty map
	}
}

// Set adds/updates a key-value pair (like `map.set(key, value)` in JS)
func (m *SafeMap) Set(key, value string) {
	m.mu.Lock()         // Exclusive lock - no other reads/writes allowed
	defer m.mu.Unlock() // Always unlock when function exits
	m.store[key] = value
}

// Get retrieves a value by key (like `map.get(key)` in JS)
func (m *SafeMap) Get(key string) (string, bool) {
	m.mu.RLock()              // Shared lock - multiple reads OK, no writes
	defer m.mu.RUnlock()      // Always unlock when function exits
	value, ok := m.store[key] // ok=true if key exists, false if not
	return value, ok
}

func main() {
	var wg sync.WaitGroup // Like Promise.all() - wait for all goroutines to finish
	smap := NewSafeMap()

	// Writers - 10 goroutines writing data (like 10 async functions)
	for i := range 10 {
		wg.Add(1)        // Tell WaitGroup to expect 1 more goroutine
		go func(i int) { // Start goroutine (like async function)
			defer wg.Done()                                              // Signal completion when done
			smap.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i)) // Write to map
		}(i) // Pass i as parameter to avoid closure issues
	}

	// Readers - 10 goroutines reading data
	for i := range 10 {
		wg.Add(1)        // Tell WaitGroup to expect 1 more goroutine
		go func(i int) { // Start another goroutine
			defer wg.Done()                                         // Signal completion when done
			if value, ok := smap.Get(fmt.Sprintf("key%d", i)); ok { // Try to read from map
				fmt.Println("Got:", value) // Print if key was found
			}
		}(i) // Pass i as parameter
	}

	wg.Wait() // Wait for all 20 goroutines to finish (like await Promise.all())
}
