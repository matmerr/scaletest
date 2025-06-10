package ciliumsteps

import (
	"context"
	"log/slog"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// InstallCiliumStep installs or upgrades Cilium using Helm
// Usage: step := &InstallCiliumStep{Namespace: "kube-system"}
//
//	err := step.Do(ctx)
type InstallCiliumStep struct {
	Namespace string
}

func (c *InstallCiliumStep) Do(ctx context.Context) error {
	namespace := c.Namespace
	if namespace == "" {
		namespace = "kube-system"
	}
	releaseName := "cilium"
	chartName := "cilium"
	repoURL := "https://helm.cilium.io/"

	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(genericclioptions.NewConfigFlags(false), namespace, "secrets", slog.Info); err != nil {
		slog.Error("Failed to initialize Helm", "err", err)
		return err
	}

	chartPathOptions := action.ChartPathOptions{RepoURL: repoURL}
	chartPath, err := chartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		slog.Error("Failed to locate chart", "err", err)
		return err
	}
	chart, err := loader.Load(chartPath)
	if err != nil {
		slog.Error("Failed to load chart", "err", err)
		return err
	}

	vals := map[string]interface{}{
		"ipam": map[string]interface{}{
			"mode": "kubernetes",
		},
		"hubble": map[string]interface{}{
			"enabled": true,
			"relay":   map[string]interface{}{"enabled": true},
			"ui":      map[string]interface{}{"enabled": true},
			"metrics": map[string]interface{}{
				"enabled": []interface{}{"dns", "drop", "tcp", "flow", "port-distribution", "icmp", "http"},
			},
		},
		"operator": map[string]interface{}{
			"enabled": true,
			"metrics": map[string]interface{}{
				"enabled": true,
			},
		},
		"metrics": map[string]interface{}{
			"enabled":        true,
			"serviceMonitor": map[string]interface{}{"enabled": true},
		},
		"prometheus": map[string]interface{}{
			"enabled": true,
			"port":    9962,
		},
		// Add or adjust more values as needed for kind compatibility
	}

	histClient := action.NewHistory(actionConfig)
	histClient.Max = 1
	_, err = histClient.Run(releaseName)

	if err != nil {
		install := action.NewInstall(actionConfig)
		install.ReleaseName = releaseName
		install.Namespace = namespace
		install.CreateNamespace = false
		install.Wait = true
		install.Timeout = 5 * time.Minute

		rel, err := install.Run(chart, vals)
		if err != nil {
			slog.Error("Failed to install Cilium", "err", err)
			return err
		}
		slog.Info("Installed Cilium release", "name", rel.Name)
	} else {
		upgrade := action.NewUpgrade(actionConfig)
		upgrade.Namespace = namespace
		upgrade.Wait = true
		upgrade.Timeout = 5 * time.Minute

		slog.Info("Preparing upgrade for release", "name", releaseName)
		rel, err := upgrade.Run(releaseName, chart, vals)
		if err != nil {
			slog.Error("Failed to upgrade Cilium", "err", err)
			return err
		}
		slog.Info("Upgraded Cilium release", "name", rel.Name)
	}

	// No need for manual waitForCiliumReady, Helm will wait
	slog.Info("Cilium is ready")
	return nil
}
