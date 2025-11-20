package handler

import "net/http"

// HTTPHandler defines the minimal interface all handlers must implement.
// Makes handlers easy to test and mock.
type HTTPHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
