package steps

import (
	context "context"
	log "log/slog"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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
		log.Error("Failed to resolve kubeconfig", "err", err)
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
		log.Error("Failed to create monitoring client", "err", err)
		return err
	}
	_, err = monClient.MonitoringV1().PodMonitors(spec.Namespace).Create(ctx, podMonitor, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			existing, getErr := monClient.MonitoringV1().PodMonitors(spec.Namespace).Get(ctx, podMonitor.Name, metav1.GetOptions{})
			if getErr != nil {
				log.Error("Failed to get existing PodMonitor", "err", getErr)
				return getErr
			}
			podMonitor.ResourceVersion = existing.ResourceVersion
			_, err = monClient.MonitoringV1().PodMonitors(spec.Namespace).Update(ctx, podMonitor, metav1.UpdateOptions{})
			if err != nil {
				log.Error("Failed to update PodMonitor", "err", err)
				return err
			}
			log.Info("Updated PodMonitor", "name", spec.Name)
		} else {
			log.Error("Failed to create PodMonitor", "err", err)
			return err
		}
	} else {
		log.Info("Created PodMonitor", "name", spec.Name)
	}
	return nil
}

type InstallCiliumPodMonitorStep struct {
	Namespace string
}

func (s *InstallCiliumPodMonitorStep) Do(ctx context.Context) error {
	return CreateOrUpdatePodMonitor(ctx, PodMonitorSpec{
		Name:      "cilium-agent-pods",
		Namespace: s.Namespace,
		Labels: map[string]string{
			"k8s-app": "cilium-agent-pods",
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
		Name:      "hubble-pods",
		Namespace: s.Namespace,
		Labels: map[string]string{
			"k8s-app": "hubble",
			"release": "prometheus",
		},
		Selector: map[string]string{
			"k8s-app": "hubble",
		},
		Port:    "metrics",
		Relabel: nil,
	})
}

// resolveKubeConfig resolves the kubeconfig for the current context.
func resolveKubeConfig() (*rest.Config, error) {
	return clientcmd.BuildConfigFromFlags("", "")
}
