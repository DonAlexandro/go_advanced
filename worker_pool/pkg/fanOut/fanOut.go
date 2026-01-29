package fanout

import (
	"fmt"
	"sync"
)

// producer generates numbers and sends them to a channel.
func producer(ch chan<- int, num int) {
	for i := range num {
		ch <- i
	}
	close(ch)
}

// worker processes numbers received from the input channel.
func worker(id int, in <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range in {
		fmt.Printf("Worker %d received %d\n", id, num)
	}
}

// fanOut starts multiple workers to process data from a single channel.
func fanOut(ch <-chan int, numWorkers int) {
	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, ch, &wg)
	}
	wg.Wait()
}

func main() {
	ch := make(chan int)

	// Start the producer in a separate goroutine
	go producer(ch, 10)

	// fanOut will block until all workers are done
	fanOut(ch, 3)

	fmt.Println("All workers finished processing")
}
