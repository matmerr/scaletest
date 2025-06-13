package scenarios

import (
	"fmt"
	"log/slog"
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
	scenarioName := os.Getenv(CL2ScenarioEnv)
	if scenarioName == "" {
		available := make([]string, 0, len(scenarioRegistry))
		for k := range scenarioRegistry {
			available = append(available, k)
		}
		slog.Error("CL2_SCENARIO not set. Please set the environment variable to one of the available scenarios.", "available", available)
		return nil, fmt.Errorf("CL2_SCENARIO not set")
	}
	if scenario, ok := scenarioRegistry[scenarioName]; ok {
		return &scenario, nil
	}
	available := make([]string, 0, len(scenarioRegistry))
	for k := range scenarioRegistry {
		available = append(available, k)
	}
	slog.Error("Scenario not found", "requested", scenarioName, "available", available)
	return nil, fmt.Errorf("Scenario not found")
}
