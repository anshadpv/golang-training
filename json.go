package main

import (
	"encoding/json"
	"fmt"
)

type person struct {
	Name   string
	Age    int
	Place  string `json:"Location"`
	Gender string
}

func main() {
	p1 := []person{
		{"Ansh", 22, "Calicut", "M"},
		{"John", 23, "America", "M"},
		{"Katheirne", 24, "Africa", "F"},
	}

	jsonData := []byte(`
	{
		"Name" : "Ansh",
		"Age" : 22,
		"Location" : "Calicut",
		"Gender" : "M"

	}
	`)

	enco, err := json.MarshalIndent(p1, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", enco)

	check := json.Valid(jsonData)
	var d person

	if check {
		fmt.Println("JSON is valid")
		json.Unmarshal(jsonData, &d)
		fmt.Printf("%#v", d)
	} else {
		fmt.Println("JSON was not valid")
	}

}
