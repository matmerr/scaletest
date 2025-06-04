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

	steps := make([]flow.Steper, 0, len(scenarios))
	for _, scenario := range scenarios {
		steps = append(steps, kb.RunKubeBurner(scenario))
	}

	setup := new(flow.Workflow).Add(
		flow.Pipe(
			// Install prerequisites
			kb.InstallKubeBurner(),
			kind.RunInstallKind(),

			//Run tests,
			new(welcome.Intro),
		),
	)

	scenarioSteps := new(flow.Workflow).Add(
		flow.Pipe(
			steps...,
		),
	)

	// add all scenarios to the root workflow
	root := new(flow.Workflow).Add(flow.Pipe(
		setup,
		scenarioSteps,
	))

	if err := root.Do(context.Background()); err != nil {
		t.Fatalf("failed to run workflow: %v", err)
	}
}
