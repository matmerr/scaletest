package utils

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func GetPackagePath(data interface{}) (string, error) {
	t := reflect.TypeOf(data)
	fmt.Println("Type Name:", t.Name())
	fmt.Println("Package Path:", t.PkgPath()) // empty if defined in main

	pkgPath := t.PkgPath()
	runningPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Extract the Go module path from the base path
	const goSrcPrefix = "/src/"
	idx := strings.Index(runningPath, goSrcPrefix)
	if idx == -1 {
		panic("invalid base path, must contain /src/")
	}
	modulePath := runningPath[idx+len(goSrcPrefix):]

	// Strip the module path from the full path
	relPath := strings.TrimPrefix(pkgPath, modulePath+"/")
	return relPath, nil
}
