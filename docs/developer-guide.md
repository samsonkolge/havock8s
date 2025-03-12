# Developer Guide

This guide provides information for developers who want to contribute to or extend havock8s.

## Architecture Overview

havock8s follows the Kubernetes Operator pattern and is built using the Operator SDK and controller-runtime libraries. The main components are:

1. **Custom Resource Definitions (CRDs)**: Define the havock8sExperiment resource
2. **Controller**: Watches for havock8sExperiment resources and reconciles their state
3. **Chaos Injectors**: Implement different types of chaos (disk failures, network latency, etc.)
4. **Safety Mechanisms**: Ensure experiments don't cause cascading failures

## Directory Structure

```
havock8s/
├── api/                  # API definitions for CRDs
│   └── v1alpha1/         # API version
├── config/               # Kubernetes manifests
├── controllers/          # Controller implementation
├── docs/                 # Documentation
├── examples/             # Example chaos experiments
├── pkg/                  # Shared packages
│   ├── chaos/            # Chaos injector implementations
│   └── utils/            # Utility functions
└── tests/                # Integration and end-to-end tests
```

## Development Environment Setup

### Prerequisites

- Go 1.23+
- Kubernetes cluster (v1.18+)
- kubectl
- kustomize
- controller-gen

### Setting Up Your Development Environment

1. Clone the repository:

```bash
git clone https://github.com/havock8s/havock8s.git
cd havock8s
```

2. Install dependencies:

```bash
go mod download
```

3. Generate code and manifests:

```bash
make generate
make manifests
```

4. Build and run the controller locally:

```bash
make install  # Install CRDs in the cluster
make run      # Run controller locally
```

## Implementing a New Chaos Type

To add a new chaos type:

1. Define the chaos type in `api/v1alpha1/havock8sexperiment_types.go`
2. Create a new injector in `pkg/chaos/`
3. Register the injector in `pkg/chaos/injector.go`
4. Update the controller to handle the new chaos type
5. Add tests for the new chaos type
6. Add documentation and examples

Example of a chaos injector:

```go
package chaos

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewChaosInjector creates a new chaos injector for the given experiment
func NewMyChaosInjector(client client.Client, log logr.Logger) Injector {
	return &myChaosInjector{
		client: client,
		log:    log,
	}
}

type myChaosInjector struct {
	client client.Client
	log    logr.Logger
}

func (i *myChaosInjector) Inject(ctx context.Context, experiment *chaosv1alpha1.havock8sExperiment) error {
	i.log.Info("Injecting chaos", "type", experiment.Spec.ChaosType)
	
	// Implement chaos injection logic here
	
	return nil
}

func (i *myChaosInjector) Cleanup(ctx context.Context, experiment *chaosv1alpha1.havock8sExperiment) error {
	i.log.Info("Cleaning up chaos", "type", experiment.Spec.ChaosType)
	
	// Implement cleanup logic here
	
	return nil
}
```

## Testing

### Running Unit Tests

```bash
make test
```

### Running Integration Tests

```bash
make test-integration
```

### Running End-to-End Tests

```bash
make test-e2e
```

## Building and Deploying

### Building the Controller Image

```bash
make docker-build IMG=your-registry/havock8s:tag
```

### Pushing the Controller Image

```bash
make docker-push IMG=your-registry/havock8s:tag
```

### Deploying to a Cluster

```bash
make deploy IMG=your-registry/havock8s:tag
```

## Debugging

### Enabling Debug Logs

Set the `--zap-log-level=debug` flag when running the controller.

### Using Delve for Debugging

```bash
dlv debug ./main.go -- --zap-log-level=debug
```

## Release Process

1. Update version in `VERSION` file
2. Update CHANGELOG.md
3. Create a new tag:

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

4. The CI/CD pipeline will build and publish the release artifacts

## Documentation

Please update the documentation when making changes:

- Update API reference when changing CRDs
- Add examples for new features
- Update the developer guide for significant changes

## Getting Help

If you need help or have questions, please:

- Open an issue on GitHub
- Join our community Slack channel
- Attend our community meetings 