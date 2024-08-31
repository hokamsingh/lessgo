package router

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/hokamsingh/lessgo/internal/core/context"
	"github.com/hokamsingh/lessgo/internal/core/middleware"
	"github.com/hokamsingh/lessgo/internal/utils"
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

// HTTPMethod represents an HTTP method as a custom type.
type HTTPMethod string

// Define constants for each HTTP method.
const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	OPTIONS HTTPMethod = "OPTIONS"
	HEAD    HTTPMethod = "HEAD"
	PATCH   HTTPMethod = "PATCH"
)

var (
	appInstance *Router
	once        sync.Once
)

// SetAppInstance sets the singleton App instance.
func SetAppInstance(app *Router) {
	once.Do(func() {
		appInstance = app
	})
}

// GetAppInstance returns the singleton App instance.
func GetApp() *Router {
	return appInstance
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
	SetAppInstance(r)
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

	// Apply the middleware to the subrouter's Mux
	for _, m := range subRouter.middleware {
		subRouter.Mux.Use(m.Handle)
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

type RateLimiterType = middleware.RateLimiterType

const (
	InMemory RateLimiterType = iota
	RedisBacked
)

// WithRateLimiter enables rate limiting middleware with the specified limit and interval.
// This option configures the rate limiter for the router.
//
// Example usage:
//
//	r := router.NewRouter(router.WithRateLimiter(100, time.Minute))
func WithInMemoryRateLimiter(NumShards int, Limit int, Interval time.Duration, CleanupInterval time.Duration) Option {
	return func(r *Router) {
		config := middleware.NewInMemoryConfig(NumShards, Limit, Interval, CleanupInterval)
		rateLimiter := middleware.NewRateLimiter(InMemory, config)
		r.Use(rateLimiter)
	}
}

// WithRateLimiter enables rate limiting middleware with the specified limit and interval.
// This option configures the rate limiter for the router.
//
// Example usage:
//
//	r := router.NewRouter(router.WithRateLimiter(100, time.Minute))
func WithRedisRateLimiter(client *redis.Client, limit int, interval time.Duration) Option {
	return func(r *Router) {
		config := middleware.NewRedisConfig(client, limit, interval)
		rateLimiter := middleware.NewRateLimiter(RedisBacked, config)
		r.Use(rateLimiter)
	}
}

// WithJSONParser enables JSON parsing middleware for request bodies.
// This option ensures that incoming JSON payloads are parsed and available in the request context.
//
// Example usage:
//
//	r := router.NewRouter(router.WithJSONParser())
func WithJSONParser(options middleware.ParserOptions) Option {
	return func(r *Router) {
		// jsonParser := middleware.MiddlewareWrapper{HandlerFunc: middleware.JSONParser}
		jsonParser := middleware.NewJsonParser(options)
		r.Use(jsonParser)
	}
}

// WithCaching is an option function that enables caching for the router using Redis.
//
// This function returns an Option that can be passed to the Router to enable
// response caching with Redis. Cached responses will be stored in Redis with a
// specified Time-To-Live (TTL), meaning they will automatically expire after the
// specified duration.
//
// Parameters:
//   - redisAddr (string): The address of the Redis server, e.g., "localhost:6379".
//   - ttl (time.Duration): The Time-To-Live for cached responses. Responses will
//     be removed from the cache after this duration.
//
// Returns:
//   - Option: An option that applies caching middleware to the router.
//
// Example usage:
//
//	router := NewRouter(
//	    WithCaching("localhost:6379", 5*time.Minute),
//	)
//
// This will enable caching for the router, storing responses in Redis for 5 minutes.
//
// Note: Ensure that the Redis server is running and accessible at the specified
// address.
func WithCaching(client *redis.Client, ttl time.Duration, cacheControl bool) Option {
	return func(r *Router) {
		caching := middleware.NewCaching(client, ttl, cacheControl)
		r.Use(caching)
	}
}

// WithCsrf is an option function that enables CSRF protection for the router.
//
// This function returns an Option that can be passed to the Router to enable
// Cross-Site Request Forgery (CSRF) protection using a middleware. The middleware
// generates and validates CSRF tokens to protect against malicious cross-origin
// requests, ensuring that requests are coming from legitimate users.
//
// Returns:
//   - Option: An option that applies CSRF protection middleware to the router.
//
// Example usage:
//
//	router := NewRouter(
//	    WithCsrf(),
//	)
//
// This will enable CSRF protection for all routes in the router.
func WithCsrf() Option {
	return func(r *Router) {
		csrf := middleware.NewCSRFProtection()
		r.Use(csrf)
	}
}

// WithXss is an option function that enables XSS protection for the router.
//
// This function returns an Option that can be passed to the Router to enable
// Cross-Site Scripting (XSS) protection using a middleware. The middleware
// helps to sanitize and filter out malicious scripts from user input, thereby
// preventing XSS attacks.
//
// Returns:
//   - Option: An option that applies XSS protection middleware to the router.
//
// Example usage:
//
//	router := NewRouter(
//	    WithXss(),
//	)
//
// This will enable XSS protection for all routes in the router, ensuring that
// user input is sanitized and secure.
func WithXss() Option {
	return func(r *Router) {
		xss := middleware.NewXSSProtection()
		r.Use(xss)
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
		// cookieParser := middleware.MiddlewareWrapper{HandlerFunc: middleware.CookieParser}
		cookieParser := middleware.NewCookieParser()
		r.Use(cookieParser)
	}
}

// WithFileUpload enables file upload middleware with the specified upload directory.
// This option configures the router to handle file uploads and save them to the given directory.
//
// Example usage:
//
//	r := router.NewRouter(router.WithFileUpload("/uploads"))
func WithFileUpload(uploadDir string, maxFileSize int64, allowedExts []string) Option {
	return func(r *Router) {
		fileUploadMiddleware := middleware.NewFileUploadMiddleware(uploadDir, maxFileSize, allowedExts)
		r.Use(fileUploadMiddleware)
	}
}

// WithTemplateRendering sets up the router to use the TemplateMiddleware for rendering HTML templates.
// It automatically loads all `.html` files from the specified directory and makes them available
// for rendering within the application's handlers.
//
// The middleware parses all `.html` files from the provided directory during initialization
// and injects the parsed templates into the request context, allowing handlers to access and render
// the templates as needed.
//
// Usage:
//
//	router := NewRouter(
//	    WithTemplateRendering("templates"), // Directory containing all .html files
//	)
//
//	router.HandleFunc("/", yourHandler)
//
// In the handler, you can retrieve and execute a template:
//
//	func yourHandler(w http.ResponseWriter, r *http.Request) {
//	    tmpl := middleware.GetTemplate(r.Context())
//	    tmpl.ExecuteTemplate(w, "index.html", nil) // Renders the index.html template
//	}
//
// Parameters:
//   - templateDir: The directory containing the `.html` files to be used as templates.
//
// Returns:
//   - Option: A function that configures the router to use the template rendering middleware.
func WithTemplateRendering(templateDir string) Option {
	return func(r *Router) {
		renderer := middleware.NewTemplateMiddleware(templateDir)
		r.Use(renderer)
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
	utils.Assert(path[0] == '/', "path must begin with '/'")
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

	server := &http.Server{
		Addr:         addr,
		Handler:      finalHandler,
		ReadTimeout:  5 * time.Second,   // Defaults timeout
		WriteTimeout: 10 * time.Second,  // Defaults timeout
		IdleTimeout:  120 * time.Second, // Defaults timeout
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
	return err
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

func (r *Router) WithContentNegotiation(next http.HandlerFunc) http.HandlerFunc {
	return ContentNegotiationHandler
}

// CustomHandler is a function type that takes a custom Context.
type CustomHandler func(ctx *context.Context)

// Server Swagger
func (r *Router) Swagger(path string, handler http.HandlerFunc) {
	r.AddRoute(path, UnWrapCustomHandler(r.withContext(UnWrapCustomHandler(handler), string(GET))))
}

func PathPrefix(path string) {

}

// Get registers a handler for GET requests.
func (r *Router) Get(path string, handler CustomHandler) *Router {
	r.AddRoute(path, UnWrapCustomHandler(r.withContext(handler, string(GET))))
	return r
}

// Post registers a handler for POST requests.
func (r *Router) Post(path string, handler CustomHandler) *Router {
	r.AddRoute(path, UnWrapCustomHandler(r.withContext(handler, string(POST))))
	return r
}

// Put registers a handler for PUT requests.
func (r *Router) Put(path string, handler CustomHandler) *Router {
	r.AddRoute(path, UnWrapCustomHandler(r.withContext(handler, string(PUT))))
	return r
}

// Delete registers a handler for DELETE requests.
func (r *Router) Delete(path string, handler CustomHandler) *Router {
	r.AddRoute(path, UnWrapCustomHandler(r.withContext(handler, string(DELETE))))
	return r
}

// Patch registers a handler for PATCH requests.
func (r *Router) Patch(path string, handler CustomHandler) *Router {
	r.AddRoute(path, UnWrapCustomHandler(r.withContext(handler, string(PATCH))))
	return r
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
//	r.AddRoute("/example", func(ctx *LessGo.Context) {
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

// Content negotiation
const (
	ContentTypeJSON = "application/json"
	ContentTypeXML  = "application/xml"
	ContentTypeHTML = "text/html"
)

func ContentNegotiationHandler(w http.ResponseWriter, r *http.Request) {
	acceptHeader := r.Header.Get("Accept")
	contentType := NegotiateContentType(acceptHeader)

	var response []byte
	var err error

	// Prepare response based on content type
	switch contentType {
	case ContentTypeJSON:
		w.Header().Set("Content-Type", ContentTypeJSON)
		response, err = json.Marshal(map[string]string{"message": "Hello, JSON!"})
	case ContentTypeXML:
		w.Header().Set("Content-Type", ContentTypeXML)
		response, err = xml.Marshal(map[string]string{"message": "Hello, XML!"})
	case ContentTypeHTML:
		w.Header().Set("Content-Type", ContentTypeHTML)
		response = []byte("<html><body><h1>Hello, HTML!</h1></body></html>")
	default:
		// If no acceptable content type is found, return 406
		http.Error(w, "Not Acceptable", http.StatusNotAcceptable)
		return
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func NegotiateContentType(acceptHeader string) string {
	// Default to JSON if nothing is specified
	if acceptHeader == "" {
		return ContentTypeJSON
	}

	// Split the Accept header into supported media types
	acceptedTypes := strings.Split(acceptHeader, ",")

	// Check for supported media types in order of preference
	for _, acceptedType := range acceptedTypes {
		acceptedType = strings.TrimSpace(strings.Split(acceptedType, ";")[0])
		switch acceptedType {
		case ContentTypeJSON:
			return ContentTypeJSON
		case ContentTypeXML:
			return ContentTypeXML
		case ContentTypeHTML:
			return ContentTypeHTML
		}
	}

	// Default to JSON if no match is found
	return ContentTypeJSON
}
