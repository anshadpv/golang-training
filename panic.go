package main

import "fmt"

func division(a int, b int) {
	defer panicHandle()

	if b == 0 {
		panic("Cannot divide by zero")
	} else {
		result := a / b
		fmt.Println(result)
	}
}

func panicHandle() {
	a := recover()
	if a != nil {
		fmt.Println("RECOVER : ", a)
	}
}

func main() {
	division(4, 2)
	division(3, 8)
	division(4, 0)
	division(6, 2)
}
