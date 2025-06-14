# Makefile for scaletest

.PHONY: generate
generate:
	@echo "Generating config files..."
	go test -v -run TestGenerate 

.PHONY: tools
tools:
	go test -v -run TestDownloadTools

.PHONY: test-cl2
test-cl2:
	@echo "Running clusterloader2 tests..."
	go test -v -run ^TestRunClusterLoader2Scenarios$$ .

.PHONY: test-kubeburner
test-kubeburner:
	@echo "Running kubeburner tests..."
	go test -v -run ^TestRunKubeBurnerScenarios$$ .
