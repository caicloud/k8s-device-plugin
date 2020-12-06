/*
Copyright 2019 The Caicloud Authors.

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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster,shortName="er"
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="NodeName",type=string,JSONPath=`.spec.nodeName`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// ExtendedResource describes a bare device.
type ExtendedResource struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired of ExtendedResource
	// +optional
	Spec ExtendedResourceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// Most recently observed status of the ExtendedResource.
	// +optional
	Status ExtendedResourceStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// ExtendedResourceSpec is the specification of the desired behavior of the ExtendedResource.
type ExtendedResourceSpec struct {
	// raw resource name. E.g.: nvidia.com/gpu
	RawResourceName string `json:"rawResourceName" protobuf:"bytes,1,opt,name=rawResourceName"`
	// device unique id
	DeviceID string `json:"deviceID" protobuf:"bytes,2,opt,name=deviceID"`
	// resource metadata received from device plugin.
	// e.g., gpuType: k80, zone: us-west1-b
	Properties map[string]string `json:"properties" protobuf:"bytes,3,opt,name=properties"`
	// NodeName constraints that limit what nodes this resource can be accessed from.
	// This field influences the scheduling of pods that use this resource.
	NodeName string `json:"nodeName" protobuf:"bytes,4,opt,name=nodeName"`
}

// ExtendedResourceStatus is the most recently observed status of the ExtendedResource.
type ExtendedResourceStatus struct {
	// Phase indicates if the compute resource is available or pending
	// +optional
	Phase ExtendedResourcePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase"`

	// Capacity represents the total resources of a ExtendedResource.
	// +optional
	Capacity v1.ResourceList `json:"capacity,omitempty" protobuf:"bytes,2,rep,name=capacity,casttype=ResourceList,castkey=ResourceName"`
	// Allocatable represents the ExtendedResource that are available for scheduling.
	// Defaults to Capacity.
	// +optional
	Allocatable v1.ResourceList `json:"allocatable,omitempty" protobuf:"bytes,3,rep,name=allocatable,casttype=ResourceList,castkey=ResourceName"`
	// Current Condition of extended resource.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,4,rep,name=conditions"`
}

// ExtendedResourcePhase defines ExtendedResource status.
type ExtendedResourcePhase string

const (
	// ExtendedResourceAvailable used for ExtendedResource that is already on a specific node.
	ExtendedResourceAvailable ExtendedResourcePhase = "Available"

	// ExtendedResourceBound used for ExtendedResource that is already using by pod.
	ExtendedResourceBound ExtendedResourcePhase = "Bound"

	// ExtendedResourcePending used for ExtendedResource that is not available now,
	// due to device plugin send device is unhealthy.
	ExtendedResourcePending ExtendedResourcePhase = "Pending"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExtendedResourceList is a collection of ExtendedResource.
type ExtendedResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of ExtendedResource.
	Items []ExtendedResource `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion

// ResourceClass defines how to choose the ExtendedResource.
type ResourceClass struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired of ResourceClass
	// +optional
	Spec ResourceClassSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// ResourceClassSpec is the specification of the desired behavior of the ResourceClass.
type ResourceClassSpec struct {
	// raw resource name. E.g.: nvidia.com/gpu
	RawResourceName string `json:"rawResourceName" protobuf:"bytes,1,opt,name=rawResourceName"`
	// defines general resource property matching constraints.
	// e.g.: zone in { us-west1-b, us-west1-c }; type: k80
	Requirements metav1.LabelSelector `json:"requirements" protobuf:"bytes,2,opt,name=requirements"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceClassList is a collection of ResourceClass.
type ResourceClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of ResourceClass.
	Items []ResourceClass `json:"items" protobuf:"bytes,2,rep,name=items"`
}
