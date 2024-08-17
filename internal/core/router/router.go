package router

import (
	"lessgo/internal/core/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	Mux        *mux.Router
	middleware []middleware.Middleware
}

func NewRouter() *Router {
	return &Router{
		Mux: mux.NewRouter(),
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
