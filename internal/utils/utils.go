// Package for internal utilities.
package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func IsValidCommand(cmd string) bool {
	validCommands := [...]string{
		"check", "doc", "generate", "help", "version",
	}

	for i, _ := range validCommands {
		if cmd == validCommands[i] {
			return true
		}
	}

	return false
}

func GetFilePathsByExtension(ext string) ([]string, error) {
	files := []string{}

	workDir, err := os.Getwd()
	if err != nil {
		return files, fmt.Errorf("get file paths by '%s' extension: %v", ext, err)
	}

	err = filepath.Walk(workDir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return files, fmt.Errorf("get file paths by '%s' extension: %v", ext, err)
	}

	return files, nil
}
