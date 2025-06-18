package scenarios

import (
	"fmt"
	"os"

	"github.com/matmerr/scaletest/pkg/yaml"
)

const (
	CL2ScenarioEnv = "CL2_SCENARIO"
)

type Scenario interface {
	GetName() string
	GetTemplates() []yaml.Template
}

type ClusterLoader2Scenario struct {
	Name      string
	Templates []yaml.Template
}

func (s ClusterLoader2Scenario) GetName() string {
	return s.Name
}

func (s ClusterLoader2Scenario) GetTemplates() []yaml.Template {
	return s.Templates
}

func GetScenarioSteps(scenarioName string) ([]yaml.Template, error) {
	if scenario, ok := scenarioRegistry[scenarioName]; ok {
		return scenario.GetTemplates(), nil
	}
	return nil, fmt.Errorf("Scenario not found")
}

// GetScenarioFromEnv returns the ClusterLoader2Scenario from the CL2_SCENARIO env var, or logs available options if not found.
func GetScenarioFromEnv() (*ClusterLoader2Scenario, error) {
	available := make([]string, 0, len(scenarioRegistry))
	for k := range scenarioRegistry {
		available = append(available, k)
	}

	scenarioName := os.Getenv(CL2ScenarioEnv)
	if scenarioName == "" {
		return nil, fmt.Errorf("%s not set, available options: %v", CL2ScenarioEnv, available)
	}
	if scenario, ok := scenarioRegistry[scenarioName]; ok {
		return &scenario, nil
	}

	return nil, fmt.Errorf("scenario not found, available options: %v", available)
}
