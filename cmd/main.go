package main

import (
	"log"

	"github.com/hokamsingh/lessgo/app/middleware"
	"github.com/hokamsingh/lessgo/app/module/test"
	"github.com/hokamsingh/lessgo/internal/core/controller"
	core "github.com/hokamsingh/lessgo/pkg/lessgo"
)

func main() {
	// Load Configuration
	cfg := core.LoadConfig()

	// Initialize Router
	r := core.NewRouter()

	// new Container init
	container := core.NewContainer()

	// register services in container
	testService := test.NewTestService()
	container.Register("test.TestService", testService)

	// Register middleware
	loggingMiddleare := middleware.NewLoggingMiddleware()
	errorMiddleware := middleware.NewErrorHandleMiddleware()
	// jwtMiddleware := middleware.NewJWTMiddleware(cfg.JwtSecret)
	r.Use(errorMiddleware)
	r.Use(loggingMiddleare)
	// r.Use(jwtMiddleware)

	// Routes

	// Register Test Module Routes
	testModule := test.NewTestModule()
	for _, ctrl := range testModule.Controllers {
		if c, ok := ctrl.(controller.Controller); ok {
			// inject dependencies
			container.Inject(c)

			// register routes
			c.RegisterRoutes(r.Mux)
		}
	}

	// Start Server
	log.Printf("Starting server on port %s in %s mode", cfg.ServerPort, cfg.Env)
	if err := r.Start(":" + cfg.ServerPort); err != nil {
		panic(err)
	}
}
