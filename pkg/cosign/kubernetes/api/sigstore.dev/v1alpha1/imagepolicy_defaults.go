package v1alpha1

import (
	"context"
)

func (ip *ClusterImagePolicy) SetDefaults(ctx context.Context) {
	ip.Spec = ClusterImagePolicySpec{
		Images: []ImagePattern{
			ImagePattern{"*"},
		},
	}
}
