package main

import "testing"

func BenchmarkStackAccess(b *testing.B) {
	// Create slice of values (stored on stack when accessed)
	points := make([]Point, 1000)
	for i := range points {
		points[i] = Point{X: float64(i), Y: float64(i + 1)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		computeDistanceStack(points)
	}
}

func BenchmarkHeapAccessWithAllocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		points := make([]*Point, 100)
		for j := range points {
			points[j] = createHeapPoint(float64(j), float64(j+1))
		}
		computeDistanceHeap(points)
	}
}
