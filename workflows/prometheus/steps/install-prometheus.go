package steps

import (
	"context"
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

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

type InstallPrometheusStep struct {
	Namespace string
}

func (c *InstallPrometheusStep) Do(ctx context.Context) error {
	namespace := "monitoring-2"
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
	_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: namespace},
		}, metav1.CreateOptions{})
		if err != nil {
			slog.Error("Failed to create namespace", "err", err)
			return err
		}
		slog.Info("Created namespace", "namespace", namespace)
	}

	// Setup Helm
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(genericclioptions.NewConfigFlags(false), namespace, "secrets", func(format string, v ...interface{}) {
		slog.Info("Helm", "msg", format, "args", v)
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
		install.Namespace = namespace
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
		upgrade.Namespace = namespace
		upgrade.Wait = true
		upgrade.Timeout = 5 * time.Minute

		rel, err := upgrade.Run(releaseName, chart, vals)
		if err != nil {
			slog.Error("Failed to upgrade release", "err", err)
			return err
		}
		slog.Info("Upgraded release", "name", rel.Name)
	}

	// Define PodMonitor object
	podMonitor := &monitoringv1.PodMonitor{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PodMonitor",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cilium-agent-pods",
			Namespace: namespace, // monitoring-2
			Labels: map[string]string{
				"k8s-app": "cilium-agent-pods",
				"release": "prometheus",
			},
		},
		Spec: monitoringv1.PodMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"k8s-app": "cilium",
				},
			},
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: []string{"kube-system"},
			},
			PodMetricsEndpoints: []monitoringv1.PodMetricsEndpoint{{
				Interval: "30s",
				Port:     ptrString("prometheus"),
				MetricRelabelConfigs: []monitoringv1.RelabelConfig{{
					SourceLabels: []monitoringv1.LabelName{"__name__"},
					Regex:        "(.*)",
					Action:       "keep",
				}},
			}},
			JobLabel: "k8s-app",
		},
	}

	// Create Prometheus Operator client
	monClient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		slog.Error("Failed to create monitoring client", "err", err)
		return err
	}
	_, err = monClient.MonitoringV1().PodMonitors(namespace).Create(context.TODO(), podMonitor, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			// Get the existing PodMonitor to get its ResourceVersion
			existing, getErr := monClient.MonitoringV1().PodMonitors(namespace).Get(context.TODO(), podMonitor.Name, metav1.GetOptions{})
			if getErr != nil {
				slog.Error("Failed to get existing PodMonitor", "err", getErr)
				return getErr
			}
			podMonitor.ResourceVersion = existing.ResourceVersion
			_, err = monClient.MonitoringV1().PodMonitors(namespace).Update(context.TODO(), podMonitor, metav1.UpdateOptions{})
			if err != nil {
				slog.Error("Failed to update PodMonitor", "err", err)
				return err
			}
			slog.Info("Updated PodMonitor for Cilium agent")
		} else {
			slog.Error("Failed to create PodMonitor", "err", err)
			return err
		}
	} else {
		slog.Info("Created PodMonitor for Cilium agent")
	}
	return nil
}

// Helper function to convert string to *string
func ptrString(s string) *string {
	return &s
}
