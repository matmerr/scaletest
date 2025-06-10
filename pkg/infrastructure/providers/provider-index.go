package providers

import (
	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/infrastructure/addons/cilium"
	"github.com/matmerr/scaletest/pkg/infrastructure/providers/azure"
	"github.com/matmerr/scaletest/pkg/infrastructure/providers/kind"
)

type Provider string

const (
	KindWithCilium     Provider = "kindwithcilium"
	AKSExistingCluster Provider = "aksexistingcluster"
)

func GetProviderSetupSteps(prov Provider) flow.AddSteps {
	return providerSetupIndex[prov]
}

var providerSetupIndex = map[Provider]flow.AddSteps{

	// creates a local kind cluster, and install cilium
	KindWithCilium: flow.Pipe(
		kind.RunDeployKind(),
		cilium.RunInstallCilium(),
	),

	// to be tested and implemented
	AKSExistingCluster: flow.Pipe(
		azure.GetExistingAzureCluster(),
	),
}
