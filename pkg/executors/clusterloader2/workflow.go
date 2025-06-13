package cl2

import (
	"context"

	flow "github.com/Azure/go-workflow"
	cl2scenarios "github.com/matmerr/scaletest/scenarios/clusterloader2"

	cl2steps "github.com/matmerr/scaletest/pkg/executors/clusterloader2/steps"
	"github.com/matmerr/scaletest/pkg/yaml"
)

type ClusterLoader2Executor struct {
	scenario cl2scenarios.Scenario

	SetupSteps flow.AddSteps

	Options ExecutorOptions
}

type ExecutorOptions struct {
	Kubeconfig string
	Provider   string
}

func NewClusterLoader2Executor(scenario cl2scenarios.Scenario, setupSteps flow.AddSteps, options ExecutorOptions) *ClusterLoader2Executor {
	return &ClusterLoader2Executor{
		scenario:   scenario,
		SetupSteps: setupSteps,
		Options:    options,
	}
}

func (k *ClusterLoader2Executor) GetScenarioTemplates() []yaml.Template {
	return k.scenario.GetTemplates()
}

func (c *ClusterLoader2Executor) GetRunWorkflow(templateConfig yaml.Template) *flow.Workflow {
	createYaml := &yaml.CreateYaml{
		InputConfig: templateConfig,
	}

	runCL2 := &cl2steps.ClusterLoader2{
		Kubeconfig: c.Options.Kubeconfig,
		Provider:   c.Options.Provider,
	}

	w := new(flow.Workflow).Add(
		flow.Step(createYaml),
		flow.Step(runCL2).DependsOn(createYaml).Input(func(ctx context.Context, g *cl2steps.ClusterLoader2) error {
			runCL2.ConfigPath = createYaml.OutputConfig
			return nil
		}),
	)

	return w
}

func (k *ClusterLoader2Executor) GetSetupWorkflow() flow.AddSteps {
	return k.SetupSteps
}

func RunInstallClusterLoader2CLI() *flow.Workflow {
	w := new(flow.Workflow).Add(
		flow.Pipe(
			&cl2steps.InstallClusterLoader2CLI{},
		),
	)

	return w
}
