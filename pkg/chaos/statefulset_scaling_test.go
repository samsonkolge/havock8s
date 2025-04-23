package chaos

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestStatefulSetScalingInjector_Inject(t *testing.T) {
	tests := []struct {
		name           string
		experiment     *chaosv1alpha1.Havock8sExperiment
		statefulSet    *appsv1.StatefulSet
		wantErr        bool
		wantReplicas   int32
		wantAnnotation bool
	}{
		{
			name: "successful scale down",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Parameters: map[string]string{
						"scaleMode":  "down",
						"scaleCount": "2",
						"scaleMin":   "1",
					},
				},
				Status: chaosv1alpha1.Havock8sExperimentStatus{
					TargetResources: []chaosv1alpha1.TargetResourceStatus{
						{
							Kind:      "StatefulSet",
							Name:      "test-sts",
							Namespace: "default",
						},
					},
				},
			},
			statefulSet: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sts",
					Namespace: "default",
				},
				Spec: appsv1.StatefulSetSpec{
					Replicas: int32Ptr(3),
				},
			},
			wantErr:        false,
			wantReplicas:   1, // 3 - 2 = 1
			wantAnnotation: true,
		},
		{
			name: "scale down below minimum",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Parameters: map[string]string{
						"scaleMode":  "down",
						"scaleCount": "3",
						"scaleMin":   "2",
					},
				},
				Status: chaosv1alpha1.Havock8sExperimentStatus{
					TargetResources: []chaosv1alpha1.TargetResourceStatus{
						{
							Kind:      "StatefulSet",
							Name:      "test-sts",
							Namespace: "default",
						},
					},
				},
			},
			statefulSet: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sts",
					Namespace: "default",
				},
				Spec: appsv1.StatefulSetSpec{
					Replicas: int32Ptr(3),
				},
			},
			wantErr:        false,
			wantReplicas:   2, // Should not go below scaleMin
			wantAnnotation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = appsv1.AddToScheme(scheme)
			_ = chaosv1alpha1.AddToScheme(scheme)

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tt.experiment, tt.statefulSet).
				Build()

			injector := &StatefulSetScalingInjector{}
			injector.SetClient(fakeClient)

			err := injector.Inject(context.Background(), tt.experiment, logr.Discard())
			if (err != nil) != tt.wantErr {
				t.Errorf("StatefulSetScalingInjector.Inject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				updatedSTS := &appsv1.StatefulSet{}
				err := fakeClient.Get(context.Background(), client.ObjectKey{
					Namespace: tt.statefulSet.Namespace,
					Name:      tt.statefulSet.Name,
				}, updatedSTS)
				if err != nil {
					t.Errorf("Failed to get updated StatefulSet: %v", err)
					return
				}

				if *updatedSTS.Spec.Replicas != tt.wantReplicas {
					t.Errorf("Replicas = %v, want %v", *updatedSTS.Spec.Replicas, tt.wantReplicas)
				}

				if tt.wantAnnotation {
					if _, ok := updatedSTS.Annotations["havock8s.io/original-replicas"]; !ok {
						t.Error("Original replicas annotation not set")
					}
				}
			}
		})
	}
}

func TestStatefulSetScalingInjector_Cleanup(t *testing.T) {
	tests := []struct {
		name        string
		experiment  *chaosv1alpha1.Havock8sExperiment
		statefulSet *appsv1.StatefulSet
		wantErr     bool
		wantReplicas int32
	}{
		{
			name: "successful cleanup",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Status: chaosv1alpha1.Havock8sExperimentStatus{
					TargetResources: []chaosv1alpha1.TargetResourceStatus{
						{
							Kind:      "StatefulSet",
							Name:      "test-sts",
							Namespace: "default",
						},
					},
				},
			},
			statefulSet: &appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-sts",
					Namespace: "default",
					Annotations: map[string]string{
						"havock8s.io/original-replicas": "3",
					},
				},
				Spec: appsv1.StatefulSetSpec{
					Replicas: int32Ptr(1),
				},
			},
			wantErr:     false,
			wantReplicas: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = appsv1.AddToScheme(scheme)
			_ = chaosv1alpha1.AddToScheme(scheme)

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tt.experiment, tt.statefulSet).
				Build()

			injector := &StatefulSetScalingInjector{}
			injector.SetClient(fakeClient)

			err := injector.Cleanup(context.Background(), tt.experiment, logr.Discard())
			if (err != nil) != tt.wantErr {
				t.Errorf("StatefulSetScalingInjector.Cleanup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				updatedSTS := &appsv1.StatefulSet{}
				err := fakeClient.Get(context.Background(), client.ObjectKey{
					Namespace: tt.statefulSet.Namespace,
					Name:      tt.statefulSet.Name,
				}, updatedSTS)
				if err != nil {
					t.Errorf("Failed to get updated StatefulSet: %v", err)
					return
				}

				if *updatedSTS.Spec.Replicas != tt.wantReplicas {
					t.Errorf("Replicas = %v, want %v", *updatedSTS.Spec.Replicas, tt.wantReplicas)
				}

				if _, ok := updatedSTS.Annotations["havock8s.io/original-replicas"]; ok {
					t.Error("Original replicas annotation still present after cleanup")
				}
			}
		})
	}
}

func int32Ptr(i int32) *int32 {
	return &i
} 