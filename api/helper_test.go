package api

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CheckFn func(t *testing.T, r *httptest.ResponseRecorder)

var (
	CheckList = func(fns ...CheckFn) []CheckFn { return fns }

	statusCode = func(want int) CheckFn {
		return func(t *testing.T, r *httptest.ResponseRecorder) {
			t.Helper()
			assert.Equal(t, want, r.Result().StatusCode)
		}
	}
)
