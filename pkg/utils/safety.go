package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SafetyChecker provides methods to check if safety conditions are met
type SafetyChecker struct {
	client client.Client
}

// NewSafetyChecker creates a new SafetyChecker
func NewSafetyChecker(c client.Client) *SafetyChecker {
	return &SafetyChecker{
		client: c,
	}
}

// CheckHealthEndpoints checks if health endpoints are responding correctly
func (s *SafetyChecker) CheckHealthEndpoints(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) (bool, string) {
	if experiment.Spec.Safety == nil || len(experiment.Spec.Safety.HealthChecks) == 0 {
		return false, ""
	}

	for _, healthCheck := range experiment.Spec.Safety.HealthChecks {
		switch healthCheck.Type {
		case "httpGet":
			healthy, reason := s.checkHTTPEndpoint(ctx, experiment, healthCheck, log)
			if !healthy {
				return true, reason // Return true to indicate rollback needed
			}
		case "tcpSocket":
			healthy, reason := s.checkTCPEndpoint(ctx, experiment, healthCheck, log)
			if !healthy {
				return true, reason // Return true to indicate rollback needed
			}
		case "exec":
			// In a real implementation, this would execute a command in a pod
			// and check the result
		}
	}

	return false, "" // All checks passed, no rollback needed
}

// checkHTTPEndpoint checks if an HTTP endpoint is healthy
func (s *SafetyChecker) checkHTTPEndpoint(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, healthCheck chaosv1alpha1.HealthCheckSpec, log logr.Logger) (bool, string) {
	// In a real implementation, this would:
	// 1. Find the pod/service with the endpoint
	// 2. Make HTTP requests to the endpoint
	// 3. Check if the response is as expected

	// For this example, we'll just simulate a check
	log.Info("Checking HTTP endpoint", "path", healthCheck.Path, "port", healthCheck.Port)

	// Create a simple HTTP client with timeout
	_ = &http.Client{
		Timeout: 5 * time.Second,
	}

	// Construct URL (in a real implementation, this would use the actual pod/service IP)
	url := fmt.Sprintf("http://example:8080%s", healthCheck.Path)
	log.Info("Would check URL", "url", url)

	// In a real implementation, we would create and execute the request
	// req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	// resp, err := client.Do(req)

	// Simulate a healthy response
	return false, ""
}

// checkTCPEndpoint checks if a TCP endpoint is reachable
func (s *SafetyChecker) checkTCPEndpoint(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, healthCheck chaosv1alpha1.HealthCheckSpec, log logr.Logger) (bool, string) {
	// In a real implementation, this would:
	// 1. Find the pod/service with the endpoint
	// 2. Try to establish a TCP connection
	// 3. Check if the connection is successful

	// For this example, we'll just simulate a check
	log.Info("Checking TCP endpoint", "port", healthCheck.Port)

	// Simulate a healthy response
	return false, ""
}

// CheckProtectedResources checks if the experiment targets any protected resources
func (s *SafetyChecker) CheckProtectedResources(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) (bool, string) {
	if experiment.Spec.Safety == nil || len(experiment.Spec.Safety.ResourceProtections) == 0 {
		return false, ""
	}

	for _, target := range experiment.Status.TargetResources {
		for _, protection := range experiment.Spec.Safety.ResourceProtections {
			switch protection.Type {
			case "Namespace":
				if target.Namespace == protection.Value {
					return true, fmt.Sprintf("Target namespace %s is protected", target.Namespace)
				}
			case "Label":
				// In a real implementation, this would check if the target has the protected label
				// For this example, we'll just simulate a check
			case "Annotation":
				// In a real implementation, this would check if the target has the protected annotation
				// For this example, we'll just simulate a check
			case "Name":
				if target.Name == protection.Value {
					return true, fmt.Sprintf("Target name %s is protected", target.Name)
				}
			}
		}
	}

	return false, "" // No protected resources found
}

// CheckMetricConditions checks if any metric-based pause conditions are met
func (s *SafetyChecker) CheckMetricConditions(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) (bool, string) {
	if experiment.Spec.Safety == nil || len(experiment.Spec.Safety.PauseConditions) == 0 {
		return false, ""
	}

	for _, condition := range experiment.Spec.Safety.PauseConditions {
		if condition.Type == "metric" && condition.MetricQuery != "" {
			// In a real implementation, this would:
			// 1. Query a metrics system (Prometheus, etc.) with the provided query
			// 2. Evaluate the result against a threshold
			// 3. Determine if the condition is met

			log.Info("Checking metric condition", "query", condition.MetricQuery)

			// For this example, we'll just simulate a check
			// In a real implementation, we would query Prometheus or another metrics system

			// Simulate that all conditions are fine
			return false, ""
		}
	}

	return false, "" // No metric conditions triggered
}

// IsProtectedPod checks if a pod is protected from chaos
func (s *SafetyChecker) IsProtectedPod(ctx context.Context, pod *corev1.Pod, log logr.Logger) (bool, string) {
	// Check for protection annotations
	if pod.Annotations != nil {
		if val, ok := pod.Annotations["statefulchaos.io/protected"]; ok && val == "true" {
			return true, "Pod has protection annotation"
		}
	}

	// Check for protection labels
	if pod.Labels != nil {
		if val, ok := pod.Labels["statefulchaos.io/protected"]; ok && val == "true" {
			return true, "Pod has protection label"
		}
	}

	// Check if pod is in kube-system namespace
	if pod.Namespace == "kube-system" {
		return true, "Pod is in kube-system namespace"
	}

	return false, ""
}
