package api

import (
	"testing"

	"github.com/padiazg/nurdsoft-challenge/internals"
	"github.com/padiazg/nurdsoft-challenge/models"
	"github.com/stretchr/testify/assert"
)

func Test_NewServer(t *testing.T) {
	type checkFn func(t *testing.T, s *Server)

	var (
		checkList = func(fns ...checkFn) []checkFn { return fns }

		checkPort = func(port int) checkFn {
			return func(t *testing.T, server *Server) {
				t.Helper()
				assert.Equal(t, port, server.config.Port)
			}
		}

		checkRoot = func(root string) checkFn {
			return func(t *testing.T, server *Server) {
				t.Helper()
				assert.Equal(t, root, server.config.Root)
			}
		}

		checkData = func(data *internals.BookList) checkFn {
			return func(t *testing.T, server *Server) {
				t.Helper()
				assert.EqualValues(t, data, server.data)
			}
		}

		emptyData = internals.NewBookList()

		tests = []struct {
			name   string
			config models.Config
			data   *internals.BookList
			checks []checkFn
		}{
			{
				name:   "default-values",
				config: models.Config{},
				data:   nil,
				checks: checkList(
					checkPort(8000),
					checkRoot("/books"),
				),
			},
			{
				name: "custom-values",
				config: models.Config{
					Port: 9000,
					Root: "/api",
				},
				data: emptyData,
				checks: checkList(
					checkPort(9000),
					checkRoot("/api"),
					checkData(emptyData),
				),
			},
		}
	)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer(&tt.config, nil)

			for _, check := range tt.checks {
				check(t, server)
			}
		})
	}
}
