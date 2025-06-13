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
)

func RunScenarios(provider providers.Provider, executor executors.Executor) error {

	// retrieve all of the yaml templates for the passed in scenario, to be handed to the executor below
	scenarioTemplates := executor.GetScenarioTemplates()

	// get a list of all scenarios paths on disk based off their types
	for _, scenario := range scenarioTemplates {
		t := reflect.TypeOf(scenario)
		pkgPath := filepath.Base(t.PkgPath())
		slog.Info("Scenario", slog.String("name", pkgPath), slog.String("path", t.PkgPath()))
	}

	// Add the executor steps to the to-be-run workflow, use this to generate the command passed to the
	// executor with the supplied config yaml
	scenarioSteps := make([]flow.Steper, 0, len(scenarioTemplates))
	for _, scenario := range scenarioTemplates {

		// get the steps within the context of the executor, at the end of the day it'll do something like
		// generate the config_generated.yaml file, and then pass that generated yaml to the
		// executors command, like "kube-burner run -c config_generated.yaml"
		scenarioSteps = append(scenarioSteps, executor.GetRunWorkflow(scenario))
	}

	// Use batch pipe to stitch setup and scenarios
	root := new(flow.Workflow).Add(

		// we use batch pipe to maintain that A->B->C for the nested steps
		flow.BatchPipe(

			// get the cluster setup/install steps for the provider,
			// usually like create cluster.
			provider.GetSteps(),

			// get the steps that were added to the executor, could be install prometheus, addons, etc.
			// non provider specific steps that are more relevant to the test execution
			executor.GetSetupWorkflow(),

			// and finally add all of the to-be-run scenario steps generated from the executor
			// to the final part of the workflow
			flow.Pipe(scenarioSteps...),
		),
	)

	// run all the steps from start to finish
	if err := root.Do(context.Background()); err != nil {
		return fmt.Errorf("failed to run workflow: %w", err)
	}

	return nil
}
