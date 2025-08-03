package main

import (
	"fmt"
	"math/rand"
	"time"
)

func worker(id int, ch chan<- string) {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		ch <- fmt.Sprintf("Worker %d: Job %d", id, i)
	}
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go worker(1, ch1)
	go worker(2, ch2)

	for i := 0; i < 10; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		}
	}
}
