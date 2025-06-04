package main

import (
	"log/slog"
	"path/filepath"
	"reflect"

	"github.com/matmerr/scaletest/scenarios"
)

func main() {
	slog.Info("List of All Available Scenarios")
	// all current scenarios in scenarios.Index
	for _, scenario := range scenarios.Index {
		t := reflect.TypeOf(scenario)
		pkgPath := filepath.Base(t.PkgPath())
		slog.Info("Scenario", slog.String("name", pkgPath), slog.String("path", t.PkgPath()))
	}
}
