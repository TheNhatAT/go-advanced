package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

// Define the biggie struct as in your example
type biggie struct {
	huge  [1e8]byte // 100 MB array
	other *biggie   // pointer to another biggie
}

type smallStruct struct {
	a, b int // Small struct with 2 integers
}

// Example function with different argument types
func myFunction(
	arg1 int, // value type - 8 bytes copied
	arg2 *int, // pointer type - 8 bytes copied (pointer value)
	arg3 biggie, // large struct - ~100MB copied (stack/heap)
	arg4 *biggie, // pointer to large struct - 8 bytes copied
	arg5 []byte, // slice - 24 bytes copied (slice header)
	arg6 *[]byte, // pointer to slice - 8 bytes copied
	arg7 chan byte, // channel - 8 bytes copied (pointer-like)
	arg8 map[string]int, // map - 8 bytes copied (pointer-like)
	arg9 func(), // function - 8 bytes copied (pointer-like)
) {
	fmt.Printf("arg1 size: %d bytes\n", unsafe.Sizeof(arg1))
	fmt.Printf("arg2 size: %d bytes\n", unsafe.Sizeof(arg2))
	fmt.Printf("arg3 size: %d bytes\n", unsafe.Sizeof(arg3))
	fmt.Printf("arg4 size: %d bytes\n", unsafe.Sizeof(arg4))
	fmt.Printf("arg5 size: %d bytes\n", unsafe.Sizeof(arg5))
	fmt.Printf("arg6 size: %d bytes\n", unsafe.Sizeof(arg6))
	fmt.Printf("arg7 size: %d bytes\n", unsafe.Sizeof(arg7))
	fmt.Printf("arg8 size: %d bytes\n", unsafe.Sizeof(arg8))
	fmt.Printf("arg9 size: %d bytes\n", unsafe.Sizeof(arg9))
}

func demonstrateMemoryBehavior() {
	var m runtime.MemStats

	// Measure initial memory
	runtime.GC()
	runtime.ReadMemStats(&m)
	initialHeap := m.HeapAlloc

	// Example 1: Basic value types
	value := 42
	pointer := &value

	// Example 2: Large struct
	largeBiggie := biggie{
		huge:  [1e8]byte{}, // This will likely be allocated on heap
		other: nil,
	}

	// Example 3: Slice examples
	slice := make([]byte, 1000)
	slicePtr := &slice

	// Example 4: Channel, Map, Function
	ch := make(chan byte, 10)
	m1 := make(map[string]int)
	fn := func() { fmt.Println("Hello") }

	// Call function with different argument types
	myFunction(value, pointer, largeBiggie, &largeBiggie,
		slice, slicePtr, ch, m1, fn)

	// Measure memory after operations
	runtime.GC()
	runtime.ReadMemStats(&m)
	finalHeap := m.HeapAlloc

	fmt.Printf("\nMemory usage: %d bytes allocated\n", finalHeap-initialHeap)
}

// Examples of different scenarios
func valueVsPointerExamples() {
	fmt.Println("=== Value vs Pointer Examples ===")

	small := smallStruct{1, 2}
	modifySmallByValue(small)
	fmt.Printf("Original small struct after value call: %+v\n", small)

	modifySmallByPointer(&small)
	fmt.Printf("Original small struct after pointer call: %+v\n", small)

	// Large struct - pointer is more efficient
	large := &biggie{other: nil}
	processLargeByPointer(large) // Efficient - only 8 bytes copied

	// Slice examples
	sliceExample()

	// Channel and Map examples
	channelMapExample()
}

func modifySmallByValue(s smallStruct) {
	s.a = 100 // This won't affect the original
}

func modifySmallByPointer(s *smallStruct) {
	s.a = 200 // This will modify the original
}

func processLargeByPointer(b *biggie) {
	// Only 8 bytes copied for the pointer
	// Access to the large struct through pointer
	b.huge[0] = 1
}

func sliceExample() {
	fmt.Println("\n=== Slice Examples ===")

	// Original slice
	original := []int{1, 2, 3, 4, 5}

	// Pass by value - slice header copied (24 bytes)
	modifySliceByValue(original)
	fmt.Printf("After value modification: %v\n", original)

	// Pass by pointer - pointer to slice copied (8 bytes)
	modifySliceByPointer(&original)
	fmt.Printf("After pointer modification: %v\n", original)
}

func modifySliceByValue(s []int) {
	// This modifies the underlying array (shared)
	if len(s) > 0 {
		s[0] = 999
	}
	// This doesn't affect the original slice header
	s = append(s, 100)
}

func modifySliceByPointer(s *[]int) {
	// This modifies the underlying array
	if len(*s) > 0 {
		(*s)[0] = 888
	}
	// This also modifies the original slice header
	*s = append(*s, 200)
}

func channelMapExample() {
	fmt.Println("\n=== Channel and Map Examples ===")

	// Channel example
	ch := make(chan int, 5)
	processChannel(ch) // Only pointer copied

	// Receive from channel
	select {
	case val := <-ch:
		fmt.Printf("Received from channel: %d\n", val)
	default:
		fmt.Println("Channel is empty")
	}

	// Map example
	m := make(map[string]int)
	processMap(m) // Only pointer copied
	fmt.Printf("Map after processing: %v\n", m)
}

func processChannel(ch chan int) {
	// Channel is reference type - changes affect original
	ch <- 42
}

func processMap(m map[string]int) {
	// Map is reference type - changes affect original
	m["key"] = 100
}

// Method receiver examples
type Counter struct {
	value int
}

// Value receiver - copy of Counter
func (c Counter) IncrementByValue() {
	c.value++ // This won't affect the original
}

// Pointer receiver - reference to Counter
func (c *Counter) IncrementByPointer() {
	c.value++ // This will modify the original
}

func receiverExample() {
	fmt.Println("\n=== Method Receiver Examples ===")

	counter := Counter{value: 0}

	counter.IncrementByValue()
	fmt.Printf("After value receiver: %d\n", counter.value) // Still 0

	counter.IncrementByPointer()
	fmt.Printf("After pointer receiver: %d\n", counter.value) // Now 1
}

func main() {
	demonstrateMemoryBehavior()
	valueVsPointerExamples()
	receiverExample()
}
