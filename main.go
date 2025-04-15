package main

import (
	"fmt"

	"github.com/padiazg/nurdsoft-challenge/internals"
	"github.com/padiazg/nurdsoft-challenge/models"
)

func main() {
	data := internals.NewBookList()

	first, err := data.Add(&models.Book{
		Title:  "A book",
		Author: "Jhon Doe",
		Price:  10.3,
		ISBN:   "ABCDERFGHT",
	})

	data.Add(&models.Book{
		Title:  "Another book",
		Author: "Jhon Doe",
		Price:  10.5,
		ISBN:   "ABCDERFGHT123235423452",
	})

	one, err := data.GetOne(first)
	if err != nil {
		fmt.Printf("from GetOne: %+v", err)
	}
	fmt.Printf("one: %+v\n", one)

	all := data.GetAll()
	fmt.Printf("all: %+v\n", all)

	data.Update(first, &models.Book{
		Title:  "A book (revised)",
		Author: "Jhon Doe",
		Price:  10.8,
		ISBN:   "0987654321",
	})

	updt, err := data.GetOne(first)
	if err != nil {
		fmt.Printf("from GetOne: %+v", err)
	}
	fmt.Printf("updated: %+v\n", updt)
}
