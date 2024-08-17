package middleware

import "net/http"

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

type BaseMiddleware struct{}

func (bm *BaseMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logic goes here
		next.ServeHTTP(w, r)
	})
}
