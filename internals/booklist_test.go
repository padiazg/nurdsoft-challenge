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
				assert.NotNil(t, bl.List)
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
				assert.Equal(t, q, len(bl.List))
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
					bl.List[1] = &models.Book{ID: 1}
					bl.Count = 1
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
					bl.List[1] = &models.Book{
						ID:     1,
						Title:  "test-book-one",
						Author: "test-author-one",
						Price:  10.0,
						ISBN:   "1234567890",
						Active: true,
					}
					bl.Count = 1
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
					bl.List = map[int32]*models.Book{
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

					bl.Count = int32(len(bl.List))
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
					bl.List = map[int32]*models.Book{}
					bl.Count = int32(0)
				},
				checks: CheckList(
					hasError(true),
				),
			},
			{
				name: "deleted",
				id:   2,
				before: func(bl *BookList) {
					bl.List = map[int32]*models.Book{
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
							Active: false,
						},
					}

					bl.Count = int32(len(bl.List))
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

func TestBookList_GetAll(t *testing.T) {
	type CheckFn func(t *testing.T, l []*models.Book, bl *BookList)

	var (
		CheckList = func(fns ...CheckFn) []CheckFn { return fns }

		expected = func(list []*models.Book) CheckFn {
			return func(t *testing.T, l []*models.Book, bl *BookList) {
				t.Helper()
				len0 := len(list)
				len1 := len(l)
				if assert.Equalf(t, len0, len1, "count differ: expected %d, got %d", len0, len1) {
					for _, b0 := range list {
						found := false
						for _, b1 := range l {
							if b0.ID == b1.ID {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("expected ID=%d to be present, not found", b0.ID)
						}
					}
				}
			}
		}

		tests = []struct {
			name   string
			before func(bl *BookList)
			checks []CheckFn
		}{
			{
				name: "success",
				before: func(bl *BookList) {
					bl.List = map[int32]*models.Book{
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

					bl.Count = int32(len(bl.List))
				},
				checks: CheckList(
					expected([]*models.Book{
						{ID: 1, Title: "test-book-one", Author: "test-author-one", Price: 10.0, ISBN: "1234567890", Active: true},
						{ID: 2, Title: "test-book-two", Author: "test-author-two", Price: 11.0, ISBN: "123456789222", Active: true},
					}),
				),
			},
			{
				name: "deleted-items",
				before: func(bl *BookList) {
					bl.List = map[int32]*models.Book{
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
							Active: false,
						},
					}

					bl.Count = int32(len(bl.List))
				},
				checks: CheckList(
					expected([]*models.Book{
						{ID: 1, Title: "test-book-one", Author: "test-author-one", Price: 10.0, ISBN: "1234567890", Active: true},
					}),
				),
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				bl  = NewBookList()
				got []*models.Book
			)

			if tt.before != nil {
				tt.before(bl)
			}

			got = bl.GetAll()
			for _, check := range tt.checks {
				check(t, got, bl)
			}
		})
	}
}

func TestBookList_Update(t *testing.T) {
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
			data   *models.Book
			before func(bl *BookList)
			checks []CheckFn
		}{
			{
				name: "update-first",
				id:   1,
				before: func(bl *BookList) {
					bl.List[1] = &models.Book{
						ID:     1,
						Title:  "test-book-one",
						Author: "test-author-one",
						Price:  10.0,
						ISBN:   "1234567890",
						Active: true,
					}
					bl.Count = 1
				},
				data: &models.Book{
					ID:     1,
					Title:  "test-book-one (updated)",
					Author: "test-author-one",
					Price:  10.0,
					ISBN:   "1234567890",
					Active: true,
				},
				checks: CheckList(
					hasError(false),
					expectedBook(&models.Book{
						ID:     1,
						Title:  "test-book-one (updated)",
						Author: "test-author-one",
						Price:  10.0,
						ISBN:   "1234567890",
						Active: true,
					}),
				),
			},
			{
				name: "not-found",
				id:   1,
				checks: CheckList(
					hasError(true),
				),
			},
			{
				name: "missing-data",
				id:   1,
				before: func(bl *BookList) {
					bl.List[1] = &models.Book{
						ID:     1,
						Title:  "test-book-one",
						Author: "test-author-one",
						Price:  10.0,
						ISBN:   "1234567890",
						Active: true,
					}
					bl.Count = 1
				},
				data: &models.Book{
					ID:     1,
					Title:  "",
					Author: "test-author-one",
					Price:  10.0,
					ISBN:   "1234567890",
					Active: true,
				},
				checks: CheckList(
					hasError(true),
				),
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				bl  = NewBookList()
				got *models.Book
			)

			if tt.before != nil {
				tt.before(bl)
			}

			got, err = bl.Update(tt.id, tt.data)
			for _, check := range tt.checks {
				check(t, got, bl)
			}
		})
	}
}

func TestBookList_Delete(t *testing.T) {
	type CheckFn func(t *testing.T, bl *BookList)

	var (
		CheckList = func(fns ...CheckFn) []CheckFn { return fns }
		err       error

		hasError = func(want bool) CheckFn {
			return func(t *testing.T, _ *BookList) {
				t.Helper()
				if want {
					assert.NotNil(t, err, "hasError: error expected, none produced")
				} else {
					assert.Nil(t, err, "hasError = [+%v], no error expected")
				}
			}
		}

		isDeleted = func(id int32) CheckFn {
			return func(t *testing.T, bl *BookList) {
				t.Helper()
				assert.False(t, bl.List[id].Active)
			}
		}

		tests = []struct {
			name   string
			id     int32
			before func(bl *BookList)
			checks []CheckFn
		}{
			{
				name: "success",
				id:   1,
				before: func(bl *BookList) {
					bl.List[1] = &models.Book{
						ID:     1,
						Title:  "test-book-one",
						Author: "test-author-one",
						Price:  10.0,
						ISBN:   "1234567890",
						Active: true,
					}
					bl.Count = 1
				},
				checks: CheckList(
					hasError(false),
					isDeleted(1),
				),
			},
			{
				name: "not-found",
				id:   1,
				checks: CheckList(
					hasError(true),
				),
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bl := NewBookList()

			if tt.before != nil {
				tt.before(bl)
			}

			err = bl.Delete(tt.id)
			for _, check := range tt.checks {
				check(t, bl)
			}
		})
	}
}
