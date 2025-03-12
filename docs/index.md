---
layout: default
title: Havock8s Documentation
---

# Havock8s

![Havock8s Logo](assets/images/logo.png)

Havock8s is a cloud-native chaos engineering framework specifically designed for stateful applications running on Kubernetes. While many existing chaos engineering tools focus on stateless microservices, Havock8s specializes in testing and enhancing the resilience of stateful components such as databases, caching systems, and persistent storage.

<div class="callout callout-info">
  <div class="callout-title">
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
    Getting Started
  </div>
  <p>Ready to start testing your stateful applications? Check out our <a href="getting-started.html">Getting Started Guide</a> to begin your chaos engineering journey.</p>
</div>

## Why Havock8s?

Traditional chaos engineering tools often focus on stateless services where instances can be easily replaced. However, stateful applications present unique challenges:

- **Data persistence requirements**
- **Complex recovery procedures**
- **State synchronization across replicas**
- **Storage dependencies**

Havock8s addresses these challenges with targeted chaos experiments designed specifically for stateful workloads.

## Key Features

<div class="docs-section">
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Customizable Failure Scenarios</h3>
    </div>
    <div class="docs-card-content">
      <p>Pre-defined chaos scenarios for stateful workloads including disk I/O failures, database connection issues, and persistent volume disruptions.</p>
      <a href="chaos-types.html">Explore Chaos Types →</a>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Kubernetes-Native Integration</h3>
    </div>
    <div class="docs-card-content">
      <p>Operates as a Kubernetes Operator using Custom Resource Definitions (CRDs) with seamless integration with Prometheus and OpenTelemetry.</p>
      <a href="installation.html">Installation Guide →</a>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Safety and Rollback Mechanisms</h3>
    </div>
    <div class="docs-card-content">
      <p>Built-in guardrails to prevent cascading failures with automated rollback features for safe production use.</p>
      <a href="api-reference.html#safety-mechanisms">Safety Mechanisms →</a>
    </div>
  </div>
</div>

<div class="docs-section">
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Extensible Plugin Architecture</h3>
    </div>
    <div class="docs-card-content">
      <p>Community-driven plugin system for adding new failure modes with extensible interfaces for custom integrations.</p>
      <a href="developer-guide.html">Developer Guide →</a>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Visualizations and Reporting</h3>
    </div>
    <div class="docs-card-content">
      <p>Integration with Grafana for experiment visualizations, automated experiment reports, and performance impact analysis dashboards.</p>
      <a href="tutorials.html">Tutorials →</a>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Production-Ready</h3>
    </div>
    <div class="docs-card-content">
      <p>Designed for real-world chaos engineering with comprehensive safety mechanisms and gradual chaos introduction features.</p>
      <a href="https://github.com/samsonkolge/havock8s/tree/master/examples">Examples →</a>
    </div>
  </div>
</div>

## Documentation

Our documentation helps you get started with Havock8s quickly and efficiently:

<div class="docs-section">
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Getting Started</h3>
    </div>
    <div class="docs-card-content">
      <p>Quick introduction to Havock8s with setup instructions</p>
      <a href="getting-started.html">Read guide →</a>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Installation Guide</h3>
    </div>
    <div class="docs-card-content">
      <p>Detailed installation instructions for different environments</p>
      <a href="installation.html">Install now →</a>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Tutorials</h3>
    </div>
    <div class="docs-card-content">
      <p>Step-by-step guides for common use cases</p>
      <a href="tutorials.html">Follow tutorials →</a>
    </div>
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

<div class="callout callout-tip">
  <div class="callout-title">
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path></svg>
    Pro Tip
  </div>
  <p>Always start with a low intensity value and gradually increase it to avoid unexpected impacts on your systems. Check out our <a href="tutorials.html">tutorials</a> for more best practices.</p>
</div>

## Community

- [GitHub Issues](https://github.com/havock8s/havock8s/issues)
- [Slack Channel](#)
- [Twitter](https://twitter.com/havock8s)

## License

Havock8s is open-source software licensed under the Apache License 2.0. 