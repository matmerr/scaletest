package steps

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type RunKubeBurner struct {
	Namespace  string
	ConfigPath string // Path to the kube-burner config file
}

func (c *RunKubeBurner) Do(ctx context.Context) error {
	// Assume scenarios.Config contains the config file path and other args

	// Example: run "kubeburner init --namespace <namespace> --config <configFile>"
	cmdArgs := []string{
		"init",
		"--config", filepath.Base(c.ConfigPath), // adjust field name as needed
	}

	// Ensure the config file is used correctly
	fmt.Printf("Running command: kube-burner %v\n", cmdArgs)
	cmd := exec.CommandContext(ctx, "kube-burner", cmdArgs...)
	cmd.Dir = filepath.Dir(c.ConfigPath) // Set the working directory to where the config file is located
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// print the whole command for debugging
	fmt.Printf("Executing command: %s\n", cmd.String())

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run kube-burner: %w", err)
	}

	return nil
}
