/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package internalversion

import (
	"context"
	time "time"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	internalinterfaces "github.com/kzz45/neverdown/pkg/client-go/informers/internalversion/internalinterfaces"
	internalversion "github.com/kzz45/neverdown/pkg/client-go/listers/v1/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	internalclientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

// RepositoryInformer provides access to a shared informer and lister for
// Repositories.
type RepositoryInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() internalversion.RepositoryLister
}

type repositoryInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewRepositoryInformer constructs a new informer for Repository type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewRepositoryInformer(client internalclientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredRepositoryInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredRepositoryInformer constructs a new informer for Repository type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredRepositoryInformer(client internalclientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.V1().Repositories(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.V1().Repositories(namespace).Watch(context.TODO(), options)
			},
		},
		&jingxv1.Repository{},
		resyncPeriod,
		indexers,
	)
}

func (f *repositoryInformer) defaultInformer(client internalclientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredRepositoryInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *repositoryInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&jingxv1.Repository{}, f.defaultInformer)
}

func (f *repositoryInformer) Lister() internalversion.RepositoryLister {
	return internalversion.NewRepositoryLister(f.Informer().GetIndexer())
}
