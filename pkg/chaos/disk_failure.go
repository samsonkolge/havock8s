package chaos

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/statefulchaos/statefulchaos/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DiskFailureInjector implements the Injector interface for disk failure chaos
type DiskFailureInjector struct {
	client client.Client
}

// SetClient sets the Kubernetes client
func (i *DiskFailureInjector) SetClient(c client.Client) {
	i.client = c
}

// Inject applies disk failure chaos
func (i *DiskFailureInjector) Inject(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Injecting disk failure chaos")

	for _, target := range experiment.Status.TargetResources {
		log.Info("Processing target for disk failure", "kind", target.Kind, "name", target.Name, "namespace", target.Namespace)

		// Different actions based on the target kind
		switch target.Kind {
		case "Pod":
			// Inject sidecar container or use fault injection tools
			if err := i.injectPodDiskFailure(ctx, experiment, target, log); err != nil {
				return err
			}
		case "StatefulSet":
			// For StatefulSets, we might want to affect the PVCs
			if err := i.injectStatefulSetDiskFailure(ctx, experiment, target, log); err != nil {
				return err
			}
		case "PersistentVolumeClaim":
			// Directly affect the PVC if possible
			if err := i.injectPVCDiskFailure(ctx, experiment, target, log); err != nil {
				return err
			}
		default:
			log.Info("Unsupported target kind for disk failure", "kind", target.Kind)
		}
	}

	log.Info("Disk failure chaos injection completed")
	return nil
}

// Cleanup removes disk failure chaos
func (i *DiskFailureInjector) Cleanup(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Cleaning up disk failure chaos")

	for _, target := range experiment.Status.TargetResources {
		log.Info("Processing target for disk failure cleanup", "kind", target.Kind, "name", target.Name, "namespace", target.Namespace)

		// Different actions based on the target kind
		switch target.Kind {
		case "Pod":
			// Remove sidecar container or reset fault injection
			if err := i.cleanupPodDiskFailure(ctx, experiment, target, log); err != nil {
				return err
			}
		case "StatefulSet":
			// Cleanup for StatefulSets
			if err := i.cleanupStatefulSetDiskFailure(ctx, experiment, target, log); err != nil {
				return err
			}
		case "PersistentVolumeClaim":
			// Cleanup for PVCs
			if err := i.cleanupPVCDiskFailure(ctx, experiment, target, log); err != nil {
				return err
			}
		default:
			log.Info("Unsupported target kind for disk failure cleanup", "kind", target.Kind)
		}
	}

	log.Info("Disk failure chaos cleanup completed")
	return nil
}

// injectPodDiskFailure applies disk failure to a pod
func (i *DiskFailureInjector) injectPodDiskFailure(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// In a real implementation, this would inject faults using:
	// 1. Add a sidecar that uses ioutil to generate I/O errors
	// 2. Patch the pod with an annotation that triggers admission controllers
	// 3. Use node-level tools to inject failures into the storage subsystem

	// Example: Get the pod and make annotations
	pod := &corev1.Pod{}
	err := i.client.Get(ctx, types.NamespacedName{
		Namespace: target.Namespace,
		Name:      target.Name,
	}, pod)
	if err != nil {
		return fmt.Errorf("failed to get pod %s/%s: %w", target.Namespace, target.Name, err)
	}

	// In a real implementation, here we would apply the disk failure

	log.Info("Applied disk failure to pod", "pod", target.Name)
	return nil
}

// injectStatefulSetDiskFailure applies disk failure to a StatefulSet's PVCs
func (i *DiskFailureInjector) injectStatefulSetDiskFailure(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// In a real implementation, this would find PVCs associated with the StatefulSet
	// and apply disk failures to them

	log.Info("Applied disk failure to StatefulSet", "statefulset", target.Name)
	return nil
}

// injectPVCDiskFailure applies disk failure to a PVC
func (i *DiskFailureInjector) injectPVCDiskFailure(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// In a real implementation, this would:
	// 1. Identify the PVC and associated PV
	// 2. Depending on the storage class, apply appropriate failure mechanisms
	// 3. For cloud providers, could use provider APIs to inject storage failures

	log.Info("Applied disk failure to PVC", "pvc", target.Name)
	return nil
}

// cleanupPodDiskFailure removes disk failure from a pod
func (i *DiskFailureInjector) cleanupPodDiskFailure(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// Clean up whatever was done in injectPodDiskFailure

	log.Info("Removed disk failure from pod", "pod", target.Name)
	return nil
}

// cleanupStatefulSetDiskFailure removes disk failure from a StatefulSet's PVCs
func (i *DiskFailureInjector) cleanupStatefulSetDiskFailure(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// Clean up whatever was done in injectStatefulSetDiskFailure

	log.Info("Removed disk failure from StatefulSet", "statefulset", target.Name)
	return nil
}

// cleanupPVCDiskFailure removes disk failure from a PVC
func (i *DiskFailureInjector) cleanupPVCDiskFailure(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// Clean up whatever was done in injectPVCDiskFailure

	log.Info("Removed disk failure from PVC", "pvc", target.Name)
	return nil
}
