package pointers

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	// --- Case 1: Without pointers (Value Semantics) ---
	p1 := Person{"John", 30}
	p2 := p1 // This creates a total CLONE of p1 on the stack.

	p2.Age = 35
	// p1.Age remains 30 because p2 is an independent copy.
	fmt.Println("p1:", p1) // p1: {John 30}
	fmt.Println("p2:", p2) // p2: {John 35}

	// --- Case 2: With pointers (Pointer Semantics) ---
	p3 := &Person{"Jane", 25} // p3 points to an address (likely on the heap)
	p4 := p3                  // p4 copies the ADDRESS, not the data.

	p4.Age = 28
	// p3.Age also becomes 28 because both variables point to the same spot.
	fmt.Println("p3:", *p3) // p3: {Jane 28}
	fmt.Println("p4:", *p4) // p4: {Jane 28}
}
