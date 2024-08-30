package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"unicode"

	"github.com/go-redis/redis/v8"
)

func GetFolderPath(folderName string) string {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicf("Failed to get folder path: %v", err)
		return ""
	}

	// Join the CWD with the folder name
	folderPath := filepath.Join(cwd, folderName)

	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		log.Panicf("folder does not exists: %v", err)
		return ""
	}
	return folderPath
}

// func RegisterModuleRoutes(container *di.Container, r *router.Router, _ interface{}) {
// 	err := container.Invoke(func(module module.IModule) {
// 		for _, ctrl := range module.GetControllers() {
// 			c, ok := ctrl.(controller.Controller)
// 			if !ok {
// 				panic(fmt.Sprintf("Controller %T does not implement controller.Controller interface", ctrl))
// 			}
// 			c.RegisterRoutes(r)
// 		}
// 	})
// 	if err != nil {
// 		panic(fmt.Sprintf("Container invocation failed: %v", err))
// 	}
// }

func init() {
	log.SetFlags(0)
	log.SetOutput(&logWriter{})
}

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	// Adding color codes to the timestamp
	coloredTimestamp := fmt.Sprintf("\033[1;33m%s\033[0m", timestamp) // Yellow color for the timestamp
	return fmt.Printf("%s %s", coloredTimestamp, bytes)
}

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Purple  = "\033[35m"
	SkyBlue = "\033[36m"
)

// GenerateRandomToken generates a random unique token of the specified length in bytes
func GenerateRandomToken(length int) (string, error) {
	// Create a byte slice to hold the random data
	token := make([]byte, length)

	// Fill the byte slice with random data
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate random token: %v", err)
	}

	// Convert the random bytes to a hexadecimal string
	return hex.EncodeToString(token), nil
}

type SizeUnit string

const (
	Bytes     SizeUnit = "bytes"
	Kilobytes SizeUnit = "kilobytes"
	Megabytes SizeUnit = "megabytes"
	Gigabytes SizeUnit = "gigabytes"
)

// ConvertToBytes converts a size given in the specified unit to bytes.
func ConvertToBytes(size float64, unit SizeUnit) (int64, error) {
	switch strings.ToLower(string(unit)) {
	case string(Bytes):
		return int64(size), nil
	case string(Kilobytes):
		return int64(size * 1024), nil
	case string(Megabytes):
		return int64(size * 1024 * 1024), nil
	case string(Gigabytes):
		return int64(size * 1024 * 1024 * 1024), nil
	default:
		return 0, errors.New("invalid size unit")
	}
}

func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func EscapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func Assert(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func NewRedisClient(redisAddr string) *redis.Client {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr, // e.g., "localhost:6379"
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	return client
}
