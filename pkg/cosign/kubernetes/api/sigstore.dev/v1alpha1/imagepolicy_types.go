package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterImagePolicy struct {
	metav1.TypeMeta `json:",inline"`
	//+optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	//+optional
	Spec ClusterImagePolicySpec `json:"spec"`
}

var (
	_ apis.Validatable = (*ClusterImagePolicy)(nil)
	_ apis.Defaultable = (*ClusterImagePolicy)(nil)
)

type ClusterImagePolicySpec struct {
	Images []ImagePattern `json:"images"`
}

type ImagePattern struct {
	Pattern string `json:"pattern"`
	//+optional
	Authorities []Authority `json:"authorities"`
}

type Authority struct {
	Key     KeyRef     `json:"key"`
	Keyless KeylessRef `json:"keyless"`
	Sources []Source   `json:"source"`
	CTLog   TLog       `json:"ctlog"`
}

type KeyRef struct {
	SecretRef SecretRef `json:"secretRef"`
	Data      string    `json:"data"`
	KMS       string    `json:"kms"`
}

type SecretRef struct {
	Name string `json:"name"`
}

type Source struct {
	OCI string `json:"oci"`
}

type TLog struct {
	URL string `json:"url"`
}

type KeylessRef struct {
	Identities []Identity `json:"identities"`
	CAKey      CAKey      `json:"ca-key"`
}

type Identity struct {
	Issuer  string `json:"issuer"`
	Subject string `json:"subject"`
}

type CAKey struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
