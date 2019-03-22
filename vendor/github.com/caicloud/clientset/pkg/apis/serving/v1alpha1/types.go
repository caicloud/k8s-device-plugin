package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ServingPlural = "servings"
)

// Serving defines a serving deployment.
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=serving
// +kubebuilder:subresource:status
type Serving struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ServingSpec   `json:"spec,omitempty"`
	Status            ServingStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=servings

// ServingList describes an array of Serving instances
type ServingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of Servings
	Items []Serving `json:"items"`
}

// ServingType is the type of serving jobs.
type ServingType string

const (
	// TensorRT is the type to serve with TensorRT. Not Implemented.
	TensorRT ServingType = "TensorRT"
	// TFServing is the type to serve with TFServing. Not Implemented.
	TFServing ServingType = "TFServing"
	// MXNetServing is the type to serve with MXNetServing. Not Implemented.
	MXNetServing ServingType = "MXNetServing"
	// GPUSharing is essentially implemented using TensorRT.
	GPUSharing ServingType = "GPUSharing"
	// GraphPipe is the type to serve with wrapped GraphPipe.
	GraphPipe ServingType = "GraphPipe"
	// Custom is the type to serve with Customized Images.
	Custom ServingType = "Custom"
)

// ServingSpec defines the specification of serving deployment.
type ServingSpec struct {
	// Resource requirements for serving instance.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Number of replicas for a serving instance. This is fixed to 1 for GPUSharing serving type.
	Replicas int32 `json:"replicas"`
	// Scaling is the configuration about how to scale the serving service.
	// If the scalingConfig is defined, replicas will not work.
	Scaling *ScalingConfig `json:"scalingConfig,omitempty"`
	// UserInfo is the information about the user.
	UserInfo *UserInformation `json:"userInfo,omitempty"`
	// StorageConfig is the configuration about the storage.
	Storage *StorageConfig `json:"storageConfig,omitempty"`
	// Models is the list of models to be served via the serving deployment.
	// In Scene, the size of the slices will be 1, while it can be larger than 1
	// in Serving.
	Models []ServingModel `json:"models,omitempty"`
	// Type of the Serving
	Type ServingType `json:"type,omitempty"`
}

// StorageConfig is the type for the configuration about storage.
type StorageConfig struct {
	// PersistentVolumeClaim is shared by all Pods in the same Serving.
	// The PVC must be ReadWriteMany in order to be used for multiple serving instances.
	// The user only needs to specify the storage size of the PVC.
	Size string `json:"size,omitempty"`
	// ClassName is the storageclass name.
	ClassName string `json:"classname,omitempty"`
}

// UserInformation is the type to store the user-related information.
// The information will be used to compose the PodSpec in the deployment.
type UserInformation struct {
	// Group defines the group that the user belongs to.
	Group string `json:"group"`
	// Username is the user's ID.
	Username string `json:"username"`
}

// ScalingConfig defines the configuration about how to scale the serving service.
type ScalingConfig struct {
	MinReplicas     *int32           `json:"minReplicas,omitempty"`
	MaxReplicas     int32            `json:"maxReplicas"`
	ResourceMetrics []ResourceMetric `json:"resourceMetric,omitempty"`
	CustomMetrics   []CustomMetric   `json:"customMetric,omitempty"`
}

// ResourceMetric specifies how to scale based on a single metric.
type ResourceMetric struct {
	Name corev1.ResourceName `json:"name,omitempty"`
	// At least one fields below should be set.
	Value              *resource.Quantity `json:"value,omitempty"`
	AverageValue       *resource.Quantity `json:"averageValue,omitempty"`
	AverageUtilization *int32             `json:"averageUtilization,omitempty"`
}

// CustomMetric defines customized resource.
type CustomMetric struct {
	// TBD
}

// +k8s:deepcopy-gen=true

// ServingStatus defines the status of serving deployment.
type ServingStatus struct {
	// InstanceStatus is the aggregated status of serving deployment.
	InstanceStatus []ServingInstanceStatus `json:"instanceStatus,omitempty"`
	// VolumeName is the name of the Volume.
	VolumeName string `json:"volumeName,omitempty"`
	// Conditions is an array of current observed job conditions.
	Conditions []ServingCondition `json:"conditions,omitempty"`
	// Represents time when the job was acknowledged by the job controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the job was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Represents last time when the job was reconciled. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`
}

// +k8s:deepcopy-gen=true

// ServingCondition describes the state of the serving service at a certain point.
type ServingCondition struct {
	// Type of job condition.
	Type ServingConditionType `json:"type"`
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

// ServingConditionType is the type of ServingCondition.
type ServingConditionType string

const (
	ServingRunning ServingConditionType = "Running"
	ServingHealth  ServingConditionType = "ModelHealth"
)

// ServingInstanceStatus defines status for a single serving instance.
type ServingInstanceStatus struct {
	// Phase of the serving instance. This is simply the phase of the corresponding pod.
	Phase corev1.PodPhase `json:"phase"`
	// Statuses of the models running in the serving instance. Model here means 'model + version'.
	ModelStatuses []ModelStatus `json:"modelStatuses"`
}

// +k8s:deepcopy-gen=true

// ServingModel is a model [with version] served with a serving instance.
// A model is uniquely identified via name and version.
type ServingModel struct {
	// Name and version, version is optional
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`

	// Time when the model is added to the serving instance.
	AdditionTime *metav1.Time `json:"additionTime,omitempty"`

	// Servering framework specific configurations
	TensorRTModelConfig *TensorRTModelConfig `json:"tensorRTConfig,omitempty"`
}

// TensorRTModelConfig is the configuration related to TensorRT.
type TensorRTModelConfig struct {
	// Number of instance per model.
	InstanceNum int32 `json:"instanceNum,omitempty"`
	// URLAlias in /api/infer/<URLAlias>/[Version(=1 in 1.4.0)]
	// URLAlias is not used in Scene Serving but TensorRT Inference Serve
	URLAlias string `json:"urlAlias,omitempty"`
}

// ModelStatus is the status of a model.
type ModelStatus struct {
	// Name and version of corresponding model.
	Name    string
	Version string

	// Health information of the model.
	Health ServingModelHealth
}

// ServingModelHealth is the health condition of a model (under a specific serving instance).
type ServingModelHealth string

const (
	ServingModelHealthy   ServingModelHealth = "Healthy"
	ServingModelUnhealthy ServingModelHealth = "Unhealthy"
)

// -----------------------------------------------------------------

const (
	ScenePlural = "scenes"
)

// Scene defines a scene service
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
type Scene struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SceneSpec   `json:"spec,omitempty"`
	Status            SceneStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SceneList describes an array of Scenes
type SceneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of Scenes
	Items []Scene `json:"items"`
}

type SceneStatus struct {
	// Total number of servings.
	// +optional
	Servings int32 `json:"servings,omitempty"`

	// Total number of ready servings.
	// +optional
	ReadyServings int32 `json:"readyServings,omitempty"`

	// Total number of unavailable servings.
	// +optional
	UnavailableServings int32 `json:"unavailableServings,omitempty"`

	// Conditions is an array of current observed job conditions.
	Conditions []SceneCondition `json:"conditions,omitempty"`

	// Represents time when the job was acknowledged by the job controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the job was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Represents last time when the job was reconciled. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`
}

type SceneCondition struct {
	// Type of job condition.
	Type SceneConditionType `json:"type"`
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

type SceneConditionType string

const (
	SceneRunning SceneConditionType = "Running"
	SceneHealth  SceneConditionType = "SceneHealth"
)

// ServingTemplateSpec describes a template to create a serving deployment.
type ServingTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ServingSpec `json:"spec,omitempty"`
}

type SceneSpec struct {
	// Resource quota for a scene. Each scene occupies a namespace, thus
	// this is the quota of a namespace.
	Quota corev1.ResourceList `json:"quota"`
	// A list of serving deployments under a scene.
	Servings []ServingTemplateSpec `json:"servings"`
	// A list of route configuration of the scene.
	Http []*HTTPRoute `json:"http"`
	// Name of default serving deployment.
	DefaultServing string
}

type HTTPRoute struct {
	// A list of rules to match requests. All matches are ORed.
	Match []*HTTPMatchRequest `json:"match,omitempty"`
	// A list of route information for matched requests.
	Route []*HTTPRouteServing `json:"route,omitempty"`
}

// HTTPMatchRequest specify rules to match requests. All rules are ANDed.
type HTTPMatchRequest struct {
	// Match headers of a request.
	Headers map[string]*StringMatch
}

type HTTPRouteServing struct {
	// Name of serving defined in []SceneSpec.Servings.
	Serving string
	// Traffic weight of the serving.
	Weight int32
}

// StringMatch defines 3 different types of matching strategy, i.e. only match prefix,
// exact string match, and regular expression match.
type StringMatch struct {
	prefix string
	exact  string
	regex  string
}
