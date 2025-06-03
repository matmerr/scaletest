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
	cmdArgs := []string{
		"init",
		"--config", filepath.Base(c.ConfigPath),
	}

	slog.Info("Running command", "cmd", "kube-burner", "args", cmdArgs)
	cmd := exec.CommandContext(ctx, "kube-burner", cmdArgs...)
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
