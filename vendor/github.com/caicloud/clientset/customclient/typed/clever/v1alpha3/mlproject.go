/*
Copyright 2020 caicloud authors. All rights reserved.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha3

import (
	"time"

	scheme "github.com/caicloud/clientset/customclient/scheme"
	v1alpha3 "github.com/caicloud/clientset/pkg/apis/clever/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MLProjectsGetter has a method to return a MLProjectInterface.
// A group's client should implement this interface.
type MLProjectsGetter interface {
	MLProjects(namespace string) MLProjectInterface
}

// MLProjectInterface has methods to work with MLProject resources.
type MLProjectInterface interface {
	Create(*v1alpha3.MLProject) (*v1alpha3.MLProject, error)
	Update(*v1alpha3.MLProject) (*v1alpha3.MLProject, error)
	UpdateStatus(*v1alpha3.MLProject) (*v1alpha3.MLProject, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha3.MLProject, error)
	List(opts v1.ListOptions) (*v1alpha3.MLProjectList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha3.MLProject, err error)
	MLProjectExpansion
}

// mLProjects implements MLProjectInterface
type mLProjects struct {
	client rest.Interface
	ns     string
}

// newMLProjects returns a MLProjects
func newMLProjects(c *CleverV1alpha3Client, namespace string) *mLProjects {
	return &mLProjects{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the mLProject, and returns the corresponding mLProject object, and an error if there is any.
func (c *mLProjects) Get(name string, options v1.GetOptions) (result *v1alpha3.MLProject, err error) {
	result = &v1alpha3.MLProject{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mlprojects").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MLProjects that match those selectors.
func (c *mLProjects) List(opts v1.ListOptions) (result *v1alpha3.MLProjectList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha3.MLProjectList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mlprojects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested mLProjects.
func (c *mLProjects) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("mlprojects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a mLProject and creates it.  Returns the server's representation of the mLProject, and an error, if there is any.
func (c *mLProjects) Create(mLProject *v1alpha3.MLProject) (result *v1alpha3.MLProject, err error) {
	result = &v1alpha3.MLProject{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("mlprojects").
		Body(mLProject).
		Do().
		Into(result)
	return
}

// Update takes the representation of a mLProject and updates it. Returns the server's representation of the mLProject, and an error, if there is any.
func (c *mLProjects) Update(mLProject *v1alpha3.MLProject) (result *v1alpha3.MLProject, err error) {
	result = &v1alpha3.MLProject{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mlprojects").
		Name(mLProject.Name).
		Body(mLProject).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *mLProjects) UpdateStatus(mLProject *v1alpha3.MLProject) (result *v1alpha3.MLProject, err error) {
	result = &v1alpha3.MLProject{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mlprojects").
		Name(mLProject.Name).
		SubResource("status").
		Body(mLProject).
		Do().
		Into(result)
	return
}

// Delete takes name of the mLProject and deletes it. Returns an error if one occurs.
func (c *mLProjects) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mlprojects").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *mLProjects) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mlprojects").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched mLProject.
func (c *mLProjects) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha3.MLProject, err error) {
	result = &v1alpha3.MLProject{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("mlprojects").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
