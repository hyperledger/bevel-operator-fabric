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

// FakeFabricExplorers implements FabricExplorerInterface
type FakeFabricExplorers struct {
	Fake *FakeHlfV1alpha1
	ns   string
}

var fabricexplorersResource = v1alpha1.SchemeGroupVersion.WithResource("fabricexplorers")

var fabricexplorersKind = v1alpha1.SchemeGroupVersion.WithKind("FabricExplorer")

// Get takes name of the fabricExplorer, and returns the corresponding fabricExplorer object, and an error if there is any.
func (c *FakeFabricExplorers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricExplorer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(fabricexplorersResource, c.ns, name), &v1alpha1.FabricExplorer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricExplorer), err
}

// List takes label and field selectors, and returns the list of FabricExplorers that match those selectors.
func (c *FakeFabricExplorers) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricExplorerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(fabricexplorersResource, fabricexplorersKind, c.ns, opts), &v1alpha1.FabricExplorerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FabricExplorerList{ListMeta: obj.(*v1alpha1.FabricExplorerList).ListMeta}
	for _, item := range obj.(*v1alpha1.FabricExplorerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested fabricExplorers.
func (c *FakeFabricExplorers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(fabricexplorersResource, c.ns, opts))

}

// Create takes the representation of a fabricExplorer and creates it.  Returns the server's representation of the fabricExplorer, and an error, if there is any.
func (c *FakeFabricExplorers) Create(ctx context.Context, fabricExplorer *v1alpha1.FabricExplorer, opts v1.CreateOptions) (result *v1alpha1.FabricExplorer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(fabricexplorersResource, c.ns, fabricExplorer), &v1alpha1.FabricExplorer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricExplorer), err
}

// Update takes the representation of a fabricExplorer and updates it. Returns the server's representation of the fabricExplorer, and an error, if there is any.
func (c *FakeFabricExplorers) Update(ctx context.Context, fabricExplorer *v1alpha1.FabricExplorer, opts v1.UpdateOptions) (result *v1alpha1.FabricExplorer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(fabricexplorersResource, c.ns, fabricExplorer), &v1alpha1.FabricExplorer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricExplorer), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFabricExplorers) UpdateStatus(ctx context.Context, fabricExplorer *v1alpha1.FabricExplorer, opts v1.UpdateOptions) (*v1alpha1.FabricExplorer, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(fabricexplorersResource, "status", c.ns, fabricExplorer), &v1alpha1.FabricExplorer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricExplorer), err
}

// Delete takes name of the fabricExplorer and deletes it. Returns an error if one occurs.
func (c *FakeFabricExplorers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(fabricexplorersResource, c.ns, name, opts), &v1alpha1.FabricExplorer{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFabricExplorers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(fabricexplorersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.FabricExplorerList{})
	return err
}

// Patch applies the patch and returns the patched fabricExplorer.
func (c *FakeFabricExplorers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricExplorer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricexplorersResource, c.ns, name, pt, data, subresources...), &v1alpha1.FabricExplorer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricExplorer), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricExplorer.
func (c *FakeFabricExplorers) Apply(ctx context.Context, fabricExplorer *hlfkungfusoftwareesv1alpha1.FabricExplorerApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricExplorer, err error) {
	if fabricExplorer == nil {
		return nil, fmt.Errorf("fabricExplorer provided to Apply must not be nil")
	}
	data, err := json.Marshal(fabricExplorer)
	if err != nil {
		return nil, err
	}
	name := fabricExplorer.Name
	if name == nil {
		return nil, fmt.Errorf("fabricExplorer.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricexplorersResource, c.ns, *name, types.ApplyPatchType, data), &v1alpha1.FabricExplorer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricExplorer), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeFabricExplorers) ApplyStatus(ctx context.Context, fabricExplorer *hlfkungfusoftwareesv1alpha1.FabricExplorerApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricExplorer, err error) {
	if fabricExplorer == nil {
		return nil, fmt.Errorf("fabricExplorer provided to Apply must not be nil")
	}
	data, err := json.Marshal(fabricExplorer)
	if err != nil {
		return nil, err
	}
	name := fabricExplorer.Name
	if name == nil {
		return nil, fmt.Errorf("fabricExplorer.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricexplorersResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1alpha1.FabricExplorer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricExplorer), err
}
