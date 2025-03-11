package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatefulChaosExperimentSpec defines the desired state of a chaos experiment
type StatefulChaosExperimentSpec struct {
	// Target defines the selection criteria for what to target with chaos
	Target TargetSpec `json:"target"`

	// ChaosType defines the type of chaos to be injected
	// +kubebuilder:validation:Enum=DiskFailure;NetworkLatency;DatabaseConnectionDisruption;PodFailure;ResourcePressure;DataCorruption;StatefulSetScaling
	ChaosType string `json:"chaosType"`

	// Duration defines how long the chaos experiment should run
	// +kubebuilder:validation:Pattern=^([0-9]+h)?([0-9]+m)?([0-9]+s)?$
	Duration string `json:"duration"`

	// Intensity defines the severity of chaos (0.0-1.0)
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	Intensity float64 `json:"intensity"`

	// Parameters holds specific configuration options for the chaos type
	// +optional
	Parameters map[string]string `json:"parameters,omitempty"`

	// Schedule defines when and how often to run the chaos experiment
	// +optional
	Schedule *ScheduleSpec `json:"schedule,omitempty"`

	// Safety defines safety mechanisms and guardrails for the experiment
	// +optional
	Safety *SafetySpec `json:"safety,omitempty"`
}

// TargetSpec defines the target selection for chaos injection
type TargetSpec struct {
	// Selector is used to select pods based on labels
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	// Name specifies the name of a specific statefulset, deployment, or pod
	// +optional
	Name string `json:"name,omitempty"`

	// Namespace limits the scope of the experiment to a specific namespace
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// TargetType defines what type of resource to target
	// +kubebuilder:validation:Enum=StatefulSet;Deployment;Pod;PersistentVolume;PersistentVolumeClaim;Service
	// +optional
	TargetType string `json:"targetType,omitempty"`

	// Mode defines how to select targets from the filtered resources
	// (one, all, random, percentage)
	// +kubebuilder:validation:Enum=One;All;Random;Percentage;Fixed
	// +optional
	Mode string `json:"mode,omitempty"`

	// Value is used in conjunction with Mode (e.g., percentage or fixed count)
	// +optional
	Value string `json:"value,omitempty"`
}

// ScheduleSpec defines when to run chaos experiments
type ScheduleSpec struct {
	// Cron expression for scheduling experiments
	// +optional
	Cron string `json:"cron,omitempty"`

	// Immediate starts the experiment as soon as it's created
	// +optional
	Immediate bool `json:"immediate,omitempty"`

	// Once runs the experiment only once
	// +optional
	Once bool `json:"once,omitempty"`
}

// SafetySpec defines safety mechanisms for chaos experiments
type SafetySpec struct {
	// AutoRollback automatically reverses chaos when conditions are met
	// +optional
	AutoRollback bool `json:"autoRollback,omitempty"`

	// HealthChecks defines endpoints to monitor during experiments
	// +optional
	HealthChecks []HealthCheckSpec `json:"healthChecks,omitempty"`

	// PauseConditions defines when to pause experiments
	// +optional
	PauseConditions []PauseConditionSpec `json:"pauseConditions,omitempty"`

	// ResourceProtections defines resources that should never be affected
	// +optional
	ResourceProtections []ProtectionSpec `json:"resourceProtections,omitempty"`
}

// HealthCheckSpec defines a health check to monitor
type HealthCheckSpec struct {
	// Type of health check (httpGet, tcpSocket, exec)
	Type string `json:"type"`

	// For httpGet health checks
	// +optional
	Path string `json:"path,omitempty"`

	// Port to check
	// +optional
	Port int32 `json:"port,omitempty"`

	// Command for exec health checks
	// +optional
	Command []string `json:"command,omitempty"`

	// FailureThreshold defines how many consecutive failures constitute unhealthy
	// +optional
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
}

// PauseConditionSpec defines when to pause an experiment
type PauseConditionSpec struct {
	// Type of condition (metric, alert, manual)
	Type string `json:"type"`

	// MetricQuery for metric-based conditions
	// +optional
	MetricQuery string `json:"metricQuery,omitempty"`

	// Threshold for metric-based conditions
	// +optional
	Threshold string `json:"threshold,omitempty"`
}

// ProtectionSpec defines protected resources
type ProtectionSpec struct {
	// Resource type to protect
	// +kubebuilder:validation:Enum=Namespace;Label;Annotation;Name
	Type string `json:"type"`

	// Value to match for protection
	Value string `json:"value"`
}

// StatefulChaosExperimentStatus defines the observed state of a chaos experiment
type StatefulChaosExperimentStatus struct {
	// Phase of the chaos experiment (Pending, Running, Completed, Failed)
	Phase string `json:"phase"`

	// StartTime when the experiment began
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// EndTime when the experiment finished
	// +optional
	EndTime *metav1.Time `json:"endTime,omitempty"`

	// Conditions represents the latest available observations of current state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// TargetResources lists the actual resources affected
	// +optional
	TargetResources []TargetResourceStatus `json:"targetResources,omitempty"`

	// FailureReason provides more information about a failure
	// +optional
	FailureReason string `json:"failureReason,omitempty"`
}

// TargetResourceStatus describes a resource affected by chaos
type TargetResourceStatus struct {
	// Kind of the target resource
	Kind string `json:"kind"`

	// Name of the target resource
	Name string `json:"name"`

	// Namespace of the target resource
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// UID of the target resource
	// +optional
	UID string `json:"uid,omitempty"`

	// Status of chaos injection for this target
	// +optional
	Status string `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.chaosType",description="Type of chaos"
// +kubebuilder:printcolumn:name="Target",type="string",JSONPath=".spec.target.targetType",description="Target type"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="Experiment phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// StatefulChaosExperiment is the Schema for the chaos experiment API
type StatefulChaosExperiment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StatefulChaosExperimentSpec   `json:"spec,omitempty"`
	Status StatefulChaosExperimentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StatefulChaosExperimentList contains a list of StatefulChaosExperiment
type StatefulChaosExperimentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StatefulChaosExperiment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StatefulChaosExperiment{}, &StatefulChaosExperimentList{})
}
