I'll convert the examples to Go to demonstrate stack vs heap access performance differences.

## Key Differences in Go

### 1. **Direct Access vs Pointer Following**

**Stack Access (Direct):**
````go
func stackExample() {
    x := 42        // Stored directly on stack
    y := x + 10    // Direct access
}
````

**Heap Access (Indirect):**
````go
func heapExample() {
    x := new(int)  // Pointer on stack, data on heap
    *x = 42
    y := *x + 10   // Two operations: load pointer, then follow it
}
````

### 2. **Memory Layout Visualization**

**Stack Layout:**
```
Stack (contiguous memory):
┌─────────────┐ ← Stack Pointer
│   var_c     │ Direct offset
├─────────────┤
│   var_b     │ Direct offset  
├─────────────┤
│   var_a     │ Direct offset
└─────────────┘
```

**Heap Layout:**
```
Stack:                    Heap (scattered memory):
┌─────────────┐          ┌─────────────┐ ← 0x1000
│  ptr_to_c   │ ────────→│   data_c    │
├─────────────┤          └─────────────┘
│  ptr_to_b   │ ──────┐   
├─────────────┤        │  ┌─────────────┐ ← 0x2500
│  ptr_to_a   │ ──┐    └─→│   data_b    │
└─────────────┘   │       └─────────────┘
                  │
                  │       ┌─────────────┐ ← 0x3200
                  └──────→│   data_a    │
                          └─────────────┘
```

### 3. **Performance Comparison Example**

````go
// Stack version - Fast
func processStackData() {
    numbers := [5]int{1, 2, 3, 4, 5}  // Array on stack
    sum := 0
    
    for _, num := range numbers {
        sum += num  // Direct memory access
    }
}

// Heap version - Slower  
func processHeapData() {
    numbers := []int{1, 2, 3, 4, 5}  // Slice backing array on heap
    sum := 0
    
    for _, num := range numbers {
        sum += num  // Must follow pointer to heap data
    }
}
````

### 4. **Struct Comparison**

````go
// Stack allocation
type Point struct {
    X, Y float64
}

func stackComputation() {
    point := Point{X: 3.0, Y: 4.0}  // Allocated on stack
    distance := math.Sqrt(point.X*point.X + point.Y*point.Y)
    _ = distance
}

// Heap allocation
func heapComputation() {
    point := &Point{X: 3.0, Y: 4.0}  // Allocated on heap
    distance := math.Sqrt(point.X*point.X + point.Y*point.Y)
    _ = distance
}
````

### 5. **Benchmark Example**

````go
package main

import (
    "math"
    "testing"
)

type Point struct {
    X, Y float64
}

// Stack-based computation
func computeDistanceStack(points []Point) float64 {
    var totalDistance float64
    for _, point := range points {
        totalDistance += math.Sqrt(point.X*point.X + point.Y*point.Y)
    }
    return totalDistance
}

// Heap-based computation
func computeDistanceHeap(points []*Point) float64 {
    var totalDistance float64
    for _, point := range points {
        totalDistance += math.Sqrt(point.X*point.X + point.Y*point.Y)
    }
    return totalDistance
}

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

func BenchmarkHeapAccess(b *testing.B) {
    // Create slice of pointers (data stored on heap)
    points := make([]*Point, 1000)
    for i := range points {
        points[i] = &Point{X: float64(i), Y: float64(i + 1)}
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        computeDistanceHeap(points)
    }
}
````

### 6. **Memory Allocation Patterns**

````go
func demonstrateAllocation() {
    // Stack allocation - Go compiler decides
    var stackInt int = 42
    stackArray := [5]int{1, 2, 3, 4, 5}
    
    // Heap allocation - explicitly using new()
    heapInt := new(int)
    *heapInt = 42
    
    // Heap allocation - taking address forces heap allocation
    localInt := 42
    heapPtr := &localInt  // This forces localInt to heap
    
    // Slice always has backing array on heap
    heapSlice := make([]int, 5)
    heapSlice[0] = 1
    
    _ = stackInt
    _ = stackArray
    _ = heapInt
    _ = heapPtr
    _ = heapSlice
}
````

### 7. **Cache Performance Example**

````go
// Sequential access (cache-friendly)
func sequentialAccess() {
    data := [1000]int{}
    for i := 0; i < 1000; i++ {
        data[i] = i  // Sequential memory access
    }
}

// Pointer chasing (cache-unfriendly)
func pointerChasing() {
    type Node struct {
        value int
        next  *Node
    }
    
    // Create linked list on heap
    var head *Node
    for i := 0; i < 1000; i++ {
        node := &Node{value: i, next: head}
        head = node
    }
    
    // Traverse (each access potentially a cache miss)
    current := head
    for current != nil {
        _ = current.value
        current = current.next
    }
}
````

## Running the Benchmark

To see the performance difference:

````bash
go test -bench=. -benchmem
````

Expected output showing stack access is faster:
```
BenchmarkStackAccess-8    100000    15000 ns/op    0 B/op    0 allocs/op
BenchmarkHeapAccess-8     50000     25000 ns/op    8000 B/op  1000 allocs/op
```

## Why Stack is Faster in Go

1. **Fewer Memory Operations**: Stack requires direct offset calculation
2. **Cache Locality**: Stack data is contiguous and cache-friendly
3. **No GC Pressure**: Stack allocation doesn't trigger garbage collection
4. **Compiler Optimizations**: Go compiler can optimize stack access patterns

The performance difference is especially noticeable in tight loops and when processing large datasets.