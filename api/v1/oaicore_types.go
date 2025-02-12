package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OAIClusterSpec defines the desired state of OAICluster
type OAIClusterSpec struct {
	AMF AMFSpec `json:"amf"`
	SMF SMFSpec `json:"smf"`
	UPF UPFSpec `json:"upf"`
	NRF NRFSpec `json:"nrf"`
}

type AMFSpec struct {
	ReplicaCount int32 `json:"replicaCount"`
}

type SMFSpec struct {
	ReplicaCount int32 `json:"replicaCount"`
}

type UPFSpec struct {
	ReplicaCount int32 `json:"replicaCount"`
}

type NRFSpec struct {
	ReplicaCount int32 `json:"replicaCount"`
}

// OAIClusterStatus defines the observed state of OAICluster
type OAIClusterStatus struct {
	AMFReady bool `json:"amfReady"`
	SMFReady bool `json:"smfReady"`
	UPFReady bool `json:"upfReady"`
	NRFReady bool `json:"nrfReady"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// OAICluster is the Schema for the oaiclusters API
type OAICluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OAIClusterSpec   `json:"spec,omitempty"`
	Status OAIClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OAIClusterList contains a list of OAICluster
type OAIClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OAICluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OAICluster{}, &OAIClusterList{})
}
