package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Price  int    `json:"price"`
}

var albums = []album{
	{"1", "Dusk till dawn", "Zayn", 100},
	{"2", "Another love", "Tom Odell", 200},
	{"3", "Blank Space", "Taylor", 50},
}

func main() {
	router := gin.Default()
	//router.Use(tokenValid)
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbum)
	router.POST("/albums", addAlbum)
	router.Run(":8088")

}

func getAlbums(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, albums)
}

func addAlbum(context *gin.Context) {
	var newAlbum album
	if err := context.BindJSON(&newAlbum); err != nil {
		fmt.Println(err)
		return
	}
	albums = append(albums, newAlbum)
	context.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbum(context *gin.Context) {
	id := context.Param("id")
	album, err := getById(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"Message": "Album not found"})
		return
	}
	context.IndentedJSON(http.StatusOK, album)
}

func getById(id string) (*album, error) {
	for i, t := range albums {
		if t.ID == id {
			return &albums[i], nil
		}
	}
	return nil, errors.New("ID not found")
}

// func tokenValid(context *gin.Context) bool {
// 	token, err := jwt.Parse(context.GetHeader("authorization"), func(t *jwt.Token) (interface{}, error) {
// 		return []byte(secretkey), nil
// 	})
// 	if err != nil {
// 		return false
// 	}
// 	return token.Valid
// }
