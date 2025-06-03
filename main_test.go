package main

import (
	"context"
	"fmt"
	"testing"

	flow "github.com/Azure/go-workflow"
	scenarios "github.com/matmerr/scaletest/workflows"
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
			scenarios.DefaultRun("testdata/yaml/longrunning"),
		),
	)

	err := root.Do(context.Background())
	if err != nil {
		t.Fatalf("failed to run workflow: %v", err)
	}
}
