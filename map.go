package main

import "fmt"

func IntMix(a, b int) int {
	if a+b < 10 {
		return a + b
	}
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Print(IntMix(1, 2))
}
