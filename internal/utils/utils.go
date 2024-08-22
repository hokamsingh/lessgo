package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
func RegisterModuleRoutes(container *di.Container, r *router.Router, _ interface{}) {
	err := container.Invoke(func(module *module.Module) {
		for _, ctrl := range module.Controllers {
			c, ok := ctrl.(controller.Controller)
			if !ok {
				panic(fmt.Sprintf("Controller %T does not implement controller.Controller interface", ctrl))
			}
			c.RegisterRoutes(r)
		}
	})
	if err != nil {
		panic(fmt.Sprintf("Container invocation failed: %v", err))
	}
}

// RegisterModules iterates over a slice of modules and registers their routes.
func RegisterModules(r *router.Router, container *di.Container, modules []module.Module) error {
	for _, module := range modules {
		RegisterModuleRoutes(container, r, module)
		log.Print("LessGo :: Registered module ", module.Name)
	}
	return nil
}
