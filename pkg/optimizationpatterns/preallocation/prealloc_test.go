package preallocation

import (
	"bytes"
	"github.com/efficientgo/core/testutil"
	"io"
	"strings"
	"testing"
)

func withoutPrealloc(b testutil.TB) {
	const size = 1e6

	var slice []string
	for i := 0; i < size; i++ {
		slice = append(slice, "something")
	}

	var slice2 []string
	for i := 0; i < size; i++ {
		slice2 = append(slice2, "something")
	}

	m := make(map[int]string)
	for i := 0; i < size; i++ {
		m[i] = "something"
	}

	buf := bytes.Buffer{}
	for i := 0; i < size; i++ {
		_ = buf.WriteByte('a')
	}

	builder := strings.Builder{}
	for i := 0; i < size; i++ {
		builder.WriteByte('a')
	}

	buf2, _ := io.ReadAll(bytes.NewReader(make([]byte, size)))
	buf3, _ := io.ReadAll(bytes.NewReader(make([]byte, size)))

	// .Test
	if !b.IsBenchmark() {
		testutil.Equals(b, slice, slice2)
		testutil.Equals(b, buf.String(), builder.String())
		testutil.Equals(b, buf2, buf3)
		testutil.Equals(b, buf2, make([]byte, size))
	}
}

// Examples of allocations that can be pre-allocated.
// Read more in "Efficient Go"; Example 11-11.
func withPrealloc(b testutil.TB) {
	const size = 1e6

	slice := make([]string, 0, size)
	for i := 0; i < size; i++ {
		slice = append(slice, "something")
	}

	slice2 := make([]string, size)
	for i := 0; i < size; i++ {
		slice2[i] = "something"
	}

	m := make(map[int]string, size)
	for i := 0; i < size; i++ {
		m[i] = "something"
	}

	buf := bytes.Buffer{}
	buf.Grow(size)
	for i := 0; i < size; i++ {
		_ = buf.WriteByte('a')
	}

	builder := strings.Builder{}
	builder.Grow(size)
	for i := 0; i < size; i++ {
		builder.WriteByte('a')
	}

	buf2, _ := ReadAll1(bytes.NewReader(make([]byte, size)), size)
	buf3, _ := ReadAll2(bytes.NewReader(make([]byte, size)), size)

	// .Test
	if !b.IsBenchmark() {
		testutil.Equals(b, slice, slice2)
		testutil.Equals(b, buf.String(), builder.String())
		testutil.Equals(b, buf2, buf3)
		testutil.Equals(b, buf2, make([]byte, size))
	}
}

// Examples of pre-allocations for standard library helper `io.ReadAll`.
// Read more in "Efficient Go"; Example 11-12.

func ReadAll1(r io.Reader, size int) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Grow(size)
	n, err := io.Copy(&buf, r)
	return buf.Bytes()[:n], err
}

func ReadAll2(r io.Reader, size int) ([]byte, error) {
	buf := make([]byte, size)
	n, err := io.ReadFull(r, buf)
	if err == io.EOF {
		err = nil
	}
	return buf[:n], err
}

func TestAllocs(t *testing.T) {
	withoutPrealloc(testutil.NewTB(t))
	withPrealloc(testutil.NewTB(t))
}

func BenchmarkAllocs(b *testing.B) {
	b.Run("WithoutPrealloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			withoutPrealloc(testutil.NewTB(b))
		}
	})
	b.Run("WithPrealloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			withPrealloc(testutil.NewTB(b))
		}
	})
}
func BenchmarkReadAlls(b *testing.B) {
	const size = int(1e6)
	inner := make([]byte, size)
	b.Run("io.ReadAll", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf, err := io.ReadAll(bytes.NewReader(inner))
			testutil.Ok(b, err)
			testutil.Equals(b, size, len(buf))
		}
	})
	b.Run("ReadAll1", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf, err := ReadAll1(bytes.NewReader(inner), size)
			testutil.Ok(b, err)
			testutil.Equals(b, size, len(buf))
		}
	})
	b.Run("ReadAll2", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			buf, err := ReadAll2(bytes.NewReader(inner), size)
			testutil.Ok(b, err)
			testutil.Equals(b, size, len(buf))
		}
	})
}

/** BenchmarkReadAlls results:
goos: darwin
goarch: arm64
pkg: github.com/efficientgo/examples/pkg/prealloc

BenchmarkReadAlls
BenchmarkReadAlls/io.ReadAll
BenchmarkReadAlls/io.ReadAll-10         	    3158	    382974 ns/op	 5241172 B/op	      29 allocs/op
BenchmarkReadAlls/ReadAll1
BenchmarkReadAlls/ReadAll1-10           	   18223	     63760 ns/op	 1007724 B/op	       4 allocs/op
BenchmarkReadAlls/ReadAll2
BenchmarkReadAlls/ReadAll2-10           	   18309	     63142 ns/op	 1007678 B/op	       3 allocs/op
*/
