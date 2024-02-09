package main

import (
	"fmt"

	"github.com/samber/lo"
)

func main() {
	slice := []int{1, 2, 3, 4, 5, 2, 2}
	fmt.Println(lo.Count(slice, 2))
	fmt.Println(lo.Sum(slice))
	if lo.Contains(slice, 3) {
		fmt.Println("Slice contains 3")
	} else {
		fmt.Println("Slice does not contain 3")
	}
}
