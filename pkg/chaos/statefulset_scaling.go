package chaos

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatefulSetScalingInjector implements the Injector interface for StatefulSet scaling chaos
type StatefulSetScalingInjector struct {
	client client.Client
}

// SetClient sets the Kubernetes client
func (i *StatefulSetScalingInjector) SetClient(c client.Client) {
	i.client = c
}

// Inject applies StatefulSet scaling chaos
func (i *StatefulSetScalingInjector) Inject(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, log logr.Logger) error {
	log.Info("Injecting StatefulSet scaling chaos")

	// Default scaling behavior is to scale down by 1
	scaleMode := "down"
	scaleCount := int32(1)
	scaleMin := int32(1)  // Minimum replicas to maintain
	scaleMax := int32(10) // Maximum replicas to scale to
	allowZero := false    // Whether to allow scaling to 0

	// Get parameters with defaults
	if val, ok := experiment.Spec.Parameters["scaleMode"]; ok {
		if val == "up" || val == "down" || val == "random" {
			scaleMode = val
		}
	}

	if val, ok := experiment.Spec.Parameters["scaleCount"]; ok {
		count, err := strconv.Atoi(val)
		if err == nil && count > 0 {
			scaleCount = int32(count)
		}
	}

	if val, ok := experiment.Spec.Parameters["scaleMin"]; ok {
		min, err := strconv.Atoi(val)
		if err == nil && min >= 0 {
			scaleMin = int32(min)
		}
	}

	if val, ok := experiment.Spec.Parameters["scaleMax"]; ok {
		max, err := strconv.Atoi(val)
		if err == nil && max > 0 {
			scaleMax = int32(max)
		}
	}

	if val, ok := experiment.Spec.Parameters["allowZero"]; ok {
		allowZero = val == "true"
		if allowZero {
			scaleMin = 0
		}
	}

	log.Info("StatefulSet scaling parameters",
		"scaleMode", scaleMode,
		"scaleCount", scaleCount,
		"scaleMin", scaleMin,
		"scaleMax", scaleMax,
		"allowZero", allowZero)

	// Apply chaos to each target
	for _, target := range experiment.Status.TargetResources {
		log.Info("Processing target for StatefulSet scaling", "kind", target.Kind, "name", target.Name, "namespace", target.Namespace)

		if target.Kind != "StatefulSet" {
			log.Info("Skipping non-StatefulSet target", "kind", target.Kind)
			continue
		}

		// Get the StatefulSet
		sts := &appsv1.StatefulSet{}
		err := i.client.Get(ctx, types.NamespacedName{
			Namespace: target.Namespace,
			Name:      target.Name,
		}, sts)
		if err != nil {
			log.Error(err, "Failed to get StatefulSet", "StatefulSet", target.Name)
			return fmt.Errorf("failed to get StatefulSet %s/%s: %w", target.Namespace, target.Name, err)
		}

		// Store original replicas in annotations for later restoration
		if sts.Annotations == nil {
			sts.Annotations = make(map[string]string)
		}

		// Save original replica count
		originalReplicas := *sts.Spec.Replicas
		sts.Annotations["havock8s.io/original-replicas"] = fmt.Sprintf("%d", originalReplicas)

		// Calculate new replica count based on scale mode
		var newReplicas int32
		switch scaleMode {
		case "down":
			// Scale down by scaleCount, but not below scaleMin
			newReplicas = originalReplicas - scaleCount
			if newReplicas < scaleMin {
				newReplicas = scaleMin
			}
		case "up":
			// Scale up by scaleCount, but not above scaleMax
			newReplicas = originalReplicas + scaleCount
			if newReplicas > scaleMax {
				newReplicas = scaleMax
			}
		case "random":
			// In a real implementation, this would select a random value
			// For simplicity, we'll just use scaleMin here
			newReplicas = scaleMin
		}

		// Update the StatefulSet with new replica count
		sts.Spec.Replicas = &newReplicas
		if err := i.client.Update(ctx, sts); err != nil {
			log.Error(err, "Failed to update StatefulSet replicas", "StatefulSet", target.Name)
			return err
		}

		log.Info("Successfully scaled StatefulSet",
			"StatefulSet", target.Name,
			"originalReplicas", originalReplicas,
			"newReplicas", newReplicas)
	}

	log.Info("StatefulSet scaling chaos injection completed")
	return nil
}

// Cleanup reverts StatefulSet scaling changes
func (i *StatefulSetScalingInjector) Cleanup(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, log logr.Logger) error {
	log.Info("Cleaning up StatefulSet scaling chaos")

	for _, target := range experiment.Status.TargetResources {
		if target.Kind != "StatefulSet" {
			continue
		}

		// Get the StatefulSet
		sts := &appsv1.StatefulSet{}
		err := i.client.Get(ctx, types.NamespacedName{
			Namespace: target.Namespace,
			Name:      target.Name,
		}, sts)
		if err != nil {
			log.Error(err, "Failed to get StatefulSet for cleanup", "StatefulSet", target.Name)
			return fmt.Errorf("failed to get StatefulSet %s/%s: %w", target.Namespace, target.Name, err)
		}

		// Check if we have stored the original replica count
		if sts.Annotations != nil {
			if originalReplicasStr, ok := sts.Annotations["havock8s.io/original-replicas"]; ok {
				originalReplicas, err := strconv.Atoi(originalReplicasStr)
				if err != nil {
					log.Error(err, "Failed to parse original replicas", "value", originalReplicasStr)
					continue
				}

				// Restore original replica count
				replicas := int32(originalReplicas)
				sts.Spec.Replicas = &replicas

				// Update the StatefulSet
				if err := i.client.Update(ctx, sts); err != nil {
					log.Error(err, "Failed to restore StatefulSet replicas", "StatefulSet", target.Name)
					return err
				}

				// Remove our annotation
				delete(sts.Annotations, "havock8s.io/original-replicas")
				if err := i.client.Update(ctx, sts); err != nil {
					log.Error(err, "Failed to update StatefulSet annotations", "StatefulSet", target.Name)
					return err
				}

				log.Info("Successfully restored StatefulSet replicas",
					"StatefulSet", target.Name,
					"replicas", originalReplicas)
			}
		}
	}

	log.Info("StatefulSet scaling chaos cleanup completed")
	return nil
}
