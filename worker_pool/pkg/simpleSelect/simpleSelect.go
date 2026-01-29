package selectPattern

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// This goroutine sends "Hello" after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "Hello"
	}()

	// This goroutine sends "World" after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- "World"
	}()

	// select waits for the first channel that is ready to communicate
	select {
	case msg1 := <-ch1:
		fmt.Println("Received message from ch1:", msg1)
	case msg2 := <-ch2:
		fmt.Println("Received message from ch2:", msg2)
	}
}
