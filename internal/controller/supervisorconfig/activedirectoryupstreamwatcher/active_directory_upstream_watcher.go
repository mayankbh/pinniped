// Copyright 2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package activedirectoryupstreamwatcher implements a controller which watches ActiveDirectoryIdentityProviders.
package activedirectoryupstreamwatcher

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/klog/v2/klogr"

	"go.pinniped.dev/generated/latest/apis/supervisor/idp/v1alpha1"
	pinnipedclientset "go.pinniped.dev/generated/latest/client/supervisor/clientset/versioned"
	idpinformers "go.pinniped.dev/generated/latest/client/supervisor/informers/externalversions/idp/v1alpha1"
	pinnipedcontroller "go.pinniped.dev/internal/controller"
	"go.pinniped.dev/internal/controller/conditionsutil"
	"go.pinniped.dev/internal/controller/supervisorconfig/upstreamwatchers"
	"go.pinniped.dev/internal/controllerlib"
	"go.pinniped.dev/internal/oidc/provider"
	"go.pinniped.dev/internal/upstreamldap"
)

const (
	activeDirectoryControllerName = "active-directory-upstream-observer"

	// Default values for active directory config.
	defaultActiveDirectoryUsernameAttributeName = "sAMAccountName"
	defaultActiveDirectoryUIDAttributeName      = "objectGUID"
)

type activeDirectoryUpstreamGenericLDAPImpl struct {
	activeDirectoryIdentityProvider v1alpha1.ActiveDirectoryIdentityProvider
}

func (g *activeDirectoryUpstreamGenericLDAPImpl) Spec() upstreamwatchers.UpstreamGenericLDAPSpec {
	return &activeDirectoryUpstreamGenericLDAPSpec{g.activeDirectoryIdentityProvider}
}

func (g *activeDirectoryUpstreamGenericLDAPImpl) Namespace() string {
	return g.activeDirectoryIdentityProvider.Namespace
}

func (g *activeDirectoryUpstreamGenericLDAPImpl) Name() string {
	return g.activeDirectoryIdentityProvider.Name
}

func (g *activeDirectoryUpstreamGenericLDAPImpl) Generation() int64 {
	return g.activeDirectoryIdentityProvider.Generation
}

func (g *activeDirectoryUpstreamGenericLDAPImpl) Status() upstreamwatchers.UpstreamGenericLDAPStatus {
	return &activeDirectoryUpstreamGenericLDAPStatus{g.activeDirectoryIdentityProvider}
}

type activeDirectoryUpstreamGenericLDAPSpec struct {
	activeDirectoryIdentityProvider v1alpha1.ActiveDirectoryIdentityProvider
}

func (s *activeDirectoryUpstreamGenericLDAPSpec) Host() string {
	return s.activeDirectoryIdentityProvider.Spec.Host
}

func (s *activeDirectoryUpstreamGenericLDAPSpec) TLSSpec() *v1alpha1.TLSSpec {
	return s.activeDirectoryIdentityProvider.Spec.TLS
}

func (s *activeDirectoryUpstreamGenericLDAPSpec) BindSecretName() string {
	return s.activeDirectoryIdentityProvider.Spec.Bind.SecretName
}

func (s *activeDirectoryUpstreamGenericLDAPSpec) UserSearch() upstreamwatchers.UpstreamGenericLDAPUserSearch {
	return &activeDirectoryUpstreamGenericLDAPUserSearch{s.activeDirectoryIdentityProvider.Spec.UserSearch}
}

func (s *activeDirectoryUpstreamGenericLDAPSpec) GroupSearch() upstreamwatchers.UpstreamGenericLDAPGroupSearch {
	return &activeDirectoryUpstreamGenericLDAPGroupSearch{s.activeDirectoryIdentityProvider.Spec.GroupSearch}
}

type activeDirectoryUpstreamGenericLDAPUserSearch struct {
	userSearch v1alpha1.ActiveDirectoryIdentityProviderUserSearch
}

func (u *activeDirectoryUpstreamGenericLDAPUserSearch) Base() string {
	return u.userSearch.Base
}

func (u *activeDirectoryUpstreamGenericLDAPUserSearch) Filter() string {
	return u.userSearch.Filter
}

func (u *activeDirectoryUpstreamGenericLDAPUserSearch) UsernameAttribute() string {
	return u.userSearch.Attributes.Username
}

func (u *activeDirectoryUpstreamGenericLDAPUserSearch) UIDAttribute() string {
	return u.userSearch.Attributes.UID
}

type activeDirectoryUpstreamGenericLDAPGroupSearch struct {
	groupSearch v1alpha1.ActiveDirectoryIdentityProviderGroupSearch
}

func (g *activeDirectoryUpstreamGenericLDAPGroupSearch) Base() string {
	return g.groupSearch.Base
}

func (g *activeDirectoryUpstreamGenericLDAPGroupSearch) Filter() string {
	return g.groupSearch.Filter
}

func (g *activeDirectoryUpstreamGenericLDAPGroupSearch) GroupNameAttribute() string {
	return g.groupSearch.Attributes.GroupName
}

type activeDirectoryUpstreamGenericLDAPStatus struct {
	activeDirectoryIdentityProvider v1alpha1.ActiveDirectoryIdentityProvider
}

func (s *activeDirectoryUpstreamGenericLDAPStatus) Conditions() []v1alpha1.Condition {
	return s.activeDirectoryIdentityProvider.Status.Conditions
}

// UpstreamActiveDirectoryIdentityProviderICache is a thread safe cache that holds a list of validated upstream LDAP IDP configurations.
type UpstreamActiveDirectoryIdentityProviderICache interface {
	SetActiveDirectoryIdentityProviders([]provider.UpstreamLDAPIdentityProviderI)
}

type activeDirectoryWatcherController struct {
	cache                                   UpstreamActiveDirectoryIdentityProviderICache
	validatedSecretVersionsCache            *upstreamwatchers.SecretVersionCache
	ldapDialer                              upstreamldap.LDAPDialer
	client                                  pinnipedclientset.Interface
	activeDirectoryIdentityProviderInformer idpinformers.ActiveDirectoryIdentityProviderInformer
	secretInformer                          corev1informers.SecretInformer
}

// New instantiates a new controllerlib.Controller which will populate the provided UpstreamActiveDirectoryIdentityProviderICache.
func New(
	idpCache UpstreamActiveDirectoryIdentityProviderICache,
	client pinnipedclientset.Interface,
	activeDirectoryIdentityProviderInformer idpinformers.ActiveDirectoryIdentityProviderInformer,
	secretInformer corev1informers.SecretInformer,
	withInformer pinnipedcontroller.WithInformerOptionFunc,
) controllerlib.Controller {
	return newInternal(
		idpCache,
		// start with an empty secretVersionCache
		upstreamwatchers.NewSecretVersionCache(),
		// nil means to use a real production dialer when creating objects to add to the cache
		nil,
		client,
		activeDirectoryIdentityProviderInformer,
		secretInformer,
		withInformer,
	)
}

// For test dependency injection purposes.
func newInternal(
	idpCache UpstreamActiveDirectoryIdentityProviderICache,
	validatedSecretVersionsCache *upstreamwatchers.SecretVersionCache,
	ldapDialer upstreamldap.LDAPDialer,
	client pinnipedclientset.Interface,
	activeDirectoryIdentityProviderInformer idpinformers.ActiveDirectoryIdentityProviderInformer,
	secretInformer corev1informers.SecretInformer,
	withInformer pinnipedcontroller.WithInformerOptionFunc,
) controllerlib.Controller {
	c := activeDirectoryWatcherController{
		cache:                                   idpCache,
		validatedSecretVersionsCache:            validatedSecretVersionsCache,
		ldapDialer:                              ldapDialer,
		client:                                  client,
		activeDirectoryIdentityProviderInformer: activeDirectoryIdentityProviderInformer,
		secretInformer:                          secretInformer,
	}
	return controllerlib.New(
		controllerlib.Config{Name: activeDirectoryControllerName, Syncer: &c},
		withInformer(
			activeDirectoryIdentityProviderInformer,
			pinnipedcontroller.MatchAnythingFilter(pinnipedcontroller.SingletonQueue()),
			controllerlib.InformerOption{},
		),
		withInformer(
			secretInformer,
			pinnipedcontroller.MatchAnySecretOfTypeFilter(upstreamwatchers.LDAPBindAccountSecretType, pinnipedcontroller.SingletonQueue()),
			controllerlib.InformerOption{},
		),
	)
}

// Sync implements controllerlib.Syncer.
func (c *activeDirectoryWatcherController) Sync(ctx controllerlib.Context) error {
	actualUpstreams, err := c.activeDirectoryIdentityProviderInformer.Lister().List(labels.Everything())
	if err != nil {
		return fmt.Errorf("failed to list ActiveDirectoryIdentityProviders: %w", err)
	}

	requeue := false
	validatedUpstreams := make([]provider.UpstreamLDAPIdentityProviderI, 0, len(actualUpstreams))
	for _, upstream := range actualUpstreams {
		valid, requestedRequeue := c.validateUpstream(ctx.Context, upstream)
		if valid != nil {
			validatedUpstreams = append(validatedUpstreams, valid)
		}
		if requestedRequeue {
			requeue = true
		}
	}

	c.cache.SetActiveDirectoryIdentityProviders(validatedUpstreams)

	if requeue {
		return controllerlib.ErrSyntheticRequeue
	}
	return nil
}

func (c *activeDirectoryWatcherController) validateUpstream(ctx context.Context, upstream *v1alpha1.ActiveDirectoryIdentityProvider) (p provider.UpstreamLDAPIdentityProviderI, requeue bool) {
	spec := upstream.Spec

	usernameAttribute := spec.UserSearch.Attributes.Username
	if len(usernameAttribute) == 0 {
		usernameAttribute = defaultActiveDirectoryUsernameAttributeName
	}
	uidAttribute := spec.UserSearch.Attributes.UID
	if len(uidAttribute) == 0 {
		uidAttribute = defaultActiveDirectoryUIDAttributeName
	}

	config := &upstreamldap.ProviderConfig{
		Name: upstream.Name,
		Host: spec.Host,
		UserSearch: upstreamldap.UserSearchConfig{
			Base:              spec.UserSearch.Base,
			Filter:            spec.UserSearch.Filter,
			UsernameAttribute: usernameAttribute,
			UIDAttribute:      uidAttribute,
		},
		GroupSearch: upstreamldap.GroupSearchConfig{
			Base:               spec.GroupSearch.Base,
			Filter:             spec.GroupSearch.Filter,
			GroupNameAttribute: spec.GroupSearch.Attributes.GroupName,
		},
		Dialer: c.ldapDialer,
	}

	conditions := upstreamwatchers.ValidateGenericLDAP(ctx, &activeDirectoryUpstreamGenericLDAPImpl{*upstream}, c.secretInformer, c.validatedSecretVersionsCache, config)

	c.updateStatus(ctx, upstream, conditions.Conditions())

	return upstreamwatchers.EvaluateConditions(conditions, config)
}

func (c *activeDirectoryWatcherController) updateStatus(ctx context.Context, upstream *v1alpha1.ActiveDirectoryIdentityProvider, conditions []*v1alpha1.Condition) {
	log := klogr.New().WithValues("namespace", upstream.Namespace, "name", upstream.Name)
	updated := upstream.DeepCopy()

	hadErrorCondition := conditionsutil.Merge(conditions, upstream.Generation, &updated.Status.Conditions, log)

	updated.Status.Phase = v1alpha1.ActiveDirectoryPhaseReady
	if hadErrorCondition {
		updated.Status.Phase = v1alpha1.ActiveDirectoryPhaseError
	}

	if equality.Semantic.DeepEqual(upstream, updated) {
		return // nothing to update
	}

	_, err := c.client.
		IDPV1alpha1().
		ActiveDirectoryIdentityProviders(upstream.Namespace).
		UpdateStatus(ctx, updated, metav1.UpdateOptions{})
	if err != nil {
		log.Error(err, "failed to update status")
	}
}