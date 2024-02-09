package main

import (
	"fmt"
)

func main() {
	var res = [...]int{1, 2, 3, 4, 5, 6, 7, 8} //array
	fmt.Println(len(res))
	fmt.Println(res)
	var res2 = [3]string{}
	fmt.Println(res2)
	res2[0] = "h"
	res2[1] = "e"
	res2[2] = "y"
	fmt.Println(res2)

	myslice := []int{10, 20, 30} //slice
	myslice = append(myslice, 40, 50)
	myslice2 := []int{100, 200, 300}
	myslice3 := append(myslice, myslice2...)
	fmt.Println(myslice)
	fmt.Println(myslice3)

}
