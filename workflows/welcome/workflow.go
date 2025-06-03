package welcome

import (
	"context"
	"log/slog"
	"os/exec"
	"strings"
)

type Intro struct {
}

// All required for a step is `Do(context.Context) error`
func (i *Intro) Do(ctx context.Context) error {
	slog.Info("starting workflow")

	// Print kind version
	kindPath := "./tools/bin/kind"
	kindVerOut, err := runAndGetStdout(kindPath, "--version")
	if err != nil {
		slog.Warn("Failed to get kind version", "err", err)
	} else {
		kindVerOut = prettySingleLine(kindVerOut)
		slog.Info("kind version", "version", kindVerOut)
	}

	// Print kube-burner version
	kubeBurnerPath := "./tools/bin/kube-burner"
	kbVerOut, err := runAndGetStdout(kubeBurnerPath, "version")
	if err != nil {
		slog.Warn("Failed to get kube-burner version", "err", err)
	} else {
		kbVerOut = prettySingleLine(kbVerOut)
		slog.Info("kube-burner version", "version", kbVerOut)
	}

	// Print go version
	goVerOut, err := runAndGetStdout("go", "version")
	if err != nil {
		slog.Warn("Failed to get go version", "err", err)
	} else {
		goVerOut = prettySingleLine(goVerOut)
		slog.Info("go version", "version", goVerOut)
	}

	return nil
}

// runAndGetStdout runs a command and returns its stdout as a string
func runAndGetStdout(bin string, arg string) (string, error) {
	cmd := exec.Command(bin, arg)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// prettySingleLine trims and replaces newlines/tabs with spaces for logging
func prettySingleLine(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", "; ")
	s = strings.ReplaceAll(s, "\t", " ")
	return s
}
