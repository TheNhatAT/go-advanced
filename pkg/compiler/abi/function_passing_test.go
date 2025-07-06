package main // Benchmark large struct by value
import "testing"

func BenchmarkLargeStructByValue(b *testing.B) {
	s := LargeStruct{
		id:   42,
		name: "test",
	}
	// Fill with some data
	for i := 0; i < 1000; i++ {
		s.data[i] = int64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processByValue(s) // Copies 8KB each call
	}
}

// Benchmark large struct by pointer
func BenchmarkLargeStructByPointer(b *testing.B) {
	s := LargeStruct{
		id:   42,
		name: "test",
	}
	// Fill with some data
	for i := 0; i < 1000; i++ {
		s.data[i] = int64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processByPointer(&s) // Copies only 8 bytes (pointer)
	}
}

// Benchmark small struct by value (should be similar performance)
func BenchmarkSmallStructByValue(b *testing.B) {
	s := SmallStruct{x: 10, y: 20}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processSmallByValue(s)
	}
}

// Benchmark small struct by pointer
func BenchmarkSmallStructByPointer(b *testing.B) {
	s := SmallStruct{x: 10, y: 20}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processSmallByPointer(&s)
	}
}

// Multiple calls benchmark to show cumulative effect
func BenchmarkMultipleCallsByValue(b *testing.B) {
	s := LargeStruct{id: 42, name: "test"}
	for i := 0; i < 1000; i++ {
		s.data[i] = int64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Multiple calls in tight loop
		_ = processByValue(s)
		_ = processByValue(s)
		_ = processByValue(s)
		_ = processByValue(s)
		_ = processByValue(s)
	}
}

func BenchmarkMultipleCallsByPointer(b *testing.B) {
	s := LargeStruct{id: 42, name: "test"}
	for i := 0; i < 1000; i++ {
		s.data[i] = int64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Multiple calls in tight loop
		_ = processByPointer(&s)
		_ = processByPointer(&s)
		_ = processByPointer(&s)
		_ = processByPointer(&s)
		_ = processByPointer(&s)
	}
}
