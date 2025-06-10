package scenarios

import (
	"github.com/matmerr/scaletest/pkg/yaml"
	uniformqps "github.com/matmerr/scaletest/scenarios/clusterloader2/uniformqps"
	netpolchurn "github.com/matmerr/scaletest/scenarios/kube-burner/netpol-churn"
)

var KubeBurnerIndex = []yaml.Template{

	// kube-burner scenarios
	netpolchurn.NewNetpolChurnConfig(),
	//apiintensive.NewApiIntensiveConfig(),

}

var ClusterLoader2Index = []yaml.Template{
	// clusterloader2 scenarios
	uniformqps.NewUniformQPSConfig(),
}
