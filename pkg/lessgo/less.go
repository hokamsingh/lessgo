package LessGo

import (
	"time"

	"github.com/hokamsingh/lessgo/internal/core/config"
	"github.com/hokamsingh/lessgo/internal/core/context"
	"github.com/hokamsingh/lessgo/internal/core/controller"
	"github.com/hokamsingh/lessgo/internal/core/di"
	"github.com/hokamsingh/lessgo/internal/core/middleware"
	"github.com/hokamsingh/lessgo/internal/core/module"
	"github.com/hokamsingh/lessgo/internal/core/router"
	"github.com/hokamsingh/lessgo/internal/core/service"
	"github.com/hokamsingh/lessgo/internal/utils"
)

// Expose core types
type Controller = controller.Controller
type BaseController = controller.BaseController
type Container = di.Container
type Middleware = middleware.Middleware
type BaseMiddleware = middleware.BaseMiddleware
type Module = module.Module
type IModule = module.IModule
type Router = router.Router
type BaseService = service.BaseService
type Service = service.Service
type CORSOptions = middleware.CORSOptions
type Context = context.Context

// Expose middleware types and functions
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
