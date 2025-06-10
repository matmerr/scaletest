# pkg/infrastructure/providers

This directory contains provider abstractions and implementations for cluster provisioning and management.

## Structure

- `workflow.go`: Defines provider enums and shared provider logic.
- `provider-index.go`: Maps provider enums to their setup workflows.
- `kind/`: Implementation and workflow steps for [kind](https://kind.sigs.k8s.io/) clusters.
- `azure/`: Implementation and workflow steps for Azure Kubernetes Service (AKS) clusters.

## Usage

Use the `Provider` type and its enums (e.g., `ProviderKindWithCilium`, `ProviderAzureExistingCluster`) to refer to supported providers in a type-safe way. Each provider subdirectory contains workflows and steps for provisioning, configuring, and managing clusters on that platform.

## Adding Providers

To add a new provider:
1. Create a new subdirectory (e.g., `mycloud/` for a new provider).
2. Implement workflows and steps for that provider.
3. Add a new enum value to `provider-index.go` and update the `ProviderSetupSteps` map.

---

This directory is part of the infrastructure automation layer for scalable Kubernetes testing and benchmarking.
