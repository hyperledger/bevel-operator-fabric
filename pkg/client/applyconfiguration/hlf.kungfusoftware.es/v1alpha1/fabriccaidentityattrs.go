// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// FabricCAIdentityAttrsApplyConfiguration represents an declarative configuration of the FabricCAIdentityAttrs type for use
// with apply.
type FabricCAIdentityAttrsApplyConfiguration struct {
	RegistrarRoles *string `json:"hf.Registrar.Roles,omitempty"`
	DelegateRoles  *string `json:"hf.Registrar.DelegateRoles,omitempty"`
	Attributes     *string `json:"hf.Registrar.Attributes,omitempty"`
	Revoker        *bool   `json:"hf.Revoker,omitempty"`
	IntermediateCA *bool   `json:"hf.IntermediateCA,omitempty"`
	GenCRL         *bool   `json:"hf.GenCRL,omitempty"`
	AffiliationMgr *bool   `json:"hf.AffiliationMgr,omitempty"`
}

// FabricCAIdentityAttrsApplyConfiguration constructs an declarative configuration of the FabricCAIdentityAttrs type for use with
// apply.
func FabricCAIdentityAttrs() *FabricCAIdentityAttrsApplyConfiguration {
	return &FabricCAIdentityAttrsApplyConfiguration{}
}

// WithRegistrarRoles sets the RegistrarRoles field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RegistrarRoles field is set to the value of the last call.
func (b *FabricCAIdentityAttrsApplyConfiguration) WithRegistrarRoles(value string) *FabricCAIdentityAttrsApplyConfiguration {
	b.RegistrarRoles = &value
	return b
}

// WithDelegateRoles sets the DelegateRoles field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DelegateRoles field is set to the value of the last call.
func (b *FabricCAIdentityAttrsApplyConfiguration) WithDelegateRoles(value string) *FabricCAIdentityAttrsApplyConfiguration {
	b.DelegateRoles = &value
	return b
}

// WithAttributes sets the Attributes field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Attributes field is set to the value of the last call.
func (b *FabricCAIdentityAttrsApplyConfiguration) WithAttributes(value string) *FabricCAIdentityAttrsApplyConfiguration {
	b.Attributes = &value
	return b
}

// WithRevoker sets the Revoker field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Revoker field is set to the value of the last call.
func (b *FabricCAIdentityAttrsApplyConfiguration) WithRevoker(value bool) *FabricCAIdentityAttrsApplyConfiguration {
	b.Revoker = &value
	return b
}

// WithIntermediateCA sets the IntermediateCA field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the IntermediateCA field is set to the value of the last call.
func (b *FabricCAIdentityAttrsApplyConfiguration) WithIntermediateCA(value bool) *FabricCAIdentityAttrsApplyConfiguration {
	b.IntermediateCA = &value
	return b
}

// WithGenCRL sets the GenCRL field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GenCRL field is set to the value of the last call.
func (b *FabricCAIdentityAttrsApplyConfiguration) WithGenCRL(value bool) *FabricCAIdentityAttrsApplyConfiguration {
	b.GenCRL = &value
	return b
}

// WithAffiliationMgr sets the AffiliationMgr field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AffiliationMgr field is set to the value of the last call.
func (b *FabricCAIdentityAttrsApplyConfiguration) WithAffiliationMgr(value bool) *FabricCAIdentityAttrsApplyConfiguration {
	b.AffiliationMgr = &value
	return b
}