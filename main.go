package main

import (
	"github.com/matmerr/scaletest/pkg/yaml"
	apiintensive "github.com/matmerr/scaletest/scenarios/kube-burner/api-intensive"
	netpolchurn "github.com/matmerr/scaletest/scenarios/kube-burner/netpol-churn"
)

// Add all kubeburner scenarios here
var Scenarios = []yaml.Template{
	netpolchurn.NewNetpolChurnConfig(),
	apiintensive.NewApiIntensiveConfig(),
}
