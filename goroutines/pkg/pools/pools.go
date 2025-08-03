package pools

import (
	"fmt"
	"sync"
)

// worker function processes jobs from a channel and sends results to another channel
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	// defer ensures wg.Done() is called when this function exits
	// This decrements the WaitGroup counter, signaling work completion
	defer wg.Done()

	// Range over the jobs channel until it's closed
	// Each worker competes to receive jobs from the shared channel
	for job := range jobs {
		// Simulate work processing - print which worker is handling which job
		fmt.Printf("Worker %d processing job %d\n", id, job)

		// Process the job (double the value) and send result to results channel
		// Multiple workers may send to results channel concurrently
		results <- job * 2
	}
	// When jobs channel is closed and drained, this worker exits
	// defer wg.Done() executes here, notifying main that this worker finished
}

func main() {
	// Constants define the workload and worker pool size
	const numJobs = 5    // Total number of jobs to process
	const numWorkers = 3 // Number of concurrent workers

	// Create buffered channels with capacity equal to number of jobs
	// Buffered channels prevent blocking when all jobs/results are sent at once
	jobs := make(chan int, numJobs)    // Channel to distribute jobs to workers
	results := make(chan int, numJobs) // Channel to collect results from workers

	// WaitGroup tracks the number of active worker goroutines
	// Ensures main doesn't exit before all workers complete
	var wg sync.WaitGroup

	// Create and start the worker pool
	for w := 1; w <= numWorkers; w++ {
		// Increment WaitGroup counter before starting each goroutine
		// This tells WaitGroup to expect one more worker to complete
		wg.Add(1)

		// Start worker goroutine with unique ID and shared channels
		// All workers share the same jobs and results channels
		go worker(w, jobs, results, &wg)
	}

	// Distribute all jobs to the workers
	for j := 1; j <= numJobs; j++ {
		// Send job number to jobs channel
		// Workers compete to receive these jobs
		jobs <- j
	}

	// Close the jobs channel to signal no more jobs will be sent
	// This causes workers' range loops to exit when channel is drained
	close(jobs)

	// Wait for all workers to complete their tasks
	// Blocks until all workers call wg.Done() (counter reaches 0)
	wg.Wait()

	// Close results channel after all workers finished
	// Safe to close because no more results will be sent
	close(results)

	// Collect and print all results from the workers
	// Range over results channel until it's closed and drained
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
	// Program completes - all jobs processed and results collected
}
