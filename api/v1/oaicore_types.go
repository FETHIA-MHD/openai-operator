package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OaiCoreSpec définit la configuration souhaitée pour OaiCore
type OaiCoreSpec struct {
	AMF AMFSpec `json:"amf,omitempty"`
	SMF SMFSpec `json:"smf,omitempty"`
	UPF UPFSpec `json:"upf,omitempty"`
	NRF NRFSpec `json:"nrf,omitempty"`
}

// AMFSpec définit la configuration pour l'AMF
type AMFSpec struct {
	Replicas int32  `json:"replicas"`
	Image    string `json:"image"`
}

// SMFSpec définit la configuration pour le SMF
type SMFSpec struct {
	Replicas int32  `json:"replicas"`
	Image    string `json:"image"`
}

// UPFSpec définit la configuration pour l'UPF
type UPFSpec struct {
	Replicas int32  `json:"replicas"`
	Image    string `json:"image"`
}

// NRFSpec définit la configuration pour le NRF
type NRFSpec struct {
	Replicas int32  `json:"replicas"`
	Image    string `json:"image"`
}

// OaiCoreStatus définit l'état observé de OaiCore
type OaiCoreStatus struct {
	AMFReady bool `json:"amfReady,omitempty"`
	SMFReady bool `json:"smfReady,omitempty"`
	UPFReady bool `json:"upfReady,omitempty"`
	NRFReady bool `json:"nrfReady,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// OaiCore est la définition principale de la ressource
type OaiCore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OaiCoreSpec   `json:"spec,omitempty"`
	Status OaiCoreStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OaiCoreList contient une liste de ressources OaiCore
type OaiCoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OaiCore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OaiCore{}, &OaiCoreList{})
}
