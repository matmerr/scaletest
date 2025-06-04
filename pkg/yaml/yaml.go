package yaml

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/matmerr/scaletest/pkg/utils"
)

type Template interface {
	GetTemplate() string
}

func CreateYamlFile(data Template) (string, error) {
	relPath, err := utils.GetPackagePath(data)
	if err != nil {
		return "", fmt.Errorf("failed to get package path: %w", err)
	}

	outputConfig := relPath + "/config_generated.yaml"

	// make output directory if it does not exist
	err = os.MkdirAll(filepath.Dir(outputConfig), os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create a new template and parse the YAML template into it
	tmpl, err := template.New("yaml").Parse(data.GetTemplate())
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// get parent folder of filename
	err = os.MkdirAll(filepath.Dir(outputConfig), os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(outputConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString("# This file was generated from the corresponding config.go, do not edit directly.\n")
	if err != nil {
		return "", fmt.Errorf("failed to write generated notice: %w", err)
	}

	// Execute the template with the data and write to file
	err = tmpl.Execute(file, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return outputConfig, nil
}
