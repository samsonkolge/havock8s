package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
	"github.com/havock8s/havock8s/pkg/chaos"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func setupScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = chaosv1alpha1.AddToScheme(scheme)
	return scheme
}

func setupFakeClient(scheme *runtime.Scheme) client.Client {
	return fake.NewClientBuilder().
		WithScheme(scheme).
		WithStatusSubresource(&chaosv1alpha1.Havock8sExperiment{}).
		Build()
}

func setupTestPod(client client.Client, name, namespace string) error {
	// First try to delete any existing pod
	existingPod := &corev1.Pod{}
	err := client.Get(context.Background(), types.NamespacedName{Name: name, Namespace: namespace}, existingPod)
	if err == nil {
		// Pod exists, delete it
		if err := client.Delete(context.Background(), existingPod); err != nil {
			return fmt.Errorf("failed to delete existing pod: %v", err)
		}
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Annotations: map[string]string{
				"havock8s.io/protected": "false",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "test",
					Image: "busybox",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 8080,
						},
					},
				},
			},
		},
	}

	// Create the pod
	if err := client.Create(context.Background(), pod); err != nil {
		return fmt.Errorf("failed to create pod: %v", err)
	}

	// Get the created pod
	createdPod := &corev1.Pod{}
	if err := client.Get(context.Background(), types.NamespacedName{Name: name, Namespace: namespace}, createdPod); err != nil {
		return fmt.Errorf("failed to get created pod: %v", err)
	}

	// Update the pod status
	createdPod.Status = corev1.PodStatus{
		Phase: corev1.PodRunning,
		Conditions: []corev1.PodCondition{
			{
				Type:               corev1.PodReady,
				Status:             corev1.ConditionTrue,
				LastTransitionTime: metav1.Now(),
			},
		},
		ContainerStatuses: []corev1.ContainerStatus{
			{
				Name:  "test",
				Ready: true,
				State: corev1.ContainerState{
					Running: &corev1.ContainerStateRunning{
						StartedAt: metav1.Now(),
					},
				},
			},
		},
	}

	// Update the pod status using the status client
	if err := client.Status().Update(context.Background(), createdPod); err != nil {
		return fmt.Errorf("failed to update pod status: %v", err)
	}

	// Verify the pod status was updated
	updatedPod := &corev1.Pod{}
	if err := client.Get(context.Background(), types.NamespacedName{Name: name, Namespace: namespace}, updatedPod); err != nil {
		return fmt.Errorf("failed to verify pod status update: %v", err)
	}

	if updatedPod.Status.Phase != corev1.PodRunning {
		return fmt.Errorf("pod status not updated correctly: expected Running, got %s", updatedPod.Status.Phase)
	}

	return nil
}

func cleanupTestResources(client client.Client) error {
	// Delete all test pods
	podList := &corev1.PodList{}
	if err := client.List(context.Background(), podList); err != nil {
		return fmt.Errorf("failed to list pods: %v", err)
	}
	for _, pod := range podList.Items {
		if err := client.Delete(context.Background(), &pod); err != nil {
			return fmt.Errorf("failed to delete pod %s: %v", pod.Name, err)
		}
	}

	// Delete all experiments
	expList := &chaosv1alpha1.Havock8sExperimentList{}
	if err := client.List(context.Background(), expList); err != nil {
		return fmt.Errorf("failed to list experiments: %v", err)
	}
	for _, exp := range expList.Items {
		if err := client.Delete(context.Background(), &exp); err != nil {
			return fmt.Errorf("failed to delete experiment %s: %v", exp.Name, err)
		}
	}
	return nil
}

func reconcileAndWait(t *testing.T, reconciler *Havock8sExperimentReconciler, req reconcile.Request, maxAttempts int, interval time.Duration) error {
	for i := 0; i < maxAttempts; i++ {
		result, err := reconciler.Reconcile(context.Background(), req)
		if err != nil {
			return err
		}
		if !result.Requeue && result.RequeueAfter == 0 {
			return nil
		}
		time.Sleep(interval)
	}
	return fmt.Errorf("max reconciliation attempts reached")
}

func TestHavock8sExperimentReconciler_Reconcile(t *testing.T) {
	// Start a mock HTTP server for health checks
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Parse the server URL to get the port
	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Failed to parse server URL: %v", err)
	}
	port, err := strconv.ParseInt(u.Port(), 10, 32)
	if err != nil {
		t.Fatalf("Failed to parse port: %v", err)
	}

	tests := []struct {
		name       string
		experiment *chaosv1alpha1.Havock8sExperiment
		wantErr    bool
		wantPhases []string
	}{
		{
			name: "successful_experiment_execution",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-experiment-1",
					Namespace: "default",
				},
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Target: chaosv1alpha1.TargetSpec{
						Name:       "test-pod",
						Namespace:  "default",
						TargetType: "Pod",
					},
					ChaosType: "pod-failure",
					Duration:  "1s",
					Intensity: 100,
					Parameters: map[string]string{
						"failureMode": "terminate",
					},
					Safety: &chaosv1alpha1.SafetySpec{
						HealthChecks: []chaosv1alpha1.HealthCheckSpec{
							{
								Type: "httpGet",
								Path: "/health",
								Port: int32(port),
							},
						},
					},
				},
				Status: chaosv1alpha1.Havock8sExperimentStatus{
					Phase: "Pending",
					TargetResources: []chaosv1alpha1.TargetResourceStatus{
						{
							Kind:      "Pod",
							Name:      "test-pod",
							Namespace: "default",
							Status:    "Targeted",
						},
					},
				},
			},
			wantErr:    false,
			wantPhases: []string{"Pending", "Running", "Completed"},
		},
		{
			name: "missing_target_pod",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-experiment-2",
					Namespace: "default",
				},
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Target: chaosv1alpha1.TargetSpec{
						Name:       "non-existent-pod",
						Namespace:  "default",
						TargetType: "Pod",
					},
					ChaosType: "pod-failure",
					Duration:  "10s",
					Intensity: 100,
					Parameters: map[string]string{
						"failureMode": "terminate",
					},
				},
				Status: chaosv1alpha1.Havock8sExperimentStatus{
					Phase: "Pending",
					TargetResources: []chaosv1alpha1.TargetResourceStatus{
						{
							Kind:      "Pod",
							Name:      "non-existent-pod",
							Namespace: "default",
							Status:    "Targeted",
						},
					},
				},
			},
			wantErr:    true,
			wantPhases: []string{"Pending", "Failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := setupScheme()
			fakeClient := setupFakeClient(scheme)
			
			// Clean up any existing resources
			if err := cleanupTestResources(fakeClient); err != nil {
				t.Fatalf("Failed to cleanup test resources: %v", err)
			}

			// Create the test pod if target name is specified
			if tt.experiment.Spec.Target.Name != "" && tt.experiment.Spec.Target.Name != "non-existent-pod" {
				if err := setupTestPod(fakeClient, tt.experiment.Spec.Target.Name, tt.experiment.Spec.Target.Namespace); err != nil {
					t.Fatalf("Failed to create test pod: %v", err)
				}

				// Verify the pod is ready
				pod := &corev1.Pod{}
				if err := fakeClient.Get(context.Background(), types.NamespacedName{
					Name:      tt.experiment.Spec.Target.Name,
					Namespace: tt.experiment.Spec.Target.Namespace,
				}, pod); err != nil {
					t.Fatalf("Failed to verify pod creation: %v", err)
				}

				if pod.Status.Phase != corev1.PodRunning {
					t.Fatalf("Pod is not in Running phase: %s", pod.Status.Phase)
				}

				isReady := false
				for _, cond := range pod.Status.Conditions {
					if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
						isReady = true
						break
					}
				}
				if !isReady {
					t.Fatalf("Pod is not ready")
				}

				t.Logf("Test pod created and ready: %s/%s", pod.Namespace, pod.Name)
			}

			// Create the experiment
			if err := fakeClient.Create(context.Background(), tt.experiment); err != nil {
				t.Fatalf("Failed to create experiment: %v", err)
			}

			// Register the pod failure injector with lowercase name
			chaos.RegisterInjector("pod-failure", &chaos.PodFailureInjector{})

			reconciler := &Havock8sExperimentReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}

			// Run reconciliation with retries
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      tt.experiment.Name,
					Namespace: tt.experiment.Namespace,
				},
			}

			err := reconcileAndWait(t, reconciler, req, 10, 200*time.Millisecond)
			if (err != nil) != tt.wantErr {
				t.Errorf("Havock8sExperimentReconciler.Reconcile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Get the updated experiment
			updatedExp := &chaosv1alpha1.Havock8sExperiment{}
			if err := fakeClient.Get(context.Background(), types.NamespacedName{
				Name:      tt.experiment.Name,
				Namespace: tt.experiment.Namespace,
			}, updatedExp); err != nil {
				t.Fatalf("Failed to get updated experiment: %v", err)
			}

			// Check if the phase is as expected
			if !tt.wantErr && updatedExp.Status.Phase != tt.wantPhases[len(tt.wantPhases)-1] {
				t.Errorf("Expected phase %s, got %s", tt.wantPhases[len(tt.wantPhases)-1], updatedExp.Status.Phase)
			}

			// For missing target pod test, verify that the experiment is marked as failed
			if tt.name == "missing_target_pod" {
				if updatedExp.Status.Phase != "Failed" {
					t.Errorf("Expected phase Failed for missing target pod, got %s", updatedExp.Status.Phase)
				}
			}

			// For successful experiment, verify the pod was deleted
			if tt.name == "successful_experiment_execution" {
				pod := &corev1.Pod{}
				err := fakeClient.Get(context.Background(), types.NamespacedName{
					Name:      tt.experiment.Spec.Target.Name,
					Namespace: tt.experiment.Spec.Target.Namespace,
				}, pod)
				if err == nil {
					t.Error("Expected pod to be deleted, but it still exists")
				}
			}
		})
	}
}