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

// MachineAutoScalingGroupsGetter has a method to return a MachineAutoScalingGroupInterface.
// A group's client should implement this interface.
type MachineAutoScalingGroupsGetter interface {
	MachineAutoScalingGroups() MachineAutoScalingGroupInterface
}

// MachineAutoScalingGroupInterface has methods to work with MachineAutoScalingGroup resources.
type MachineAutoScalingGroupInterface interface {
	Create(*v1beta1.MachineAutoScalingGroup) (*v1beta1.MachineAutoScalingGroup, error)
	Update(*v1beta1.MachineAutoScalingGroup) (*v1beta1.MachineAutoScalingGroup, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1beta1.MachineAutoScalingGroup, error)
	List(opts v1.ListOptions) (*v1beta1.MachineAutoScalingGroupList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.MachineAutoScalingGroup, err error)
	MachineAutoScalingGroupExpansion
}

// machineAutoScalingGroups implements MachineAutoScalingGroupInterface
type machineAutoScalingGroups struct {
	client rest.Interface
}

// newMachineAutoScalingGroups returns a MachineAutoScalingGroups
func newMachineAutoScalingGroups(c *ResourceV1beta1Client) *machineAutoScalingGroups {
	return &machineAutoScalingGroups{
		client: c.RESTClient(),
	}
}

// Get takes name of the machineAutoScalingGroup, and returns the corresponding machineAutoScalingGroup object, and an error if there is any.
func (c *machineAutoScalingGroups) Get(name string, options v1.GetOptions) (result *v1beta1.MachineAutoScalingGroup, err error) {
	result = &v1beta1.MachineAutoScalingGroup{}
	err = c.client.Get().
		Resource("machineautoscalinggroups").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MachineAutoScalingGroups that match those selectors.
func (c *machineAutoScalingGroups) List(opts v1.ListOptions) (result *v1beta1.MachineAutoScalingGroupList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.MachineAutoScalingGroupList{}
	err = c.client.Get().
		Resource("machineautoscalinggroups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested machineAutoScalingGroups.
func (c *machineAutoScalingGroups) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("machineautoscalinggroups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a machineAutoScalingGroup and creates it.  Returns the server's representation of the machineAutoScalingGroup, and an error, if there is any.
func (c *machineAutoScalingGroups) Create(machineAutoScalingGroup *v1beta1.MachineAutoScalingGroup) (result *v1beta1.MachineAutoScalingGroup, err error) {
	result = &v1beta1.MachineAutoScalingGroup{}
	err = c.client.Post().
		Resource("machineautoscalinggroups").
		Body(machineAutoScalingGroup).
		Do().
		Into(result)
	return
}

// Update takes the representation of a machineAutoScalingGroup and updates it. Returns the server's representation of the machineAutoScalingGroup, and an error, if there is any.
func (c *machineAutoScalingGroups) Update(machineAutoScalingGroup *v1beta1.MachineAutoScalingGroup) (result *v1beta1.MachineAutoScalingGroup, err error) {
	result = &v1beta1.MachineAutoScalingGroup{}
	err = c.client.Put().
		Resource("machineautoscalinggroups").
		Name(machineAutoScalingGroup.Name).
		Body(machineAutoScalingGroup).
		Do().
		Into(result)
	return
}

// Delete takes name of the machineAutoScalingGroup and deletes it. Returns an error if one occurs.
func (c *machineAutoScalingGroups) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("machineautoscalinggroups").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *machineAutoScalingGroups) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("machineautoscalinggroups").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched machineAutoScalingGroup.
func (c *machineAutoScalingGroups) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.MachineAutoScalingGroup, err error) {
	result = &v1beta1.MachineAutoScalingGroup{}
	err = c.client.Patch(pt).
		Resource("machineautoscalinggroups").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
