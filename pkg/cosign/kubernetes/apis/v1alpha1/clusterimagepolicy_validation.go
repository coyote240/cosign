package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
)

func (ip *ClusterImagePolicy) Validate(ctx context.Context) *apis.FieldError {
	return nil
}
