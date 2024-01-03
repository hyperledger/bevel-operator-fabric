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

// FabricOperationsConsolesGetter has a method to return a FabricOperationsConsoleInterface.
// A group's client should implement this interface.
type FabricOperationsConsolesGetter interface {
	FabricOperationsConsoles(namespace string) FabricOperationsConsoleInterface
}

// FabricOperationsConsoleInterface has methods to work with FabricOperationsConsole resources.
type FabricOperationsConsoleInterface interface {
	Create(ctx context.Context, fabricOperationsConsole *v1alpha1.FabricOperationsConsole, opts v1.CreateOptions) (*v1alpha1.FabricOperationsConsole, error)
	Update(ctx context.Context, fabricOperationsConsole *v1alpha1.FabricOperationsConsole, opts v1.UpdateOptions) (*v1alpha1.FabricOperationsConsole, error)
	UpdateStatus(ctx context.Context, fabricOperationsConsole *v1alpha1.FabricOperationsConsole, opts v1.UpdateOptions) (*v1alpha1.FabricOperationsConsole, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.FabricOperationsConsole, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.FabricOperationsConsoleList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricOperationsConsole, err error)
	Apply(ctx context.Context, fabricOperationsConsole *hlfkungfusoftwareesv1alpha1.FabricOperationsConsoleApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperationsConsole, err error)
	ApplyStatus(ctx context.Context, fabricOperationsConsole *hlfkungfusoftwareesv1alpha1.FabricOperationsConsoleApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperationsConsole, err error)
	FabricOperationsConsoleExpansion
}

// fabricOperationsConsoles implements FabricOperationsConsoleInterface
type fabricOperationsConsoles struct {
	client rest.Interface
	ns     string
}

// newFabricOperationsConsoles returns a FabricOperationsConsoles
func newFabricOperationsConsoles(c *HlfV1alpha1Client, namespace string) *fabricOperationsConsoles {
	return &fabricOperationsConsoles{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the fabricOperationsConsole, and returns the corresponding fabricOperationsConsole object, and an error if there is any.
func (c *fabricOperationsConsoles) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricOperationsConsole, err error) {
	result = &v1alpha1.FabricOperationsConsole{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricoperationsconsoles").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FabricOperationsConsoles that match those selectors.
func (c *fabricOperationsConsoles) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricOperationsConsoleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.FabricOperationsConsoleList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricoperationsconsoles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested fabricOperationsConsoles.
func (c *fabricOperationsConsoles) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("fabricoperationsconsoles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a fabricOperationsConsole and creates it.  Returns the server's representation of the fabricOperationsConsole, and an error, if there is any.
func (c *fabricOperationsConsoles) Create(ctx context.Context, fabricOperationsConsole *v1alpha1.FabricOperationsConsole, opts v1.CreateOptions) (result *v1alpha1.FabricOperationsConsole, err error) {
	result = &v1alpha1.FabricOperationsConsole{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("fabricoperationsconsoles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOperationsConsole).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a fabricOperationsConsole and updates it. Returns the server's representation of the fabricOperationsConsole, and an error, if there is any.
func (c *fabricOperationsConsoles) Update(ctx context.Context, fabricOperationsConsole *v1alpha1.FabricOperationsConsole, opts v1.UpdateOptions) (result *v1alpha1.FabricOperationsConsole, err error) {
	result = &v1alpha1.FabricOperationsConsole{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricoperationsconsoles").
		Name(fabricOperationsConsole.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOperationsConsole).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *fabricOperationsConsoles) UpdateStatus(ctx context.Context, fabricOperationsConsole *v1alpha1.FabricOperationsConsole, opts v1.UpdateOptions) (result *v1alpha1.FabricOperationsConsole, err error) {
	result = &v1alpha1.FabricOperationsConsole{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricoperationsconsoles").
		Name(fabricOperationsConsole.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOperationsConsole).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the fabricOperationsConsole and deletes it. Returns an error if one occurs.
func (c *fabricOperationsConsoles) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricoperationsconsoles").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *fabricOperationsConsoles) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricoperationsconsoles").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched fabricOperationsConsole.
func (c *fabricOperationsConsoles) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricOperationsConsole, err error) {
	result = &v1alpha1.FabricOperationsConsole{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("fabricoperationsconsoles").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricOperationsConsole.
func (c *fabricOperationsConsoles) Apply(ctx context.Context, fabricOperationsConsole *hlfkungfusoftwareesv1alpha1.FabricOperationsConsoleApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperationsConsole, err error) {
	if fabricOperationsConsole == nil {
		return nil, fmt.Errorf("fabricOperationsConsole provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricOperationsConsole)
	if err != nil {
		return nil, err
	}
	name := fabricOperationsConsole.Name
	if name == nil {
		return nil, fmt.Errorf("fabricOperationsConsole.Name must be provided to Apply")
	}
	result = &v1alpha1.FabricOperationsConsole{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricoperationsconsoles").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *fabricOperationsConsoles) ApplyStatus(ctx context.Context, fabricOperationsConsole *hlfkungfusoftwareesv1alpha1.FabricOperationsConsoleApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperationsConsole, err error) {
	if fabricOperationsConsole == nil {
		return nil, fmt.Errorf("fabricOperationsConsole provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricOperationsConsole)
	if err != nil {
		return nil, err
	}

	name := fabricOperationsConsole.Name
	if name == nil {
		return nil, fmt.Errorf("fabricOperationsConsole.Name must be provided to Apply")
	}

	result = &v1alpha1.FabricOperationsConsole{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricoperationsconsoles").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
