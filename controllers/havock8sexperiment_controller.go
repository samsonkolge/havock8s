package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
	"github.com/havock8s/havock8s/pkg/chaos"
	"github.com/havock8s/havock8s/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Havock8sExperimentReconciler reconciles a Havock8sExperiment object
type Havock8sExperimentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=chaos.havock8s.io,resources=havock8sexperiments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=chaos.havock8s.io,resources=havock8sexperiments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=chaos.havock8s.io,resources=havock8sexperiments/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;delete
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=persistentvolumes,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch

// Reconcile handles the reconciliation of Havock8sExperiment resources
func (r *Havock8sExperimentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Havock8sExperiment instance
	experiment := &chaosv1alpha1.Havock8sExperiment{}
	if err := r.Get(ctx, req.NamespacedName, experiment); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle experiment deletion
	if !experiment.DeletionTimestamp.IsZero() {
		return r.handleExperimentDeletion(ctx, experiment, logger)
	}

	// Process experiment based on its phase
	switch experiment.Status.Phase {
	case "":
		// Initialize new experiment
		return r.initializeExperiment(ctx, experiment, logger)
	case "Pending":
		// Process pending experiment
		return r.processPendingExperiment(ctx, experiment, logger)
	case "Running":
		// Process running experiment
		return r.processRunningExperiment(ctx, experiment, logger)
	default:
		// No action needed for completed or failed experiments
		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager
func (r *Havock8sExperimentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&chaosv1alpha1.Havock8sExperiment{}).
		Complete(r)
}

// handleExperimentDeletion handles cleanup when an experiment is being deleted
func (r *Havock8sExperimentReconciler) handleExperimentDeletion(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, logger logr.Logger) (ctrl.Result, error) {
	// Perform cleanup
	injector, err := chaos.GetInjector(experiment.Spec.ChaosType)
	if err != nil {
		return ctrl.Result{}, err
	}

	injector.SetClient(r.Client)
	if err := injector.Cleanup(ctx, experiment, logger); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// initializeExperiment initializes a new experiment
func (r *Havock8sExperimentReconciler) initializeExperiment(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, logger logr.Logger) (ctrl.Result, error) {
	// Set initial phase
	experiment.Status.Phase = "Pending"
	experiment.Status.StartTime = &metav1.Time{Time: time.Now()}

	if err := r.Status().Update(ctx, experiment); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

// processPendingExperiment processes an experiment in the Pending phase
func (r *Havock8sExperimentReconciler) processPendingExperiment(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, logger logr.Logger) (ctrl.Result, error) {
	// Check if target exists
	targetExists, err := utils.CheckTargetExists(ctx, r.Client, experiment.Spec.Target)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !targetExists {
		experiment.Status.Phase = "Failed"
		experiment.Status.FailureReason = "Target resource not found"
		if err := r.Status().Update(ctx, experiment); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, fmt.Errorf("target resource not found")
	}

	// Check safety conditions
	safetyChecker := utils.NewSafetyChecker(r.Client)
	if shouldRollback, reason := safetyChecker.CheckSafety(ctx, experiment, logger); shouldRollback {
		experiment.Status.Phase = "Failed"
		experiment.Status.FailureReason = reason
		if err := r.Status().Update(ctx, experiment); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, fmt.Errorf(reason)
	}

	// Set up target resources in status
	if experiment.Spec.Target.TargetType == "Pod" {
		// Get the pod to get its UID
		pod := &corev1.Pod{}
		err := r.Get(ctx, types.NamespacedName{
			Namespace: experiment.Spec.Target.Namespace,
			Name:      experiment.Spec.Target.Name,
		}, pod)
		if err != nil {
			return ctrl.Result{}, err
		}

		experiment.Status.TargetResources = []chaosv1alpha1.TargetResourceStatus{
			{
				Kind:      "Pod",
				Name:      pod.Name,
				Namespace: pod.Namespace,
				UID:       string(pod.UID),
				Status:    "Targeted",
			},
		}
		if err := r.Status().Update(ctx, experiment); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Start chaos injection
	injector, err := chaos.GetInjector(experiment.Spec.ChaosType)
	if err != nil {
		experiment.Status.Phase = "Failed"
		experiment.Status.FailureReason = err.Error()
		if err := r.Status().Update(ctx, experiment); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	injector.SetClient(r.Client)
	if err := injector.Inject(ctx, experiment, logger); err != nil {
		experiment.Status.Phase = "Failed"
		experiment.Status.FailureReason = err.Error()
		if err := r.Status().Update(ctx, experiment); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	// Update status to Running
	experiment.Status.Phase = "Running"
	if err := r.Status().Update(ctx, experiment); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Second * 30}, nil
}

// processRunningExperiment processes an experiment in the Running phase
func (r *Havock8sExperimentReconciler) processRunningExperiment(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, logger logr.Logger) (ctrl.Result, error) {
	// Check if experiment duration has elapsed
	duration, err := time.ParseDuration(experiment.Spec.Duration)
	if err != nil {
		return ctrl.Result{}, err
	}

	// If StartTime is nil, set it to now
	if experiment.Status.StartTime == nil {
		experiment.Status.StartTime = &metav1.Time{Time: time.Now()}
		if err := r.Status().Update(ctx, experiment); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: time.Second * 30}, nil
	}

	if time.Since(experiment.Status.StartTime.Time) > duration {
		// Clean up chaos
		injector, err := chaos.GetInjector(experiment.Spec.ChaosType)
		if err != nil {
			return ctrl.Result{}, err
		}

		injector.SetClient(r.Client)
		if err := injector.Cleanup(ctx, experiment, logger); err != nil {
			return ctrl.Result{}, err
		}

		// Update status to Completed
		experiment.Status.Phase = "Completed"
		experiment.Status.EndTime = &metav1.Time{Time: time.Now()}
		if err := r.Status().Update(ctx, experiment); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// Check if target still exists
	targetExists, err := utils.CheckTargetExists(ctx, r.Client, experiment.Spec.Target)
	if err != nil {
		return ctrl.Result{}, err
	}

	// If target doesn't exist and it's a pod failure experiment, this is expected
	if !targetExists && experiment.Spec.ChaosType == "pod-failure" {
		// Update status to Completed since the pod has been successfully deleted
		experiment.Status.Phase = "Completed"
		experiment.Status.EndTime = &metav1.Time{Time: time.Now()}
		if err := r.Status().Update(ctx, experiment); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Continue monitoring
	return ctrl.Result{RequeueAfter: time.Second * 30}, nil
}

// Helper functions
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}

// cleanupExperiment performs cleanup for an experiment
func (r *Havock8sExperimentReconciler) cleanupExperiment(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, logger logr.Logger) error {
	injector, err := chaos.GetInjector(experiment.Spec.ChaosType)
	if err != nil {
		return err
	}

	injector.SetClient(r.Client)
	if err := injector.Cleanup(ctx, experiment, logger); err != nil {
		return err
	}

	experiment.Status.Phase = "Completed"
	experiment.Status.EndTime = &metav1.Time{Time: time.Now()}
	if err := r.Status().Update(ctx, experiment); err != nil {
		return err
	}

	return nil
}
