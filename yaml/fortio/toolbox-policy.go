package fortio

type ToolboxFQDNPolicy struct {
	Name                     string
	Namespace                string
	AppLabelForPolicyToApply string
}

func (f ToolboxFQDNPolicy) GetTemplate() string {
	return ToolboxFQDNPolicyTemplate
}

const ToolboxFQDNPolicyTemplate = `
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  endpointSelector:
    matchLabels:
      app: {{ .AppLabelForPolicyToApply }}
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
          role: server
  ingress:
  - fromEntities:
    - all
`
