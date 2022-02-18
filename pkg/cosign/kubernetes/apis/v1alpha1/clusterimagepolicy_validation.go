package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-sigstore-dev-v1alpha1-clusterimagepolicy,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=sigstore.dev,resources=clusterimagepolicies,versions=v1alpha1,name=validation.clusterimagepolicy.sigstore.dev,sideEffects=None,admissionReviewVersions=v1

func (ip *ClusterImagePolicy) Validate(ctx context.Context) *apis.FieldError {
	return nil
}
