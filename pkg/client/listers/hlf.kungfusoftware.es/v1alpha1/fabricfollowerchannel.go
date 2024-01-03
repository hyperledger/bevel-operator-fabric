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

// FabricFollowerChannelLister helps list FabricFollowerChannels.
// All objects returned here must be treated as read-only.
type FabricFollowerChannelLister interface {
	// List lists all FabricFollowerChannels in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.FabricFollowerChannel, err error)
	// Get retrieves the FabricFollowerChannel from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.FabricFollowerChannel, error)
	FabricFollowerChannelListerExpansion
}

// fabricFollowerChannelLister implements the FabricFollowerChannelLister interface.
type fabricFollowerChannelLister struct {
	indexer cache.Indexer
}

// NewFabricFollowerChannelLister returns a new FabricFollowerChannelLister.
func NewFabricFollowerChannelLister(indexer cache.Indexer) FabricFollowerChannelLister {
	return &fabricFollowerChannelLister{indexer: indexer}
}

// List lists all FabricFollowerChannels in the indexer.
func (s *fabricFollowerChannelLister) List(selector labels.Selector) (ret []*v1alpha1.FabricFollowerChannel, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.FabricFollowerChannel))
	})
	return ret, err
}

// Get retrieves the FabricFollowerChannel from the index for a given name.
func (s *fabricFollowerChannelLister) Get(name string) (*v1alpha1.FabricFollowerChannel, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("fabricfollowerchannel"), name)
	}
	return obj.(*v1alpha1.FabricFollowerChannel), nil
}
