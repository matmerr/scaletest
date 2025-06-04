package steps

import (
	context "context"
	"fmt"
	"os"

	"log/slog"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type InstallCiliumPodMonitorStep struct {
	Namespace string
}

func (s *InstallCiliumPodMonitorStep) Do(ctx context.Context) error {
	return CreateOrUpdatePodMonitor(ctx, PodMonitorSpec{
		Name:      "cilium-metrics",
		Namespace: s.Namespace,
		Labels: map[string]string{
			"k8s-app": "cilium-metrics",
			"release": "prometheus",
		},
		Selector: map[string]string{
			"k8s-app": "cilium",
		},
		Port: "prometheus",
		Relabel: []monitoringv1.RelabelConfig{{
			SourceLabels: []monitoringv1.LabelName{"__name__"},
			Regex:        "(.*)",
			Action:       "keep",
		}},
	})
}

type InstallHubblePodMonitorStep struct {
	Namespace string
}

func (s *InstallHubblePodMonitorStep) Do(ctx context.Context) error {
	return CreateOrUpdatePodMonitor(ctx, PodMonitorSpec{
		Name:      "hubble-metrics",
		Namespace: s.Namespace,
		Labels: map[string]string{
			"k8s-app": "hubble-metrics",
			"release": "prometheus",
		},
		Selector: map[string]string{
			"k8s-app": "cilium",
		},
		Port: "hubble-metrics",
		Relabel: []monitoringv1.RelabelConfig{{
			SourceLabels: []monitoringv1.LabelName{"__name__"},
			Regex:        "(.*)",
			Action:       "keep",
		}},
	})
}

type PodMonitorSpec struct {
	Name      string
	Namespace string
	Labels    map[string]string
	Selector  map[string]string
	Port      string
	Relabel   []monitoringv1.RelabelConfig
}

func CreateOrUpdatePodMonitor(ctx context.Context, spec PodMonitorSpec) error {
	cfg, err := resolveKubeConfig()
	if err != nil {
		slog.Error("Failed to resolve kubeconfig", "err", err)
		return err
	}
	podMonitor := &monitoringv1.PodMonitor{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PodMonitor",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Name,
			Namespace: spec.Namespace,
			Labels:    spec.Labels,
		},
		Spec: monitoringv1.PodMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: spec.Selector,
			},
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: []string{"kube-system"},
			},
			PodMetricsEndpoints: []monitoringv1.PodMetricsEndpoint{{
				Interval:             "30s",
				Port:                 ptrString(spec.Port),
				MetricRelabelConfigs: spec.Relabel,
			}},
			JobLabel: "k8s-app",
		},
	}
	monClient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		slog.Error("Failed to create monitoring client", "err", err)
		return err
	}
	_, err = monClient.MonitoringV1().PodMonitors(spec.Namespace).Create(ctx, podMonitor, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			existing, getErr := monClient.MonitoringV1().PodMonitors(spec.Namespace).Get(ctx, podMonitor.Name, metav1.GetOptions{})
			if getErr != nil {
				slog.Error("Failed to get existing PodMonitor", "err", getErr)
				return getErr
			}
			podMonitor.ResourceVersion = existing.ResourceVersion
			_, err = monClient.MonitoringV1().PodMonitors(spec.Namespace).Update(ctx, podMonitor, metav1.UpdateOptions{})
			if err != nil {
				slog.Error("Failed to update PodMonitor", "err", err)
				return err
			}
			slog.Info("Updated PodMonitor", "name", spec.Name)
		} else {
			slog.Error("Failed to create PodMonitor", "err", err)
			return err
		}
	} else {
		slog.Info("Created PodMonitor", "name", spec.Name)
	}
	return nil
}

// resolveKubeConfig robustly resolves the kubeconfig for the current context (in-cluster, KUBECONFIG, or default path).
func resolveKubeConfig() (*rest.Config, error) {
	// Try in-cluster config first
	cfg, err := rest.InClusterConfig()
	if err == nil {
		slog.Info("Using in-cluster kubeconfig")
		return cfg, nil
	}

	// Try KUBECONFIG env var
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err == nil {
			slog.Info("Using kubeconfig from KUBECONFIG", "path", kubeconfig)
			return cfg, nil
		}
		slog.Warn("Failed to use KUBECONFIG", "path", kubeconfig, "err", err)
	}

	// Try default path
	home, err := os.UserHomeDir()
	if err == nil {
		defaultKubeconfig := home + "/.kube/config"
		cfg, err := clientcmd.BuildConfigFromFlags("", defaultKubeconfig)
		if err == nil {
			slog.Info("Using default kubeconfig path", "path", defaultKubeconfig)
			return cfg, nil
		}
		slog.Warn("Failed to use default kubeconfig path", "path", defaultKubeconfig, "err", err)
	}

	return nil, fmt.Errorf("could not resolve kubeconfig: tried in-cluster, KUBECONFIG, and default path")
}
