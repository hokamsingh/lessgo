package test

import "github.com/hokamsingh/lessgo/internal/core/module"

func NewTestModule() *module.Module {
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
