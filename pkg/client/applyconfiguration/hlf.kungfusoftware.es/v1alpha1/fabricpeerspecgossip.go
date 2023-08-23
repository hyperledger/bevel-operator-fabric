/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

// FabricPeerSpecGossipApplyConfiguration represents an declarative configuration of the FabricPeerSpecGossip type for use
// with apply.
type FabricPeerSpecGossipApplyConfiguration struct {
	ExternalEndpoint  *string `json:"externalEndpoint,omitempty"`
	Bootstrap         *string `json:"bootstrap,omitempty"`
	Endpoint          *string `json:"endpoint,omitempty"`
	UseLeaderElection *bool   `json:"useLeaderElection,omitempty"`
	OrgLeader         *bool   `json:"orgLeader,omitempty"`
}

// FabricPeerSpecGossipApplyConfiguration constructs an declarative configuration of the FabricPeerSpecGossip type for use with
// apply.
func FabricPeerSpecGossip() *FabricPeerSpecGossipApplyConfiguration {
	return &FabricPeerSpecGossipApplyConfiguration{}
}

// WithExternalEndpoint sets the ExternalEndpoint field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ExternalEndpoint field is set to the value of the last call.
func (b *FabricPeerSpecGossipApplyConfiguration) WithExternalEndpoint(value string) *FabricPeerSpecGossipApplyConfiguration {
	b.ExternalEndpoint = &value
	return b
}

// WithBootstrap sets the Bootstrap field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Bootstrap field is set to the value of the last call.
func (b *FabricPeerSpecGossipApplyConfiguration) WithBootstrap(value string) *FabricPeerSpecGossipApplyConfiguration {
	b.Bootstrap = &value
	return b
}

// WithEndpoint sets the Endpoint field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Endpoint field is set to the value of the last call.
func (b *FabricPeerSpecGossipApplyConfiguration) WithEndpoint(value string) *FabricPeerSpecGossipApplyConfiguration {
	b.Endpoint = &value
	return b
}

// WithUseLeaderElection sets the UseLeaderElection field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UseLeaderElection field is set to the value of the last call.
func (b *FabricPeerSpecGossipApplyConfiguration) WithUseLeaderElection(value bool) *FabricPeerSpecGossipApplyConfiguration {
	b.UseLeaderElection = &value
	return b
}

// WithOrgLeader sets the OrgLeader field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the OrgLeader field is set to the value of the last call.
func (b *FabricPeerSpecGossipApplyConfiguration) WithOrgLeader(value bool) *FabricPeerSpecGossipApplyConfiguration {
	b.OrgLeader = &value
	return b
}
