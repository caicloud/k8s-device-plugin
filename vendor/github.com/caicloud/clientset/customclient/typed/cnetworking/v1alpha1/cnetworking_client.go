/*
Copyright 2020 caicloud authors. All rights reserved.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/caicloud/clientset/customclient/scheme"
	v1alpha1 "github.com/caicloud/clientset/pkg/apis/cnetworking/v1alpha1"
	rest "k8s.io/client-go/rest"
)

type CnetworkingV1alpha1Interface interface {
	RESTClient() rest.Interface
	NetworkPoliciesGetter
}

// CnetworkingV1alpha1Client is used to interact with features provided by the cnetworking.caicloud.io group.
type CnetworkingV1alpha1Client struct {
	restClient rest.Interface
}

func (c *CnetworkingV1alpha1Client) NetworkPolicies(namespace string) NetworkPolicyInterface {
	return newNetworkPolicies(c, namespace)
}

// NewForConfig creates a new CnetworkingV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*CnetworkingV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &CnetworkingV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new CnetworkingV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *CnetworkingV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new CnetworkingV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *CnetworkingV1alpha1Client {
	return &CnetworkingV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *CnetworkingV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
