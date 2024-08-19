package utils

import (
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

// registerModuleRoutes is a helper function to register routes for a module.
func RegisterModuleRoutes(container *di.Container, r *router.Router, _ interface{}) error {
	return container.Invoke(func(module *module.Module) {
		for _, ctrl := range module.Controllers {
			if c, ok := ctrl.(controller.Controller); ok {
				c.RegisterRoutes(r.Mux)
			}
		}
	})
}
