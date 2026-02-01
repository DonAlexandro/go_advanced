package benchmarktest

// Fibonacci calculates the nth number in the Fibonacci sequence using recursion.
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}
