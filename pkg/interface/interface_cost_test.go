package _interface

import (
	"io"
	"testing"
)

// Simple benchmark comparing interface vs direct calls
func BenchmarkInterfaceCallSimple(b *testing.B) {
	var z zero
	var g getter
	g = z
	b.Run("via interface", func(b *testing.B) {
		total := 0
		for i := 0; i < b.N; i++ {
			total += g.get()
		}
		if total > 0 {
			b.Logf("total is %d", total)
		}
	})
	b.Run("direct", func(b *testing.B) {
		total := 0
		for i := 0; i < b.N; i++ {
			total += z.get()
		}
		if total > 0 {
			b.Logf("total is %d", total)
		}
	})
}

// Benchmark showing allocation overhead with interfaces
func BenchmarkInterfaceAlloc(b *testing.B) {
	var z zeroReader
	var r io.Reader
	r = z
	b.Run("via interface", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var buf [7]byte
			r.Read(buf[:])
		}
	})
	b.Run("direct", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var buf [7]byte
			z.Read(buf[:])
		}
	})
}
