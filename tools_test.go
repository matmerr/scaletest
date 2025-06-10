package main

import (
	"context"
	"testing"

	flow "github.com/Azure/go-workflow"
	cl2 "github.com/matmerr/scaletest/pkg/executors/clusterloader2"
	kb "github.com/matmerr/scaletest/pkg/executors/kube-burner"
	kind "github.com/matmerr/scaletest/pkg/infrastructure/providers/kind"
	"github.com/matmerr/scaletest/pkg/yaml"
	scenarios "github.com/matmerr/scaletest/scenarios"
)

func TestGenerate(t *testing.T) {
	root := new(flow.Workflow)

	steps := make([]flow.Steper, 0, len(scenarios.KubeBurnerIndex)+len(scenarios.ClusterLoader2Index))
	for _, scenario := range scenarios.KubeBurnerIndex {
		steps = append(steps, yaml.GenerateYaml(scenario))
	}

	for _, scenario := range scenarios.ClusterLoader2Index {
		steps = append(steps, yaml.GenerateYaml(scenario))
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
