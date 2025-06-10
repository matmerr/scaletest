package azure

import (
	flow "github.com/Azure/go-workflow"
	azuresteps "github.com/matmerr/scaletest/pkg/infrastructure/providers/azure/steps"
)

func GetExistingAzureCluster() *flow.Workflow {
	w := new(flow.Workflow)
	w.Add(
		flow.Pipe(
			&azuresteps.GetKubeConfig{},
		),
	)

	return w
}
