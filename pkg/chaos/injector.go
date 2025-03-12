package chaos

import (
	"context"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
)

// Injector is the interface that defines chaos injection and cleanup methods
type Injector interface {
	// Inject applies chaos to the target resources
	Inject(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error

	// Cleanup removes chaos from the target resources
	Cleanup(ctx context.Context, experiment *chaosv1alpha1.StatefulChaosExperiment, log logr.Logger) error
}

// GetInjector returns the appropriate injector for the given chaos type
func GetInjector(chaosType string) (Injector, error) {
	switch chaosType {
	case "DiskFailure":
		return &DiskFailureInjector{}, nil
	case "NetworkLatency":
		return &NetworkLatencyInjector{}, nil
	case "PodFailure":
		return &PodFailureInjector{}, nil
	case "StatefulSetScaling":
		return &StatefulSetScalingInjector{}, nil
	default:
		return nil, nil // Return nil to indicate unsupported chaos type
	}
}
