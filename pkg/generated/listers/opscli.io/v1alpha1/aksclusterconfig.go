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
	v1alpha1 "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// AKSClusterConfigLister helps list AKSClusterConfigs.
type AKSClusterConfigLister interface {
	// List lists all AKSClusterConfigs in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.AKSClusterConfig, err error)
	// Get retrieves the AKSClusterConfig from the index for a given name.
	Get(name string) (*v1alpha1.AKSClusterConfig, error)
	AKSClusterConfigListerExpansion
}

// aKSClusterConfigLister implements the AKSClusterConfigLister interface.
type aKSClusterConfigLister struct {
	indexer cache.Indexer
}

// NewAKSClusterConfigLister returns a new AKSClusterConfigLister.
func NewAKSClusterConfigLister(indexer cache.Indexer) AKSClusterConfigLister {
	return &aKSClusterConfigLister{indexer: indexer}
}

// List lists all AKSClusterConfigs in the indexer.
func (s *aKSClusterConfigLister) List(selector labels.Selector) (ret []*v1alpha1.AKSClusterConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.AKSClusterConfig))
	})
	return ret, err
}

// Get retrieves the AKSClusterConfig from the index for a given name.
func (s *aKSClusterConfigLister) Get(name string) (*v1alpha1.AKSClusterConfig, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("aksclusterconfig"), name)
	}
	return obj.(*v1alpha1.AKSClusterConfig), nil
}
