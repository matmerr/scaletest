# Makefile for scaletest

.PHONY: generate
generate:
	@echo "Generating scenario YAMLs using go generate..."
	go test -v -run TestGenerate 

.PHONY: tools
tools:
	go test -v -run TestDownloadTools

.PHONY: test-cl2
# Run only the TestRunClusterLoader2Scenarios test from main_test.go
# Usage: make test-cl2

test-cl2:
	go test -v -run ^TestRunClusterLoader2Scenarios$$ ./main_test.go
