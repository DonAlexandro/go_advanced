package pkg

import "time"

// This is a simple Go program that demonstrates the use of goroutines.
func sayHello() {
	// This function prints a greeting message.
	println("Hello, World!")
}

func main() {
	go sayHello()

	time.Sleep(1 * time.Second) // Wait for the goroutine to finish
}
