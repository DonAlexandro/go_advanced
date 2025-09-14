package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetTxtFiles returns a list of all .txt file paths in the given directory
func GetTxtFiles(directoryPath string) ([]string, error) {
	var txtFiles []string

	// Check if directory exists
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", directoryPath)
	}

	// Walk through the directory
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a file and has .txt extension
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".txt" {
			txtFiles = append(txtFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	return txtFiles, nil
}
