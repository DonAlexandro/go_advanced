package errorhandling

import (
	"errors"
	"fmt"
)

// Handling errors in goroutines requires careful management since errors cannot be returned directly from a goroutine to its caller.
// One approach is to use channels for error reporting.
func workerWithError(jobs <-chan int, results chan<- int, errChan chan<- error) {
	for job := range jobs {
		if job == 2 {
			errChan <- errors.New("job 2 failed")
			continue
		}
		results <- job * 2
	}
}

func main() {
	jobs := make(chan int, 5)
	results := make(chan int, 5)
	errChan := make(chan error)

	for w := 1; w <= 3; w++ {
		go workerWithError(jobs, results, errChan)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	for range 5 {
		select {
		case result := <-results:
			fmt.Println("Result:", result)
		case err := <-errChan:
			fmt.Println("Error:", err)
		}
	}
}
