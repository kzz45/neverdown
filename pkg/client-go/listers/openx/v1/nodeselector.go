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

// NodeSelectorLister helps list NodeSelectors.
// All objects returned here must be treated as read-only.
type NodeSelectorLister interface {
	// List lists all NodeSelectors in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.NodeSelector, err error)
	// NodeSelectors returns an object that can list and get NodeSelectors.
	NodeSelectors(namespace string) NodeSelectorNamespaceLister
	NodeSelectorListerExpansion
}

// nodeSelectorLister implements the NodeSelectorLister interface.
type nodeSelectorLister struct {
	indexer cache.Indexer
}

// NewNodeSelectorLister returns a new NodeSelectorLister.
func NewNodeSelectorLister(indexer cache.Indexer) NodeSelectorLister {
	return &nodeSelectorLister{indexer: indexer}
}

// List lists all NodeSelectors in the indexer.
func (s *nodeSelectorLister) List(selector labels.Selector) (ret []*v1.NodeSelector, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.NodeSelector))
	})
	return ret, err
}

// NodeSelectors returns an object that can list and get NodeSelectors.
func (s *nodeSelectorLister) NodeSelectors(namespace string) NodeSelectorNamespaceLister {
	return nodeSelectorNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// NodeSelectorNamespaceLister helps list and get NodeSelectors.
// All objects returned here must be treated as read-only.
type NodeSelectorNamespaceLister interface {
	// List lists all NodeSelectors in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.NodeSelector, err error)
	// Get retrieves the NodeSelector from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.NodeSelector, error)
	NodeSelectorNamespaceListerExpansion
}

// nodeSelectorNamespaceLister implements the NodeSelectorNamespaceLister
// interface.
type nodeSelectorNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all NodeSelectors in the indexer for a given namespace.
func (s nodeSelectorNamespaceLister) List(selector labels.Selector) (ret []*v1.NodeSelector, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.NodeSelector))
	})
	return ret, err
}

// Get retrieves the NodeSelector from the indexer for a given namespace and name.
func (s nodeSelectorNamespaceLister) Get(name string) (*v1.NodeSelector, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("nodeselector"), name)
	}
	return obj.(*v1.NodeSelector), nil
}
