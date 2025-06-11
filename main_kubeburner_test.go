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

	// specify the cluster provider environment we want to use, in this case a kind cluster with Cilium
	kindCluster := providers.KindWithCilium

	// get the scenario for netpol churn config, in this case the netpol scenario
	scenario := kbscenario.NetpolChurnConfig

	// create a new Kube-Burner executor, which will run the scenarios
	kbexec := kb.NewKubeBurnerExecutor(

		// pass the scenario to the executor
		scenario,

		// here we can specify any dependencies to install, and/or addons we want to install
		flow.Pipe(

			// need the kube-burner CLI to run the scenarios
			kb.RunInstallKubeBurnerCLI(),

			// unlike ClusterLoader2, Kube-Burner does not install Prometheus by default,
			// so we need to install it here
			prom.RunDeployPrometheus(),

			// print the versions
			new(welcome.Intro),
		),
	)

	// Run the scenarios with the specified cluster provider, the executor, and
	err := RunScenarios(kindCluster, kbexec)
	if err != nil {
		t.Fatalf("failed to run Kube-Burner scenarios: %v", err)
	}
}
