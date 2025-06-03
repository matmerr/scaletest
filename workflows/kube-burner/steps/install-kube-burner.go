package steps

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type InstallKubeBurner struct {
	Version string
}

func (s *InstallKubeBurner) Do(ctx context.Context) error {
	binDir := filepath.Join(".", "tools", "bin")
	binPath := filepath.Join(binDir, "kube-burner")
	if _, err := os.Stat(binPath); err == nil {
		slog.Info("kube-burner already installed", "path", binPath)
		return nil
	}
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		slog.Error("Failed to create bin directory", "err", err)
		return err
	}
	version := s.Version
	if version == "" {
		version = "1.16.0"
	}
	tarURL := "https://github.com/kube-burner/kube-burner/releases/download/v" + version + "/kube-burner-V" + version + "-linux-x86_64.tar.gz"
	tmpTar := filepath.Join(binDir, "kube-burner.tar.gz")

	slog.Info("Downloading kube-burner tarball", "url", tarURL)
	resp, err := http.Get(tarURL)
	if err != nil {
		slog.Error("Failed to download kube-burner tarball", "err", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("Failed to download kube-burner: bad status", "status", resp.Status)
		return err
	}
	out, err := os.Create(tmpTar)
	if err != nil {
		slog.Error("Failed to create temp tar file for kube-burner", "err", err)
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, resp.Body); err != nil {
		slog.Error("Failed to write kube-burner tarball", "err", err)
		return err
	}
	out.Close()

	// Extract kube-burner binary from tar.gz
	if err := extractKubeBurnerBinary(tmpTar, binPath); err != nil {
		slog.Error("Failed to extract kube-burner binary", "err", err)
		return err
	}
	if err := os.Chmod(binPath, 0o755); err != nil {
		slog.Error("Failed to chmod kube-burner binary", "err", err)
		return err
	}
	slog.Info("kube-burner installed", "path", binPath)

	// Remove the kube-burner tar.gz after extracting the binary
	if err := os.Remove(tmpTar); err != nil && !os.IsNotExist(err) {
		slog.Warn("Failed to remove kube-burner.tar.gz after extraction", "path", tmpTar, "err", err)
	} else {
		slog.Info("Removed kube-burner.tar.gz after extraction", "path", tmpTar)
	}
	return nil
}

// extractKubeBurnerBinary extracts kube-burner from a tar.gz archive to the given path
func extractKubeBurnerBinary(tarPath, binPath string) error {
	f, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer f.Close()
	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("Error reading tar archive", "err", err)
			return err
		}
		if hdr.Typeflag == tar.TypeReg && (hdr.Name == "kube-burner" || filepath.Base(hdr.Name) == "kube-burner") {
			out, err := os.Create(binPath)
			if err != nil {
				slog.Error("Failed to create kube-burner binary file", "err", err)
				return err
			}
			defer out.Close()
			if _, err := io.Copy(out, tr); err != nil {
				slog.Error("Failed to extract kube-burner binary", "err", err)
				return err
			}
			slog.Info("Extracted kube-burner binary", "path", binPath)
			return nil
		}
	}
	slog.Error("kube-burner binary not found in archive")
	return fmt.Errorf("kube-burner binary not found in archive")
}
