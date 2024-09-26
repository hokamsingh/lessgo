package src

import (
	"github.com/hokamsingh/lessgo/app/src/test"
	"github.com/hokamsingh/lessgo/app/src/upload"
	user "github.com/hokamsingh/lessgo/app/src/user"
	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type RootModule struct {
	LessGo.Module
}

func NewRootModule(r *LessGo.Router) *RootModule {
	// Initialize and collect all modules
	modules := []LessGo.IModule{
		test.NewTestModule(),
		upload.NewUploadModule(),
		user.NewUserModule(),
	}

	// Register all modules
	LessGo.RegisterModules(r, modules)
	service := NewRootService()
	controller := NewRootController(service, "/")
	return &RootModule{
		Module: *LessGo.NewModule("Root", []interface{}{controller}, []interface{}{service}, modules),
	}
}
