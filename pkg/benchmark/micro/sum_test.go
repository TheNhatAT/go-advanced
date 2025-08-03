package micro

import (
	"fmt"
	"github.com/efficientgo/core/testutil"
	"github.com/felixge/fgprof"
	"os"
	"testing"
)

// options cli for running benchmarks:
// $ go test -run '^$' -bench '^BenchmarkSum$' -> default options
// $ go test -run '^$' -bench '^BenchmarkSum$' -benchtime 10s -> for specific time duration
// $ go test -run '^$' -bench '^BenchmarkSum$' -benchtime 100x -> for specific number of iterations
// $ go test -run '^$' -bench '^BenchmarkSum$' -benchtime 1s -count 5 -> for calculate variance between runs
// one line shell command to run benchmark:
/**
export ver=v1 && \
	go test -run '^$' -bench '^BenchmarkSum$' -benchtime 10s -count 6 \
		-cpu 4 \
		-benchmem \
		-memprofile=./benchmarkresult/${ver}.mem.pprof -cpuprofile=./benchmarkresult/${ver}.cpu.pprof \
	| tee ./benchmarkresult/${ver}.txt
*/
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/benchmarkresult/v1.txt
*/
func BenchmarkSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Sum("testdata/test.2000000.txt")
	}
}

/**
export ver=v2 && \
	go test -run '^$' -bench '^BenchmarkSum2' -benchtime 10s -count 6 \
		-cpu 4 \
		-benchmem \
		-memprofile=./benchmarkresult/${ver}.mem.pprof -cpuprofile=./benchmarkresult/${ver}.cpu.pprof \
	| tee ./benchmarkresult/${ver}.txt
*/
// after benchmark, for running benchstat, go into v2.txt and rename the BenchmarkSum2 -> BenchmarkSum
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/benchmarkresult/v1.txt ./pkg/benchmark/micro/benchmarkresult/v2.txt
*/
func BenchmarkSum2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Sum2("testdata/test.2000000.txt")
	}
}

/**
export ver=v3 && \
	go test -run '^$' -bench '^BenchmarkSum3' -benchtime 10s -count 6 \
		-cpu 4 \
		-benchmem \
		-memprofile=./benchmarkresult/${ver}.mem.pprof -cpuprofile=./benchmarkresult/${ver}.cpu.pprof \
	| tee ./benchmarkresult/${ver}.txt
*/
// after benchmark, for running benchstat, go into v3.txt and rename the BenchmarkSum3 -> BenchmarkSum
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/benchmarkresult/v1.txt ./pkg/benchmark/micro/benchmarkresult/v3.txt
*/
func BenchmarkSum3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Sum3("testdata/test.2000000.txt")
	}
}

/**
export ver=v4 && \
	go test -run '^$' -bench '^BenchmarkSum4' -benchtime 10s -count 6 \
		-cpu 4 \
		-benchmem \
		-memprofile=./benchmarkresult/${ver}.mem.pprof -cpuprofile=./benchmarkresult/${ver}.cpu.pprof \
	| tee ./benchmarkresult/${ver}.txt
*/
// after benchmark, for running benchstat, go into v4.txt and rename the BenchmarkSum4 -> BenchmarkSum
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/benchmarkresult/v1.txt ./pkg/benchmark/micro/benchmarkresult/v4.txt
*/
func BenchmarkSum4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Sum4("testdata/test.2000000.txt")
	}
}

/**
export ver=v5 && \
	go test -run '^$' -bench '^BenchmarkSum5' -benchtime 10s -count 6 \
		-cpu 4 \
		-benchmem \
		-memprofile=./benchmarkresult/${ver}.mem.pprof -cpuprofile=./benchmarkresult/${ver}.cpu.pprof \
	| tee ./benchmarkresult/${ver}.txt
*/
// after benchmark, for running benchstat, go into v5.txt and rename the BenchmarkSum5 -> BenchmarkSum
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/benchmarkresult/v1.txt ./pkg/benchmark/micro/benchmarkresult/v5.txt
*/
func BenchmarkSum5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Sum5("testdata/test.2000000.txt")
	}
}

/**
export ver=v6 && \
	go test -run '^$' -bench '^BenchmarkSum6' -benchtime 10s -count 6 \
		-cpu 4 \
		-benchmem \
		-memprofile=./benchmarkresult/${ver}.mem.pprof -cpuprofile=./benchmarkresult/${ver}.cpu.pprof \
	| tee ./benchmarkresult/${ver}.txt
*/
// after benchmark, for running benchstat, go into v6.txt and rename the BenchmarkSum6 -> BenchmarkSum
/**
using benchstat for visualization:
$ gvm use go1.24.1
$ benchstat ./pkg/benchmark/micro/benchmarkresult/v1.txt ./pkg/benchmark/micro/benchmarkresult/v6.txt
*/
func BenchmarkSum6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Sum6("testdata/test.2000000.txt")
	}
}

// For correctness of the benchmark, we should use testutil.TB interface.
// unittest
func TestBenchmarkSum_unittest(t *testing.T) {
	benchmarkSum(testutil.NewTB(t))
}

// benchmark
func BenchmarkSum_benchmark(b *testing.B) {
	benchmarkSum(testutil.NewTB(b))
}

func benchmarkSum(tb testutil.TB) {
	for i := 0; i < tb.N(); i++ {
		ret, err := Sum("testdata/test.2000000.txt")
		testutil.Ok(tb, err)
		if !tb.IsBenchmark() {
			// More expensive result checks can be here.
			testutil.Equals(tb, int64(6221600000), ret)
		}
	}
}

// example of well-documented benchmark
//
// BenchmarkSum assesses `Sum` function.
// NOTE(bwplotka): Test it with a maximum of 4 CPU cores, given we don't allocate
// more in our production containers.
//
// Recommended run options:
/*
export ver=v1 && go test \
	-run '^$' -bench '^BenchmarkSum$' \
	-benchtime 10s -count 6 -cpu 4 -benchmem \
	-memprofile=${ver}.mem.pprof -cpuprofile=${ver}.cpu.pprof \
| tee ${ver}.txt
*/
func BenchmarkSum_well_document(b *testing.B) {
	// Create 7.55 MB file with 2 million lines.
	fn := "testdata/test.2000000.txt"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Sum(fn)
		testutil.Ok(b, err)
	}
}

// example of table tests for different input sizes
func BenchmarkSum_table_test(b *testing.B) {
	for _, tcase := range []struct {
		numLines int
	}{
		{numLines: 0},
		{numLines: 1e2},
		{numLines: 1e4},
		{numLines: 1e6},
		{numLines: 2e6},
	} {
		b.Run(fmt.Sprintf("lines-%d", tcase.numLines), func(b *testing.B) {
			b.ReportAllocs() // go test ignores any benchmark methods outside b.Run => remember to repeat them here
			//fn := lazyCreateTestInput(tb, tcase.numLines)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Sum("fn")
				testutil.Ok(b, err)
			}
		})
	}
}

// BenchmarkSum_fgprof recommended run options:
//
/**
$ export ver=v1fg && go test -run '^$' -bench '^BenchmarkSum_fgprof' \
  -benchtime 10s -count 6 -cpu 4 | tee ./benchmarkresult/${ver}.txt
*/
// Read more in "Efficient Go"; Example 10-2.
func BenchmarkSum_fgprof(b *testing.B) {
	f, err := os.Create("./benchmarkresult/v1fg_fgprof.pprof")
	testutil.Ok(b, err)

	defer func() { testutil.Ok(b, f.Close()) }()

	closeFn := fgprof.Start(f, fgprof.FormatPprof)
	BenchmarkSum(b)
	testutil.Ok(b, closeFn())
}

// BenchmarkSum5_fgprof recommended run options:
//
/**
$ export ver=v5fg && go test -run '^$' -bench '^BenchmarkSum5_fgprof' \
  -benchtime 10s -count 6 -cpu 4 | tee ./benchmarkresult/${ver}.txt
*/
// Read more in "Efficient Go"; Example 10-2.
func BenchmarkSum5_fgprof(b *testing.B) {
	f, err := os.Create("./benchmarkresult/v5fg_fgprof.pprof")
	testutil.Ok(b, err)

	defer func() { testutil.Ok(b, f.Close()) }()

	closeFn := fgprof.Start(f, fgprof.FormatPprof)
	BenchmarkSum5(b)
	testutil.Ok(b, closeFn())
}

// BenchmarkSum6_fgprof recommended run options:
//
/**
$ export ver=v6fg && go test -run '^$' -bench '^BenchmarkSum6_fgprof' \
  -benchtime 10s -count 6 -cpu 4 | tee ./benchmarkresult/${ver}.txt
*/
// Read more in "Efficient Go"; Example 10-2.
func BenchmarkSum6_fgprof(b *testing.B) {
	f, err := os.Create("./benchmarkresult/v6fg_fgprof.pprof")
	testutil.Ok(b, err)

	defer func() { testutil.Ok(b, f.Close()) }()

	closeFn := fgprof.Start(f, fgprof.FormatPprof)
	BenchmarkSum6(b)
	testutil.Ok(b, closeFn())
}
