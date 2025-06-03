package scenarios

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/yaml"
	netpolchurn "github.com/matmerr/scaletest/scenarios/kube-burner/netpol-churn"
	"github.com/matmerr/scaletest/steps"
)

func DefaultRun(yamlDirectory string) *flow.Workflow {
	applyYamls := flow.Func("apply all yamls", func(ctx context.Context) error {
		fmt.Println("this is where we apply all yamls")
		return nil
	})

	//installPrometheusStep := &steps.InstallPrometheusStep{
	//	Namespace: "monitoring-2",
	//}

	runKubeBurner := &steps.RunKubeBurner{
		Namespace: "kube-burner",
	}

	w := new(flow.Workflow)
	w.Add(
		flow.Pipe(
			//installPrometheusStep,
			applyYamls,
		),
		// ensure generateYamls is called before applyYamls
		//flow.Steps(applyYamls).DependsOn(installPrometheusStep),

		// applyYamls will need retry``
		flow.Step(applyYamls).
			Retry(func(ro *flow.RetryOption) {
				ro.Attempts = 3 // retry 3 times
			}).
			Timeout(10*time.Minute), // timeout after 10 minutes

		// use Input to change step at runtime
		flow.Step(runKubeBurner).Input(func(ctx context.Context, g *steps.RunKubeBurner) error {

			template := netpolchurn.Config{}

			t := reflect.TypeOf(template)
			fmt.Println("Type Name:", t.Name())
			fmt.Println("Package Path:", t.PkgPath()) // empty if defined in main

			pkgPath := t.PkgPath()
			runningPath, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %w", err)
			}

			// Extract the Go module path from the base path
			const goSrcPrefix = "/src/"
			idx := strings.Index(runningPath, goSrcPrefix)
			if idx == -1 {
				panic("invalid base path, must contain /src/")
			}
			modulePath := runningPath[idx+len(goSrcPrefix):]

			// Strip the module path from the full path
			relPath := strings.TrimPrefix(pkgPath, modulePath+"/")

			fmt.Println(relPath)

			inputPath := runningPath + "/" + relPath

			outputPath := "./output/" + relPath
			generatedConfigPath := inputPath + "/generated" + "/config.yaml"

			template.ResultsDirectory = outputPath
			template.TemplateDirectory = inputPath + "/templates"
			template.MetricsConfigDirectory = inputPath + "/metrics"

			// make output directory if it does not exist
			err = os.MkdirAll(outputPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			err = yaml.CreateYamlFile(generatedConfigPath, template)
			if err != nil {
				return fmt.Errorf("failed to create YAML file: %w", err)
			}

			g.Template = template

			return nil
		}),
	)

	return w
}
