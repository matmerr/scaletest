package main

import (
	"testing"

	flow "github.com/Azure/go-workflow"
	kbscenario "github.com/matmerr/scaletest/scenarios/kube-burner"

	kb "github.com/matmerr/scaletest/pkg/executors/kube-burner"

	prom "github.com/matmerr/scaletest/pkg/infrastructure/addons/prometheus"
	"github.com/matmerr/scaletest/pkg/infrastructure/providers"
	"github.com/matmerr/scaletest/pkg/welcome"
)

func TestRunKubeBurnerScenarios(t *testing.T) {
	clusterProvider, err := providers.GetClusterProviderFromEnv()
	if err != nil {
		t.Fatalf("failed to get cluster provider from environment: %v", err)
	}

	scenario, err := kbscenario.GetScenarioFromEnv()
	if err != nil {
		t.Fatalf("failed to get scenario from environment: %v", err)
	}

	// create a new Kube-Burner executor, which will run the scenarios
	kbexec := kb.NewKubeBurnerExecutor(
		scenario,
		flow.Pipe(
			kb.RunInstallKubeBurnerCLI(),
			prom.RunDeployPrometheus(),
			new(welcome.Intro),
		),
	)

	err = RunScenarios(clusterProvider, kbexec)
	if err != nil {
		t.Fatalf("failed to run Kube-Burner scenarios: %v", err)
	}
}
