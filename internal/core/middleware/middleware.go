/*
Package middleware provides a simple interface and base implementation for creating HTTP middlewares.

This package defines a `Middleware` interface with a `Handle` method that allows chaining HTTP handlers. It also includes a `BaseMiddleware` struct that provides a basic implementation of the middleware pattern.

Usage:

	import (
		"net/http"
		"github.com/hokamsingh/lessgo/pkg/lessgo/middleware"
	)

	func main() {
		mw := &middleware.BaseMiddleware{}

		mux := http.NewServeMux()
		mux.Handle("/", mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})))

		http.ListenAndServe(":8080", mux)
	}
*/
package middleware

import "net/http"

// Middleware defines the interface for HTTP middlewares.
// Implementers should provide a `Handle` method that takes an `http.Handler` and returns a new `http.Handler`.
// This allows for wrapping existing handlers with additional functionality.
type Middleware interface {
	Handle(next http.Handler) http.Handler
}

// BaseMiddleware provides a basic implementation of the Middleware interface.
// It allows chaining of HTTP handlers by passing the request to the next handler in the chain.
//
// Example:
//
//	mw := &middleware.BaseMiddleware{}
//	http.Handle("/", mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	    w.Write([]byte("Hello, World!"))
//	})))
//
//	http.ListenAndServe(":8080", nil)
type BaseMiddleware struct{}

// Handle is the method that wraps an existing `http.Handler` with additional functionality.
// This method calls the next handler in the chain after executing any custom logic.
//
// Example:
//
//	mw := &middleware.BaseMiddleware{}
//
//	http.Handle("/", mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	    log.Println("Request received")
//	    w.W
