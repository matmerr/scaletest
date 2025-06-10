package steps

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/matmerr/scaletest/pkg/infrastructure/addons/prometheus/podmonitors"
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

	cmdArgs := []string{
		"--v=2",
		"--kubeconfig", c.Kubeconfig,
		"--testconfig", c.ConfigPath,
		"--provider", c.Provider,
		"--report-dir", outputDir,
		"--apiserver-pprof-by-client-enabled", "false",
		"--enable-prometheus-server",
		"--prometheus-storage-class-provisioner", "standard",
		"--prometheus-pvc-storage-class", "standard",
		"--prometheus-storage-class-volume-type", "standard",
		"--prometheus-additional-monitors-path", podmonitors.PodMonitorDirectory,
		"--prometheus-scrape-master-kubelets=false",
		"--prometheus-scrape-kubelets=false",
		"--prometheus-scrape-metrics-server=false",
		"--prometheus-scrape-kube-state-metrics=false",
		"--prometheus-scrape-kubelets=false",
		"--prometheus-scrape-kube-proxy=false",
	}

	deleteCtx, deleteCancel := context.WithCancel(ctx)
	deletionDone := make(chan struct{})
	go func() {
		defer close(deletionDone)
		deleteMasterServiceMonitor(deleteCtx, c.Kubeconfig)
	}()

	slog.Info("Running command", "cmd", binPath, "args", cmdArgs)
	cmd := exec.CommandContext(ctx, binPath, cmdArgs...)
	cmd.Dir = filepath.Dir(c.ConfigPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	slog.Info("Executing command string", "cmd", cmd.String())

	err = cmd.Run()
	deleteCancel()
	<-deletionDone

	if err != nil {
		slog.Error("failed to run clusterloader2", "err", err)
		return err
	}
	return nil
}

// deleteMasterServiceMonitor watches for the ServiceMonitor named "master" in the monitoring namespace and deletes it as soon as it appears.
func deleteMasterServiceMonitor(ctx context.Context, kubeconfig string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		slog.Error("Failed to build kubeconfig for client-go", "err", err)
		return
	}
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		slog.Error("Failed to create dynamic client", "err", err)
		return
	}
	gvr := schema.GroupVersionResource{
		Group:    "monitoring.coreos.com",
		Version:  "v1",
		Resource: "servicemonitors",
	}
	resource := dyn.Resource(gvr).Namespace("monitoring")
	watcher, err := resource.Watch(ctx, v1.ListOptions{})
	if err != nil {
		slog.Error("Failed to start watch on ServiceMonitors", "err", err)
		return
	}
	defer watcher.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-watcher.ResultChan():
			if event.Type == watch.Added || event.Type == watch.Modified {
				obj := event.Object.(v1.Object)
				if obj.GetName() == "master" {

					slog.Info("Deleting ServiceMonitor 'master', cl2 creates it and it's not relevant")
					err := resource.Delete(ctx, "master", v1.DeleteOptions{})
					if err != nil && !os.IsNotExist(err) {
						slog.Error("Failed to delete ServiceMonitor 'master'", "err", err)
					}
				}
			}
		}
	}
}
