package kind

import (
	flow "github.com/Azure/go-workflow"
	kindsteps "github.com/matmerr/scaletest/workflows/infra/kind/steps"
)

func RunDeployKind() *flow.Workflow {
	w := new(flow.Workflow)
	w.Add(
		flow.Pipe(
			&kindsteps.CreateKindCluster{},
			&kindsteps.GetKindClusterKubeConfig{},
			&kindsteps.InstallCiliumStep{},
		),
	)

	return w
}
