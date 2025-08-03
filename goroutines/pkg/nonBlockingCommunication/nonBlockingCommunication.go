package nonBlockingCommunication

import (
	"fmt"
)

// KEY INSIGHT: default cases make channel operations non-blocking
// Without default, each select would block indefinitely waiting for channel activity

func main() {
	// Create two unbuffered channels
	// Unbuffered channels block sender until receiver is ready (synchronous)
	messages := make(chan string) // Channel for string messages
	signals := make(chan bool)    // Channel for boolean signals

	// SELECT BLOCK 1: Non-blocking receive attempt
	select {
	// Try to receive a message from messages channel
	// This would block if there was no default case
	case msg := <-messages:
		fmt.Println("Received message:", msg)
	// default case executes immediately if no other case is ready
	// Since no goroutine is sending to messages, this case will execute
	default:
		fmt.Println("No message received")
	}

	// Prepare a message to send
	msg := "Hi"

	// SELECT BLOCK 2: Non-blocking send attempt
	select {
	// Try to send message to messages channel
	// This would block if there was no receiver and no default case
	case messages <- msg:
		fmt.Println("Sent message:", msg)
	// default case executes if send would block
	// Since no goroutine is receiving from messages, this case will execute
	default:
		fmt.Println("No message sent")
	}

	// SELECT BLOCK 3: Non-blocking multi-channel receive attempt
	select {
	// Try to receive from messages channel
	// Even though we attempted to send earlier, it failed due to no receiver
	case msg := <-messages:
		fmt.Println("Received message:", msg)
	// Try to receive from signals channel
	// No goroutine has sent any boolean signals
	case sig := <-signals:
		fmt.Println("Received signal:", sig)
	// default case executes when neither channel has data available
	// Since both channels are empty, this case will execute
	default:
		fmt.Println("No activity")
	}

	// EXECUTION FLOW ANALYSIS:
	// 1. First select: prints "No message received" (no data in messages)
	// 2. Second select: prints "No message sent" (no receiver for messages)
	// 3. Third select: prints "No activity" (no data in either channel)
}
