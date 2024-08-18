package LessGo

import (
	"time"

	"github.com/hokamsingh/lessgo/internal/core/config"
	"github.com/hokamsingh/lessgo/internal/core/controller"
	"github.com/hokamsingh/lessgo/internal/core/di"
	"github.com/hokamsingh/lessgo/internal/core/middleware"
	"github.com/hokamsingh/lessgo/internal/core/module"
	"github.com/hokamsingh/lessgo/internal/core/router"
	"github.com/hokamsingh/lessgo/internal/core/service"
)

// Expose core types
type Controller = controller.Controller
type BaseController = controller.BaseController
type Container = di.Container
type Middleware = middleware.Middleware
type BaseMiddleware = middleware.BaseMiddleware
type Module = module.Module
type Router = router.Router
type BaseService = service.BaseService
type Service = service.Service
type CORSOptions = middleware.CORSOptions

// Expose middleware types and functions
type CORSMiddleware = middleware.CORSMiddleware
type RateLimiterMiddleware = middleware.RateLimiter
type FileUploadMiddleware = middleware.FileUploadMiddleware

// LoadConfig loads the configuration
func LoadConfig() config.Config {
	config := config.LoadConfig()
	return *config
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return di.NewContainer()
}

// NewModule creates a new module
func NewModule(name string, controllers []interface{}, services []interface{}) *Module {
	return module.NewModule(name, controllers, services)
}

// NewRouter creates a new Router with optional configuration
func NewRouter(options ...router.Option) *Router {
	return router.NewRouter(options...)
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
