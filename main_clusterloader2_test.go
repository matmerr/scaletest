package main

import (
	"testing"

	flow "github.com/Azure/go-workflow"
	cl2scenarios "github.com/matmerr/scaletest/scenarios/clusterloader2"
	"k8s.io/client-go/tools/clientcmd"

	cl2 "github.com/matmerr/scaletest/pkg/executors/clusterloader2"
	"github.com/matmerr/scaletest/pkg/infrastructure/providers"

	"github.com/matmerr/scaletest/pkg/welcome"
)

func TestRunClusterLoader2Scenarios(t *testing.T) {

	// this is a Clusterloader2 test, so we need to set up the ClusterLoader2 executor
	cl2exec := cl2.NewClusterLoader2Executor(

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

	// specify the cluster provider we want to, since we specified it above
	kindCluster := providers.KindWithCilium

	// we want to run the UniformQPS scenario from ClusterLoader2
	scenario := cl2scenarios.UniformQPS

	// kick off the run scenarios
	err := RunScenarios(kindCluster, cl2exec, cl2scenarios.GetScenarioSteps(scenario))
	if err != nil {
		t.Fatalf("failed to run ClusterLoader2 scenarios: %v", err)
	}
}
