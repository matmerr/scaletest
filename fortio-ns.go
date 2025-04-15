package main

import (
	"fmt"
	"os"
	"text/template"
)

type Namespace struct {
	Name string
}

const namespaceDeploymentTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .Name }}

`

func CreateNamespaces(directory string, numberOfNamespaces int) {
	for i := 0; i < numberOfNamespaces; i++ {
		namespace := Namespace{
			Name: "fortio-" + fmt.Sprint(i),
		}

		// Create a new template and parse the YAML template into it
		tmpl, err := template.New("namespace").Parse(namespaceDeploymentTemplate)
		if err != nil {
			panic(err)
		}

		if err := os.MkdirAll(fmt.Sprintf("%s/%d/", directory, i), os.ModePerm); err != nil {
			panic(err)
		}

		file, _ := os.Create(fmt.Sprintf("%s/%d/0-ns-%d.yaml", directory, i, i))
		defer file.Close()

		// Execute the template with the data and write to standard output
		err = tmpl.Execute(file, namespace)
		if err != nil {
			panic(err)
		}
	}
}
