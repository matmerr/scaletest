package scenarios

import (
	"context"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/yaml"
	kb "github.com/matmerr/scaletest/workflows/kube-burner/steps"
)

func RunKubeBurner(templateConfig yaml.Template) *flow.Workflow {
	createYaml := &kb.CreateYaml{
		InputConfig: templateConfig,
	}

	runKubeBurner := &kb.KubeBurner{
		Namespace: "kube-burner",
	}

	w := new(flow.Workflow).Add(
		flow.Step(createYaml),
		flow.Step(runKubeBurner).DependsOn(createYaml).Input(func(ctx context.Context, g *kb.KubeBurner) error {
			runKubeBurner.ConfigPath = createYaml.OutputConfig
			return nil
		}),
	)

	return w
}

func InstallKubeBurner() *flow.Workflow {
	w := new(flow.Workflow).Add(
		flow.Pipe(
			&kb.InstallKubeBurner{},
		),
	)

	return w
}

func GenerateYaml(templateConfig yaml.Template) *flow.Workflow {
	w := new(flow.Workflow).Add(
		flow.Step(&kb.CreateYaml{
			InputConfig: templateConfig,
		}),
	)

	return w
}
