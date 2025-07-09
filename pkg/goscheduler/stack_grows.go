package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func main() {
	debug.SetGCPercent(-1) // Disable GC for clearer example

	fmt.Printf("=== Stack Growth Demonstration ===\n")

	go func() {
		showStackInfo("Initial", 0)
		recursiveFunction(1, 50)
		runtime.Gosched()
		runtime.GC()
	}()

	time.Sleep(2 * time.Second) // Allow time for goroutine to run
}

func recursiveFunction(depth, maxDepth int) {
	// Large local array to consume stack space
	var localArray [500]int64 // 4KB per call
	localArray[0] = int64(depth)

	if depth <= 5 || depth%10 == 0 {
		showStackInfo("Recursive", depth)
	}

	if depth < maxDepth {
		recursiveFunction(depth+1, maxDepth)
	}
}

func showStackInfo(context string, depth int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("%s (depth %d): Stack usage pattern\n", context, depth)
	fmt.Printf("  Stack in use: %d KB\n", m.StackInuse/1024)
	fmt.Printf("  Stack from system: %d KB\n", m.StackSys/1024)
}
