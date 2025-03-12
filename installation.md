# Installation Guide

This guide provides detailed instructions for installing havock8s in your Kubernetes cluster.

## Prerequisites

- Kubernetes cluster (v1.18+)
- kubectl configured to communicate with your cluster
- Helm v3 (optional, for chart-based installation)

## Installation Methods

### Method 1: Using kubectl

This is the simplest method to install havock8s:

```bash
kubectl apply -f https://raw.githubusercontent.com/havock8s/havock8s/main/config/install.yaml
```

This will:
1. Create the necessary Custom Resource Definitions (CRDs)
2. Set up RBAC permissions
3. Deploy the havock8s controller

### Method 2: Using Helm

For more flexibility and easier upgrades, you can use Helm:

```bash
# Add the havock8s Helm repository
helm repo add havock8s https://havock8s.github.io/charts

# Update your Helm repositories
helm repo update

# Install havock8s
helm install havock8s havock8s/havock8s
```

#### Customizing the Helm Installation

You can customize the installation by providing a values file:

```bash
helm install havock8s havock8s/havock8s -f my-values.yaml
```

Example `my-values.yaml`:

```yaml
controller:
  resources:
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 100m
      memory: 128Mi

monitoring:
  prometheus:
    enabled: true
  
tracing:
  opentelemetry:
    enabled: true
```

### Method 3: Building from Source

For development or customization:

```bash
# Clone the repository
git clone https://github.com/havock8s/havock8s.git
cd havock8s

# Build and deploy
make deploy
```

## Verifying the Installation

Check that the havock8s controller is running:

```bash
kubectl get pods -n havock8s-system
```

You should see the havock8s controller pod running:

```
NAME                                        READY   STATUS    RESTARTS   AGE
havock8s-controller-manager-xxxx-yyyy       1/1     Running   0          1m
```

Verify the CRDs are installed:

```bash
kubectl get crds | grep havock8s
```

You should see:

```
havock8sexperiments.chaos.havock8s.io
```

## Uninstalling

### Using kubectl

```bash
kubectl delete -f https://raw.githubusercontent.com/havock8s/havock8s/main/config/install.yaml
```

### Using Helm

```bash
helm uninstall havock8s
```

## Next Steps

- [Getting Started Guide](getting-started.md)
- [Configuration Options](configuration.md)
- [Example Scenarios](../examples/README.md) 