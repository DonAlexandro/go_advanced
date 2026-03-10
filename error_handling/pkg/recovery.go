package main

import "fmt"

func recoverPanic() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from panic:", err)
		}
	}()

	panic("Something went wrong!")
}
