package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func GetPackagePath(data interface{}) (string, error) {
	t := reflect.TypeOf(data)
	pkgPath := t.PkgPath()
	runningPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Try to find the module root by looking for go.mod
	modRoot := runningPath
	for {
		if _, err := os.Stat(filepath.Join(modRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(modRoot)
		if parent == modRoot {
			// Reached root, not found
			modRoot = runningPath
			break
		}
		modRoot = parent
	}

	// If pkgPath is empty (main), just use current dir
	if pkgPath == "" {
		return runningPath, nil
	}

	// Try to find the relative path from module root to the package
	rel := pkgPath
	if strings.HasPrefix(pkgPath, "github.com/") {
		// Remove module path prefix if present
		parts := strings.SplitN(pkgPath, "/", 4)
		if len(parts) == 4 {
			rel = parts[3]
		}
	}
	return filepath.Join(modRoot, rel), nil
}
