package main

import (
	"context"
	"testing"

	flow "github.com/Azure/go-workflow"
	kind "github.com/matmerr/scaletest/workflows/infra/kind"
	kb "github.com/matmerr/scaletest/workflows/kube-burner"
)

func TestGenerate(t *testing.T) {
	root := new(flow.Workflow)

	steps := make([]flow.Steper, 0, len(scenarios))
	for _, scenario := range scenarios {
		steps = append(steps, kb.GenerateYaml(scenario))
	}

	root.Add(flow.Pipe(steps...))

	err := root.Do(context.Background())
	if err != nil {
		t.Fatalf("failed to run generate configs: %v", err)
	}
}

func TestDownloadTools(t *testing.T) {
	root := new(flow.Workflow).Add(
		flow.Pipe(
			kb.InstallKubeBurner(),
			kind.RunInstallKind(),
		),
	)
	err := root.Do(context.Background())
	if err != nil {
		t.Fatalf("failed download tools: %v", err)
	}
}
