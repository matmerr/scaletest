package main

import (
	"os"
	"testing"
	"text/template"
)

func TestTest(t *testing.T) {
	server := FortioServerDeployment{
		Name:                "fortio-server-1",
		Namespace:           "fortio-2",
		Replicas:            1,
		ServiceBackendLabel: "fortio-service-1",
		AppLabel:            "fortio-server-1",
		NodeSelector:        "scenario: podcount",
	}

	service := FortioService{
		Name:                "fortio-service-1",
		Namespace:           "fortio-2",
		TargetPort:          "8080",
		ServiceBackendLabel: "fortio-service-1",
	}

	client := FortioClientDeployment{
		Name:         "fortio-client-1",
		Namespace:    "fortio-2",
		Replicas:     1,
		RequestURL:   service.Name,
		RequestPort:  "8080",
		AppLabel:     "fortio-client-1",
		QPS:          "1000",
		NodeSelector: "scenario: podcount",
	}

	namespace := Namespace{
		Name: "fortio-2",
	}

	// Create a new template and parse the YAML template into it
	tmpl, err := template.New("deployment").Parse(fortioServerDeploymentTemplate)
	if err != nil {
		panic(err)
	}

	file, _ := os.Create("./yaml/server.yaml")
	defer file.Close()

	// Execute the template with the data and write to standard output
	err = tmpl.Execute(file, server)
	if err != nil {
		panic(err)
	}
	// #########################
	// Create a new template and parse the YAML template into it
	tmpl, err = template.New("client").Parse(fortioClientDeploymentTemplate)
	if err != nil {
		panic(err)
	}

	file, _ = os.Create("./yaml/client.yaml")
	defer file.Close()

	// Execute the template with the data and write to standard output
	err = tmpl.Execute(file, client)
	if err != nil {
		panic(err)
	}
	// #########################
	// #########################
	// Create a new template and parse the YAML template into it
	tmpl, err = template.New("ns").Parse(namespaceDeploymentTemplate)
	if err != nil {
		panic(err)
	}

	file, _ = os.Create("./yaml/ns.yaml")
	defer file.Close()

	// Execute the template with the data and write to standard output
	err = tmpl.Execute(file, namespace)
	if err != nil {
		panic(err)
	}

	// #########################
	// Create a new template and parse the YAML template into it
	tmpl, err = template.New("svc").Parse(fortioServiceTemplate)
	if err != nil {
		panic(err)
	}

	file, _ = os.Create("./yaml/0-svc.yaml")
	defer file.Close()

	// Execute the template with the data and write to standard output
	err = tmpl.Execute(file, service)
	if err != nil {
		panic(err)
	}
}
