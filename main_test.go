package main

import (
	"context"
	"fmt"
	"testing"

	flow "github.com/Azure/go-workflow"
	"github.com/matmerr/scaletest/scenarios/longrunning"
)

type Welcome struct {
}

// All required for a step is `Do(context.Context) error`
func (i *Welcome) Do(ctx context.Context) error {
	fmt.Println("starting workflow")
	return nil
}

func TestWorkflow(t *testing.T) {
	root := new(flow.Workflow).Add(
		flow.Pipe(
			new(Welcome),
			longrunning.LargeScaleWorkflow("./output"),
		),
	)

	err := root.Do(context.Background())
	if err != nil {
		t.Fatalf("failed to run workflow: %v", err)
	}
}
