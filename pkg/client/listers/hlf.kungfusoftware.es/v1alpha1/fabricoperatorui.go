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

// FabricOperatorUILister helps list FabricOperatorUIs.
// All objects returned here must be treated as read-only.
type FabricOperatorUILister interface {
	// List lists all FabricOperatorUIs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricOperatorUI, err error)
	// FabricOperatorUIs returns an object that can list and get FabricOperatorUIs.
	FabricOperatorUIs(namespace string) FabricOperatorUINamespaceLister
	FabricOperatorUIListerExpansion
}

// fabricOperatorUILister implements the FabricOperatorUILister interface.
type fabricOperatorUILister struct {
	indexer cache.Indexer
}

// NewFabricOperatorUILister returns a new FabricOperatorUILister.
func NewFabricOperatorUILister(indexer cache.Indexer) FabricOperatorUILister {
	return &fabricOperatorUILister{indexer: indexer}
}

// List lists all FabricOperatorUIs in the indexer.
func (s *fabricOperatorUILister) List(selector labels.Selector) (ret []*v1alpha1.FabricOperatorUI, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricOperatorUI))
	})
	return ret, err
}

// FabricOperatorUIs returns an object that can list and get FabricOperatorUIs.
func (s *fabricOperatorUILister) FabricOperatorUIs(namespace string) FabricOperatorUINamespaceLister {
	return fabricOperatorUINamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FabricOperatorUINamespaceLister helps list and get FabricOperatorUIs.
// All objects returned here must be treated as read-only.
type FabricOperatorUINamespaceLister interface {
	// List lists all FabricOperatorUIs in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricOperatorUI, err error)
	// Get retrieves the FabricOperatorUI from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.FabricOperatorUI, error)
	FabricOperatorUINamespaceListerExpansion
}

// fabricOperatorUINamespaceLister implements the FabricOperatorUINamespaceLister
// interface.
type fabricOperatorUINamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all FabricOperatorUIs in the indexer for a given namespace.
func (s fabricOperatorUINamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.FabricOperatorUI, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricOperatorUI))
	})
	return ret, err
}

// Get retrieves the FabricOperatorUI from the indexer for a given namespace and name.
func (s fabricOperatorUINamespaceLister) Get(name string) (*v1alpha1.FabricOperatorUI, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("fabricoperatorui"), name)
	}
	return obj.(*v1alpha1.FabricOperatorUI), nil
}
