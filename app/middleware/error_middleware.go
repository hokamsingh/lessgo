package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

type ErrorMiddleware struct{}

func NewErrorHandleMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

func (em *ErrorMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("An error occurred: %v", err)
				log.Printf("Stack trace:\n%s\n", debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
