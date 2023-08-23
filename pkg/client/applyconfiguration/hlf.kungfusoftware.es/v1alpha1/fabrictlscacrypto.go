/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricTLSCACryptoApplyConfiguration represents an declarative configuration of the FabricTLSCACrypto type for use
// with apply.
type FabricTLSCACryptoApplyConfiguration struct {
	Key        *string                               `json:"key,omitempty"`
	Cert       *string                               `json:"cert,omitempty"`
	ClientAuth *FabricCAClientAuthApplyConfiguration `json:"clientAuth,omitempty"`
}

// FabricTLSCACryptoApplyConfiguration constructs an declarative configuration of the FabricTLSCACrypto type for use with
// apply.
func FabricTLSCACrypto() *FabricTLSCACryptoApplyConfiguration {
	return &FabricTLSCACryptoApplyConfiguration{}
}

// WithKey sets the Key field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Key field is set to the value of the last call.
func (b *FabricTLSCACryptoApplyConfiguration) WithKey(value string) *FabricTLSCACryptoApplyConfiguration {
	b.Key = &value
	return b
}

// WithCert sets the Cert field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Cert field is set to the value of the last call.
func (b *FabricTLSCACryptoApplyConfiguration) WithCert(value string) *FabricTLSCACryptoApplyConfiguration {
	b.Cert = &value
	return b
}

// WithClientAuth sets the ClientAuth field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ClientAuth field is set to the value of the last call.
func (b *FabricTLSCACryptoApplyConfiguration) WithClientAuth(value *FabricCAClientAuthApplyConfiguration) *FabricTLSCACryptoApplyConfiguration {
	b.ClientAuth = value
	return b
}
