package fortio

type ToolboxDeployment struct {
	Name         string
	Namespace    string
	Replicas     int
	NodeSelector string
}

func (f ToolboxDeployment) GetTemplate() string {
	return ToolboxDeploymentTemplate
}

const ToolboxDeploymentTemplate = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      role: toolbox
  template:
    metadata:
      labels:
        role: toolbox
    spec:
      nodeSelector:
        {{ .NodeSelector }}
      containers:
        - name: toolbox
          image: acnpublic.azurecr.io/toolbox:latest
          ports:
            - containerPort: 8080
`
