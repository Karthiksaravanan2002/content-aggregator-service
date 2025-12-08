package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadRequest(t *testing.T) {
	err := errors.New("invalid input")
	app := BadRequest(err)

	assert.Equal(t, 400, app.StatusCode())
	assert.Equal(t, "BAD_REQUEST", app.Code())
	assert.Equal(t, "invalid input", app.Message())

	// original error should be preserved
	assert.Equal(t, err, unwrap(app))
}

func TestInternalError(t *testing.T) {
	err := errors.New("system failure")
	app := InternalError(err)

	assert.Equal(t, 500, app.StatusCode())
	assert.Equal(t, "INTERNAL_SERVICE_ERROR", app.Code())
	assert.Equal(t, "system failure", app.Message())
	assert.Equal(t, err, unwrap(app))
}

func TestBadGateway(t *testing.T) {
	err := errors.New("provider timeout")
	app := BadGateway(err)

	assert.Equal(t, 502, app.StatusCode())
	assert.Equal(t, "EXTERNAL_SERVICE_ERROR", app.Code())
	assert.Equal(t, "provider timeout", app.Message())
	assert.Equal(t, err, unwrap(app))
}

func TestProviderError(t *testing.T) {
	err := errors.New("twitch unreachable")
	app := ProviderError(503, err)

	assert.Equal(t, 503, app.StatusCode())
	assert.Equal(t, "PROVIDER_ERROR", app.Code())
	assert.Equal(t, "twitch unreachable", app.Message())
	assert.Equal(t, err, unwrap(app))
}

func TestNewErr(t *testing.T) {
	err := errors.New("just wrapping")
	app := NewErr(err)

	// Your NewErr gives NO status/code/message → ensure that is true
	assert.Nil(t, unwrapStatusCode(app))
	assert.Nil(t, unwrapCode(app))
	assert.Nil(t, unwrapMessage(app))

	assert.Equal(t, err, unwrap(app))
}



func unwrap(app AppError) error {
	// appError embeds error, so .Error() should match original
	// But to get original error value, we need a type assertion.
	type causer interface{ Unwrap() error }

	if u, ok := app.(causer); ok {
		return u.Unwrap()
	}
	// If Unwrap() not implemented, we fallback to comparing Error() string
	return errors.New(app.Error())
}

func unwrapStatusCode(app AppError) *int {
	code := app.StatusCode()
	if code == 0 {
		return nil
	}
	return &code
}

func unwrapCode(app AppError) *string {
	c := app.Code()
	if c == "" {
		return nil
	}
	return &c
}

func unwrapMessage(app AppError) *string {
	m := app.Message()
	if m == "" {
		return nil
	}
	return &m
}
