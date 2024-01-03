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

// FabricChaincodeLister helps list FabricChaincodes.
// All objects returned here must be treated as read-only.
type FabricChaincodeLister interface {
	// List lists all FabricChaincodes in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricChaincode, err error)
	// FabricChaincodes returns an object that can list and get FabricChaincodes.
	FabricChaincodes(namespace string) FabricChaincodeNamespaceLister
	FabricChaincodeListerExpansion
}

// fabricChaincodeLister implements the FabricChaincodeLister interface.
type fabricChaincodeLister struct {
	indexer cache.Indexer
}

// NewFabricChaincodeLister returns a new FabricChaincodeLister.
func NewFabricChaincodeLister(indexer cache.Indexer) FabricChaincodeLister {
	return &fabricChaincodeLister{indexer: indexer}
}

// List lists all FabricChaincodes in the indexer.
func (s *fabricChaincodeLister) List(selector labels.Selector) (ret []*v1alpha1.FabricChaincode, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricChaincode))
	})
	return ret, err
}

// FabricChaincodes returns an object that can list and get FabricChaincodes.
func (s *fabricChaincodeLister) FabricChaincodes(namespace string) FabricChaincodeNamespaceLister {
	return fabricChaincodeNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FabricChaincodeNamespaceLister helps list and get FabricChaincodes.
// All objects returned here must be treated as read-only.
type FabricChaincodeNamespaceLister interface {
	// List lists all FabricChaincodes in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricChaincode, err error)
	// Get retrieves the FabricChaincode from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.FabricChaincode, error)
	FabricChaincodeNamespaceListerExpansion
}

// fabricChaincodeNamespaceLister implements the FabricChaincodeNamespaceLister
// interface.
type fabricChaincodeNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all FabricChaincodes in the indexer for a given namespace.
func (s fabricChaincodeNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.FabricChaincode, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricChaincode))
	})
	return ret, err
}

// Get retrieves the FabricChaincode from the indexer for a given namespace and name.
func (s fabricChaincodeNamespaceLister) Get(name string) (*v1alpha1.FabricChaincode, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("fabricchaincode"), name)
	}
	return obj.(*v1alpha1.FabricChaincode), nil
}
