package steps

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

type InstallClusterLoader2CLI struct {
	Version string
}

func (s *InstallClusterLoader2CLI) Do(ctx context.Context) error {
	binDir := filepath.Join(".", "tools", "bin")
	binPath := filepath.Join("tools", "bin", "clusterloader2")
	if _, err := os.Stat(binPath); err == nil {
		slog.Info("clusterloader2 already installed", "path", binPath)
		return nil
	}
	toolsSrcDir := filepath.Join(".", "tools", "src")
	repoDir := filepath.Join(toolsSrcDir, "perf-tests")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		slog.Error("Failed to create bin directory", "err", err)
		return err
	}
	if err := os.MkdirAll(toolsSrcDir, 0o755); err != nil {
		slog.Error("Failed to create tools/src directory", "err", err)
		return err
	}
	if _, err := os.Stat(repoDir); err == nil {
		slog.Info("perf-tests repo already cloned", "dir", repoDir)
	} else {
		slog.Info("Cloning perf-tests repo for clusterloader2 build", "dir", repoDir)
		cloneCmd := []string{"git", "clone", "--depth=1", "https://github.com/kubernetes/perf-tests.git", repoDir}
		if err := runCmd(ctx, cloneCmd, ""); err != nil {
			slog.Error("Failed to clone perf-tests repo", "err", err)
			return err
		}
	}
	// run go mod tidy in the clusterloader2 source directory before building
	slog.Info("Running go mod tidy", "dir", "tools/src/perf-tests/clusterloader2")
	if err := runCmd(ctx, []string{"go", "mod", "tidy"}, "tools/src/perf-tests/clusterloader2"); err != nil {
		slog.Error("Failed to run go mod tidy in clusterloader2 source", "err", err)
		return err
	}
	// build clusterloader2 from the correct absolute path
	slog.Info("Running go build", "cmd", "go build -o clusterloader2 cmd/clusterloader.go", "dir", "tools/src/perf-tests/clusterloader2")
	buildCmd := []string{"go", "build", "-o", "clusterloader2", "cmd/clusterloader.go"}
	if err := runCmd(ctx, buildCmd, "tools/src/perf-tests/clusterloader2"); err != nil {
		slog.Error("Failed to build clusterloader2", "err", err)
		return err
	}
	// move the built binary to tools/bin/clusterloader2
	srcBin := filepath.Join("tools", "src", "perf-tests", "clusterloader2", "clusterloader2")
	dstBin := filepath.Join("tools", "bin", "clusterloader2")
	if err := os.Rename(srcBin, dstBin); err != nil {
		slog.Error("Failed to move clusterloader2 binary to bin", "err", err)
		return err
	}
	if err := os.Chmod(dstBin, 0o755); err != nil {
		slog.Error("Failed to chmod clusterloader2 binary", "err", err)
		return err
	}
	slog.Info("clusterloader2 installed", "path", dstBin)
	return nil
}

// runCmd runs a command with optional working directory
func runCmd(ctx context.Context, cmd []string, dir string) error {
	if len(cmd) == 0 {
		return fmt.Errorf("no command provided")
	}
	// Print the full command and working directory context
	slog.Info("Executing command", "cmd", cmd, "dir", dir)
	c := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	if dir != "" {
		c.Dir = dir
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
