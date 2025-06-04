# Makefile for scaletest

.PHONY: generate
generate:
	@echo "Generating scenario YAMLs using go generate..."
	go test -run TestGenerate

.PHONY: tools
tools: 
	go test -run TestDownloadTools
