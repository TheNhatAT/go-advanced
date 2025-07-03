package _interface

// Basic interface and implementation
type getter interface {
	get() int
}

type zero struct{}

// using go:noinline to prevent inlining of the method
// for benchmarking purposes
//
//go:noinline
func (z zero) get() int {
	return 0
}

// io.Reader implementation that fills buffer with zeros
type zeroReader struct{}

func (z zeroReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}
