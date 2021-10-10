# Set the shell
SHELL := /bin/bash

# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)

# Set our default go compiler
GO := go

.PHONY: all
all: hello

.PHONY: fmt
fmt: ## Verifies all files have been `gofmt`ed.
	@echo "+ $@"
	@if [[ ! -z "$(shell gofmt -s -l . | grep -v '.pb.go:' | grep -v '.twirp.go:' | grep -v vendor | tee /dev/stderr)" ]]; then \
		exit 1; \
	fi

.PHONY: vet
vet: ## Verifies `go vet` passes.
	@echo "+ $@"
	@if [[ ! -z "$(shell $(GO) vet $(shell $(GO) list ./... | grep -v vendor) | tee /dev/stderr)" ]]; then \
		exit 1; \
	fi

# This rule tells make how to build hello from hello.cpp
hello: src/hello.c
	gcc -o hello src/hello.c

.PHONY: install
install:
	mkdir -p bin
	cp -p hello bin

.PHONY: clean 
clean:
	rm -f hello
