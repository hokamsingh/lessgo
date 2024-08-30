package discovery

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/hokamsingh/lessgo/internal/core/module"
)

// DiscoverModules scans the src directory for module files and compiles them into shared object files (.so).
func DiscoverModules() ([]func() module.IModule, error) {
	var modules []func() module.IModule

	// Define paths
	srcDir := "app/src"
	pluginDir := "app/plugins"

	// Ensure the plugin directory exists
	err := os.MkdirAll(pluginDir, 0750)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugins directory: %w", err)
	}

	// Walk through the src directory to find *_module.go files
	err = filepath.WalkDir(srcDir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), "_module.go") {
			moduleName := strings.TrimSuffix(info.Name(), ".go")
			moduleFile := path
			moduleSOFile := filepath.Join(pluginDir, moduleName+".so")

			// Compile the Go file into a shared object
			if err := compileModule(moduleFile, moduleSOFile); err != nil {
				return fmt.Errorf("error compiling module %s: %w", moduleName, err)
			}

			// Load the compiled module
			mod, err := loadModule(moduleSOFile, moduleName)
			if err != nil {
				return fmt.Errorf("error loading module %s: %w", moduleName, err)
			}
			modules = append(modules, mod)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error discovering modules: %w", err)
	}

	return modules, nil
}

// compileModule compiles a Go source file into a shared object (.so) file.
func compileModule(moduleFile, moduleSOFile string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", moduleSOFile, moduleFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error compiling module %s: %s", moduleFile, string(output))
	}
	return nil
}

// loadModule dynamically loads a compiled module (.so file).
func loadModule(moduleSOFile, moduleName string) (func() module.IModule, error) {
	plg, err := plugin.Open(moduleSOFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open shared object file: %w", err)
	}

	// Look up the module's initialization function
	sym, err := plg.Lookup("New" + moduleName)
	if err != nil {
		return nil, fmt.Errorf("failed to find symbol for module %s: %w", moduleName, err)
	}

	// Assert the symbol to a function type
	initFunc, ok := sym.(func() module.IModule)
	if !ok {
		return nil, errors.New("symbol is not of expected function type")
	}

	return initFunc, nil
}
