// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "go.pinniped.dev/generated/latest/client/concierge/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// TokenCredentialRequests returns a TokenCredentialRequestInformer.
	TokenCredentialRequests() TokenCredentialRequestInformer
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

// TokenCredentialRequests returns a TokenCredentialRequestInformer.
func (v *version) TokenCredentialRequests() TokenCredentialRequestInformer {
	return &tokenCredentialRequestInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
