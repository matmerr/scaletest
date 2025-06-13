package scenarios

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/matmerr/scaletest/pkg/yaml"
)

type Scenario interface {
	GetName() string
	GetTemplates() []yaml.Template
}

type KubeBurnerScenario struct {
	Name      string
	Templates []yaml.Template
}

func (s KubeBurnerScenario) GetName() string {
	return s.Name
}

func (s KubeBurnerScenario) GetTemplates() []yaml.Template {
	return s.Templates
}

// GetScenarioFromEnv returns the KubeBurnerScenario from the KB_SCENARIO env var, or logs available options if not found.
func GetScenarioFromEnv() (*KubeBurnerScenario, error) {
	scenarioName := os.Getenv(KBScenarioEnv)
	if scenarioName == "" {
		available := make([]string, 0, len(scenarioRegistry))
		for k := range scenarioRegistry {
			available = append(available, k)
		}
		slog.Error(KBScenarioEnv+" not set. Please set the environment variable to one of the available scenarios.", "available", available)
		return nil, fmt.Errorf("%s not set", KBScenarioEnv)
	}
	if scenario, ok := scenarioRegistry[scenarioName]; ok {
		return &scenario, nil
	}
	available := make([]string, 0, len(scenarioRegistry))
	for k := range scenarioRegistry {
		available = append(available, k)
	}
	slog.Error("Scenario not found", "requested", scenarioName, "available", available)
	return nil, fmt.Errorf("scenario not found")
}
