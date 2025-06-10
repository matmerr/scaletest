package yaml

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"reflect"

	flow "github.com/Azure/go-workflow"
)

func GenerateYaml(templateConfig Template) *flow.Workflow {
	w := new(flow.Workflow).Add(
		flow.Step(&CreateYaml{
			InputConfig: templateConfig,
		}),
	)

	return w
}

type CreateYaml struct {
	InputConfig  Template
	OutputConfig string
}

func (c *CreateYaml) Do(ctx context.Context) error {
	configpath, err := CreateYamlFile(c.InputConfig)
	if err != nil {
		return fmt.Errorf("failed to create YAML file: %w", err)
	}

	t := reflect.TypeOf(c.InputConfig)
	pkgPath := filepath.Base(t.PkgPath())

	slog.Info("config generated:", slog.String("scenario", pkgPath), slog.String("path", configpath))
	c.OutputConfig = configpath
	return nil
}

func NewCreateYaml(inputConfig Template) *CreateYaml {
	return &CreateYaml{
		InputConfig: inputConfig,
	}
}
