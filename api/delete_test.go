package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/padiazg/nurdsoft-challenge/internals"
	"github.com/padiazg/nurdsoft-challenge/models"
)

func Test_deleteHandlerFn(t *testing.T) {
	var (
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

			server.router.ServeHTTP(recorder,
				httptest.NewRequest("DELETE", fmt.Sprintf("/books/%d", tt.id), nil),
			)
			for _, check := range tt.checks {
				check(t, recorder)
			}
		})
	}
}
