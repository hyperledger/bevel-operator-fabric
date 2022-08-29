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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFabricPeers implements FabricPeerInterface
type FakeFabricPeers struct {
	Fake *FakeHlfV1alpha1
	ns   string
}

var fabricpeersResource = schema.GroupVersionResource{Group: "hlf.kungfusoftware.es", Version: "v1alpha1", Resource: "fabricpeers"}

var fabricpeersKind = schema.GroupVersionKind{Group: "hlf.kungfusoftware.es", Version: "v1alpha1", Kind: "FabricPeer"}

// Get takes name of the fabricPeer, and returns the corresponding fabricPeer object, and an error if there is any.
func (c *FakeFabricPeers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricPeer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(fabricpeersResource, c.ns, name), &v1alpha1.FabricPeer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricPeer), err
}

// List takes label and field selectors, and returns the list of FabricPeers that match those selectors.
func (c *FakeFabricPeers) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricPeerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(fabricpeersResource, fabricpeersKind, c.ns, opts), &v1alpha1.FabricPeerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FabricPeerList{ListMeta: obj.(*v1alpha1.FabricPeerList).ListMeta}
	for _, item := range obj.(*v1alpha1.FabricPeerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested fabricPeers.
func (c *FakeFabricPeers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(fabricpeersResource, c.ns, opts))

}

// Create takes the representation of a fabricPeer and creates it.  Returns the server's representation of the fabricPeer, and an error, if there is any.
func (c *FakeFabricPeers) Create(ctx context.Context, fabricPeer *v1alpha1.FabricPeer, opts v1.CreateOptions) (result *v1alpha1.FabricPeer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(fabricpeersResource, c.ns, fabricPeer), &v1alpha1.FabricPeer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricPeer), err
}

// Update takes the representation of a fabricPeer and updates it. Returns the server's representation of the fabricPeer, and an error, if there is any.
func (c *FakeFabricPeers) Update(ctx context.Context, fabricPeer *v1alpha1.FabricPeer, opts v1.UpdateOptions) (result *v1alpha1.FabricPeer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(fabricpeersResource, c.ns, fabricPeer), &v1alpha1.FabricPeer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricPeer), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFabricPeers) UpdateStatus(ctx context.Context, fabricPeer *v1alpha1.FabricPeer, opts v1.UpdateOptions) (*v1alpha1.FabricPeer, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(fabricpeersResource, "status", c.ns, fabricPeer), &v1alpha1.FabricPeer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricPeer), err
}

// Delete takes name of the fabricPeer and deletes it. Returns an error if one occurs.
func (c *FakeFabricPeers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(fabricpeersResource, c.ns, name, opts), &v1alpha1.FabricPeer{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFabricPeers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(fabricpeersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.FabricPeerList{})
	return err
}

// Patch applies the patch and returns the patched fabricPeer.
func (c *FakeFabricPeers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricPeer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricpeersResource, c.ns, name, pt, data, subresources...), &v1alpha1.FabricPeer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricPeer), err
}
