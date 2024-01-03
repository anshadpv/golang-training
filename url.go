package main

import (
	"fmt"
	"net/url"
)

const myUrl = "https://lco.dev:3000/learn?coursename=reactjs&paymentid=ghbj456ghb"

func main() {
	fmt.Println("Handling URLs")
	fmt.Println(myUrl)
	result, err := url.Parse(myUrl)
	if err != nil {
		panic(err)
	}
	// Printing the parsed URL details
	fmt.Println(result.Scheme)
	fmt.Println(result.Host)
	fmt.Println(result.Path)
	fmt.Println(result.Port())
	fmt.Println(result.RawQuery)

	qparams := result.Query()
	fmt.Println("\nParameters in Query : ")
	for key, values := range qparams {
		fmt.Printf("Key:%s \t Values:%v\n", key, values)
	}

	pathsOfUrl := &url.URL{
		Scheme:   "http",
		Host:     "www.example.com",
		Path:     "/test",
		RawQuery: "Key=Value",
	}
	fmt.Println("Path : ", pathsOfUrl.String())
	fmt.Println("Full Path : ", pathsOfUrl.RequestURI())
}
