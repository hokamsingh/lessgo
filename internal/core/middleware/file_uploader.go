package middleware

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// FileUploadMiddleware handles file uploads and saves them to the specified directory
type FileUploadMiddleware struct {
	uploadDir string
}

// NewFileUploadMiddleware creates a new instance of FileUploadMiddleware
func NewFileUploadMiddleware(uploadDir string) *FileUploadMiddleware {
	return &FileUploadMiddleware{uploadDir: uploadDir}
}

// Handle is the middleware function that processes file uploads
func (f *FileUploadMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32 MB limit
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get file from form data
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to get file from form data", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Extract file extension
		ext := filepath.Ext(fileHeader.Filename)
		if ext == "" {
			ext = ".bin" // Default extension if none provided
		}

		// Create file on server with timestamp and original extension
		filePath := filepath.Join(f.uploadDir, generateFileName()+ext)
		destFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
		defer destFile.Close()

		// Copy file content
		if _, err := io.Copy(destFile, file); err != nil {
			http.Error(w, "Unable to copy file content", http.StatusInternalServerError)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// generateFileName generates a unique file name with timestamp
func generateFileName() string {
	return time.Now().Format("20060102150405")
}
