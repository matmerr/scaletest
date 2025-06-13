package providers

import (
	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/infrastructure/addons/cilium"
	"github.com/matmerr/scaletest/pkg/infrastructure/providers/azure"
	"github.com/matmerr/scaletest/pkg/infrastructure/providers/kind"
)

var providerSetupIndex = map[string]*ClusterProvider{
	"kindwithcilium": {
		name: "Kind with Cilium",
		steps: flow.Pipe(
			kind.RunDeployKind(),
			cilium.RunInstallCilium()),
	},

	"aksexistingcluster": {
		name: "AKS Existing Cluster",
		steps: flow.Pipe(
			azure.GetExistingAzureCluster(),
		),
	},
}
