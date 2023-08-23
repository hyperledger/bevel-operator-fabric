/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricMainChannelExternalOrdererNodeApplyConfiguration represents an declarative configuration of the FabricMainChannelExternalOrdererNode type for use
// with apply.
type FabricMainChannelExternalOrdererNodeApplyConfiguration struct {
	Host      *string `json:"host,omitempty"`
	AdminPort *int    `json:"port,omitempty"`
}

// FabricMainChannelExternalOrdererNodeApplyConfiguration constructs an declarative configuration of the FabricMainChannelExternalOrdererNode type for use with
// apply.
func FabricMainChannelExternalOrdererNode() *FabricMainChannelExternalOrdererNodeApplyConfiguration {
	return &FabricMainChannelExternalOrdererNodeApplyConfiguration{}
}

// WithHost sets the Host field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Host field is set to the value of the last call.
func (b *FabricMainChannelExternalOrdererNodeApplyConfiguration) WithHost(value string) *FabricMainChannelExternalOrdererNodeApplyConfiguration {
	b.Host = &value
	return b
}

// WithAdminPort sets the AdminPort field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AdminPort field is set to the value of the last call.
func (b *FabricMainChannelExternalOrdererNodeApplyConfiguration) WithAdminPort(value int) *FabricMainChannelExternalOrdererNodeApplyConfiguration {
	b.AdminPort = &value
	return b
}
