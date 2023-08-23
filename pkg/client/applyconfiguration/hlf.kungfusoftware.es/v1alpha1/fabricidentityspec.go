/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricIdentitySpecApplyConfiguration represents an declarative configuration of the FabricIdentitySpec type for use
// with apply.
type FabricIdentitySpecApplyConfiguration struct {
	Cahost       *string                                   `json:"cahost,omitempty"`
	Caname       *string                                   `json:"caname,omitempty"`
	Caport       *int                                      `json:"caport,omitempty"`
	Catls        *CatlsApplyConfiguration                  `json:"catls,omitempty"`
	Enrollid     *string                                   `json:"enrollid,omitempty"`
	Enrollsecret *string                                   `json:"enrollsecret,omitempty"`
	MSPID        *string                                   `json:"mspid,omitempty"`
	Register     *FabricIdentityRegisterApplyConfiguration `json:"register,omitempty"`
}

// FabricIdentitySpecApplyConfiguration constructs an declarative configuration of the FabricIdentitySpec type for use with
// apply.
func FabricIdentitySpec() *FabricIdentitySpecApplyConfiguration {
	return &FabricIdentitySpecApplyConfiguration{}
}

// WithCahost sets the Cahost field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Cahost field is set to the value of the last call.
func (b *FabricIdentitySpecApplyConfiguration) WithCahost(value string) *FabricIdentitySpecApplyConfiguration {
	b.Cahost = &value
	return b
}

// WithCaname sets the Caname field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Caname field is set to the value of the last call.
func (b *FabricIdentitySpecApplyConfiguration) WithCaname(value string) *FabricIdentitySpecApplyConfiguration {
	b.Caname = &value
	return b
}

// WithCaport sets the Caport field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Caport field is set to the value of the last call.
func (b *FabricIdentitySpecApplyConfiguration) WithCaport(value int) *FabricIdentitySpecApplyConfiguration {
	b.Caport = &value
	return b
}

// WithCatls sets the Catls field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Catls field is set to the value of the last call.
func (b *FabricIdentitySpecApplyConfiguration) WithCatls(value *CatlsApplyConfiguration) *FabricIdentitySpecApplyConfiguration {
	b.Catls = value
	return b
}

// WithEnrollid sets the Enrollid field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Enrollid field is set to the value of the last call.
func (b *FabricIdentitySpecApplyConfiguration) WithEnrollid(value string) *FabricIdentitySpecApplyConfiguration {
	b.Enrollid = &value
	return b
}

// WithEnrollsecret sets the Enrollsecret field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Enrollsecret field is set to the value of the last call.
func (b *FabricIdentitySpecApplyConfiguration) WithEnrollsecret(value string) *FabricIdentitySpecApplyConfiguration {
	b.Enrollsecret = &value
	return b
}

// WithMSPID sets the MSPID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MSPID field is set to the value of the last call.
func (b *FabricIdentitySpecApplyConfiguration) WithMSPID(value string) *FabricIdentitySpecApplyConfiguration {
	b.MSPID = &value
	return b
}

// WithRegister sets the Register field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Register field is set to the value of the last call.
func (b *FabricIdentitySpecApplyConfiguration) WithRegister(value *FabricIdentityRegisterApplyConfiguration) *FabricIdentitySpecApplyConfiguration {
	b.Register = value
	return b
}
