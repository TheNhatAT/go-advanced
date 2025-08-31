// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package thingsmustbeclosed

import (
	"net/http"

	"github.com/efficientgo/core/errcapture"
	"github.com/efficientgo/core/errors"
)

// Examples of code which is leaking resources, because of not exhausted readers.
// Read more in "Efficient Go"; Example 11-10.

func handleResp_Wrong(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return errors.Newf("got non-200 response; code: %v", resp.StatusCode)
	}
	// NOTE: leaks two goroutines:
	// 1. net/http.(*persistConn).writeLoop
	// 2. net/http.(*persistConn).readLoop
	// see BenchmarkClient's results in http_exhaust_test.go.
	return nil
}

func handleResp_StillWrong(resp *http.Response) error {
	defer func() {
		_ = resp.Body.Close()
	}()
	// NOTE: stop the leaks but still don't read bytes from the body
	// => net/http implementations can block the TCP connection if we don't fully exhaust the body
	// => may not be able to re-use a persistent TCP connection to the server for a subsequent 'keep-alive' request.
	// ref: https://pkg.go.dev/net/http#Client.Do
	if resp.StatusCode != http.StatusOK {
		return errors.Newf("got non-200 response; code: %v", resp.StatusCode)
	}
	return nil
}

func handleResp_Better(resp *http.Response) (err error) {
	defer errcapture.ExhaustClose(&err, resp.Body, "close")
	if resp.StatusCode != http.StatusOK {
		return errors.Newf("got non-200 response; code: %v", resp.StatusCode)
	}
	return nil
}
