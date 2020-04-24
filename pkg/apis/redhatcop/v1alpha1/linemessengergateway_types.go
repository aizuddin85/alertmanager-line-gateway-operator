package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LineMessengerGatewaySpec defines the desired state of LineMessengerGateway
type LineMessengerGatewaySpec struct {
	// Deployment size
	Size int32 `json:"size"`
	Image string `json:"image"`
}

// LineMessengerGatewayStatus defines the observed state of LineMessengerGateway
type LineMessengerGatewayStatus struct {
	// Pods status
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LineMessengerGateway is the Schema for the linemessengergateways API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=linemessengergateways,scope=Namespaced
type LineMessengerGateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LineMessengerGatewaySpec   `json:"spec,omitempty"`
	Status LineMessengerGatewayStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LineMessengerGatewayList contains a list of LineMessengerGateway
type LineMessengerGatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LineMessengerGateway `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LineMessengerGateway{}, &LineMessengerGatewayList{})
}
