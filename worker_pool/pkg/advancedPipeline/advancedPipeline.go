package advancedPipeline

import (
	"fmt"
	"math/rand"
)

// Source generates a stream of random numbers and sends them to a channel.
func Source(count int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for range count {
			out <- rand.Intn(100)
		}
	}()
	return out
}

// Filter filters out numbers from the input channel that are less than the threshold.
func Filter(in <-chan int, threshold int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			if num >= threshold {
				out <- num
			}
		}
	}()
	return out
}

// Aggregator calculates the sum of all numbers received from the input channel.
func Aggregator(in <-chan int) int {
	sum := 0
	for num := range in {
		sum += num
	}
	return sum
}

func main() {
	// Create pipeline stages
	numbers := Source(10)
	filteredNumbers := Filter(numbers, 50)

	// Execute aggregation
	sum := Aggregator(filteredNumbers)

	// Print the result
	fmt.Println("Sum of filtered numbers:", sum)
}
