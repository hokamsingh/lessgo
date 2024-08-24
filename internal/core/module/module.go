/*
Package module provides the definition of modules, which can encapsulate controllers and services in an application.

This package defines the `IModule` interface and a `Module` struct that allows for organizing and managing controllers, services, and submodules. It provides methods to retrieve the name, controllers, and services associated with a module.

Usage:

	import (
		"github.com/hokamsingh/lessgo/pkg/lessgo/module"
	)

	func main() {
		ctrl1 := &MyController{}
		svc1 := &MyService{}

		mod := module.NewModule(
			"MyModule",
			[]interface{}{ctrl1},
			[]interface{}{svc1},
			nil,
		)

		fmt.Println(mod.GetName())           // Outputs: MyModule
		fmt.Println(mod.GetControllers())    // Outputs: [<controller>]
		fmt.Println(mod.GetServices())       // Outputs: [<service>]
	}
*/
package module

// IModule defines the interface for a module in the application.
// Modules are responsible for managing controllers and services and can include other submodules.
// Implementers of this interface must provide methods to get the module's name, controllers, and services.
type IModule interface {
	GetName() string
	GetControllers() []interface{}
	GetServices() []interface{}
}

// Module represents a module in the application.
// It holds the name, a list of controllers, services, and any submodules.
// The module can be used to organize and group related functionality.
type Module struct {
	Name        string
	submodules  []IModule
	Controllers []interface{}
	Services    []interface{}
}

// NewModule creates a new instance of `Module` with the specified name, controllers, services, and submodules.
//
// Example:
//
//	ctrl1 := &MyController{}
//	svc1 := &MyService{}
//
//	mod := module.NewModule(
//		"MyModule",
//		[]interface{}{ctrl1},
//		[]interface{}{svc1},
//		nil,
//	)
//
//	fmt.Println(mod.GetName())           // Outputs: MyModule
//	fmt.Println(mod.GetControllers())    // Outputs: [<controller>]
//	fmt.Println(mod.GetServices())       // Outputs: [<service>]
func NewModule(name string, controllers []interface{}, services []interface{}, submodules []IModule) *Module {
	return &Module{
		Name:        name,
		Controllers: controllers,
		Services:    services,
		submodules:  submodules,
	}
}

// GetName returns the name of the module.
//
// Example:
//
//	mod := module.NewModule("MyModule", nil, nil, nil)
//	fmt.Println(mod.GetName()) // Outputs: MyModule
func (m *Module) GetName() string {
	return m.Name
}

// GetControllers returns a list of controllers associated with the module.
//
// Example:
//
//	ctrl1 := &MyController{}
//	mod := module.NewModule("MyModule", []interface{}{ctrl1}, nil, nil)
//	fmt.Println(mod.GetControllers()) // Outputs: [<controller>]
func (m *Module) GetControllers() []interface{} {
	return m.Controllers
}

// GetServices returns a list of services associated with the module.
//
// Example:
//
//	svc1 := &MyService{}
//	mod := module.NewModule("MyModule", nil, []interface{}{svc1}, nil)
//	fmt.Println(mod.GetServices()) // Outputs: [<service>]
func (m *Module) GetServices() []interface{} {
	return m.Services
}
