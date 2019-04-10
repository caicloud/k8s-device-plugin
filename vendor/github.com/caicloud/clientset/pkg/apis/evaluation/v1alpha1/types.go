/*
Copyright 2019 caicloud authors. All rights reserved.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EvalJob is a specification for a EvalJob resource
type EvalJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the EvalJob.
	Spec EvalJobSpec `json:"spec"`

	// Most recently observed status of the EvalJob.
	Status EvalJobStatus `json:"status"`
}

// EvalJobSpec is the spec for a EvalJob resource
// default cleanPodPolicy = all
type EvalJobSpec struct {
	// CleanPodPolicy defines the policy to kill pods after EvalJob is
	// succeeded.
	// Default to None.
	CleanPodPolicy *CleanPodPolicy `json:"cleanPodPolicy,omitempty"`

	// PersistentVolumeClaim is shared by all Pods in the same evaluation.
	// The user only needs to specify the size of the PVC.
	StorageSize string `json:"storageSize"`

	// Models contains all models that will be evaluated.
	Models []EvalModel `json:"models,omitempty"`

	// Functions contains all evaluation functions.
	Functions []EvalFunc `json:"functions,omitempty"`

	// Template specified the pod template.
	// This v1alpha1 version, we only care about 'labels','annotations','env' and 'resources'.
	// We will overwrite the image, command, args when really creating pod.
	Template corev1.PodTemplateSpec `json:"template,omitempty"`
}

// EvalModel is the descriptor of a model
type EvalModel struct {
	// Name of model.
	Name string `json:"name"`
	// Version of model.
	Version string `json:"version"`
	// ServingImage is a formatted string, like "serving:v0.1"，
	// specified the image serving will use.
	ServingImage string `json:"servingImage"`
	// Framework of the model.
	Framework string `json:"framework"`
}

// EvalFunc is the descriptor of an evaluation function
type EvalFunc struct {
	// Name of evaluation function.
	Name string `json:"name"`
	// URL where to find the function code.
	URL string `json:"url"`
	// Protocol to get the function code.
	Protocol string `json:"protocol"`
}

// +k8s:deepcopy-gen=true

// EvalJobStatus represents the current observed state of the EvalJob.
type EvalJobStatus struct {
	// The phase of EvalJob.
	Phase JobPhaseType `json:"phase,omitempty"`
	// The status of each worker pod.
	WorkersStatus EvalWorkersStatus `json:"workersStatus,omitempty"`
	// Represents the lastest available observations of a EvalJob's current state.
	Conditions []EvalJobCondition `json:"conditions,omitempty"`
	// The time that worker pods are created.
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// The end time of this evaljob finished.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
}

type EvalJobConditionType string

// These are valid conditions of an EvalJob
const (
	// Processing means the evaljob has been created correctly
	// and it is processing tasks in the right way.
	JobProcessing EvalJobConditionType = "Processing"
	// OutOfStorage means the size of pvc is smaller than the job need
	// or there is no available pvc.
	JobOutOfStorage EvalJobConditionType = "OutOfStorage"
)

// +k8s:deepcopy-gen=true
// EvalJobCondition describe the state of a evaljob at a certain point.
type EvalJobCondition struct {
	// Type of job condition.
	Type EvalJobConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

// EvalWorkersStatus represents the current observed state of all worker pods.
type EvalWorkersStatus struct {
	// The total number of worker pods
	TotalWorkersNum int32 `json:"totalWorkersNum"`

	// Workers contains a map of all workers.
	// key: modelName/modleVersion
	// value: Pod.Status.Phase
	Workers map[string]corev1.PodPhase `json:"workers"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// EvalJobList is a list of EvlJob resources
type EvalJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []EvalJob `json:"items"`
}

type JobPhaseType string

// These are valid phases of an EvalJob.
const (
	// JobPending means the job is creating pvc.
	JobPending JobPhaseType = "Pending"

	// JobLauncing means launcher pod has been created.
	JobLaunching JobPhaseType = "Launching"

	// JobLaunced means launcher pod has been completed.
	JobLaunched JobPhaseType = "Launched"

	// JobRunning means one or more worker pods had been running.
	JobRunning JobPhaseType = "Running"

	// JobSucceeded means all pods of this job
	// had been completed.
	JobSucceeded JobPhaseType = "Succeeded"

	// JobFailed means launcher pod has been failed.
	JobFailed JobPhaseType = "Failed"
)

// CleanPodPolicy describes how to deal with pods when the EvalJob is finished.
type CleanPodPolicy string

const (
	// User hasn't defined the cleanPodPolicy, do nothing.
	CleanPodPolicyUndefined CleanPodPolicy = ""
	// CleanPodPolcicyAll means the controller will delete all pods after it finished.
	CleanPodPolicyAll CleanPodPolicy = "All"
	// CleanPodPolicyRunning means the controller only delete pods that are still running after it finished.
	CleanPodPolicyRunning CleanPodPolicy = "Running"
)
