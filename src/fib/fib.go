package main

import "fmt"

// 返回一个“返回int的函数”
func fibonacci() func() int {
	x1 := 0
	x2 := 0
	sum := 0
	return func() int {
		if x1+x2 == 0 {
			x2 = 1
			return 0
		} else if x1+x2 == 1 && sum == 0 {
			x1 = 1
			sum = 1
			return 1
		} else if sum == 1 {
			sum = 0
			return 1
		} else {
			sum = x1 + x2
			x1 = x2
			x2 = sum
			return sum
		}
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
