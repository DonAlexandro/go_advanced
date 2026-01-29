package simplePipeline

import "fmt"

// Generator generates a sequence of integers and sends them to a channel.
func Generator(count int) <-chan int {
	out := make(chan int, count)
	go func() {
		defer close(out)
		for i := 1; i <= count; i++ {
			out <- i
		}
	}()
	return out
}

// Doubler doubles each integer received from the input channel.
func Doubler(in <-chan int) <-chan int {
	out := make(chan int, cap(in))
	go func() {
		defer close(out)
		for num := range in {
			out <- num * 2
		}
	}()
	return out
}

// Printer prints each integer received from the input channel.
func Printer(in <-chan int) {
	for num := range in {
		fmt.Println(num)
	}
}

func main() {
	// Create pipeline stages
	numbers := Generator(5)
	doubledNumbers := Doubler(numbers)

	// Start the pipeline
	Printer(doubledNumbers)
}
