package scenarios

import (
	"context"

	flow "github.com/Azure/go-workflow"
	kbsteps "github.com/matmerr/scaletest/pkg/executors/kube-burner/steps"
	"github.com/matmerr/scaletest/pkg/yaml"
)

type KubeBurnerExecutor struct {
	SetupSteps flow.AddSteps
}

func NewKubeBurnerExecutor(setupSteps flow.AddSteps) *KubeBurnerExecutor {
	return &KubeBurnerExecutor{
		SetupSteps: setupSteps,
	}
}

func (k *KubeBurnerExecutor) GetRunWorkflow(templateConfig yaml.Template) *flow.Workflow {
	createYaml := &yaml.CreateYaml{
		InputConfig: templateConfig,
	}

	runKubeBurner := &kbsteps.KubeBurner{
		Namespace: "kube-burner",
	}

	w := new(flow.Workflow).Add(
		flow.Step(createYaml),
		flow.Step(runKubeBurner).DependsOn(createYaml).Input(func(ctx context.Context, g *kbsteps.KubeBurner) error {
			runKubeBurner.ConfigPath = createYaml.OutputConfig
			return nil
		}),
	)

	return w
}

func (k *KubeBurnerExecutor) GetSetupWorkflow() flow.AddSteps {
	return k.SetupSteps
}

func RunInstallKubeBurnerCLI() *flow.Workflow {
	w := new(flow.Workflow).Add(
		flow.Pipe(
			&kbsteps.InstallKubeBurnerCLI{},
		),
	)

	return w
}
