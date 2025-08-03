package bufferedchannel

import (
	"fmt"
)

func main() {
	// Create a buffered channel with capacity of 2 string messages
	// The number 2 is the buffer size - channel can hold 2 values without blocking
	// Buffered channels allow asynchronous communication between goroutines
	messages := make(chan string, 2)

	// Send first message to the channel
	// This does NOT block because buffer has space (0/2 -> 1/2)
	// With unbuffered channel, this would block until someone receives
	messages <- "Buffered"

	// Send second message to the channel
	// This also does NOT block because buffer still has space (1/2 -> 2/2)
	// Now the buffer is full - a third send would block
	messages <- "Channel"

	// Receive and print the first message
	// This gets "Buffered" (FIFO - First In, First Out)
	// Buffer state changes from full to (1/2)
	fmt.Println(<-messages)

	// Receive and print the second message
	// This gets "Channel"
	// Buffer state changes from (1/2) to empty (0/2)
	fmt.Println(<-messages)

	// At this point, the channel buffer is empty
	// Any additional receive operation would block until someone sends
}
