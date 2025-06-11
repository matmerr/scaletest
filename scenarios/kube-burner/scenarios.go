package scenarios

import (
	"context"
	"fmt"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/yaml"
	netpolchurn "github.com/matmerr/scaletest/scenarios/kube-burner/netpol-churn"
)

type Scenario string

const (
	NetpolChurnConfig Scenario = "netpolchurn"
	APIIntensive      Scenario = "apiintensive"
)

var providerSetupIndex = map[Scenario][]yaml.Template{
	NetpolChurnConfig: {
		netpolchurn.NewNetpolChurnConfig(),
	},

	APIIntensive: {
		// Placeholder for future API intensive scenario
		// apiintensive.NewApiIntensiveConfig(),
	},
}

func GetScenarioSteps(kbs Scenario) ([]yaml.Template, error) {
	steps, ok := providerSetupIndex[kbs]
	if !ok {
		return nil, fmt.Errorf("unknown scenario: %s", kbs)
	}
	return steps, nil
}

// GenerateAllScenarioYAML generates YAML files for all defined scenarios in the provider setup index.
// each scenario will have a config_generated.yaml file created in the current working directory.
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
