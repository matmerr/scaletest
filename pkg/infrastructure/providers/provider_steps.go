package providers

import (
	"fmt"
	"log/slog"
	"os"

	flow "github.com/Azure/go-workflow"
)

const (
	EnvClusterProvider string = "CLUSTER_PROVIDER"
)

type Provider interface {
	Name() string
	GetSteps() flow.AddSteps
}

type ClusterProvider struct {
	name  string
	steps flow.AddSteps
}

func (p *ClusterProvider) Name() string {
	return p.name
}

func (p *ClusterProvider) GetSteps() flow.AddSteps {
	return p.steps
}

func GetClusterProviderFromEnv() (*ClusterProvider, error) {
	clusterProviderName := os.Getenv(EnvClusterProvider)

	if provider, exists := providerSetupIndex[clusterProviderName]; exists {
		return provider, nil
	}
	available := make([]string, 0, len(providerSetupIndex))
	for k := range providerSetupIndex {
		available = append(available, k)
	}
	slog.Warn("ClusterProvider not found", "requested", clusterProviderName, "available", available)
	return nil, fmt.Errorf("ClusterProvider not found")

}
