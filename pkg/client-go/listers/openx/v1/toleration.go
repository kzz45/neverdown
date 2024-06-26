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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// TolerationLister helps list Tolerations.
// All objects returned here must be treated as read-only.
type TolerationLister interface {
	// List lists all Tolerations in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Toleration, err error)
	// Tolerations returns an object that can list and get Tolerations.
	Tolerations(namespace string) TolerationNamespaceLister
	TolerationListerExpansion
}

// tolerationLister implements the TolerationLister interface.
type tolerationLister struct {
	indexer cache.Indexer
}

// NewTolerationLister returns a new TolerationLister.
func NewTolerationLister(indexer cache.Indexer) TolerationLister {
	return &tolerationLister{indexer: indexer}
}

// List lists all Tolerations in the indexer.
func (s *tolerationLister) List(selector labels.Selector) (ret []*v1.Toleration, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Toleration))
	})
	return ret, err
}

// Tolerations returns an object that can list and get Tolerations.
func (s *tolerationLister) Tolerations(namespace string) TolerationNamespaceLister {
	return tolerationNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// TolerationNamespaceLister helps list and get Tolerations.
// All objects returned here must be treated as read-only.
type TolerationNamespaceLister interface {
	// List lists all Tolerations in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Toleration, err error)
	// Get retrieves the Toleration from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Toleration, error)
	TolerationNamespaceListerExpansion
}

// tolerationNamespaceLister implements the TolerationNamespaceLister
// interface.
type tolerationNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Tolerations in the indexer for a given namespace.
func (s tolerationNamespaceLister) List(selector labels.Selector) (ret []*v1.Toleration, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Toleration))
	})
	return ret, err
}

// Get retrieves the Toleration from the indexer for a given namespace and name.
func (s tolerationNamespaceLister) Get(name string) (*v1.Toleration, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("toleration"), name)
	}
	return obj.(*v1.Toleration), nil
}
