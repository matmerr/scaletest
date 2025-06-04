package main

import (
	"context"
	"log/slog"
	"path/filepath"
	"reflect"
	"testing"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/scenarios"
	kind "github.com/matmerr/scaletest/workflows/infra/kind"
	kb "github.com/matmerr/scaletest/workflows/kube-burner"
	prom "github.com/matmerr/scaletest/workflows/prometheus"
	"github.com/matmerr/scaletest/workflows/welcome"
)

func TestWorkflow(t *testing.T) {
	// Setup steps (prereqs, cluster, intro)
	setup := flow.Pipe(
		kb.InstallKubeBurner(),
		kind.RunInstallKind(),
		prom.RunConfigurePrometheus(),
		new(welcome.Intro),
	)

	// print all scenario names to console
	for _, scenario := range scenarios.Index {
		t := reflect.TypeOf(scenario)
		pkgPath := filepath.Base(t.PkgPath())
		slog.Info("Scenario", slog.String("name", pkgPath), slog.String("path", t.PkgPath()))
	}

	// Scenario steps (run kube-burner for each scenario)
	scenarioSteps := make([]flow.Steper, 0, len(scenarios.Index))
	for _, scenario := range scenarios.Index {
		scenarioSteps = append(scenarioSteps, kb.RunKubeBurner(scenario))
	}

	scenarioPipe := flow.Pipe(scenarioSteps...)

	// Use batch pipe to stitch setup and scenarios
	root := new(flow.Workflow).Add(
		flow.BatchPipe(
			setup,
			scenarioPipe,
		),
	)

	if err := root.Do(context.Background()); err != nil {
		t.Fatalf("failed to run workflow: %v", err)
	}
}
