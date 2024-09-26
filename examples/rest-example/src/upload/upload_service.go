package upload

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type IUploadService interface{}

type UploadService struct {
	UploadDir string
	LessGo.BaseService
}

func NewUploadService(uploadDir string) *UploadService {
	return &UploadService{UploadDir: uploadDir}
}

func (s *UploadService) SaveFile(file http.File, fileName string) (string, error) {
	filePath := filepath.Join(s.UploadDir, fileName)
	cleanFilePath := filepath.Clean(filePath)
	if !strings.HasPrefix(cleanFilePath, s.UploadDir) {
		return "", errors.New("invalid file path")
	}
	destFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
