package main

func passByValue(a [5]int) {
	a[0] = 100     // This won't affect the original array
	println(&a[0]) // Print address of the local copy
}

func main_3() {
	// Example usage of passByValue
	arr := [5]int{1, 2, 3, 4, 5}
	passByValue(arr)
	// arr remains unchanged: [1, 2, 3, 4, 5]
	println(arr[0])  // Output: 1
	println(&arr[0]) // Print address of the original array
}
