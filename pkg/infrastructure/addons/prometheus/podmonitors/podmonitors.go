package podmonitors

import (
	context "context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	goosruntime "runtime"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "embed"
)

var PodMonitorDirectory = func() string {
	_, file, _, ok := goosruntime.Caller(0)
	if !ok {
		return filepath.Join(".", "pkg", "infrastructure", "prometheus", "podmonitors")
	}

	return filepath.Dir(file)
}()

//go:embed cilium.yaml
var ciliumPodMonitorYAML []byte

//go:embed hubble.yaml
var hubblePodMonitorYAML []byte

type InstallCiliumPodMonitorStep struct {
	Namespace string
}

func (s *InstallCiliumPodMonitorStep) Do(ctx context.Context) error {
	return applyPodMonitorYAML(ctx, ciliumPodMonitorYAML, s.Namespace)
}

type InstallHubblePodMonitorStep struct {
	Namespace string
}

func (s *InstallHubblePodMonitorStep) Do(ctx context.Context) error {
	return applyPodMonitorYAML(ctx, hubblePodMonitorYAML, s.Namespace)
}

// applyPodMonitorYAML applies a PodMonitor YAML manifest to the given namespace
func applyPodMonitorYAML(ctx context.Context, yamlData []byte, namespace string) error {
	cfg, err := resolveKubeConfig()
	if err != nil {
		slog.Error("Failed to resolve kubeconfig", "err", err)
		return err
	}
	monClient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		slog.Error("Failed to create monitoring client", "err", err)
		return err
	}

	// Decode YAML to PodMonitor object
	jsonData, err := yaml.ToJSON(yamlData)
	if err != nil {
		slog.Error("Failed to convert YAML to JSON", "err", err)
		return err
	}
	scheme := runtime.NewScheme()
	_ = monitoringv1.AddToScheme(scheme)
	serializer := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme, scheme, json.SerializerOptions{Yaml: false, Pretty: false, Strict: false})
	obj, _, err := serializer.Decode(jsonData, nil, nil)
	if err != nil {
		slog.Error("Failed to decode PodMonitor JSON", "err", err)
		return err
	}
	pm, ok := obj.(*monitoringv1.PodMonitor)
	if !ok {
		return fmt.Errorf("decoded object is not a PodMonitor")
	}
	pm.Namespace = namespace // override namespace

	_, err = monClient.MonitoringV1().PodMonitors(namespace).Create(ctx, pm, metav1.CreateOptions{})
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			existing, getErr := monClient.MonitoringV1().PodMonitors(namespace).Get(ctx, pm.Name, metav1.GetOptions{})
			if getErr != nil {
				slog.Error("Failed to get existing PodMonitor", "err", getErr)
				return getErr
			}
			pm.ResourceVersion = existing.ResourceVersion
			_, err = monClient.MonitoringV1().PodMonitors(namespace).Update(ctx, pm, metav1.UpdateOptions{})
			if err != nil {
				slog.Error("Failed to update PodMonitor", "err", err)
				return err
			}
			slog.Info("Updated PodMonitor from YAML", "name", pm.Name)
		} else {
			slog.Error("Failed to create PodMonitor from YAML", "err", err)
			return err
		}
	} else {
		slog.Info("Created PodMonitor from YAML", "name", pm.Name)
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
