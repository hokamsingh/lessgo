package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config map[string]string

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

func (c Config) Get(key, defaultValue string) string {
	if value, exists := c[key]; exists {
		return value
	}
	return defaultValue
}
