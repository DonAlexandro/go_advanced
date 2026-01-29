package advancedSelect

import (
	"fmt"
)

// producer generates numbers and sends them to a channel.
func producer(ch chan<- int, num int) {
	for i := range num {
		ch <- i
	}
	close(ch)
}

func read(ch1, ch2 chan int) chan int {
	result := make(chan int)

	// Fan-in: multiplexing multiple channels into one
	go func() {
		for {
			// Stop when both channels are set to nil
			if ch1 == nil && ch2 == nil {
				close(result)
				break
			}

			select {
			case num, ok := <-ch1:
				if !ok {
					ch1 = nil
					continue
				}
				result <- num
			case num, ok := <-ch2:
				if !ok {
					ch2 = nil
					continue
				}
				result <- num
			}
		}
	}()

	return result
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	// Start producers
	go producer(ch1, 5)
	go producer(ch2, 5)

	result := read(ch1, ch2)

	// Consume aggregated results
	for num := range result {
		fmt.Println("Received number:", num)
	}
}
