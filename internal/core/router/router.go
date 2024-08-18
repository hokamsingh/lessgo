package router

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/hokamsingh/lessgo/internal/core/middleware"
)

type Router struct {
	Mux        *mux.Router
	middleware []middleware.Middleware
}

func NewRouter() *Router {
	return &Router{
		Mux:        mux.NewRouter(),
	}
}

func (r *Router) Use(m middleware.Middleware) {
	r.middleware = append(r.middleware, m)
}

func (r *Router) AddRoute(path string, handler http.HandlerFunc) {
	r.Mux.HandleFunc(path, handler)
}

func (r *Router) Start(addr string) error {
	// apply middlewares
	finalHandler := http.Handler(r.Mux)
	for _, m := range r.middleware {
		finalHandler = m.Handle(finalHandler)
	}
	return http.ListenAndServe(addr, finalHandler)
}

// withErrorHandling wraps the given handler with error handling middleware
func (r *Router) withErrorHandling(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("An error occurred: %v", err)
				log.Printf("Stack trace:\n%s\n", debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

func (r *Router) WithLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		func() {
			log.Printf("Recieved %s %s", r.Method, r.URL.Path)
		}()
		next(w, r)
	}
}
