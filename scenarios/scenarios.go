package scenarios

import (
	"github.com/matmerr/scaletest/pkg/yaml"
	netpolchurn "github.com/matmerr/scaletest/scenarios/kube-burner/netpol-churn"
)

var Index = []yaml.Template{
	netpolchurn.NewNetpolChurnConfig(),
	//apiintensive.NewApiIntensiveConfig(),
}
