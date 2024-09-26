// main.go

// @title My API
// @version 1.0
// @description This is a sample server.
// @host localhost:8080
// @BasePath /

package main

import (
	"log"
	"time"

	"github.com/hokamsingh/lessgo/app/src"
	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

func main() {
	// Load Configuration
	cfg := LessGo.LoadConfig()
	serverPort := cfg.Get("SERVER_PORT", "8080")
	env := cfg.Get("ENV", "development")
	addr := ":" + serverPort

	// CORS Options
	corsOptions := LessGo.NewCorsOptions(
		[]string{"*"}, // Allow all origins
		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed methods
		[]string{"Content-Type", "Authorization"},           // Allowed headers
	)

	// Parser Options
	size := LessGo.ConvertToBytes(int64(1024), LessGo.Kilobytes)
	parserOptions := LessGo.NewParserOptions(size * 5)

	// redis client
	rClient := LessGo.NewRedisClient("localhost:6379")

	// Initialize App with Middlewares
	App := LessGo.App(
		LessGo.WithCORS(*corsOptions),
		// LessGo.WithInMemoryRateLimiter(4, 50, 1*time.Second, 5*time.Minute), // Rate limiter
		// LessGo.WithRedisRateLimiter("localhost:6379", 10, time.Minute*5),
		LessGo.WithJSONParser(*parserOptions),
		LessGo.WithCookieParser(),                        // Cookie parser
		LessGo.WithCsrf(),                                // CSRF protection middleware
		LessGo.WithXss(),                                 // XSS protection middleware
		LessGo.WithCaching(rClient, 5*time.Minute, true), // Caching middleware using Redis
		LessGo.WithRedisRateLimiter(rClient, 100, 1*time.Second),
		// LessGo.WithFileUpload("uploads"), // Uncomment if you want to handle file uploads
	)

	// Serve Static Files
	folderPath := LessGo.GetFolderPath("uploads")
	App.ServeStatic("/static/", folderPath)

	// Register dependencies
	dependencies := []interface{}{src.NewRootService, src.NewRootModule}
	LessGo.RegisterDependencies(dependencies)

	// Root Module
	rootModule := src.NewRootModule(App)
	LessGo.RegisterModules(App, []LessGo.IModule{rootModule})

	// Example Route
	App.Get("/ping", func(ctx *LessGo.Context) {
		ctx.Send("pong")
	})

	// Start the server
	log.Printf("Starting server on port %s in %s mode", serverPort, env)
	// LessGo.PProfiling()
	httpCfg := LessGo.NewHttpConfig()
	if err := App.Listen(addr, httpCfg); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
