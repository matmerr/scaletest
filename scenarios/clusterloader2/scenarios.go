package scenarios

import (
	"context"
	"fmt"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/yaml"
	uniformqps "github.com/matmerr/scaletest/scenarios/clusterloader2/uniformqps"
)

type Scenario string

const (
	UniformQPS Scenario = "UniformQPS"
)

// preserve mapping of string to scenarios, which may result in 1:many scenarios by the same name
// later on
var providerSetupIndex = map[Scenario][]yaml.Template{
	UniformQPS: {
		uniformqps.NewUniformQPSConfig(),
	},
}

func GetScenarioSteps(cl2s Scenario) ([]yaml.Template, error) {
	steps, ok := providerSetupIndex[cl2s]
	if !ok {
		return nil, fmt.Errorf("unknown scenario: %s", cl2s)
	}
	return steps, nil
}

func GenerateAllScenarioYAML() error {
	steps := make([]flow.Steper, 0, len(providerSetupIndex))
	for _, scenario := range providerSetupIndex {
		for _, template := range scenario {
			steps = append(steps, yaml.GenerateYaml(template))
		}
	}
	root := new(flow.Workflow)
	root.Add(flow.Pipe(steps...))

	err := root.Do(context.Background())
	if err != nil {
		return fmt.Errorf("failed to generate kube-burner scenario YAML files: %w", err)
	}
	return nil
}
