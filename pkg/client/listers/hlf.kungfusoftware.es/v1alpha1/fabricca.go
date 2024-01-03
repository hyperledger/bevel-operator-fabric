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

// FabricCALister helps list FabricCAs.
// All objects returned here must be treated as read-only.
type FabricCALister interface {
	// List lists all FabricCAs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricCA, err error)
	// FabricCAs returns an object that can list and get FabricCAs.
	FabricCAs(namespace string) FabricCANamespaceLister
	FabricCAListerExpansion
}

// fabricCALister implements the FabricCALister interface.
type fabricCALister struct {
	indexer cache.Indexer
}

// NewFabricCALister returns a new FabricCALister.
func NewFabricCALister(indexer cache.Indexer) FabricCALister {
	return &fabricCALister{indexer: indexer}
}

// List lists all FabricCAs in the indexer.
func (s *fabricCALister) List(selector labels.Selector) (ret []*v1alpha1.FabricCA, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricCA))
	})
	return ret, err
}

// FabricCAs returns an object that can list and get FabricCAs.
func (s *fabricCALister) FabricCAs(namespace string) FabricCANamespaceLister {
	return fabricCANamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FabricCANamespaceLister helps list and get FabricCAs.
// All objects returned here must be treated as read-only.
type FabricCANamespaceLister interface {
	// List lists all FabricCAs in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricCA, err error)
	// Get retrieves the FabricCA from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.FabricCA, error)
	FabricCANamespaceListerExpansion
}

// fabricCANamespaceLister implements the FabricCANamespaceLister
// interface.
type fabricCANamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all FabricCAs in the indexer for a given namespace.
func (s fabricCANamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.FabricCA, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricCA))
	})
	return ret, err
}

// Get retrieves the FabricCA from the indexer for a given namespace and name.
func (s fabricCANamespaceLister) Get(name string) (*v1alpha1.FabricCA, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("fabricca"), name)
	}
	return obj.(*v1alpha1.FabricCA), nil
}
