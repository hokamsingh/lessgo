package middleware

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type FileUploadMiddleware struct {
	uploadDir   string
	maxFileSize int64    // Maximum file size in bytes
	allowedExts []string // Allowed file extensions
}

// NewFileUploadMiddleware creates a new instance of FileUploadMiddleware
func NewFileUploadMiddleware(uploadDir string, maxFileSize int64, allowedExts []string) *FileUploadMiddleware {
	// Ensure the upload directory exists
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	if len(allowedExts) == 0 {
		allowedExts = []string{".jpg"} // Default allowed extension if none provided
	}

	return &FileUploadMiddleware{
		uploadDir:   uploadDir,
		maxFileSize: maxFileSize,
		allowedExts: allowedExts,
	}
}

// Handle is the middleware function that processes file uploads
func (f *FileUploadMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(f.maxFileSize); err != nil {
			http.Error(w, "File too large or unable to parse form", http.StatusBadRequest)
			log.Printf("Error parsing form: %v", err)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to get file from form data", http.StatusBadRequest)
			log.Printf("Error retrieving file: %v", err)
			return
		}
		defer file.Close()

		// Validate file size
		if fileHeader.Size > f.maxFileSize {
			http.Error(w, "File size exceeds limit", http.StatusRequestEntityTooLarge)
			return
		}

		// Validate file extension
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !f.isAllowedExt(ext) {
			http.Error(w, "File type not allowed", http.StatusUnsupportedMediaType)
			return
		}

		// Generate a unique file name
		fileName := generateFileName() + ext
		filePath := filepath.Join(f.uploadDir, fileName)

		// Create the file
		cleanFilePath := filepath.Clean(filePath)
		if !strings.HasPrefix(cleanFilePath, f.uploadDir) {
			log.Panic("invalid file path")
			log.Printf("Error creating file: %v", err)
			return
		}

		destFile, err := os.Create(cleanFilePath)
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			log.Printf("Error creating file: %v", err)
			return
		}
		defer destFile.Close()

		// Copy file content
		if _, err := io.Copy(destFile, file); err != nil {
			http.Error(w, "Unable to copy file content", http.StatusInternalServerError)
			log.Printf("Error copying file content: %v", err)
			return
		}

		// Optionally, you can add a response to inform about the successful upload
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "File uploaded successfully: %s", fileName)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// isAllowedExt checks if the file extension is allowed
func (f *FileUploadMiddleware) isAllowedExt(ext string) bool {
	for _, allowedExt := range f.allowedExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// generateFileName generates a unique file name using UUID
func generateFileName() string {
	return uuid.New().String()
}
