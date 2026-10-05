package middleware

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

// Chain wraps h so that middlewares run in the given order: the first one is
// the outermost.
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range slices.Backward(middlewares) {
		h = middleware(h)
	}

	return h
}
