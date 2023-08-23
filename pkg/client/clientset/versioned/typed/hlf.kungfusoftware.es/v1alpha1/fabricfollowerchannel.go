/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	hlfkungfusoftwareesv1alpha1 "github.com/kfsoftware/hlf-operator/pkg/client/applyconfiguration/hlf.kungfusoftware.es/v1alpha1"
	scheme "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// FabricFollowerChannelsGetter has a method to return a FabricFollowerChannelInterface.
// A group's client should implement this interface.
type FabricFollowerChannelsGetter interface {
	FabricFollowerChannels() FabricFollowerChannelInterface
}

// FabricFollowerChannelInterface has methods to work with FabricFollowerChannel resources.
type FabricFollowerChannelInterface interface {
	Create(ctx context.Context, fabricFollowerChannel *v1alpha1.FabricFollowerChannel, opts v1.CreateOptions) (*v1alpha1.FabricFollowerChannel, error)
	Update(ctx context.Context, fabricFollowerChannel *v1alpha1.FabricFollowerChannel, opts v1.UpdateOptions) (*v1alpha1.FabricFollowerChannel, error)
	UpdateStatus(ctx context.Context, fabricFollowerChannel *v1alpha1.FabricFollowerChannel, opts v1.UpdateOptions) (*v1alpha1.FabricFollowerChannel, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.FabricFollowerChannel, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.FabricFollowerChannelList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricFollowerChannel, err error)
	Apply(ctx context.Context, fabricFollowerChannel *hlfkungfusoftwareesv1alpha1.FabricFollowerChannelApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricFollowerChannel, err error)
	ApplyStatus(ctx context.Context, fabricFollowerChannel *hlfkungfusoftwareesv1alpha1.FabricFollowerChannelApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricFollowerChannel, err error)
	FabricFollowerChannelExpansion
}

// fabricFollowerChannels implements FabricFollowerChannelInterface
type fabricFollowerChannels struct {
	client rest.Interface
}

// newFabricFollowerChannels returns a FabricFollowerChannels
func newFabricFollowerChannels(c *HlfV1alpha1Client) *fabricFollowerChannels {
	return &fabricFollowerChannels{
		client: c.RESTClient(),
	}
}

// Get takes name of the fabricFollowerChannel, and returns the corresponding fabricFollowerChannel object, and an error if there is any.
func (c *fabricFollowerChannels) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricFollowerChannel, err error) {
	result = &v1alpha1.FabricFollowerChannel{}
	err = c.client.Get().
		Resource("fabricfollowerchannels").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FabricFollowerChannels that match those selectors.
func (c *fabricFollowerChannels) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricFollowerChannelList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.FabricFollowerChannelList{}
	err = c.client.Get().
		Resource("fabricfollowerchannels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested fabricFollowerChannels.
func (c *fabricFollowerChannels) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("fabricfollowerchannels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a fabricFollowerChannel and creates it.  Returns the server's representation of the fabricFollowerChannel, and an error, if there is any.
func (c *fabricFollowerChannels) Create(ctx context.Context, fabricFollowerChannel *v1alpha1.FabricFollowerChannel, opts v1.CreateOptions) (result *v1alpha1.FabricFollowerChannel, err error) {
	result = &v1alpha1.FabricFollowerChannel{}
	err = c.client.Post().
		Resource("fabricfollowerchannels").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricFollowerChannel).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a fabricFollowerChannel and updates it. Returns the server's representation of the fabricFollowerChannel, and an error, if there is any.
func (c *fabricFollowerChannels) Update(ctx context.Context, fabricFollowerChannel *v1alpha1.FabricFollowerChannel, opts v1.UpdateOptions) (result *v1alpha1.FabricFollowerChannel, err error) {
	result = &v1alpha1.FabricFollowerChannel{}
	err = c.client.Put().
		Resource("fabricfollowerchannels").
		Name(fabricFollowerChannel.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricFollowerChannel).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *fabricFollowerChannels) UpdateStatus(ctx context.Context, fabricFollowerChannel *v1alpha1.FabricFollowerChannel, opts v1.UpdateOptions) (result *v1alpha1.FabricFollowerChannel, err error) {
	result = &v1alpha1.FabricFollowerChannel{}
	err = c.client.Put().
		Resource("fabricfollowerchannels").
		Name(fabricFollowerChannel.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricFollowerChannel).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the fabricFollowerChannel and deletes it. Returns an error if one occurs.
func (c *fabricFollowerChannels) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("fabricfollowerchannels").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *fabricFollowerChannels) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("fabricfollowerchannels").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched fabricFollowerChannel.
func (c *fabricFollowerChannels) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricFollowerChannel, err error) {
	result = &v1alpha1.FabricFollowerChannel{}
	err = c.client.Patch(pt).
		Resource("fabricfollowerchannels").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricFollowerChannel.
func (c *fabricFollowerChannels) Apply(ctx context.Context, fabricFollowerChannel *hlfkungfusoftwareesv1alpha1.FabricFollowerChannelApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricFollowerChannel, err error) {
	if fabricFollowerChannel == nil {
		return nil, fmt.Errorf("fabricFollowerChannel provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricFollowerChannel)
	if err != nil {
		return nil, err
	}
	name := fabricFollowerChannel.Name
	if name == nil {
		return nil, fmt.Errorf("fabricFollowerChannel.Name must be provided to Apply")
	}
	result = &v1alpha1.FabricFollowerChannel{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("fabricfollowerchannels").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *fabricFollowerChannels) ApplyStatus(ctx context.Context, fabricFollowerChannel *hlfkungfusoftwareesv1alpha1.FabricFollowerChannelApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricFollowerChannel, err error) {
	if fabricFollowerChannel == nil {
		return nil, fmt.Errorf("fabricFollowerChannel provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricFollowerChannel)
	if err != nil {
		return nil, err
	}

	name := fabricFollowerChannel.Name
	if name == nil {
		return nil, fmt.Errorf("fabricFollowerChannel.Name must be provided to Apply")
	}

	result = &v1alpha1.FabricFollowerChannel{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("fabricfollowerchannels").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
