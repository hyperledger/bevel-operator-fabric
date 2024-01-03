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

// FabricOrderingServicesGetter has a method to return a FabricOrderingServiceInterface.
// A group's client should implement this interface.
type FabricOrderingServicesGetter interface {
	FabricOrderingServices(namespace string) FabricOrderingServiceInterface
}

// FabricOrderingServiceInterface has methods to work with FabricOrderingService resources.
type FabricOrderingServiceInterface interface {
	Create(ctx context.Context, fabricOrderingService *v1alpha1.FabricOrderingService, opts v1.CreateOptions) (*v1alpha1.FabricOrderingService, error)
	Update(ctx context.Context, fabricOrderingService *v1alpha1.FabricOrderingService, opts v1.UpdateOptions) (*v1alpha1.FabricOrderingService, error)
	UpdateStatus(ctx context.Context, fabricOrderingService *v1alpha1.FabricOrderingService, opts v1.UpdateOptions) (*v1alpha1.FabricOrderingService, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.FabricOrderingService, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.FabricOrderingServiceList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricOrderingService, err error)
	Apply(ctx context.Context, fabricOrderingService *hlfkungfusoftwareesv1alpha1.FabricOrderingServiceApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOrderingService, err error)
	ApplyStatus(ctx context.Context, fabricOrderingService *hlfkungfusoftwareesv1alpha1.FabricOrderingServiceApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOrderingService, err error)
	FabricOrderingServiceExpansion
}

// fabricOrderingServices implements FabricOrderingServiceInterface
type fabricOrderingServices struct {
	client rest.Interface
	ns     string
}

// newFabricOrderingServices returns a FabricOrderingServices
func newFabricOrderingServices(c *HlfV1alpha1Client, namespace string) *fabricOrderingServices {
	return &fabricOrderingServices{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the fabricOrderingService, and returns the corresponding fabricOrderingService object, and an error if there is any.
func (c *fabricOrderingServices) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricOrderingService, err error) {
	result = &v1alpha1.FabricOrderingService{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricorderingservices").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FabricOrderingServices that match those selectors.
func (c *fabricOrderingServices) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricOrderingServiceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.FabricOrderingServiceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricorderingservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested fabricOrderingServices.
func (c *fabricOrderingServices) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("fabricorderingservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a fabricOrderingService and creates it.  Returns the server's representation of the fabricOrderingService, and an error, if there is any.
func (c *fabricOrderingServices) Create(ctx context.Context, fabricOrderingService *v1alpha1.FabricOrderingService, opts v1.CreateOptions) (result *v1alpha1.FabricOrderingService, err error) {
	result = &v1alpha1.FabricOrderingService{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("fabricorderingservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOrderingService).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a fabricOrderingService and updates it. Returns the server's representation of the fabricOrderingService, and an error, if there is any.
func (c *fabricOrderingServices) Update(ctx context.Context, fabricOrderingService *v1alpha1.FabricOrderingService, opts v1.UpdateOptions) (result *v1alpha1.FabricOrderingService, err error) {
	result = &v1alpha1.FabricOrderingService{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricorderingservices").
		Name(fabricOrderingService.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOrderingService).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *fabricOrderingServices) UpdateStatus(ctx context.Context, fabricOrderingService *v1alpha1.FabricOrderingService, opts v1.UpdateOptions) (result *v1alpha1.FabricOrderingService, err error) {
	result = &v1alpha1.FabricOrderingService{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricorderingservices").
		Name(fabricOrderingService.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOrderingService).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the fabricOrderingService and deletes it. Returns an error if one occurs.
func (c *fabricOrderingServices) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricorderingservices").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *fabricOrderingServices) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricorderingservices").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched fabricOrderingService.
func (c *fabricOrderingServices) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricOrderingService, err error) {
	result = &v1alpha1.FabricOrderingService{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("fabricorderingservices").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricOrderingService.
func (c *fabricOrderingServices) Apply(ctx context.Context, fabricOrderingService *hlfkungfusoftwareesv1alpha1.FabricOrderingServiceApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOrderingService, err error) {
	if fabricOrderingService == nil {
		return nil, fmt.Errorf("fabricOrderingService provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricOrderingService)
	if err != nil {
		return nil, err
	}
	name := fabricOrderingService.Name
	if name == nil {
		return nil, fmt.Errorf("fabricOrderingService.Name must be provided to Apply")
	}
	result = &v1alpha1.FabricOrderingService{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricorderingservices").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *fabricOrderingServices) ApplyStatus(ctx context.Context, fabricOrderingService *hlfkungfusoftwareesv1alpha1.FabricOrderingServiceApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOrderingService, err error) {
	if fabricOrderingService == nil {
		return nil, fmt.Errorf("fabricOrderingService provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricOrderingService)
	if err != nil {
		return nil, err
	}

	name := fabricOrderingService.Name
	if name == nil {
		return nil, fmt.Errorf("fabricOrderingService.Name must be provided to Apply")
	}

	result = &v1alpha1.FabricOrderingService{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricorderingservices").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
