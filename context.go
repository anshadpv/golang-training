package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	ctxtimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	go doSomething(ctxtimeout, ch)

	select {
	case <-ctxtimeout.Done():
		fmt.Println("Context Cancelled : ", ctxtimeout.Err())
	case result := <-ch:
		fmt.Println("Recieved : %s\n", result)
	}
}

func doSomething(ctx context.Context, ch chan string) {
	fmt.Println("Do something...")
	time.Sleep(5 * time.Second)
	fmt.Println("Wake Up!!")
	ch <- "Did Something"

}
