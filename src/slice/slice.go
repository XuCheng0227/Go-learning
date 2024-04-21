package main

import "fmt"

func Pic(dx, dy int) [][]uint8 {
	s := make([][]uint8, dx)
	for i := range s {
		s[i] = make([]uint8, dy)
	}

	for x, slice := range s {
		for y := range slice {
			s[x][y] = uint8((x + y) / 2)
		}
	}
	return s
}

func main() {
	a := Pic(6, 5)
	a[3][2] = 100
	for x, slice := range a {
		for y := range slice {
			fmt.Printf("The (%d, %d) value is: %d ", x, y, a[x][y])
		}
		fmt.Printf("\n")
	}
	fmt.Print(a[5][4])
	fmt.Print(a)

}
