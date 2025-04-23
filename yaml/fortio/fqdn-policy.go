package fortio

type ClientFQDNPolicy struct {
	Name                   string
	Namespace              string
	TargetPort             string
	AppLabel               string
	TargetServiceName      string
	TargetServiceNamespace string
	ToPort                 int
}

func (f ClientFQDNPolicy) GetTemplate() string {
	return ciliumFQDNPolicyTemplate
}

const ciliumFQDNPolicyTemplate = `
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  endpointSelector:
    matchLabels:
      app: {{ .AppLabel }}
  egress:
    - toFQDNs:
        - matchName: {{ .TargetServiceName }}.{{ .TargetServiceNamespace }}.svc.cluster.local
      toPorts:
        - ports:
            - port: {{ .ToPort }}
              protocol: TCP
  egressDeny:
    - toEntities:
        - all

`
