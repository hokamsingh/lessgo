package config

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	ServerPort string
	Env        string
	JwtSecret  string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		Env:        getEnv("ENV", "development"),
		JwtSecret:  getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
