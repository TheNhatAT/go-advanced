package main

// Large struct to demonstrate the difference
type LargeStruct struct {
	data [1000]int64 // 8KB of data
	id   int
	name string
}

// Small struct for comparison
type SmallStruct struct {
	x, y int
}

// Function that takes large struct by value
func processByValue(s LargeStruct) int64 {
	sum := int64(0)
	for _, v := range s.data {
		sum += v
	}
	return sum + int64(s.id)
}

// Function that takes large struct by pointer
func processByPointer(s *LargeStruct) int64 {
	sum := int64(0)
	for _, v := range s.data {
		sum += v
	}
	return sum + int64(s.id)
}

// Functions for small struct comparison
func processSmallByValue(s SmallStruct) int {
	return s.x + s.y
}

func processSmallByPointer(s *SmallStruct) int {
	return s.x + s.y
}
