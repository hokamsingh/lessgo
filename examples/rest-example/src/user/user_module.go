package user

import (
	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type UserModule struct {
	LessGo.Module
}

func NewUserModule() *UserModule {
	userService := NewUserService()
	userController := NewUserController(userService, "/users")
	return &UserModule{
		Module: *LessGo.NewModule("User",
			[]interface{}{userController}, // Controllers
			[]interface{}{userService},    // Services
			[]LessGo.IModule{},
		),
	}
}
