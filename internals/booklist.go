package internals

import (
	"fmt"
	"sync"

	"github.com/padiazg/nurdsoft-challenge/models"
)

type BookList struct {
	list  map[int32]*models.Book
	lock  sync.Mutex
	count int32
}

func NewBookList() *BookList {
	book_list := &BookList{}
	book_list.list = map[int32]*models.Book{}
	return book_list
}

func (bl *BookList) Add(data *models.Book) (int32, error) {
	// only two fields checked to trigger some error
	if data.Title == "" || data.Author == "" {
		return 0, fmt.Errorf("some field data missing")
	}

	bl.lock.Lock()
	defer bl.lock.Unlock()

	bl.count++
	data.ID = bl.count
	data.Active = true
	bl.list[bl.count] = data

	return bl.count, nil
}

func (bl *BookList) GetOne(id int32) (*models.Book, error) {
	if book, ok := bl.list[id]; ok {
		return book, nil
	}

	return nil, fmt.Errorf("%d not found", id)
}

func (bl *BookList) GetAll() []*models.Book {
	res := []*models.Book{}

	for _, book := range bl.list {
		if book.Active {
			res = append(res, book)
		}
	}

	return res
}

func (bl *BookList) Update(id int32, data *models.Book) (*models.Book, error) {
	book, err := bl.GetOne(id)
	if err != nil {
		return nil, err
	}

	bl.lock.Lock()
	defer bl.lock.Unlock()

	book.Title = data.Title
	book.Author = data.Author
	book.Price = data.Price
	book.ISBN = data.ISBN

	return book, nil
}

func (bl *BookList) Delete(id int32) error {
	book, err := bl.GetOne(id)
	if err != nil {
		return err
	}

	bl.lock.Lock()
	defer bl.lock.Unlock()
	book.Active = false

	return nil
}
