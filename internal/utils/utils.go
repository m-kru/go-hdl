// Package for internal utilities.
package utils

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
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

var validLangs [1]string = [...]string{"vhdl"}

func ValidLangs() [1]string {
	return validLangs
}

func IsValidLang(lang string) bool {
	for i, _ := range validLangs {
		if lang == validLangs[i] {
			return true
		}
	}

	return false
}

func GetFilePathsByExtension(ext string, workDir string) ([]string, error) {
	files := []string{}

	err := filepath.Walk(workDir, func(path string, info fs.FileInfo, err error) error {
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

// IsIgnoredVHDLFile returns true if given file should be ignored.
// For example, it may be a Xilinx encrypted file.
// In such case there is no point in analyzing its content.
func IsIgnoredVHDLFile(filepath string) bool {
	// Ignore Xilinx encrypted files.
	if strings.HasSuffix(filepath, "_rfs.vhd") {
		return true
	}
	return false
}
