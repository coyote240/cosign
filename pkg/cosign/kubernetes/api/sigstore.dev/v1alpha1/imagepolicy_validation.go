package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

func (ip *ClusterImagePolicy) Validate(ctx context.Context) *apis.FieldError {
	return ip.Spec.Validate(ctx).ViaField("metadata")
}

func (metadata *metav1.ObjectMeta) Validate(ctx context.Context) *apis.FieldError {
	if metadata.Name != "image-policy" {
		return apis.ErrInvalidValue(metadata.Name, "metadata.name", "metadata.name must be \"image-policy\"")
	}
	return nil
}
