---
layout: documentation
title: Installation Guide
prev_page: getting-started.html
prev_title: Getting Started
next_page: chaos-types.html
next_title: Chaos Types
---

# Installation Guide

<div class="callout callout-info">
  <div class="callout-title">
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
    Installation Options
  </div>
  <p>This guide provides detailed instructions for installing Havock8s in your Kubernetes cluster using multiple methods.</p>
</div>

## Prerequisites

Before installing Havock8s, ensure your environment meets these requirements:

<div class="article-section">
  <ul>
    <li><strong>Kubernetes cluster</strong> - version 1.18 or higher</li>
    <li><strong>kubectl</strong> - configured to communicate with your cluster</li>
    <li><strong>Helm v3</strong> - (optional) for chart-based installation</li>
    <li><strong>Cluster admin access</strong> - required for creating CRDs and RBAC roles</li>
  </ul>
</div>

## Installation Methods

Havock8s can be installed using several methods. Choose the one that best fits your workflow and requirements.

### Method 1: Using kubectl

This is the simplest method to quickly install Havock8s:

```bash
kubectl apply -f https://raw.githubusercontent.com/havock8s/havock8s/main/config/install.yaml
```

When you run this command, the following resources will be created:

1. **Custom Resource Definitions (CRDs)** - Define the chaos experiments
2. **RBAC Resources** - Service accounts, roles, and role bindings
3. **Deployments** - The Havock8s controller manager
4. **Services** - For monitoring and metrics

### Method 2: Using Helm

For more flexibility, easier upgrades, and configuration management, Helm is recommended:

```bash
# Add the Havock8s Helm repository
helm repo add havock8s https://charts.havock8s.io

# Update your Helm repositories
helm repo update

# Install Havock8s
helm install havock8s havock8s/havock8s
```

#### Customizing the Helm Installation

Havock8s's Helm chart offers extensive customization options. Create a `values.yaml` file:

```yaml
controller:
  # Resource limits and requests
  resources:
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 100m
      memory: 128Mi
  
  # Controller settings
  logLevel: info
  replicaCount: 1

# Integration with monitoring systems
monitoring:
  prometheus:
    enabled: true
    serviceMonitor:
      enabled: true
  
# Distributed tracing
tracing:
  opentelemetry:
    enabled: true
    endpoint: "otel-collector:4317"

# Security settings
security:
  podSecurityContext:
    runAsNonRoot: true
    runAsUser: 1000
  containerSecurityContext:
    allowPrivilegeEscalation: false
```

Install with your custom values:

```bash
helm install havock8s havock8s/havock8s -f values.yaml
```

<div class="callout callout-tip">
  <div class="callout-title">
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path></svg>
    Helm Pro Tip
  </div>
  <p>For production deployments, always create a dedicated values file with your customizations. This makes upgrades and configuration tracking much easier.</p>
</div>

### Method 3: Building from Source

For development, customization, or contributing to Havock8s:

```bash
# Clone the repository
git clone https://github.com/havock8s/havock8s.git
cd havock8s

# Build the controller
make build

# Deploy to the cluster
make deploy
```

This method requires:
- Go 1.17+
- Docker or compatible container runtime
- Access to push images to a container registry

## Verification and Troubleshooting

### Verifying the Installation

After installation, verify that all components are running properly:

```bash
# Check controller pod status
kubectl get pods -n havock8s-system
```

You should see output similar to:

```
NAME                                        READY   STATUS    RESTARTS   AGE
havock8s-controller-manager-xxxx-yyyy       1/1     Running   0          1m
```

Verify the Custom Resource Definitions are installed:

```bash
kubectl get crds | grep havock8s
```

Expected output:

```
havock8sexperiments.chaos.havock8s.io
```

### Common Installation Issues

<div class="article-section">
  <h4>Controller Pod Not Starting</h4>
  <p>If the controller pod isn't starting, check the pod logs:</p>
  <pre><code>kubectl logs -n havock8s-system deploy/havock8s-controller-manager</code></pre>
  
  <h4>CRDs Not Being Created</h4>
  <p>Ensure you have cluster-admin privileges to create CRDs:</p>
  <pre><code>kubectl auth can-i create crds</code></pre>
  
  <h4>Webhook Configuration Issues</h4>
  <p>If you're using webhooks, check that the certificates are correctly set up:</p>
  <pre><code>kubectl get validatingwebhookconfiguration -l app=havock8s</code></pre>
</div>

## Uninstalling

### Using kubectl

```bash
kubectl delete -f https://raw.githubusercontent.com/havock8s/havock8s/main/config/install.yaml
```

### Using Helm

```bash
helm uninstall havock8s
```

### Manual Cleanup (if needed)

If some resources remain after uninstallation:

```bash
# Remove CRDs
kubectl delete crd havock8sexperiments.chaos.havock8s.io

# Remove namespace
kubectl delete namespace havock8s-system

# Remove RBAC resources
kubectl delete clusterrole havock8s-manager-role
kubectl delete clusterrolebinding havock8s-manager-rolebinding
```

## Configurations for Different Environments

<div class="docs-section">
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Development</h3>
    </div>
    <div class="docs-card-content">
      <p>For development or testing, use minimal resources:</p>
      <pre><code>helm install havock8s havock8s/havock8s \
  --set controller.resources.requests.cpu=50m \
  --set controller.resources.requests.memory=64Mi</code></pre>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Production</h3>
    </div>
    <div class="docs-card-content">
      <p>For production, enable high availability and monitoring:</p>
      <pre><code>helm install havock8s havock8s/havock8s \
  --set controller.replicaCount=2 \
  --set monitoring.prometheus.enabled=true</code></pre>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>Restricted Environments</h3>
    </div>
    <div class="docs-card-content">
      <p>For environments with strict security policies:</p>
      <pre><code>helm install havock8s havock8s/havock8s \
  --set security.podSecurityContext.runAsNonRoot=true \
  --set security.containerSecurityContext.readOnlyRootFilesystem=true</code></pre>
    </div>
  </div>
</div>

## Next Steps

Now that you have Havock8s installed, you can:

- [Get started with your first chaos experiment](getting-started.html)
- [Learn about different chaos types](chaos-types.html)
- [Explore tutorials for common scenarios](tutorials.html) 