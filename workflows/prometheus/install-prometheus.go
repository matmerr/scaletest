package steps

import (
	"k8s.io/client-go/tools/clientcmd"

	flow "github.com/Azure/go-workflow"
	utils "github.com/matmerr/scaletest/pkg/utils"
	promsteps "github.com/matmerr/scaletest/workflows/prometheus/steps"
)

func RunConfigurePrometheus() *flow.Workflow {
	w := new(flow.Workflow)
	w.Add(
		flow.Pipe(
			&promsteps.InstallPrometheusStep{
				Namespace: "monitoring-2",
			},
			&promsteps.InstallCiliumPodMonitorStep{
				Namespace: "monitoring-2",
			},
			&promsteps.InstallHubblePodMonitorStep{
				Namespace: "monitoring-2",
			},
			&utils.PortForward{
				Namespace:          "monitoring-2",
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
