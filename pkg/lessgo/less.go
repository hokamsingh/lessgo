package LessGo

import (
	"github.com/hokamsingh/lessgo/internal/core/config"
	"github.com/hokamsingh/lessgo/internal/core/controller"
	"github.com/hokamsingh/lessgo/internal/core/di"
	"github.com/hokamsingh/lessgo/internal/core/middleware"
	"github.com/hokamsingh/lessgo/internal/core/module"
	"github.com/hokamsingh/lessgo/internal/core/router"
	"github.com/hokamsingh/lessgo/internal/core/service"
)

type Controller = controller.Controller
type BaseController = controller.BaseController

func LoadConfig() config.Config {
	config := config.LoadConfig()
	return *config
}

type Container = di.Container

func NewContainer() *Container {
	return di.NewContainer()
}

type Middleware = middleware.Middleware
type BaseMiddleware = middleware.BaseMiddleware

type Module = module.Module

func NewModule(name string, controllers []interface{}, services []interface{}) *Module {
	return module.NewModule(name, controllers, services)
}

type Router = router.Router

func NewRouter(options []router.Option) *Router {
	return router.NewRouter(options...)
}

type BaseService = service.BaseService
type Service = service.Service
