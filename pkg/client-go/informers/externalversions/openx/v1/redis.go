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

package v1

import (
	"context"
	time "time"

	openxv1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	versioned "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned"
	internalinterfaces "github.com/kzz45/neverdown/pkg/client-go/informers/externalversions/internalinterfaces"
	v1 "github.com/kzz45/neverdown/pkg/client-go/listers/openx/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// RedisInformer provides access to a shared informer and lister for
// Redises.
type RedisInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.RedisLister
}

type redisInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewRedisInformer constructs a new informer for Redis type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewRedisInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredRedisInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredRedisInformer constructs a new informer for Redis type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredRedisInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OpenxV1().Redises(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OpenxV1().Redises(namespace).Watch(context.TODO(), options)
			},
		},
		&openxv1.Redis{},
		resyncPeriod,
		indexers,
	)
}

func (f *redisInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredRedisInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *redisInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&openxv1.Redis{}, f.defaultInformer)
}

func (f *redisInformer) Lister() v1.RedisLister {
	return v1.NewRedisLister(f.Informer().GetIndexer())
}
