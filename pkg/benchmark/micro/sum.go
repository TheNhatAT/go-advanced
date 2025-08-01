package micro

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

// Sum is a naive implementation and algorithm for summing integers from file.
// Read more in "Efficient Go"; Example 4-1.
// using "ps -p <PID> -o pid,rss,vsz to check RSS and VSZ memory usage.
func Sum(fileName string) (ret int64, _ error) {
	fmt.Println("PID", os.Getpid())
	b, err := os.ReadFile(fileName) // RSS ~= 3MB
	if err != nil {
		return 0, err
	} // RSS ~= 11MB
	for _, line := range bytes.Split(b, []byte("\n")) { // RSS ~= 59MB
		if len(line) == 0 {
			// Empty line at the end.
			continue
		}

		num, err := strconv.ParseInt(string(line), 10, 64)
		if err != nil {
			return 0, err
		}

		ret += num
	}
	return ret, nil // RSS ~= 59MB
}

// Sum2 is sum with optimized the first latency + CPU bottleneck bytes.Split.
// bytes.Split look complex to hande different cases. It allocates a lot causing  It looks like the algo is simple enough to just
// implement on our own (tried scanner := bufio.NewScanner(f) but it's slower).
// 30% less latency and 5x less memory than Sum.
// Read more in "Efficient Go"; Example 10-3.
func Sum2(fileName string) (ret int64, _ error) {
	fmt.Println("PID", os.Getpid())
	b, err := os.ReadFile(fileName) // RSS ~= 3MB
	if err != nil {
		return 0, err
	}

	var last int
	for i := 0; i < len(b); i++ { // RSS ~= 11MB
		if b[i] != '\n' {
			continue
		}
		num, err := strconv.ParseInt(string(b[last:i]), 10, 64)
		if err != nil {
			return 0, err
		}

		ret += num
		last = i + 1
	}
	return ret, nil // RSS ~= 11MB
}
