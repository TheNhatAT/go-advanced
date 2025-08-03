package micro

import (
	"github.com/efficientgo/core/testutil"
	"github.com/felixge/fgprof"
	"os"
	"testing"
)

/**
export ver=v1 && \
	go test -run '^$' -bench '^BenchmarkConcurrentSum1' -benchtime 10s -count 6 \
		-cpu 4 \
		-benchmem \
		-memprofile=./concurrentbenchmarkresult/${ver}.mem.pprof -cpuprofile=./concurrentbenchmarkresult/${ver}.cpu.pprof \
	| tee ./concurrentbenchmarkresult/${ver}.txt
*/
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/concurrentbenchmarkresult/v1.txt ./pkg/benchmark/micro/concurrentbenchmarkresult/v1.txt
*/
func BenchmarkConcurrentSum1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ConcurrentSum1("testdata/test.2000000.txt")
	}
}

/**
export ver=v2 && \
	go test -run '^$' -bench '^BenchmarkConcurrentSum2' -benchtime 10s -count 6 \
		-cpu 4 \
		-benchmem \
		-memprofile=./concurrentbenchmarkresult/${ver}.mem.pprof -cpuprofile=./concurrentbenchmarkresult/${ver}.cpu.pprof \
	| tee ./concurrentbenchmarkresult/${ver}.txt
*/
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/concurrentbenchmarkresult/v2.txt ./pkg/benchmark/micro/concurrentbenchmarkresult/v2.txt
*/
func BenchmarkConcurrentSum2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ConcurrentSum2("testdata/test.2000000.txt", 4)
	}
}

/**
export ver=v3 && \
	go test -run '^$' -bench '^BenchmarkConcurrentSum3' -benchtime 10s -count 6 \
		-cpu 4 \
		-benchmem \
		-memprofile=./concurrentbenchmarkresult/${ver}.mem.pprof -cpuprofile=./concurrentbenchmarkresult/${ver}.cpu.pprof \
	| tee ./concurrentbenchmarkresult/${ver}.txt
*/
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/concurrentbenchmarkresult/v3.txt ./pkg/benchmark/micro/concurrentbenchmarkresult/v3.txt
*/
func BenchmarkConcurrentSum3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ConcurrentSum3("testdata/test.2000000.txt", 4)
	}
}

//
/**
$ export ver=v3_fg && go test -run '^$' -bench '^BenchmarkConcurrentSum3_fgprof' \
  -benchtime 10s -count 6 -cpu 4 | tee ./concurrentbenchmarkresult/${ver}.txt
*/
func BenchmarkConcurrentSum3_fgprof(b *testing.B) {
	f, err := os.Create("concurrentbenchmarkresult/v3_fg.pprof")
	testutil.Ok(b, err)

	defer func() { testutil.Ok(b, f.Close()) }()

	closeFunc := fgprof.Start(f, fgprof.FormatPprof)
	BenchmarkConcurrentSum3(b)
	testutil.Ok(b, closeFunc())
}

/**
export ver=v4 && \
	go test -run '^$' -bench '^BenchmarkConcurrentSum4' -benchtime 10s -count 6 \
		-cpu 4 \
		-benchmem \
		-memprofile=./concurrentbenchmarkresult/${ver}.mem.pprof -cpuprofile=./concurrentbenchmarkresult/${ver}.cpu.pprof \
	| tee ./concurrentbenchmarkresult/${ver}.txt
*/
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/concurrentbenchmarkresult/v4.txt ./pkg/benchmark/micro/concurrentbenchmarkresult/v4.txt
*/
func BenchmarkConcurrentSum4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ConcurrentSum4("testdata/test.2000000.txt", 4)
	}
}

//
/**
$ export ver=v4_fg && go test -run '^$' -bench '^BenchmarkConcurrentSum4_fgprof' \
  -benchtime 10s -count 6 -cpu 4 | tee ./concurrentbenchmarkresult/${ver}.txt
*/
func BenchmarkConcurrentSum4_fgprof(b *testing.B) {
	f, err := os.Create("concurrentbenchmarkresult/v4_fg.pprof")
	testutil.Ok(b, err)

	defer func() { testutil.Ok(b, f.Close()) }()

	closeFunc := fgprof.Start(f, fgprof.FormatPprof)
	BenchmarkConcurrentSum4(b)
	testutil.Ok(b, closeFunc())
}
