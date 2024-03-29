/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */
// Code generated by client-gen. DO NOT EDIT.

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

// FakeFabricChaincodeTemplates implements FabricChaincodeTemplateInterface
type FakeFabricChaincodeTemplates struct {
	Fake *FakeHlfV1alpha1
	ns   string
}

var fabricchaincodetemplatesResource = v1alpha1.SchemeGroupVersion.WithResource("fabricchaincodetemplates")

var fabricchaincodetemplatesKind = v1alpha1.SchemeGroupVersion.WithKind("FabricChaincodeTemplate")

// Get takes name of the fabricChaincodeTemplate, and returns the corresponding fabricChaincodeTemplate object, and an error if there is any.
func (c *FakeFabricChaincodeTemplates) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricChaincodeTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(fabricchaincodetemplatesResource, c.ns, name), &v1alpha1.FabricChaincodeTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincodeTemplate), err
}

// List takes label and field selectors, and returns the list of FabricChaincodeTemplates that match those selectors.
func (c *FakeFabricChaincodeTemplates) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricChaincodeTemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(fabricchaincodetemplatesResource, fabricchaincodetemplatesKind, c.ns, opts), &v1alpha1.FabricChaincodeTemplateList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FabricChaincodeTemplateList{ListMeta: obj.(*v1alpha1.FabricChaincodeTemplateList).ListMeta}
	for _, item := range obj.(*v1alpha1.FabricChaincodeTemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested fabricChaincodeTemplates.
func (c *FakeFabricChaincodeTemplates) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(fabricchaincodetemplatesResource, c.ns, opts))

}

// Create takes the representation of a fabricChaincodeTemplate and creates it.  Returns the server's representation of the fabricChaincodeTemplate, and an error, if there is any.
func (c *FakeFabricChaincodeTemplates) Create(ctx context.Context, fabricChaincodeTemplate *v1alpha1.FabricChaincodeTemplate, opts v1.CreateOptions) (result *v1alpha1.FabricChaincodeTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(fabricchaincodetemplatesResource, c.ns, fabricChaincodeTemplate), &v1alpha1.FabricChaincodeTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincodeTemplate), err
}

// Update takes the representation of a fabricChaincodeTemplate and updates it. Returns the server's representation of the fabricChaincodeTemplate, and an error, if there is any.
func (c *FakeFabricChaincodeTemplates) Update(ctx context.Context, fabricChaincodeTemplate *v1alpha1.FabricChaincodeTemplate, opts v1.UpdateOptions) (result *v1alpha1.FabricChaincodeTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(fabricchaincodetemplatesResource, c.ns, fabricChaincodeTemplate), &v1alpha1.FabricChaincodeTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincodeTemplate), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFabricChaincodeTemplates) UpdateStatus(ctx context.Context, fabricChaincodeTemplate *v1alpha1.FabricChaincodeTemplate, opts v1.UpdateOptions) (*v1alpha1.FabricChaincodeTemplate, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(fabricchaincodetemplatesResource, "status", c.ns, fabricChaincodeTemplate), &v1alpha1.FabricChaincodeTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincodeTemplate), err
}

// Delete takes name of the fabricChaincodeTemplate and deletes it. Returns an error if one occurs.
func (c *FakeFabricChaincodeTemplates) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(fabricchaincodetemplatesResource, c.ns, name, opts), &v1alpha1.FabricChaincodeTemplate{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFabricChaincodeTemplates) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(fabricchaincodetemplatesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.FabricChaincodeTemplateList{})
	return err
}

// Patch applies the patch and returns the patched fabricChaincodeTemplate.
func (c *FakeFabricChaincodeTemplates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricChaincodeTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricchaincodetemplatesResource, c.ns, name, pt, data, subresources...), &v1alpha1.FabricChaincodeTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincodeTemplate), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricChaincodeTemplate.
func (c *FakeFabricChaincodeTemplates) Apply(ctx context.Context, fabricChaincodeTemplate *hlfkungfusoftwareesv1alpha1.FabricChaincodeTemplateApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricChaincodeTemplate, err error) {
	if fabricChaincodeTemplate == nil {
		return nil, fmt.Errorf("fabricChaincodeTemplate provided to Apply must not be nil")
	}
	data, err := json.Marshal(fabricChaincodeTemplate)
	if err != nil {
		return nil, err
	}
	name := fabricChaincodeTemplate.Name
	if name == nil {
		return nil, fmt.Errorf("fabricChaincodeTemplate.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricchaincodetemplatesResource, c.ns, *name, types.ApplyPatchType, data), &v1alpha1.FabricChaincodeTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincodeTemplate), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeFabricChaincodeTemplates) ApplyStatus(ctx context.Context, fabricChaincodeTemplate *hlfkungfusoftwareesv1alpha1.FabricChaincodeTemplateApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricChaincodeTemplate, err error) {
	if fabricChaincodeTemplate == nil {
		return nil, fmt.Errorf("fabricChaincodeTemplate provided to Apply must not be nil")
	}
	data, err := json.Marshal(fabricChaincodeTemplate)
	if err != nil {
		return nil, err
	}
	name := fabricChaincodeTemplate.Name
	if name == nil {
		return nil, fmt.Errorf("fabricChaincodeTemplate.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricchaincodetemplatesResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1alpha1.FabricChaincodeTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricChaincodeTemplate), err
}
