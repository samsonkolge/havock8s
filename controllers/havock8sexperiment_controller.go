package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// havock8sExperimentReconciler reconciles a havock8sExperiment object
type havock8sExperimentReconciler struct {
	client.Client
	Log    logr.Logger
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

// Reconcile handles the main reconciliation loop for havock8sExperiment resources
func (r *havock8sExperimentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("havock8sexperiment", req.NamespacedName)

	// Fetch the havock8sExperiment instance
	experiment := &chaosv1alpha1.havock8sExperiment{}
	if err := r.Get(ctx, req.NamespacedName, experiment); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted
			return ctrl.Result{}, nil
		}
		// Error reading the object
		log.Error(err, "Failed to get havock8sExperiment")
		return ctrl.Result{}, err
	}

	// Initialize finalizer
	finalizerName := "chaos.havock8s.io/finalizer"

	// Examine if the object is being deleted
	if experiment.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then add the finalizer
		if !containsString(experiment.ObjectMeta.Finalizers, finalizerName) {
			experiment.ObjectMeta.Finalizers = append(experiment.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, experiment); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(experiment.ObjectMeta.Finalizers, finalizerName) {
			// Handle experiment deletion
			return r.handleExperimentDeletion(ctx, experiment, log)
		}
		// Our finalizer has finished, so the reconciler loop can stop
		return ctrl.Result{}, nil
	}

	// Main reconciliation logic
	switch experiment.Status.Phase {
	case "":
		// Initialize the experiment
		return r.initializeExperiment(ctx, experiment, log)
	case "Pending":
		// Process pending experiment
		return r.processPendingExperiment(ctx, experiment, log)
	case "Running":
		// Process running experiment
		return r.processRunningExperiment(ctx, experiment, log)
	case "Completed":
		// Already completed, nothing to do
		return ctrl.Result{}, nil
	case "Failed":
		// Already failed, nothing to do
		return ctrl.Result{}, nil
	default:
		log.Info("Unknown experiment phase", "phase", experiment.Status.Phase)
		return ctrl.Result{}, nil
	}
}

// Functions like initializeExperiment, processPendingExperiment, etc. would follow.
// They would be similar to the previous implementation but with updated types.

// SetupWithManager sets up the controller with the Manager.
func (r *havock8sExperimentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&chaosv1alpha1.havock8sExperiment{}).
		Complete(r)
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
