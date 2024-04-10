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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	scheme "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// GroupVersionKindRulesGetter has a method to return a GroupVersionKindRuleInterface.
// A group's client should implement this interface.
type GroupVersionKindRulesGetter interface {
	GroupVersionKindRules(namespace string) GroupVersionKindRuleInterface
}

// GroupVersionKindRuleInterface has methods to work with GroupVersionKindRule resources.
type GroupVersionKindRuleInterface interface {
	Create(ctx context.Context, groupVersionKindRule *v1.GroupVersionKindRule, opts metav1.CreateOptions) (*v1.GroupVersionKindRule, error)
	Update(ctx context.Context, groupVersionKindRule *v1.GroupVersionKindRule, opts metav1.UpdateOptions) (*v1.GroupVersionKindRule, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.GroupVersionKindRule, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.GroupVersionKindRuleList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.GroupVersionKindRule, err error)
	GroupVersionKindRuleExpansion
}

// groupVersionKindRules implements GroupVersionKindRuleInterface
type groupVersionKindRules struct {
	client rest.Interface
	ns     string
}

// newGroupVersionKindRules returns a GroupVersionKindRules
func newGroupVersionKindRules(c *RbacV1Client, namespace string) *groupVersionKindRules {
	return &groupVersionKindRules{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the groupVersionKindRule, and returns the corresponding groupVersionKindRule object, and an error if there is any.
func (c *groupVersionKindRules) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.GroupVersionKindRule, err error) {
	result = &v1.GroupVersionKindRule{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("groupversionkindrules").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of GroupVersionKindRules that match those selectors.
func (c *groupVersionKindRules) List(ctx context.Context, opts metav1.ListOptions) (result *v1.GroupVersionKindRuleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.GroupVersionKindRuleList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("groupversionkindrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested groupVersionKindRules.
func (c *groupVersionKindRules) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("groupversionkindrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a groupVersionKindRule and creates it.  Returns the server's representation of the groupVersionKindRule, and an error, if there is any.
func (c *groupVersionKindRules) Create(ctx context.Context, groupVersionKindRule *v1.GroupVersionKindRule, opts metav1.CreateOptions) (result *v1.GroupVersionKindRule, err error) {
	result = &v1.GroupVersionKindRule{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("groupversionkindrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(groupVersionKindRule).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a groupVersionKindRule and updates it. Returns the server's representation of the groupVersionKindRule, and an error, if there is any.
func (c *groupVersionKindRules) Update(ctx context.Context, groupVersionKindRule *v1.GroupVersionKindRule, opts metav1.UpdateOptions) (result *v1.GroupVersionKindRule, err error) {
	result = &v1.GroupVersionKindRule{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("groupversionkindrules").
		Name(groupVersionKindRule.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(groupVersionKindRule).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the groupVersionKindRule and deletes it. Returns an error if one occurs.
func (c *groupVersionKindRules) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("groupversionkindrules").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *groupVersionKindRules) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("groupversionkindrules").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched groupVersionKindRule.
func (c *groupVersionKindRules) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.GroupVersionKindRule, err error) {
	result = &v1.GroupVersionKindRule{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("groupversionkindrules").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}