# scaletest

## Overview

This project automates the deployment and benchmarking of Kubernetes clusters using kind, Cilium, Prometheus, and kube-burner. It is designed for robust, repeatable, and CI-friendly performance testing and metrics collection.

## Usage

### Prerequisites
- Go 1.20+
- GNU Make
- Docker (for kind clusters)

### Generating Scenario YAMLs
To generate scenario YAML files for all registered scenarios:

```sh
make generate
```
This will run `go generate ./...` and invoke the scenario YAML generation logic for all scenarios listed in `main.go`.

### Running the Full Workflow
To run the full test workflow (including cluster setup, Prometheus, Cilium, kube-burner, etc.):

```sh
go test -v
```
This will execute the `TestWorkflow` in `main_test.go`, which runs the entire flow end-to-end.

## Adding New Scenarios

1. **Create your scenario**
   - Add a new package under `scenarios/kube-burner/your-scenario/`.
   - Implement a Go type that satisfies the `yaml.YamlGenerator` interface (i.e., has a `GetTemplate()` method returning the scenario YAML template).
   - Provide any supporting files (templates, metrics, etc.) in your scenario directory.

2. **Register your scenario**
   - Open `main.go`.
   - Import your scenario package.
   - Add your scenario to the `scenarios` slice, e.g.:
     ```go
     import yourscenario "github.com/matmerr/scaletest/scenarios/kube-burner/your-scenario"
     // ...
     var scenarios = []yaml.Template{
         netpolchurn.NewNetpolChurnConfig(),
         apiintensive.NewApiIntensiveConfig(),
         yourscenario.NewYourScenarioConfig(),
     }
     ```

3. **Generate YAMLs**
   - Run `make generate` to create the YAML for your new scenario.

## General Flow of kube-burner

1. **YAML Generation**
   - Scenario Go structs/templates are rendered to YAML using the generator logic.
   - Generated YAMLs are written to each scenario's directory with a notice at the top.

2. **Cluster Setup**
   - The workflow installs kind, creates a cluster, installs Cilium, and sets up Prometheus.

3. **Benchmark Execution**
   - kube-burner is run with the generated scenario YAML.
   - Metrics are scraped from Prometheus and written to output files.

4. **Results**
   - Output metrics and summaries are written to the `output/` directory of each scenario.
   - Logs and errors are available for debugging.

## Troubleshooting
- If you see errors about missing scenarios, ensure your scenario is registered in `main.go`.
- If kube-burner fails with object verification errors, check your scenario YAML and cluster state.
- For CI integration, see `.github/workflows/go-kind-test.yml`.

---

For more details, see the code comments and each scenario's directory.
