package chaos

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Injector defines the interface for chaos injection
type Injector interface {
	// SetClient sets the Kubernetes client
	SetClient(c client.Client)

	// Inject applies chaos to the target resources
	Inject(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, log logr.Logger) error

	// Cleanup removes chaos from the target resources
	Cleanup(ctx context.Context, experiment *chaosv1alpha1.Havock8sExperiment, log logr.Logger) error
}

var injectors = make(map[string]Injector)

// RegisterInjector registers a new chaos injector
func RegisterInjector(name string, injector Injector) {
	injectors[name] = injector
}

// GetInjector returns a chaos injector by name
func GetInjector(name string) (Injector, error) {
	injector, ok := injectors[name]
	if !ok {
		return nil, fmt.Errorf("no injector registered for chaos type: %s", name)
	}
	return injector, nil
}
