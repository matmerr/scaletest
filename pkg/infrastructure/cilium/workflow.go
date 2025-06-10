package cilium

import (
	flow "github.com/Azure/go-workflow"
	ciliumsteps "github.com/matmerr/scaletest/pkg/infrastructure/cilium/steps"
)

func RunInstallCilium() *flow.Workflow {
	w := new(flow.Workflow)
	w.Add(
		flow.Pipe(
			&ciliumsteps.InstallCiliumStep{},
		),
	)

	return w
}
