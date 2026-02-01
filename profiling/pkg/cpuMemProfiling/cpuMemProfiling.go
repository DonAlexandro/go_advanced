package cpumemprofiling

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// generateData creates a large slice of random integers.
func generateData() []int {
	data := make([]int, 1000000)
	for i := range data {
		data[i] = rand.Intn(1000)
	}
	return data
}

// processCPUIntensive performs a simple summation to use CPU cycles.
func processCPUIntensive(data []int) {
	sum := 0
	for _, num := range data {
		sum += num
	}
	fmt.Println("Sum:", sum)
}

// processMemoryIntensive modifies the data in place.
func processMemoryIntensive(data []int) {
	for i := range data {
		data[i] *= 2
	}
}

func main() {
	// 1. CPU Profiling
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Error creating CPU profile file:", err)
		return
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		fmt.Println("Error starting CPU profiling:", err)
		return
	}
	defer pprof.StopCPUProfile()

	// Perform operations
	data := generateData()
	processCPUIntensive(data)
	processMemoryIntensive(data)

	// Wait for a few seconds to capture steady state
	time.Sleep(3 * time.Second)

	// 2. Memory (Heap) Profiling
	memFile, err := os.Create("mem.prof")
	if err != nil {
		fmt.Println("Error creating memory profile file:", err)
		return
	}
	defer memFile.Close()

	runtime.GC() // Optional: force GC to get a clean view of memory in use
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		fmt.Println("Error writing memory profile:", err)
		return
	}

	fmt.Println("Profile data written to cpu.prof and mem.prof")
}
