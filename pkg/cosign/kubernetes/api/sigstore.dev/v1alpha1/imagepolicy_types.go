package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clusterimagepolicies,scope=Cluster
// +kubebuilder:storageversion
type ClusterImagePolicy struct {
	metav1.TypeMeta `json:",inline"`
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Optional
	Spec ClusterImagePolicySpec `json:"spec"`
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

	// +kubebuilder:validation:MinItems=1
	Authorities []Authority `json:"authorities"`
}

// +kubebuilder:validation:MinProperties=1
type Authority struct {
	Key     KeyRef     `json:"key"`
	Keyless KeylessRef `json:"keyless"`

	// +kubebuilder:validation:Optional
	Sources []Source `json:"source"`

	// +kubebuilder:validation:Optional
	CTLog TLog `json:"ctlog"`
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
