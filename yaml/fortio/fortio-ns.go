package fortio

type Namespace struct {
	Name string
}

func (n *Namespace) GetTemplate() string {
	return namespaceDeploymentTemplate
}

const namespaceDeploymentTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .Name }}

`
