package router

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
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

// SubRouter creates a subrouter with the given path prefix
func (r *Router) SubRouter(pathPrefix string) *Router {
	subRouter := &Router{
		Mux:        r.Mux.PathPrefix(pathPrefix).Subrouter(),
		middleware: []middleware.Middleware{},
	}
	return subRouter
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

// HTTPError represents an error with an associated HTTP status code.
type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

/*
withErrorHandling wraps the given HTTP handler function with centralized error handling.

This middleware intercepts any panics that occur during the execution of the handler function,
and handles them based on their type. If the panic is of type *HTTPError, it sends an HTTP response
with the specified status code and message. For other types of panics, it logs the error and stack trace,
and sends a generic "Internal Server Error" response.

Example:

	handler := r.withErrorHandling(func(w http.ResponseWriter, r *http.Request) {
		// Example: Trigger a Bad Request error
		panic(LessGo.NewHTTPError(http.StatusBadRequest, "Bad Request: missing parameters"))
	})

	// Use the handler in your router
	r.Mux.Handle("/example", handler)
*/
func (r *Router) withErrorHandling(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				switch e := err.(type) {
				case *HTTPError:
					log.Printf("HTTP error occurred: %v", e)
					http.Error(w, e.Message, e.Code)
				default:
					log.Printf("An unexpected error occurred: %v", err)
					log.Printf("Stack trace:\n%s\n", debug.Stack())
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()
		next(w, req)
	}
}

// withLogging logs the request method and path
func (r *Router) withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

// ServeStatic creates a file server handler to serve static files
func ServeStatic(pathPrefix, dir string) http.Handler {
	// Resolve the absolute path for debugging
	absPath, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path: %v", err)
	}
	log.Printf("Serving static files from: %s", absPath)

	fs := http.FileServer(http.Dir(absPath))
	return http.StripPrefix(pathPrefix, fs)
}
