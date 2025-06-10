package kindsteps

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type InstallKindStep struct {
	Version string // e.g. v0.29.0
}

func (s *InstallKindStep) Do(ctx context.Context) error {
	version := s.Version
	if version == "" {
		version = "v0.29.0"
	}
	binDir := "./tools/bin"
	kindPath := filepath.Join(binDir, "kind")

	if _, err := os.Stat(kindPath); err == nil {
		slog.Info("kind binary already exists", "path", kindPath)
		return nil
	}

	if err := os.MkdirAll(binDir, 0o755); err != nil {
		slog.Error("Failed to create bin dir", "err", err)
		return err
	}

	url := "https://github.com/kubernetes-sigs/kind/releases/download/" + version + "/kind-linux-amd64"
	tarGzPath := filepath.Join(binDir, "kind-linux-amd64")

	slog.Info("Downloading kind", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Failed to download kind", "err", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		slog.Error("Failed to download kind: bad status", "status", resp.Status)
		return errors.New("failed to download kind: " + resp.Status)
	}

	out, err := os.Create(tarGzPath)
	if err != nil {
		slog.Error("Failed to create file for kind", "err", err)
		return err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		os.Remove(tarGzPath)
		slog.Error("Failed to write kind binary", "err", err)
		return err
	}
	out.Close()

	// Move to final path and make executable
	if err := os.Rename(tarGzPath, kindPath); err != nil {
		slog.Error("Failed to move kind binary", "err", err)
		return err
	}
	if err := os.Chmod(kindPath, 0o755); err != nil {
		slog.Error("Failed to chmod kind binary", "err", err)
		return err
	}

	slog.Info("Installed kind binary", "path", kindPath)
	return nil
}
