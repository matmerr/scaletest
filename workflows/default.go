package scenarios

import (
	"context"
	"fmt"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/yaml"
	"github.com/matmerr/scaletest/steps"
)

type CreateYaml struct {
	InputConfig  yaml.YamlGenerator
	OutputConfig string
}

func (c *CreateYaml) Do(ctx context.Context) error {
	configpath, err := yaml.CreateYamlFile(c.InputConfig)
	if err != nil {
		return fmt.Errorf("failed to create YAML file: %w", err)
	}
	c.OutputConfig = configpath
	return nil
}

func DefaultRun(yamlDirectory yaml.YamlGenerator) *flow.Workflow {
	createYaml := &CreateYaml{
		InputConfig: yamlDirectory,
	}

	runKubeBurner := &steps.RunKubeBurner{
		Namespace: "kube-burner",
	}

	w := new(flow.Workflow)
	w.Add(
		flow.Step(createYaml),
		flow.Step(runKubeBurner).DependsOn(createYaml).Input(func(ctx context.Context, g *steps.RunKubeBurner) error {
			runKubeBurner.ConfigPath = createYaml.OutputConfig
			return nil
		}),
	)

	return w
}
