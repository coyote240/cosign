// +kubebuilder:validation:Optional
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clusterimagepolicies,scope=Cluster
// +kubebuilder:storageversion
type ClusterImagePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterImagePolicySpec `json:"spec,omitempty"`
}

var (
	_ apis.Validatable = (*ClusterImagePolicy)(nil)
	_ apis.Defaultable = (*ClusterImagePolicy)(nil)
)

type ClusterImagePolicySpec struct {
	Images []ImagePattern `json:"images,omitempty"`
}

type ImagePattern struct {
	Pattern     string      `json:"pattern,omitempty"`
	Authorities []Authority `json:"authorities,anyOf,omitempty"`
}

type Authority struct {
	Key     KeyRef     `json:"key,omitempty"`
	Keyless KeylessRef `json:"keyless,omitempty"`
	Sources []Source   `json:"source,anyOf,omitempty"`
	CTLog   TLog       `json:"ctlog,omitempty"`
}

type KeyRef struct {
	SecretRef SecretRef `json:"secretRef,omitempty"`
	Data      string    `json:"data,omitempty"`
	KMS       string    `json:"kms,omitempty"`
}

type SecretRef struct {
	Name string `json:"name,omitempty"`
}

type Source struct {
	OCI string `json:"oci"`
}

type TLog struct {
	URL string `json:"url,omitempty"`
}

type KeylessRef struct {
	Identities []Identity `json:"identities,anyOf,omitempty"`
	CAKey      CAKey      `json:"ca-key,omitempty"`
}

type Identity struct {
	Issuer  string `json:"issuer,omitempty"`
	Subject string `json:"subject,omitempty"`
}

type CAKey struct {
	Name string `json:"name,omitempty"`
	Data string `json:"data,omitempty"`
}
