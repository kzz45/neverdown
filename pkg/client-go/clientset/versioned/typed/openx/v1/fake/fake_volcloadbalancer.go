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

package fake

import (
	"context"

	v1 "github.com/kzz45/neverdown/pkg/apis/openx/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVolcLoadBalancers implements VolcLoadBalancerInterface
type FakeVolcLoadBalancers struct {
	Fake *FakeOpenxV1
	ns   string
}

var volcloadbalancersResource = v1.SchemeGroupVersion.WithResource("volcloadbalancers")

var volcloadbalancersKind = v1.SchemeGroupVersion.WithKind("VolcLoadBalancer")

// Get takes name of the volcLoadBalancer, and returns the corresponding volcLoadBalancer object, and an error if there is any.
func (c *FakeVolcLoadBalancers) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.VolcLoadBalancer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(volcloadbalancersResource, c.ns, name), &v1.VolcLoadBalancer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.VolcLoadBalancer), err
}

// List takes label and field selectors, and returns the list of VolcLoadBalancers that match those selectors.
func (c *FakeVolcLoadBalancers) List(ctx context.Context, opts metav1.ListOptions) (result *v1.VolcLoadBalancerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(volcloadbalancersResource, volcloadbalancersKind, c.ns, opts), &v1.VolcLoadBalancerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.VolcLoadBalancerList{ListMeta: obj.(*v1.VolcLoadBalancerList).ListMeta}
	for _, item := range obj.(*v1.VolcLoadBalancerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested volcLoadBalancers.
func (c *FakeVolcLoadBalancers) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(volcloadbalancersResource, c.ns, opts))

}

// Create takes the representation of a volcLoadBalancer and creates it.  Returns the server's representation of the volcLoadBalancer, and an error, if there is any.
func (c *FakeVolcLoadBalancers) Create(ctx context.Context, volcLoadBalancer *v1.VolcLoadBalancer, opts metav1.CreateOptions) (result *v1.VolcLoadBalancer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(volcloadbalancersResource, c.ns, volcLoadBalancer), &v1.VolcLoadBalancer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.VolcLoadBalancer), err
}

// Update takes the representation of a volcLoadBalancer and updates it. Returns the server's representation of the volcLoadBalancer, and an error, if there is any.
func (c *FakeVolcLoadBalancers) Update(ctx context.Context, volcLoadBalancer *v1.VolcLoadBalancer, opts metav1.UpdateOptions) (result *v1.VolcLoadBalancer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(volcloadbalancersResource, c.ns, volcLoadBalancer), &v1.VolcLoadBalancer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.VolcLoadBalancer), err
}

// Delete takes name of the volcLoadBalancer and deletes it. Returns an error if one occurs.
func (c *FakeVolcLoadBalancers) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(volcloadbalancersResource, c.ns, name, opts), &v1.VolcLoadBalancer{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVolcLoadBalancers) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(volcloadbalancersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1.VolcLoadBalancerList{})
	return err
}

// Patch applies the patch and returns the patched volcLoadBalancer.
func (c *FakeVolcLoadBalancers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.VolcLoadBalancer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(volcloadbalancersResource, c.ns, name, pt, data, subresources...), &v1.VolcLoadBalancer{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.VolcLoadBalancer), err
}
