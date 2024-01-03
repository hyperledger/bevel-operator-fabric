/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricCACryptoApplyConfiguration represents an declarative configuration of the FabricCACrypto type for use
// with apply.
type FabricCACryptoApplyConfiguration struct {
	Key   *string `json:"key,omitempty"`
	Cert  *string `json:"cert,omitempty"`
	Chain *string `json:"chain,omitempty"`
}

// FabricCACryptoApplyConfiguration constructs an declarative configuration of the FabricCACrypto type for use with
// apply.
func FabricCACrypto() *FabricCACryptoApplyConfiguration {
	return &FabricCACryptoApplyConfiguration{}
}

// WithKey sets the Key field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Key field is set to the value of the last call.
func (b *FabricCACryptoApplyConfiguration) WithKey(value string) *FabricCACryptoApplyConfiguration {
	b.Key = &value
	return b
}

// WithCert sets the Cert field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Cert field is set to the value of the last call.
func (b *FabricCACryptoApplyConfiguration) WithCert(value string) *FabricCACryptoApplyConfiguration {
	b.Cert = &value
	return b
}

// WithChain sets the Chain field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Chain field is set to the value of the last call.
func (b *FabricCACryptoApplyConfiguration) WithChain(value string) *FabricCACryptoApplyConfiguration {
	b.Chain = &value
	return b
}
