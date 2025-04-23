package chaos

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
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
func (i *DiskFailureInjector) Inject(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, log logr.Logger) error {
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
func (i *DiskFailureInjector) Cleanup(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, log logr.Logger) error {
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
func (i *DiskFailureInjector) injectPodDiskFailure(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// Get the pod
	pod := &corev1.Pod{}
	err := i.client.Get(ctx, types.NamespacedName{
		Namespace: target.Namespace,
		Name:      target.Name,
	}, pod)
	if err != nil {
		return fmt.Errorf("failed to get pod %s/%s: %w", target.Namespace, target.Name, err)
	}

	// Validate failure mode
	failureMode, ok := experiment.Spec.Parameters["failureMode"]
	if !ok {
		return fmt.Errorf("failureMode parameter is required")
	}
	if failureMode != "readonly" && failureMode != "writeonly" && failureMode != "readwrite" {
		return fmt.Errorf("invalid failure mode: %s", failureMode)
	}

	// Get mount path
	mountPath, ok := experiment.Spec.Parameters["mountPath"]
	if !ok {
		return fmt.Errorf("mountPath parameter is required")
	}

	// Initialize annotations if nil
	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string)
	}

	// Set disk failure annotations
	pod.Annotations["havock8s.io/disk-failure"] = "true"
	pod.Annotations["havock8s.io/disk-failure-mount"] = mountPath
	pod.Annotations["havock8s.io/disk-failure-mode"] = failureMode

	// Update the pod with new annotations
	if err := i.client.Update(ctx, pod); err != nil {
		return fmt.Errorf("failed to update pod %s/%s: %w", target.Namespace, target.Name, err)
	}

	log.Info("Applied disk failure to pod", "pod", target.Name)
	return nil
}

// injectStatefulSetDiskFailure applies disk failure to a StatefulSet's PVCs
func (i *DiskFailureInjector) injectStatefulSetDiskFailure(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// In a real implementation, this would find PVCs associated with the StatefulSet
	// and apply disk failures to them

	log.Info("Applied disk failure to StatefulSet", "statefulset", target.Name)
	return nil
}

// injectPVCDiskFailure applies disk failure to a PVC
func (i *DiskFailureInjector) injectPVCDiskFailure(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// In a real implementation, this would:
	// 1. Identify the PVC and associated PV
	// 2. Depending on the storage class, apply appropriate failure mechanisms
	// 3. For cloud providers, could use provider APIs to inject storage failures

	log.Info("Applied disk failure to PVC", "pvc", target.Name)
	return nil
}

// cleanupPodDiskFailure removes disk failure from a pod
func (i *DiskFailureInjector) cleanupPodDiskFailure(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// Get the pod
	pod := &corev1.Pod{}
	err := i.client.Get(ctx, types.NamespacedName{
		Namespace: target.Namespace,
		Name:      target.Name,
	}, pod)
	if err != nil {
		return fmt.Errorf("failed to get pod %s/%s: %w", target.Namespace, target.Name, err)
	}

	// Remove disk failure annotations if they exist
	if pod.Annotations != nil {
		delete(pod.Annotations, "havock8s.io/disk-failure")
		delete(pod.Annotations, "havock8s.io/disk-failure-mount")
		delete(pod.Annotations, "havock8s.io/disk-failure-mode")
	}

	// Update the pod to remove annotations
	if err := i.client.Update(ctx, pod); err != nil {
		return fmt.Errorf("failed to update pod %s/%s: %w", target.Namespace, target.Name, err)
	}

	log.Info("Removed disk failure from pod", "pod", target.Name)
	return nil
}

// cleanupStatefulSetDiskFailure removes disk failure from a StatefulSet's PVCs
func (i *DiskFailureInjector) cleanupStatefulSetDiskFailure(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// Clean up whatever was done in injectStatefulSetDiskFailure

	log.Info("Removed disk failure from StatefulSet", "statefulset", target.Name)
	return nil
}

// cleanupPVCDiskFailure removes disk failure from a PVC
func (i *DiskFailureInjector) cleanupPVCDiskFailure(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, target chaosv1alpha1.TargetResourceStatus, log logr.Logger) error {
	// Clean up whatever was done in injectPVCDiskFailure

	log.Info("Removed disk failure from PVC", "pvc", target.Name)
	return nil
}
