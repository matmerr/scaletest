# Makefile for scaletest

.PHONY: generate
generate:
	@echo "Generating scenario YAMLs using go generate..."
	go test ./tools -run TestGenerate

.PHONY: tools
tools: 
	go test ./tools -run TestDownloadTools
