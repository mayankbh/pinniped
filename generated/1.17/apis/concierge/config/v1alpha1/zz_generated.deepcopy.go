// +build !ignore_autogenerated

// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CredentialIssuer) DeepCopyInto(out *CredentialIssuer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CredentialIssuer.
func (in *CredentialIssuer) DeepCopy() *CredentialIssuer {
	if in == nil {
		return nil
	}
	out := new(CredentialIssuer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CredentialIssuer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CredentialIssuerKubeConfigInfo) DeepCopyInto(out *CredentialIssuerKubeConfigInfo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CredentialIssuerKubeConfigInfo.
func (in *CredentialIssuerKubeConfigInfo) DeepCopy() *CredentialIssuerKubeConfigInfo {
	if in == nil {
		return nil
	}
	out := new(CredentialIssuerKubeConfigInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CredentialIssuerList) DeepCopyInto(out *CredentialIssuerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CredentialIssuer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CredentialIssuerList.
func (in *CredentialIssuerList) DeepCopy() *CredentialIssuerList {
	if in == nil {
		return nil
	}
	out := new(CredentialIssuerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CredentialIssuerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CredentialIssuerStatus) DeepCopyInto(out *CredentialIssuerStatus) {
	*out = *in
	if in.Strategies != nil {
		in, out := &in.Strategies, &out.Strategies
		*out = make([]CredentialIssuerStrategy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.KubeConfigInfo != nil {
		in, out := &in.KubeConfigInfo, &out.KubeConfigInfo
		*out = new(CredentialIssuerKubeConfigInfo)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CredentialIssuerStatus.
func (in *CredentialIssuerStatus) DeepCopy() *CredentialIssuerStatus {
	if in == nil {
		return nil
	}
	out := new(CredentialIssuerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CredentialIssuerStrategy) DeepCopyInto(out *CredentialIssuerStrategy) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CredentialIssuerStrategy.
func (in *CredentialIssuerStrategy) DeepCopy() *CredentialIssuerStrategy {
	if in == nil {
		return nil
	}
	out := new(CredentialIssuerStrategy)
	in.DeepCopyInto(out)
	return out
}
