package upload

import (
	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type UploadModule struct {
	LessGo.Module
}

func NewUploadModule() *UploadModule {
	service := NewUploadService("uploads")
	controller := NewUploadController(service, "/upload")
	return &UploadModule{
		Module: *LessGo.NewModule("Upload",
			[]interface{}{controller},
			[]interface{}{service},
			[]LessGo.IModule{},
		),
	}
}
