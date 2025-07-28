package micro

import (
	"math"
	"testing"
)

const m1 = 0x5555555555555555
const m2 = 0x3333333333333333
const m4 = 0x0f0f0f0f0f0f0f0f
const h01 = 0x0101010101010101

func popcnt(x uint64) uint64 {
	x -= (x >> 1) & m1
	x = (x & m2) + ((x >> 2) & m2)
	x = (x + (x >> 4)) & m4
	return (x * h01) >> 56
}

// using "go test -run '^$' -bench '^BenchmarkPopcnt$' -benchtime 10s -count 6 -cpu 4 -benchmem"
// to run the benchmark.
func BenchmarkPopcnt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcnt(math.MaxUint64)
	}
}

/** result of the benchmark: -> too good to be true, dead code elimination kicks in
goos: darwin
goarch: arm64
pkg: go-advanced/pkg/benchmark/micro
BenchmarkPopcnt-4       1000000000               0.8735 ns/op          0 B/op          0 allocs/op
BenchmarkPopcnt-4       1000000000               0.8508 ns/op          0 B/op          0 allocs/op
BenchmarkPopcnt-4       1000000000               0.8495 ns/op          0 B/op          0 allocs/op
BenchmarkPopcnt-4       1000000000               0.8511 ns/op          0 B/op          0 allocs/op
BenchmarkPopcnt-4       1000000000               0.8554 ns/op          0 B/op          0 allocs/op
BenchmarkPopcnt-4       1000000000               0.8516 ns/op          0 B/op          0 allocs/op
*/

/** result of running: go test -run '^$' -gcflags="-S" -bench '^BenchmarkPopcnt$' => dead code elimination kicks in
...
0x0000 00000 (.../dead_code_elimination_test.go:22)        PCDATA  $3, $1
0x0000 00000 (.../dead_code_elimination_test.go:22)        MOVD    ZR, R1
0x0004 00004 (.../dead_code_elimination_test.go:23)        JMP     12
0x0008 00008 (.../dead_code_elimination_test.go:23)        ADD     $1, R1, R1
0x000c 00012 (.../dead_code_elimination_test.go:23)        MOVD    416(R0), R2
0x0010 00016 (.../dead_code_elimination_test.go:23)        CMP     R2, R1
0x0014 00020 (.../dead_code_elimination_test.go:23)        BLT     8
0x0018 00024 (.../dead_code_elimination_test.go:26)        RET     (R30)
...
*/

/** result of running: go test -run '^$' -gcflags="-S -N" -bench '^BenchmarkPopcnt$' => force turn off compiler optimizations
0x0000 00000 (.../dead_code_elimination_test.go:22)        TEXT    go-advanced/pkg/benchmark/micro.BenchmarkPopcnt(SB), NOSPLIT|LEAF|ABIInternal, $48-8
0x0000 00000 (.../dead_code_elimination_test.go:22)        MOVD.W  R30, -48(RSP)
0x0004 00004 (.../dead_code_elimination_test.go:22)        MOVD    R29, -8(RSP)
0x0008 00008 (.../dead_code_elimination_test.go:22)        SUB     $8, RSP, R29
0x000c 00012 (.../dead_code_elimination_test.go:22)        FUNCDATA        $0, gclocals·wgcWObbY2HYnK2SU/U22lA==(SB)
0x000c 00012 (.../dead_code_elimination_test.go:22)        FUNCDATA        $1, gclocals·J5F+7Qw7O7ve2QcWC7DpeQ==(SB)
0x000c 00012 (.../dead_code_elimination_test.go:22)        FUNCDATA        $5, go-advanced/pkg/benchmark/micro.BenchmarkPopcnt.arginfo1(SB)
0x000c 00012 (.../dead_code_elimination_test.go:22)        MOVD    R0, go-advanced/pkg/benchmark/micro.b(FP)
0x0010 00016 (.../dead_code_elimination_test.go:23)        MOVD    ZR, go-advanced/pkg/benchmark/micro.i-8(SP)
0x0014 00020 (.../dead_code_elimination_test.go:23)        JMP     24
0x0018 00024 (.../dead_code_elimination_test.go:23)        MOVD    go-advanced/pkg/benchmark/micro.b(FP), R0
0x001c 00028 (.../dead_code_elimination_test.go:23)        PCDATA  $0, $-2
0x001c 00028 (.../dead_code_elimination_test.go:23)        MOVB    (R0), R27
0x0020 00032 (.../dead_code_elimination_test.go:23)        PCDATA  $0, $-1
0x0020 00032 (.../dead_code_elimination_test.go:23)        MOVD    go-advanced/pkg/benchmark/micro.i-8(SP), R1
0x0024 00036 (.../dead_code_elimination_test.go:23)        MOVD    416(R0), R0
0x0028 00040 (.../dead_code_elimination_test.go:23)        CMP     R1, R0
0x002c 00044 (.../dead_code_elimination_test.go:23)        BGT     52
0x0030 00048 (.../dead_code_elimination_test.go:23)        JMP     116
0x0034 00052 (.../dead_code_elimination_test.go:24)        MOVD    $-1, R0
0x0038 00056 (.../dead_code_elimination_test.go:24)        MOVD    R0, go-advanced/pkg/benchmark/micro.x-16(SP)
0x003c 00060 (<unknown line number>)    NOP
0x003c 00060 (.../dead_code_elimination_test.go:14)        MOVD    $-6148914691236517206, R1
0x0040 00064 (.../dead_code_elimination_test.go:14)        MOVD    R1, go-advanced/pkg/benchmark/micro.x-16(SP)
0x0044 00068 (.../dead_code_elimination_test.go:15)        MOVD    $4919131752989213764, R2
0x0048 00072 (.../dead_code_elimination_test.go:15)        MOVD    R2, go-advanced/pkg/benchmark/micro.x-16(SP)
0x004c 00076 (.../dead_code_elimination_test.go:16)        MOVD    $578721382704613384, R3
0x0050 00080 (.../dead_code_elimination_test.go:16)        MOVD    R3, go-advanced/pkg/benchmark/micro.x-16(SP)
0x0054 00084 (.../dead_code_elimination_test.go:24)        MOVD    $64, R4
0x0058 00088 (.../dead_code_elimination_test.go:24)        MOVD    R4, go-advanced/pkg/benchmark/micro.~R0-24(SP)
0x005c 00092 (.../dead_code_elimination_test.go:24)        JMP     96
0x0060 00096 (.../dead_code_elimination_test.go:23)        JMP     100
0x0064 00100 (.../dead_code_elimination_test.go:23)        MOVD    go-advanced/pkg/benchmark/micro.i-8(SP), R5
0x0068 00104 (.../dead_code_elimination_test.go:23)        ADD     $1, R5, R5
0x006c 00108 (.../dead_code_elimination_test.go:23)        MOVD    R5, go-advanced/pkg/benchmark/micro.i-8(SP)
0x0070 00112 (.../dead_code_elimination_test.go:23)        JMP     24
0x0074 00116 (.../dead_code_elimination_test.go:26)        ADD     $48, RSP
0x0078 00120 (.../dead_code_elimination_test.go:26)        SUB     $8, RSP, R29
0x007c 00124 (.../dead_code_elimination_test.go:26)        RET     (R30)
*/
