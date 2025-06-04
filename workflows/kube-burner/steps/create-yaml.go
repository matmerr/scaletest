package steps

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"reflect"

	"github.com/matmerr/scaletest/pkg/yaml"
)

type CreateYaml struct {
	InputConfig  yaml.Template
	OutputConfig string
}

func (c *CreateYaml) Do(ctx context.Context) error {
	configpath, err := yaml.CreateYamlFile(c.InputConfig)
	if err != nil {
		return fmt.Errorf("failed to create YAML file: %w", err)
	}

	t := reflect.TypeOf(c.InputConfig)
	pkgPath := filepath.Base(t.PkgPath())

	slog.Info("config generated:", slog.String("scenario", pkgPath),  slog.String("path", configpath))
	c.OutputConfig = configpath
	return nil
}

func NewCreateYaml(inputConfig yaml.Template) *CreateYaml {
	return &CreateYaml{
		InputConfig: inputConfig,
	}
}
