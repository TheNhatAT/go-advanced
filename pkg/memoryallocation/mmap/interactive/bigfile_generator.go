package main

import (
	"crypto/rand"
	"fmt"
	"os"
)

// generate test686mbfile.out file with 686MB of random data.
func main_1() {
	const (
		filename   = "test686mbfile.out"
		fileSize   = 686 * 1024 * 1024 // 686 MB in bytes
		bufferSize = 1024 * 1024       // 1 MB buffer
	)

	fmt.Printf("Generating %s with %d MB of random data...\n", filename, fileSize/(1024*1024))

	// Create the output file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("Error closing file: %v\n", closeErr)
		}
	}()

	// Create a buffer for writing random data
	buffer := make([]byte, bufferSize)
	bytesWritten := int64(0)

	for bytesWritten < fileSize {
		// Calculate how much to write in this iteration
		remaining := fileSize - bytesWritten
		writeSize := bufferSize
		if remaining < int64(bufferSize) {
			writeSize = int(remaining)
			buffer = buffer[:writeSize]
		}

		// Fill buffer with random data
		_, err := rand.Read(buffer)
		if err != nil {
			fmt.Printf("Error generating random data: %v\n", err)
			os.Exit(1)
		}

		// Write buffer to file
		n, err := file.Write(buffer)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			os.Exit(1)
		}

		bytesWritten += int64(n)

		// Show progress every 50MB
		if bytesWritten%(50*1024*1024) == 0 || bytesWritten == fileSize {
			progress := float64(bytesWritten) / float64(fileSize) * 100
			fmt.Printf("Progress: %.1f%% (%d MB written)\n", progress, bytesWritten/(1024*1024))
		}
	}

	// Ensure all data is written to disk
	err = file.Sync()
	if err != nil {
		fmt.Printf("Error syncing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s (%d bytes)\n", filename, bytesWritten)
}
