package test

import (
	"github.com/hokamsingh/lessgo/internal/core/module"
	core "github.com/hokamsingh/lessgo/pkg/lessgo"
)

func NewTestModule() *core.Module {
	// Create the service first
	testService := NewTestService()

	// Create the controller with the service
	testController := NewTestController(testService)

	// Return the module with the controller and service
	return module.NewModule("ExampleModule",
		[]interface{}{testController}, // Controllers
		[]interface{}{testService},    // Services
	)
}
