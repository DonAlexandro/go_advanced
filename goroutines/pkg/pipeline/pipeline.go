package pipeline

import (
	"fmt"
	"math/rand"
	"time"
)

// PIPELINE PATTERN: Data flows through multiple processing stages

// generateNumbers creates a pipeline stage that produces random numbers
func generateNumbers(count int, rng *rand.Rand) <-chan int {
	// Create unbuffered output channel
	// Unbuffered channel provides backpressure - blocks if consumer is slow
	out := make(chan int)

	// Start goroutine to produce numbers asynchronously
	go func() {
		// defer ensures channel is closed when producer finishes
		// Closing signals downstream consumers that no more data is coming
		defer close(out)

		// Generate 'count' random numbers
		for range count {
			// Send random number (0-99) to output channel
			// This blocks until downstream stage is ready to receive
			out <- rng.Intn(100)
		}
		// After loop, defer close(out) executes
		// This tells downstream stage "no more numbers coming"
	}()

	// Return receive-only channel immediately (non-blocking)
	// Caller can start consuming while producer runs in background
	return out
}

// squareNumbers creates a pipeline stage that transforms input numbers
func squareNumbers(in <-chan int) <-chan int {
	// Create unbuffered output channel for squared results
	out := make(chan int)

	// Start goroutine to process input stream asynchronously
	go func() {
		// defer ensures output channel is properly closed
		// Critical for pipeline termination propagation
		defer close(out)

		// Process each number from input channel until it's closed
		// Range automatically handles channel closing - exits when in is closed and drained
		for num := range in {
			// Transform input (square the number) and send to output
			// This creates a 1:1 transformation in the pipeline
			out <- num * num
		}
		// When input channel closes and all numbers processed,
		// defer close(out) signals next stage that transformation is complete
	}()

	// Return receive-only channel immediately (non-blocking)
	// Allows caller to start consuming transformed data
	return out
}

func main() {
	// Create random number generator with current time seed
	// Ensures different random sequences on each program run
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// PIPELINE CONSTRUCTION: Chain processing stages together

	// Stage 1: Generate 10 random numbers (Producer)
	// Returns immediately with channel, generation happens in background
	numbers := generateNumbers(10, rng)

	// Stage 2: Square each number from stage 1 (Transformer)
	// Takes output from generateNumbers as input
	// Returns immediately with channel, processing happens in background
	squares := squareNumbers(numbers)

	// PIPELINE EXECUTION: Consume final results (Consumer)
	// Range over final stage output until channel is closed
	// This drives the entire pipeline - backpressure flows upstream
	for square := range squares {
		// Print each squared result as it becomes available
		// Pipeline executes concurrently: generating, squaring, and printing happen simultaneously
		fmt.Println(square)
	}

	// PIPELINE FLOW ANALYSIS:
	// 1. generateNumbers goroutine produces random number → sends to numbers channel
	// 2. squareNumbers goroutine receives from numbers → squares it → sends to squares channel
	// 3. main goroutine receives from squares → prints result
	// 4. Process repeats until generateNumbers closes numbers channel
	// 5. squareNumbers detects closed numbers, processes remaining, closes squares
	// 6. main detects closed squares, exits range loop, program terminates
	//
	// CONCURRENCY BENEFITS:
	// - Overlapped execution: while printing result N, generating number N+1, squaring number N+2
	// - Memory efficiency: only one number in each stage at a time (unbuffered channels)
	// - Backpressure handling: slow consumer automatically slows down entire pipeline
}
