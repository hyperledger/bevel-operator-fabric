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

// FabricOperatorAPILister helps list FabricOperatorAPIs.
// All objects returned here must be treated as read-only.
type FabricOperatorAPILister interface {
	// List lists all FabricOperatorAPIs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricOperatorAPI, err error)
	// FabricOperatorAPIs returns an object that can list and get FabricOperatorAPIs.
	FabricOperatorAPIs(namespace string) FabricOperatorAPINamespaceLister
	FabricOperatorAPIListerExpansion
}

// fabricOperatorAPILister implements the FabricOperatorAPILister interface.
type fabricOperatorAPILister struct {
	indexer cache.Indexer
}

// NewFabricOperatorAPILister returns a new FabricOperatorAPILister.
func NewFabricOperatorAPILister(indexer cache.Indexer) FabricOperatorAPILister {
	return &fabricOperatorAPILister{indexer: indexer}
}

// List lists all FabricOperatorAPIs in the indexer.
func (s *fabricOperatorAPILister) List(selector labels.Selector) (ret []*v1alpha1.FabricOperatorAPI, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricOperatorAPI))
	})
	return ret, err
}

// FabricOperatorAPIs returns an object that can list and get FabricOperatorAPIs.
func (s *fabricOperatorAPILister) FabricOperatorAPIs(namespace string) FabricOperatorAPINamespaceLister {
	return fabricOperatorAPINamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FabricOperatorAPINamespaceLister helps list and get FabricOperatorAPIs.
// All objects returned here must be treated as read-only.
type FabricOperatorAPINamespaceLister interface {
	// List lists all FabricOperatorAPIs in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricOperatorAPI, err error)
	// Get retrieves the FabricOperatorAPI from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.FabricOperatorAPI, error)
	FabricOperatorAPINamespaceListerExpansion
}

// fabricOperatorAPINamespaceLister implements the FabricOperatorAPINamespaceLister
// interface.
type fabricOperatorAPINamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all FabricOperatorAPIs in the indexer for a given namespace.
func (s fabricOperatorAPINamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.FabricOperatorAPI, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricOperatorAPI))
	})
	return ret, err
}

// Get retrieves the FabricOperatorAPI from the indexer for a given namespace and name.
func (s fabricOperatorAPINamespaceLister) Get(name string) (*v1alpha1.FabricOperatorAPI, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("fabricoperatorapi"), name)
	}
	return obj.(*v1alpha1.FabricOperatorAPI), nil
}
