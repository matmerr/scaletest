package steps

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/matmerr/scaletest/pkg/yaml"
)

type RunKubeBurner struct {
	Namespace  string
	ConfigPath string // Path to the kube-burner config file
	Template   yaml.YamlGenerator
}

func (c *RunKubeBurner) Do(ctx context.Context) error {
	// Assume scenarios.Config contains the config file path and other args

	// Example: run "kubeburner init --namespace <namespace> --config <configFile>"
	cmdArgs := []string{
		"init",
		"--namespace", c.Namespace,
		"--config", c.ConfigPath, // adjust field name as needed
	}

	fmt.Printf("Running command: kube-burner %v\n", cmdArgs)
	cmd := exec.CommandContext(ctx, "kube-burner", cmdArgs...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run kube-burner: %w", err)
	}

	return nil
}
