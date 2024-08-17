package main

import (
	"lessgo/app/middleware"
	"lessgo/app/module/test"
	"lessgo/internal/core/config"
	"lessgo/internal/core/controller"
	"lessgo/internal/core/di"
	"lessgo/internal/core/router"
	"log"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Router
	r := router.NewRouter()

	// new Container init
	container := di.NewContainer()

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
