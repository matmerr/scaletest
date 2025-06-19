package scenarios

import (
	"github.com/matmerr/scaletest/pkg/yaml"
	uniformqps "github.com/matmerr/scaletest/scenarios/clusterloader2/configs/uniformqps"
	largepodcount "github.com/matmerr/scaletest/scenarios/clusterloader2/configs/largepodcount"
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
	"largepodcount": {
		Name: "largepodcount",
		Templates: []yaml.Template{
			largepodcount.NewLargePodCountConfig(),
		},
	},
}
