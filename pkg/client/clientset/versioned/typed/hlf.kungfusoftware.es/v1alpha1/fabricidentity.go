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

// FabricIdentitiesGetter has a method to return a FabricIdentityInterface.
// A group's client should implement this interface.
type FabricIdentitiesGetter interface {
	FabricIdentities(namespace string) FabricIdentityInterface
}

// FabricIdentityInterface has methods to work with FabricIdentity resources.
type FabricIdentityInterface interface {
	Create(ctx context.Context, fabricIdentity *v1alpha1.FabricIdentity, opts v1.CreateOptions) (*v1alpha1.FabricIdentity, error)
	Update(ctx context.Context, fabricIdentity *v1alpha1.FabricIdentity, opts v1.UpdateOptions) (*v1alpha1.FabricIdentity, error)
	UpdateStatus(ctx context.Context, fabricIdentity *v1alpha1.FabricIdentity, opts v1.UpdateOptions) (*v1alpha1.FabricIdentity, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.FabricIdentity, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.FabricIdentityList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricIdentity, err error)
	Apply(ctx context.Context, fabricIdentity *hlfkungfusoftwareesv1alpha1.FabricIdentityApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricIdentity, err error)
	ApplyStatus(ctx context.Context, fabricIdentity *hlfkungfusoftwareesv1alpha1.FabricIdentityApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricIdentity, err error)
	FabricIdentityExpansion
}

// fabricIdentities implements FabricIdentityInterface
type fabricIdentities struct {
	client rest.Interface
	ns     string
}

// newFabricIdentities returns a FabricIdentities
func newFabricIdentities(c *HlfV1alpha1Client, namespace string) *fabricIdentities {
	return &fabricIdentities{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the fabricIdentity, and returns the corresponding fabricIdentity object, and an error if there is any.
func (c *fabricIdentities) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricIdentity, err error) {
	result = &v1alpha1.FabricIdentity{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricidentities").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FabricIdentities that match those selectors.
func (c *fabricIdentities) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricIdentityList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.FabricIdentityList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricidentities").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested fabricIdentities.
func (c *fabricIdentities) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("fabricidentities").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a fabricIdentity and creates it.  Returns the server's representation of the fabricIdentity, and an error, if there is any.
func (c *fabricIdentities) Create(ctx context.Context, fabricIdentity *v1alpha1.FabricIdentity, opts v1.CreateOptions) (result *v1alpha1.FabricIdentity, err error) {
	result = &v1alpha1.FabricIdentity{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("fabricidentities").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricIdentity).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a fabricIdentity and updates it. Returns the server's representation of the fabricIdentity, and an error, if there is any.
func (c *fabricIdentities) Update(ctx context.Context, fabricIdentity *v1alpha1.FabricIdentity, opts v1.UpdateOptions) (result *v1alpha1.FabricIdentity, err error) {
	result = &v1alpha1.FabricIdentity{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricidentities").
		Name(fabricIdentity.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricIdentity).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *fabricIdentities) UpdateStatus(ctx context.Context, fabricIdentity *v1alpha1.FabricIdentity, opts v1.UpdateOptions) (result *v1alpha1.FabricIdentity, err error) {
	result = &v1alpha1.FabricIdentity{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricidentities").
		Name(fabricIdentity.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricIdentity).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the fabricIdentity and deletes it. Returns an error if one occurs.
func (c *fabricIdentities) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricidentities").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *fabricIdentities) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricidentities").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched fabricIdentity.
func (c *fabricIdentities) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricIdentity, err error) {
	result = &v1alpha1.FabricIdentity{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("fabricidentities").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricIdentity.
func (c *fabricIdentities) Apply(ctx context.Context, fabricIdentity *hlfkungfusoftwareesv1alpha1.FabricIdentityApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricIdentity, err error) {
	if fabricIdentity == nil {
		return nil, fmt.Errorf("fabricIdentity provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricIdentity)
	if err != nil {
		return nil, err
	}
	name := fabricIdentity.Name
	if name == nil {
		return nil, fmt.Errorf("fabricIdentity.Name must be provided to Apply")
	}
	result = &v1alpha1.FabricIdentity{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricidentities").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *fabricIdentities) ApplyStatus(ctx context.Context, fabricIdentity *hlfkungfusoftwareesv1alpha1.FabricIdentityApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricIdentity, err error) {
	if fabricIdentity == nil {
		return nil, fmt.Errorf("fabricIdentity provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricIdentity)
	if err != nil {
		return nil, err
	}

	name := fabricIdentity.Name
	if name == nil {
		return nil, fmt.Errorf("fabricIdentity.Name must be provided to Apply")
	}

	result = &v1alpha1.FabricIdentity{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricidentities").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
