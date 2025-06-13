package main

import (
	"testing"

	flow "github.com/Azure/go-workflow"
	"k8s.io/client-go/tools/clientcmd"

	cl2 "github.com/matmerr/scaletest/pkg/executors/clusterloader2"
	"github.com/matmerr/scaletest/pkg/infrastructure/providers"

	"github.com/matmerr/scaletest/pkg/welcome"
	cl2scenarios "github.com/matmerr/scaletest/scenarios/clusterloader2"
)

func TestRunClusterLoader2Scenarios(t *testing.T) {

	clusterProvider, err := providers.GetClusterProviderFromEnv()
	if err != nil {
		t.Fatalf("failed to get cluster provider from environment: %v", err)
	}

	scenario, err := cl2scenarios.GetScenarioFromEnv()
	if err != nil {
		t.Fatalf("failed to get scenario from environment: %v", err)
	}

	// this is a Clusterloader2 test, so we need to set up the ClusterLoader2 executor
	cl2exec := cl2.NewClusterLoader2Executor(

		// pass the scenario to the executor here
		scenario,

		// here we can specify any dependencies to install, and/or addons we want to install
		flow.Pipe(

			// install clusterloader2 CLI since this is a clusterloader2 scenario
			cl2.RunInstallClusterLoader2CLI(),

			// print tool versions
			new(welcome.Intro),
		),
		cl2.ExecutorOptions{
			Kubeconfig: clientcmd.RecommendedHomeFile,
			Provider:   "kind",
		},
	)

	// kick off the run scenarios
	err = RunScenarios(clusterProvider, cl2exec)
	if err != nil {
		t.Fatalf("failed to run ClusterLoader2 scenarios: %v", err)
	}
}
