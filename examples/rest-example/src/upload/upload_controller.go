package upload

import (
	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type UploadController struct {
	Path    string
	Service UploadService
}

func NewUploadController(service *UploadService, path string) *UploadController {
	// if !LessGo.ValidatePath(path) {
	// 	log.Fatalf("Invalid path provided: %s", path)
	// }
	return &UploadController{Path: path, Service: *service}
}

func (uc *UploadController) RegisterRoutes(r *LessGo.Router) {
	size := int64(5 * 1024 * 1024) // 5 mb in bytes
	ur := r.SubRouter(uc.Path, LessGo.WithFileUpload("uploads", size, []string{".jpg", ".png"}))

	ur.Post("/files", func(ctx *LessGo.Context) {
		ctx.Send("file saved")
	})
}
