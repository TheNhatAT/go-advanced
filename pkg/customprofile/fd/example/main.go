// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package main

import (
	"go-advanced/pkg/customprofile/fd"
	"io"
	"log"
	"sync"
)

// Example application instrumented with custom fd customprofile from Example 9-1.
// Read more in "Efficient Go"; Example 9-2.

type TestApp struct {
	files []io.ReadCloser
}

func (a *TestApp) Close() {
	for _, cl := range a.files {
		// (2) ensured file will be closed and removed from the profile.
		_ = cl.Close() // TODO: Check error.
	}
	a.files = a.files[:0]
}

func (a *TestApp) open(name string) {
	// (1) open file with fd.Open func with start recording
	// it in profile as a side effect.
	f, _ := fd.Open(name) // TODO: Check error.
	a.files = append(a.files, f)
}

func (a *TestApp) OpenSingleFile(name string) {
	a.open(name)
}

func (a *TestApp) OpenTenFiles(name string) {
	for i := 0; i < 10; i++ {
		a.open(name)
	}
}

func (a *TestApp) Open100FilesConcurrently(name string) {
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			a.OpenTenFiles(name)
			wg.Done()
		}()
	}
	wg.Wait()
}

func main() {
	a := &TestApp{}
	defer a.Close()

	// No matter how many files we opened in the past...
	for i := 0; i < 10; i++ {
		// (3) first open ten files closing them, repeated 10 times.
		a.OpenTenFiles("/dev/null")
		a.Close()
	}

	// ...after last close, only files below will be used in customprofile.
	f, _ := fd.Open("/dev/null") // TODO: Check error.
	a.files = append(a.files, f)

	a.OpenSingleFile("/dev/null")
	a.OpenTenFiles("/dev/null")
	a.Open100FilesConcurrently("/dev/null")

	// (4) take a snapshot of the situation in form of fd.inuse profile.
	if err := fd.Write("fd.pprof"); err != nil {
		log.Fatal(err)
	}
}

/**
- using `go tool pprof -lines -http :8080 fd.pprof`  or  `http://localhost:8080/ui/?g=lines` for visualizing the profile in lines granularity.
- using `go tool pprof -http :8080 fd.pprof` or  `http://localhost:8080/ui/?g=functions` for visualizing the profile in functions granularity.
*/
