package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	GroupName = "sigstore.dev"
)

var SchemeGroupVersion = schema.GroupVersion{Group: "sigstore.dev", Version: "v1alpha1"}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
