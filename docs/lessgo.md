
## Package Overview

### Configuration

- **`LessGo.LoadConfig()`**: Loads the configuration settings. Typically used to configure server parameters like port and environment.
- **`cfg.GetInt(key, default)`**: Retrieves an integer configuration value.
- **`cfg.Get(key, default)`**: Retrieves a string configuration value.

### Middleware and Options

- **`LessGo.NewCorsOptions(origins, methods, headers)`**: Creates new CORS options for handling cross-origin requests.
- **`LessGo.NewParserOptions(maxSize)`**: Configures options for JSON parsing, including maximum size of request bodies.
- **`LessGo.NewRedisClient(address)`**: Creates a new Redis client instance.
- **`LessGo.WithCORS(options)`**: Adds CORS middleware with the provided options.
- **`LessGo.WithJSONParser(options)`**: Adds JSON parsing middleware with specified options.
- **`LessGo.WithCookieParser()`**: Adds middleware for parsing cookies.
- **`LessGo.WithCsrf()`**: Adds CSRF protection middleware.
- **`LessGo.WithXss()`**: Adds XSS protection middleware.
- **`LessGo.WithCaching(client, duration, enable)`**: Adds caching middleware using Redis.
- **`LessGo.WithRedisRateLimiter(address, limit, duration)`**: Adds rate limiting middleware with Redis.

### Application Initialization

- **`LessGo.App(middlewares...)`**: Initializes a new application instance with the provided middlewares.
- **`App.ServeStatic(path, folderPath)`**: Configures the application to serve static files from a specified folder.
- **`LessGo.RegisterDependencies(dependencies)`**: Registers dependencies for dependency injection.
- **`LessGo.RegisterModules(app, modules)`**: Registers application modules with the framework.

### Routes and Server

- **`App.Get(route, handler)`**: Registers a GET route with a specified handler.
- **`App.Listen(address)`**: Starts the server and listens on the specified address.

### Example Usage
```go
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
	serverPort := cfg.GetInt("SERVER_PORT", 8080)
	env := cfg.Get("ENV", "development")
	addr := ":" + string(rune(serverPort))

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
	if err := App.Listen(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
```
In the example provided, the server is configured to use various middleware options, serve static files, and register dependencies and modules. The server listens on a specified port and handles requests with CSRF protection and other middleware features.
```