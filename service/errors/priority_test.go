package errors_test

import (
	"net/http"
	"testing"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"

	"github.com/stretchr/testify/assert"
)

func Test_ErrorPriority(t *testing.T) {

	tests := []struct {
		name     string
		status   int
		expected int
	}{
		{"BadRequest → highest priority", 400, 1},
		{"Unauthorized", 401, 2},
		{"Forbidden", 403, 3},
		{"NotFound", 404, 4},
		{"UnprocessableEntity", 422, 5},
		{"TooManyRequests", 429, 6},

		{"BadGateway → provider failure", 502, 7},
		{"ServiceUnavailable", 503, 8},

		{"InternalServerError → lowest", 500, 9},

		{"Unknown status → fallback", 418, 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, errors.ErrorPriority(tt.status))
		})
	}
}

func Test_PickBetter(t *testing.T) {

	makeErr := func(status int) errors.AppError {
		return &mockErr{status: status}
	}

	t.Run("current nil → return new", func(t *testing.T) {
		newErr := makeErr(400)
		got := errors.PickBetter(nil, newErr)
		assert.Equal(t, newErr, got)
	})

	t.Run("new nil → keep current", func(t *testing.T) {
		current := makeErr(400)
		got := errors.PickBetter(current, nil)
		assert.Equal(t, current, got)
	})

	t.Run("new has higher priority → replace", func(t *testing.T) {
		current := makeErr(http.StatusInternalServerError) // priority 9
		newErr := makeErr(http.StatusBadRequest)           // priority 1

		got := errors.PickBetter(current, newErr)
		assert.Equal(t, newErr, got)
	})

	t.Run("current has higher priority → keep", func(t *testing.T) {
		current := makeErr(http.StatusBadRequest) // priority 1
		newErr := makeErr(http.StatusInternalServerError)

		got := errors.PickBetter(current, newErr)
		assert.Equal(t, current, got)
	})

	t.Run("equal priority → keep current", func(t *testing.T) {
		current := makeErr(http.StatusBadGateway)
		newErr := makeErr(http.StatusBadGateway)

		got := errors.PickBetter(current, newErr)
		assert.Equal(t, current, got)
	})
}

type mockErr struct {
	status int
}

func (m *mockErr) Error() string   { return "mock" }
func (m *mockErr) StatusCode() int { return m.status }
func (m *mockErr) Code() string    { return "MOCK" }
func (m *mockErr) Message() string { return "mock error" }
