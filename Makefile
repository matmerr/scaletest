# Makefile for scaletest

.PHONY: generate
generate:
	@echo "Generating scenario YAMLs using go generate..."
	go test -v -run TestGenerate 

.PHONY: tools
tools:
	go test -v -run TestDownloadTools

.PHONY: test-cl2
test-cl2:
	@echo "Running clusterloader2 tests..."
	go test -v -run ^TestRunClusterLoader2Scenarios$$ ./main_clusterloader2_test.go

.PHONY: test-kubeburner
test-kubeburner:
	@echo "Running kubeburner tests..."
	go test -v -run ^TestRunKubeburnerScenarios$$ ./main_kubeburner_test.go
