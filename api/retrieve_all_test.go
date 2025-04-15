package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/padiazg/nurdsoft-challenge/internals"
	"github.com/padiazg/nurdsoft-challenge/models"
	"github.com/stretchr/testify/assert"
)

type CheckFn func(t *testing.T, r *httptest.ResponseRecorder)

func Test_retrieveAllHandlerFn(t *testing.T) {
	var (
		CheckList = func(fns ...CheckFn) []CheckFn { return fns }

		statusCode = func(want int) CheckFn {
			return func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()
				assert.Equal(t, want, r.Result().StatusCode)
			}
		}

		expected = func(list []*models.Book) CheckFn {
			return func(t *testing.T, r *httptest.ResponseRecorder) {
				t.Helper()

				data := []*models.Book{}
				json.Unmarshal(r.Body.Bytes(), &data)

				len0 := len(list)
				len1 := len(data)

				if assert.Equalf(t, len0, len1, "count differ: expected %d, got %d", len0, len1) {
					for _, b0 := range list {
						found := false
						for _, b1 := range data {
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
			before func(bl *internals.BookList)
			checks []CheckFn
		}{
			{
				name: "success",
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
					expected([]*models.Book{
						{ID: 1, Title: "test-book-one", Author: "test-author-one", Price: 10.0, ISBN: "1234567890", Active: true},
					}),
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

			server.router.ServeHTTP(recorder, httptest.NewRequest("GET", "/books", nil))
			for _, check := range tt.checks {
				check(t, recorder)
			}
		})
	}
}
