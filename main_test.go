package main

import (
	"context"
	"testing"

	flow "github.com/Azure/go-workflow"
	kind "github.com/matmerr/scaletest/workflows/infra/kind"
	kb "github.com/matmerr/scaletest/workflows/kube-burner"
	"github.com/matmerr/scaletest/workflows/welcome"
)

func TestWorkflow(t *testing.T) {
	// Setup steps (prereqs, cluster, intro)
	setup := flow.Pipe(
		kb.InstallKubeBurner(),
		kind.RunInstallKind(),
		new(welcome.Intro),
	)

	// Scenario steps (run kube-burner for each scenario)
	scenarioSteps := make([]flow.Steper, 0, len(Scenarios))
	for _, scenario := range Scenarios {
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
