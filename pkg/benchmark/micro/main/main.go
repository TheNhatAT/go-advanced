package main

import "go-advanced/pkg/benchmark/micro"

func main() {
	//_, _ = micro.Sum("./pkg/benchmark/micro/testdata/test.2000000.txt")
	_, _ = micro.Sum2("./pkg/benchmark/micro/testdata/test.2000000.txt")
}

/**
RSS of micro.Sum ~= 58MB
RSS of micro.Sum2 ~= 11MB
*/
