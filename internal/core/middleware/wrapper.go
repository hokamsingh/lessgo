package middleware

import "net/http"

// MiddlewareWrapper wraps a function to match the Middleware interface
type MiddlewareWrapper struct {
	HandlerFunc func(next http.Handler) http.Handler
}

// Handle implements the Middleware interface
func (mw MiddlewareWrapper) Handle(next http.Handler) http.Handler {
	return mw.HandlerFunc(next)
}
