/*
Package config provides functionality to load and manage application configuration
from environment variables and .env files.

The package leverages the `godotenv` package to load environment variables from a `.env` file into
the application, allowing easy configuration management.

Usage:

	import "your/package/path/config"

	func main() {
		cfg := config.LoadConfig()
		port := cfg.Get("PORT", "8080")
		// Use the `port` variable...
	}
*/
package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config represents a map of configuration key-value pairs loaded from the environment.
type Config map[string]string

// LoadConfig loads environment variables into a Config map. It first attempts to load a `.env` file
// using the `godotenv` package. If no `.env` file is found, it logs a message but continues to load
// environment variables from the system.
// Example:
//
//	cfg := config.LoadConfig()
//	fmt.Println(cfg["PORT"])
func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	config := make(Config)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			config[pair[0]] = pair[1]
		}
	}

	return config
}

// Get retrieves a value from the Config map based on the provided key. If the key does not exist
// in the Config, the function returns the specified default value.
//
// Example:
//
//	cfg := config.LoadConfig()
//	port := cfg.Get("PORT", "8080")
//	fmt.Println("Server will run on port:", port)
func (c Config) Get(key, defaultValue string) string {
	if value, exists := c[key]; exists {
		return value
	}
	return defaultValue
}
