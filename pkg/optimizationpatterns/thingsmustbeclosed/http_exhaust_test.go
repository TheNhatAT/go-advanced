// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package thingsmustbeclosed

import (
	"net/http"
	"testing"

	"github.com/efficientgo/core/testutil"
	"go.uber.org/goleak"
)

func BenchmarkClient(b *testing.B) {
	defer goleak.VerifyNone(
		b,
		goleak.IgnoreTopFunction("testing.(*B).run1"),
		goleak.IgnoreTopFunction("testing.(*B).doBench"),
	)
	c := &http.Client{}

	// NOTE: client runs some goroutines for each TCP connection we want to keep alive and reuse.
	// => call CloseIdleConnections for closing them to detect any leaks of code might introduce.
	defer c.CloseIdleConnections()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := c.Get("http://google.com")
		testutil.Ok(b, err)
		testutil.Ok(b, handleResp_Wrong(resp))
	}
}

/** Result:
BenchmarkClient
    http_exhaust_test.go:28: found unexpected goroutines:
        [Goroutine 35 in state select, with net/http.(*persistConn).writeLoop on top of the stack:
        net/http.(*persistConn).writeLoop(0x140000ee900)
        	/Users/g2-nhat.nguyen-dev/sdk/go1.20/src/net/http/transport.go:2410 +0x9c
        created by net/http.(*Transport).dialConn
        	/Users/g2-nhat.nguyen-dev/sdk/go1.20/src/net/http/transport.go:1766 +0x119c
         Goroutine 34 in state select, with net/http.(*persistConn).readLoop on top of the stack:
        net/http.(*persistConn).readLoop(0x140000ee900)
        	/Users/g2-nhat.nguyen-dev/sdk/go1.20/src/net/http/transport.go:2227 +0xba0
        created by net/http.(*Transport).dialConn
        	/Users/g2-nhat.nguyen-dev/sdk/go1.20/src/net/http/transport.go:1765 +0x1154
        ]
--- FAIL: BenchmarkClient

*/
