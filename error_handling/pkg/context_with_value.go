package main

import (
	"context"
	"fmt"
)

// Define a private type for keys to prevent collisions with other packages
type key int

// iota helps define unique constants
const requestIDKey key = iota

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// GetRequestID retrieves the request ID from the context safely
func GetRequestID(ctx context.Context) string {
	val, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return "unknown" // Return a default value if key not found
	}
	return val
}

func contextWithValueExample() {
	// Root context
	ctx := context.Background()

	// Wrap context with metadata
	ctx = WithRequestID(ctx, "req-12345")

	// Pass context to other parts of your app
	fmt.Println("Current Request ID:", GetRequestID(ctx))
}
