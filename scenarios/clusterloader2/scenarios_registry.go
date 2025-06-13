package scenarios

import (
	"github.com/matmerr/scaletest/pkg/yaml"
	"github.com/matmerr/scaletest/scenarios/clusterloader2/networkload"
	uniformqps "github.com/matmerr/scaletest/scenarios/clusterloader2/uniformqps"
)

// preserve mapping of string to scenarios, which may result in 1:many scenarios by the same name
// later on
var scenarioRegistry = map[string]ClusterLoader2Scenario{
	"UniformQPS": {
		Name: "UniformQPS",
		Templates: []yaml.Template{
			uniformqps.NewUniformQPSConfig(),
		},
	},
	"HighTrafficLoad": {
		Name: "HighTrafficLoad",
		Templates: []yaml.Template{
			networkload.NewNetworkLoadConfig(),
		},
	},
}
