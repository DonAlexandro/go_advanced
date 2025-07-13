package main

// This is a simple Go program that demonstrates the use of goroutines.
// func sayHello() {
// 	// This function prints a greeting message.
// 	println("Hello, World!")
// }

// func main() {
// 	go sayHello()

// 	time.Sleep(1 * time.Second) // Wait for the goroutine to finish
// }

// This program demonstrates how to use channels to communicate between goroutines.
// func sum(a, b int, resultChan chan int) int {
// 	sum := a + b

// 	resultChan <- sum // Send the result to the channel

// 	return sum
// }

// func main() {
// 	resultChan := make(chan int)

// 	// Start the goroutine to calculate the sum
// 	go sum(3, 4, resultChan)

// 	// Wait for the result from the channel
// 	result := <-resultChan

// 	println("The sum is:", result) // Print the result
// }

// This program demonstrates how to use goroutines and channels to process jobs concurrently.
// import (
// 	"fmt"
// 	"sync"
// )

// func worker(id int /* Receive-only channel */, jobs <-chan int /* Send-only channel */, results chan<- int, wg *sync.WaitGroup) {
// 	defer wg.Done() // Notify the WaitGroup that this goroutine is done

// 	for job := range jobs {
// 		fmt.Printf("Worker %d processing job %d\n", id, job)

// 		results <- job * 2 // Simulate some work by doubling the job value
// 	}
// }

// func main() {
// 	const numJobs = 5
// 	const numWorkers = 3

// 	jobs := make(chan int, numJobs)
// 	results := make(chan int, numJobs)

// 	var wg sync.WaitGroup

// 	for w := 1; w <= numWorkers; w++ {
// 		wg.Add(1)

// 		go worker(w, jobs, results, &wg)
// 	}

// 	for j := 1; j <= numJobs; j++ {
// 		jobs <- j
// 	}

// 	close(jobs)

// 	wg.Wait()

// 	close(results)

// 	for result := range results {
// 		fmt.Println("Result:", result)
// 	}
// }
