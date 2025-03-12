# Getting Started with havock8s

This guide will help you get started with havock8s, a cloud-native chaos engineering framework for stateful applications running on Kubernetes.

## Prerequisites

Before you begin, make sure you have the following:

- A Kubernetes cluster (v1.18+)
- kubectl configured to communicate with your cluster
- Helm v3 (optional, for chart-based installation)

## Installation

### Using kubectl

```bash
kubectl apply -f https://raw.githubusercontent.com/havock8s/havock8s/main/config/install.yaml
```

### Using Helm

```bash
helm repo add havock8s https://havock8s.github.io/charts
helm install havock8s havock8s/havock8s
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

## Your First Chaos Experiment

Create a file named `first-experiment.yaml` with the following content:

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

Apply the experiment:

```bash
kubectl apply -f first-experiment.yaml
```

Monitor the experiment:

```bash
kubectl get havock8sexperiments
```

## Next Steps

- Learn about [different chaos types](chaos-types.md)
- Explore [advanced configuration options](advanced-configuration.md)
- Set up [monitoring and observability](monitoring.md)
- Check out [example scenarios](../examples/README.md) 