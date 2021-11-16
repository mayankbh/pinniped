// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package supervisor

import (
	"io/ioutil"
	"os"
	"testing"

	"k8s.io/utils/pointer"

	"github.com/stretchr/testify/require"

	"go.pinniped.dev/internal/here"
)

func TestFromPath(t *testing.T) {
	tests := []struct {
		name       string
		yaml       string
		wantConfig *Config
		wantError  string
	}{
		{
			name: "Happy",
			yaml: here.Doc(`
				---
				apiGroupSuffix: some.suffix.com
				labels:
				  myLabelKey1: myLabelValue1
				  myLabelKey2: myLabelValue2
				names:
				  defaultTLSCertificateSecret: my-secret-name
				listenPort: 12345
			`),
			wantConfig: &Config{
				APIGroupSuffix: pointer.StringPtr("some.suffix.com"),
				Labels: map[string]string{
					"myLabelKey1": "myLabelValue1",
					"myLabelKey2": "myLabelValue2",
				},
				NamesConfig: NamesConfigSpec{
					DefaultTLSCertificateSecret: "my-secret-name",
				},
				ListenPort: pointer.IntPtr(12345),
			},
		},
		{
			name: "When only the required fields are present, causes other fields to be defaulted",
			yaml: here.Doc(`
				---
				names:
				  defaultTLSCertificateSecret: my-secret-name
			`),
			wantConfig: &Config{
				APIGroupSuffix: pointer.StringPtr("pinniped.dev"),
				Labels:         map[string]string{},
				NamesConfig: NamesConfigSpec{
					DefaultTLSCertificateSecret: "my-secret-name",
				},
				ListenPort: pointer.IntPtr(8443),
			},
		},
		{
			name: "Missing defaultTLSCertificateSecret name",
			yaml: here.Doc(`
				---
			`),
			wantError: "validate names: missing required names: defaultTLSCertificateSecret",
		},
		{
			name: "apiGroupSuffix is prefixed with '.'",
			yaml: here.Doc(`
				---
				apiGroupSuffix: .starts.with.dot
				names:
				  defaultTLSCertificateSecret: my-secret-name
			`),
			wantError: "validate apiGroupSuffix: a lowercase RFC 1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character (e.g. 'example.com', regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*')",
		},
		{
			name: "when an invalid port is provided",
			yaml: here.Doc(`
				---
				names:
				  defaultTLSCertificateSecret: my-secret-name
				listenPort: 2000000
			`),
			wantError: "validate listenPort: must be between 1 and 65535, inclusive",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			// Write yaml to temp file
			f, err := ioutil.TempFile("", "pinniped-test-config-yaml-*")
			require.NoError(t, err)
			defer func() {
				err := os.Remove(f.Name())
				require.NoError(t, err)
			}()
			_, err = f.WriteString(test.yaml)
			require.NoError(t, err)
			err = f.Close()
			require.NoError(t, err)

			// Test FromPath()
			config, err := FromPath(f.Name())

			if test.wantError != "" {
				require.EqualError(t, err, test.wantError)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.wantConfig, config)
			}
		})
	}
}
