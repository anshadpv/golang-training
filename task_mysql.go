package main

import (
	"fmt"
)

type mySQL struct{}

func (m *mySQL) Insert(id int, data string) error {
	for i, _ := range maap {
		if id == i {
			fmt.Println("mySQL : ID already taken.")
			return nil
		}
	}
	maap[id] = data
	fmt.Println("mySQL : Data Inserted.")
	return nil
}

func (m *mySQL) Update(id int, data string) error {
	for i, _ := range maap {
		if id == i {
			maap[id] = data
			fmt.Println("mySQL : Data Updated.")
			return nil
		}
	}
	fmt.Println("mySQL : No data with this id.")
	return nil
}
