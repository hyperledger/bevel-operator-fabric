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

// FabricChaincodesGetter has a method to return a FabricChaincodeInterface.
// A group's client should implement this interface.
type FabricChaincodesGetter interface {
	FabricChaincodes(namespace string) FabricChaincodeInterface
}

// FabricChaincodeInterface has methods to work with FabricChaincode resources.
type FabricChaincodeInterface interface {
	Create(ctx context.Context, fabricChaincode *v1alpha1.FabricChaincode, opts v1.CreateOptions) (*v1alpha1.FabricChaincode, error)
	Update(ctx context.Context, fabricChaincode *v1alpha1.FabricChaincode, opts v1.UpdateOptions) (*v1alpha1.FabricChaincode, error)
	UpdateStatus(ctx context.Context, fabricChaincode *v1alpha1.FabricChaincode, opts v1.UpdateOptions) (*v1alpha1.FabricChaincode, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.FabricChaincode, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.FabricChaincodeList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricChaincode, err error)
	Apply(ctx context.Context, fabricChaincode *hlfkungfusoftwareesv1alpha1.FabricChaincodeApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricChaincode, err error)
	ApplyStatus(ctx context.Context, fabricChaincode *hlfkungfusoftwareesv1alpha1.FabricChaincodeApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricChaincode, err error)
	FabricChaincodeExpansion
}

// fabricChaincodes implements FabricChaincodeInterface
type fabricChaincodes struct {
	client rest.Interface
	ns     string
}

// newFabricChaincodes returns a FabricChaincodes
func newFabricChaincodes(c *HlfV1alpha1Client, namespace string) *fabricChaincodes {
	return &fabricChaincodes{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the fabricChaincode, and returns the corresponding fabricChaincode object, and an error if there is any.
func (c *fabricChaincodes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricChaincode, err error) {
	result = &v1alpha1.FabricChaincode{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricchaincodes").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FabricChaincodes that match those selectors.
func (c *fabricChaincodes) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricChaincodeList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.FabricChaincodeList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricchaincodes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested fabricChaincodes.
func (c *fabricChaincodes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("fabricchaincodes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a fabricChaincode and creates it.  Returns the server's representation of the fabricChaincode, and an error, if there is any.
func (c *fabricChaincodes) Create(ctx context.Context, fabricChaincode *v1alpha1.FabricChaincode, opts v1.CreateOptions) (result *v1alpha1.FabricChaincode, err error) {
	result = &v1alpha1.FabricChaincode{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("fabricchaincodes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricChaincode).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a fabricChaincode and updates it. Returns the server's representation of the fabricChaincode, and an error, if there is any.
func (c *fabricChaincodes) Update(ctx context.Context, fabricChaincode *v1alpha1.FabricChaincode, opts v1.UpdateOptions) (result *v1alpha1.FabricChaincode, err error) {
	result = &v1alpha1.FabricChaincode{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricchaincodes").
		Name(fabricChaincode.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricChaincode).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *fabricChaincodes) UpdateStatus(ctx context.Context, fabricChaincode *v1alpha1.FabricChaincode, opts v1.UpdateOptions) (result *v1alpha1.FabricChaincode, err error) {
	result = &v1alpha1.FabricChaincode{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricchaincodes").
		Name(fabricChaincode.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricChaincode).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the fabricChaincode and deletes it. Returns an error if one occurs.
func (c *fabricChaincodes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricchaincodes").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *fabricChaincodes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricchaincodes").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched fabricChaincode.
func (c *fabricChaincodes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricChaincode, err error) {
	result = &v1alpha1.FabricChaincode{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("fabricchaincodes").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricChaincode.
func (c *fabricChaincodes) Apply(ctx context.Context, fabricChaincode *hlfkungfusoftwareesv1alpha1.FabricChaincodeApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricChaincode, err error) {
	if fabricChaincode == nil {
		return nil, fmt.Errorf("fabricChaincode provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricChaincode)
	if err != nil {
		return nil, err
	}
	name := fabricChaincode.Name
	if name == nil {
		return nil, fmt.Errorf("fabricChaincode.Name must be provided to Apply")
	}
	result = &v1alpha1.FabricChaincode{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricchaincodes").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *fabricChaincodes) ApplyStatus(ctx context.Context, fabricChaincode *hlfkungfusoftwareesv1alpha1.FabricChaincodeApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricChaincode, err error) {
	if fabricChaincode == nil {
		return nil, fmt.Errorf("fabricChaincode provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricChaincode)
	if err != nil {
		return nil, err
	}

	name := fabricChaincode.Name
	if name == nil {
		return nil, fmt.Errorf("fabricChaincode.Name must be provided to Apply")
	}

	result = &v1alpha1.FabricChaincode{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricchaincodes").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
