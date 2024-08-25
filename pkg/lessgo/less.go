package LessGo

import (
	"time"

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

// VARS
var (
	app = router.GetApp()
)

func GetApp() *Router {
	return app
}

// LoadConfig loads the configuration
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

func App(options ...router.Option) *Router {
	return router.NewRouter(options...)
}

// New Cors Options init
func NewCorsOptions(origins []string, methods []string, headers []string) *CORSOptions {
	return middleware.NewCorsOptions(origins, methods, headers)
}

// Expose middleware options
func WithCORS(options middleware.CORSOptions) router.Option {
	return router.WithCORS(options)
}

func WithRateLimiter(limit int, interval time.Duration) router.Option {
	return router.WithRateLimiter(limit, interval)
}

func WithJSONParser() router.Option {
	return router.WithJSONParser()
}

func WithCookieParser() router.Option {
	return router.WithCookieParser()
}

func WithFileUpload(uploadDir string) router.Option {
	return router.WithFileUpload(uploadDir)
}

// // ServeStatic creates a file server handler to serve static files
// func ServeStatic(pathPrefix, dir string) http.Handler {
// 	return router.ServeStatic(pathPrefix, dir)
// }

func GetFolderPath(folderName string) (string, error) {
	return utils.GetFolderPath(folderName)
}

func RegisterModuleRoutes(r *router.Router, module module.Module) {
	utils.RegisterModuleRoutes(r, &module)
}

// RegisterModules iterates over a slice of modules and registers their routes.
func RegisterModules(r *router.Router, modules []IModule) error {
	return utils.RegisterModules(r, modules)
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

// RegisterDependencies registers dependencies into container
func RegisterDependencies(dependencies []interface{}) {
	utils.RegisterDependencies(dependencies)
}
