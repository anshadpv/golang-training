package main

import (
	"fmt"

	"github.com/jenazads/gotry"
)

func main() {
	var obj interface{}
	obj = 2

	gotry.Try(func() {
		text := obj.(string)
		fmt.Println("Try ---> ", text)
	}).Catch(func(e gotry.Exception) {
		fmt.Println("Catch ---> exception catched #1:", e)

		gotry.Try(func() {
			gotry.Throw("New exception here")
		}).Catch(func(e gotry.Exception) {
			fmt.Println("Catch ---> exception catched #2:", e)

		}).Finally(func() {
			fmt.Println("Finally ---> This always print after all try block #2")
		})
	}).Finally(func() {
		fmt.Println("Finally ---> This always print after all try block #1")
	})
}
