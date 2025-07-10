package concurrency

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate API response time
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "success"}`))
	}))
}

func BenchmarkCallAPIsSequentially(b *testing.B) {
	server := setupMockServer()
	defer server.Close()

	caller := NewConcurrentAPICaller(5 * time.Second)
	urls := make([]string, 10)
	for i := range urls {
		urls[i] = fmt.Sprintf("%s/api/%d", server.URL, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		caller.CallAPIsSequentially(ctx, urls)
	}
}

func BenchmarkCallAPIsConcurrently(b *testing.B) {
	server := setupMockServer()
	defer server.Close()

	caller := NewConcurrentAPICaller(5 * time.Second)
	urls := make([]string, 10)
	for i := range urls {
		urls[i] = fmt.Sprintf("%s/api/%d", server.URL, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		caller.CallAPIsConcurrently(ctx, urls)
	}
}

func BenchmarkCallAPIsConcurrentlyWithWorkerPool(b *testing.B) {
	server := setupMockServer()
	defer server.Close()

	caller := NewConcurrentAPICaller(5 * time.Second)
	urls := make([]string, 10)
	for i := range urls {
		urls[i] = fmt.Sprintf("%s/api/%d", server.URL, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		caller.CallAPIsConcurrentlyWithWorkerPool(ctx, urls, 5)
	}
}

// Benchmark with different URL counts
func BenchmarkSequential_5URLs(b *testing.B) {
	benchmarkWithURLCount(b, 5, "sequential")
}

func BenchmarkConcurrent_5URLs(b *testing.B) {
	benchmarkWithURLCount(b, 5, "concurrent")
}

func BenchmarkWorkerPool_5URLs(b *testing.B) {
	benchmarkWithURLCount(b, 5, "worker")
}

func BenchmarkSequential_20URLs(b *testing.B) {
	benchmarkWithURLCount(b, 20, "sequential")
}

func BenchmarkConcurrent_20URLs(b *testing.B) {
	benchmarkWithURLCount(b, 20, "concurrent")
}

func BenchmarkWorkerPool_20URLs(b *testing.B) {
	benchmarkWithURLCount(b, 20, "worker")
}

func BenchmarkSequential_100URLs(b *testing.B) {
	benchmarkWithURLCount(b, 100, "sequential")
}

func BenchmarkConcurrent_100URLs(b *testing.B) {
	benchmarkWithURLCount(b, 100, "concurrent")
}

func BenchmarkWorkerPool_100URLs(b *testing.B) {
	benchmarkWithURLCount(b, 100, "worker")
}
func benchmarkWithURLCount(b *testing.B, urlCount int, method string) {
	server := setupMockServer()
	defer server.Close()

	caller := NewConcurrentAPICaller(5 * time.Second)
	urls := make([]string, urlCount)
	for i := range urls {
		urls[i] = fmt.Sprintf("%s/api/%d", server.URL, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		switch method {
		case "sequential":
			caller.CallAPIsSequentially(ctx, urls)
		case "concurrent":
			caller.CallAPIsConcurrently(ctx, urls)
		case "worker":
			caller.CallAPIsConcurrentlyWithWorkerPool(ctx, urls, 5)
		}
	}
}
