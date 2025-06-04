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
	root := new(flow.Workflow)

	stepKubeBurner := flow.Step(kb.InstallKubeBurner())
	prev := root.Add(stepKubeBurner)

	stepKind := flow.Step(kind.RunInstallKind()).DependsOn(prev)
	prev = root.Add(stepKind)

	stepIntro := flow.Step(new(welcome.Intro)).DependsOn(prev)
	prev = root.Add(stepIntro)

	for _, scenario := range Scenarios {
		scenarioStep := flow.Step(kb.RunKubeBurner(scenario)).DependsOn(prev)
		prev = root.Add(scenarioStep)
	}

	if err := root.Do(context.Background()); err != nil {
		t.Fatalf("failed to run workflow: %v", err)
	}
}
