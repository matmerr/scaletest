package fortio

type ClientFQDNPolicy struct {
	Name                string
	Namespace           string
	TargetPort          string
	AppLabel            string
	ServiceBackendLabel string
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
    - toEndpoints:
      - matchLabels:
          k8s:io.kubernetes.pod.namespace: kube-system
          k8s:k8s-app: node-local-dns
      toPorts:
      - ports:
        - port: "53"
          protocol: ANY
        rules:
            dns:
            - matchPattern: '*'
    - toFQDNs:
        - matchPattern: '*'
    - toEndpoints:
        - matchLabels:
            svc: {{ .ServiceBackendLabel }}

  ingress:
  - fromEntities:
    - all
`
