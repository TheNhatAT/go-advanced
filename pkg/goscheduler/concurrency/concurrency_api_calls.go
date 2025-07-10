package concurrency

import (
	"context"
	"io"
	"net/http"
	"runtime"
	"sync"
	"time"
)

func setupGoMaxProcs() {
	// Set GOMAXPROCS to 1 to test I/O bound tasks
	// this is the variable to define how many virtual CPUs  (1-1 with OS threads) can execute Go code simultaneously.
	runtime.GOMAXPROCS(1)

	// when increasing the GOMAXPROCS, the performance of concurrently calling APIs does not improve
	// because the bottleneck is the I/O operations, not CPU processing.
	// but when running concurrently, we improve significantly the performance of it
	// because we can make multiple I/O requests in parallel (the hardware can handle multiple network requests at once)
}

// APIResponse represents the response from an API call
type APIResponse struct {
	URL      string
	Status   int
	Body     string
	Duration time.Duration
	Error    error
}

// ConcurrentAPICaller handles concurrent API calls
type ConcurrentAPICaller struct {
	client  *http.Client
	timeout time.Duration
}

// NewConcurrentAPICaller creates a new API caller with timeout
func NewConcurrentAPICaller(timeout time.Duration) *ConcurrentAPICaller {
	return &ConcurrentAPICaller{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// CallAPI makes a single API call
func (c *ConcurrentAPICaller) CallAPI(ctx context.Context, url string) APIResponse {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return APIResponse{
			URL:      url,
			Duration: time.Since(start),
			Error:    err,
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return APIResponse{
			URL:      url,
			Duration: time.Since(start),
			Error:    err,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{
			URL:      url,
			Status:   resp.StatusCode,
			Duration: time.Since(start),
			Error:    err,
		}
	}

	return APIResponse{
		URL:      url,
		Status:   resp.StatusCode,
		Body:     string(body),
		Duration: time.Since(start),
	}
}

// CallAPIsSequentially makes API calls one by one
func (c *ConcurrentAPICaller) CallAPIsSequentially(ctx context.Context, urls []string) []APIResponse {
	setupGoMaxProcs()
	var responses []APIResponse
	for _, url := range urls {
		response := c.CallAPI(ctx, url)
		responses = append(responses, response)
	}

	return responses
}

// CallAPIsConcurrently makes API calls concurrently using goroutines
func (c *ConcurrentAPICaller) CallAPIsConcurrently(ctx context.Context, urls []string) []APIResponse {
	setupGoMaxProcs()
	responses := make([]APIResponse, len(urls))
	var wg sync.WaitGroup

	for i, url := range urls {
		wg.Add(1)
		go func(index int, apiURL string) {
			defer wg.Done()
			responses[index] = c.CallAPI(ctx, apiURL)
		}(i, url)
	}

	wg.Wait()
	return responses
}

// CallAPIsConcurrentlyWithWorkerPool uses a worker pool pattern
func (c *ConcurrentAPICaller) CallAPIsConcurrentlyWithWorkerPool(ctx context.Context, urls []string, numWorkers int) []APIResponse {
	setupGoMaxProcs()
	urlChan := make(chan string, len(urls))
	resultChan := make(chan APIResponse, len(urls))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urlChan {
				select {
				case <-ctx.Done():
					return
				default:
					response := c.CallAPI(ctx, url)
					resultChan <- response
				}
			}
		}()
	}

	// Send URLs to workers
	go func() {
		defer close(urlChan)
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return
			case urlChan <- url:
			}
		}
	}()

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var responses []APIResponse
	for response := range resultChan {
		responses = append(responses, response)
	}

	return responses
}

// go test -bench=. -benchmem
