// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	hlfkungfusoftwareesv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	versioned "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	internalinterfaces "github.com/kfsoftware/hlf-operator/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kfsoftware/hlf-operator/pkg/client/listers/hlf.kungfusoftware.es/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// FabricOrderingServiceInformer provides access to a shared informer and lister for
// FabricOrderingServices.
type FabricOrderingServiceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.FabricOrderingServiceLister
}

type fabricOrderingServiceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewFabricOrderingServiceInformer constructs a new informer for FabricOrderingService type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFabricOrderingServiceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredFabricOrderingServiceInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredFabricOrderingServiceInformer constructs a new informer for FabricOrderingService type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredFabricOrderingServiceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HlfV1alpha1().FabricOrderingServices(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HlfV1alpha1().FabricOrderingServices(namespace).Watch(context.TODO(), options)
			},
		},
		&hlfkungfusoftwareesv1alpha1.FabricOrderingService{},
		resyncPeriod,
		indexers,
	)
}

func (f *fabricOrderingServiceInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredFabricOrderingServiceInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *fabricOrderingServiceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&hlfkungfusoftwareesv1alpha1.FabricOrderingService{}, f.defaultInformer)
}

func (f *fabricOrderingServiceInformer) Lister() v1alpha1.FabricOrderingServiceLister {
	return v1alpha1.NewFabricOrderingServiceLister(f.Informer().GetIndexer())
}
