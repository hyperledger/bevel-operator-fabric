/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

import (
	v1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// FabricExplorerLister helps list FabricExplorers.
// All objects returned here must be treated as read-only.
type FabricExplorerLister interface {
	// List lists all FabricExplorers in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricExplorer, err error)
	// FabricExplorers returns an object that can list and get FabricExplorers.
	FabricExplorers(namespace string) FabricExplorerNamespaceLister
	FabricExplorerListerExpansion
}

// fabricExplorerLister implements the FabricExplorerLister interface.
type fabricExplorerLister struct {
	indexer cache.Indexer
}

// NewFabricExplorerLister returns a new FabricExplorerLister.
func NewFabricExplorerLister(indexer cache.Indexer) FabricExplorerLister {
	return &fabricExplorerLister{indexer: indexer}
}

// List lists all FabricExplorers in the indexer.
func (s *fabricExplorerLister) List(selector labels.Selector) (ret []*v1alpha1.FabricExplorer, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricExplorer))
	})
	return ret, err
}

// FabricExplorers returns an object that can list and get FabricExplorers.
func (s *fabricExplorerLister) FabricExplorers(namespace string) FabricExplorerNamespaceLister {
	return fabricExplorerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FabricExplorerNamespaceLister helps list and get FabricExplorers.
// All objects returned here must be treated as read-only.
type FabricExplorerNamespaceLister interface {
	// List lists all FabricExplorers in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricExplorer, err error)
	// Get retrieves the FabricExplorer from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.FabricExplorer, error)
	FabricExplorerNamespaceListerExpansion
}

// fabricExplorerNamespaceLister implements the FabricExplorerNamespaceLister
// interface.
type fabricExplorerNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all FabricExplorers in the indexer for a given namespace.
func (s fabricExplorerNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.FabricExplorer, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricExplorer))
	})
	return ret, err
}

// Get retrieves the FabricExplorer from the indexer for a given namespace and name.
func (s fabricExplorerNamespaceLister) Get(name string) (*v1alpha1.FabricExplorer, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("fabricexplorer"), name)
	}
	return obj.(*v1alpha1.FabricExplorer), nil
}
