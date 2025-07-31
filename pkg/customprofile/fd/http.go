// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package fd

import (
	"errors"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/felixge/fgprof"
)

// ExampleHTTP is an example of HTTP exposure of Go profiles.
// Read more in "Efficient Go"; Example 9-5.
func ExampleHTTP() {
	m := http.NewServeMux()
	m.HandleFunc("/debug/pprof/", pprof.Index)
	m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.HandleFunc("/debug/fgprof/profile", fgprof.Handler().ServeHTTP)
	m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	srv := http.Server{Handler: m}

	// Start server with port 8080.
	srv.Addr = ":8080"
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		// Handle error, e.g., log it.
		log.Fatalf("Failed to start server: %v", err)
	}
}
