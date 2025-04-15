package main

import (
	"context"
	"fmt"
	"testing"

	flow "github.com/Azure/go-workflow"

	"github.com/matmerr/scaletest/scenarios/largescale"
)

type Init struct {
}

// All required for a step is `Do(context.Context) error`
func (i *Init) Do(ctx context.Context) error {
	fmt.Println("Init")
	return nil
}

func TestWorkflow(t *testing.T) {
	root := new(flow.Workflow).Add(
		flow.Pipe(
			new(Init),
			largescale.LargeScaleWorkflow(),
		),
	)

	err := root.Do(context.Background())
	if err != nil {
		t.Fatalf("failed to run workflow: %v", err)
	}
}
