// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package concurrency

import (
	"sync"
	"testing"
)

func BenchmarkSharingWithAtomic(b *testing.B) {
	// Setup deterministic random function for consistent benchmarking
	var counter int64
	var mu sync.Mutex
	randInt64 = func() int64 {
		mu.Lock()
		defer mu.Unlock()
		counter++
		return counter
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sharingWithAtomic()
	}
}

func BenchmarkSharingWithMutex(b *testing.B) {
	// Setup deterministic random function for consistent benchmarking
	var counter int64
	var mu sync.Mutex
	randInt64 = func() int64 {
		mu.Lock()
		defer mu.Unlock()
		counter++
		return counter
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sharingWithMutex()
	}
}

func BenchmarkSharingWithChannel(b *testing.B) {
	// Setup deterministic random function for consistent benchmarking
	var counter int64
	var mu sync.Mutex
	randInt64 = func() int64 {
		mu.Lock()
		defer mu.Unlock()
		counter++
		return counter
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sharingWithChannel()
	}
}

func BenchmarkSharingWithShardedSpace(b *testing.B) {
	// Setup deterministic random function for consistent benchmarking
	var counter int64
	var mu sync.Mutex
	randInt64 = func() int64 {
		mu.Lock()
		defer mu.Unlock()
		counter++
		return counter
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sharingWithShardedSpace()
	}
}

// Comparative benchmark to measure all concurrency patterns in parallel
func BenchmarkConcurrencyPatterns(b *testing.B) {
	// Setup deterministic random function
	var counter int64
	var mu sync.Mutex
	randInt64 = func() int64 {
		mu.Lock()
		defer mu.Unlock()
		counter++
		return counter
	}

	b.Run("Atomic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sharingWithAtomic()
		}
	})

	b.Run("Mutex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sharingWithMutex()
		}
	})

	b.Run("Channel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sharingWithChannel()
		}
	})

	b.Run("ShardedSpace", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sharingWithShardedSpace()
		}
	})
}
