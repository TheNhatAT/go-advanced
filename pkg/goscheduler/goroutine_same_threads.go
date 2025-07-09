package main

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
)

func main_1() {
	runtime.GOMAXPROCS(2) // Limit to 2 OS threads for demonstration

	var wg sync.WaitGroup

	// Goroutine that will block on network I/O
	wg.Add(1)
	go func() {
		defer wg.Done()

		threadID := getThreadID()
		fmt.Printf("I/O Goroutine 1 starts on thread %d\n", threadID)

		// This will block on network I/O
		resp, err := http.Get("https://httpbin.org/delay/1")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		threadIDAfter := getThreadID()
		fmt.Printf("I/O Goroutine 1 resumes on thread %d (same: %v)\n",
			threadIDAfter, threadID == threadIDAfter)
	}()

	// Other goroutines to keep threads busy
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// CPU-bound work
			for j := 0; j < 1000000; j++ {
				_ = j * j
			}

			fmt.Printf("CPU goroutine %d on thread %d\n", id, getThreadID())
		}(i)
	}

	wg.Wait()
}

func getThreadID() int {
	// Simplified thread ID detection (not production code)
	buf := make([]byte, 64)
	n := runtime.Stack(buf, false)
	// Parse thread ID from stack trace (simplified)
	return int(buf[n-1]) // This is just for demonstration
}

/**
sample result:
	Goroutine 1 starts on thread 103
  	CPU goroutine 2 on thread 103
	CPU goroutine 0 on thread 103
	CPU goroutine 1 on thread 103
	Goroutine 1 resumes on thread 103 (same: true)
meaning: when the goroutine 1 blocks on network I/O, another goroutine can run on the same thread for prevent context switching.
*/
