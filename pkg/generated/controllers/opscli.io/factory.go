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

package opscli

import (
	"context"
	"time"

	clientset "github.com/181192/ops-cli/pkg/generated/clientset/versioned"
	scheme "github.com/181192/ops-cli/pkg/generated/clientset/versioned/scheme"
	informers "github.com/181192/ops-cli/pkg/generated/informers/externalversions"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/schemes"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func init() {
	scheme.AddToScheme(schemes.All)
}

type Factory struct {
	synced            bool
	informerFactory   informers.SharedInformerFactory
	clientset         clientset.Interface
	controllerManager *generic.ControllerManager
	threadiness       map[schema.GroupVersionKind]int
}

func NewFactoryFromConfigOrDie(config *rest.Config) *Factory {
	f, err := NewFactoryFromConfig(config)
	if err != nil {
		panic(err)
	}
	return f
}

func NewFactoryFromConfig(config *rest.Config) (*Factory, error) {
	return NewFactoryFromConfigWithOptions(config, nil)
}

func NewFactoryFromConfigWithNamespace(config *rest.Config, namespace string) (*Factory, error) {
	return NewFactoryFromConfigWithOptions(config, &FactoryOptions{
		Namespace: namespace,
	})
}

type FactoryOptions struct {
	Namespace string
	Resync    time.Duration
}

func NewFactoryFromConfigWithOptions(config *rest.Config, opts *FactoryOptions) (*Factory, error) {
	if opts == nil {
		opts = &FactoryOptions{}
	}

	cs, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	resync := opts.Resync
	if resync == 0 {
		resync = 2 * time.Hour
	}

	if opts.Namespace == "" {
		informerFactory := informers.NewSharedInformerFactory(cs, resync)
		return NewFactory(cs, informerFactory), nil
	}

	informerFactory := informers.NewSharedInformerFactoryWithOptions(cs, resync, informers.WithNamespace(opts.Namespace))
	return NewFactory(cs, informerFactory), nil
}

func NewFactory(clientset clientset.Interface, informerFactory informers.SharedInformerFactory) *Factory {
	return &Factory{
		threadiness:       map[schema.GroupVersionKind]int{},
		controllerManager: &generic.ControllerManager{},
		clientset:         clientset,
		informerFactory:   informerFactory,
	}
}

func (c *Factory) Controllers() map[schema.GroupVersionKind]*generic.Controller {
	return c.controllerManager.Controllers()
}

func (c *Factory) SetThreadiness(gvk schema.GroupVersionKind, threadiness int) {
	c.threadiness[gvk] = threadiness
}

func (c *Factory) Sync(ctx context.Context) error {
	c.informerFactory.Start(ctx.Done())
	c.informerFactory.WaitForCacheSync(ctx.Done())
	return nil
}

func (c *Factory) Start(ctx context.Context, defaultThreadiness int) error {
	if err := c.Sync(ctx); err != nil {
		return err
	}

	return c.controllerManager.Start(ctx, defaultThreadiness, c.threadiness)
}

func (c *Factory) Opscli() Interface {
	return New(c.controllerManager, c.informerFactory.Opscli(), c.clientset)
}
