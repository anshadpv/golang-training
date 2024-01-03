package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Redis:", pong)

	err = client.Set("example-key", "Hello, Redis!", 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	val, err := client.Get("example-key").Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Value:", val)

	// Use Redis list to push and pop elements
	err = client.LPush("example-list", "item1", "item2", "item3").Err()
	if err != nil {
		log.Fatal(err)
	}

	listVal, err := client.LRange("example-list", 0, 2).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("List:", listVal)

	err = client.FlushDB().Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Redis database flushed successfully.")

	// Close the connection when done
	err = client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
