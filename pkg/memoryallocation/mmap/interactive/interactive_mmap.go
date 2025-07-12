// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/efficientgo/examples/pkg/memory/mmap"
)

// Creating memory mapping for 600MB of file.
// Read more in "Efficient Go"; Example 5-3.

func runMmap() {
	fmt.Println("PID", os.Getpid())

	// TODO(bwplotka): Create big file here, so we can play with it - there is no need to upload so big file to GitHub.

	// Mmap 600 MB of 686MB file.
	f, err := mmap.OpenFileBacked("test686mbfile.out", 600*1024*1024)
	if err != nil {
		log.Fatal(err)
	}

	// Check out:
	// ls -l /proc/<PID>/map_files
	// watch -n 1 'export PID=$(ps -ax --format=pid,command | grep "exe/interactive" | grep -v "grep" | head -n 1 | cut -d" " -f2) && ps -ax --format=pid,rss,vsz | grep $PID && cat /proc/$PID/smaps | grep -A22 test686mbfile | grep Rss'
	b := f.Bytes()

	fmt.Println("1")
	fmt.Scanln() // wait for Enter Key

	fmt.Println("Reading 5000 index", b[5000])

	fmt.Println("2")
	fmt.Scanln() // wait for Enter Key

	fmt.Println("Reading 100 000 index", b[100000])

	fmt.Println("3")
	fmt.Scanln() // wait for Enter Key

	fmt.Println("Reading 104 000 index", b[104000])

	fmt.Println("4")
	fmt.Scanln() // wait for Enter Key

	fmt.Println("Unmapping")
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Finish")
	fmt.Scanln() // wait for Enter Key

}

func bookExample1() error {
	// Mmap 600MB of 686MB file.
	f, err := mmap.OpenFileBacked("test686mbfile.out", 600*1024*1024)
	if err != nil {
		return err
	}
	b := f.Bytes() // RSS ~= 4MB

	// At this point we can see symlink to test686mbfile.out file in /proc/<PID>/map_files.
	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 test686mbfile | grep Rss` shows 0KB.
	// NOTE exact RSS can vary due to state of memory, OS, OS version etc.
	fmt.Println("Reading 5000 index", b[5000]) // RSS ~= 4MB

	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 test686mbfile | grep Rss` shows 48-70KB.

	fmt.Println("Reading 100 000 index", b[100000]) // RSS ~= 4MB

	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 test686mbfile | grep Rss` shows 100-126KB.

	fmt.Println("Reading 104 000 index", b[104000]) // RSS ~= 4MB

	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 test686mbfile | grep Rss` shows 100-126KB (same).

	if err := f.Close(); err != nil {
		return err
	}

	// If we would pause the program now `cat /proc/<PID>/smaps | grep -A22 test686mbfile | grep Rss` shows nothing, RSS freed.
	return nil // still RSS ~= 4MB
}

func main_3() {
	pid := os.Getpid()
	fmt.Printf("Current process PID: %d\n", pid)
	_ = bookExample1()
	fmt.Println("Run open example")
}
