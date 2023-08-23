/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricCACRLApplyConfiguration represents an declarative configuration of the FabricCACRL type for use
// with apply.
type FabricCACRLApplyConfiguration struct {
	Expiry *string `json:"expiry,omitempty"`
}

// FabricCACRLApplyConfiguration constructs an declarative configuration of the FabricCACRL type for use with
// apply.
func FabricCACRL() *FabricCACRLApplyConfiguration {
	return &FabricCACRLApplyConfiguration{}
}

// WithExpiry sets the Expiry field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Expiry field is set to the value of the last call.
func (b *FabricCACRLApplyConfiguration) WithExpiry(value string) *FabricCACRLApplyConfiguration {
	b.Expiry = &value
	return b
}
