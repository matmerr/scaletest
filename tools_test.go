package main

import (
	"context"
	"testing"

	flow "github.com/Azure/go-workflow"
	cl2 "github.com/matmerr/scaletest/pkg/executors/clusterloader2"
	kb "github.com/matmerr/scaletest/pkg/executors/kube-burner"
	kind "github.com/matmerr/scaletest/pkg/infrastructure/providers/kind"
	cl2scenarios "github.com/matmerr/scaletest/scenarios/clusterloader2"
	kbscenarios "github.com/matmerr/scaletest/scenarios/kube-burner"
)

func TestGenerate(t *testing.T) {
	kbscenarios.GenerateAllScenarioYAML()
	cl2scenarios.GenerateAllScenarioYAML()
}

func TestDownloadTools(t *testing.T) {
	root := new(flow.Workflow).Add(
		flow.Pipe(
			kb.RunInstallKubeBurnerCLI(),
			kind.RunInstallKindCLI(),
			cl2.RunInstallClusterLoader2CLI(),
		),
	)

	err := root.Do(context.Background())
	if err != nil {
		t.Fatalf("failed download tools: %v", err)
	}
}
