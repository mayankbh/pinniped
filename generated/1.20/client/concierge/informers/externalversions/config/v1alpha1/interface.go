// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "go.pinniped.dev/generated/1.20/client/concierge/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// CredentialIssuers returns a CredentialIssuerInformer.
	CredentialIssuers() CredentialIssuerInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// CredentialIssuers returns a CredentialIssuerInformer.
func (v *version) CredentialIssuers() CredentialIssuerInformer {
	return &credentialIssuerInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
