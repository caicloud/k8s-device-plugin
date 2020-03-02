/*
Copyright 2020 caicloud authors. All rights reserved.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	scheme "github.com/caicloud/clientset/customclient/scheme"
	v1alpha1 "github.com/caicloud/clientset/pkg/apis/workload/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// WorkloadRevisionsGetter has a method to return a WorkloadRevisionInterface.
// A group's client should implement this interface.
type WorkloadRevisionsGetter interface {
	WorkloadRevisions(namespace string) WorkloadRevisionInterface
}

// WorkloadRevisionInterface has methods to work with WorkloadRevision resources.
type WorkloadRevisionInterface interface {
	Create(*v1alpha1.WorkloadRevision) (*v1alpha1.WorkloadRevision, error)
	Update(*v1alpha1.WorkloadRevision) (*v1alpha1.WorkloadRevision, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.WorkloadRevision, error)
	List(opts v1.ListOptions) (*v1alpha1.WorkloadRevisionList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.WorkloadRevision, err error)
	WorkloadRevisionExpansion
}

// workloadRevisions implements WorkloadRevisionInterface
type workloadRevisions struct {
	client rest.Interface
	ns     string
}

// newWorkloadRevisions returns a WorkloadRevisions
func newWorkloadRevisions(c *WorkloadV1alpha1Client, namespace string) *workloadRevisions {
	return &workloadRevisions{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the workloadRevision, and returns the corresponding workloadRevision object, and an error if there is any.
func (c *workloadRevisions) Get(name string, options v1.GetOptions) (result *v1alpha1.WorkloadRevision, err error) {
	result = &v1alpha1.WorkloadRevision{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("workloadrevisions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of WorkloadRevisions that match those selectors.
func (c *workloadRevisions) List(opts v1.ListOptions) (result *v1alpha1.WorkloadRevisionList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.WorkloadRevisionList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("workloadrevisions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested workloadRevisions.
func (c *workloadRevisions) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("workloadrevisions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a workloadRevision and creates it.  Returns the server's representation of the workloadRevision, and an error, if there is any.
func (c *workloadRevisions) Create(workloadRevision *v1alpha1.WorkloadRevision) (result *v1alpha1.WorkloadRevision, err error) {
	result = &v1alpha1.WorkloadRevision{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("workloadrevisions").
		Body(workloadRevision).
		Do().
		Into(result)
	return
}

// Update takes the representation of a workloadRevision and updates it. Returns the server's representation of the workloadRevision, and an error, if there is any.
func (c *workloadRevisions) Update(workloadRevision *v1alpha1.WorkloadRevision) (result *v1alpha1.WorkloadRevision, err error) {
	result = &v1alpha1.WorkloadRevision{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("workloadrevisions").
		Name(workloadRevision.Name).
		Body(workloadRevision).
		Do().
		Into(result)
	return
}

// Delete takes name of the workloadRevision and deletes it. Returns an error if one occurs.
func (c *workloadRevisions) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("workloadrevisions").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *workloadRevisions) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("workloadrevisions").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched workloadRevision.
func (c *workloadRevisions) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.WorkloadRevision, err error) {
	result = &v1alpha1.WorkloadRevision{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("workloadrevisions").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
