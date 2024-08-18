package utils

import (
	"os"
	"path/filepath"
)

func GetFolderPath(folderName string) (string, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Join the CWD with the folder name
	folderPath := filepath.Join(cwd, folderName)

	// Check if the folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return "", err
	}

	return folderPath, nil
}
