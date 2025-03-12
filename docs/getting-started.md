---
layout: documentation
title: Getting Started with Havock8s
next_page: installation.html
next_title: Installation Guide
---

# Getting Started with Havock8s

<div class="callout callout-info">
  <div class="callout-title">
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
    Quick Start
  </div>
  <p>This guide will help you quickly get started with Havock8s, a cloud-native chaos engineering framework for stateful applications running on Kubernetes.</p>
</div>

## What is Havock8s?

Havock8s is specifically designed to help you test the resilience of stateful applications in Kubernetes environments. While many chaos engineering tools focus on stateless services, Havock8s specializes in testing applications that maintain state, such as:

- Databases (PostgreSQL, MySQL, MongoDB, etc.)
- Message queues and streaming platforms (Kafka, RabbitMQ)
- Caching systems (Redis, Memcached)
- Stateful applications with persistent storage requirements

## Why Chaos Engineering for Stateful Applications?

Stateful applications require special attention during chaos testing because:

1. **Data persistence is critical** - Data loss can be catastrophic
2. **Recovery procedures are complex** - Restoring state after failure can involve multiple steps
3. **State synchronization** - Ensuring replicas maintain consistent state after disruptions
4. **Storage dependencies** - Interaction with storage systems adds another failure point

Havock8s provides targeted chaos experiments that safely test these specific challenges.

## Prerequisites

Before you begin, ensure you have:

<div class="article-section">
  <ul>
    <li>A Kubernetes cluster (v1.18+)</li>
    <li>kubectl configured to communicate with your cluster</li>
    <li>Helm v3 (optional, for chart-based installation)</li>
    <li>A test environment that resembles your production environment</li>
  </ul>
</div>

## Quick Installation

You can install Havock8s using either kubectl or Helm:

### Using kubectl

```bash
kubectl apply -f https://raw.githubusercontent.com/havock8s/havock8s/main/config/install.yaml
```

### Using Helm

```bash
helm repo add havock8s https://havock8s.github.io/charts
helm install havock8s havock8s/havock8s
```

For more detailed installation instructions, including configuration options and advanced setups, visit our [Installation Guide](installation.html).

## Verifying the Installation

After installation, verify that the Havock8s controller is running properly:

```bash
kubectl get pods -n havock8s-system
```

You should see output similar to:

```
NAME                                        READY   STATUS    RESTARTS   AGE
havock8s-controller-manager-xxxx-yyyy       1/1     Running   0          1m
```

## Your First Chaos Experiment

Let's create a simple chaos experiment that simulates disk I/O failures for a PostgreSQL database:

<div class="callout callout-warning">
  <div class="callout-title">
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
    Important
  </div>
  <p>Always run chaos experiments in a test environment first before trying them in production.</p>
</div>

### 1. Create the Experiment Definition

Create a file named `postgres-disk-failure.yaml` with the following content:

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

### 2. Apply the Experiment

```bash
kubectl apply -f postgres-disk-failure.yaml
```

### 3. Monitor the Experiment

Track the status of your experiment:

```bash
kubectl get havock8sexperiments
```

### 4. View Detailed Results

Get detailed information about your experiment:

```bash
kubectl describe havock8sexperiment postgres-disk-failure
```

## Understanding the Results

When your chaos experiment completes, you'll see results that help you understand how your application behaved during the failure scenario:

- **Recovery Time**: How long it took your system to recover after the chaos was stopped
- **Error Rate**: Any increase in error rates during the chaos period
- **Performance Impact**: Changes in latency or throughput
- **Failure Modes**: Specific ways your application failed or degraded during chaos

These metrics help you identify resilience gaps in your stateful applications.

## Next Steps

Now that you've run your first experiment, explore these resources to deepen your Havock8s expertise:

<div class="docs-section">
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Explore Chaos Types</h3>
    </div>
    <div class="docs-card-content">
      <p>Discover the different types of chaos experiments available in Havock8s</p>
      <a href="chaos-types.html">View Chaos Types →</a>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Follow Tutorials</h3>
    </div>
    <div class="docs-card-content">
      <p>Step-by-step guides for common chaos testing scenarios</p>
      <a href="tutorials.html">Explore Tutorials →</a>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>API Reference</h3>
    </div>
    <div class="docs-card-content">
      <p>Detailed documentation of the Havock8s API</p>
      <a href="api-reference.html">View API Reference →</a>
    </div>
  </div>
</div> 