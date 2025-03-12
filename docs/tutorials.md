---
layout: default
title: Tutorials
---

# Tutorials

This section provides step-by-step tutorials for common Havock8s use cases. These tutorials will help you get started with chaos engineering for your stateful applications.

## Tutorial List

<div class="docs-section">
  <div class="docs-card">
    <h3><a href="#postgres-disk-failure">Testing PostgreSQL Resilience to Disk Failures</a></h3>
    <p>Learn how to test a PostgreSQL database's resilience to disk I/O failures.</p>
  </div>
  <div class="docs-card">
    <h3><a href="#redis-network-partition">Simulating Network Partitions in Redis Cluster</a></h3>
    <p>Test how a Redis cluster handles network partitions between nodes.</p>
  </div>
  <div class="docs-card">
    <h3><a href="#mongodb-pod-failure">Testing MongoDB Replica Set Recovery</a></h3>
    <p>Verify that a MongoDB replica set can recover from pod failures.</p>
  </div>
</div>

<h2 id="postgres-disk-failure">Testing PostgreSQL Resilience to Disk Failures</h2>

This tutorial demonstrates how to test a PostgreSQL database's resilience to disk I/O failures.

### Prerequisites

- Kubernetes cluster with Havock8s installed
- PostgreSQL StatefulSet running in your cluster
- Basic understanding of PostgreSQL architecture

### Step 1: Deploy a PostgreSQL StatefulSet

If you don't already have PostgreSQL running, you can deploy it using the following manifest:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  labels:
    app: postgres
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:14
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_PASSWORD
          value: password
        - name: PGDATA
          value: /var/lib/postgresql/data/pgdata
        volumeMounts:
        - name: postgres-data
          mountPath: /var/lib/postgresql/data
        readinessProbe:
          exec:
            command:
            - pg_isready
          initialDelaySeconds: 5
          periodSeconds: 10
  volumeClaimTemplates:
  - metadata:
      name: postgres-data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
  clusterIP: None
```

Apply this manifest:

```bash
kubectl apply -f postgres.yaml
```

### Step 2: Create a Test Database and Table

Connect to the PostgreSQL pod and create a test database and table:

```bash
# Connect to the PostgreSQL pod
kubectl exec -it postgres-0 -- bash

# Connect to PostgreSQL
psql -U postgres

# Create a test database
CREATE DATABASE testdb;
\c testdb

# Create a test table
CREATE TABLE test_data (
  id SERIAL PRIMARY KEY,
  data TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

# Insert some test data
INSERT INTO test_data (data) 
SELECT 'Test data ' || i 
FROM generate_series(1, 1000) AS i;

# Exit PostgreSQL and the pod
\q
exit
```

### Step 3: Create a Monitoring Pod

Deploy a simple pod that will continuously query the database to monitor its availability:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: postgres-monitor
spec:
  containers:
  - name: monitor
    image: postgres:14
    command:
    - /bin/bash
    - -c
    - |
      while true; do
        PGPASSWORD=password psql -h postgres -U postgres -d testdb -c "SELECT COUNT(*) FROM test_data;" || echo "Query failed"
        sleep 5
      done
```

Apply this manifest:

```bash
kubectl apply -f postgres-monitor.yaml
```

### Step 4: Create a Disk Failure Experiment

Now, create a Havock8s experiment to simulate disk I/O failures:

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
  parameters:
    mode: WriteFailure
    intensity: 0.3
    targetVolumes:
      - postgres-data
  safety:
    autoRollback: true
    healthChecks:
      - type: exec
        command:
          - pg_isready
        failureThreshold: 5
```

Apply this experiment:

```bash
kubectl apply -f postgres-disk-failure.yaml
```

### Step 5: Monitor the Experiment

Watch the logs of the monitoring pod to see how PostgreSQL handles the disk failures:

```bash
kubectl logs -f postgres-monitor
```

You should see some queries failing with errors related to disk I/O.

### Step 6: Analyze the Results

After the experiment completes, check the PostgreSQL logs to see how it handled the disk failures:

```bash
kubectl logs postgres-0
```

Look for error messages, recovery attempts, and any data integrity issues.

### Step 7: Verify Data Integrity

Connect to PostgreSQL again and verify that the data is still intact:

```bash
kubectl exec -it postgres-0 -- bash
psql -U postgres -d testdb
SELECT COUNT(*) FROM test_data;
\q
exit
```

### Conclusion

This tutorial demonstrated how to test PostgreSQL's resilience to disk failures using Havock8s. You can modify the experiment parameters to test different failure scenarios, such as read failures, latency, or varying intensities.

<h2 id="redis-network-partition">Simulating Network Partitions in Redis Cluster</h2>

This tutorial shows how to test a Redis cluster's behavior during network partitions.

### Prerequisites

- Kubernetes cluster with Havock8s installed
- Redis cluster running in your Kubernetes cluster
- Basic understanding of Redis cluster architecture

### Step 1: Deploy a Redis Cluster

If you don't already have a Redis cluster running, you can deploy one using the following manifest:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis
  labels:
    app: redis
spec:
  serviceName: redis
  replicas: 3
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:6
        command:
        - redis-server
        - --cluster-enabled yes
        - --cluster-config-file /data/nodes.conf
        - --cluster-node-timeout 5000
        ports:
        - containerPort: 6379
        volumeMounts:
        - name: redis-data
          mountPath: /data
  volumeClaimTemplates:
  - metadata:
      name: redis-data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
---
apiVersion: v1
kind: Service
metadata:
  name: redis
spec:
  selector:
    app: redis
  ports:
  - port: 6379
    targetPort: 6379
  clusterIP: None
```

Apply this manifest:

```bash
kubectl apply -f redis.yaml
```

### Step 2: Initialize the Redis Cluster

Initialize the Redis cluster:

```bash
# Get the pod IPs
PODS=$(kubectl get pods -l app=redis -o jsonpath='{range.items[*]}{.status.podIP}{" "}{end}')
POD_IPS=($PODS)

# Create the cluster
kubectl exec -it redis-0 -- redis-cli --cluster create ${POD_IPS[0]}:6379 ${POD_IPS[1]}:6379 ${POD_IPS[2]}:6379 --cluster-replicas 0
```

### Step 3: Create a Network Partition Experiment

Create a Havock8s experiment to simulate a network partition between Redis nodes:

```yaml
apiVersion: chaos.havock8s.io/v1alpha1
kind: havock8sExperiment
metadata:
  name: redis-network-partition
spec:
  target:
    selector:
      app: redis
    mode: Fixed
    value: "1"  # Target only one pod
  chaosType: NetworkLatency
  duration: 3m
  parameters:
    packetLoss: 1.0  # 100% packet loss = complete partition
    targetPorts:
      - 6379
  safety:
    autoRollback: true
```

Apply this experiment:

```bash
kubectl apply -f redis-network-partition.yaml
```

### Step 4: Monitor the Cluster State

Watch the Redis cluster state during the experiment:

```bash
# In one terminal, watch the cluster info
kubectl exec -it redis-0 -- watch -n 1 redis-cli cluster info

# In another terminal, watch the cluster nodes
kubectl exec -it redis-0 -- watch -n 1 redis-cli cluster nodes
```

### Step 5: Test Cluster Availability

During the experiment, test if the cluster is still available for reads and writes:

```bash
# Connect to a Redis pod
kubectl exec -it redis-1 -- bash

# Use redis-cli to set and get values
redis-cli -c set testkey "Hello Chaos"
redis-cli -c get testkey

# Exit the pod
exit
```

### Step 6: Analyze the Results

After the experiment completes, analyze how the Redis cluster handled the network partition:

- Did the cluster detect the partition?
- How long did it take to detect the partition?
- Did the cluster continue to serve requests?
- Did any data loss occur?

### Conclusion

This tutorial demonstrated how to test a Redis cluster's resilience to network partitions using Havock8s. You can modify the experiment parameters to test different network failure scenarios, such as partial packet loss, latency, or targeting specific nodes.

<h2 id="mongodb-pod-failure">Testing MongoDB Replica Set Recovery</h2>

This tutorial shows how to test a MongoDB replica set's ability to recover from pod failures.

### Prerequisites

- Kubernetes cluster with Havock8s installed
- MongoDB replica set running in your Kubernetes cluster
- Basic understanding of MongoDB replica set architecture

### Step 1: Deploy a MongoDB Replica Set

If you don't already have a MongoDB replica set running, you can deploy one using the following manifest:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mongodb
  labels:
    app: mongodb
spec:
  serviceName: mongodb
  replicas: 3
  selector:
    matchLabels:
      app: mongodb
  template:
    metadata:
      labels:
        app: mongodb
    spec:
      containers:
      - name: mongodb
        image: mongo:4.4
        command:
        - mongod
        - --replSet
        - rs0
        - --bind_ip_all
        ports:
        - containerPort: 27017
        volumeMounts:
        - name: mongodb-data
          mountPath: /data/db
  volumeClaimTemplates:
  - metadata:
      name: mongodb-data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: 1Gi
---
apiVersion: v1
kind: Service
metadata:
  name: mongodb
spec:
  selector:
    app: mongodb
  ports:
  - port: 27017
    targetPort: 27017
  clusterIP: None
```

Apply this manifest:

```bash
kubectl apply -f mongodb.yaml
```

### Step 2: Initialize the MongoDB Replica Set

Initialize the MongoDB replica set:

```bash
# Connect to the first MongoDB pod
kubectl exec -it mongodb-0 -- bash

# Connect to MongoDB shell
mongo

# Initialize the replica set
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "mongodb-0.mongodb:27017" },
    { _id: 1, host: "mongodb-1.mongodb:27017" },
    { _id: 2, host: "mongodb-2.mongodb:27017" }
  ]
})

# Check the replica set status
rs.status()

# Create a test database and collection
use testdb
db.testcollection.insertMany([
  { name: "Document 1", value: 1 },
  { name: "Document 2", value: 2 },
  { name: "Document 3", value: 3 }
])

# Exit MongoDB and the pod
exit
exit
```

### Step 3: Create a Pod Failure Experiment

Create a Havock8s experiment to simulate pod failures in the MongoDB replica set:

```yaml
apiVersion: chaos.havock8s.io/v1alpha1
kind: havock8sExperiment
metadata:
  name: mongodb-pod-failure
spec:
  target:
    selector:
      app: mongodb
  chaosType: PodFailure
  duration: 5m
  parameters:
    failureType: Kill
    podIndexes: [0]  # Target the primary node (mongodb-0)
  safety:
    autoRollback: true
```

Apply this experiment:

```bash
kubectl apply -f mongodb-pod-failure.yaml
```

### Step 4: Monitor the Replica Set Status

Watch the replica set status during the experiment:

```bash
# Connect to one of the remaining MongoDB pods
kubectl exec -it mongodb-1 -- mongo

# Check the replica set status
rs.status()

# Keep checking the status to see the election of a new primary
```

### Step 5: Test Replica Set Availability

During the experiment, test if the replica set is still available for reads and writes:

```bash
# In the MongoDB shell
use testdb

# Try to write data (this should work once a new primary is elected)
db.testcollection.insertOne({ name: "Document during chaos", value: 42 })

# Read data
db.testcollection.find()
```

### Step 6: Analyze the Results

After the experiment completes and the failed pod recovers, analyze how the MongoDB replica set handled the failure:

- How quickly was a new primary elected?
- Was there any data loss?
- Did the failed pod successfully rejoin the replica set?
- Were there any unexpected behaviors?

### Conclusion

This tutorial demonstrated how to test a MongoDB replica set's resilience to pod failures using Havock8s. You can modify the experiment parameters to test different failure scenarios, such as multiple pod failures, repeated failures, or targeting specific nodes in the replica set.

## Next Steps

Now that you've completed these tutorials, you can:

1. Explore other chaos types available in Havock8s
2. Create more complex experiments with multiple chaos types
3. Integrate Havock8s experiments into your CI/CD pipeline
4. Develop custom chaos types for your specific applications

Check out the [API Reference](api-reference.html) and [Developer Guide](developer-guide.html) for more information. 