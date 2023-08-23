/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricCAClientAuthApplyConfiguration represents an declarative configuration of the FabricCAClientAuth type for use
// with apply.
type FabricCAClientAuthApplyConfiguration struct {
	Type     *string  `json:"type,omitempty"`
	CertFile []string `json:"cert_file,omitempty"`
}

// FabricCAClientAuthApplyConfiguration constructs an declarative configuration of the FabricCAClientAuth type for use with
// apply.
func FabricCAClientAuth() *FabricCAClientAuthApplyConfiguration {
	return &FabricCAClientAuthApplyConfiguration{}
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *FabricCAClientAuthApplyConfiguration) WithType(value string) *FabricCAClientAuthApplyConfiguration {
	b.Type = &value
	return b
}

// WithCertFile adds the given value to the CertFile field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the CertFile field.
func (b *FabricCAClientAuthApplyConfiguration) WithCertFile(values ...string) *FabricCAClientAuthApplyConfiguration {
	for i := range values {
		b.CertFile = append(b.CertFile, values[i])
	}
	return b
}
