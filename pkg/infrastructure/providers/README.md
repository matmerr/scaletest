# pkg/infrastructure/providers

This directory contains provider abstractions and implementations for cluster provisioning and management. Providers are used throughout the codebase to enable scenario and executor logic to work across different platforms (e.g., kind, Azure).

## Structure

- `provider_registry.go`: Maps provider names to their setup workflows in a registry-driven design.
- `provider_steps.go`: Provider interface, ClusterProvider struct, and provider selection logic (including error handling and available options logging).
- `kind/`: Implementation and workflow steps for [kind](https://kind.sigs.k8s.io/) clusters.
- `azure/`: Implementation and workflow steps for Azure Kubernetes Service (AKS) clusters.

## Usage

Use the provider registry and `GetClusterProviderFromEnv` to select providers in a robust, error-checked way. Provider selection is now modular and validated against the registry. If an invalid provider is specified, available options are logged and the test fails early.

Provider options (see registry for full list):
- `kindwithcilium`: Local kind cluster with Cilium
- `aksexistingcluster`: Use an existing Azure AKS cluster

Provider selection is integrated with the matrix-based executor workflow in CI, allowing scenarios to be run against different providers as needed.

## Adding Providers

To add a new provider:
1. Create a new subdirectory (e.g., `mycloud/` for a new provider).
2. Implement workflows and steps for that provider.
3. Register the provider in `provider_registry.go`.

---

This directory is part of the modular provider layer for scalable Kubernetes benchmarking and validation.
