package lifecycleofgoroutines

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// Example case of leak in HTTP handlers.
// Read more in "Efficient Go"; Example 11-2.

func ComplexComputation() int {
	time.Sleep(1 * time.Second) // Computation.
	time.Sleep(1 * time.Second) // Cleanup.
	return 4
}

func Handle_VeryWrong(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int)

	go func() {
		defer close(respCh)
		respCh <- ComplexComputation()
	}()

	// Some other work...

	select {
	case <-r.Context().Done():
		/** NOTE: we only control the lifecycle only in a good case (when no cancellation occurs).
		 * If the request context is done, we return without caring about the above goroutine lifecycle
		 * => we don't stop it, don't wait for it => it can be a permanent leak - ComplexComputation() will be starved
		 * as no one reads from the respCh channel.
		 */
		return
	case resp := <-respCh:
		_, _ = w.Write([]byte(strconv.Itoa(resp)))
		return
	}
}

// More examples of leaks in HTTP handlers.
// Read more in "Efficient Go"; Example 11-5.

func Handle_Wrong(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int, 1) // NOTE: a channel with a buffer for one message -> allow computation goroutine
	// to push one message to the channel without waiting for someone to read it
	// => if we cancel and wait some time -> the "left behind" goroutine will eventually finish

	go func() {
		defer close(respCh)
		respCh <- ComplexComputation()
	}()

	// Some other work...

	select {
	case <-r.Context().Done():
		// NOTE: but the goroutine continues running ComplexComputation() even after the request is cancelled
		// => waiting unaccounted resource usage & hard to trace because the memory allocation is allocated after
		// the request is done.
		return
	case resp := <-respCh:
		_, _ = w.Write([]byte(strconv.Itoa(resp)))
		return
	}
}

func ComplexComputationWithCtx(ctx context.Context) (ret int) {
	select {
	case <-ctx.Done():
	case <-time.After(1 * time.Second): // Computation.
		ret = 4
	}

	// NOTE: even if we cancel the context, we still do some cleanup.
	// => still consuming resources (CPU time, memory, etc.) for cleanup work that may not be necessary.
	time.Sleep(1 * time.Second) // Cleanup.
	return ret
}

func Handle_AlsoWrong(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int, 1)

	go func() {
		defer close(respCh)
		respCh <- ComplexComputationWithCtx(r.Context()) // NOTE: accept context as a parameter
		// => cancels computation when no longer needed
	}()

	// Some other work...

	select {
	case <-r.Context().Done():
		return
	case resp := <-respCh:
		_, _ = w.Write([]byte(strconv.Itoa(resp)))
		return
	}
}

// Recommended code that does not leak.
// Read more in "Efficient Go"; Example 11-6.

func Handle_Better(w http.ResponseWriter, r *http.Request) {
	respCh := make(chan int)

	go func() {
		defer close(respCh)
		respCh <- ComplexComputationWithCtx(r.Context())
	}()

	// Some other work...

	// NOTE: always reading from the channel -> wait for the goroutine stop
	// => respond to cancel as quickly as possible
	resp := <-respCh
	if r.Context().Err() != nil {
		return
	}

	_, _ = w.Write([]byte(strconv.Itoa(resp)))
}
