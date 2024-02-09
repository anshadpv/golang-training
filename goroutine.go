package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("GoRoutine example")
	var wg sync.WaitGroup
	wg.Add(2)
	go test1(&wg)
	go test2(&wg)
	wg.Wait()

}

func test1(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("In the first thread")
}
func test2(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("In the second thread")
}
