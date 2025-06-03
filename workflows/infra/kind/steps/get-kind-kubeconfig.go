package kindsteps

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	kindcluster "sigs.k8s.io/kind/pkg/cluster"
)

// Step to get the kubeconfig for a kind cluster
// Usage: step := &GetKindClusterKubeConfig{Name: "my-cluster"}
//
//	err := step.Do(ctx)
type GetKindClusterKubeConfig struct {
	Name           string // Name of the kind cluster
	KubeconfigPath string // Output path for kubeconfig (optional)
}

func (s *GetKindClusterKubeConfig) Do(ctx context.Context) error {
	provider := kindcluster.NewProvider()
	kubeconfig, err := provider.KubeConfig(s.Name, false)
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig for cluster %s: %w", s.Name, err)
	}
	if s.KubeconfigPath != "" {
		if err := os.WriteFile(s.KubeconfigPath, []byte(kubeconfig), 0600); err != nil {
			return fmt.Errorf("failed to write kubeconfig to %s: %w", s.KubeconfigPath, err)
		}
		slog.Info("Kubeconfig written", "path", s.KubeconfigPath)
	} else {
		slog.Info("Kubeconfig output", "kubeconfig", kubeconfig)
	}
	return nil
}
