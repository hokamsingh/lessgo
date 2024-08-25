package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hokamsingh/lessgo/internal/core/controller"
	"github.com/hokamsingh/lessgo/internal/core/di"
	"github.com/hokamsingh/lessgo/internal/core/module"
	"github.com/hokamsingh/lessgo/internal/core/router"
)

func GetFolderPath(folderName string) (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Join the CWD with the folder name
	folderPath := filepath.Join(cwd, folderName)

	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return "", err
	}

	return folderPath, nil
}

// RegisterModuleRoutes is a helper function to register routes for a module.
// It will panic if there is an error during registration or if a controller does not implement the required interface.
func RegisterModuleRoutes(r *router.Router, m module.IModule) {
	for _, ctrl := range m.GetControllers() {
		c, ok := ctrl.(controller.Controller)
		if !ok {
			panic(fmt.Sprintf("Controller %T does not implement controller.Controller interface", ctrl))
		}
		c.RegisterRoutes(r)
	}
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

// RegisterModules iterates over a slice of modules and registers their routes.
func RegisterModules(r *router.Router, modules []module.IModule) error {
	for _, module := range modules {
		RegisterModuleRoutes(r, module)
		log.Println(fmt.Sprintf("%sLessGo :: Registered module %s%s%s", Green, Yellow, module.GetName(), Reset))
	}
	return nil
}

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

// RegisterDependencies registers dependencies into container
func RegisterDependencies(dependencies []interface{}) {
	container := di.NewContainer()
	for _, dep := range dependencies {
		if err := container.Register(dep); err != nil {
			log.Fatalf("Error registering dependencies: %v", err)
		}
	}
}
