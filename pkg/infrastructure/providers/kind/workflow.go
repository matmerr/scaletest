package kind

import (
	flow "github.com/Azure/go-workflow"
	kindsteps "github.com/matmerr/scaletest/pkg/infrastructure/providers/kind/steps"
)

func RunDeployKind() *flow.Workflow {
	w := new(flow.Workflow)
	w.Add(
		flow.Pipe(
			&kindsteps.CreateKindCluster{},
			&kindsteps.GetKindClusterKubeConfig{},
		),
	)

	return w
}

func RunInstallKindCLI() *flow.Workflow {
	w := new(flow.Workflow)
	w.Add(
		flow.Pipe(
			&kindsteps.InstallKindStep{},
		),
	)

	return w
}
