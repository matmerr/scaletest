package steps

import (
	context "context"
	log "log/slog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

type InstallCiliumPodMonitorStep struct {
	Namespace string
	KubeConfig *rest.Config
}

func (s *InstallCiliumPodMonitorStep) Do(ctx context.Context) error {
	return ensureCiliumPodMonitor(ctx, s.KubeConfig, s.Namespace)
}

func ensureCiliumPodMonitor(ctx context.Context, cfg *rest.Config, namespace string) error {
	podMonitor := &monitoringv1.PodMonitor{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PodMonitor",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cilium-agent-pods",
			Namespace: namespace,
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

	monClient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		log.Error("Failed to create monitoring client", "err", err)
		return err
	}
	_, err = monClient.MonitoringV1().PodMonitors(namespace).Create(ctx, podMonitor, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			existing, getErr := monClient.MonitoringV1().PodMonitors(namespace).Get(ctx, podMonitor.Name, metav1.GetOptions{})
			if getErr != nil {
				log.Error("Failed to get existing PodMonitor", "err", getErr)
				return getErr
			}
			podMonitor.ResourceVersion = existing.ResourceVersion
			_, err = monClient.MonitoringV1().PodMonitors(namespace).Update(ctx, podMonitor, metav1.UpdateOptions{})
			if err != nil {
				log.Error("Failed to update PodMonitor", "err", err)
				return err
			}
			log.Info("Updated PodMonitor for Cilium agent")
		} else {
			log.Error("Failed to create PodMonitor", "err", err)
			return err
		}
	} else {
		log.Info("Created PodMonitor for Cilium agent")
	}
	return nil
}


type InstallHubblePodMonitorStep struct {
	Namespace string
	KubeConfig *rest.Config
}

func (s *InstallHubblePodMonitorStep) Do(ctx context.Context) error {
	return ensureHubblePodMonitor(ctx, s.KubeConfig, s.Namespace)
}

func ensureHubblePodMonitor(ctx context.Context, cfg *rest.Config, namespace string) error {
	podMonitor := &monitoringv1.PodMonitor{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.coreos.com/v1",
			Kind:       "PodMonitor",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hubble-pods",
			Namespace: namespace,
			Labels: map[string]string{
				"k8s-app": "hubble",
				"release": "prometheus",
			},
		},
		Spec: monitoringv1.PodMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"k8s-app": "hubble",
				},
			},
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: []string{"kube-system"},
			},
			PodMetricsEndpoints: []monitoringv1.PodMetricsEndpoint{{
				Interval: "30s",
				Port:     ptrString("metrics"),
			}},
			JobLabel: "k8s-app",
		},
	}

	monClient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		log.Error("Failed to create monitoring client", "err", err)
		return err
	}
	_, err = monClient.MonitoringV1().PodMonitors(namespace).Create(ctx, podMonitor, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			existing, getErr := monClient.MonitoringV1().PodMonitors(namespace).Get(ctx, podMonitor.Name, metav1.GetOptions{})
			if getErr != nil {
				log.Error("Failed to get existing Hubble PodMonitor", "err", getErr)
				return getErr
			}
			podMonitor.ResourceVersion = existing.ResourceVersion
			_, err = monClient.MonitoringV1().PodMonitors(namespace).Update(ctx, podMonitor, metav1.UpdateOptions{})
			if err != nil {
				log.Error("Failed to update Hubble PodMonitor", "err", err)
				return err
			}
			log.Info("Updated PodMonitor for Hubble")
		} else {
			log.Error("Failed to create Hubble PodMonitor", "err", err)
			return err
		}
	} else {
		log.Info("Created PodMonitor for Hubble")
	}
	return nil
}
