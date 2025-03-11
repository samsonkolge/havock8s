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

// PodFailureInjector implements the Injector interface for pod failure chaos
type PodFailureInjector struct {
	client client.Client
}

// SetClient sets the Kubernetes client
func (i *PodFailureInjector) SetClient(c client.Client) {
	i.client = c
}

// Inject applies pod failure chaos
func (i *PodFailureInjector) Inject(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Injecting pod failure chaos")

	// Get parameters with defaults
	gracePeriod := int64(0) // Default to immediate termination
	forceDelete := false
	podCount := 1

	// Override defaults with experiment parameters if provided
	if val, ok := experiment.Spec.Parameters["gracePeriodSeconds"]; ok {
		var err error
		var period int64
		_, err = fmt.Sscanf(val, "%d", &period)
		if err == nil && period >= 0 {
			gracePeriod = period
		}
	}

	if val, ok := experiment.Spec.Parameters["forceDelete"]; ok {
		forceDelete = val == "true"
	}

	if val, ok := experiment.Spec.Parameters["podCount"]; ok {
		var count int
		_, err := fmt.Sscanf(val, "%d", &count)
		if err == nil && count > 0 {
			podCount = count
		}
	}

	log.Info("Pod failure parameters",
		"gracePeriodSeconds", gracePeriod,
		"forceDelete", forceDelete,
		"podCount", podCount)

	// Track affected pods to record in status
	podsTerminated := 0

	for _, target := range experiment.Status.TargetResources {
		log.Info("Processing target for pod failure", "kind", target.Kind, "name", target.Name, "namespace", target.Namespace)

		switch target.Kind {
		case "Pod":
			if err := i.terminatePod(ctx, experiment, target, gracePeriod, forceDelete, log); err != nil {
				return err
			}
			podsTerminated++
			if podsTerminated >= podCount {
				log.Info("Reached desired pod termination count", "count", podCount)
				return nil
			}

		case "StatefulSet":
			// Find pods belonging to this StatefulSet and terminate them
			pods, err := i.findStatefulSetPods(ctx, target.Namespace, target.Name)
			if err != nil {
				log.Error(err, "Failed to find pods for StatefulSet", "statefulset", target.Name)
				return err
			}

			for _, pod := range pods {
				podTarget := chaosv1alpha1.TargetResourceStatus{
					Kind:      "Pod",
					Name:      pod.Name,
					Namespace: pod.Namespace,
					UID:       string(pod.UID),
					Status:    "Targeted",
				}

				if err := i.terminatePod(ctx, experiment, podTarget, gracePeriod, forceDelete, log); err != nil {
					return err
				}

				podsTerminated++
				if podsTerminated >= podCount {
					log.Info("Reached desired pod termination count", "count", podCount)
					return nil
				}
			}

		default:
			log.Info("Unsupported target kind for pod failure", "kind", target.Kind)
		}
	}

	log.Info("Pod failure chaos injection completed", "podsTerminated", podsTerminated)
	return nil
}

// Cleanup is a no-op for pod failure as pods are already terminated
func (i *PodFailureInjector) Cleanup(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("No cleanup needed for pod failure, Kubernetes will recreate StatefulSet pods")
	return nil
}

// terminatePod terminates a specific pod
func (i *PodFailureInjector) terminatePod(
	ctx context.Context,
	experiment *chaosv1alpha1.StatefulChaosExperiment,
	target chaosv1alpha1.TargetResourceStatus,
	gracePeriod int64,
	forceDelete bool,
	log logr.Logger,
) error {
	log.Info("Terminating pod", "pod", target.Name, "namespace", target.Namespace)

	// Get the pod
	pod := &corev1.Pod{}
	err := i.client.Get(ctx, types.NamespacedName{
		Namespace: target.Namespace,
		Name:      target.Name,
	}, pod)
	if err != nil {
		return fmt.Errorf("failed to get pod %s/%s: %w", target.Namespace, target.Name, err)
	}

	// Record pod state before deletion
	log.Info("Pod details before termination",
		"pod", pod.Name,
		"phase", pod.Status.Phase,
		"containers", len(pod.Spec.Containers))

	// Delete the pod
	deleteOptions := client.DeleteOptions{}
	if gracePeriod >= 0 {
		deleteOptions.GracePeriodSeconds = &gracePeriod
	}

	// For force delete, we would add appropriate options in a real implementation
	// In Kubernetes, force delete typically means deleting without finalizers

	// Delete the pod with specified options
	if err := i.client.Delete(ctx, pod, &deleteOptions); err != nil {
		log.Error(err, "Failed to delete pod", "pod", target.Name)
		return err
	}

	log.Info("Successfully terminated pod", "pod", target.Name)
	return nil
}

// findStatefulSetPods finds all pods belonging to a StatefulSet
func (i *PodFailureInjector) findStatefulSetPods(ctx context.Context, namespace, statefulSetName string) ([]corev1.Pod, error) {
	podList := &corev1.PodList{}

	// In a real implementation, we would use proper label selectors based on
	// StatefulSet labels. For this example, we use a simple naming convention
	// that StatefulSet pods follow: <statefulset-name>-<ordinal>

	// List pods in the namespace
	err := i.client.List(ctx, podList, client.InNamespace(namespace))
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}

	var statefulSetPods []corev1.Pod
	for _, pod := range podList.Items {
		// Check if pod is owned by the StatefulSet
		for _, owner := range pod.OwnerReferences {
			if owner.Kind == "StatefulSet" && owner.Name == statefulSetName {
				statefulSetPods = append(statefulSetPods, pod)
				break
			}
		}
	}

	return statefulSetPods, nil
}
