package netpolchurn

import (
	"os"
	"strconv"
)

func NewNetpolChurnConfig() Config {
	jobIterationsStr := os.Getenv("NETPOL_CHURN_JOB_ITERATIONS")
	jobIterations := 5

	if jobIterationsStr != "" {
		if val, err := strconv.Atoi(jobIterationsStr); err == nil {
			jobIterations = val
		}
	}

	return Config{
		JobIterations: jobIterations,
	}
}

type Config struct {
	JobIterations int `yaml:"jobIterations,omitempty"` // Number of iterations for each job
}

func (f Config) GetTemplate() string {
	return configTemplate
}

const configTemplate = `
metricsEndpoints:
  - endpoint: http://localhost:9090 # URL to your Prometheus server
    step: 30s # Scrape interval (optional, default is 30s)
    skipTLSVerify: true # Skip TLS certificate verification (optional)
    metrics:
      - ./metrics/metrics-cilium.yaml # Reference to your custom metrics profile file
    indexer:
      type: local # Store results locally (can also be "opensearch" or "elastic")
      metricsDirectory: ./output/

jobs:
  - name: network-policy-perf-pods
    namespace: network-policy-perf
    jobIterations: {{ .JobIterations }}
    qps: 20
    burst: 20
    namespacedIterations: true
    podWait: false
    waitWhenFinished: true
    preLoadImages: false
    preLoadPeriod: 1s
    skipIndexing: true
    namespaceLabels:
      kube-burner.io/skip-networkpolicy-latency: true
      security.openshift.io/scc.podSecurityLabelSync: false
      pod-security.kubernetes.io/enforce: privileged
      pod-security.kubernetes.io/audit: privileged
      pod-security.kubernetes.io/warn: privileged
    objects:
      - objectTemplate: ./templates/pod.yml
        replicas: 2

      - objectTemplate: ./templates/np-deny-all.yml
        replicas: 1

      - objectTemplate: ./templates/np-allow-from-proxy.yml
        replicas: 1

  - name: network-policy-perf
    namespace: network-policy-perf
    jobIterations: 9
    qps: 20
    burst: 20
    namespacedIterations: true
    podWait: false
    waitWhenFinished: true
    preLoadImages: false
    preLoadPeriod: 15s
    jobPause: 15s
    cleanup: false
    namespaceLabels:
      security.openshift.io/scc.podSecurityLabelSync: false
      pod-security.kubernetes.io/enforce: privileged
      pod-security.kubernetes.io/audit: privileged
      pod-security.kubernetes.io/warn: privileged
    objects:
      - objectTemplate: ./templates/ingress-np.yml
        replicas: 1
        inputVars:
          namespaces: 9
          pods_per_namespace: 2
          netpols_per_namespace: 1
          local_pods: 1
          pod_selectors: 1
          single_ports: 1
          port_ranges: 1
          peer_namespaces: 2
          peer_pods: 2
          cidr_rules: 1

`
