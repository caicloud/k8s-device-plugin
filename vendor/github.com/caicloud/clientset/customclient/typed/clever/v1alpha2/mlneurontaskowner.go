/*
Copyright 2020 caicloud authors. All rights reserved.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha2

import (
	"time"

	scheme "github.com/caicloud/clientset/customclient/scheme"
	v1alpha2 "github.com/caicloud/clientset/pkg/apis/clever/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MLNeuronTaskOwnersGetter has a method to return a MLNeuronTaskOwnerInterface.
// A group's client should implement this interface.
type MLNeuronTaskOwnersGetter interface {
	MLNeuronTaskOwners(namespace string) MLNeuronTaskOwnerInterface
}

// MLNeuronTaskOwnerInterface has methods to work with MLNeuronTaskOwner resources.
type MLNeuronTaskOwnerInterface interface {
	Create(*v1alpha2.MLNeuronTaskOwner) (*v1alpha2.MLNeuronTaskOwner, error)
	Update(*v1alpha2.MLNeuronTaskOwner) (*v1alpha2.MLNeuronTaskOwner, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha2.MLNeuronTaskOwner, error)
	List(opts v1.ListOptions) (*v1alpha2.MLNeuronTaskOwnerList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.MLNeuronTaskOwner, err error)
	MLNeuronTaskOwnerExpansion
}

// mLNeuronTaskOwners implements MLNeuronTaskOwnerInterface
type mLNeuronTaskOwners struct {
	client rest.Interface
	ns     string
}

// newMLNeuronTaskOwners returns a MLNeuronTaskOwners
func newMLNeuronTaskOwners(c *CleverV1alpha2Client, namespace string) *mLNeuronTaskOwners {
	return &mLNeuronTaskOwners{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the mLNeuronTaskOwner, and returns the corresponding mLNeuronTaskOwner object, and an error if there is any.
func (c *mLNeuronTaskOwners) Get(name string, options v1.GetOptions) (result *v1alpha2.MLNeuronTaskOwner, err error) {
	result = &v1alpha2.MLNeuronTaskOwner{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mlneurontaskowners").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MLNeuronTaskOwners that match those selectors.
func (c *mLNeuronTaskOwners) List(opts v1.ListOptions) (result *v1alpha2.MLNeuronTaskOwnerList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha2.MLNeuronTaskOwnerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mlneurontaskowners").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested mLNeuronTaskOwners.
func (c *mLNeuronTaskOwners) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("mlneurontaskowners").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a mLNeuronTaskOwner and creates it.  Returns the server's representation of the mLNeuronTaskOwner, and an error, if there is any.
func (c *mLNeuronTaskOwners) Create(mLNeuronTaskOwner *v1alpha2.MLNeuronTaskOwner) (result *v1alpha2.MLNeuronTaskOwner, err error) {
	result = &v1alpha2.MLNeuronTaskOwner{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("mlneurontaskowners").
		Body(mLNeuronTaskOwner).
		Do().
		Into(result)
	return
}

// Update takes the representation of a mLNeuronTaskOwner and updates it. Returns the server's representation of the mLNeuronTaskOwner, and an error, if there is any.
func (c *mLNeuronTaskOwners) Update(mLNeuronTaskOwner *v1alpha2.MLNeuronTaskOwner) (result *v1alpha2.MLNeuronTaskOwner, err error) {
	result = &v1alpha2.MLNeuronTaskOwner{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mlneurontaskowners").
		Name(mLNeuronTaskOwner.Name).
		Body(mLNeuronTaskOwner).
		Do().
		Into(result)
	return
}

// Delete takes name of the mLNeuronTaskOwner and deletes it. Returns an error if one occurs.
func (c *mLNeuronTaskOwners) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mlneurontaskowners").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *mLNeuronTaskOwners) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mlneurontaskowners").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched mLNeuronTaskOwner.
func (c *mLNeuronTaskOwners) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.MLNeuronTaskOwner, err error) {
	result = &v1alpha2.MLNeuronTaskOwner{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("mlneurontaskowners").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
