package controller

import (
	"github.com/gorilla/mux"
)

// BASE
type Controller interface {
	RegisterRoutes(mux *mux.Router)
}

type BaseController struct{}

func (bc *BaseController) RegisterRoutes(mux *mux.Router) {

}
