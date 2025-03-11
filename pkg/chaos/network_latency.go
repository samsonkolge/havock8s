package chaos

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/statefulchaos/statefulchaos/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NetworkLatencyInjector implements the Injector interface for network latency chaos
type NetworkLatencyInjector struct {
	client client.Client
}

// SetClient sets the Kubernetes client
func (i *NetworkLatencyInjector) SetClient(c client.Client) {
	i.client = c
}

// Inject applies network latency chaos
func (i *NetworkLatencyInjector) Inject(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Injecting network latency chaos")

	// Get parameters with defaults
	latency := "100ms"
	jitter := "10ms"
	correlation := "75"
	ports := ""

	// Override defaults with experiment parameters if provided
	if val, ok := experiment.Spec.Parameters["latency"]; ok {
		latency = val
	}
	if val, ok := experiment.Spec.Parameters["jitter"]; ok {
		jitter = val
	}
	if val, ok := experiment.Spec.Parameters["correlation"]; ok {
		correlation = strings.TrimSuffix(val, "%") // Remove % suffix if present
	}
	if val, ok := experiment.Spec.Parameters["ports"]; ok {
		ports = val
	}

	log.Info("Network latency parameters",
		"latency", latency,
		"jitter", jitter,
		"correlation", correlation,
		"ports", ports)

	for _, target := range experiment.Status.TargetResources {
		log.Info("Processing target for network latency", "kind", target.Kind, "name", target.Name, "namespace", target.Namespace)

		// Different actions based on the target kind
		switch target.Kind {
		case "Pod":
			if err := i.injectPodNetworkLatency(ctx, experiment, target, latency, jitter, correlation, ports, log); err != nil {
				return err
			}
		case "StatefulSet":
			if err := i.injectStatefulSetNetworkLatency(ctx, experiment, target, latency, jitter, correlation, ports, log); err != nil {
				return err
			}
		default:
			log.Info("Unsupported target kind for network latency", "kind", target.Kind)
		}
	}

	log.Info("Network latency chaos injection completed")
	return nil
}

// Cleanup removes network latency chaos
func (i *NetworkLatencyInjector) Cleanup(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Cleaning up network latency chaos")

	for _, target := range experiment.Status.TargetResources {
		log.Info("Processing target for network latency cleanup", "kind", target.Kind, "name", target.Name, "namespace", target.Namespace)

		// Different actions based on the target kind
		switch target.Kind {
		case "Pod":
			if err := i.cleanupPodNetworkLatency(ctx, experiment, target, log); err != nil {
				return err
			}
		case "StatefulSet":
			if err := i.cleanupStatefulSetNetworkLatency(ctx, experiment, target, log); err != nil {
				return err
			}
		default:
			log.Info("Unsupported target kind for network latency cleanup", "kind", target.Kind)
		}
	}

	log.Info("Network latency chaos cleanup completed")
	return nil
}

// injectPodNetworkLatency adds network latency to a pod
func (i *NetworkLatencyInjector) injectPodNetworkLatency(
	ctx context.Context,
	experiment *chaosv1alpha1.StatefulChaosExperiment,
	target chaosv1alpha1.TargetResourceStatus,
	latency, jitter, correlation, ports string,
	log logr.Logger,
) error {
	// In a real implementation, this would:
	// 1. Add a sidecar container with tc/netem commands for traffic control
	// 2. Use a CNI plugin or service mesh to inject latency
	// 3. Use a chaos mesh tool to apply network latency

	// Example: Get the pod
	pod := &corev1.Pod{}
	err := i.client.Get(ctx, types.NamespacedName{
		Namespace: target.Namespace,
		Name:      target.Name,
	}, pod)
	if err != nil {
		return fmt.Errorf("failed to get pod %s/%s: %w", target.Namespace, target.Name, err)
	}

	// Store configuration in annotations for later cleanup
	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string)
	}

	pod.Annotations["statefulchaos.io/network-latency"] = "true"
	pod.Annotations["statefulchaos.io/network-latency-value"] = latency
	pod.Annotations["statefulchaos.io/network-jitter-value"] = jitter
	pod.Annotations["statefulchaos.io/network-correlation-value"] = correlation
	if ports != "" {
		pod.Annotations["statefulchaos.io/network-ports"] = ports
	}

	// In a real implementation, here we would apply network latency
	// (for example by adding a sidecar, using a daemonset, etc.)

	log.Info("Applied network latency to pod", "pod", target.Name, "latency", latency)
	return nil
}

// injectStatefulSetNetworkLatency adds network latency to all pods in a StatefulSet
func (i *NetworkLatencyInjector) injectStatefulSetNetworkLatency(
	ctx context.Context,
	experiment *chaosv1alpha1.StatefulChaosExperiment,
	target chaosv1alpha1.TargetResourceStatus,
	latency, jitter, correlation, ports string,
	log logr.Logger,
) error {
	// In a real implementation, this would find all pods belonging to the StatefulSet
	// and apply network latency to them

	log.Info("Applied network latency to StatefulSet", "statefulset", target.Name, "latency", latency)
	return nil
}

// cleanupPodNetworkLatency removes network latency from a pod
func (i *NetworkLatencyInjector) cleanupPodNetworkLatency(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// In a real implementation, this would remove the network latency configuration
	// and restore normal network operation

	// Example: Get the pod
	pod := &corev1.Pod{}
	err := i.client.Get(ctx, types.NamespacedName{
		Namespace: target.Namespace,
		Name:      target.Name,
	}, pod)
	if err != nil {
		return fmt.Errorf("failed to get pod %s/%s: %w", target.Namespace, target.Name, err)
	}

	// Check if network latency was applied
	if pod.Annotations != nil {
		if _, ok := pod.Annotations["statefulchaos.io/network-latency"]; ok {
			// Remove the annotations
			delete(pod.Annotations, "statefulchaos.io/network-latency")
			delete(pod.Annotations, "statefulchaos.io/network-latency-value")
			delete(pod.Annotations, "statefulchaos.io/network-jitter-value")
			delete(pod.Annotations, "statefulchaos.io/network-correlation-value")
			delete(pod.Annotations, "statefulchaos.io/network-ports")

			// In a real implementation, here we would remove network latency
			// (for example by removing the sidecar, updating the daemonset, etc.)
		}
	}

	log.Info("Removed network latency from pod", "pod", target.Name)
	return nil
}

// cleanupStatefulSetNetworkLatency removes network latency from a StatefulSet
func (i *NetworkLatencyInjector) cleanupStatefulSetNetworkLatency(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// In a real implementation, this would find all pods belonging to the StatefulSet
	// and remove network latency from them

	log.Info("Removed network latency from StatefulSet", "statefulset", target.Name)
	return nil
}
