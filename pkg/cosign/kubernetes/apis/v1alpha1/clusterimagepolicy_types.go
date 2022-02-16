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
	metav1.TypeMeta `json:",inline"`

	// +kubebuilder:validation:Required
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterImagePolicySpec `json:"spec"`
}

var (
	_ apis.Validatable = (*ClusterImagePolicy)(nil)
	_ apis.Defaultable = (*ClusterImagePolicy)(nil)
)

type ClusterImagePolicySpec struct {
	// +kubebuilder:validation:MinItems=1
	Images []ImagePattern `json:"images"`
}

type ImagePattern struct {
	// +kubebuilder:validation:Required
	Pattern string `json:"pattern"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Authorities []Authority `json:"authorities"`
}

// +kubebuilder:validation:MinProperties=1
type Authority struct {
	Key     KeyRef     `json:"key"`
	Keyless KeylessRef `json:"keyless"`
	Sources []Source   `json:"source"`
	CTLog   TLog       `json:"ctlog"`
}

// +kubebuilder:validation:MinProperties=1
type KeyRef struct {
	SecretRef SecretRef `json:"secretRef"`
	Data      string    `json:"data"`
	KMS       string    `json:"kms"`
}

type SecretRef struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

type Source struct {
	// +kubebuilder:validation:Required
	OCI string `json:"oci"`
}

type TLog struct {
	// +kubebuilder:validation:Required
	URL string `json:"url"`
}

// +kubebuilder:validation:MaxProperties=1
type KeylessRef struct {
	Identities []Identity `json:"identities"`
	CAKey      CAKey      `json:"ca-key"`
}

type Identity struct {
	// +kubebuilder:validation:Required
	Issuer string `json:"issuer"`

	// +kubebuilder:validation:Required
	Subject string `json:"subject"`
}

type CAKey struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	Data string `json:"data"`
}
