/*
Package controller provides a base structure and interface for defining and registering routes in the application.

This package defines the Controller interface that all controllers must implement to register their routes,
as well as a BaseController struct that provides a default implementation of the interface.
*/
package controller

import (
	"github.com/hokamsingh/lessgo/internal/core/router"
)

// Controller defines the interface that all controllers in the application must implement.
// Any controller that implements this interface must define the RegisterRoutes method,
// which is responsible for setting up the necessary routes for the controller.
type Controller interface {
	RegisterRoutes(r *router.Router)
}

// BaseController provides a default implementation of the Controller interface.
// It can be embedded in other controllers to inherit its default behavior,
// or overridden with custom implementations.
type BaseController struct {
}

// RegisterRoutes is the default implementation of the Controller interface's method.
// This method can be overridden by embedding BaseController in another struct
// and defining a custom implementation.
//
// Example
//
//	type TestController struct {
//		LessGo.BaseController
//		Path    string
//		Service TestService
//	}
//
//	func NewTestController(service *TestService, path string) *TestController {
//		return &TestController{
//			Service: *service,
//			Path:    path,
//		}
//	}
//
//	func (tc *TestController) RegisterRoutes(r *LessGo.Router) {
//		tr := r.SubRouter(tc.Path)
//		tr.Get("/ping", func(ctx *LessGo.Context) {
//			ctx.Send("pong")
//		})
//	}
func (bc *BaseController) RegisterRoutes(r *router.Router) {

}
