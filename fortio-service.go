package main

import (
	"os"
	"text/template"
)

type FortioService struct {
	Name                string
	Namespace           string
	TargetPort          string
	ServiceBackendLabel string
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

func CreateServiceDeployments(filename string, deployment *FortioService) {
	// Create a new template and parse the YAML template into it
	tmpl, err := template.New("service").Parse(fortioServiceTemplate)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Execute the template with the data and write to standard output
	err = tmpl.Execute(file, deployment)
	if err != nil {
		panic(err)
	}
}
