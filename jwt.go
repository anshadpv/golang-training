package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
	UserId int
}

const secretkey = "ansh"
const tokenexp = int64(time.Hour * 3)

func main() {

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenexp,
		},
		UserId: 66,
	})

	token, err := t.SignedString([]byte(secretkey))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(token)

	fmt.Println(getUserId(token))

	fmt.Println(tokenisValid(token))

}

func getUserId(tokenstr string) int {
	claims := &Claims{}
	jwt.ParseWithClaims(tokenstr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretkey), nil
	})
	return claims.UserId
}

func tokenisValid(tokenstr string) bool {
	token, err := jwt.Parse(tokenstr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretkey), nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}
