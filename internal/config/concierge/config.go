// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package concierge contains functionality to load/store Config's from/to
// some source.
package concierge

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"

	"go.pinniped.dev/internal/constable"
	"go.pinniped.dev/internal/groupsuffix"
	"go.pinniped.dev/internal/plog"
)

const (
	aboutAYear   = 60 * 60 * 24 * 365
	about9Months = 60 * 60 * 24 * 30 * 9

	defaultConciergeListenPort = 8443
)

// FromPath loads an Config from a provided local file path, inserts any
// defaults (from the Config documentation), and verifies that the config is
// valid (per the Config documentation).
//
// Note! The Config file should contain base64-encoded WebhookCABundle data.
// This function will decode that base64-encoded data to PEM bytes to be stored
// in the Config.
func FromPath(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("decode yaml: %w", err)
	}

	maybeSetAPIDefaults(&config.APIConfig)
	maybeSetAPIGroupSuffixDefault(&config.APIGroupSuffix)
	maybeSetKubeCertAgentDefaults(&config.KubeCertAgentConfig)
	maybeSetListenPort(&config.ListenPort)

	if err := validateAPI(&config.APIConfig); err != nil {
		return nil, fmt.Errorf("validate api: %w", err)
	}

	if err := validateAPIGroupSuffix(*config.APIGroupSuffix); err != nil {
		return nil, fmt.Errorf("validate apiGroupSuffix: %w", err)
	}

	if err := validateNames(&config.NamesConfig); err != nil {
		return nil, fmt.Errorf("validate names: %w", err)
	}

	if err := validatePort(*config.ListenPort); err != nil {
		return nil, fmt.Errorf("validate listenPort: %w", err)
	}

	if err := plog.ValidateAndSetLogLevelGlobally(config.LogLevel); err != nil {
		return nil, fmt.Errorf("validate log level: %w", err)
	}

	if config.Labels == nil {
		config.Labels = make(map[string]string)
	}

	return &config, nil
}

func maybeSetAPIDefaults(apiConfig *APIConfigSpec) {
	if apiConfig.ServingCertificateConfig.DurationSeconds == nil {
		apiConfig.ServingCertificateConfig.DurationSeconds = pointer.Int64Ptr(aboutAYear)
	}

	if apiConfig.ServingCertificateConfig.RenewBeforeSeconds == nil {
		apiConfig.ServingCertificateConfig.RenewBeforeSeconds = pointer.Int64Ptr(about9Months)
	}
}

func maybeSetAPIGroupSuffixDefault(apiGroupSuffix **string) {
	if *apiGroupSuffix == nil {
		*apiGroupSuffix = pointer.StringPtr(groupsuffix.PinnipedDefaultSuffix)
	}
}

func maybeSetKubeCertAgentDefaults(cfg *KubeCertAgentSpec) {
	if cfg.NamePrefix == nil {
		cfg.NamePrefix = pointer.StringPtr("pinniped-kube-cert-agent-")
	}

	if cfg.Image == nil {
		cfg.Image = pointer.StringPtr("debian:latest")
	}
}

func maybeSetListenPort(listenPort **int) {
	if *listenPort == nil {
		*listenPort = pointer.IntPtr(defaultConciergeListenPort)
	}
}

func validateNames(names *NamesConfigSpec) error {
	missingNames := []string{}
	if names == nil {
		names = &NamesConfigSpec{}
	}
	if names.ServingCertificateSecret == "" {
		missingNames = append(missingNames, "servingCertificateSecret")
	}
	if names.CredentialIssuer == "" {
		missingNames = append(missingNames, "credentialIssuer")
	}
	if names.APIService == "" {
		missingNames = append(missingNames, "apiService")
	}
	if names.ImpersonationLoadBalancerService == "" {
		missingNames = append(missingNames, "impersonationLoadBalancerService")
	}
	if names.ImpersonationClusterIPService == "" {
		missingNames = append(missingNames, "impersonationClusterIPService")
	}
	if names.ImpersonationTLSCertificateSecret == "" {
		missingNames = append(missingNames, "impersonationTLSCertificateSecret")
	}
	if names.ImpersonationCACertificateSecret == "" {
		missingNames = append(missingNames, "impersonationCACertificateSecret")
	}
	if names.ImpersonationSignerSecret == "" {
		missingNames = append(missingNames, "impersonationSignerSecret")
	}
	if names.AgentServiceAccount == "" {
		missingNames = append(missingNames, "agentServiceAccount")
	}
	if len(missingNames) > 0 {
		return constable.Error("missing required names: " + strings.Join(missingNames, ", "))
	}
	return nil
}

func validateAPI(apiConfig *APIConfigSpec) error {
	if *apiConfig.ServingCertificateConfig.DurationSeconds < *apiConfig.ServingCertificateConfig.RenewBeforeSeconds {
		return constable.Error("durationSeconds cannot be smaller than renewBeforeSeconds")
	}

	if *apiConfig.ServingCertificateConfig.RenewBeforeSeconds <= 0 {
		return constable.Error("renewBefore must be positive")
	}

	return nil
}

func validateAPIGroupSuffix(apiGroupSuffix string) error {
	return groupsuffix.Validate(apiGroupSuffix)
}

func validatePort(listenPort int) error {
	if result := validation.IsValidPortNum(listenPort); result != nil {
		return errors.New(strings.Join(result, " "))
	}
	return nil
}
