package main

import (
	"context"
	"testing"

	flow "github.com/Azure/go-workflow"
	netpolchurn "github.com/matmerr/scaletest/scenarios/kube-burner/netpol-churn"
	kind "github.com/matmerr/scaletest/workflows/infra/kind"
	kb "github.com/matmerr/scaletest/workflows/kube-burner"
	prom "github.com/matmerr/scaletest/workflows/prometheus"
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
			kind.RunDeployKind(),
			prom.RunConfigurePrometheus(),
			kb.RunKubeBurner(netpolchurn.NewNetpolChurnConfig()),
		),
	)

	if err := root.Do(context.Background()); err != nil {
		t.Fatalf("failed to run workflow: %v", err)
	}
}
