# StatefulChaos

StatefulChaos is a cloud-native chaos engineering framework specifically designed for stateful applications running on Kubernetes. While many existing chaos engineering tools focus on stateless microservices, StatefulChaos specializes in testing and enhancing the resilience of stateful components such as databases, caching systems, and persistent storage.

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
kubectl apply -f https://raw.githubusercontent.com/statefulchaos/statefulchaos/main/config/install.yaml
```

#### Using Helm

```bash
helm repo add statefulchaos https://statefulchaos.github.io/charts
helm install statefulchaos statefulchaos/statefulchaos
```

## Quick Start

To run a simple chaos experiment against a stateful application:

```yaml
apiVersion: chaos.statefulchaos.io/v1alpha1
kind: StatefulChaosExperiment
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

For detailed documentation, examples, and guides, visit [the official documentation](https://statefulchaos.io/docs).

## Contributing

Contributions are welcome! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

StatefulChaos is open-source software licensed under the Apache License 2.0. 