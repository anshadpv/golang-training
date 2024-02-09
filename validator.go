package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Person struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=50"`
	Phone string `validate:"required,len"`
}

var validate *validator.Validate

func PhoneLen(fl validator.FieldLevel) bool {
	if fl.Field().Len() == 10 {
		return true
	}
	return false
}

func main() {
	person := Person{Name: "Anshad", Email: "anshadpv@gmail.com", Age: 22, Phone: "8157005774"}
	validate := validator.New()
	validate.RegisterValidation("len", PhoneLen)
	err := validate.Struct(person)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Person is valid")
}
