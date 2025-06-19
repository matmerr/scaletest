package largepodcount

import "github.com/matmerr/scaletest/scenarios/clusterloader2/modules"

func NewLargePodCountConfig() Config {
	return Config{
		NamespaceCount:      10, // Default number of namespaces
		PodsPerNamespace:    1000, // Default pods per namespace
		ModulePath:          modules.RelativePath(),
	}
}

type Config struct {
	NamespaceCount    int    `yaml:"namespaceCount,omitempty"`
	PodsPerNamespace  int    `yaml:"podsPerNamespace,omitempty"`
	ModulePath        string `yaml:"modulePath,omitempty"`
}

func (f Config) GetTemplate() string {
	return configTemplate
}

const configTemplate = `
name: largepodcount

namespace:
  number: {{ .NamespaceCount }}

steps:
- module:
    path: {{ .ModulePath }}/cilium.yaml
    params:
      action: start

- name: Create pods
  phases:
  - replicasPerNamespace: {{ .PodsPerNamespace }}
    objectBundle:
      - apiVersion: v1
        kind: Pod
        metadata:
          generateName: largepod-
        spec:
          containers:
          - name: pause
            image: k8s.gcr.io/pause:3.2

- name: Start measurements
  measurements:
    - name: PodStartupLatency
      params:
        action: start
`
