package main

import (
	"context"
	"fmt"
	"time"
)

func contextWithTimeoutExample() {
	// Create a context that automatically cancels after 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// Always call cancel to release resources, even if timeout happens first
	defer cancel()

	select {
	case <-time.After(5 * time.Second):
		// This block would execute if the operation finished within 5 seconds
		fmt.Println("Operation completed successfully")
	case <-ctx.Done():
		// This block executes when the 3-second timeout is reached
		// ctx.Err() will contain "context deadline exceeded"
		fmt.Printf("Stopped: %v\n", ctx.Err())
	}
}
