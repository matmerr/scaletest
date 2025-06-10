# scaletest

## Overview

This project automates the deployment and benchmarking of Kubernetes clusters using modular providers (kind, Azure), scenario-driven test logic, and reusable CI workflows. It is designed for robust, repeatable, and CI-friendly performance testing and metrics collection.

## Project Layout

- `.github/workflows/`: All CI/CD workflows, including reusable workflows for kube-burner and clusterloader2. Main entry points for CI are here.
- `pkg/infrastructure/providers/`: Provider abstractions and implementations (kind, Azure, etc.), including enums and workflow logic for cluster provisioning.
- `pkg/executors/`: Executors for running specific test tools (e.g., kube-burner, clusterloader2) in a provider-agnostic way.
- `scenarios/`: Scenario definitions and scenario-specific assets for each test type (e.g., `scenarios/kube-burner/`, `scenarios/clusterloader2/`).
- `Makefile`: Build, test, and utility targets for local development and CI.
- `.gitignore`: Ignores all `output/` directories under scenarios and other generated files.

## Usage

### Prerequisites
- Go 1.20+
- GNU Make
- Docker (for kind clusters)
- Azure SDK for Go (for Azure provider, if used)

### Downloading Tools
To download required binaries (kind, kube-burner, etc.):

```sh
make tools
```

### Generating Scenario YAMLs
To generate scenario YAML files for all registered scenarios:

```sh
make generate
```

### Running Tests and Workflows
To run all tests (including end-to-end workflows):

```sh
go test -v ./...
```

Or use the Makefile for specific targets:

```sh
make test-kubeburner
make test-clusterloader2
```

### CI/CD
- All CI workflows are defined in `.github/workflows/`.
- Reusable workflows for kube-burner and clusterloader2 are referenced by main test workflows.
- To add or update CI logic, edit or add workflows in `.github/workflows/` and use `workflow_call` for reuse.

## Adding New Scenarios

1. **Create your scenario**
   - Add a new package under the appropriate scenario directory (e.g., `scenarios/kube-burner/your-scenario/`).
   - Implement a Go type that satisfies the scenario interface (e.g., `yaml.YamlGenerator`).
   - Provide supporting files (templates, metrics, etc.) in your scenario directory.

2. **Register your scenario**
   - Open the scenario registry (e.g., `scenarios/scenarios.go`).
   - Import your scenario package and add it to the `Index` slice.

3. **Generate YAMLs**
   - Run `make generate` to create the YAML for your new scenario.

## Adding New Providers

1. **Create a new provider**
   - Add a new subdirectory under `pkg/infrastructure/providers/` (e.g., `mycloud/`).
   - Implement provider logic and workflows.
   - Add a new enum value and update the provider setup mapping in `provider-index.go` and `workflow.go`.

## Output and Results
- Scenario outputs are written to `output/` directories under each scenario.
- CI artifacts and logs are available via GitHub Actions.

## Troubleshooting
- Ensure scenarios are registered in the scenario registry.
- For CI issues, check workflow logs in `.github/workflows/` and ensure reusable workflow references are correct.
- For provider issues, check provider registration and implementation in `pkg/infrastructure/providers/`.

---

For more details, see code comments, scenario directories, and provider documentation in `pkg/infrastructure/providers/`.
