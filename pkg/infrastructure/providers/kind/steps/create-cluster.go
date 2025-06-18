package kindsteps

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
	kindcluster "sigs.k8s.io/kind/pkg/cluster"
)

// CreateKindCluster is a workflow step to deploy a kind cluster
// Usage: step := &CreateKindCluster{Name: "my-cluster"}
//
//	err := step.Do(ctx)
type CreateKindCluster struct {
	Name string // Name of the kind cluster
}

func (s *CreateKindCluster) Do(ctx context.Context) error {
	name := s.Name
	if name == "" {
		name = "kind"
	}
	slog.Info("Creating kind cluster", "name", name)

	provider := kindcluster.NewProvider()
	clusters, err := provider.List()
	if err == nil {
		for _, c := range clusters {
			if c == name {
				slog.Info("Kind cluster already exists, using existing cluster", "name", name)
				return nil
			}
		}
	}

	// Cilium-ready config: disables default CNI, sets kubeProxyMode: none, and uses extraMounts for containerd
	kindConfig := map[string]interface{}{
		"kind":       "Cluster",
		"apiVersion": "kind.x-k8s.io/v1alpha4",
		"networking": map[string]interface{}{
			"disableDefaultCNI": true,
		},
		"nodes": []map[string]interface{}{
			{
				"role": "control-plane",
			},
			{
				"role": "control-plane",
			},
			{"role": "worker"},
			{"role": "worker"},
			{"role": "worker"},
			{"role": "worker"},
			{"role": "worker"},
			{"role": "worker"},
			{"role": "worker"},
		},
	}

	configBytes, err := yaml.Marshal(kindConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal kind config: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "kind-config-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temp file for kind config: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write(configBytes); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write kind config: %w", err)
	}
	tmpFile.Close()

	if err := provider.Create(name, kindcluster.CreateWithConfigFile(tmpFile.Name())); err != nil {
		if err.Error() == fmt.Sprintf("node(s) already exist for a cluster with the name \"%s\"", name) ||
			// fallback for default cluster name 'kind'
			(err.Error() == "node(s) already exist for a cluster with the name \"kind\"") {
			slog.Info("Kind cluster already exists (detected by error), using existing cluster", "name", name)
			return nil
		}
		return fmt.Errorf("failed to create kind cluster: %w", err)
	}
	slog.Info("Kind cluster created", "name", name)
	return nil
}
