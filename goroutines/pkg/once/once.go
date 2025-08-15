package once

import (
	"fmt"
	"sync"
)

func main() {
	// sync.Once ensures a function is executed only once across multiple goroutines
	// Even if called concurrently from different goroutines, the wrapped function runs exactly once
	var once sync.Once

	// Anonymous function that we want to execute only once
	initFunc := func() {
		fmt.Println("Initialization done")
	}

	// WaitGroup coordinates waiting for multiple goroutines to complete
	// Acts as a counter: increment before starting goroutines, decrement when they finish
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1) // Increment counter - expecting one more goroutine to complete

		go func() { // Launch goroutine (concurrent execution)
			defer wg.Done()   // Decrement counter when this goroutine exits
			once.Do(initFunc) // Execute initFunc, but only once across all goroutines
		}()
	}

	wg.Wait() // Block until all goroutines call wg.Done() (counter reaches zero)
}
