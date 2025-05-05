# havock8s

```
 ██╗  ██╗  █████╗  ██╗   ██╗  ██████╗   ██████╗ ██╗  ██╗  █████╗  ███████╗
 ██║  ██║ ██╔══██╗ ██║   ██║ ██╔═══██╗ ██╔════╝ ██║ ██╔╝ ██╔══██╗ ██╔════╝
 ███████║ ███████║ ██║   ██║ ██║   ██║ ██║      █████╔╝  ╚█████╔╝ ███████╗
 ██╔══██║ ██╔══██║ ╚██╗ ██╔╝ ██║   ██║ ██║      ██╔═██╗  ██╔══██╗ ╚════██║
 ██║  ██║ ██║  ██║  ╚████╔╝  ╚██████╔╝ ╚██████╗ ██║  ██╗ ╚█████╔╝ ███████║
 ╚═╝  ╚═╝ ╚═╝  ╚═╝   ╚═══╝    ╚═════╝   ╚═════╝ ╚═╝  ╚═╝  ╚════╝  ╚══════╝
```

[![Build Status](https://github.com/havock8s/havock8s/workflows/CI/badge.svg)](https://github.com/havock8s/havock8s/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/havock8s/havock8s)](https://goreportcard.com/report/github.com/havock8s/havock8s)
[![GoDoc](https://pkg.go.dev/badge/github.com/havock8s/havock8s)](https://pkg.go.dev/github.com/havock8s/havock8s)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/havock8s/havock8s)](go.mod)

havock8s is a cloud-native chaos engineering framework specifically designed for stateful applications running on Kubernetes. While many existing chaos engineering tools focus on stateless microservices, havock8s specializes in testing and enhancing the resilience of stateful components such as databases, caching systems, and persistent storage.

## Features

### Customizable Failure Scenarios
- Pre-defined chaos scenarios for stateful workloads (disk I/O failures, database connection issues, persistent volume disruptions)
- Custom experiment framework for defining specialized failure modes
- Targeted chaos injection for specific components of stateful applications

### Kubernetes-Native Integration
- Operates as a Kubernetes Operator using Custom Resource Definitions (CRDs)
- Seamless integration with Prometheus for monitoring and OpenTelemetry for tracing
- Works with standard Kubernetes resources and StatefulSets

### Safety and Rollback Mechanisms
- Built-in guardrails to prevent cascading failures
- Automated rollback features for safe production use
- Gradual chaos intensity adjustments

### Extensible Plugin Architecture
- Community-driven plugin system for adding new failure modes
- Extensible interfaces for custom integrations
- Open-source foundation for collaboration

### Visualizations and Reporting
- Integration with Grafana for experiment visualizations
- Automated experiment reports
- Performance impact analysis dashboards

## Getting Started

### Prerequisites
- Kubernetes cluster (v1.18+)
- kubectl configured to communicate with your cluster
- Helm v3 (optional, for chart-based installation)

### Installation

#### Using kubectl

```bash
kubectl apply -f https://raw.githubusercontent.com/havock8s/havock8s/main/config/install.yaml
```

#### Using Helm

```bash
helm repo add havock8s https://charts.havock8s.io
helm install havock8s havock8s/havock8s
```

## Quick Start

To run a simple chaos experiment against a stateful application:

```yaml
apiVersion: chaos.havock8s.io/v1alpha1
kind: havock8sExperiment
metadata:
  name: postgres-disk-failure
spec:
  target:
    selector:
      app: postgres
  chaosType: DiskFailure
  duration: 5m
  intensity: 0.3  # 30% of I/O operations will fail
  safety:
    autoRollback: true
    healthChecks:
      - type: httpGet
        path: /health
        port: 8080
        failureThreshold: 3
```

## Documentation

For detailed documentation, examples, and guides, visit [the official documentation](https://docs.havock8s.io).

## Contributing

Contributions are welcome! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

havock8s is open-source software licensed under the Apache License 2.0.

## Chaos Experiments

### Pod Failure Chaos
The pod failure chaos experiment successfully terminated a target pod in the test namespace. This experiment type is designed to test application resilience by simulating pod failures in a controlled manner.

#### Example Usage

1. Create a test pod:
```bash
kubectl run test-pod --image=nginx -n test
```

2. Apply the chaos experiment:
```yaml
apiVersion: chaos.havock8s.io/v1alpha1
kind: Havock8sExperiment
metadata:
  name: test-pod-failure
  namespace: test
spec:
  chaosType: PodFailure
  duration: 5m
  intensity: 1
  safety:
    autoRollback: true
  target:
    name: test-pod
    namespace: test
    targetType: Pod
```

#### Experiment Behavior
- **Target Selection**: The experiment identifies and targets the specified pod by name and namespace
- **Termination**: Pods are terminated immediately (gracePeriodSeconds: 0) or with a specified grace period
- **Duration**: The experiment maintains the "Running" state for the specified duration
- **Auto Rollback**: If configured (autoRollback: true), the experiment will automatically restore the pod after completion

#### Monitoring
You can monitor the experiment progress using:
```bash
# Check experiment status
kubectl get havock8sexperiment test-pod-failure -n test -o yaml

# View controller logs
kubectl logs -n havock8s-system deployment/havock8s-controller-manager
```

#### Expected Outcomes
- The target pod will be terminated
- The experiment will show "Running" status
- Target resources will be marked as "Targeted" in the experiment status
- Controller logs will show successful execution of the chaos injection

#### Troubleshooting
If the experiment doesn't work as expected:
1. Verify the target pod exists and is running
2. Check controller logs for any errors
3. Ensure the experiment has the correct target name and namespace
4. Verify the controller has necessary permissions to delete pods

The controller logs showed successful execution of the chaos injection, with proper tracking of the target pod and experiment status. 