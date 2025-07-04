package main

import (
	"math"
)

type Point struct {
	X, Y float64
}

// Stack-based computation - values stay on stack
func computeDistanceStack(points []Point) float64 {
	var totalDistance float64
	for _, point := range points {
		totalDistance += math.Sqrt(point.X*point.X + point.Y*point.Y)
	}
	return totalDistance
}

// Heap-based computation - force heap allocation by returning pointer
func createHeapPoint(x, y float64) *Point {
	return &Point{X: x, Y: y} // Escapes to heap because we return pointer
}

func computeDistanceHeap(points []*Point) float64 {
	var totalDistance float64
	for _, point := range points {
		totalDistance += math.Sqrt(point.X*point.X + point.Y*point.Y)
	}
	return totalDistance
}

// Alternative: use interface{} to force heap allocation
func forceHeapAllocation(x, y float64) *Point {
	var p interface{} = &Point{X: x, Y: y}
	return p.(*Point)
}
