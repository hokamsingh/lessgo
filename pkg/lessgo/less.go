/*
Package LessGo is a minimalist web framework for building fast, scalable, and lightweight web applications in Go.

LessGo is designed to be simple and easy to use, providing essential features for web development without the overhead of large, complex frameworks. It emphasizes speed, flexibility, and a small footprint, making it ideal for developers who want to build web applications quickly while maintaining full control over their projects.

# Features

- **Routing**: LessGo provides a flexible and powerful routing mechanism for handling HTTP requests.
- **Middleware**: Support for middleware allows you to add custom functionality to your request pipeline.
- **Content Negotiation**: Built-in support for content negotiation, enabling your API to serve different content types like JSON, XML, etc.
- **Environment Configuration**: Load environment variables and .env files easily for configuration management.
- **Pluggable Architecture**: Extend LessGo with custom plugins and middleware for additional functionality.
- **CORS Support**: Configure CORS settings for your API.
- **Redis Integration**: Easily integrate Redis for caching, rate limiting, and other use cases.
- **Static File Serving**: Serve static files like HTML, CSS, JavaScript, or images.
- **Security**: CSRF and XSS protection are built-in to enhance the security of your applications.
- **Rate Limiting**: Implement rate limiting to protect your application from abuse.

# Usage

Here's an example of how to use the LessGo framework in a basic web server setup:

	package main

	import (
		"log"
		"time"
		"github.com/yourusername/LessGo/app/src"
		LessGo "github.com/yourusername/LessGo/pkg/lessgo"
	)

	func main() {
		// Configuration setup
		cfg := LessGo.LoadConfig()
		serverPort := cfg.Get("SERVER_PORT", "8080")
		env := cfg.Get("ENV", "development")
		addr := ":" + serverPort

		// Define CORS options
		corsOptions := LessGo.NewCorsOptions(
			[]string{"*"},
			[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			[]string{"Content-Type", "Authorization"},
		)

		// Initialize Redis client
		rClient := LessGo.NewRedisClient("localhost:6379")

		// Initialize app with middlewares
		App := LessGo.App(
			LessGo.WithCORS(*corsOptions),
			LessGo.WithJSONParser(LessGo.NewParserOptions(5*1024*1024)), // 5MB limit
			LessGo.WithCookieParser(),
			LessGo.WithCsrf(),
			LessGo.WithXss(),
			LessGo.WithCaching(rClient, 5*time.Minute, true),
			LessGo.WithRedisRateLimiter(rClient, 100, 1*time.Second),
		)

		// Serve static files
		App.ServeStatic("/static/", LessGo.GetFolderPath("uploads"))

		// Register modules and dependencies
		LessGo.RegisterDependencies([]interface{}{src.NewRootService, src.NewRootModule})
		LessGo.RegisterModules(App, []LessGo.IModule{src.NewRootModule(App)})

		// Example route
		App.Get("/ping", func(ctx *LessGo.Context) {
			ctx.Send("pong")
		})

		// Start the server
		log.Printf("Starting server on port %s in %s mode", serverPort, env)
		if err := App.Listen(addr); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}

For more detailed documentation and examples, please visit the [official LessGo documentation](https://github.com/hokamsingh).

# Package Structure

- **LessGo**: The core of the framework, providing the main application structure and utilities.
- **LessGo/context**: Handles the request context, providing methods to respond with different content types.
- **LessGo/middleware**: Contains built-in middleware functions for request handling.
- **LessGo/config**: Manages configuration loading from environment variables and .env files.
*/
package LessGo

import (
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hokamsingh/lessgo/internal/core/concurrency"
	"github.com/hokamsingh/lessgo/internal/core/config"
	"github.com/hokamsingh/lessgo/internal/core/context"
	"github.com/hokamsingh/lessgo/internal/core/controller"
	"github.com/hokamsingh/lessgo/internal/core/di"
	"github.com/hokamsingh/lessgo/internal/core/discovery"
	"github.com/hokamsingh/lessgo/internal/core/middleware"
	"github.com/hokamsingh/lessgo/internal/core/module"
	"github.com/hokamsingh/lessgo/internal/core/router"
	"github.com/hokamsingh/lessgo/internal/core/service"
	"github.com/hokamsingh/lessgo/internal/core/websocket"
	"github.com/hokamsingh/lessgo/internal/utils"
)

// Version
const Version = "v1.0.2"

// Expose core types

// Controller defines the interface that all controllers in the application must implement.
// Any controller that implements this interface must define the RegisterRoutes method,
// which is responsible for setting up the necessary routes for the controller.
type Controller = controller.Controller

// BaseController provides a default implementation of the Controller interface.
// It can be embedded in other controllers to inherit its default behavior,
// or overridden with custom implementations.
type BaseController = controller.BaseController

// Container wraps the `dig.Container` and provides methods for registering and invoking dependencies.
// This struct serves as the main entry point for setting up and managing dependency injection within the application.
type Container = di.Container

// Middleware defines the interface for HTTP middlewares.
// Implementers should provide a `Handle` method that takes an `http.Handler` and returns a new `http.Handler`.
// This allows for wrapping existing handlers with additional functionality.
type Middleware = middleware.Middleware

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
type BaseMiddleware = middleware.BaseMiddleware

// Module represents a module in the application.
// It holds the name, a list of controllers, services, and any submodules.
// The module can be used to organize and group related functionality.
type Module = module.Module

// IModule defines the interface for a module in the application.
// Modules are responsible for managing controllers and services and can include other submodules.
// Implementers of this interface must provide methods to get the module's name, controllers, and services.
type IModule = module.IModule

// Router represents an HTTP router with middleware support and error handling.
type Router = router.Router

// BaseService provides a default implementation of the Service interface.
// This struct can be embedded in other service implementations to inherit
// common functionalities or to be extended with custom methods.
type BaseService = service.BaseService

// Service defines the interface for all services in the application.
// Implementations of this interface can provide specific functionalities
// required by different parts of the application.
type Service = service.Service

// CORSOptions defines the configuration for the CORS middleware
type CORSOptions = middleware.CORSOptions

// Context holds the request and response writer and provides utility methods.
type Context = context.Context

type WebSocketServer = websocket.WebSocketServer

// Expose middleware types and functions

// CORSMiddleware is the middleware that handles CORS
type CORSMiddleware = middleware.CORSMiddleware

type RateLimiterMiddleware = middleware.RateLimiter
type FileUploadMiddleware = middleware.FileUploadMiddleware
type Config = config.Config

// VARS
var (
	app = router.GetApp()
)

func GetApp() *Router {
	return app
}

// LoadConfig loads the ENV configurations
func LoadConfig() config.Config {
	config := config.LoadConfig()
	return config
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return di.NewContainer()
}

// NewModule creates a new module
func NewModule(name string, controllers []interface{}, services []interface{}, submodules []IModule) *Module {
	return module.NewModule(name, controllers, services, submodules)
}

// NewRouter creates a new Router with optional configuration
func NewRouter(options ...router.Option) *Router {
	return router.NewRouter(options...)
}

// App creates a new app with optional configuration. You can pass options like WithCORS or WithJSONParser to configure the app.
func App(options ...router.Option) *Router {
	return router.NewRouter(options...)
}

// New Cors Options.
//
// Example (default to)
//
//	 corsOptions := LessGo.NewCorsOptions(
//		[]string{"*"}, // Allow all origins
//		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed methods
//		[]string{"Content-Type", "Authorization"},           // Allowed headers
//
// )
func NewCorsOptions(origins []string, methods []string, headers []string) *CORSOptions {
	var defCorsOpts = middleware.CORSOptions{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}
	if len(origins) == 0 {
		origins = append(origins, defCorsOpts.AllowedOrigins...)
	}
	if len(headers) == 0 {
		headers = append(headers, defCorsOpts.AllowedHeaders...)
	}
	if len(methods) == 0 {
		methods = append(methods, defCorsOpts.AllowedMethods...)
	}
	return middleware.NewCorsOptions(origins, methods, headers)
}

// WithCORS enables CORS middleware with specific options.
// This option configures the CORS settings for the router.
//
// Example usage:
//
//	r := router.NewRouter(router.WithCORS(middleware.CORSOptions{...}))
func WithCORS(options middleware.CORSOptions) router.Option {
	return router.WithCORS(options)
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
func WithInMemoryRateLimiter(NumShards int, Limit int, Interval time.Duration, CleanupInterval time.Duration) router.Option {
	return router.WithInMemoryRateLimiter(NumShards, Limit, Interval, CleanupInterval)
}

// WithRateLimiter enables rate limiting middleware with the specified limit and interval.
// This option configures the rate limiter for the router.
//
// Example usage:
//
//	r := router.NewRouter(router.WithRateLimiter(100, time.Minute))
func WithRedisRateLimiter(client *redis.Client, limit int, interval time.Duration) router.Option {
	return router.WithRedisRateLimiter(client, limit, interval)
}

type ParserOptions = middleware.ParserOptions

// Parser options. set default size
func NewParserOptions(size int64) *ParserOptions {
	return middleware.NewParserOptions(size)
}

// WithJSONParser enables JSON parsing middleware for request bodies.
// This option ensures that incoming JSON payloads are parsed and available in the request context.
//
// Example usage:
//
//	r := router.NewRouter(router.WithJSONParser())
func WithJSONParser(options ParserOptions) router.Option {
	return router.WithJSONParser(options)
}

// WithCookieParser enables cookie parsing middleware.
// This option ensures that cookies are parsed and available in the request context.
//
// Example usage:
//
//	r := router.NewRouter(router.WithCookieParser())
func WithCookieParser() router.Option {
	return router.WithCookieParser()
}

// WithFileUpload enables file upload middleware with the specified upload directory.
// This option configures the router to handle file uploads and save them to the given directory.
//
// Example usage:
//
//	r := router.NewRouter(router.WithFileUpload("/uploads"))
func WithFileUpload(uploadDir string, maxFileSize int64, allowedExts []string) router.Option {
	return router.WithFileUpload(uploadDir, maxFileSize, allowedExts)
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
func WithCaching(redisClient *redis.Client, ttl time.Duration, cacheControl bool) router.Option {
	return router.WithCaching(redisClient, ttl, cacheControl)
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
func WithCsrf() router.Option {
	return router.WithCsrf()
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
func WithXss() router.Option {
	return router.WithXss()
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
func WithTemplateRendering(templateDir string) router.Option {
	return router.WithTemplateRendering(templateDir)
}

func RegisterModules(r *router.Router, modules []module.IModule) error {
	return di.RegisterModules(r, modules)
}

func RegisterDependencies(dependencies []interface{}) {
	di.RegisterDependencies(dependencies)
}

// Resolves the path of specified folder
func GetFolderPath(folderName string) string {
	return utils.GetFolderPath(folderName)
}

func GenerateRandomToken(len int) (string, error) {
	return utils.GenerateRandomToken(len)
}

func DiscoverModules() ([]func() IModule, error) {
	return discovery.DiscoverModules()
}

func NewWebSocketServer() *WebSocketServer {
	return websocket.NewWebSocketServer()
}

// TASKS
type TaskBuilder = concurrency.TaskBuilder

const Parallel = 0
const Sequential = 1

func NewTaskBuilder(mode int) *TaskBuilder {
	return concurrency.NewTaskBuilder(concurrency.ExecutionMode(mode))
}

type SizeUnit string

const (
	Bytes     SizeUnit = "bytes"
	Kilobytes SizeUnit = "kilobytes"
	Megabytes SizeUnit = "megabytes"
	Gigabytes SizeUnit = "gigabytes"
)

// Convert size to bytes
//
// # Example
//
// const (
//
//	Bytes     SizeUnit = "bytes"
//	Kilobytes SizeUnit = "kilobytes"
//	Megabytes SizeUnit = "megabytes"
//	Gigabytes SizeUnit = "gigabytes"
//
// )
func ConvertToBytes(size int64, unit SizeUnit) int64 {
	s, err := utils.ConvertToBytes(float64(size), utils.SizeUnit(unit))
	if err != nil {
		log.Fatalf("Failed to convert bytes: %v", err)
	}
	return int64(s)
}

func NewRedisClient(redisAddr string) *redis.Client {
	return utils.NewRedisClient(redisAddr)
}
