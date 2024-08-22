package controller

import (
	"github.com/hokamsingh/lessgo/internal/core/router"
)

// BASE
type Controller interface {
	RegisterRoutes(r *router.Router)
}

type BaseController struct{}

func (bc *BaseController) RegisterRoutes(r *router.Router) {

}
