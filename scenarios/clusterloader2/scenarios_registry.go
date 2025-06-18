package scenarios

import (
	"github.com/matmerr/scaletest/pkg/yaml"
	uniformqps "github.com/matmerr/scaletest/scenarios/clusterloader2/configs/uniformqps"
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
}
