/*
Copyright 2020 caicloud authors. All rights reserved.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"time"

	scheme "github.com/caicloud/clientset/customclient/scheme"
	v1beta1 "github.com/caicloud/clientset/pkg/apis/resource/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// InfraNetworksGetter has a method to return a InfraNetworkInterface.
// A group's client should implement this interface.
type InfraNetworksGetter interface {
	InfraNetworks() InfraNetworkInterface
}

// InfraNetworkInterface has methods to work with InfraNetwork resources.
type InfraNetworkInterface interface {
	Create(*v1beta1.InfraNetwork) (*v1beta1.InfraNetwork, error)
	Update(*v1beta1.InfraNetwork) (*v1beta1.InfraNetwork, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1beta1.InfraNetwork, error)
	List(opts v1.ListOptions) (*v1beta1.InfraNetworkList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.InfraNetwork, err error)
	InfraNetworkExpansion
}

// infraNetworks implements InfraNetworkInterface
type infraNetworks struct {
	client rest.Interface
}

// newInfraNetworks returns a InfraNetworks
func newInfraNetworks(c *ResourceV1beta1Client) *infraNetworks {
	return &infraNetworks{
		client: c.RESTClient(),
	}
}

// Get takes name of the infraNetwork, and returns the corresponding infraNetwork object, and an error if there is any.
func (c *infraNetworks) Get(name string, options v1.GetOptions) (result *v1beta1.InfraNetwork, err error) {
	result = &v1beta1.InfraNetwork{}
	err = c.client.Get().
		Resource("infranetworks").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of InfraNetworks that match those selectors.
func (c *infraNetworks) List(opts v1.ListOptions) (result *v1beta1.InfraNetworkList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.InfraNetworkList{}
	err = c.client.Get().
		Resource("infranetworks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested infraNetworks.
func (c *infraNetworks) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("infranetworks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a infraNetwork and creates it.  Returns the server's representation of the infraNetwork, and an error, if there is any.
func (c *infraNetworks) Create(infraNetwork *v1beta1.InfraNetwork) (result *v1beta1.InfraNetwork, err error) {
	result = &v1beta1.InfraNetwork{}
	err = c.client.Post().
		Resource("infranetworks").
		Body(infraNetwork).
		Do().
		Into(result)
	return
}

// Update takes the representation of a infraNetwork and updates it. Returns the server's representation of the infraNetwork, and an error, if there is any.
func (c *infraNetworks) Update(infraNetwork *v1beta1.InfraNetwork) (result *v1beta1.InfraNetwork, err error) {
	result = &v1beta1.InfraNetwork{}
	err = c.client.Put().
		Resource("infranetworks").
		Name(infraNetwork.Name).
		Body(infraNetwork).
		Do().
		Into(result)
	return
}

// Delete takes name of the infraNetwork and deletes it. Returns an error if one occurs.
func (c *infraNetworks) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("infranetworks").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *infraNetworks) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("infranetworks").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched infraNetwork.
func (c *infraNetworks) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.InfraNetwork, err error) {
	result = &v1beta1.InfraNetwork{}
	err = c.client.Patch(pt).
		Resource("infranetworks").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
