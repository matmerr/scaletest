package scenarios

import (
	"github.com/matmerr/scaletest/pkg/yaml"
	netpolchurn "github.com/matmerr/scaletest/scenarios/kube-burner/netpol-churn"
)

const (
	KBScenarioEnv = "KB_SCENARIO"
)

var scenarioRegistry = map[string]KubeBurnerScenario{
	"netpolchurn": {
		Name:      "netpolchurn",
		Templates: []yaml.Template{netpolchurn.NewNetpolChurnConfig()},
	},
	// Add more scenarios here as needed
}
