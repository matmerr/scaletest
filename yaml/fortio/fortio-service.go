package fortio

type FortioService struct {
	Name                string
	Namespace           string
	TargetPort          string
	ServiceBackendLabel string
}

func (f FortioService) GetTemplate() string {
	return fortioServiceTemplate
}

const fortioServiceTemplate = `
apiVersion: v1
kind: Service
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
  labels:
    svc: {{ .Name }}
spec:
  ports:
    - port: {{ .TargetPort }}
      protocol: TCP
      targetPort: {{ .TargetPort }}
  selector:
    svc: {{ .ServiceBackendLabel }}
  type: ClusterIP

`
