# StatefulChaos Developer Guide

This guide explains how to extend the StatefulChaos operator with new chaos types, enhance existing functionality, and contribute to the project.

## Architecture Overview

StatefulChaos follows the Kubernetes Operator pattern with a controller-based architecture:

```
┌─────────────────────────┐
│   StatefulChaos CRD     │
│ StatefulChaosExperiment │
└───────────┬─────────────┘
            │
            ▼
┌─────────────────────────┐
│ StatefulChaos Controller│
└───────────┬─────────────┘
            │
      ┌─────┴─────┐
      │           │
      ▼           ▼
┌──────────┐ ┌────────────────┐
│ Injectors │ │ Safety Monitors│
└──────────┘ └────────────────┘
```

- **Custom Resource Definition (CRD)**: Defines the `StatefulChaosExperiment` resource
- **Controller**: Reconciles the desired state with the actual state
- **Injectors**: Implement various chaos types (disk failures, network latency, etc.)
- **Safety Monitors**: Implement safety mechanisms to prevent cascading failures

## Adding a New Chaos Type

### 1. Define the Chaos Type in the CRD

First, add your new chaos type to the enum in `api/v1alpha1/statefulchaosexperiment_types.go`:

```go
// ChaosType defines the type of chaos to be injected
// +kubebuilder:validation:Enum=DiskFailure;NetworkLatency;DatabaseConnectionDisruption;PodFailure;ResourcePressure;DataCorruption;StatefulSetScaling;YourNewChaosType
ChaosType string `json:"chaosType"`
```

### 2. Implement the Injector

Create a new file in `pkg/chaos/your_chaos_type.go` that implements both injection and cleanup:

```go
package chaos

import (
    "context"
    "github.com/go-logr/logr"
    chaosv1alpha1 "github.com/statefulchaos/statefulchaos/api/v1alpha1"
)

// YourChaosTypeInjector implements chaos for your specific use case
type YourChaosTypeInjector struct {
    // Add necessary fields
}

// Inject applies your chaos type to the target
func (i *YourChaosTypeInjector) Inject(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
    // Implement your chaos injection logic
    return nil
}

// Cleanup removes your chaos type from the target
func (i *YourChaosTypeInjector) Cleanup(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
    // Implement your cleanup logic
    return nil
}
```

### 3. Register Your Chaos Type in the Controller

Modify the `injectChaos` and `removeChaos` methods in `controllers/statefulchaosexperiment_controller.go`:

```go
// injectChaos applies the specified chaos to target resources
func (r *StatefulChaosExperimentReconciler) injectChaos(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
    chaosType := experiment.Spec.ChaosType
    log.Info("Injecting chaos", "type", chaosType)

    // Apply different chaos types
    switch chaosType {
        // Existing chaos types...
        case "YourNewChaosType":
            return r.injectYourChaos(ctx, experiment, log)
        default:
            return fmt.Errorf("unsupported chaos type: %s", chaosType)
    }
}

// Add a new method for your chaos type
func (r *StatefulChaosExperimentReconciler) injectYourChaos(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
    injector := &chaos.YourChaosTypeInjector{}
    return injector.Inject(ctx, experiment, log)
}

// Don't forget to add the cleanup method as well
func (r *StatefulChaosExperimentReconciler) removeYourChaos(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
    injector := &chaos.YourChaosTypeInjector{}
    return injector.Cleanup(ctx, experiment, log)
}
```

### 4. Create Examples

Add an example in the `examples/` directory:

```yaml
apiVersion: chaos.statefulchaos.io/v1alpha1
kind: StatefulChaosExperiment
metadata:
  name: your-chaos-example
spec:
  target:
    selector:
      matchLabels:
        app: your-app
    targetType: StatefulSet
  chaosType: YourNewChaosType
  duration: 5m
  intensity: 0.5
  parameters:
    yourParam1: "value1"
    yourParam2: "value2"
```

## Implementing a New Safety Feature

### 1. Extend the Safety Spec

Add your new safety feature to the SafetySpec in `api/v1alpha1/statefulchaosexperiment_types.go`:

```go
// SafetySpec defines safety mechanisms for chaos experiments
type SafetySpec struct {
    // Existing fields...
    
    // YourNewSafety defines your new safety feature
    // +optional
    YourNewSafety *YourNewSafetySpec `json:"yourNewSafety,omitempty"`
}

// YourNewSafetySpec defines the parameters for your safety feature
type YourNewSafetySpec struct {
    // Define your safety feature fields
    Enabled bool `json:"enabled"`
    Parameter string `json:"parameter"`
}
```

### 2. Implement the Safety Logic

Add your safety check to the `shouldRollbackExperiment` method in `controllers/statefulchaosexperiment_controller.go`:

```go
// shouldRollbackExperiment checks if experiment should be rolled back based on safety conditions
func (r *StatefulChaosExperimentReconciler) shouldRollbackExperiment(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment) (bool, string) {
    // Existing checks...
    
    // Check your new safety feature
    if experiment.Spec.Safety != nil && experiment.Spec.Safety.YourNewSafety != nil && experiment.Spec.Safety.YourNewSafety.Enabled {
        // Implement your safety check
        if shouldStop, reason := r.checkYourSafety(ctx, experiment); shouldStop {
            return true, reason
        }
    }
    
    return false, ""
}

// checkYourSafety implements your custom safety check
func (r *StatefulChaosExperimentReconciler) checkYourSafety(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment) (bool, string) {
    // Implement your safety check logic
    return false, ""
}
```

## Testing Your Changes

### Unit Tests

Create or modify tests in the `controllers/statefulchaosexperiment_controller_test.go` file:

```go
func TestYourNewChaosType(t *testing.T) {
    // Implement your test
}
```

### Integration Tests

Create an integration test in the `test/` directory.

### Manual Testing

1. Build and deploy the operator:
   ```bash
   make docker-build
   make deploy
   ```

2. Create a test experiment:
   ```bash
   kubectl apply -f examples/your-example.yaml
   ```

3. Monitor experiment progress:
   ```bash
   kubectl get sce
   kubectl describe sce your-chaos-example
   ```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for your changes
5. Run tests: `make test`
6. Submit a pull request

## Best Practices

1. **Safety First**: Always implement proper safety mechanisms for new chaos types
2. **Logging**: Add detailed logs for debugging
3. **Error Handling**: Handle all possible error conditions
4. **Idempotency**: Ensure your controller is idempotent
5. **Documentation**: Document your chaos type and add examples 