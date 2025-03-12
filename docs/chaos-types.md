---
layout: default
title: Chaos Types
---

# Chaos Types

Havock8s provides a variety of chaos types specifically designed for testing the resilience of stateful applications. Each chaos type targets different aspects of your stateful workloads.

## Available Chaos Types

<div class="docs-section">
  <div class="docs-card">
    <h3>DiskFailure</h3>
    <p>Simulates disk I/O failures, read/write errors, and latency for persistent volumes.</p>
    <a href="#diskfailure">Learn more →</a>
  </div>
  <div class="docs-card">
    <h3>NetworkLatency</h3>
    <p>Introduces network latency, packet loss, and connection disruptions between stateful components.</p>
    <a href="#networklatency">Learn more →</a>
  </div>
  <div class="docs-card">
    <h3>DatabaseConnectionDisruption</h3>
    <p>Simulates database connection failures, query timeouts, and connection pool exhaustion.</p>
    <a href="#databaseconnectiondisruption">Learn more →</a>
  </div>
  <div class="docs-card">
    <h3>PodFailure</h3>
    <p>Causes pods in StatefulSets to fail, restart, or become unresponsive.</p>
    <a href="#podfailure">Learn more →</a>
  </div>
  <div class="docs-card">
    <h3>ResourcePressure</h3>
    <p>Creates CPU, memory, or disk pressure on stateful workloads.</p>
    <a href="#resourcepressure">Learn more →</a>
  </div>
  <div class="docs-card">
    <h3>DataCorruption</h3>
    <p>Simulates data corruption scenarios in persistent volumes or databases.</p>
    <a href="#datacorruption">Learn more →</a>
  </div>
</div>

## Detailed Chaos Type Reference

<h3 id="diskfailure">DiskFailure</h3>

Simulates various types of disk failures for persistent volumes attached to stateful applications.

#### Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `intensity` | Float | Percentage of I/O operations that will fail (0.0-1.0) | 0.2 |
| `mode` | String | Type of failure: `ReadFailure`, `WriteFailure`, `ReadWriteFailure`, `Latency` | `ReadWriteFailure` |
| `latency` | String | When mode is `Latency`, the amount of latency to add (e.g., "100ms") | "100ms" |
| `targetVolumes` | Array | List of PVCs to target (if empty, all volumes attached to target pods) | [] |

#### Example

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
  intensity: 0.3
  parameters:
    mode: WriteFailure
    targetVolumes:
      - postgres-data-pvc
```

<h3 id="networklatency">NetworkLatency</h3>

Introduces network latency, packet loss, and connection disruptions between stateful components.

#### Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `latency` | String | Amount of latency to add (e.g., "200ms") | "100ms" |
| `jitter` | String | Variation in latency (e.g., "50ms") | "30ms" |
| `packetLoss` | Float | Percentage of packets to drop (0.0-1.0) | 0.0 |
| `correlation` | Float | Correlation between successive packet losses (0.0-1.0) | 0.0 |
| `targetPorts` | Array | List of ports to affect (if empty, all ports) | [] |
| `targetIPs` | Array | List of destination IPs to affect (if empty, all IPs) | [] |

#### Example

```yaml
apiVersion: chaos.havock8s.io/v1alpha1
kind: havock8sExperiment
metadata:
  name: redis-network-latency
spec:
  target:
    selector:
      app: redis
  chaosType: NetworkLatency
  duration: 10m
  parameters:
    latency: "300ms"
    jitter: "100ms"
    packetLoss: 0.05
    targetPorts:
      - 6379
```

<h3 id="databaseconnectiondisruption">DatabaseConnectionDisruption</h3>

Simulates database connection failures, query timeouts, and connection pool exhaustion.

#### Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `failureType` | String | Type of disruption: `ConnectionFailure`, `QueryTimeout`, `ConnectionPoolExhaustion` | `ConnectionFailure` |
| `failureRate` | Float | Percentage of connections/queries to affect (0.0-1.0) | 0.3 |
| `targetPort` | Integer | Database port to target | 5432 |
| `queryPattern` | String | When type is `QueryTimeout`, regex pattern of queries to affect | "" |

#### Example

```yaml
apiVersion: chaos.havock8s.io/v1alpha1
kind: havock8sExperiment
metadata:
  name: mysql-connection-disruption
spec:
  target:
    selector:
      app: mysql
  chaosType: DatabaseConnectionDisruption
  duration: 15m
  parameters:
    failureType: ConnectionPoolExhaustion
    failureRate: 0.5
    targetPort: 3306
```

<h3 id="podfailure">PodFailure</h3>

Causes pods in StatefulSets to fail, restart, or become unresponsive.

#### Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `failureType` | String | Type of failure: `Kill`, `Restart`, `Unresponsive` | `Restart` |
| `podIndexes` | Array | For StatefulSets, indexes of pods to target (if empty, random selection) | [] |
| `count` | Integer | Number of pods to affect | 1 |
| `interval` | String | For multiple pods, interval between failures | "10s" |

#### Example

```yaml
apiVersion: chaos.havock8s.io/v1alpha1
kind: havock8sExperiment
metadata:
  name: zookeeper-pod-failure
spec:
  target:
    selector:
      app: zookeeper
  chaosType: PodFailure
  duration: 5m
  parameters:
    failureType: Kill
    podIndexes: [0, 2]  # Target first and third pods in the StatefulSet
```

<h3 id="resourcepressure">ResourcePressure</h3>

Creates CPU, memory, or disk pressure on stateful workloads.

#### Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `resourceType` | String | Type of resource pressure: `CPU`, `Memory`, `Disk` | `CPU` |
| `intensity` | Float | Percentage of resource to consume (0.0-1.0) | 0.8 |
| `workers` | Integer | Number of worker processes to create pressure | 1 |
| `path` | String | For disk pressure, path to fill | "/data" |

#### Example

```yaml
apiVersion: chaos.havock8s.io/v1alpha1
kind: havock8sExperiment
metadata:
  name: elasticsearch-memory-pressure
spec:
  target:
    selector:
      app: elasticsearch
  chaosType: ResourcePressure
  duration: 8m
  parameters:
    resourceType: Memory
    intensity: 0.7
    workers: 2
```

<h3 id="datacorruption">DataCorruption</h3>

Simulates data corruption scenarios in persistent volumes or databases.

#### Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `corruptionType` | String | Type of corruption: `BitFlip`, `Truncate`, `Append`, `Replace` | `BitFlip` |
| `path` | String | Path to target files | "/data" |
| `filePattern` | String | Pattern of files to corrupt | "*.db" |
| `percentage` | Float | Percentage of matching files to corrupt (0.0-1.0) | 0.1 |
| `bytes` | Integer | Number of bytes to corrupt | 1 |

#### Example

```yaml
apiVersion: chaos.havock8s.io/v1alpha1
kind: havock8sExperiment
metadata:
  name: mongodb-data-corruption
spec:
  target:
    selector:
      app: mongodb
  chaosType: DataCorruption
  duration: 1m  # Short duration as this is destructive
  parameters:
    corruptionType: BitFlip
    path: "/data/db"
    filePattern: "*.wt"
    percentage: 0.05
    bytes: 2
  safety:
    autoRollback: true
```

## Safety Considerations

When running chaos experiments, especially those targeting stateful applications, it's important to implement proper safety measures:

<div class="warning">
  <strong>Warning:</strong> Some chaos types like DataCorruption can cause permanent data loss if not properly configured with safety mechanisms.
</div>

Always use the safety parameters in your experiments:

```yaml
safety:
  autoRollback: true  # Automatically roll back if health checks fail
  healthChecks:
    - type: httpGet
      path: /health
      port: 8080
      failureThreshold: 3
  maxTargetPods: 1  # Limit the number of pods affected
  targetPercentage: 30  # Only target 30% of eligible pods
```

## Creating Custom Chaos Types

Havock8s is designed to be extensible. See the [Developer Guide](developer-guide.html) for instructions on creating custom chaos types for your specific stateful application needs. 