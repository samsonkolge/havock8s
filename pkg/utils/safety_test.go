package utils

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

func TestSafetyChecker_CheckProtectedResources(t *testing.T) {
	tests := []struct {
		name       string
		experiment *chaosv1alpha1.Havock8sExperiment
		pod        *corev1.Pod
		wantRollback bool
		wantReason   string
	}{
		{
			name: "protected namespace",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Target: chaosv1alpha1.TargetSpec{
						TargetType: "Pod",
						Name:      "test-pod",
						Namespace: "kube-system",
					},
				},
			},
			wantRollback: true,
			wantReason:   "Target namespace kube-system is protected",
		},
		{
			name: "protected pod",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Target: chaosv1alpha1.TargetSpec{
						TargetType: "Pod",
						Name:      "test-pod",
						Namespace: "default",
					},
				},
			},
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						"havock8s.io/protected": "true",
					},
				},
			},
			wantRollback: true,
			wantReason:   "Pod has protection annotation",
		},
		{
			name: "unprotected pod",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Target: chaosv1alpha1.TargetSpec{
						TargetType: "Pod",
						Name:      "test-pod",
						Namespace: "default",
					},
				},
			},
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
			},
			wantRollback: false,
			wantReason:   "",
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

			s := NewSafetyChecker(fakeClient)
			gotRollback, gotReason := s.CheckProtectedResources(context.Background(), tt.experiment, logr.Discard())

			if gotRollback != tt.wantRollback {
				t.Errorf("SafetyChecker.CheckProtectedResources() rollback = %v, want %v", gotRollback, tt.wantRollback)
			}
			if gotReason != tt.wantReason {
				t.Errorf("SafetyChecker.CheckProtectedResources() reason = %v, want %v", gotReason, tt.wantReason)
			}
		})
	}
}

func TestSafetyChecker_CheckHealthEndpoints(t *testing.T) {
	tests := []struct {
		name       string
		experiment *chaosv1alpha1.Havock8sExperiment
		wantRollback bool
		wantReason   string
	}{
		{
			name: "no health checks",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{},
			},
			wantRollback: false,
			wantReason:   "",
		},
		{
			name: "with health checks",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Safety: &chaosv1alpha1.SafetySpec{
						HealthChecks: []chaosv1alpha1.HealthCheckSpec{
							{
								Type: "httpGet",
								Path: "/health",
								Port: 8080,
							},
						},
					},
				},
			},
			wantRollback: true, // Since we're not actually running a server
			wantReason:   "Health check failed for endpoint /health:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = corev1.AddToScheme(scheme)
			_ = chaosv1alpha1.AddToScheme(scheme)

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tt.experiment).
				Build()

			s := NewSafetyChecker(fakeClient)
			gotRollback, gotReason := s.CheckHealthEndpoints(context.Background(), tt.experiment, logr.Discard())

			if gotRollback != tt.wantRollback {
				t.Errorf("SafetyChecker.CheckHealthEndpoints() rollback = %v, want %v", gotRollback, tt.wantRollback)
			}
			if gotReason != tt.wantReason {
				t.Errorf("SafetyChecker.CheckHealthEndpoints() reason = %v, want %v", gotReason, tt.wantReason)
			}
		})
	}
}

func TestSafetyChecker_CheckMetricConditions(t *testing.T) {
	tests := []struct {
		name       string
		experiment *chaosv1alpha1.Havock8sExperiment
		wantRollback bool
		wantReason   string
	}{
		{
			name: "no metric conditions",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{},
			},
			wantRollback: false,
			wantReason:   "",
		},
		{
			name: "with metric conditions",
			experiment: &chaosv1alpha1.Havock8sExperiment{
				Spec: chaosv1alpha1.Havock8sExperimentSpec{
					Safety: &chaosv1alpha1.SafetySpec{
						PauseConditions: []chaosv1alpha1.PauseConditionSpec{
							{
								Type:        "metric",
								MetricQuery: "mongodb_connections",
								Threshold:   "100",
							},
						},
					},
				},
			},
			wantRollback: false, // Since our mock always returns 0
			wantReason:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = corev1.AddToScheme(scheme)
			_ = chaosv1alpha1.AddToScheme(scheme)

			fakeClient := fake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(tt.experiment).
				Build()

			s := NewSafetyChecker(fakeClient)
			gotRollback, gotReason := s.CheckMetricConditions(context.Background(), tt.experiment, logr.Discard())

			if gotRollback != tt.wantRollback {
				t.Errorf("SafetyChecker.CheckMetricConditions() rollback = %v, want %v", gotRollback, tt.wantRollback)
			}
			if gotReason != tt.wantReason {
				t.Errorf("SafetyChecker.CheckMetricConditions() reason = %v, want %v", gotReason, tt.wantReason)
			}
		})
	}
} 