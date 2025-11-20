package server

import "net/http"

// Middleware type
type Middleware func(http.Handler) http.Handler

// Chain middleware like Gin
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
