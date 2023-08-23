/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricOperatorAPINetworkConfigApplyConfiguration represents an declarative configuration of the FabricOperatorAPINetworkConfig type for use
// with apply.
type FabricOperatorAPINetworkConfigApplyConfiguration struct {
	SecretName *string `json:"secretName,omitempty"`
	Key        *string `json:"key,omitempty"`
}

// FabricOperatorAPINetworkConfigApplyConfiguration constructs an declarative configuration of the FabricOperatorAPINetworkConfig type for use with
// apply.
func FabricOperatorAPINetworkConfig() *FabricOperatorAPINetworkConfigApplyConfiguration {
	return &FabricOperatorAPINetworkConfigApplyConfiguration{}
}

// WithSecretName sets the SecretName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SecretName field is set to the value of the last call.
func (b *FabricOperatorAPINetworkConfigApplyConfiguration) WithSecretName(value string) *FabricOperatorAPINetworkConfigApplyConfiguration {
	b.SecretName = &value
	return b
}

// WithKey sets the Key field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Key field is set to the value of the last call.
func (b *FabricOperatorAPINetworkConfigApplyConfiguration) WithKey(value string) *FabricOperatorAPINetworkConfigApplyConfiguration {
	b.Key = &value
	return b
}
