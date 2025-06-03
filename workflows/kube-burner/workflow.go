package scenarios

import (
	"context"
	"fmt"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/yaml"
	kb "github.com/matmerr/scaletest/workflows/kube-burner/steps"
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

func RunKubeBurner(yamlDirectory yaml.YamlGenerator) *flow.Workflow {
	createYaml := &CreateYaml{
		InputConfig: yamlDirectory,
	}

	runKubeBurner := &kb.KubeBurner{
		Namespace: "kube-burner",
	}

	w := new(flow.Workflow)
	w.Add(
		flow.Step(createYaml),
		flow.Step(runKubeBurner).DependsOn(createYaml).Input(func(ctx context.Context, g *kb.KubeBurner) error {
			runKubeBurner.ConfigPath = createYaml.OutputConfig
			return nil
		}),
	)

	return w
}
