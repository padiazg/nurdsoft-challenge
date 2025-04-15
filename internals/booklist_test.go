package internals

import (
	"testing"

	"github.com/padiazg/nurdsoft-challenge/models"
	"github.com/stretchr/testify/assert"
)

func TestNewBookList(t *testing.T) {
	type CheckFn func(t *testing.T, bl *BookList)

	var (
		CheckList = func(fns ...CheckFn) []CheckFn { return fns }

		checkNotNull = func() CheckFn {
			return func(t *testing.T, bl *BookList) {
				t.Helper()
				assert.NotNil(t, bl)
			}
		}

		checkListInstantiated = func() CheckFn {
			return func(t *testing.T, bl *BookList) {
				t.Helper()
				assert.NotNil(t, bl.list)
			}
		}

		tests = []struct {
			name   string
			checks []CheckFn
		}{
			{
				name: "success",
				checks: CheckList(
					checkNotNull(),
					checkListInstantiated(),
				),
			},
		}
	)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := NewBookList()
			for _, check := range tt.checks {
				check(t, got)
			}
		})
	}
}

func TestBookList_Add(t *testing.T) {
	type CheckFn func(t *testing.T, id int32, bl *BookList)

	var (
		CheckList = func(fns ...CheckFn) []CheckFn { return fns }
		err       error

		hasError = func(want bool) CheckFn {
			return func(t *testing.T, _ int32, _ *BookList) {
				t.Helper()
				if want {
					assert.NotNil(t, err, "hasError: error expected, none produced")
				} else {
					assert.Nil(t, err, "hasError = [+%v], no error expected")
				}
			}
		}

		checkCount = func(q int) CheckFn {
			return func(t *testing.T, _ int32, bl *BookList) {
				t.Helper()
				assert.Equal(t, q, len(bl.list))
			}
		}

		expectedID = func(expected_id int32) CheckFn {
			return func(t *testing.T, id int32, _ *BookList) {
				t.Helper()
				assert.Equal(t, expected_id, id)
			}
		}

		tests = []struct {
			name   string
			book   *models.Book
			before func(bl *BookList)
			checks []CheckFn
		}{
			{
				name: "first-added",
				book: &models.Book{Title: "test-book-1", Author: "author"},
				checks: CheckList(
					hasError(false),
					checkCount(1),
					expectedID(1),
				),
			},
			{
				name: "second-added",
				book: &models.Book{Title: "test-book-1", Author: "author"},
				before: func(bl *BookList) {
					bl.list[1] = &models.Book{ID: 1}
					bl.count = 1
				},
				checks: CheckList(
					hasError(false),
					checkCount(2),
					expectedID(2),
				),
			},
			{
				name: "missing-fields",
				book: &models.Book{},
				checks: CheckList(
					hasError(true),
				),
			},
		}
	)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var (
				bl  = NewBookList()
				got int32
			)

			if tt.before != nil {
				tt.before(bl)
			}

			got, err = bl.Add(tt.book)
			for _, check := range tt.checks {
				check(t, got, bl)
			}
		})
	}
}

func TestBookList_GetOne(t *testing.T) {
	type CheckFn func(t *testing.T, book *models.Book, bl *BookList)

	var (
		CheckList = func(fns ...CheckFn) []CheckFn { return fns }
		err       error

		hasError = func(want bool) CheckFn {
			return func(t *testing.T, _ *models.Book, _ *BookList) {
				t.Helper()
				if want {
					assert.NotNil(t, err, "hasError: error expected, none produced")
				} else {
					assert.Nil(t, err, "hasError = [+%v], no error expected")
				}
			}
		}

		expectedBook = func(expected *models.Book) CheckFn {
			return func(t *testing.T, book *models.Book, bl *BookList) {
				assert.EqualValues(t, expected, book)
			}
		}

		tests = []struct {
			name   string
			id     int32
			before func(bl *BookList)
			checks []CheckFn
		}{
			{
				name: "get-first",
				id:   1,
				before: func(bl *BookList) {
					bl.list[1] = &models.Book{
						ID:     1,
						Title:  "test-book-one",
						Author: "test-author-one",
						Price:  10.0,
						ISBN:   "1234567890",
						Active: true,
					}
					bl.count = 1
				},
				checks: CheckList(
					hasError(false),
					expectedBook(&models.Book{
						ID:     1,
						Title:  "test-book-one",
						Author: "test-author-one",
						Price:  10.0,
						ISBN:   "1234567890",
						Active: true,
					}),
				),
			},
			{
				name: "get-second",
				id:   2,
				before: func(bl *BookList) {
					bl.list = map[int32]*models.Book{
						1: {
							ID:     1,
							Title:  "test-book-one",
							Author: "test-author-one",
							Price:  10.0,
							ISBN:   "1234567890",
							Active: true,
						},
						2: {
							ID:     2,
							Title:  "test-book-two",
							Author: "test-author-two",
							Price:  11.0,
							ISBN:   "123456789222",
							Active: true,
						},
					}

					bl.count = int32(len(bl.list))
				},
				checks: CheckList(
					hasError(false),
					expectedBook(&models.Book{
						ID:     2,
						Title:  "test-book-two",
						Author: "test-author-two",
						Price:  11.0,
						ISBN:   "123456789222",
						Active: true,
					}),
				),
			},
			{
				name: "not-found",
				id:   3,
				before: func(bl *BookList) {
					bl.list = map[int32]*models.Book{}
					bl.count = int32(0)
				},
				checks: CheckList(
					hasError(true),
				),
			},
		}
	)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var (
				bl  = NewBookList()
				got *models.Book
			)

			if tt.before != nil {
				tt.before(bl)
			}

			got, err = bl.GetOne(tt.id)
			for _, check := range tt.checks {
				check(t, got, bl)
			}
		})
	}
}
