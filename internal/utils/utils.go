package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
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
		Addr:     redisAddr, // e.g., "localhost:6379"
		Password: "secret",
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	return client
}

// GenerateSalt creates a random salt of the given length.
func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// HashPassword generates a hash for the given data using the provided salt.
func HashPassword(data string, salt string, length int) (string, error) {
	if len(salt) == 0 {
		return "", errors.New("salt cannot be empty")
	}
	hashed := fmt.Sprintf("%x", data+salt) // Placeholder, replace with real hash function
	return hashed[:length], nil
}

// GenerateRandomToken creates a random token of the specified length.
func GenerateRandomToken(length int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	token := make([]byte, length)
	for i := range token {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		token[i] = letters[num.Int64()]
	}
	return string(token), nil
}

// MsToDay converts milliseconds to days.
func MsToDay(ms int64) int64 {
	return ms / (24 * 60 * 60 * 1000)
}

// MsToMin converts milliseconds to minutes.
func MsToMin(ms int64) int64 {
	return ms / (60 * 1000)
}

// MsToHr converts milliseconds to hours.
func MsToHr(ms int64) int64 {
	return ms / (60 * 60 * 1000)
}

// MsToSec converts milliseconds to seconds.
func MsToSec(ms int64) int64 {
	return ms / 1000
}

// MsToHuman converts milliseconds to a human-readable string.
func MsToHuman(ms int64, maxUnit string) string {
	var days, hours, minutes, seconds int64
	switch maxUnit {
	case "day":
		days = MsToDay(ms)
		hours = MsToHr(ms) % 24
		minutes = MsToMin(ms) % 60
		seconds = MsToSec(ms) % 60
	case "hour":
		hours = MsToHr(ms)
		minutes = MsToMin(ms) % 60
		seconds = MsToSec(ms) % 60
	case "minute":
		minutes = MsToMin(ms)
		seconds = MsToSec(ms) % 60
	case "second":
		seconds = MsToSec(ms)
	}

	return fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", days, hours, minutes, seconds)
}

// Sleep pauses the execution for the given number of milliseconds.
func Sleep(ms int64) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// Retryable retries the provided function on failure with backoff.
func Retryable(fn func() error, retries int, backoffType string, delay time.Duration) error {
	var err error
	for i := 0; i <= retries; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if backoffType == "exponential" {
			time.Sleep(delay * time.Duration(math.Pow(2, float64(i))))
		} else {
			time.Sleep(delay)
		}
	}
	return err
}

// GenerateRange creates a range of numbers from start to end.
func GenerateRange(start, end int) []int {
	arr := make([]int, end-start+1)
	for i := range arr {
		arr[i] = start + i
	}
	return arr
}

// GetRandomIndex generates a random index between min and max.
func GetRandomIndex(min, max int) (int, error) {
	randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		return 0, err
	}
	return int(randIndex.Int64()) + min, nil
}

// ShuffleNumbers randomly shuffles a slice of numbers.
func ShuffleNumbers(numbers []int) ([]int, error) {
	shuffled := make([]int, len(numbers))
	copy(shuffled, numbers)
	for i := len(shuffled) - 1; i > 0; i-- {
		j, err := GetRandomIndex(0, i+1)
		if err != nil {
			return nil, err
		}
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled, nil
}
