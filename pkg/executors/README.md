# pkg/executors

This directory contains executor implementations for running specific test tools and workflows in a provider-agnostic way. Executors encapsulate the logic for tools such as kube-burner and clusterloader2, and are designed to be modular and extensible.

## Structure

- `kube-burner/`: Executor logic and workflow steps for running kube-burner scenarios. Accepts scenario structs from the scenario registry.
- `clusterloader2/`: Executor logic and workflow steps for running clusterloader2 scenarios. Accepts scenario structs from the scenario registry.
- `executor.go`: Common interfaces and types for executors.

## Usage

Executors are invoked by the main test workflows and CI jobs. Each executor is responsible for orchestrating its tool's lifecycle, scenario execution, and artifact collection. Executors are referenced in the CI workflows and can be extended to support additional tools. Executors now accept scenario structs and provider objects from the registries, ensuring robust error handling and modularity.

## Adding Executors

To add a new executor:
1. Create a new subdirectory (e.g., `my-executor/`).
2. Implement the workflow and steps for that executor.
3. Register the executor in the CI workflows.

---

This directory is part of the modular test execution layer for scalable Kubernetes benchmarking and validation.
