---
layout: default
title: Havock8s Documentation
---

# Havock8s

![Havock8s Logo](assets/images/logo.png)

Havock8s is a cloud-native chaos engineering framework specifically designed for stateful applications running on Kubernetes. While many existing chaos engineering tools focus on stateless microservices, Havock8s specializes in testing and enhancing the resilience of stateful components such as databases, caching systems, and persistent storage.

## Why Havock8s?

Traditional chaos engineering tools often focus on stateless services where instances can be easily replaced. However, stateful applications present unique challenges:

- Data persistence requirements
- Complex recovery procedures
- State synchronization across replicas
- Storage dependencies

Havock8s addresses these challenges with targeted chaos experiments designed specifically for stateful workloads.

## Key Features

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

## Documentation

<div class="docs-section">
  <div class="docs-card">
    <h3><a href="getting-started.html">Getting Started</a></h3>
    <p>Quick introduction to Havock8s with setup instructions</p>
  </div>
  <div class="docs-card">
    <h3><a href="installation.html">Installation Guide</a></h3>
    <p>Detailed installation instructions for different environments</p>
  </div>
  <div class="docs-card">
    <h3><a href="chaos-types.html">Chaos Types</a></h3>
    <p>Overview of available chaos scenarios and their parameters</p>
  </div>
  <div class="docs-card">
    <h3><a href="api-reference.html">API Reference</a></h3>
    <p>Complete reference of all API objects and their fields</p>
  </div>
  <div class="docs-card">
    <h3><a href="developer-guide.html">Developer Guide</a></h3>
    <p>Instructions for developers who want to contribute or extend</p>
  </div>
  <div class="docs-card">
    <h3><a href="tutorials.html">Tutorials</a></h3>
    <p>Step-by-step guides for common use cases</p>
  </div>
</div>

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

## Community

- [GitHub Issues](https://github.com/havock8s/havock8s/issues)
- [Slack Channel](#)
- [Twitter](https://twitter.com/havock8s)

## License

Havock8s is open-source software licensed under the Apache License 2.0. 