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

// FabricIdentityLister helps list FabricIdentities.
// All objects returned here must be treated as read-only.
type FabricIdentityLister interface {
	// List lists all FabricIdentities in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricIdentity, err error)
	// FabricIdentities returns an object that can list and get FabricIdentities.
	FabricIdentities(namespace string) FabricIdentityNamespaceLister
	FabricIdentityListerExpansion
}

// fabricIdentityLister implements the FabricIdentityLister interface.
type fabricIdentityLister struct {
	indexer cache.Indexer
}

// NewFabricIdentityLister returns a new FabricIdentityLister.
func NewFabricIdentityLister(indexer cache.Indexer) FabricIdentityLister {
	return &fabricIdentityLister{indexer: indexer}
}

// List lists all FabricIdentities in the indexer.
func (s *fabricIdentityLister) List(selector labels.Selector) (ret []*v1alpha1.FabricIdentity, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricIdentity))
	})
	return ret, err
}

// FabricIdentities returns an object that can list and get FabricIdentities.
func (s *fabricIdentityLister) FabricIdentities(namespace string) FabricIdentityNamespaceLister {
	return fabricIdentityNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FabricIdentityNamespaceLister helps list and get FabricIdentities.
// All objects returned here must be treated as read-only.
type FabricIdentityNamespaceLister interface {
	// List lists all FabricIdentities in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricIdentity, err error)
	// Get retrieves the FabricIdentity from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.FabricIdentity, error)
	FabricIdentityNamespaceListerExpansion
}

// fabricIdentityNamespaceLister implements the FabricIdentityNamespaceLister
// interface.
type fabricIdentityNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all FabricIdentities in the indexer for a given namespace.
func (s fabricIdentityNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.FabricIdentity, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricIdentity))
	})
	return ret, err
}

// Get retrieves the FabricIdentity from the indexer for a given namespace and name.
func (s fabricIdentityNamespaceLister) Get(name string) (*v1alpha1.FabricIdentity, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("fabricidentity"), name)
	}
	return obj.(*v1alpha1.FabricIdentity), nil
}
