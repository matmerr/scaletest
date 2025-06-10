package prom

import (
	"k8s.io/client-go/tools/clientcmd"

	flow "github.com/Azure/go-workflow"
	podmonitors "github.com/matmerr/scaletest/pkg/infrastructure/addons/prometheus/podmonitors"
	promsteps "github.com/matmerr/scaletest/pkg/infrastructure/addons/prometheus/steps"
	utils "github.com/matmerr/scaletest/pkg/utils"
)

func RunDeployPrometheus() *flow.Workflow {
	w := new(flow.Workflow)
	w.Add(
		flow.Pipe(
			&promsteps.InstallPrometheusStep{
				Namespace: "monitoring",
			},
			&podmonitors.InstallCiliumPodMonitorStep{
				Namespace: "monitoring",
			},
			&podmonitors.InstallHubblePodMonitorStep{
				Namespace: "monitoring",
			},
			&utils.PortForward{
				Namespace:          "monitoring",
				LabelSelector:      "app.kubernetes.io/name=prometheus",
				LocalPort:          "9090",
				RemotePort:         "9090",
				Endpoint:           "http://localhost:9090/metrics",
				KubeConfigFilePath: clientcmd.RecommendedHomeFile,
			},
		),
	)

	return w
}
