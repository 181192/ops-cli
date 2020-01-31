/*
Copyright 2020 github.com/181192.

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

// Code generated by main. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterConfigs implements ClusterConfigInterface
type FakeClusterConfigs struct {
	Fake *FakeOpscliV1alpha1
}

var clusterconfigsResource = schema.GroupVersionResource{Group: "opscli.io", Version: "v1alpha1", Resource: "clusterconfigs"}

var clusterconfigsKind = schema.GroupVersionKind{Group: "opscli.io", Version: "v1alpha1", Kind: "ClusterConfig"}

// Get takes name of the clusterConfig, and returns the corresponding clusterConfig object, and an error if there is any.
func (c *FakeClusterConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.ClusterConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusterconfigsResource, name), &v1alpha1.ClusterConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterConfig), err
}

// List takes label and field selectors, and returns the list of ClusterConfigs that match those selectors.
func (c *FakeClusterConfigs) List(opts v1.ListOptions) (result *v1alpha1.ClusterConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusterconfigsResource, clusterconfigsKind, opts), &v1alpha1.ClusterConfigList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ClusterConfigList{ListMeta: obj.(*v1alpha1.ClusterConfigList).ListMeta}
	for _, item := range obj.(*v1alpha1.ClusterConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterConfigs.
func (c *FakeClusterConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clusterconfigsResource, opts))
}

// Create takes the representation of a clusterConfig and creates it.  Returns the server's representation of the clusterConfig, and an error, if there is any.
func (c *FakeClusterConfigs) Create(clusterConfig *v1alpha1.ClusterConfig) (result *v1alpha1.ClusterConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clusterconfigsResource, clusterConfig), &v1alpha1.ClusterConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterConfig), err
}

// Update takes the representation of a clusterConfig and updates it. Returns the server's representation of the clusterConfig, and an error, if there is any.
func (c *FakeClusterConfigs) Update(clusterConfig *v1alpha1.ClusterConfig) (result *v1alpha1.ClusterConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clusterconfigsResource, clusterConfig), &v1alpha1.ClusterConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterConfig), err
}

// Delete takes name of the clusterConfig and deletes it. Returns an error if one occurs.
func (c *FakeClusterConfigs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(clusterconfigsResource, name), &v1alpha1.ClusterConfig{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clusterconfigsResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.ClusterConfigList{})
	return err
}

// Patch applies the patch and returns the patched clusterConfig.
func (c *FakeClusterConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ClusterConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clusterconfigsResource, name, pt, data, subresources...), &v1alpha1.ClusterConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterConfig), err
}