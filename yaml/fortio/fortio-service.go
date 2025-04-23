package fortio

type FortioService struct {
	Name                string
	Namespace           string
	TargetPort          int
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
spec:
  ports:
    - port: {{ .TargetPort }}
      protocol: TCP
      targetPort: {{ .TargetPort }}
  selector:
    svc: {{ .ServiceBackendLabel }}
  type: ClusterIP

`
