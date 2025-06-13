package scenarios

import (
	"context"
	"fmt"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/yaml"
)

func GenerateAllScenarioYAML() error {
	steps := make([]flow.Steper, 0, len(scenarioRegistry))
	for _, scenario := range scenarioRegistry {
		for _, template := range scenario.GetTemplates() {
			steps = append(steps, yaml.GenerateYaml(template))
		}
	}
	root := new(flow.Workflow)
	root.Add(flow.Pipe(steps...))

	err := root.Do(context.Background())
	if err != nil {
		return fmt.Errorf("failed to generate clusterloader2 scenario YAML files: %w", err)
	}
	return nil
}
