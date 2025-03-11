package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	chaosv1alpha1 "github.com/statefulchaos/statefulchaos/api/v1alpha1"
)

// StatefulChaosExperimentReconciler reconciles a StatefulChaosExperiment object
type StatefulChaosExperimentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=chaos.statefulchaos.io,resources=statefulchaosexperiments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=chaos.statefulchaos.io,resources=statefulchaosexperiments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=chaos.statefulchaos.io,resources=statefulchaosexperiments/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;delete
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=persistentvolumes,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch

// Reconcile handles the main reconciliation loop for StatefulChaosExperiment resources
func (r *StatefulChaosExperimentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("statefulchaosexperiment", req.NamespacedName)

	// Fetch the StatefulChaosExperiment instance
	experiment := &chaosv1alpha1.StatefulChaosExperiment{}
	err := r.Get(ctx, req.NamespacedName, experiment)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return. Created objects are automatically garbage collected
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		log.Error(err, "Failed to get StatefulChaosExperiment")
		return ctrl.Result{}, err
	}

	// Add finalizer if not present
	if !containsString(experiment.Finalizers, chaosExperimentFinalizer) {
		log.Info("Adding finalizer")
		experiment.Finalizers = append(experiment.Finalizers, chaosExperimentFinalizer)
		if err := r.Update(ctx, experiment); err != nil {
			log.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Check if the experiment is being deleted
	if !experiment.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleExperimentDeletion(ctx, experiment, log)
	}

	// Process the experiment based on its current phase
	switch experiment.Status.Phase {
	case "":
		// New experiment, initialize it
		return r.initializeExperiment(ctx, experiment, log)
	case "Pending":
		// Experiment is waiting to start
		return r.processPendingExperiment(ctx, experiment, log)
	case "Running":
		// Experiment is running, check if it needs to be stopped
		return r.processRunningExperiment(ctx, experiment, log)
	case "Completed", "Failed":
		// Experiment has finished, nothing to do
		return ctrl.Result{}, nil
	default:
		log.Info("Unknown experiment phase", "phase", experiment.Status.Phase)
		return ctrl.Result{}, nil
	}
}

// initializeExperiment sets up a new experiment
func (r *StatefulChaosExperimentReconciler) initializeExperiment(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) (ctrl.Result, error) {
	log.Info("Initializing new experiment")

	// Initialize experiment status
	experiment.Status.Phase = "Pending"

	// Set target resources based on selector
	targetResources, err := r.identifyTargetResources(ctx, experiment)
	if err != nil {
		log.Error(err, "Failed to identify target resources")
		experiment.Status.Phase = "Failed"
		experiment.Status.FailureReason = fmt.Sprintf("Failed to identify targets: %v", err)
		if updateErr := r.Status().Update(ctx, experiment); updateErr != nil {
			log.Error(updateErr, "Failed to update experiment status")
		}
		return ctrl.Result{}, err
	}

	experiment.Status.TargetResources = targetResources

	// Set timestamp for pending state
	now := metav1.Now()
	experiment.Status.StartTime = &now

	// Update status
	if err := r.Status().Update(ctx, experiment); err != nil {
		log.Error(err, "Failed to update experiment status")
		return ctrl.Result{}, err
	}

	// If immediate execution is requested, move to Running right away
	if experiment.Spec.Schedule != nil && experiment.Spec.Schedule.Immediate {
		return ctrl.Result{Requeue: true}, nil
	}

	// Check if we need to schedule the experiment
	if experiment.Spec.Schedule != nil && experiment.Spec.Schedule.Cron != "" {
		// Handle cron scheduling logic - for now, we'll just requeue
		// In a real implementation, use a proper cron scheduler
		return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
	}

	// Default to immediate execution
	return ctrl.Result{Requeue: true}, nil
}

// processPendingExperiment handles experiments in Pending phase
func (r *StatefulChaosExperimentReconciler) processPendingExperiment(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) (ctrl.Result, error) {
	log.Info("Processing pending experiment")

	// Move to Running state
	experiment.Status.Phase = "Running"

	// Add chaos-specific condition
	condition := metav1.Condition{
		Type:               "ChaosInjected",
		Status:             metav1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             "ChaosStarted",
		Message:            fmt.Sprintf("Started chaos experiment of type %s", experiment.Spec.ChaosType),
	}
	experiment.Status.Conditions = append(experiment.Status.Conditions, condition)

	if err := r.Status().Update(ctx, experiment); err != nil {
		log.Error(err, "Failed to update experiment status")
		return ctrl.Result{}, err
	}

	// Inject chaos based on the chaos type
	if err := r.injectChaos(ctx, experiment, log); err != nil {
		log.Error(err, "Failed to inject chaos")
		experiment.Status.Phase = "Failed"
		experiment.Status.FailureReason = fmt.Sprintf("Chaos injection failed: %v", err)
		if updateErr := r.Status().Update(ctx, experiment); updateErr != nil {
			log.Error(updateErr, "Failed to update experiment status after chaos injection failure")
		}
		return ctrl.Result{}, err
	}

	// Calculate experiment duration and set up requeue
	duration, err := time.ParseDuration(experiment.Spec.Duration)
	if err != nil {
		log.Error(err, "Invalid duration format", "duration", experiment.Spec.Duration)
		duration = 5 * time.Minute // Default to 5 minutes if invalid
	}

	return ctrl.Result{RequeueAfter: duration}, nil
}

// processRunningExperiment handles experiments in Running phase
func (r *StatefulChaosExperimentReconciler) processRunningExperiment(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) (ctrl.Result, error) {
	log.Info("Processing running experiment")

	// Check if safety conditions require stopping the experiment
	if experiment.Spec.Safety != nil && experiment.Spec.Safety.AutoRollback {
		if shouldRollback, reason := r.shouldRollbackExperiment(ctx, experiment); shouldRollback {
			log.Info("Auto-rollback triggered", "reason", reason)
			if err := r.removeChaos(ctx, experiment, log); err != nil {
				log.Error(err, "Failed to remove chaos during rollback")
			}
			experiment.Status.Phase = "Failed"
			experiment.Status.FailureReason = fmt.Sprintf("Auto-rollback triggered: %s", reason)
			experiment.Status.EndTime = &metav1.Time{Time: time.Now()}
			if err := r.Status().Update(ctx, experiment); err != nil {
				log.Error(err, "Failed to update experiment status")
			}
			return ctrl.Result{}, nil
		}
	}

	// Check if experiment has timed out
	if experiment.Status.StartTime != nil {
		duration, err := time.ParseDuration(experiment.Spec.Duration)
		if err != nil {
			log.Error(err, "Invalid duration format", "duration", experiment.Spec.Duration)
			duration = 5 * time.Minute // Default to 5 minutes if invalid
		}

		if time.Since(experiment.Status.StartTime.Time) > duration {
			log.Info("Experiment duration reached, cleaning up")
			// Remove chaos
			if err := r.removeChaos(ctx, experiment, log); err != nil {
				log.Error(err, "Failed to remove chaos")
				return ctrl.Result{}, err
			}

			// Update status
			experiment.Status.Phase = "Completed"
			experiment.Status.EndTime = &metav1.Time{Time: time.Now()}

			condition := metav1.Condition{
				Type:               "ChaosRemoved",
				Status:             metav1.ConditionTrue,
				LastTransitionTime: metav1.Now(),
				Reason:             "ExperimentCompleted",
				Message:            "Chaos experiment completed successfully",
			}
			experiment.Status.Conditions = append(experiment.Status.Conditions, condition)

			if err := r.Status().Update(ctx, experiment); err != nil {
				log.Error(err, "Failed to update experiment status")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}
	}

	// Continue monitoring
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

// handleExperimentDeletion processes the deletion of an experiment
func (r *StatefulChaosExperimentReconciler) handleExperimentDeletion(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) (ctrl.Result, error) {
	log.Info("Handling experiment deletion")

	// Remove chaos if experiment is still running
	if experiment.Status.Phase == "Running" {
		if err := r.removeChaos(ctx, experiment, log); err != nil {
			log.Error(err, "Failed to remove chaos during deletion")
			return ctrl.Result{}, err
		}
	}

	// Remove finalizer
	experiment.Finalizers = removeString(experiment.Finalizers, chaosExperimentFinalizer)
	if err := r.Update(ctx, experiment); err != nil {
		log.Error(err, "Failed to remove finalizer")
		return ctrl.Result{}, err
	}

	log.Info("Experiment cleanup completed")
	return ctrl.Result{}, nil
}

// identifyTargetResources finds the resources that match the target selector
func (r *StatefulChaosExperimentReconciler) identifyTargetResources(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment) ([]chaosv1alpha1.TargetResourceStatus, error) {
	var targetResources []chaosv1alpha1.TargetResourceStatus

	// Determine target namespace
	namespace := experiment.Spec.Target.Namespace
	if namespace == "" {
		namespace = experiment.Namespace
	}

	// Target by name
	if experiment.Spec.Target.Name != "" {
		targetType := experiment.Spec.Target.TargetType
		if targetType == "" {
			targetType = "StatefulSet" // Default to StatefulSet if not specified
		}

		switch targetType {
		case "StatefulSet":
			sts := &appsv1.StatefulSet{}
			err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: experiment.Spec.Target.Name}, sts)
			if err != nil {
				return nil, fmt.Errorf("failed to get StatefulSet %s: %w", experiment.Spec.Target.Name, err)
			}
			targetResources = append(targetResources, chaosv1alpha1.TargetResourceStatus{
				Kind:      "StatefulSet",
				Name:      sts.Name,
				Namespace: sts.Namespace,
				UID:       string(sts.UID),
				Status:    "Targeted",
			})

		case "Pod":
			pod := &corev1.Pod{}
			err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: experiment.Spec.Target.Name}, pod)
			if err != nil {
				return nil, fmt.Errorf("failed to get Pod %s: %w", experiment.Spec.Target.Name, err)
			}
			targetResources = append(targetResources, chaosv1alpha1.TargetResourceStatus{
				Kind:      "Pod",
				Name:      pod.Name,
				Namespace: pod.Namespace,
				UID:       string(pod.UID),
				Status:    "Targeted",
			})

		case "PersistentVolumeClaim":
			pvc := &corev1.PersistentVolumeClaim{}
			err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: experiment.Spec.Target.Name}, pvc)
			if err != nil {
				return nil, fmt.Errorf("failed to get PVC %s: %w", experiment.Spec.Target.Name, err)
			}
			targetResources = append(targetResources, chaosv1alpha1.TargetResourceStatus{
				Kind:      "PersistentVolumeClaim",
				Name:      pvc.Name,
				Namespace: pvc.Namespace,
				UID:       string(pvc.UID),
				Status:    "Targeted",
			})

		default:
			return nil, fmt.Errorf("unsupported target type %s", targetType)
		}

		return targetResources, nil
	}

	// Target by selector
	if experiment.Spec.Target.Selector != nil {
		labelSelector, err := metav1.LabelSelectorAsSelector(experiment.Spec.Target.Selector)
		if err != nil {
			return nil, fmt.Errorf("invalid label selector: %w", err)
		}

		targetType := experiment.Spec.Target.TargetType
		if targetType == "" {
			targetType = "StatefulSet" // Default to StatefulSet if not specified
		}

		switch targetType {
		case "StatefulSet":
			stsList := &appsv1.StatefulSetList{}
			if err := r.List(ctx, stsList, client.InNamespace(namespace), client.MatchingLabelsSelector{Selector: labelSelector}); err != nil {
				return nil, fmt.Errorf("failed to list StatefulSets: %w", err)
			}
			for _, sts := range stsList.Items {
				targetResources = append(targetResources, chaosv1alpha1.TargetResourceStatus{
					Kind:      "StatefulSet",
					Name:      sts.Name,
					Namespace: sts.Namespace,
					UID:       string(sts.UID),
					Status:    "Targeted",
				})
			}

		case "Pod":
			podList := &corev1.PodList{}
			if err := r.List(ctx, podList, client.InNamespace(namespace), client.MatchingLabelsSelector{Selector: labelSelector}); err != nil {
				return nil, fmt.Errorf("failed to list Pods: %w", err)
			}
			for _, pod := range podList.Items {
				targetResources = append(targetResources, chaosv1alpha1.TargetResourceStatus{
					Kind:      "Pod",
					Name:      pod.Name,
					Namespace: pod.Namespace,
					UID:       string(pod.UID),
					Status:    "Targeted",
				})
			}

		default:
			return nil, fmt.Errorf("unsupported target type %s for selector", targetType)
		}
	}

	// Handle the Mode/Value parameters (all, one, random, percentage)
	if experiment.Spec.Target.Mode != "" {
		switch experiment.Spec.Target.Mode {
		case "One":
			if len(targetResources) > 0 {
				// Just keep the first resource
				targetResources = targetResources[:1]
			}
		case "All":
			// Already have all resources, nothing to do
		case "Random":
			// Implementation would normally select a random resource
			// For simplicity, just take the first one here
			if len(targetResources) > 0 {
				targetResources = targetResources[:1]
			}
		case "Percentage", "Fixed":
			// Implementation would filter resources based on percentage/count
			// For simplicity, not fully implemented here
		}
	}

	if len(targetResources) == 0 {
		return nil, fmt.Errorf("no matching resources found for target criteria")
	}

	return targetResources, nil
}

// injectChaos applies the specified chaos to target resources
func (r *StatefulChaosExperimentReconciler) injectChaos(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	chaosType := experiment.Spec.ChaosType
	log.Info("Injecting chaos", "type", chaosType)

	// Apply different chaos types
	switch chaosType {
	case "DiskFailure":
		return r.injectDiskFailure(ctx, experiment, log)
	case "NetworkLatency":
		return r.injectNetworkLatency(ctx, experiment, log)
	case "PodFailure":
		return r.injectPodFailure(ctx, experiment, log)
	case "StatefulSetScaling":
		return r.injectStatefulSetScaling(ctx, experiment, log)
	default:
		return fmt.Errorf("unsupported chaos type: %s", chaosType)
	}
}

// Placeholder implementations of chaos injection methods
// These would be implemented with actual chaos injection logic in a real project

func (r *StatefulChaosExperimentReconciler) injectDiskFailure(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Simulating disk failure")
	// Implementation would inject I/O errors or latency into storage
	return nil
}

func (r *StatefulChaosExperimentReconciler) injectNetworkLatency(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Simulating network latency")
	// Implementation would add network latency/packet loss
	return nil
}

func (r *StatefulChaosExperimentReconciler) injectPodFailure(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Simulating pod failure")

	// Get the list of target pods
	for _, target := range experiment.Status.TargetResources {
		if target.Kind == "Pod" {
			pod := &corev1.Pod{}
			err := r.Get(ctx, types.NamespacedName{Namespace: target.Namespace, Name: target.Name}, pod)
			if err != nil {
				log.Error(err, "Failed to get target pod", "pod", target.Name)
				continue
			}

			// Delete the pod to simulate failure
			if err := r.Delete(ctx, pod); err != nil {
				log.Error(err, "Failed to delete pod", "pod", target.Name)
				return err
			}
			log.Info("Successfully deleted pod", "pod", target.Name)
		}
	}

	return nil
}

func (r *StatefulChaosExperimentReconciler) injectStatefulSetScaling(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Simulating StatefulSet scaling issues")

	// Get the list of target statefulsets
	for _, target := range experiment.Status.TargetResources {
		if target.Kind == "StatefulSet" {
			sts := &appsv1.StatefulSet{}
			err := r.Get(ctx, types.NamespacedName{Namespace: target.Namespace, Name: target.Name}, sts)
			if err != nil {
				log.Error(err, "Failed to get target StatefulSet", "statefulset", target.Name)
				continue
			}

			// Store the original replica count for later restoration
			originalReplicas := *sts.Spec.Replicas

			// Scale down the StatefulSet
			desiredReplicas := int32(1)
			if originalReplicas > 1 {
				desiredReplicas = originalReplicas - 1
			}

			sts.Spec.Replicas = &desiredReplicas
			if err := r.Update(ctx, sts); err != nil {
				log.Error(err, "Failed to scale down StatefulSet", "statefulset", target.Name)
				return err
			}
			log.Info("Successfully scaled down StatefulSet", "statefulset", target.Name,
				"originalReplicas", originalReplicas, "newReplicas", desiredReplicas)

			// Store original replicas in annotations for later restoration
			if sts.Annotations == nil {
				sts.Annotations = make(map[string]string)
			}
			sts.Annotations["statefulchaos.io/original-replicas"] = fmt.Sprintf("%d", originalReplicas)
			if err := r.Update(ctx, sts); err != nil {
				log.Error(err, "Failed to update StatefulSet annotations", "statefulset", target.Name)
				return err
			}
		}
	}

	return nil
}

// removeChaos cleans up any chaos injected by the experiment
func (r *StatefulChaosExperimentReconciler) removeChaos(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	chaosType := experiment.Spec.ChaosType
	log.Info("Removing chaos", "type", chaosType)

	// Remove different chaos types
	switch chaosType {
	case "DiskFailure":
		return r.removeDiskFailure(ctx, experiment, log)
	case "NetworkLatency":
		return r.removeNetworkLatency(ctx, experiment, log)
	case "PodFailure":
		// Nothing to do, Kubernetes will recreate pods
		return nil
	case "StatefulSetScaling":
		return r.removeStatefulSetScaling(ctx, experiment, log)
	default:
		return fmt.Errorf("unsupported chaos type: %s", chaosType)
	}
}

// Placeholder implementations of chaos removal methods
func (r *StatefulChaosExperimentReconciler) removeDiskFailure(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Removing disk failure simulation")
	// Implementation would remove I/O errors or latency from storage
	return nil
}

func (r *StatefulChaosExperimentReconciler) removeNetworkLatency(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Removing network latency simulation")
	// Implementation would remove network latency/packet loss
	return nil
}

func (r *StatefulChaosExperimentReconciler) removeStatefulSetScaling(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error {
	log.Info("Removing StatefulSet scaling issues")

	// Get the list of target statefulsets
	for _, target := range experiment.Status.TargetResources {
		if target.Kind == "StatefulSet" {
			sts := &appsv1.StatefulSet{}
			err := r.Get(ctx, types.NamespacedName{Namespace: target.Namespace, Name: target.Name}, sts)
			if err != nil {
				log.Error(err, "Failed to get target StatefulSet", "statefulset", target.Name)
				continue
			}

			// Get the original replica count from annotations
			if originalReplicasStr, ok := sts.Annotations["statefulchaos.io/original-replicas"]; ok {
				var originalReplicas int32
				_, err := fmt.Sscanf(originalReplicasStr, "%d", &originalReplicas)
				if err != nil {
					log.Error(err, "Failed to parse original replicas", "value", originalReplicasStr)
					continue
				}

				// Restore the original replica count
				sts.Spec.Replicas = &originalReplicas
				if err := r.Update(ctx, sts); err != nil {
					log.Error(err, "Failed to restore StatefulSet replicas", "statefulset", target.Name)
					return err
				}
				log.Info("Successfully restored StatefulSet replicas", "statefulset", target.Name, "replicas", originalReplicas)

				// Remove our annotation
				delete(sts.Annotations, "statefulchaos.io/original-replicas")
				if err := r.Update(ctx, sts); err != nil {
					log.Error(err, "Failed to update StatefulSet annotations", "statefulset", target.Name)
					return err
				}
			}
		}
	}

	return nil
}

// shouldRollbackExperiment checks if experiment should be rolled back based on safety conditions
func (r *StatefulChaosExperimentReconciler) shouldRollbackExperiment(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment) (bool, string) {
	if experiment.Spec.Safety == nil || len(experiment.Spec.Safety.HealthChecks) == 0 {
		return false, ""
	}

	// Implement health check logic here
	// For example, checking HTTP endpoints, metrics, etc.

	// For now, we'll return a simple "no rollback needed"
	return false, ""
}

// SetupWithManager sets up the controller with the Manager
func (r *StatefulChaosExperimentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&chaosv1alpha1.StatefulChaosExperiment{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}

// Helper constants and functions
const chaosExperimentFinalizer = "chaos.statefulchaos.io/finalizer"

// Helper functions for finalizer management
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) []string {
	var result []string
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}
