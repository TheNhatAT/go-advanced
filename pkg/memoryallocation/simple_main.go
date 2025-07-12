package main

import "time"

func main() {
	println("hello world")
	//go sampleGoroutine()
	//go sampleGoroutine()
	//println("main function finished")
}

func sampleGoroutine() {
	// This is a sample goroutine function.
	// It does nothing but serves as an example.
	println("This is a sample goroutine.")
	time.Sleep(2 * time.Second) // Simulate some work
}

// using "ps -p <PID> -o pid,rss,vsz" to check RSS and VSZ memory usage.
/**

ps -p 45595 -o pid,rss,vsz
  PID    RSS      VSZ
45595   3040 411314848

=> RSS: 3040 kB
=> VSZ: 411314848 kB
*/
