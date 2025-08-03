package pkg

// This program demonstrates how to use channels to communicate between goroutines.
func sum(a, b int, resultChan chan int) int {
	sum := a + b

	resultChan <- sum // Send the result to the channel

	return sum
}

func main() {
	resultChan := make(chan int)

	// Start the goroutine to calculate the sum
	go sum(3, 4, resultChan)

	// Wait for the result from the channel
	result := <-resultChan

	println("The sum is:", result) // Print the result
}
