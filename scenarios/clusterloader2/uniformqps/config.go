package unifromqps

func NewUniformQPSConfig() Config {
	return Config{
		ReplicasPerNamespace: 1, // Default number of replicas per namespace
	}
}

type Config struct {
	ReplicasPerNamespace int `yaml:"replicasPerNamespace,omitempty"` // Number of replicas per namespace
}

func (f Config) GetTemplate() string {
	return configTemplate
}

const configTemplate = `
name: test

namespace:
  number: 1

tuningSets:
- name: Uniform1qps
  qpsLoad:
    qps: 1

steps:
- name: Start measurements
  measurements:
  - Identifier: PodStartupLatency
    Method: PodStartupLatency
    Params:
      action: start
      labelSelector: group = test-pod
      threshold: 20s
  - Identifier: WaitForControlledPodsRunning
    Method: WaitForControlledPodsRunning
    Params:
      action: start
      apiVersion: apps/v1
      kind: Deployment
      labelSelector: group = test-deployment
      operationTimeout: 120s
  - Identifier: CiliumBPFMapPressure
    Method: GenericPrometheusQuery
    Params:
      action: start
      metricName: Cilium BPF Map Pressure
      metricVersion: v1
      unit: "%"
      dimensions:
        - map_name
      queries:
        - name: Max BPF Map Pressure
          query: max(cilium_bpf_map_pressure)
          threshold: 90
- name: Create deployment
  phases:
  - namespaceRange:
      min: 1
      max: 1
    replicasPerNamespace: {{ .ReplicasPerNamespace }}
    tuningSet: Uniform1qps
    objectBundle:
    - basename: test-deployment
      objectTemplatePath: "./templates/deployment.yaml"
      templateFillMap:
        Replicas: 10
- name: Wait for pods to be running
  measurements:
  - Identifier: WaitForControlledPodsRunning
    Method: WaitForControlledPodsRunning
    Params:
      action: gather
- name: Measure pod startup latency
  measurements:
  - Identifier: PodStartupLatency
    Method: PodStartupLatency
    Params:
      action: gather
- name: Measure Cilium Metrics
  measurements:
    - Identifier: CiliumBPFMapPressure
      Method: GenericPrometheusQuery
      Params:
        action: gather
        metricName: Cilium BPF Map Pressure
        metricVersion: v1
        unit: "%"
        dimensions:
          - map_name
        queries:
          - name: Max BPF Map Pressure
            query: max(cilium_bpf_map_pressure)
            threshold: 90

`
