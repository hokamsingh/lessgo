package router

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
	"github.com/hokamsingh/lessgo/internal/core/middleware"
)

type Router struct {
	Mux        *mux.Router
	middleware []middleware.Middleware
}

type Option func(*Router)

// Default CORS options
var defaultCORSOptions = middleware.CORSOptions{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders: []string{"Content-Type", "Authorization"},
}

// NewRouter creates a new Router with optional configuration
func NewRouter(options ...Option) *Router {
	r := &Router{
		Mux:        mux.NewRouter(),
		middleware: []middleware.Middleware{},
	}
	for _, opt := range options {
		opt(r)
	}
	// Apply default CORS options
	// r.Use(middleware.NewCORSMiddleware(defaultCORSOptions))
	return r
}

// WithCORS enables CORS middleware with specific options
func WithCORS(options middleware.CORSOptions) Option {
	return func(r *Router) {
		corsMiddleware := middleware.NewCORSMiddleware(options)
		r.Use(corsMiddleware)
	}
}

func WithRateLimiter(limit int, interval time.Duration) Option {
	return func(r *Router) {
		rateLimiter := middleware.NewRateLimiter(limit, interval)
		r.Use(rateLimiter)
	}
}

func WithJSONParser() Option {
	return func(r *Router) {
		jsonParser := middleware.MiddlewareWrapper{HandlerFunc: middleware.JSONParser}
		r.Use(jsonParser)
	}
}

func WithCookieParser() Option {
	return func(r *Router) {
		cookieParser := middleware.MiddlewareWrapper{HandlerFunc: middleware.CookieParser}
		r.Use(cookieParser)
	}
}

func WithFileUpload(uploadDir string) Option {
	return func(r *Router) {
		fileUploadMiddleware := middleware.NewFileUploadMiddleware(uploadDir)
		r.Use(fileUploadMiddleware)
	}
}

func (r *Router) Use(m middleware.Middleware) {
	r.middleware = append(r.middleware, m)
}

func (r *Router) AddRoute(path string, handler http.HandlerFunc) {
	// Apply logging and error handling to the handler
	handler = r.withErrorHandling(handler)
	handler = r.withLogging(handler)
	r.Mux.HandleFunc(path, handler)
}

func (r *Router) Start(addr string) error {
	// Apply middlewares
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

// withLogging logs the request method and path
func (r *Router) withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
