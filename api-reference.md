# API Reference

This document provides a detailed reference for the havock8s API.

## Custom Resource Definitions

havock8s defines the following Custom Resource Definitions (CRDs):

### havock8sExperiment

The main CRD for defining chaos experiments.

```yaml
apiVersion: chaos.havock8s.io/v1alpha1
kind: havock8sExperiment
metadata:
  name: string
  namespace: string
spec:
  # Target defines what to target with chaos
  target:
    # Select resources by label
    selector:
      matchLabels:
        key: value
    # Or specify a specific resource
    name: string
    namespace: string
    # Type of resource to target
    targetType: StatefulSet|Deployment|Pod|PersistentVolume|PersistentVolumeClaim|Service
    # How to select targets from filtered resources
    mode: One|All|Random|Percentage|Fixed
    value: string  # Used with mode (e.g., "50%" for Percentage)

  # Type of chaos to inject
  chaosType: DiskFailure|NetworkLatency|DatabaseConnectionDisruption|PodFailure|ResourcePressure|DataCorruption|StatefulSetScaling
  
  # Duration of the experiment
  duration: string  # e.g., "5m", "1h30m"
  
  # Intensity of chaos (0.0-1.0)
  intensity: float
  
  # Optional parameters specific to the chaos type
  parameters:
    key1: value1
    key2: value2
  
  # Optional scheduling
  schedule:
    cron: string  # Cron expression
    immediate: boolean
    once: boolean
  
  # Safety mechanisms
  safety:
    autoRollback: boolean
    healthChecks:
      - type: httpGet|tcpSocket|exec
        path: string  # For httpGet
        port: integer
        command: [string]  # For exec
        failureThreshold: integer
    pauseConditions:
      - type: metric|alert|manual
        metricQuery: string
        threshold: string
    resourceProtections:
      - type: Namespace|Label|Annotation|Name
        value: string
status:
  phase: Pending|Running|Completed|Failed
  startTime: time
  endTime: time
  conditions:
    - type: string
      status: string
      reason: string
      message: string
      lastTransitionTime: time
  targetResources:
    - kind: string
      name: string
      namespace: string
      uid: string
      status: string
  failureReason: string
```

## Chaos Types

havock8s supports the following chaos types:

### DiskFailure

Simulates disk I/O failures for stateful applications.

Parameters:
- `errorRate`: Percentage of I/O operations that will fail
- `latency`: Additional latency to add to I/O operations
- `corruptData`: Whether to corrupt data instead of failing operations

### NetworkLatency

Introduces network latency between components.

Parameters:
- `latency`: Amount of latency to add (e.g., "100ms")
- `jitter`: Variation in latency (e.g., "10ms")
- `correlation`: Percentage correlation between successive packets

### DatabaseConnectionDisruption

Disrupts connections to databases.

Parameters:
- `connectionFailureRate`: Percentage of new connections that will fail
- `queryTimeout`: Timeout for queries
- `dropConnections`: Whether to drop existing connections

### PodFailure

Causes pods to fail or restart.

Parameters:
- `killMode`: How to terminate pods (graceful, force, etc.)
- `restartOnly`: Whether to only restart pods instead of killing them

### ResourcePressure

Applies CPU, memory, or disk pressure.

Parameters:
- `cpuLoad`: CPU load to generate (0-100)
- `memoryConsumption`: Memory to consume (e.g., "512Mi")
- `diskFill`: Disk space to fill (e.g., "1Gi")

### DataCorruption

Corrupts data in persistent volumes.

Parameters:
- `corruptionType`: Type of corruption (bit flip, zero out, etc.)
- `dataPattern`: Pattern to use for corruption
- `recoverable`: Whether corruption is recoverable

### StatefulSetScaling

Scales StatefulSets up and down.

Parameters:
- `minReplicas`: Minimum number of replicas
- `maxReplicas`: Maximum number of replicas
- `scalingInterval`: Time between scaling operations 