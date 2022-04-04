// Copyright 2022 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clusterimagepolicy

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/sigstore/cosign/pkg/apis/cosigned/v1alpha1"
	"knative.dev/pkg/logging"
)

// ClusterImagePolicy defines the images that go through verification
// and the authorities used for verification.
// This is the internal representation of the external v1alpha1.ClusterImagePolicy.
// KeyRef does not store secretRefs in internal representation.
// KeyRef does store parsed publicKeys from Data in internal representation.
type ClusterImagePolicy struct {
	Images      []v1alpha1.ImagePattern `json:"images"`
	Authorities []Authority             `json:"authorities"`
}

type Authority struct {
	// +optional
	Key *KeyRef `json:"key,omitempty"`
	// +optional
	Keyless *v1alpha1.KeylessRef `json:"keyless,omitempty"`
	// +optional
	Sources []v1alpha1.Source `json:"source,omitempty"`
	// +optional
	CTLog *v1alpha1.TLog `json:"ctlog,omitempty"`
}

// This references a public verification key stored in
// a secret in the cosign-system namespace.
type KeyRef struct {
	// Data contains the inline public key
	// +optional
	Data string `json:"data,omitempty"`
	// KMS contains the KMS url of the public key
	// +optional
	KMS string `json:"kms,omitempty"`
	// +optional
	PublicKeys []*ecdsa.PublicKey `json:"publicKeys,omitempty"`
}

func ConvertClusterImagePolicyV1alpha1ToInternal(ctx context.Context, in *v1alpha1.ClusterImagePolicy) (*ClusterImagePolicy, error) {
	outAuthorities := make([]Authority, 0)
	for _, authority := range in.Spec.Authorities {
		outAuthority, err := convertAuthorityV1Alpha1ToInternal(ctx, &authority)
		if err != nil {
			return nil, err
		}
		outAuthorities = append(outAuthorities, outAuthority)
	}

	return &ClusterImagePolicy{
		Images:      in.Spec.Images,
		Authorities: outAuthorities,
	}, nil
}

func convertAuthorityV1Alpha1ToInternal(ctx context.Context, in *v1alpha1.Authority) (Authority, error) {
	keyRef, err := convertKeyRefV1Alpha1ToInternal(ctx, in.Key)
	if err != nil {
		return Authority{}, err
	}
	return Authority{
		Key:     keyRef,
		Keyless: in.Keyless,
		Sources: in.Sources,
		CTLog:   in.CTLog,
	}, nil
}

func convertKeyRefV1Alpha1ToInternal(ctx context.Context, in *v1alpha1.KeyRef) (*KeyRef, error) {
	if in == nil {
		return nil, nil
	}

	var publicKeys []*ecdsa.PublicKey
	var err error
	if in.Data != "" {
		publicKeys, err = ConvertKeyDataToPublicKeys(ctx, in.Data)
		if err != nil {
			return nil, err
		}
	}

	return &KeyRef{
		KMS:        in.KMS,
		PublicKeys: publicKeys,
	}, nil
}

func ConvertKeyDataToPublicKeys(ctx context.Context, pubKey string) ([]*ecdsa.PublicKey, error) {
	keys := []*ecdsa.PublicKey{}

	logging.FromContext(ctx).Debugf("Got public key: %v", pubKey)

	pems := parsePems([]byte(pubKey))
	for _, p := range pems {
		key, err := x509.ParsePKIXPublicKey(p.Bytes)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key.(*ecdsa.PublicKey))
	}
	return keys, nil
}

func parsePems(b []byte) []*pem.Block {
	p, rest := pem.Decode(b)
	if p == nil {
		return nil
	}
	pems := []*pem.Block{p}

	if rest != nil {
		return append(pems, parsePems(rest)...)
	}
	return pems
}
