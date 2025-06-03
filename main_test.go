package main

import (
	"context"
	"log/slog"
	"testing"

	flow "github.com/Azure/go-workflow"
	netpolchurn "github.com/matmerr/scaletest/scenarios/kube-burner/netpol-churn"
	kind "github.com/matmerr/scaletest/workflows/infra/kind"
	kb "github.com/matmerr/scaletest/workflows/kube-burner"
	prom "github.com/matmerr/scaletest/workflows/prometheus"
)

type Welcome struct {
}

// All required for a step is `Do(context.Context) error`
func (i *Welcome) Do(ctx context.Context) error {
	slog.Info("starting workflow")
	return nil
}

func TestWorkflow(t *testing.T) {
	root := new(flow.Workflow).Add(
		flow.Pipe(
			&Welcome{},
			//&steps.InstallPrometheusStep{},
			kind.RunDeployKind(),
			prom.RunConfigurePrometheus(),
			kb.RunKubeBurner(netpolchurn.NewNetpolChurnConfig()),
		),
	)

	err := root.Do(context.Background())
	if err != nil {
		slog.Error("failed to run workflow", "err", err)
		t.Fatalf("failed to run workflow: %v", err)
	}
}
