# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=havock8s
BINARY_UNIX=$(BINARY_NAME)_unix

# Linter parameters
LINTER=$(shell go env GOPATH)/bin/golangci-lint
LINTER_VERSION=v1.55.2

# Build parameters
BUILD_DIR=build
DIST_DIR=dist

# GitOps parameters
KUSTOMIZE=$(shell go env GOPATH)/bin/kustomize

.PHONY: all build clean test lint install-linter gitops-apply gitops-diff gitops-destroy

all: clean lint test build

build:
	@echo "Building binary..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

lint:
	@echo "Running linter..."
	$(LINTER) run

install-linter:
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(LINTER_VERSION)

# Development tools
.PHONY: tools
tools: install-linter

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

# GitOps targets
.PHONY: gitops-apply gitops-diff gitops-destroy

gitops-apply: manifests kustomize
	$(KUSTOMIZE) build config/default | kubectl apply -f -

gitops-diff: manifests kustomize
	$(KUSTOMIZE) build config/default | kubectl diff -f -

gitops-destroy: manifests kustomize
	$(KUSTOMIZE) build config/default | kubectl delete -f -

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all           - Clean, lint, test, and build"
	@echo "  build         - Build the binary"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  lint          - Run linter"
	@echo "  install-linter - Install golangci-lint"
	@echo "  tools         - Install development tools"
	@echo "  docker-build  - Build Docker image"
	@echo "  gitops-apply  - Apply GitOps manifests"
	@echo "  gitops-diff   - Diff GitOps manifests"
	@echo "  gitops-destroy - Destroy GitOps manifests"
	@echo "  help          - Show this help message" 