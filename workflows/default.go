package scenarios

import (
	"context"
	"fmt"
	"time"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/steps"
)

func DefaultRun(yamlDirectory string) *flow.Workflow {
	applyYamls := flow.Func("apply all yamls", func(ctx context.Context) error {
		fmt.Println("this is where we apply all yamls")
		return nil
	})

	installPrometheusStep := &steps.InstallPrometheusStep{
		Namespace: "monitoring-2",
	}

	runKubeBurner := &steps.RunKubeBurner{
		Namespace: "kube-burner",
	}

	w := new(flow.Workflow)
	w.Add(
		flow.Pipe(
			installPrometheusStep,
			applyYamls,
		),
		// ensure generateYamls is called before applyYamls
		flow.Steps(applyYamls).DependsOn(installPrometheusStep),

		// applyYamls will need retry``
		flow.Step(applyYamls).
			Retry(func(ro *flow.RetryOption) {
				ro.Attempts = 3 // retry 3 times
			}).
			Timeout(10*time.Minute), // timeout after 10 minutes

		// use Input to change step at runtime
		flow.Step(runKubeBurner).Input(func(ctx context.Context, g *steps.RunKubeBurner) error {
			return nil
		}),
	)

	return w
}
