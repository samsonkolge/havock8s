package chaos

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestPodFailureInjector_Inject(t *testing.T) {
	tests := []struct {
		name       string
		experiment *chaosv1alpha1.Havock8sExperiment
		pod        *corev1.Pod
		wantErr    bool
		wantPhase  corev1.PodPhase
	}{
		{
			name: "successful pod failure injection",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Parameters: map[string]string{
						"failureMode": "crash",
					},
				},
				Status: chaosv1alpha1.Havock8sExperimentStatus{
					TargetResources: []chaosv1alpha1.TargetResourceStatus{
						{
							Kind:      "Pod",
							Name:      "test-pod",
							Namespace: "default",
						},
					},
				},
			},
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
				},
			},
			wantErr:   false,
			wantPhase: corev1.PodFailed,
		},
		{
			name: "missing pod",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Parameters: map[string]string{
						"failureMode": "crash",
					},
				},
				Status: chaosv1alpha1.Havock8sExperimentStatus{
					TargetResources: []chaosv1alpha1.TargetResourceStatus{
						{
							Kind:      "Pod",
							Name:      "non-existent-pod",
							Namespace: "default",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid failure mode",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Parameters: map[string]string{
						"failureMode": "invalid",
					},
				},
				Status: chaosv1alpha1.Havock8sExperimentStatus{
					TargetResources: []chaosv1alpha1.TargetResourceStatus{
						{
							Kind:      "Pod",
							Name:      "test-pod",
							Namespace: "default",
						},
					},
				},
			},
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = corev1.AddToScheme(scheme)
			_ = chaosv1alpha1.AddToScheme(scheme)

			objs := []client.Object{tt.experiment}
			if tt.pod != nil {
				objs = append(objs, tt.pod)
			}

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(objs...).
				Build()

			injector := &PodFailureInjector{}
			injector.SetClient(fakeClient)

			err := injector.Inject(context.Background(), tt.experiment, logr.Discard())
			if (err != nil) != tt.wantErr {
				t.Errorf("PodFailureInjector.Inject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.pod != nil {
				updatedPod := &corev1.Pod{}
				err := fakeClient.Get(context.Background(), client.ObjectKey{
					Namespace: tt.pod.Namespace,
					Name:      tt.pod.Name,
				}, updatedPod)

				// For any failure mode, we expect the pod to be deleted
				if client.IgnoreNotFound(err) != nil {
					t.Errorf("Unexpected error getting pod: %v", err)
				} else if err == nil {
					t.Error("Pod should have been deleted")
				}
			}
		})
	}
}

func TestPodFailureInjector_Cleanup(t *testing.T) {
	tests := []struct {
		name       string
		experiment *chaosv1alpha1.Havock8sExperiment
		pod        *corev1.Pod
		wantErr    bool
	}{
		{
			name: "successful cleanup",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Status: chaosv1alpha1.Havock8sExperimentStatus{
					TargetResources: []chaosv1alpha1.TargetResourceStatus{
						{
							Kind:      "Pod",
							Name:      "test-pod",
							Namespace: "default",
						},
					},
				},
			},
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						"havock8s.io/pod-failure": "true",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = corev1.AddToScheme(scheme)
			_ = chaosv1alpha1.AddToScheme(scheme)

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tt.experiment, tt.pod).
				Build()

			injector := &PodFailureInjector{}
			injector.SetClient(fakeClient)

			err := injector.Cleanup(context.Background(), tt.experiment, logr.Discard())
			if (err != nil) != tt.wantErr {
				t.Errorf("PodFailureInjector.Cleanup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				updatedPod := &corev1.Pod{}
				err := fakeClient.Get(context.Background(), client.ObjectKey{
					Namespace: tt.pod.Namespace,
					Name:      tt.pod.Name,
				}, updatedPod)
				if err != nil {
					t.Errorf("Failed to get updated pod: %v", err)
					return
				}

				if _, ok := updatedPod.Annotations["havock8s.io/pod-failure"]; ok {
					t.Error("Pod failure annotation still present after cleanup")
				}
			}
		})
	}
}
