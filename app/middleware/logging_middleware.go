package middleware

import (
	"log"
	"net/http"

	core "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type LoggingMiddleare struct {
	core.BaseMiddleware
}

func NewLoggingMiddleware() *LoggingMiddleare {
	return &LoggingMiddleare{}
}

func (lm *LoggingMiddleare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request recieved: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
