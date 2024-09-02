package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config represents a map of configuration key-value pairs loaded from the environment.
type Config map[string]string

// LoadConfig loads environment variables into a Config map. It first attempts to load a `.env` file
// using the `godotenv` package. If no `.env` file is found, it logs a message but continues to load
// environment variables from the system.
func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found in current directory: %v", err)

		// Attempt to find the .env file in the root directory
		rootDir, err := findRootDir()
		if err != nil {
			log.Printf("Unable to determine the root directory: %v", err)
		} else {
			envFile := filepath.Join(rootDir, ".env")
			if err := godotenv.Load(envFile); err != nil {
				log.Printf("No .env file found in root directory: %v", err)
			} else {
				log.Printf(".env file loaded from root directory: %s", envFile)
			}
		}
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

// Get retrieves a string value from the Config map based on the provided key. If the key does not exist
// in the Config, the function returns the specified default value.
func (c Config) Get(key, defaultValue string) string {
	if value, exists := c[key]; exists {
		return value
	}
	return defaultValue
}

// GetInt retrieves an integer value from the Config map based on the provided key. If the key does not exist
// or cannot be converted to an integer, the function returns the specified default value.
func (c Config) GetInt(key string, defaultValue int) int {
	if value, exists := c[key]; exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Invalid integer for key %s: %v", key, value)
	}
	return defaultValue
}

// GetBool retrieves a boolean value from the Config map based on the provided key. If the key does not exist
// or cannot be converted to a boolean, the function returns the specified default value.
func (c Config) GetBool(key string, defaultValue bool) bool {
	if value, exists := c[key]; exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
		log.Printf("Invalid boolean for key %s: %v", key, value)
	}
	return defaultValue
}

// GetFloat64 retrieves a float64 value from the Config map based on the provided key. If the key does not exist
// or cannot be converted to a float64, the function returns the specified default value.
func (c Config) GetFloat64(key string, defaultValue float64) float64 {
	if value, exists := c[key]; exists {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
		log.Printf("Invalid float for key %s: %v", key, value)
	}
	return defaultValue
}

// Validate checks that all the provided keys are present in the Config map. If any key is missing, it logs
// a fatal error and exits the program. This ensures that required configuration is always set.
func (c Config) Validate(requiredKeys ...string) {
	for _, key := range requiredKeys {
		if _, exists := c[key]; !exists {
			log.Fatalf("Missing required environment variable: %s", key)
		}
	}
}

// Reload reloads the configuration from the environment variables and `.env` file. This can be useful
// if the environment variables might change during runtime and you need to refresh the configuration.
func (c *Config) Reload() {
	*c = LoadConfig()
}

// MergeWithDefaults merges the current configuration with a default configuration map. If a key exists
// in both, the current configuration's value is preserved.
func (c Config) MergeWithDefaults(defaults Config) Config {
	merged := make(Config)

	for k, v := range defaults {
		merged[k] = v
	}

	for k, v := range c {
		merged[k] = v
	}

	return merged
}

// FilterByPrefix returns a new Config map containing only the keys that start with the specified prefix.
// The prefix is removed from the keys in the returned map.
func (c Config) FilterByPrefix(prefix string) Config {
	filtered := make(Config)

	for k, v := range c {
		if strings.HasPrefix(k, prefix) {
			filtered[strings.TrimPrefix(k, prefix)] = v
		}
	}

	return filtered
}

// findRootDir attempts to find the root directory by walking up from the current directory
// until it finds a directory containing a .env file, or it reaches the system root.
func findRootDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(currentDir, ".env")); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir { // Reached the system root
			break
		}
		currentDir = parentDir
	}

	return "", os.ErrNotExist
}
