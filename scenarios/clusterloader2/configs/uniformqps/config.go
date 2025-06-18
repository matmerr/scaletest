package unifromqps

import "github.com/matmerr/scaletest/scenarios/clusterloader2/modules"

func NewUniformQPSConfig() Config {
	return Config{
		ReplicasPerNamespace: 1, // Default number of replicas per namespace
		ModulePath:           modules.RelativePath(),
	}
}

type Config struct {
	ReplicasPerNamespace int    `yaml:"replicasPerNamespace,omitempty"` // Number of replicas per namespace
	ModulePath           string `yaml:"modulePath,omitempty"`           // Path to the module directory
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

- module:
    path: {{ .ModulePath }}/cilium.yaml
    params:
      action: start

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
        Replicas: 150
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

- name: Sleep
  measurements:
  - Identifier: WaitAfterExec
    Method: Sleep
    Params:
      duration: 1m

- module:
    path: {{ .ModulePath }}/cilium.yaml
    params:
      action: gather

`
