package executors

import (
	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/yaml"
)

type Executor interface {
	// RunWorkflow returns a workflow to run the executor
	GetRunWorkflow(templateConfig yaml.Template) *flow.Workflow
	GetSetupWorkflow() flow.AddSteps
	GetScenarioTemplates() ([]yaml.Template, error)
}
