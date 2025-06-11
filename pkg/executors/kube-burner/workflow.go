package scenarios

import (
	"context"

	flow "github.com/Azure/go-workflow"
	kbsteps "github.com/matmerr/scaletest/pkg/executors/kube-burner/steps"
	"github.com/matmerr/scaletest/pkg/yaml"
	kbscenarios "github.com/matmerr/scaletest/scenarios/kube-burner"
)

type KubeBurnerExecutor struct {
	scenario kbscenarios.Scenario

	SetupSteps flow.AddSteps
}

func NewKubeBurnerExecutor(scenario kbscenarios.Scenario, executorSetupSteps flow.AddSteps) *KubeBurnerExecutor {
	return &KubeBurnerExecutor{
		scenario:   scenario,
		SetupSteps: executorSetupSteps,
	}
}

func (k *KubeBurnerExecutor) GetScenarioTemplates() ([]yaml.Template, error) {
	return kbscenarios.GetScenarioSteps(k.scenario)
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
