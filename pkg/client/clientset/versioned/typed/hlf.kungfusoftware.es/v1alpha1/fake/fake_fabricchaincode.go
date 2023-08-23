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

// FakeFabricChaincodes implements FabricChaincodeInterface
type FakeFabricChaincodes struct {
	Fake *FakeHlfV1alpha1
	ns   string
}

var fabricchaincodesResource = v1alpha1.SchemeGroupVersion.WithResource("fabricchaincodes")

var fabricchaincodesKind = v1alpha1.SchemeGroupVersion.WithKind("FabricChaincode")

// Get takes name of the fabricChaincode, and returns the corresponding fabricChaincode object, and an error if there is any.
func (c *FakeFabricChaincodes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricChaincode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(fabricchaincodesResource, c.ns, name), &v1alpha1.FabricChaincode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincode), err
}

// List takes label and field selectors, and returns the list of FabricChaincodes that match those selectors.
func (c *FakeFabricChaincodes) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricChaincodeList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(fabricchaincodesResource, fabricchaincodesKind, c.ns, opts), &v1alpha1.FabricChaincodeList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FabricChaincodeList{ListMeta: obj.(*v1alpha1.FabricChaincodeList).ListMeta}
	for _, item := range obj.(*v1alpha1.FabricChaincodeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested fabricChaincodes.
func (c *FakeFabricChaincodes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(fabricchaincodesResource, c.ns, opts))

}

// Create takes the representation of a fabricChaincode and creates it.  Returns the server's representation of the fabricChaincode, and an error, if there is any.
func (c *FakeFabricChaincodes) Create(ctx context.Context, fabricChaincode *v1alpha1.FabricChaincode, opts v1.CreateOptions) (result *v1alpha1.FabricChaincode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(fabricchaincodesResource, c.ns, fabricChaincode), &v1alpha1.FabricChaincode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincode), err
}

// Update takes the representation of a fabricChaincode and updates it. Returns the server's representation of the fabricChaincode, and an error, if there is any.
func (c *FakeFabricChaincodes) Update(ctx context.Context, fabricChaincode *v1alpha1.FabricChaincode, opts v1.UpdateOptions) (result *v1alpha1.FabricChaincode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(fabricchaincodesResource, c.ns, fabricChaincode), &v1alpha1.FabricChaincode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincode), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFabricChaincodes) UpdateStatus(ctx context.Context, fabricChaincode *v1alpha1.FabricChaincode, opts v1.UpdateOptions) (*v1alpha1.FabricChaincode, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(fabricchaincodesResource, "status", c.ns, fabricChaincode), &v1alpha1.FabricChaincode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincode), err
}

// Delete takes name of the fabricChaincode and deletes it. Returns an error if one occurs.
func (c *FakeFabricChaincodes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(fabricchaincodesResource, c.ns, name, opts), &v1alpha1.FabricChaincode{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFabricChaincodes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(fabricchaincodesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.FabricChaincodeList{})
	return err
}

// Patch applies the patch and returns the patched fabricChaincode.
func (c *FakeFabricChaincodes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricChaincode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricchaincodesResource, c.ns, name, pt, data, subresources...), &v1alpha1.FabricChaincode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincode), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricChaincode.
func (c *FakeFabricChaincodes) Apply(ctx context.Context, fabricChaincode *hlfkungfusoftwareesv1alpha1.FabricChaincodeApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricChaincode, err error) {
	if fabricChaincode == nil {
		return nil, fmt.Errorf("fabricChaincode provided to Apply must not be nil")
	}
	data, err := json.Marshal(fabricChaincode)
	if err != nil {
		return nil, err
	}
	name := fabricChaincode.Name
	if name == nil {
		return nil, fmt.Errorf("fabricChaincode.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricchaincodesResource, c.ns, *name, types.ApplyPatchType, data), &v1alpha1.FabricChaincode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincode), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeFabricChaincodes) ApplyStatus(ctx context.Context, fabricChaincode *hlfkungfusoftwareesv1alpha1.FabricChaincodeApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricChaincode, err error) {
	if fabricChaincode == nil {
		return nil, fmt.Errorf("fabricChaincode provided to Apply must not be nil")
	}
	data, err := json.Marshal(fabricChaincode)
	if err != nil {
		return nil, err
	}
	name := fabricChaincode.Name
	if name == nil {
		return nil, fmt.Errorf("fabricChaincode.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricchaincodesResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1alpha1.FabricChaincode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincode), err
}
