package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/padiazg/nurdsoft-challenge/internals"
	"github.com/padiazg/nurdsoft-challenge/models"
	"github.com/stretchr/testify/assert"
)

func Test_updateHandlerFn(t *testing.T) {
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
			name    string
			id      int32
			payload string
			before  func(bl *internals.BookList)
			checks  []CheckFn
		}{
			{
				name: "success",
				id:   1,
				before: func(bl *internals.BookList) {
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
				payload: `{"Title":"test-book-one (updated)", "Author":"test-author-one", "Price": 10.0, "ISBN": "1234567890", "Active": true}`,
				checks: CheckList(
					statusCode(http.StatusCreated),
					expected(&models.Book{ID: 1, Title: "test-book-one (updated)", Author: "test-author-one", Price: 10.0, ISBN: "1234567890", Active: true}),
				),
			},
			{
				name: "not-found",
				id:   2,
				before: func(bl *internals.BookList) {
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
				payload: `{"Title":"test-book-one (updated)", "Author":"test-author-one", "Price": 10.0, "ISBN": "1234567890", "Active": true}`,
				checks: CheckList(
					statusCode(http.StatusNotFound),
				),
			},
			{
				name: "bad-request",
				id:   1,
				before: func(bl *internals.BookList) {
					bl.List[1] = &models.Book{
						ID:     1,
						Title:  "test-book-one ",
						Author: "test-author-one",
						Price:  10.0,
						ISBN:   "1234567890",
						Active: true,
					}
					bl.Count = 1
				},
				payload: `{"Title":"", "Author":"test-author-one", "Price": 10.0, "ISBN": "1234567890", "Active": true}`,
				checks: CheckList(
					statusCode(http.StatusBadRequest),
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

			req := httptest.NewRequest("PUT", fmt.Sprintf("/books/%d", tt.id), strings.NewReader(tt.payload))
			req.Header.Add("Content-Type", "application/json")

			server.router.ServeHTTP(recorder, req)
			for _, check := range tt.checks {
				check(t, recorder)
			}
		})
	}
}
