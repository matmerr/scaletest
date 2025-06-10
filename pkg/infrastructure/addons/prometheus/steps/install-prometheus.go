package steps

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type InstallPrometheusStep struct {
	Namespace string
}

func (c *InstallPrometheusStep) Do(ctx context.Context) error {
	releaseName := "prometheus"
	chartName := "kube-prometheus-stack"
	repoName := "prometheus-community"
	repoURL := "https://prometheus-community.github.io/helm-charts"

	// Load kubeconfig and create clientset
	kubeconfigPath := clientcmd.RecommendedHomeFile
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		slog.Error("Failed to load kubeconfig", "err", err)
		return err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		slog.Error("Failed to create clientset", "err", err)
		return err
	}

	// Create namespace if not exists
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), c.Namespace, metav1.GetOptions{})
	if err != nil {
		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: c.Namespace},
		}, metav1.CreateOptions{})
		if err != nil {
			slog.Error("Failed to create namespace", "err", err)
			return err
		}
		slog.Info("Created namespace", "namespace", c.Namespace)
	}

	// Setup Helm
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(genericclioptions.NewConfigFlags(false), c.Namespace, "secrets", func(format string, v ...interface{}) {
		// Instead of format string, use structured logging
		if len(v) > 0 {
			slog.Info("Helm", "msg", fmt.Sprintf(format, v...))
		} else {
			slog.Info("Helm", "msg", format)
		}
	}); err != nil {
		slog.Error("Failed to initialize Helm", "err", err)
		return err
	}

	// Add Helm repo
	chartRepoEntry := &repo.Entry{Name: repoName, URL: repoURL}
	repoFilePath := filepath.Join(os.TempDir(), "helm-repos.yaml")
	repoFile := repo.NewFile()
	repoFile.Update(chartRepoEntry)
	if err := repoFile.WriteFile(repoFilePath, 0644); err != nil {
		slog.Error("Failed to write repo file", "err", err)
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
		"alertmanager": map[string]interface{}{
			"enabled": false,
		},
		"pushgateway": map[string]interface{}{
			"enabled": false,
		},
		"nodeExporter": map[string]interface{}{
			"enabled": false,
		},
		"server": map[string]interface{}{
			"enabled": true,
		},
		"prometheus": map[string]interface{}{
			"enabled": true,
			"service": map[string]interface{}{
				"type": "ClusterIP",
			},
			"prometheusSpec": map[string]interface{}{
				"scrapeInterval":      "15s",
				"evaluationInterval":  "15s",
				"scrape_interval":     "15s",
				"evaluation_interval": "15s",
			},
			"serviceMonitor": map[string]interface{}{
				"selfMonitor":               false,
				"additionalServiceMonitors": []interface{}{},
			},
		},
		"kubeApiServer":         map[string]interface{}{"enabled": false},
		"kubelet":               map[string]interface{}{"enabled": false},
		"kubeControllerManager": map[string]interface{}{"enabled": false},
		"coreDns":               map[string]interface{}{"enabled": false},
		"kubeEtcd":              map[string]interface{}{"enabled": false},
		"kubeScheduler":         map[string]interface{}{"enabled": false},
		"kubeProxy":             map[string]interface{}{"enabled": false},
		"kubeStateMetrics":      map[string]interface{}{"enabled": false},
		"prometheusOperator":    map[string]interface{}{"enabled": true, "serviceMonitor": map[string]interface{}{"selfMonitor": false}},
	}

	histClient := action.NewHistory(actionConfig)
	histClient.Max = 1
	_, err = histClient.Run(releaseName)

	if err != nil {
		// Not found, install
		install := action.NewInstall(actionConfig)
		install.ReleaseName = releaseName
		install.Namespace = c.Namespace
		install.CreateNamespace = false
		install.Wait = true
		install.Timeout = 5 * time.Minute

		rel, err := install.Run(chart, vals)
		if err != nil {
			slog.Error("Failed to install release", "err", err)
			return err
		}
		slog.Info("Installed release", "name", rel.Name)
	} else {
		// Found, upgrade
		upgrade := action.NewUpgrade(actionConfig)
		upgrade.Namespace = c.Namespace
		upgrade.Wait = true
		upgrade.Timeout = 5 * time.Minute

		rel, err := upgrade.Run(releaseName, chart, vals)
		if err != nil {
			slog.Error("Failed to upgrade release", "err", err)
			return err
		}
		slog.Info("Upgraded release", "name", rel.Name)
	}

	// Define and create PodMonitor for Cilium agent
	// (moved to its own step)
	return nil
}
