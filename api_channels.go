package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Album struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Price  int    `json:"price"`
}

func fetchDataFromAPI(url string, resultChannel chan<- []Album) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error fetching data:", err)
		close(resultChannel)
		return
	}
	defer resp.Body.Close()

	var data []Album
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		close(resultChannel)
		return
	}

	resultChannel <- data
	close(resultChannel)
}

func main() {

	apiURL := "http://localhost:8088/albums"

	resultChannel := make(chan []Album)

	go fetchDataFromAPI(apiURL, resultChannel)

	for albums := range resultChannel {
		for _, album := range albums {
			fmt.Printf("Album:\nID: %s\nTitle: %s\nArtist: %s\nPrice: %d\n",
				album.ID, album.Title, album.Artist, album.Price)
		}
	}
}
