package main

import (
	"fmt"
)

func reciever(ch chan person) {
	defer fmt.Println("Value recieved")
	fmt.Println(<-ch)
}

type person struct {
	name string
	age  int
}

func main() {
	fmt.Println("Creating channel")
	ch := make(chan person) //buffered channel
	// var ch chan string
	go reciever(ch)
	ch <- person{"ansh", 22}
	fmt.Println("Finished")
}
