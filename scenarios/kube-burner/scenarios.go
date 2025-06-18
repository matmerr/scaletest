package scenarios

import (
	"fmt"
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
	available := make([]string, 0, len(scenarioRegistry))
	for k := range scenarioRegistry {
		available = append(available, k)
	}

	scenarioName := os.Getenv(KBScenarioEnv)
	if scenarioName == "" {
		return nil, fmt.Errorf("%s not set, available options: %v", KBScenarioEnv, available)
	}
	if scenario, ok := scenarioRegistry[scenarioName]; ok {
		return &scenario, nil
	}

	return nil, fmt.Errorf("scenario not found, available options: %v", available)
}
