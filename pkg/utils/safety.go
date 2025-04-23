package utils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SafetyChecker provides safety mechanisms for chaos experiments
type SafetyChecker struct {
	client client.Client
}

// NewSafetyChecker creates a new SafetyChecker instance
func NewSafetyChecker(c client.Client) *SafetyChecker {
	return &SafetyChecker{client: c}
}

// CheckSafety performs all safety checks for an experiment
func (s *SafetyChecker) CheckSafety(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, logger logr.Logger) (bool, string) {
	// Check protected resources
	if shouldRollback, reason := s.CheckProtectedResources(ctx, experiment, logger); shouldRollback {
		return true, reason
	}

	// Check health endpoints
	if shouldRollback, reason := s.CheckHealthEndpoints(ctx, experiment, logger); shouldRollback {
		return true, reason
	}

	// Check metric conditions
	if shouldRollback, reason := s.CheckMetricConditions(ctx, experiment, logger); shouldRollback {
		return true, reason
	}

	return false, ""
}

// CheckProtectedResources verifies that no protected resources will be affected
func (s *SafetyChecker) CheckProtectedResources(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, logger logr.Logger) (bool, string) {
	// Check if target namespace is protected
	if experiment.Spec.Target.Namespace == "kube-system" {
		return true, "Target namespace kube-system is protected"
	}

	// Check if target pod has protection annotation
	if experiment.Spec.Target.TargetType == "Pod" {
		pod := &corev1.Pod{}
		err := s.client.Get(ctx, types.NamespacedName{
			Name:      experiment.Spec.Target.Name,
			Namespace: experiment.Spec.Target.Namespace,
		}, pod)
		if err != nil {
			logger.Error(err, "Failed to get target pod")
			return true, "Failed to verify pod protection status"
		}

		if pod.Annotations["havock8s.io/protected"] == "true" {
			return true, "Pod has protection annotation"
		}
	}

	return false, ""
}

// CheckHealthEndpoints verifies that health check endpoints are responding
func (s *SafetyChecker) CheckHealthEndpoints(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, logger logr.Logger) (bool, string) {
	if experiment.Spec.Safety == nil || len(experiment.Spec.Safety.HealthChecks) == 0 {
		return false, ""
	}

	for _, check := range experiment.Spec.Safety.HealthChecks {
		switch check.Type {
		case "httpGet":
			if !s.checkHTTPEndpoint(check.Path, check.Port) {
				return true, fmt.Sprintf("Health check failed for endpoint %s:%d", check.Path, check.Port)
			}
		case "tcpSocket":
			if !s.checkTCPEndpoint(check.Port) {
				return true, fmt.Sprintf("Health check failed for TCP port %d", check.Port)
			}
		}
	}

	return false, ""
}

// CheckMetricConditions verifies that metric-based conditions are met
func (s *SafetyChecker) CheckMetricConditions(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, logger logr.Logger) (bool, string) {
	if experiment.Spec.Safety == nil || len(experiment.Spec.Safety.PauseConditions) == 0 {
		return false, ""
	}

	for _, condition := range experiment.Spec.Safety.PauseConditions {
		if condition.Type == "metric" {
			// Example: Check MongoDB connections
			if condition.MetricQuery == "mongodb_connections" {
				threshold, err := strconv.ParseFloat(condition.Threshold, 64)
				if err != nil {
					logger.Error(err, "Failed to parse threshold")
					continue
				}

				currentValue := s.getMetricValue(condition.MetricQuery)
				if currentValue > threshold {
					return true, fmt.Sprintf("Metric %s exceeds threshold: %v > %v",
						condition.MetricQuery, currentValue, threshold)
				}
			}
		}
	}

	return false, ""
}

// Helper functions

func (s *SafetyChecker) checkHTTPEndpoint(path string, port int32) bool {
	url := fmt.Sprintf("http://localhost:%d%s", port, path)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (s *SafetyChecker) checkTCPEndpoint(port int32) bool {
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func (s *SafetyChecker) getMetricValue(query string) float64 {
	// This is a placeholder. In a real implementation, this would query
	// a metrics system like Prometheus
	return 0.0
}



