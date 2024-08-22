package router

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
	"github.com/hokamsingh/lessgo/internal/core/context"
	"github.com/hokamsingh/lessgo/internal/core/middleware"
)

// Router represents an HTTP router with middleware support and error handling.
type Router struct {
	Mux        *mux.Router
	middleware []middleware.Middleware
}

// Option is a function that configures a Router.
type Option func(*Router)

// Default CORS options
var _ = middleware.CORSOptions{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders: []string{"Content-Type", "Authorization"},
}

// NewRouter creates a new Router with optional configuration.
// You can pass options like WithCORS or WithJSONParser to configure the router.
//
// Example usage:
//
//	r := router.NewRouter(
//		router.WithCORS(middleware.CORSOptions{}),
//		router.WithJSONParser(),
//	)
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

// SubRouter creates a subrouter with the given path prefix.
//
// Example usage:
//
//	subRouter := r.SubRouter("/api")
//	subRouter.AddRoute("/ping", handler)
func (r *Router) SubRouter(pathPrefix string, options ...Option) *Router {
	subRouter := &Router{
		Mux:        r.Mux.PathPrefix(pathPrefix).Subrouter(),
		middleware: append([]middleware.Middleware{}, r.middleware...),
	}
	// Apply options to the subrouter
	for _, opt := range options {
		opt(subRouter)
	}
	return subRouter
}

// WithCORS enables CORS middleware with specific options.
// This option configures the CORS settings for the router.
//
// Example usage:
//
//	r := router.NewRouter(router.WithCORS(middleware.CORSOptions{...}))
func WithCORS(options middleware.CORSOptions) Option {
	return func(r *Router) {
		corsMiddleware := middleware.NewCORSMiddleware(options)
		r.Use(corsMiddleware)
	}
}

// WithRateLimiter enables rate limiting middleware with the specified limit and interval.
// This option configures the rate limiter for the router.
//
// Example usage:
//
//	r := router.NewRouter(router.WithRateLimiter(100, time.Minute))
func WithRateLimiter(limit int, interval time.Duration) Option {
	return func(r *Router) {
		rateLimiter := middleware.NewRateLimiter(limit, interval)
		r.Use(rateLimiter)
	}
}

// WithJSONParser enables JSON parsing middleware for request bodies.
// This option ensures that incoming JSON payloads are parsed and available in the request context.
//
// Example usage:
//
//	r := router.NewRouter(router.WithJSONParser())
func WithJSONParser() Option {
	return func(r *Router) {
		jsonParser := middleware.MiddlewareWrapper{HandlerFunc: middleware.JSONParser}
		r.Use(jsonParser)
	}
}

// WithCookieParser enables cookie parsing middleware.
// This option ensures that cookies are parsed and available in the request context.
//
// Example usage:
//
//	r := router.NewRouter(router.WithCookieParser())
func WithCookieParser() Option {
	return func(r *Router) {
		cookieParser := middleware.MiddlewareWrapper{HandlerFunc: middleware.CookieParser}
		r.Use(cookieParser)
	}
}

// WithFileUpload enables file upload middleware with the specified upload directory.
// This option configures the router to handle file uploads and save them to the given directory.
//
// Example usage:
//
//	r := router.NewRouter(router.WithFileUpload("/uploads"))
func WithFileUpload(uploadDir string) Option {
	return func(r *Router) {
		fileUploadMiddleware := middleware.NewFileUploadMiddleware(uploadDir)
		r.Use(fileUploadMiddleware)
	}
}

// Use adds a middleware to the router's middleware stack.
//
// Example usage:
//
//	r.Use(middleware.LoggingMiddleware{})
func (r *Router) Use(m middleware.Middleware) {
	r.middleware = append(r.middleware, m)
}

// AddRoute adds a route with the given path and handler function.
// This method applies context, error handling, and logging to the handler.
//
// Example usage:
//
//	r.AddRoute("/ping", func(ctx *context.Context) {
//		ctx.JSON(http.StatusOK, map[string]string{"message": "pong"})
//	})
func (r *Router) AddRoute(path string, handler CustomHandler) {
	// Create an HTTP handler function that uses the custom context
	handlerFunc := WrapCustomHandler(handler)
	// Wrap the handler function with error handling and logging
	handlerFunc = r.withErrorHandling(handlerFunc)
	handlerFunc = r.withLogging(handlerFunc)
	r.Mux.HandleFunc(path, handlerFunc)
}

// Start starts the HTTP server on the specified address.
// It applies all middleware and listens for incoming requests.
//
// Example usage:
//
//	err := r.Start(":8080")
//	if err != nil {
//		log.Fatalf("Server failed: %v", err)
//	}
func (r *Router) Start(addr string) error {
	// Apply middlewares
	finalHandler := http.Handler(r.Mux)
	for _, m := range r.middleware {
		finalHandler = m.Handle(finalHandler)
	}
	return http.ListenAndServe(addr, finalHandler)
}

// Start http server
func (r *Router) Listen(addr string) error {
	return r.Start(addr)
}

// HTTPError represents an error with an associated HTTP status code.
type HTTPError struct {
	Code    int
	Message string
}

// Error returns a string representation of the HTTPError.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

// NewHTTPError creates a new HTTPError instance with the given status code and message.
//
// Example usage:
//
//	err := NewHTTPError(http.StatusBadRequest, "Bad Request: missing parameters")
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

// withLogging logs the request method and path.
func (r *Router) withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

// CustomHandler is a function type that takes a custom Context.
type CustomHandler func(ctx *context.Context)

// Get registers a handler for GET requests.
func (r *Router) Get(path string, handler CustomHandler) {
	r.AddRoute(path, UnWrapCustomHandler(r.withContext(handler, "GET")))
}

// Post registers a handler for POST requests.
func (r *Router) Post(path string, handler CustomHandler) {
	r.AddRoute(path, UnWrapCustomHandler(r.withContext(handler, "POST")))
}

// Put registers a handler for PUT requests.
func (r *Router) Put(path string, handler CustomHandler) {
	r.AddRoute(path, UnWrapCustomHandler(r.withContext(handler, "PUT")))
}

// Delete registers a handler for DELETE requests.
func (r *Router) Delete(path string, handler CustomHandler) {
	r.AddRoute(path, UnWrapCustomHandler(r.withContext(handler, "DELETE")))
}

// Patch registers a handler for PATCH requests.
func (r *Router) Patch(path string, handler CustomHandler) {
	r.AddRoute(path, UnWrapCustomHandler(r.withContext(handler, "PATCH")))
}

// WrapCustomHandler converts a CustomHandler to http.HandlerFunc.
func WrapCustomHandler(handler CustomHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(r, w)
		handler(ctx)
	}
}

// UnWrapCustomHandler converts a http.HandlerFunc to CustomHandler.
func UnWrapCustomHandler(handler http.HandlerFunc) CustomHandler {
	return func(ctx *context.Context) {
		handler.ServeHTTP(ctx.Res, ctx.Req)
	}
}

// withContext wraps the given handler with a custom context.
// This provides utility methods for handling requests and responses.
// It transforms the original handler to use the custom Context.
//
// Example usage:
//
//	r.AddRoute("/example", func(ctx *context.Context) {
//		ctx.JSON(http.StatusOK, map[string]string{"message": "Hello, world!"})
//	})
func (r *Router) withContext(next CustomHandler, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		ctx := context.NewContext(req, w)
		next(ctx)
	}
}

// ServeStatic creates a file server handler to serve static files from the given directory.
// The pathPrefix is stripped from the request URL before serving the file.
//
// Example usage:
//
//	 r := LessGo.NewRouter(
//			LessGo.WithCORS(*corsOptions),
//			LessGo.WithRateLimiter(100, 1*time.Minute),
//			LessGo.WithJSONParser(),
//			LessGo.WithCookieParser(),
//		)
//	r.ServeStatic("/static/", "/path/to/static/files"))
func (r *Router) ServeStatic(pathPrefix, dir string) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path: %v", err)
	}
	fs := http.FileServer(http.Dir(absPath))
	r.Mux.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix, fs))
}
