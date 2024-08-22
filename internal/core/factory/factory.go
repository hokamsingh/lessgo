package factory

import (
	"github.com/hokamsingh/lessgo/internal/core/di"
	"github.com/hokamsingh/lessgo/internal/core/router"
)

// App represents the main application structure
type App struct {
	Router    *router.Router
	Container *di.Container
}

// NewApp creates a new App instance
func NewApp(router *router.Router, container *di.Container) *App {
	return &App{
		Router:    router,
		Container: container,
	}
}

// Start the HTTP server on the specified address
func (app *App) Start(addr string) error {
	return app.Router.Listen(addr)
}
