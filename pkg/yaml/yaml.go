package yaml

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type YamlGenerator interface {
	GetTemplate() string
}

func CreateYamlFile(filename string, data YamlGenerator) error {
	// Create a new template and parse the YAML template into it
	tmpl, err := template.New("yaml").Parse(data.GetTemplate())
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// get parent folder of filename
	err = os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Execute the template with the data and write to standard output
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
