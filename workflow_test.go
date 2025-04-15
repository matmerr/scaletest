package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	flow "github.com/Azure/go-workflow"
	steps "github.com/matmerr/scaletest/steps"
)

func TestWorkflow(t *testing.T) {
	var (
		generateYamls = &steps.GenerateYamlsStep{
			Directory:                     "./output/",
			Namespaces:                    35,
			ServerDeploymentsPerNamespace: 5,
			ServerReplicasPerDeployment:   150,
			ServerServicesPerNamespace:    5,
			ClientDeploymentsPerNamespace: 5,
			ClientReplicasPerDeployment:   150,
		}
		applyYamls = flow.Func("apply all yamls", func(ctx context.Context) error {
			fmt.Println("apply all yamls")
			return nil
		})
	)
	// compose steps into a workflow!
	w := new(flow.Workflow)
	w.Add(
		flow.Steps(applyYamls).DependsOn(generateYamls),

		// other configurations, like retry, timeout, condition, etc.
		flow.Step(applyYamls).
			Retry(func(ro *flow.RetryOption) {
				ro.Attempts = 3 // retry 3 times
			}).
			Timeout(10*time.Minute), // timeout after 10 minutes

		// use Input to change step at runtime
		flow.Step(generateYamls).Input(func(ctx context.Context, g *steps.GenerateYamlsStep) error {
			g.Directory = "./output/"
			return nil
		}),
	)
	// execute the workflow and block until all steps are terminated
	err := w.Do(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Workflow completed successfully")
	}
}
