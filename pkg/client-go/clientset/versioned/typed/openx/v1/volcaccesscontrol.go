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

	v1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	scheme "github.com/kzz45/neverdown/pkg/client-go/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VolcAccessControlsGetter has a method to return a VolcAccessControlInterface.
// A group's client should implement this interface.
type VolcAccessControlsGetter interface {
	VolcAccessControls(namespace string) VolcAccessControlInterface
}

// VolcAccessControlInterface has methods to work with VolcAccessControl resources.
type VolcAccessControlInterface interface {
	Create(ctx context.Context, volcAccessControl *v1.VolcAccessControl, opts metav1.CreateOptions) (*v1.VolcAccessControl, error)
	Update(ctx context.Context, volcAccessControl *v1.VolcAccessControl, opts metav1.UpdateOptions) (*v1.VolcAccessControl, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.VolcAccessControl, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.VolcAccessControlList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.VolcAccessControl, err error)
	VolcAccessControlExpansion
}

// volcAccessControls implements VolcAccessControlInterface
type volcAccessControls struct {
	client rest.Interface
	ns     string
}

// newVolcAccessControls returns a VolcAccessControls
func newVolcAccessControls(c *OpenxV1Client, namespace string) *volcAccessControls {
	return &volcAccessControls{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the volcAccessControl, and returns the corresponding volcAccessControl object, and an error if there is any.
func (c *volcAccessControls) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.VolcAccessControl, err error) {
	result = &v1.VolcAccessControl{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("volcaccesscontrols").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VolcAccessControls that match those selectors.
func (c *volcAccessControls) List(ctx context.Context, opts metav1.ListOptions) (result *v1.VolcAccessControlList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.VolcAccessControlList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("volcaccesscontrols").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested volcAccessControls.
func (c *volcAccessControls) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("volcaccesscontrols").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a volcAccessControl and creates it.  Returns the server's representation of the volcAccessControl, and an error, if there is any.
func (c *volcAccessControls) Create(ctx context.Context, volcAccessControl *v1.VolcAccessControl, opts metav1.CreateOptions) (result *v1.VolcAccessControl, err error) {
	result = &v1.VolcAccessControl{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("volcaccesscontrols").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(volcAccessControl).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a volcAccessControl and updates it. Returns the server's representation of the volcAccessControl, and an error, if there is any.
func (c *volcAccessControls) Update(ctx context.Context, volcAccessControl *v1.VolcAccessControl, opts metav1.UpdateOptions) (result *v1.VolcAccessControl, err error) {
	result = &v1.VolcAccessControl{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("volcaccesscontrols").
		Name(volcAccessControl.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(volcAccessControl).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the volcAccessControl and deletes it. Returns an error if one occurs.
func (c *volcAccessControls) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("volcaccesscontrols").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *volcAccessControls) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("volcaccesscontrols").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched volcAccessControl.
func (c *volcAccessControls) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.VolcAccessControl, err error) {
	result = &v1.VolcAccessControl{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("volcaccesscontrols").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}