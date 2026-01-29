package timeoutselect

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	// This goroutine simulates a long-running process (2 seconds)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "Hello"
	}()

	// select will pick whichever case is ready first
	select {
	case msg := <-ch:
		// This will execute if the message arrives within 1 second
		fmt.Println("Received message:", msg)
	case <-time.After(1 * time.Second):
		// This will execute if 1 second passes before ch receives a value
		fmt.Println("Timeout: Didn't receive message in time")
	}
}
