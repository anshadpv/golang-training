package main

import (
	"fmt"
	"time"

	"github.com/bluele/gcache"
)

func main() {

	gc := gcache.New(10).Expiration(3 * time.Second).LRU().Build()

	gc.Set("key 1", "value 1")
	val, err := gc.Get("key 1")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	gc.SetWithExpire("key 2", "value 2", time.Second*8)

	exist := gc.Has("key 1")
	fmt.Println("Key exist : ", exist)

	<-time.After(4 * time.Second)
	size := gc.Len(false)
	fmt.Println("Size : ", size)

	gc.Set("key 3", "value 3")
	gc.Set("key 4", "value 4")

	keys := gc.Keys(true)

	for _, key := range keys {
		val, err := gc.Get(key)
		if err != nil {
			panic(err)
		}
		fmt.Println(val)
	}

	gc.Purge()
	size = gc.Len(true)
	fmt.Println("Size : ", size)

}
