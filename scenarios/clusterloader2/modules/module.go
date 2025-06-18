package modules

import (
	"embed"
	"path/filepath"
	"runtime"

	"github.com/matmerr/scaletest/pkg/utils"
)

type modulefs struct {
	embed.FS
}

var m modulefs

// ensure the measurement yamls here don't get renamed, since configs may depend on them
// if they do, then this is the warning to update those configs

//go:embed cilium.yaml
var _ []byte

func FullPath() string {
	path, err := utils.GetPackagePath(m)
	if err != nil {
		return ""
	}
	return path
}

func RelativePath() string {
	// Get the path to this file (module.go)
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	// Get the path to the caller
	_, callerFile, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	// Get the directory of this package and the caller
	thisPkgDir := filepath.Dir(thisFile)
	callerDir := filepath.Dir(callerFile)

	// Compute the relative path from caller to this package
	relPath, err := filepath.Rel(callerDir, thisPkgDir)
	if err != nil {
		return ""
	}
	return relPath
}
