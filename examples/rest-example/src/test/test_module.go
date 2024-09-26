package test

import (
	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type TestModule struct {
	LessGo.Module
}

// NewTestModule creates a new instance of TestModule
func NewTestModule() *TestModule {
	testService := NewTestService()
	testController := NewTestController(testService, "/test")

	return &TestModule{
		Module: *LessGo.NewModule( // You need to initialize the embedded Module field
			"Test",                        // Name of the module
			[]interface{}{testController}, // Controllers
			[]interface{}{testService},    // Services
			[]LessGo.IModule{},
		),
	}
}
