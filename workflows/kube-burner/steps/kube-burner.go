package steps

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

type KubeBurner struct {
	Namespace  string
	ConfigPath string // Path to the kube-burner config file
}

func (c *KubeBurner) Do(ctx context.Context) error {
	// Always use absolute path for the kube-burner binary
	binPath, err := filepath.Abs(filepath.Join("tools", "bin", "kube-burner"))
	if err != nil {
		slog.Error("Failed to get absolute path for kube-burner binary", "err", err)
		return err
	}
	cmdArgs := []string{
		"init",
		"--config", filepath.Base(c.ConfigPath),
	}

	slog.Info("Running command", "cmd", binPath, "args", cmdArgs)
	cmd := exec.CommandContext(ctx, binPath, cmdArgs...)
	cmd.Dir = filepath.Dir(c.ConfigPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	slog.Info("Executing command string", "cmd", cmd.String())

	if err := cmd.Run(); err != nil {
		slog.Error("failed to run kube-burner", "err", err)
		return err
	}

	return nil
}
