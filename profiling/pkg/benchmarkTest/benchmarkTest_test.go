package benchmarktest

import "testing"

// BenchmarkFibonacci measures the performance of the Fibonacci function.
func BenchmarkFibonacci(b *testing.B) {
	// b.N is automatically adjusted by the testing framework until the benchmark
	// lasts long enough to be statistically significant.
	for b.Loop() {
		Fibonacci(10)
	}
}
