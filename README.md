# scaletest

## Overview

This project automates the deployment and benchmarking of Kubernetes clusters using modular providers (kind, Azure), scenario-driven test logic, and matrix-based CI workflows. It is designed for robust, repeatable, and CI-friendly performance testing and metrics collection.

## Project Layout

- `.github/workflows/`: All CI/CD workflows, including matrix-based workflows for executors (kube-burner, clusterloader2). Main entry points for CI are here.
- `pkg/infrastructure/providers/`: Provider abstractions and implementations (kind, Azure, etc.), including provider registry and workflow logic for cluster provisioning. Provider selection is now registry-driven and error-checked.
- `pkg/executors/`: Executors for running specific test tools (e.g., kube-burner, clusterloader2) in a provider-agnostic way. Executors are modular and accept scenario structs from the scenario registry.
- `scenarios/`: Scenario definitions and scenario-specific assets for each test type (e.g., `scenarios/kube-burner/`, `scenarios/clusterloader2/`). Scenarios are registered in a struct-based registry and selected via environment variable.
- `Makefile`: Build, test, and utility targets for local development and CI.
- `.gitignore`: Ignores all `output/` directories under scenarios and other generated files.
- `artifacts/`: Output and log artifacts from test runs and CI jobs.

## Usage

### Prerequisites
- Go 1.20+
- GNU Make
- Docker (for kind clusters)
- Azure CLI (for Azure provider, if used)

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

- Scenario and provider selection is now handled via environment variables (`KB_SCENARIO`, `CL2_SCENARIO`, `CLUSTER_PROVIDER`) and validated against the scenario/provider registries. If an invalid value is provided, available options are logged and the test fails early.
- All scenario and provider logic is modular, registry-driven, and CI-friendly.

---

For more details, see the `README.md` files in each subdirectory.
