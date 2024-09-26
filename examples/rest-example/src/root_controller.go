package src

import LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"

type RootController struct {
	Path    string
	Service RootService
}

func NewRootController(s *RootService, path string) *RootController {
	return &RootController{
		Path:    path,
		Service: *s,
	}
}

func (rc *RootController) RegisterRoutes(r *LessGo.Router) {
	// r.Get("/hello", func(ctx *LessGo.Context) {
	// 	ctx.Send("Hello world")
	// })
}
