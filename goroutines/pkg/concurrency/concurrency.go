package concurrency

// This program demonstrates how to use goroutines and channels to process jobs concurrently.
import (
	"fmt"
	"sync"
)

func worker(id int /* Receive-only channel */, jobs <-chan int /* Send-only channel */, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done() // Notify the WaitGroup that this goroutine is done

	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)

		results <- job * 2 // Simulate some work by doubling the job value
	}
}

func main() {
	const numJobs = 5
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)

		go worker(w, jobs, results, &wg)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	close(jobs)

	wg.Wait()

	close(results)

	for result := range results {
		fmt.Println("Result:", result)
	}
}
