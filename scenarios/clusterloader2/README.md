# clusterloader2 Scenarios Directory

This directory contains scenarios and supporting files for running clusterloader2-based tests.

## Directory Layout

- `scenarios/` — Scenario definitions and configurations for different test cases.
- `modules/` — Go modules and resources used as clusterloader2 modules (e.g., embedded YAMLs, helpers).


## Notes

- The `modules` directory is specifically for clusterloader2 modules, which may be used for measurements
- Docs for such modules are here: https://pkg.go.dev/k8s.io/perf-tests/clusterloader2#section-readme

## clusterloader2 examples:
- https://github.com/kubernetes/perf-tests/tree/master/clusterloader2/testing/load
