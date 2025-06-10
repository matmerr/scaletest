package steps

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/matmerr/scaletest/pkg/infrastructure/prometheus/podmonitors"
)

type ClusterLoader2 struct {
	ConfigPath string // Path to the clusterloader2 config file

	Kubeconfig string
	Provider   string
}

func (c *ClusterLoader2) Do(ctx context.Context) error {
	binPath, err := filepath.Abs(filepath.Join("tools", "bin", "clusterloader2"))
	if err != nil {
		slog.Error("Failed to get absolute path for clusterloader2 binary", "err", err)
		return err
	}

	// Determine output directory based on scenario config path
	outputDir := filepath.Join(filepath.Dir(c.ConfigPath), "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		slog.Error("Failed to create output directory", "dir", outputDir, "err", err)
		return err
	}

	// Only add --enable-prometheus-server if you want clusterloader2 to deploy its own Prometheus
	// Remove --prometheus-url, as this flag is not supported by your clusterloader2 build
	cmdArgs := []string{
		"--kubeconfig", c.Kubeconfig,
		"--testconfig", c.ConfigPath,
		"--provider", c.Provider,
		"--report-dir", outputDir,
		"--enable-prometheus-server",
		"--prometheus-storage-class-provisioner", "standard",
		"--prometheus-pvc-storage-class", "standard",
		"--prometheus-storage-class-volume-type", "standard",
		"--prometheus-additional-monitors-path", podmonitors.PodMonitorDirectory,
	}

	slog.Info("Running command", "cmd", binPath, "args", cmdArgs)
	cmd := exec.CommandContext(ctx, binPath, cmdArgs...)
	cmd.Dir = filepath.Dir(c.ConfigPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	slog.Info("Executing command string", "cmd", cmd.String())

	if err := cmd.Run(); err != nil {
		slog.Error("failed to run clusterloader2", "err", err)
		return err
	}
	return nil
}
