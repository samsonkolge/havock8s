---
layout: documentation
title: API Reference
prev_page: chaos-types.html
prev_title: Chaos Types
next_page: tutorials.html
next_title: Tutorials
---

# Havock8s API Reference

<div class="callout callout-info">
  <div class="callout-title">
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
    API Documentation
  </div>
  <p>This reference documents the Kubernetes Custom Resource Definitions (CRDs) that make up the Havock8s API.</p>
</div>

## Overview

Havock8s extends Kubernetes with custom resources to define and manage chaos experiments. The primary API resource is the `havock8sExperiment`, which defines a chaos experiment against stateful applications.

## API Versions

| API Version | Description | Status |
|-------------|-------------|--------|
| `chaos.havock8s.io/v1alpha1` | Initial API version | Available |
| `chaos.havock8s.io/v1beta1` | Beta API with additional features | In development |

## havock8sExperiment

The `havock8sExperiment` resource defines a chaos experiment to be executed against target applications.

### YAML Structure

```yaml
apiVersion: chaos.havock8s.io/v1alpha1
kind: havock8sExperiment
metadata:
  name: example-experiment
spec:
  # Target selection - required
  target:
    selector:
      app: example-app
    namespaces:
      - default
    
  # Chaos type - required
  chaosType: DiskFailure
  
  # Experiment schedule - required
  duration: 10m
  startTime: "2023-06-15T08:00:00Z"  # Optional, if not specified starts immediately
  timeoutSeconds: 600  # Optional, defaults to duration
  
  # Chaos parameters - varies by chaosType
  intensity: 0.3
  mode: error
  
  # Safety settings - recommended
  safety:
    autoRollback: true
    healthChecks:
      - type: httpGet
        path: /health
        port: 8080
        periodSeconds: 10
        failureThreshold: 3
    maxTargetPods: 1
    targetPercentage: 30
```

### Spec Fields

<div class="article-section">
  <table>
    <thead>
      <tr>
        <th>Field</th>
        <th>Type</th>
        <th>Description</th>
        <th>Required</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td><code>target</code></td>
        <td>Object</td>
        <td>Defines which resources will be targeted by the experiment</td>
        <td>Yes</td>
      </tr>
      <tr>
        <td><code>chaosType</code></td>
        <td>String</td>
        <td>The type of chaos to be injected (e.g., DiskFailure, NetworkLatency)</td>
        <td>Yes</td>
      </tr>
      <tr>
        <td><code>duration</code></td>
        <td>String</td>
        <td>How long the chaos experiment will run (e.g., "5m", "1h")</td>
        <td>Yes</td>
      </tr>
      <tr>
        <td><code>startTime</code></td>
        <td>String</td>
        <td>When to start the experiment (RFC3339 format)</td>
        <td>No</td>
      </tr>
      <tr>
        <td><code>timeoutSeconds</code></td>
        <td>Integer</td>
        <td>Maximum time the experiment can run before forced termination</td>
        <td>No</td>
      </tr>
      <tr>
        <td><code>safety</code></td>
        <td>Object</td>
        <td>Safety guardrails for the experiment</td>
        <td>No</td>
      </tr>
      <tr>
        <td colspan="4">Additional fields depend on the <code>chaosType</code></td>
      </tr>
    </tbody>
  </table>
</div>

### Target Specification

The `target` field defines which resources will be affected by the chaos experiment:

```yaml
target:
  # Select resources by label
  selector:
    app: postgres
    tier: database
  
  # Limit to specific namespaces (optional)
  namespaces:
    - default
    - database
  
  # Exclude specific resources by label (optional)
  excludeSelector:
    role: primary
```

<div class="article-section">
  <table>
    <thead>
      <tr>
        <th>Field</th>
        <th>Type</th>
        <th>Description</th>
        <th>Required</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td><code>selector</code></td>
        <td>Object</td>
        <td>Label selector to identify target resources</td>
        <td>Yes</td>
      </tr>
      <tr>
        <td><code>namespaces</code></td>
        <td>Array</td>
        <td>List of namespaces to limit targets to</td>
        <td>No</td>
      </tr>
      <tr>
        <td><code>excludeSelector</code></td>
        <td>Object</td>
        <td>Label selector to exclude resources</td>
        <td>No</td>
      </tr>
    </tbody>
  </table>
</div>

### Safety Mechanisms

<div id="safety-mechanisms"></div>

The `safety` field defines guardrails to prevent the chaos experiment from causing excessive damage:

```yaml
safety:
  # Automatically roll back if health checks fail
  autoRollback: true
  
  # Health checks to monitor during experiment
  healthChecks:
    - type: httpGet
      path: /health
      port: 8080
      periodSeconds: 10
      initialDelaySeconds: 20
      failureThreshold: 3
  
  # Limit the number of affected pods
  maxTargetPods: 1
  
  # Only affect a percentage of eligible pods
  targetPercentage: 30
```

<div class="article-section">
  <table>
    <thead>
      <tr>
        <th>Field</th>
        <th>Type</th>
        <th>Description</th>
        <th>Default</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td><code>autoRollback</code></td>
        <td>Boolean</td>
        <td>Whether to automatically stop experiment if health checks fail</td>
        <td><code>false</code></td>
      </tr>
      <tr>
        <td><code>healthChecks</code></td>
        <td>Array</td>
        <td>List of health checks to perform during the experiment</td>
        <td><code>[]</code></td>
      </tr>
      <tr>
        <td><code>maxTargetPods</code></td>
        <td>Integer</td>
        <td>Maximum number of pods that can be affected</td>
        <td>No limit</td>
      </tr>
      <tr>
        <td><code>targetPercentage</code></td>
        <td>Integer</td>
        <td>Percentage of eligible pods to target (1-100)</td>
        <td><code>100</code></td>
      </tr>
    </tbody>
  </table>
</div>

<div class="callout callout-warning">
  <div class="callout-title">
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>
    Important Safety Note
  </div>
  <p>Always configure appropriate safety mechanisms when running chaos experiments against production or production-like environments. Without proper safety settings, chaos experiments can cause unintended service disruptions.</p>
</div>

## Status

The `status` field is updated by the Havock8s controller to reflect the current state of the experiment:

```yaml
status:
  phase: Running
  startTime: "2023-06-15T08:00:00Z"
  endTime: null
  targetPods:
    - name: postgres-0
      namespace: default
      status: Affected
    - name: postgres-1
      namespace: default
      status: Affected
  observations:
    - time: "2023-06-15T08:02:14Z"
      message: "Disk write latency increased to 250ms"
  healthStatus:
    healthy: true
    message: "All health checks passing"
```

<div class="article-section">
  <table>
    <thead>
      <tr>
        <th>Field</th>
        <th>Type</th>
        <th>Description</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td><code>phase</code></td>
        <td>String</td>
        <td>Current phase of the experiment (Pending, Running, Completed, Failed, Stopped)</td>
      </tr>
      <tr>
        <td><code>startTime</code></td>
        <td>String</td>
        <td>When the experiment actually started</td>
      </tr>
      <tr>
        <td><code>endTime</code></td>
        <td>String</td>
        <td>When the experiment ended (null if still running)</td>
      </tr>
      <tr>
        <td><code>targetPods</code></td>
        <td>Array</td>
        <td>List of pods affected by the experiment</td>
      </tr>
      <tr>
        <td><code>observations</code></td>
        <td>Array</td>
        <td>Observations recorded during the experiment</td>
      </tr>
      <tr>
        <td><code>healthStatus</code></td>
        <td>Object</td>
        <td>Current health status of the targets</td>
      </tr>
    </tbody>
  </table>
</div>

## Common API Operations

### Creating an Experiment

```bash
kubectl apply -f experiment.yaml
```

### Viewing Experiment Status

```bash
kubectl get havock8sexperiment my-experiment -o yaml
```

### Monitoring Experiment Progress

```bash
kubectl describe havock8sexperiment my-experiment
```

### Stopping an Experiment Early

```bash
kubectl annotate havock8sexperiment my-experiment chaos.havock8s.io/stop=true
```

## Working with the API Programmatically

If you want to use the Havock8s API programmatically, you can use the Kubernetes client libraries:

<div class="docs-section">
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Go Client</h3>
    </div>
    <div class="docs-card-content">
      <p>Use the generated client in Go:</p>
      <pre><code>import (
  chaosv1alpha1 "github.com/samsonkolge/havock8s/api/v1alpha1"
  "k8s.io/client-go/kubernetes/scheme"
)

// Create a new experiment
experiment := &chaosv1alpha1.havock8sExperiment{
  // ...
}</code></pre>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Python Client</h3>
    </div>
    <div class="docs-card-content">
      <p>Use the Kubernetes Python client:</p>
      <pre><code>from kubernetes import client, config

config.load_kube_config()
api = client.CustomObjectsApi()

# Create a new experiment
experiment = {
  "apiVersion": "chaos.havock8s.io/v1alpha1",
  "kind": "havock8sExperiment",
  # ...
}

api.create_namespaced_custom_object(
  group="chaos.havock8s.io",
  version="v1alpha1",
  namespace="default",
  plural="havock8sexperiments",
  body=experiment
)</code></pre>
    </div>
  </div>
</div>

## API Extensions

The Havock8s API is extensible through custom resource definitions. For information on extending the API with custom chaos types, see the [Developer Guide](developer-guide.html). 