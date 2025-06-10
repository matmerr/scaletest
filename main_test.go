package main

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"reflect"
	"testing"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/executors"
	"github.com/matmerr/scaletest/pkg/yaml"
	"github.com/matmerr/scaletest/scenarios"
	"k8s.io/client-go/tools/clientcmd"

	cl2 "github.com/matmerr/scaletest/pkg/executors/clusterloader2"
	kb "github.com/matmerr/scaletest/pkg/executors/kube-burner"

	"github.com/matmerr/scaletest/pkg/infrastructure/cilium"
	prom "github.com/matmerr/scaletest/pkg/infrastructure/prometheus"
	"github.com/matmerr/scaletest/pkg/welcome"
)

func TestRunClusterLoader2Scenarios(t *testing.T) {
	cl2exec := cl2.NewClusterLoader2Executor(
		flow.Pipe(
			//cl2.RunInstallClusterLoader2CLI(),
			//kind.RunInstallKindCLI(),
			//kind.RunDeployKind(),
			//cilium.RunInstallCilium(),
			//prom.RunDeployPrometheus(),
			new(welcome.Intro),
		),
		cl2.ExecutorOptions{
			Kubeconfig: clientcmd.RecommendedHomeFile,
			Provider:   "kind",
		},
	)

	err := RunScenarios(cl2exec, scenarios.ClusterLoader2Index)
	if err != nil {
		t.Fatalf("failed to run ClusterLoader2 scenarios: %v", err)
	}
}

func TestRunKubeBurnerScenarios(t *testing.T) {
	kbexec := kb.NewKubeBurnerExecutor(
		flow.Pipe(
			kb.RunInstallKubeBurnerCLI(),
			//kind.RunInstallKind(),
			//kind.RunDeployKind(),
			cilium.RunInstallCilium(),
			prom.RunDeployPrometheus(),
			new(welcome.Intro),
		),
	)

	err := RunScenarios(kbexec, scenarios.KubeBurnerIndex)
	if err != nil {
		t.Fatalf("failed to run Kube-Burner scenarios: %v", err)
	}
}

func RunScenarios(executor executors.Executor, scenarios []yaml.Template) error {
	for _, scenario := range scenarios {
		t := reflect.TypeOf(scenario)
		pkgPath := filepath.Base(t.PkgPath())
		slog.Info("Scenario", slog.String("name", pkgPath), slog.String("path", t.PkgPath()))
	}

	// Scenario steps (run kube-burner for each scenario)
	scenarioSteps := make([]flow.Steper, 0, len(scenarios))
	for _, scenario := range scenarios {
		scenarioSteps = append(scenarioSteps, executor.GetRunWorkflow(scenario))
	}

	scenarioPipe := flow.Pipe(scenarioSteps...)

	// Use batch pipe to stitch setup and scenarios
	root := new(flow.Workflow).Add(
		flow.BatchPipe(
			executor.GetSetupWorkflow(),
			scenarioPipe,
		),
	)

	if err := root.Do(context.Background()); err != nil {
		return fmt.Errorf("failed to run workflow: %w", err)
	}

	return nil
}
