# scenarios

This directory contains scenario definitions and assets for each supported test executor. Scenarios define the workload, configuration, and templates used by executors such as kube-burner and clusterloader2.

## Structure

- `kube-burner/`: Scenarios, templates, and metrics for kube-burner-based tests. Scenarios are registered in a struct-based registry and selected via the `KB_SCENARIO` environment variable.
- `clusterloader2/`: Scenarios and templates for clusterloader2-based tests. Scenarios are registered in a struct-based registry and selected via the `CL2_SCENARIO` environment variable.
- Each scenario subdirectory contains a `config.go` (scenario definition), `config_generated.yaml` (generated YAML), and supporting files (templates, metrics, etc.).

## Usage

Scenarios are registered in their respective `scenarios.go` files and referenced by executors. Scenario YAMLs are generated via `make generate` and are used as input for test runs. Scenario selection is now robust and error-checked: if an invalid scenario is specified, available options are logged and the test fails early. Outputs are collected and zipped as artifacts in CI workflows.

## Adding Scenarios

To add a new scenario:
1. Create a new subdirectory under the appropriate executor (e.g., `kube-burner/my-scenario/`).
2. Implement the scenario definition in Go and provide any required templates or metrics.
3. Register the scenario in the relevant `scenarios.go` file using the struct-based registry.
4. Run `make generate` to produce the scenario YAML.

---

This directory is the canonical source for all test scenarios and assets used in automated benchmarking and validation.
