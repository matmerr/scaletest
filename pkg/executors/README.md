# pkg/executors

This directory contains executor implementations for running specific test tools and workflows in a provider-agnostic way. Executors encapsulate the logic for tools such as kube-burner and clusterloader2, and are designed to be modular and extensible.

## Structure

- `kube-burner/`: Executor logic and workflow steps for running kube-burner scenarios.
- `clusterloader2/`: Executor logic and workflow steps for running clusterloader2 scenarios.
- `executor.go`: Common interfaces and types for executors.

## Usage

Executors are invoked by the main test workflows and CI matrix jobs. Each executor is responsible for orchestrating its tool's lifecycle, scenario execution, and artifact collection. Executors are referenced in the CI matrix (see `.github/workflows/go-kind-test.yml`) and can be extended to support additional tools.

## Adding Executors

To add a new executor:
1. Create a new subdirectory (e.g., `my-executor/`).
2. Implement the workflow and steps for that executor.
3. Add the executor to the matrix in `.github/workflows/go-kind-test.yml`.

---

This directory is part of the modular test execution layer for scalable Kubernetes benchmarking and validation.
