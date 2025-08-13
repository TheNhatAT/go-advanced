package optimizationpatterns

import (
	"testing"
)

func TestGcLargeHeapPointer(t *testing.T) {
	pointer8GB()

	// result:
	/**
	=== RUN   TestGcLargeHeapPointer
	PID: 87824
	GC took 77.844667ms
	GC took 24.907417ms
	GC took 26.334833ms
	GC took 25.111333ms
	GC took 30.450625ms
	GC took 41.127458ms
	GC took 31.541792ms
	GC took 25.049333ms
	GC took 24.485166ms
	GC took 25.500125ms
	*/
}

func TestGcLargeHeapValue(t *testing.T) {
	value8GB()

	// result:
	/**
	=== RUN   TestGcLargeHeapValue
	PID: 87927
	GC took 408.375µs
	GC took 166.583µs
	GC took 153.958µs
	GC took 160.875µs
	GC took 255.417µs
	GC took 185.834µs
	GC took 188.875µs
	GC took 179.958µs
	GC took 183.583µs
	GC took 230.375µs
	*/
}

func TestGcLargeHeapOsDirectlyMapMemory(t *testing.T) {
	osDirectlyMapMemory8GB()

	// result:
	/**
	=== RUN   TestGcLargeHeapOsDirectlyMapMemory
	GC took 614.542µs
	GC took 242µs
	GC took 244.75µs
	GC took 229.625µs
	GC took 198µs
	GC took 226.541µs
	GC took 258.709µs
	GC took 173.459µs
	GC took 153.25µs
	GC took 202.375µs
	*/
}

// == benchmarks ==

/*
*
	export ver=v1 && \
		go test -run '^$' -bench '^BenchmarkGcLargeHeapPointer' -benchtime 10s -count 6 \
			-cpu 4 \
			-benchmem \
			-memprofile=./benchmarkresult/${ver}.mem.pprof -cpuprofile=./benchmarkresult/${ver}.cpu.pprof \
		| tee ./benchmarkresult/${ver}.txt
*/
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/benchmarkresult/v1.txt ./pkg/benchmark/micro/benchmarkresult/v1.txt
*/
func BenchmarkGcLargeHeapPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pointer8GB()
	}
}

/*
*
	export ver=v2 && \
		go test -run '^$' -bench '^BenchmarkGcLargeHeapValue' -benchtime 10s -count 6 \
			-cpu 4 \
			-benchmem \
			-memprofile=./benchmarkresult/${ver}.mem.pprof -cpuprofile=./benchmarkresult/${ver}.cpu.pprof \
		| tee ./benchmarkresult/${ver}.txt
*/
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/benchmarkresult/v2.txt ./pkg/benchmark/micro/benchmarkresult/v2.txt
*/
func BenchmarkGcLargeHeapValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		value8GB()
	}
}

/*
*
	export ver=v3 && \
		go test -run '^$' -bench '^BenchmarkGcLargeHeapOffHeap' -benchtime 10s -count 6 \
			-cpu 4 \
			-benchmem \
			-memprofile=./benchmarkresult/${ver}.mem.pprof -cpuprofile=./benchmarkresult/${ver}.cpu.pprof \
		| tee ./benchmarkresult/${ver}.txt
*/
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/benchmarkresult/v3.txt ./pkg/benchmark/micro/benchmarkresult/v3.txt
*/
func BenchmarkGcLargeHeapOffHeap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		osDirectlyMapMemory8GB()
	}
}
