package main

func add(a, b int) int {
	return a + b
}

func main() {
	result := add(5, 10)
	_ = result
}

// using "go build -gcflags="-S" sample.go > assembly.txt 2>&1" command to generate assembly code
// using "GOSSAFUNC=main go build" command to generate SSA form
// using go build -gcflags="-m" sample.go to check for inlining and escape analysis
