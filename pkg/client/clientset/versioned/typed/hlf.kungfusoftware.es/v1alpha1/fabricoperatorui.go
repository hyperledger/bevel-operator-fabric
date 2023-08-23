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

// FabricOperatorUIsGetter has a method to return a FabricOperatorUIInterface.
// A group's client should implement this interface.
type FabricOperatorUIsGetter interface {
	FabricOperatorUIs(namespace string) FabricOperatorUIInterface
}

// FabricOperatorUIInterface has methods to work with FabricOperatorUI resources.
type FabricOperatorUIInterface interface {
	Create(ctx context.Context, fabricOperatorUI *v1alpha1.FabricOperatorUI, opts v1.CreateOptions) (*v1alpha1.FabricOperatorUI, error)
	Update(ctx context.Context, fabricOperatorUI *v1alpha1.FabricOperatorUI, opts v1.UpdateOptions) (*v1alpha1.FabricOperatorUI, error)
	UpdateStatus(ctx context.Context, fabricOperatorUI *v1alpha1.FabricOperatorUI, opts v1.UpdateOptions) (*v1alpha1.FabricOperatorUI, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.FabricOperatorUI, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.FabricOperatorUIList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricOperatorUI, err error)
	Apply(ctx context.Context, fabricOperatorUI *hlfkungfusoftwareesv1alpha1.FabricOperatorUIApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperatorUI, err error)
	ApplyStatus(ctx context.Context, fabricOperatorUI *hlfkungfusoftwareesv1alpha1.FabricOperatorUIApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperatorUI, err error)
	FabricOperatorUIExpansion
}

// fabricOperatorUIs implements FabricOperatorUIInterface
type fabricOperatorUIs struct {
	client rest.Interface
	ns     string
}

// newFabricOperatorUIs returns a FabricOperatorUIs
func newFabricOperatorUIs(c *HlfV1alpha1Client, namespace string) *fabricOperatorUIs {
	return &fabricOperatorUIs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the fabricOperatorUI, and returns the corresponding fabricOperatorUI object, and an error if there is any.
func (c *fabricOperatorUIs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricOperatorUI, err error) {
	result = &v1alpha1.FabricOperatorUI{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricoperatoruis").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FabricOperatorUIs that match those selectors.
func (c *fabricOperatorUIs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricOperatorUIList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.FabricOperatorUIList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("fabricoperatoruis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested fabricOperatorUIs.
func (c *fabricOperatorUIs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("fabricoperatoruis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a fabricOperatorUI and creates it.  Returns the server's representation of the fabricOperatorUI, and an error, if there is any.
func (c *fabricOperatorUIs) Create(ctx context.Context, fabricOperatorUI *v1alpha1.FabricOperatorUI, opts v1.CreateOptions) (result *v1alpha1.FabricOperatorUI, err error) {
	result = &v1alpha1.FabricOperatorUI{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("fabricoperatoruis").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOperatorUI).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a fabricOperatorUI and updates it. Returns the server's representation of the fabricOperatorUI, and an error, if there is any.
func (c *fabricOperatorUIs) Update(ctx context.Context, fabricOperatorUI *v1alpha1.FabricOperatorUI, opts v1.UpdateOptions) (result *v1alpha1.FabricOperatorUI, err error) {
	result = &v1alpha1.FabricOperatorUI{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricoperatoruis").
		Name(fabricOperatorUI.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOperatorUI).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *fabricOperatorUIs) UpdateStatus(ctx context.Context, fabricOperatorUI *v1alpha1.FabricOperatorUI, opts v1.UpdateOptions) (result *v1alpha1.FabricOperatorUI, err error) {
	result = &v1alpha1.FabricOperatorUI{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("fabricoperatoruis").
		Name(fabricOperatorUI.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(fabricOperatorUI).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the fabricOperatorUI and deletes it. Returns an error if one occurs.
func (c *fabricOperatorUIs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricoperatoruis").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *fabricOperatorUIs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("fabricoperatoruis").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched fabricOperatorUI.
func (c *fabricOperatorUIs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricOperatorUI, err error) {
	result = &v1alpha1.FabricOperatorUI{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("fabricoperatoruis").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricOperatorUI.
func (c *fabricOperatorUIs) Apply(ctx context.Context, fabricOperatorUI *hlfkungfusoftwareesv1alpha1.FabricOperatorUIApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperatorUI, err error) {
	if fabricOperatorUI == nil {
		return nil, fmt.Errorf("fabricOperatorUI provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricOperatorUI)
	if err != nil {
		return nil, err
	}
	name := fabricOperatorUI.Name
	if name == nil {
		return nil, fmt.Errorf("fabricOperatorUI.Name must be provided to Apply")
	}
	result = &v1alpha1.FabricOperatorUI{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricoperatoruis").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *fabricOperatorUIs) ApplyStatus(ctx context.Context, fabricOperatorUI *hlfkungfusoftwareesv1alpha1.FabricOperatorUIApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOperatorUI, err error) {
	if fabricOperatorUI == nil {
		return nil, fmt.Errorf("fabricOperatorUI provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(fabricOperatorUI)
	if err != nil {
		return nil, err
	}

	name := fabricOperatorUI.Name
	if name == nil {
		return nil, fmt.Errorf("fabricOperatorUI.Name must be provided to Apply")
	}

	result = &v1alpha1.FabricOperatorUI{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("fabricoperatoruis").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
