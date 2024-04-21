package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0
	init_val := 0.0
	for math.Abs(init_val-z) > 0.01 {
		init_val = z
		z -= (z*z - x) / (2 * z)
		// fmt.Println(z)
		// fmt.Println(init_val)
		// fmt.Println(math.Abs(init_val - z))
	}

	return z
}

func main() {
	fmt.Println("Final Result")
	fmt.Println(Sqrt(2))

}
