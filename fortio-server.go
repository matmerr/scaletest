package main

import (
	"os"
	"text/template"
)

type FortioServerDeployment struct {
	Name                string
	Namespace           string
	Replicas            int
	ServiceBackendLabel string
	AppLabel            string
	NodeSelector        string
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

func CreateServerDeployments(filename string, deployment *FortioServerDeployment) {
	// Create a new template and parse the YAML template into it
	tmpl, err := template.New("server").Parse(fortioServerDeploymentTemplate)
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
