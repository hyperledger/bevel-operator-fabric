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

// FabricOperatorAPIsGetter has a method to return a FabricOperatorAPIInterface.
// A group's client should implement this interface.
type FabricOperatorAPIsGetter interface {
	FabricOperatorAPIs(namespace string) FabricOperatorAPIInterface
}

// FabricOperatorAPIInterface has methods to work with FabricOperatorAPI resources.
type FabricOperatorAPIInterface interface {
	Create(ctx context.Context, fabricOperatorAPI *v1alpha1.FabricOperatorAPI, opts v1.CreateOptions) (*v1alpha1.FabricOperatorAPI, error)
	Update(ctx context.Context, fabricOperatorAPI *v1alpha1.FabricOperatorAPI, opts v1.UpdateOptions) (*v1alpha1.FabricOperatorAPI, error)
	UpdateStatus(ctx context.Context, fabricOperatorAPI *v1alpha1.FabricOperatorAPI, opts v1.UpdateOptions) (*v1alpha1.FabricOperatorAPI, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.FabricOperatorAPI, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.FabricOperatorAPIList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricOperatorAPI, err error)
	Apply(ctx context.Context, fabricOperatorAPI *hlfkungfusoftwareesv1alpha1.FabricOperatorAPIApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperatorAPI, err error)
	ApplyStatus(ctx context.Context, fabricOperatorAPI *hlfkungfusoftwareesv1alpha1.FabricOperatorAPIApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperatorAPI, err error)
	FabricOperatorAPIExpansion
}

// fabricOperatorAPIs implements FabricOperatorAPIInterface
type fabricOperatorAPIs struct {
	client rest.Interface
	ns     string
}

// newFabricOperatorAPIs returns a FabricOperatorAPIs
func newFabricOperatorAPIs(c *HlfV1alpha1Client, namespace string) *fabricOperatorAPIs {
	return &fabricOperatorAPIs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the fabricOperatorAPI, and returns the corresponding fabricOperatorAPI object, and an error if there is any.
func (c *fabricOperatorAPIs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricOperatorAPI, err error) {
	result = &v1alpha1.FabricOperatorAPI{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricoperatorapis").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FabricOperatorAPIs that match those selectors.
func (c *fabricOperatorAPIs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricOperatorAPIList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.FabricOperatorAPIList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricoperatorapis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested fabricOperatorAPIs.
func (c *fabricOperatorAPIs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("fabricoperatorapis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a fabricOperatorAPI and creates it.  Returns the server's representation of the fabricOperatorAPI, and an error, if there is any.
func (c *fabricOperatorAPIs) Create(ctx context.Context, fabricOperatorAPI *v1alpha1.FabricOperatorAPI, opts v1.CreateOptions) (result *v1alpha1.FabricOperatorAPI, err error) {
	result = &v1alpha1.FabricOperatorAPI{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("fabricoperatorapis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOperatorAPI).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a fabricOperatorAPI and updates it. Returns the server's representation of the fabricOperatorAPI, and an error, if there is any.
func (c *fabricOperatorAPIs) Update(ctx context.Context, fabricOperatorAPI *v1alpha1.FabricOperatorAPI, opts v1.UpdateOptions) (result *v1alpha1.FabricOperatorAPI, err error) {
	result = &v1alpha1.FabricOperatorAPI{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricoperatorapis").
		Name(fabricOperatorAPI.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOperatorAPI).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *fabricOperatorAPIs) UpdateStatus(ctx context.Context, fabricOperatorAPI *v1alpha1.FabricOperatorAPI, opts v1.UpdateOptions) (result *v1alpha1.FabricOperatorAPI, err error) {
	result = &v1alpha1.FabricOperatorAPI{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricoperatorapis").
		Name(fabricOperatorAPI.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOperatorAPI).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the fabricOperatorAPI and deletes it. Returns an error if one occurs.
func (c *fabricOperatorAPIs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricoperatorapis").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *fabricOperatorAPIs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricoperatorapis").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched fabricOperatorAPI.
func (c *fabricOperatorAPIs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricOperatorAPI, err error) {
	result = &v1alpha1.FabricOperatorAPI{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("fabricoperatorapis").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricOperatorAPI.
func (c *fabricOperatorAPIs) Apply(ctx context.Context, fabricOperatorAPI *hlfkungfusoftwareesv1alpha1.FabricOperatorAPIApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperatorAPI, err error) {
	if fabricOperatorAPI == nil {
		return nil, fmt.Errorf("fabricOperatorAPI provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricOperatorAPI)
	if err != nil {
		return nil, err
	}
	name := fabricOperatorAPI.Name
	if name == nil {
		return nil, fmt.Errorf("fabricOperatorAPI.Name must be provided to Apply")
	}
	result = &v1alpha1.FabricOperatorAPI{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricoperatorapis").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *fabricOperatorAPIs) ApplyStatus(ctx context.Context, fabricOperatorAPI *hlfkungfusoftwareesv1alpha1.FabricOperatorAPIApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperatorAPI, err error) {
	if fabricOperatorAPI == nil {
		return nil, fmt.Errorf("fabricOperatorAPI provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricOperatorAPI)
	if err != nil {
		return nil, err
	}

	name := fabricOperatorAPI.Name
	if name == nil {
		return nil, fmt.Errorf("fabricOperatorAPI.Name must be provided to Apply")
	}

	result = &v1alpha1.FabricOperatorAPI{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricoperatorapis").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
