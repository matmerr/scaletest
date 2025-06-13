package networkload

type Config struct {
	FortioNamespaces                           int    `yaml:"fortioNamespaces,omitempty"`
	APIServerCallsPerSecond                    int    `yaml:"apiServerCallsPerSecond,omitempty"`
	FortioServerDeployments                    int    `yaml:"fortioServerDeployments,omitempty"`
	FortioClientDeployments                    int    `yaml:"fortioClientDeployments,omitempty"`
	FortioServerReplicasPerDeployment          int    `yaml:"fortioServerReplicasPerDeployment,omitempty"`
	FortioClientReplicasPerDeployment          int    `yaml:"fortioClientReplicasPerDeployment,omitempty"`
	FortioServerDeploymentReplicasPerNamespace int    `yaml:"fortioServerDeploymentReplicasPerNamespace,omitempty"`
	FortioClientDeploymentReplicasPerNamespace int    `yaml:"fortioClientDeploymentReplicasPerNamespace,omitempty"`
	FortioServiceReplicasPerNamespace          int    `yaml:"fortioServiceReplicasPerNamespace,omitempty"`
	FortioClientQueriesPerSecond               int    `yaml:"fortioClientQueriesPerSecond,omitempty"` // Queries per second for Fortio client
	GroupName                                  string `yaml:"groupName,omitempty"`
	OperationTimeout                           string `yaml:"operationTimeout,omitempty"` // Duration for operation timeouts, e.g., "5m" for 5 minutes
}

func NewNetworkLoadConfig() Config {
	return Config{
		FortioNamespaces:                           1,
		APIServerCallsPerSecond:                    10,
		FortioServerDeployments:                    1,
		FortioClientDeployments:                    1,
		FortioServerReplicasPerDeployment:          1,
		FortioClientReplicasPerDeployment:          1,
		FortioServerDeploymentReplicasPerNamespace: 1,
		FortioClientDeploymentReplicasPerNamespace: 1,
		FortioServiceReplicasPerNamespace:          1,
		FortioClientQueriesPerSecond:               100, // Default queries per second for Fortio client
		GroupName:                                  "fortio",
		OperationTimeout:                           "5m", // Default operation timeout
	}
}

func (f Config) GetTemplate() string {
	return configTemplate
}

const configTemplate = `
name: load-config

namespace:
  number: {{ .FortioNamespaces }}
  prefix: ns
  deleteAutomanagedNamespaces: false
  enableExistingNamespaces: true

tuningSets:
  - name: Sequence
    parallelismLimitedLoad:
      parallelismLimit: 1
  - name: DeploymentCreateQps
    qpsLoad:
      qps: {{ .APIServerCallsPerSecond }}
  - name: DeploymentDeleteQps
    qpsLoad:
      qps: {{ .APIServerCallsPerSecond }}

steps:
  - name: Log - fortioNamespaces={{.FortioNamespaces}}, fortioServerDeployments={{.FortioServerDeployments}}, fortioClientDeployments={{.FortioClientDeployments}}, fortioServerReplicasPerDeployment={{.FortioServerReplicasPerDeployment}}, fortioServerDeploymentReplicasPerNamespace={{.FortioServerDeploymentReplicasPerNamespace}}, fortioServiceReplicasPerNamespace=={{.FortioServiceReplicasPerNamespace}}
    measurements:
    - Identifier: Dummy
      Method: Sleep
      Params:
        action: start
        duration: 1ms

  - module:
      path: /modules/measurements.yaml
      params:
        action: start
        group: {{ .GroupName }}

  - module:
      path: /modules/cilium-measurements.yaml
      params:
        action: start

# create resources that won't change
  - module:
      path: /modules/services.yaml
      params:
        actionName: "Creating"
        namespaces: {{ .FortioNamespaces }}
        fortioServiceReplicasPerNamespace: {{ .FortioServiceReplicasPerNamespace }}

  - module:
      path: /modules/ciliumnetworkpolicy.yaml
      params:
        actionName: "Creating"
        namespaces: {{ .FortioNamespaces }}
        Group: {{ .GroupName }}
        cnpsPerNamespace: 1

# deployment creation
  - module:
      path: /modules/reconcile-objects.yaml
      params:
        actionName: "create"
        tuningSet: DeploymentCreateQps
        operationTimeout: {{ .OperationTimeout }}
        Group: {{ .GroupName }}
        namespaces: {{ .FortioNamespaces }}
        fortioServerDeployments: {{ .FortioServerDeployments }}
        fortioClientDeployments: {{ .FortioClientDeployments }}
        fortioServerDeploymentReplicasPerNamespace: {{ .FortioServerDeploymentReplicasPerNamespace }}
        fortioClientDeploymentReplicasPerNamespace: {{ .FortioClientDeploymentReplicasPerNamespace }}
        fortioServerReplicasPerDeployment: {{ .FortioServerReplicasPerDeployment }}
        fortioClientReplicasPerDeployment: {{ .FortioClientReplicasPerDeployment }}
        fortioClientQueriesPerSecond: {{ .FortioClientQueriesPerSecond }}
        deploymentLabel: start

# FIXME sleep intervals
  - name: Log - fortioNamespaces={{ .FortioNamespaces }}, fortioServerDeployments={{ .FortioServerDeployments }}, fortioClientDeployments={{ .FortioClientDeployments }}, fortioServerReplicasPerDeployment={{ .FortioServerReplicasPerDeployment }}
    measurements:
    - Identifier: Dummy
      Method: Sleep
      Params:
        action: start
        duration: 1m30s

  - module:
      path: /modules/cilium-measurements.yaml
      params:
        action: gather
`
