package datastructures

import "fmt"

func main() {
	// --- Case 1: Slice inside a Map (Sequential Memory) ---
	// Good for: Iterating over values, keeping order, memory density.
	sliceMap := make(map[int][]int)
	sliceMap[1] = append(sliceMap[1], 10)
	sliceMap[1] = append(sliceMap[1], 20)
	fmt.Println("sliceMap:", sliceMap)

	// --- Case 2: Map inside a Map (Hashed Memory) ---
	// Good for: "Is 10 inside key 1?" (O(1) lookups), deduplication.
	mapMap := make(map[int]map[int]bool)
	mapMap[1] = make(map[int]bool)
	mapMap[1][10] = true
	mapMap[1][20] = true
	fmt.Println("mapMap:", mapMap)
}
