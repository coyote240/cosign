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
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/sigstore/cosign/pkg/apis/cosigned/v1alpha1"
	"github.com/sigstore/cosign/pkg/apis/utils"
	sigs "github.com/sigstore/cosign/pkg/signature"
	"github.com/sigstore/sigstore/pkg/signature/kms"
	signatureoptions "github.com/sigstore/sigstore/pkg/signature/options"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/system"
	"knative.dev/pkg/tracker"
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
	Keyless *KeylessRef `json:"keyless,omitempty"`
	// +optional
	Sources []v1alpha1.Source `json:"source,omitempty"`
	// +optional
	CTLog *v1alpha1.TLog `json:"ctlog,omitempty"`
}

// This references a public verification key stored in
// a secret in the cosign-system namespace.
type KeyRef struct {
	// TODO(dennyhoang): Remove Data
	// Data contains the inline public key
	// +optional
	Data string `json:"data,omitempty"`
	// KMS contains the KMS url of the public key
	// +optional
	PublicKeys []*ecdsa.PublicKey `json:"publicKeys,omitempty"`
}

type KeylessRef struct {
	// +optional
	URL *apis.URL `json:"url,omitempty"`
	// +optional
	Identities []v1alpha1.Identity `json:"identities,omitempty"`
	// +optional
	CACert *KeyRef `json:"ca-cert,omitempty"`
}

func (k *KeyRef) UnmarshalJSON(data []byte) error {
	var publicKeys []*ecdsa.PublicKey
	var err error

	publicKey := string(data)

	if publicKey != "" {
		publicKeys, err = ConvertKeyDataToPublicKeys(context.Background(), publicKey)
		if err != nil {
			return err
		}
	}

	k.PublicKeys = publicKeys

	return nil
}

func ConvertClusterImagePolicyV1alpha1ToInternal(ctx context.Context, in *v1alpha1.ClusterImagePolicy, rtracker tracker.Interface, secretLister corev1listers.SecretLister) (*ClusterImagePolicy, error) {
	copyIn := in.DeepCopy()
	outAuthorities := make([]Authority, 0)
	for _, authority := range copyIn.Spec.Authorities {
		outAuthority, err := convertAuthorityV1Alpha1ToInternal(ctx, copyIn, authority, rtracker, secretLister)
		if err != nil {
			return nil, err
		}
		outAuthorities = append(outAuthorities, *outAuthority)
	}

	return &ClusterImagePolicy{
		Images:      copyIn.Spec.Images,
		Authorities: outAuthorities,
	}, nil
}

func convertAuthorityV1Alpha1ToInternal(ctx context.Context, cipIn *v1alpha1.ClusterImagePolicy, in v1alpha1.Authority, rtracker tracker.Interface, secretLister corev1listers.SecretLister) (*Authority, error) {
	var publicKey string
	var err error

	// Key Handle
	if in.Key != nil {
		publicKey = in.Key.Data
	}

	if in.Key != nil && in.Key.SecretRef != nil {
		if publicKey, err = dataAndTrackSecret(ctx, cipIn, in.Key, rtracker, secretLister); err != nil {
			logging.FromContext(ctx).Errorf("Failed to read secret %q: %v", in.Key.SecretRef.Name, err)
			return nil, err
		}
	} else if in.Key != nil && in.Key.KMS != "" {
		if strings.Contains(in.Key.KMS, "://") {
			publicKey, err = GetKMSPublicKey(ctx, in.Key.KMS)
			if err != nil {
				return nil, err
			}
		}
	}

	// Populate KeyRef with publicKeys
	var keyRef *KeyRef
	if publicKey != "" {
		keyRef, err = convertKeyRefV1Alpha1ToInternal(ctx, publicKey)
	}
	if err != nil {
		return nil, err
	}

	// Populate KeylessRef
	var keylessRef *KeylessRef
	if in.Keyless != nil {
		keylessRef, err = convertKeylessRefV1Alpha1ToInternal(ctx, in.Keyless, cipIn, rtracker, secretLister)
		if err != nil {
			return nil, err
		}
	}

	return &Authority{
		Key:     keyRef,
		Keyless: keylessRef,
		Sources: in.Sources,
		CTLog:   in.CTLog,
	}, nil
}

func convertKeyRefV1Alpha1ToInternal(ctx context.Context, publicKey string) (*KeyRef, error) {
	var publicKeys []*ecdsa.PublicKey
	var err error

	if publicKey != "" {
		publicKeys, err = ConvertKeyDataToPublicKeys(ctx, publicKey)
		if err != nil {
			return nil, err
		}
	}

	return &KeyRef{
		PublicKeys: publicKeys,
		Data:       publicKey,
	}, nil
}

func convertKeylessRefV1Alpha1ToInternal(ctx context.Context, in *v1alpha1.KeylessRef, cipIn *v1alpha1.ClusterImagePolicy, rtracker tracker.Interface, secretLister corev1listers.SecretLister) (*KeylessRef, error) {
	var publicKey string
	var keyRef *KeyRef
	var err error

	if in != nil &&
		in.CACert != nil &&
		in.CACert.SecretRef != nil {
		publicKey, err = dataAndTrackSecret(ctx, cipIn, in.CACert, rtracker, secretLister)
		if err != nil {
			logging.FromContext(ctx).Errorf("Failed to read secret %q: %v", in.CACert.SecretRef.Name, err)
			return nil, err
		}
	}

	if publicKey != "" {
		keyRef, err = convertKeyRefV1Alpha1ToInternal(ctx, publicKey)
		if err != nil {
			return nil, err
		}
	}

	return &KeylessRef{
		URL:        in.URL,
		Identities: in.Identities,
		CACert:     keyRef,
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

// dataAndTrackSecret will take in a KeyRef and tries to read the Secret, finding the
// first key from it and return the data.
// Additionally, we set up a tracker so we will be notified if the secret
// is modified.
// There's still some discussion about how to handle multiple keys in a secret
// for now, just grab one from it. For reference, the discussion is here:
// TODO(vaikas): https://github.com/sigstore/cosign/issues/1573
func dataAndTrackSecret(ctx context.Context, cip *v1alpha1.ClusterImagePolicy, keyref *v1alpha1.KeyRef, rtracker tracker.Interface, secretLister corev1listers.SecretLister) (string, error) {
	if err := rtracker.TrackReference(tracker.Reference{
		APIVersion: "v1",
		Kind:       "Secret",
		Namespace:  system.Namespace(),
		Name:       keyref.SecretRef.Name,
	}, cip); err != nil {
		return "", fmt.Errorf("failed to track changes to secret %q : %w", keyref.SecretRef.Name, err)
	}
	secret, err := secretLister.Secrets(system.Namespace()).Get(keyref.SecretRef.Name)
	if err != nil {
		return "", err
	}
	if len(secret.Data) == 0 {
		return "", fmt.Errorf("secret %q contains no data", keyref.SecretRef.Name)
	}
	if len(secret.Data) > 1 {
		return "", fmt.Errorf("secret %q contains multiple data entries, only one is supported", keyref.SecretRef.Name)
	}
	for k, v := range secret.Data {
		logging.FromContext(ctx).Infof("inlining secret %q key %q", keyref.SecretRef.Name, k)
		if !utils.IsValidKey(v) {
			return "", fmt.Errorf("secret %q contains an invalid public key", keyref.SecretRef.Name)
		}

		return string(v), nil
	}
	return "", nil
}

// GetKMSPublicKey returns the public key as a string from the configured KMS service using the key ID
func GetKMSPublicKey(ctx context.Context, keyID string) (string, error) {
	kmsSigner, err := kms.Get(ctx, keyID, crypto.SHA256)
	if err != nil {
		logging.FromContext(ctx).Errorf("Failed to read KMS key ID %q: %v", keyID, err)
		return "", err
	}
	pemBytes, err := sigs.PublicKeyPem(kmsSigner, signatureoptions.WithContext(ctx))
	if err != nil {
		return "", err
	}
	return string(pemBytes), nil
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
