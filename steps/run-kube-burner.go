package steps

import (
	"context"
	"fmt"
	"os/exec"
)

type RunKubeBurner struct {
	Namespace  string
	ConfigPath string
}

func (c *RunKubeBurner) Do(ctx context.Context) error {
	// Assume scenarios.Config contains the config file path and other args

	// Example: run "kubeburner init --namespace <namespace> --config <configFile>"
	cmdArgs := []string{
		"init",
		"--namespace", c.Namespace,
		"--config", c.ConfigPath, // adjust field name as needed
	}

	fmt.Printf("Running command: kubeburner %v\n", cmdArgs)
	cmd := exec.CommandContext(ctx, "kubeburner", cmdArgs...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run kubeburner: %w", err)
	}

	return nil
}
