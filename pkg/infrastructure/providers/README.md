# pkg/infrastructure/providers

This directory contains provider abstractions and implementations for cluster provisioning and management. Providers are used throughout the codebase to enable scenario and executor logic to work across different platforms (e.g., kind, Azure).

## Structure

- `workflow.go`: Defines provider enums and shared provider logic.
- `provider-index.go`: Maps provider enums to their setup workflows.
- `kind/`: Implementation and workflow steps for [kind](https://kind.sigs.k8s.io/) clusters.
- `azure/`: Implementation and workflow steps for Azure Kubernetes Service (AKS) clusters.

## Usage

Use the `Provider` type and its enums (e.g., `ProviderKindWithCilium`, `ProviderAzureExistingCluster`) to refer to supported providers in a type-safe way. All scenario and executor logic references these abstractions to ensure modularity and extensibility. Each provider subdirectory contains workflows and steps for provisioning, configuring, and managing clusters on that platform.

## Adding Providers

To add a new provider:
1. Create a new subdirectory (e.g., `mycloud/` for a new provider).
2. Implement workflows and steps for that provider.
3. Add a new enum value to `provider-index.go` and update the `ProviderSetupSteps` map in `workflow.go`.

This directory is the canonical place for adding new providers. All scenario execution and CI workflows reference these abstractions to ensure consistent cluster setup and management.

---

This directory is part of the infrastructure automation layer for scalable Kubernetes testing and benchmarking.
