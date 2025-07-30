// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package fd

import (
	"os"
	"runtime/pprof"
)

// TODO:NOTE: (1) register profiles with name (unique) - fd.inuse name to indicate that the profile
// tracks in-use file descriptors.
var fdProfile = pprof.NewProfile("fd.inuse")

// File is a wrapper on os.File that tracks file descriptor lifetime.
type File struct {
	*os.File
}

// Open opens file and tracks it in the `fd` customprofile`.
// NOTE(bwplotka): We could use finalizers here, but explicit Close is more reliable and accurate.
// Unfortunately it also changes type which might be dropped accidentally.
func Open(name string) (*File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	//TODO:NOTE: (2) Add method records the object with second argument
	// tells how many calls to skip in the stack trace.
	fdProfile.Add(f, 2)
	return &File{File: f}, nil
}

// Close closes files and updates customprofile.
func (f *File) Close() error {
	//TODO:NOTE: (3) remote the object when the file is closed
	// using the same inner *os.File -> pprof package can track and find
	// the object that opened
	defer fdProfile.Remove(f.File)
	return f.File.Close()
}

// Write saves the customprofile of the currently open file descriptors in to file in pprof format.
func Write(profileOutPath string) error {
	out, err := os.Create(profileOutPath) // For simplicity, we don't include this file in customprofile.
	if err != nil {
		return err
	}
	//TODO:NOTE: (4) writes bytes of a full pprof file into a provided writer
	if err := fdProfile.WriteTo(out, 0); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func CreateTemp(dir, pattern string) (*File, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return nil, err
	}
	fdProfile.Add(f, 2)
	return &File{File: f}, nil
}
