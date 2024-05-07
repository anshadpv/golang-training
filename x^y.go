package main

import (
	"fmt"
	"math"
)

func isExponential(num float64) bool {
	if num == 0 || num == 1 {
		return false
	}
	if num > 0 && num < 1 {
		for num < 1 {
			num *= 10
		}
	}

	var x float64
	for x = 2; x <= math.Sqrt(num); x++ {
		k := num
		for y := 1; k > 1; y++ {
			k /= x
			if k == 1 {
				return true
			}
		}

	}
	return false

}

func main() {
	// Test cases
	numbers := []float64{0, 0.04, 0.09, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 1000, 49, 64, 81, 10000}
	for _, num := range numbers {
		if isExponential(num) {
			fmt.Printf("%f can be expressed in exponential form\n", num)
		} else {
			fmt.Printf("%f cannot be expressed in exponential form\n", num)
		}
	}
}
