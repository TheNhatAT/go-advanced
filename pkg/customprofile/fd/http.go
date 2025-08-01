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
	// (1) allow registering HTTP server handlers on specific HTTP paths
	// importing `_ "net/http/pprof"` will register standard profiles in the default global mux (http.DefaultServeMux) by default
	// or manually like in this example.
	m := http.NewServeMux()

	// (2) exposes a root HTML index page that lists quick statistics and links to profilers registered using pprof.NewProfile.
	m.HandleFunc("/debug/pprof/", pprof.Index)

	// (3) standard Go CPU is not using pprof.Profile, have to register explicitly.
	m.HandleFunc("/debug/pprof/profile", pprof.Profile) // CPU profile

	// (4) used for third-party profilers, like fgprof, which profile the Off-CPU time.
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

/** example usecase of HTTP handlers:
(1) using `http://<address>/debug/pprof/heap?debug=1` for simple text html page
(2) using `http://<address>/debug/pprof/heap` for download profile
(3) using `go tool pprof -http :8080 http://<address>/debug/pprof/heap` for directly visualizing without download profile
(4) using another server to collect those profiles to a dedicated database periodically

(1) using `http://localhost:8081/ui/top?si=alloc_space&g=lines` for showing in lines granularity
*/
