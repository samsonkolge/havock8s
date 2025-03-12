---
layout: documentation
title: Chaos Types
prev_page: installation.html
prev_title: Installation Guide
next_page: api-reference.html
next_title: API Reference
---

# Chaos Types

<div class="callout callout-info">
  <div class="callout-title">
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
    Understanding Chaos Types
  </div>
  <p>Havock8s offers a variety of chaos types specifically designed to test the resilience of stateful applications on Kubernetes.</p>
</div>

## Introduction to Chaos Types

Havock8s provides specialized chaos experiments for stateful applications. Each chaos type targets a specific aspect of your application's infrastructure or behavior to help you discover how your system responds to different failure modes.

<div class="article-section">
  <p>When selecting a chaos type, consider:</p>
  <ul>
    <li>What aspect of your application you want to test</li>
    <li>The potential impact on your production environment</li>
    <li>The specific failure modes you want to simulate</li>
    <li>The expected behavior of your application under these conditions</li>
  </ul>
</div>

## Available Chaos Types

Havock8s offers the following chaos types, organized by category:

### Storage Chaos

<div class="docs-section">
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>DiskFailure</h3>
    </div>
    <div class="docs-card-content">
      <p>Simulates disk I/O errors, latency, or complete disk failures to test how your application handles storage problems.</p>
      <h4>Parameters:</h4>
      <ul>
        <li><strong>mode</strong>: latency, error, or corruption</li>
        <li><strong>intensity</strong>: 0.0-1.0 (percentage of operations affected)</li>
        <li><strong>devices</strong>: List of target devices (optional)</li>
        <li><strong>paths</strong>: List of target filesystem paths (optional)</li>
      </ul>
      <h4>Example:</h4>
      <pre><code>spec:
  chaosType: DiskFailure
  mode: error
  intensity: 0.3
  paths: ["/data", "/var/lib/postgresql/data"]</code></pre>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>VolumeFailure</h3>
    </div>
    <div class="docs-card-content">
      <p>Simulates Kubernetes persistent volume failures, including unmounting, permissions issues, or capacity problems.</p>
      <h4>Parameters:</h4>
      <ul>
        <li><strong>failureType</strong>: unmount, permission, or capacity</li>
        <li><strong>duration</strong>: How long the failure persists</li>
        <li><strong>volumeClaimNames</strong>: Target PVCs (optional)</li>
      </ul>
      <h4>Example:</h4>
      <pre><code>spec:
  chaosType: VolumeFailure
  failureType: unmount
  duration: 2m
  volumeClaimNames: ["postgres-data"]</code></pre>
    </div>
  </div>
</div>

### Network Chaos

<div class="docs-section">
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>NetworkPartition</h3>
    </div>
    <div class="docs-card-content">
      <p>Creates network partitions between stateful components, simulating network splits in distributed systems.</p>
      <h4>Parameters:</h4>
      <ul>
        <li><strong>direction</strong>: ingress, egress, or both</li>
        <li><strong>targetLabels</strong>: Labels of pods to partition</li>
        <li><strong>duration</strong>: How long the partition lasts</li>
      </ul>
      <h4>Example:</h4>
      <pre><code>spec:
  chaosType: NetworkPartition
  direction: both
  targetLabels:
    role: "master"
  duration: 5m</code></pre>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>NetworkLatency</h3>
    </div>
    <div class="docs-card-content">
      <p>Introduces latency into network connections between stateful components or between apps and databases.</p>
      <h4>Parameters:</h4>
      <ul>
        <li><strong>latency</strong>: Added network latency (e.g., "100ms")</li>
        <li><strong>jitter</strong>: Variation in latency (e.g., "20ms")</li>
        <li><strong>correlation</strong>: 0-100% correlation between packets</li>
        <li><strong>ports</strong>: Target ports (optional)</li>
      </ul>
      <h4>Example:</h4>
      <pre><code>spec:
  chaosType: NetworkLatency
  latency: "200ms"
  jitter: "50ms"
  correlation: 80
  ports: [5432, 6379]</code></pre>
    </div>
  </div>
</div>

### State Corruption Chaos

<div class="docs-section">
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>DataCorruption</h3>
    </div>
    <div class="docs-card-content">
      <p>Simulates data corruption in databases or other stateful components to test data integrity mechanisms.</p>
      <h4>Parameters:</h4>
      <ul>
        <li><strong>corruptionType</strong>: byte, schema, or record</li>
        <li><strong>target</strong>: Specific files or database objects</li>
        <li><strong>percentage</strong>: Amount of data to corrupt (0-100%)</li>
      </ul>
      <h4>Example:</h4>
      <pre><code>spec:
  chaosType: DataCorruption
  corruptionType: record
  target: "users_table"
  percentage: 5</code></pre>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>StateDelay</h3>
    </div>
    <div class="docs-card-content">
      <p>Introduces delays in state synchronization between replicas to test convergence behavior.</p>
      <h4>Parameters:</h4>
      <ul>
        <li><strong>delay</strong>: Time to delay state propagation</li>
        <li><strong>replicaSelector</strong>: Which replicas to affect</li>
        <li><strong>stateOperations</strong>: Types of operations to delay</li>
      </ul>
      <h4>Example:</h4>
      <pre><code>spec:
  chaosType: StateDelay
  delay: "30s"
  stateOperations: ["write", "sync"]
  replicaSelector:
    role: "slave"</code></pre>
    </div>
  </div>
</div>

### Resource Chaos

<div class="docs-section">
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>ResourcePressure</h3>
    </div>
    <div class="docs-card-content">
      <p>Induces CPU, memory, or IO pressure on stateful components to test performance degradation scenarios.</p>
      <h4>Parameters:</h4>
      <ul>
        <li><strong>resourceType</strong>: cpu, memory, or io</li>
        <li><strong>intensity</strong>: 0-100% of resource to consume</li>
        <li><strong>duration</strong>: How long to apply the pressure</li>
      </ul>
      <h4>Example:</h4>
      <pre><code>spec:
  chaosType: ResourcePressure
  resourceType: memory
  intensity: 80
  duration: "10m"</code></pre>
    </div>
  </div>
  
  <div class="docs-card">
    <div class="docs-card-header">
      <h3>ConnectionOverload</h3>
    </div>
    <div class="docs-card-content">
      <p>Simulates too many connections to stateful services like databases to test connection pooling and limits.</p>
      <h4>Parameters:</h4>
      <ul>
        <li><strong>connectionCount</strong>: Number of connections to create</li>
        <li><strong>connectionType</strong>: idle, active, or mixed</li>
        <li><strong>port</strong>: Target port number</li>
      </ul>
      <h4>Example:</h4>
      <pre><code>spec:
  chaosType: ConnectionOverload
  connectionCount: 1000
  connectionType: mixed
  port: 5432</code></pre>
    </div>
  </div>
</div>

## Creating Custom Chaos Types

Havock8s is extensible - you can create custom chaos types for your specific needs:

<div class="callout callout-tip">
  <div class="callout-title">
    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"></path></svg>
    Custom Extensions
  </div>
  <p>Refer to our <a href="developer-guide.html">Developer Guide</a> for detailed instructions on creating custom chaos type plugins.</p>
</div>

## Designing Effective Chaos Experiments

When designing chaos experiments, follow these best practices:

1. **Start small** - Begin with low intensity values and short durations
2. **Increase gradually** - Slowly increase the chaos intensity to find breaking points
3. **Monitor carefully** - Always have monitoring in place when running experiments
4. **Focus on recovery** - The goal is to understand and improve recovery mechanisms
5. **Document findings** - Capture all findings and improvements from each experiment

## Next Steps

- [See detailed examples in our Tutorials](tutorials.html)
- [Consult the API Reference](api-reference.html) for detailed parameter information
- [Learn how to combine multiple chaos types](developer-guide.html#combining-chaos-types) 