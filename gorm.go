package main

import "gorm.io/gorm"

type User struct {
	ID   uint
	Name string
}

func main() {
	db, err := gorm.Open("mysql", "root:msf@12345@/class?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user User
	db.Where("name = ?", "John").First(&user)
}
