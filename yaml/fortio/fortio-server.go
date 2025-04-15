package fortio

type FortioServerDeployment struct {
	Name                string
	Namespace           string
	Replicas            int
	ServiceBackendLabel string
	AppLabel            string
	NodeSelector        string
}

func (f FortioServerDeployment) GetTemplate() string {
	return fortioServerDeploymentTemplate
}

const fortioServerDeploymentTemplate = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      role: server
      svc: {{ .ServiceBackendLabel }}
      app: {{ .AppLabel }}
  template:
    metadata:
      labels:
        role: server
        svc: {{ .ServiceBackendLabel }}
        app: {{ .AppLabel }}
    spec:
      nodeSelector:
        {{ .NodeSelector }}
      containers:
        - name: fortio
          image: acnpublic.azurecr.io/fortio:latest
          args: ["server"]
          ports:
            - containerPort: 8080
`
