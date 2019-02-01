package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Annotation constants
const (
	// ClusterIDLabel is the label that a machineset must have to identify the
	// cluster to which it belongs.
	ClusterIDLabel   = "sigs.k8s.io/cluster-api-cluster"
	MachineRoleLabel = "sigs.k8s.io/cluster-api-machine-role"
	MachineTypeLabel = "sigs.k8s.io/cluster-api-machine-type"
)

// HyperkitMachineProviderConfig is the type that will be embedded in a Machine.Spec.ProviderSpec field
// for an Hyperkit instance.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HyperkitMachineProviderConfig struct {
	metav1.TypeMeta `json:",inline"`

	DomainMemory             int        `json:"domainMemory"`
	DomainVcpu               int        `json:"domainVcpu"`
	IgnKey                   string     `json:"ignKey"`
	Ignition                 *Ignition  `json:"ignition"`
	NetworkInterfaceName     string     `json:"networkInterfaceName"`
	NetworkInterfaceHostname string     `json:"networkInterfaceHostname"`
	NetworkInterfaceAddress  string     `json:"networkInterfaceAddress"`
	NetworkUUID              string     `json:"networkUUID"`
	Autostart                bool       `json:"autostart"`
	HyperkitBinaryPath       string     `json:"hyperkitBinaryPath"`
	HyperkitStateDir         string     `json:"hyperkitStateDir"`
}

// Ignition contains location of ignition to be run during bootstrapping
type Ignition struct {
	// Ignition config to be run during bootstrapping
	UserDataSecret string `json:"userDataSecret"`
}

// HyperkitClusterProviderConfig is the type that will be embedded in a Cluster.Spec.ProviderSpec field.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HyperkitClusterProviderConfig struct {
	metav1.TypeMeta `json:",inline"`
}

// HyperkitMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains Hyperkit-specific status information.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HyperkitMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// InstanceID is the instance ID of the machine created in Hyperkit
	InstanceID *string `json:"instanceID"`

	// InstanceState is the state of the HyperKit instance for this machine
	InstanceState *string `json:"instanceState"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []HyperkitMachineProviderCondition `json:"conditions"`
}

// HyperkitMachineProviderConditionType is a valid value for HyperkitMachineProviderCondition.Type
type HyperkitMachineProviderConditionType string

// Valid conditions for an Hyperkit machine instance
const (
	// MachineCreated indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreated HyperkitMachineProviderConditionType = "MachineCreated"
)

// HyperkitMachineProviderCondition is a condition in a HyperkitMachineProviderStatus
type HyperkitMachineProviderCondition struct {
	// Type is the type of the condition.
	Type HyperkitMachineProviderConditionType `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message"`
}

// HyperkitClusterProviderStatus is the type that will be embedded in a Cluster.Status.ProviderStatus field.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HyperkitClusterProviderStatus struct {
	metav1.TypeMeta `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// HyperkitMachineProviderConfigList contains a list of HyperkitMachineProviderConfig
type HyperkitMachineProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HyperkitMachineProviderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HyperkitMachineProviderConfig{}, &HyperkitMachineProviderConfigList{}, &HyperkitMachineProviderStatus{})
}
