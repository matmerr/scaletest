package main

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"reflect"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/pkg/executors"
	"github.com/matmerr/scaletest/pkg/infrastructure/providers"
	"github.com/matmerr/scaletest/pkg/yaml"
)

func RunScenarios(provider providers.Provider, executor executors.Executor, scenarios []yaml.Template) error {

	// get a list of all scenarios paths on disk based off their types
	for _, scenario := range scenarios {
		t := reflect.TypeOf(scenario)
		pkgPath := filepath.Base(t.PkgPath())
		slog.Info("Scenario", slog.String("name", pkgPath), slog.String("path", t.PkgPath()))
	}

	// Add the executor steps to the to-be-run workflow
	scenarioSteps := make([]flow.Steper, 0, len(scenarios))
	for _, scenario := range scenarios {

		// get the steps within the context of the executor, at the end of the day it'll do something like
		// generate the config_generated.yaml file, and then pass that generated yaml to the
		// executors command, like "kube-burner run -c config_generated.yaml"
		scenarioSteps = append(scenarioSteps, executor.GetRunWorkflow(scenario))
	}

	// stitch all of the scenario steps together into a single pipe
	scenarioPipe := flow.Pipe(scenarioSteps...)

	// Use batch pipe to stitch setup and scenarios
	root := new(flow.Workflow).Add(

		// we use batch pipe to maintain that A->B->C for the nested steps
		flow.BatchPipe(

			// get the cluster setup/install steps for the provider
			providers.GetProviderSetupSteps(provider),

			// get the steps that were added to the executor, could be install prometheus, addons, etc.
			executor.GetSetupWorkflow(),

			// get the all of the scenarios as steps
			scenarioPipe,
		),
	)

	// run all the steps from start to finish
	if err := root.Do(context.Background()); err != nil {
		return fmt.Errorf("failed to run workflow: %w", err)
	}

	return nil
}
