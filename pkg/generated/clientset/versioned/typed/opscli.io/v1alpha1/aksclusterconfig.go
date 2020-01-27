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

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	scheme "github.com/181192/ops-cli/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AKSClusterConfigsGetter has a method to return a AKSClusterConfigInterface.
// A group's client should implement this interface.
type AKSClusterConfigsGetter interface {
	AKSClusterConfigs() AKSClusterConfigInterface
}

// AKSClusterConfigInterface has methods to work with AKSClusterConfig resources.
type AKSClusterConfigInterface interface {
	Create(*v1alpha1.AKSClusterConfig) (*v1alpha1.AKSClusterConfig, error)
	Update(*v1alpha1.AKSClusterConfig) (*v1alpha1.AKSClusterConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.AKSClusterConfig, error)
	List(opts v1.ListOptions) (*v1alpha1.AKSClusterConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AKSClusterConfig, err error)
	AKSClusterConfigExpansion
}

// aKSClusterConfigs implements AKSClusterConfigInterface
type aKSClusterConfigs struct {
	client rest.Interface
}

// newAKSClusterConfigs returns a AKSClusterConfigs
func newAKSClusterConfigs(c *OpscliV1alpha1Client) *aKSClusterConfigs {
	return &aKSClusterConfigs{
		client: c.RESTClient(),
	}
}

// Get takes name of the aKSClusterConfig, and returns the corresponding aKSClusterConfig object, and an error if there is any.
func (c *aKSClusterConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.AKSClusterConfig, err error) {
	result = &v1alpha1.AKSClusterConfig{}
	err = c.client.Get().
		Resource("aksclusterconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AKSClusterConfigs that match those selectors.
func (c *aKSClusterConfigs) List(opts v1.ListOptions) (result *v1alpha1.AKSClusterConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.AKSClusterConfigList{}
	err = c.client.Get().
		Resource("aksclusterconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested aKSClusterConfigs.
func (c *aKSClusterConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("aksclusterconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a aKSClusterConfig and creates it.  Returns the server's representation of the aKSClusterConfig, and an error, if there is any.
func (c *aKSClusterConfigs) Create(aKSClusterConfig *v1alpha1.AKSClusterConfig) (result *v1alpha1.AKSClusterConfig, err error) {
	result = &v1alpha1.AKSClusterConfig{}
	err = c.client.Post().
		Resource("aksclusterconfigs").
		Body(aKSClusterConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a aKSClusterConfig and updates it. Returns the server's representation of the aKSClusterConfig, and an error, if there is any.
func (c *aKSClusterConfigs) Update(aKSClusterConfig *v1alpha1.AKSClusterConfig) (result *v1alpha1.AKSClusterConfig, err error) {
	result = &v1alpha1.AKSClusterConfig{}
	err = c.client.Put().
		Resource("aksclusterconfigs").
		Name(aKSClusterConfig.Name).
		Body(aKSClusterConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the aKSClusterConfig and deletes it. Returns an error if one occurs.
func (c *aKSClusterConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("aksclusterconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *aKSClusterConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("aksclusterconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched aKSClusterConfig.
func (c *aKSClusterConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AKSClusterConfig, err error) {
	result = &v1alpha1.AKSClusterConfig{}
	err = c.client.Patch(pt).
		Resource("aksclusterconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
