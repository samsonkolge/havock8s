package utils

import (
	"context"

	chaosv1alpha1 "github.com/havock8s/havock8s/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CheckTargetExists verifies if the target resource exists
func CheckTargetExists(ctx context.Context, c client.Client, target chaosv1alpha1.TargetSpec) (bool, error) {
	switch target.TargetType {
	case "Pod":
		pod := &corev1.Pod{}
		err := c.Get(ctx, types.NamespacedName{
			Name:      target.Name,
			Namespace: target.Namespace,
		}, pod)
		if err != nil {
			if client.IgnoreNotFound(err) == nil {
				return false, nil
			}
			return false, err
		}
		return true, nil
	default:
		return false, nil
	}
} 