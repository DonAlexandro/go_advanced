package fanOut

import (
	"fmt"
	"sync"
)

// worker function processes jobs from a shared channel
func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	// defer ensures wg.Done() executes when worker function exits
	// This decrements the WaitGroup counter, signaling work completion
	// Critical: must be called to prevent main goroutine from waiting forever
	defer wg.Done()

	// Range over jobs channel until it's closed and drained
	// Multiple workers compete to receive jobs from the same channel
	// Each job is processed by exactly one worker (load balancing)
	for job := range jobs {
		// Simulate job processing - print worker ID and job number
		// In real applications, this would be actual work (computation, I/O, etc.)
		fmt.Printf("Worker %d processing job %d\n", id, job)

		// NOTE: This simplified version doesn't collect results
		// Workers just process and print, no return values
	}
	// When jobs channel is closed and this worker has processed all available jobs,
	// the range loop exits and defer wg.Done() executes
}

func main() {
	// Define workload and worker pool size
	const numJobs = 5    // Total jobs to distribute among workers
	const numWorkers = 3 // Number of concurrent worker goroutines

	// Create buffered channel with capacity equal to number of jobs
	// Buffer size = numJobs allows sending all jobs without blocking
	// Workers can process jobs at their own pace
	jobs := make(chan int, numJobs)

	// WaitGroup synchronizes main with worker goroutines
	// Prevents main from exiting before all workers finish
	var wg sync.WaitGroup

	// Create and start worker pool
	for w := 1; w <= numWorkers; w++ {
		// Increment WaitGroup counter BEFORE starting goroutine
		// This prevents race condition where main calls wg.Wait() before worker calls wg.Add()
		wg.Add(1)

		// Start worker goroutine with unique ID and shared jobs channel
		// All workers receive from the same channel, creating natural load balancing
		go worker(w, jobs, &wg)
	}

	// Distribute all jobs to the worker pool
	for j := 1; j <= numJobs; j++ {
		// Send job number to jobs channel
		// Buffered channel allows all jobs to be queued immediately
		jobs <- j
	}

	// Close jobs channel to signal "no more jobs coming"
	// This causes workers' range loops to exit after processing remaining jobs
	// Critical: must close channel or workers will wait forever
	close(jobs)

	// Wait for all workers to complete their assigned jobs
	// Blocks until all workers call wg.Done() (counter reaches 0)
	// Ensures program doesn't exit while work is still being processed
	wg.Wait()

	// At this point, all jobs have been processed by the worker pool
	// Main goroutine can safely exit, knowing all work is complete
}
