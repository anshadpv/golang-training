package main

import (
	"fmt"
)

type msSQL struct{}

func (m *msSQL) Insert(id int, data string) error {
	for i, _ := range maap {
		if id == i {
			fmt.Println("msSQL : ID already taken.")
			return nil
		}
	}
	maap[id] = data
	fmt.Println("msSQL : Data Inserted.")
	return nil
}

func (m *msSQL) Update(id int, data string) error {
	for i, _ := range maap {
		if id == i {
			fmt.Println("msSQL : Data Updated.")
			return nil
		}
	}
	fmt.Println("msSQL : No data with this id.")
	return nil
}
