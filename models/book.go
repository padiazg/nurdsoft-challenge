package models

type Book struct {
	ID     int32
	Title  string
	Author string
	Price  float32
	ISBN   string
	Active bool
}
