package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/padiazg/nurdsoft-challenge/internals"
	"github.com/padiazg/nurdsoft-challenge/models"
	"github.com/stretchr/testify/assert"
)

func Test_retrieveOneHandlerFn(t *testing.T) {
	var (
		expected = func(book *models.Book) CheckFn {
			return func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()

				data := &models.Book{}
				json.Unmarshal(r.Body.Bytes(), &data)

				assert.EqualValues(t, book, data)
			}
		}

		tests = []struct {
			name   string
			id     int32
			before func(bl *internals.BookList)
			checks []CheckFn
		}{
			{
				name: "success",
				id:   1,
				before: func(bl *internals.BookList) {
					bl.List = map[int32]*models.Book{
						1: {
							ID:     1,
							Title:  "test-book-one",
							Author: "test-author-one",
							Price:  10.0,
							ISBN:   "1234567890",
							Active: true,
						},
					}

					bl.Count = int32(len(bl.List))
				},
				checks: CheckList(
					statusCode(http.StatusOK),
					expected(&models.Book{ID: 1, Title: "test-book-one", Author: "test-author-one", Price: 10.0, ISBN: "1234567890", Active: true}),
				),
			},
			{
				name: "not-found",
				id:   2,
				before: func(bl *internals.BookList) {
					bl.List = map[int32]*models.Book{
						1: {
							ID:     1,
							Title:  "test-book-one",
							Author: "test-author-one",
							Price:  10.0,
							ISBN:   "1234567890",
							Active: true,
						},
					}

					bl.Count = int32(len(bl.List))
				},
				checks: CheckList(
					statusCode(http.StatusNotFound),
				),
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var (
				server   = NewServer(&models.Config{}, nil)
				recorder = httptest.NewRecorder()
			)

			if tt.before != nil {
				tt.before(server.data)
			}

			server.router.ServeHTTP(recorder, httptest.NewRequest("GET", fmt.Sprintf("/books/%d", tt.id), nil))
			for _, check := range tt.checks {
				check(t, recorder)
			}
		})
	}
}
