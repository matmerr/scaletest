package scenarios

import (
	"github.com/matmerr/scaletest/pkg/yaml"
	"github.com/matmerr/scaletest/scenarios/clusterloader2/networkload"
	uniformqps "github.com/matmerr/scaletest/scenarios/clusterloader2/uniformqps"
)

// preserve mapping of string to scenarios, which may result in 1:many scenarios by the same name
// later on
var scenarioRegistry = map[string]ClusterLoader2Scenario{
	"uniformqps": {
		Name: "uniformqps",
		Templates: []yaml.Template{
			uniformqps.NewUniformQPSConfig(),
		},
	},
	"hightrafficload": {
		Name: "hightrafficload",
		Templates: []yaml.Template{
			networkload.NewNetworkLoadConfig(),
		},
	},
}
