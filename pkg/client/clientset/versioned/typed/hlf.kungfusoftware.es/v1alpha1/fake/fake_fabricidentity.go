/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	hlfkungfusoftwareesv1alpha1 "github.com/kfsoftware/hlf-operator/pkg/client/applyconfiguration/hlf.kungfusoftware.es/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFabricIdentities implements FabricIdentityInterface
type FakeFabricIdentities struct {
	Fake *FakeHlfV1alpha1
	ns   string
}

var fabricidentitiesResource = v1alpha1.SchemeGroupVersion.WithResource("fabricidentities")

var fabricidentitiesKind = v1alpha1.SchemeGroupVersion.WithKind("FabricIdentity")

// Get takes name of the fabricIdentity, and returns the corresponding fabricIdentity object, and an error if there is any.
func (c *FakeFabricIdentities) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricIdentity, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(fabricidentitiesResource, c.ns, name), &v1alpha1.FabricIdentity{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricIdentity), err
}

// List takes label and field selectors, and returns the list of FabricIdentities that match those selectors.
func (c *FakeFabricIdentities) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricIdentityList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(fabricidentitiesResource, fabricidentitiesKind, c.ns, opts), &v1alpha1.FabricIdentityList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FabricIdentityList{ListMeta: obj.(*v1alpha1.FabricIdentityList).ListMeta}
	for _, item := range obj.(*v1alpha1.FabricIdentityList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested fabricIdentities.
func (c *FakeFabricIdentities) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(fabricidentitiesResource, c.ns, opts))

}

// Create takes the representation of a fabricIdentity and creates it.  Returns the server's representation of the fabricIdentity, and an error, if there is any.
func (c *FakeFabricIdentities) Create(ctx context.Context, fabricIdentity *v1alpha1.FabricIdentity, opts v1.CreateOptions) (result *v1alpha1.FabricIdentity, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(fabricidentitiesResource, c.ns, fabricIdentity), &v1alpha1.FabricIdentity{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricIdentity), err
}

// Update takes the representation of a fabricIdentity and updates it. Returns the server's representation of the fabricIdentity, and an error, if there is any.
func (c *FakeFabricIdentities) Update(ctx context.Context, fabricIdentity *v1alpha1.FabricIdentity, opts v1.UpdateOptions) (result *v1alpha1.FabricIdentity, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(fabricidentitiesResource, c.ns, fabricIdentity), &v1alpha1.FabricIdentity{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricIdentity), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFabricIdentities) UpdateStatus(ctx context.Context, fabricIdentity *v1alpha1.FabricIdentity, opts v1.UpdateOptions) (*v1alpha1.FabricIdentity, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(fabricidentitiesResource, "status", c.ns, fabricIdentity), &v1alpha1.FabricIdentity{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricIdentity), err
}

// Delete takes name of the fabricIdentity and deletes it. Returns an error if one occurs.
func (c *FakeFabricIdentities) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(fabricidentitiesResource, c.ns, name, opts), &v1alpha1.FabricIdentity{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFabricIdentities) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(fabricidentitiesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.FabricIdentityList{})
	return err
}

// Patch applies the patch and returns the patched fabricIdentity.
func (c *FakeFabricIdentities) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricIdentity, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricidentitiesResource, c.ns, name, pt, data, subresources...), &v1alpha1.FabricIdentity{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricIdentity), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricIdentity.
func (c *FakeFabricIdentities) Apply(ctx context.Context, fabricIdentity *hlfkungfusoftwareesv1alpha1.FabricIdentityApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricIdentity, err error) {
	if fabricIdentity == nil {
		return nil, fmt.Errorf("fabricIdentity provided to Apply must not be nil")
	}
	data, err := json.Marshal(fabricIdentity)
	if err != nil {
		return nil, err
	}
	name := fabricIdentity.Name
	if name == nil {
		return nil, fmt.Errorf("fabricIdentity.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricidentitiesResource, c.ns, *name, types.ApplyPatchType, data), &v1alpha1.FabricIdentity{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricIdentity), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeFabricIdentities) ApplyStatus(ctx context.Context, fabricIdentity *hlfkungfusoftwareesv1alpha1.FabricIdentityApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricIdentity, err error) {
	if fabricIdentity == nil {
		return nil, fmt.Errorf("fabricIdentity provided to Apply must not be nil")
	}
	data, err := json.Marshal(fabricIdentity)
	if err != nil {
		return nil, err
	}
	name := fabricIdentity.Name
	if name == nil {
		return nil, fmt.Errorf("fabricIdentity.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricidentitiesResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1alpha1.FabricIdentity{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricIdentity), err
}
