package fanIn

import (
	"fmt"
	"sync"
)

// FAN-IN PATTERN: Distribute work from single source to multiple workers

// worker function processes jobs and sends results back
func worker(jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	// defer ensures wg.Done() executes when worker exits
	// This decrements WaitGroup counter, signaling this worker completed
	defer wg.Done()

	// Process jobs until channel is closed and drained
	// Multiple workers compete for jobs from the same input channel
	for job := range jobs {
		// Process job (double the value) and send result to output channel
		// All workers send to the same results channel (fan-in of results)
		results <- job * 2
	}
	// Worker exits when no more jobs available
	// defer wg.Done() executes here, notifying main of completion
}

func main() {
	// Configuration constants for the worker pool
	const numJobs = 5    // Total number of jobs to process
	const numWorkers = 3 // Number of parallel workers

	// Create buffered channels for job distribution and result collection
	// Buffer size = numJobs allows non-blocking send/receive operations
	jobs := make(chan int, numJobs)    // Input channel: distributes jobs to workers
	results := make(chan int, numJobs) // Output channel: collects results from workers

	// WaitGroup synchronizes main thread with worker goroutines
	// Ensures main doesn't exit before all workers complete
	var wg sync.WaitGroup

	// PHASE 1: Start worker pool (FAN-OUT begins)
	for w := 1; w <= numWorkers; w++ {
		// Increment counter BEFORE starting worker to avoid race condition
		// Tell WaitGroup to expect one more worker to complete
		wg.Add(1)

		// Start worker goroutine with shared channels
		// All workers read from same jobs channel (competitive consumption)
		// All workers write to same results channel (concurrent aggregation)
		go worker(jobs, results, &wg)
	}

	// PHASE 2: Distribute jobs to worker pool
	for j := 1; j <= numJobs; j++ {
		// Send job number to jobs channel
		// Workers compete to receive and process these jobs
		jobs <- j
	}

	// Signal end of job distribution
	// Close jobs channel to tell workers "no more jobs coming"
	// Workers will exit their range loops when channel is drained
	close(jobs)

	// PHASE 3: Wait for all workers to complete processing
	// Block until all workers call wg.Done() (counter reaches 0)
	// This ensures all jobs are processed before collecting results
	wg.Wait()

	// PHASE 4: Signal end of result collection
	// Safe to close results channel because all workers finished
	// No more results will be sent to the channel
	close(results)

	// PHASE 5: Collect and display all results (FAN-IN of results)
	// Range over results channel until it's closed and empty
	// Order of results is non-deterministic due to concurrent processing
	for result := range results {
		fmt.Println("Result:", result)
	}

	// EXECUTION SUMMARY:
	// 1. Jobs 1,2,3,4,5 distributed among 3 workers
	// 2. Workers process jobs concurrently (multiply by 2)
	// 3. Results 2,4,6,8,10 collected (order may vary)
	// 4. All goroutines cleaned up properly
}
