package longrunning

import (
	"context"
	"fmt"
	"time"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/scenarios/longrunning/steps"
)

func LargeScaleWorkflow(yamlDirectory string) *flow.Workflow {
	// clean the directory first
	cleanDirectory := &steps.CleanDirectoryStep{
		Directory: yamlDirectory,
	}

	generateYamls := &steps.GenerateYamlsStep{
		Directory:  yamlDirectory,
		Namespaces: 35,

		ServerDeploymentsPerNamespace: 5,
		ServerReplicasPerDeployment:   150,

		ServerServicesPerNamespace: 5,

		ClientDeploymentsPerNamespace: 5,
		ClientReplicasPerDeployment:   150, // normally 150

		ClientQPS: 1600,
	}

	applyYamls := flow.Func("apply all yamls", func(ctx context.Context) error {
		fmt.Println("this is where we apply all yamls")
		return nil
	})

	w := new(flow.Workflow)
	w.Add(

		flow.Pipe(
			cleanDirectory,
			generateYamls,
			applyYamls,
		),
		// ensure generateYamls is called before applyYamls
		flow.Steps(applyYamls).DependsOn(generateYamls),

		// applyYamls will need retry
		flow.Step(applyYamls).
			Retry(func(ro *flow.RetryOption) {
				ro.Attempts = 3 // retry 3 times
			}).
			Timeout(10*time.Minute), // timeout after 10 minutes

		// use Input to change step at runtime
		flow.Step(generateYamls).Input(func(ctx context.Context, g *steps.GenerateYamlsStep) error {
			return nil
		}),
	)
	return w
}
