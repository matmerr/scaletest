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
	go test -count=1 -v -run ^TestRunClusterLoader2Scenarios$$ .

.PHONY: test-kb
test-kb:
	@echo "Running kubeburner tests..."
	go test -count=1 -v -run ^TestRunKubeBurnerScenarios$$ .
