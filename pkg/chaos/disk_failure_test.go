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

func TestDiskFailureInjector_Inject(t *testing.T) {
	tests := []struct {
		name       string
		experiment *chaosv1alpha1.Havock8sExperiment
		pod        *corev1.Pod
		wantErr    bool
		wantMount  string
		wantMode   string
	}{
		{
			name: "successful disk failure injection",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Parameters: map[string]string{
						"mountPath": "/data",
						"failureMode": "readonly",
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
			wantErr:   false,
			wantMount: "/data",
			wantMode:  "readonly",
		},
		{
			name: "missing pod",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Parameters: map[string]string{
						"mountPath": "/data",
						"failureMode": "readonly",
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
			wantErr: true,
		},
		{
			name: "invalid failure mode",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Parameters: map[string]string{
						"mountPath": "/data",
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
			wantErr: true,
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

			injector := &DiskFailureInjector{}
			injector.SetClient(fakeClient)

			err := injector.Inject(context.Background(), tt.experiment, logr.Discard())
			if (err != nil) != tt.wantErr {
				t.Errorf("DiskFailureInjector.Inject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.pod != nil {
				updatedPod := &corev1.Pod{}
				err := fakeClient.Get(context.Background(), client.ObjectKey{
					Namespace: tt.pod.Namespace,
					Name:      tt.pod.Name,
				}, updatedPod)
				if err != nil {
					t.Errorf("Failed to get updated pod: %v", err)
					return
				}

				if updatedPod.Annotations["havock8s.io/disk-failure"] != "true" {
					t.Error("Disk failure annotation not set")
				}
				if updatedPod.Annotations["havock8s.io/disk-failure-mount"] != tt.wantMount {
					t.Errorf("Mount path = %v, want %v", 
						updatedPod.Annotations["havock8s.io/disk-failure-mount"], 
						tt.wantMount)
				}
				if updatedPod.Annotations["havock8s.io/disk-failure-mode"] != tt.wantMode {
					t.Errorf("Failure mode = %v, want %v", 
						updatedPod.Annotations["havock8s.io/disk-failure-mode"], 
						tt.wantMode)
				}
			}
		})
	}
}

func TestDiskFailureInjector_Cleanup(t *testing.T) {
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
						"havock8s.io/disk-failure":        "true",
						"havock8s.io/disk-failure-mount":  "/data",
						"havock8s.io/disk-failure-mode":   "readonly",
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

			injector := &DiskFailureInjector{}
			injector.SetClient(fakeClient)

			err := injector.Cleanup(context.Background(), tt.experiment, logr.Discard())
			if (err != nil) != tt.wantErr {
				t.Errorf("DiskFailureInjector.Cleanup() error = %v, wantErr %v", err, tt.wantErr)
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

				if _, ok := updatedPod.Annotations["havock8s.io/disk-failure"]; ok {
					t.Error("Disk failure annotation still present after cleanup")
				}
			}
		})
	}
} 