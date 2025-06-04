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
	root := new(flow.Workflow).Add(
		flow.Pipe(
			// Install prerequisites
			kb.InstallKubeBurner(),
			kind.RunInstallKind(),

			//Run tests,
			new(welcome.Intro),
		),
	)

	steps := make([]flow.Steper, 0, len(scenarios))
	for _, scenario := range scenarios {
		steps = append(steps, kb.GenerateYaml(scenario))
	}

	flow.Pipe(steps...)

	if err := root.Do(context.Background()); err != nil {
		t.Fatalf("failed to run workflow: %v", err)
	}
}
