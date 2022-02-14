package v1alpha1

import (
	"context"
)

// +kubebuilder:webhook:verbs=create;update,path=/default-sigstore-dev-v1alpha1-clusterimagepolicy,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=sigstore.dev,resources=clusterimagepolicies,versions=v1alpha1,name=defaulting.clusterimagepolicy.sigstore.dev,sideEffects=None,admissionReviewVersions=v1;v1beta1

func (ip *ClusterImagePolicy) SetDefaults(ctx context.Context) {
	/*
		ip.Spec = ClusterImagePolicySpec{
			Images: []ImagePattern{
				ImagePattern{"*"},
			},
		}
	*/
}
