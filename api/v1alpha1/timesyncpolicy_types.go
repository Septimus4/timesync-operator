/*
Copyright 2025.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TimeSyncPolicySpec defines the desired state of TimeSyncPolicy.
type TimeSyncPolicySpec struct {
	NamespaceSelector metav1.LabelSelector `json:"namespaceSelector"`
	Enable            bool                 `json:"enable"`
	Image             string               `json:"image"`
}

// TimeSyncPolicyStatus defines the observed state of TimeSyncPolicy.
type TimeSyncPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	MatchedNamespaces int `json:"matchedNamespaces"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// TimeSyncPolicy is the Schema for the timesyncpolicies API.
type TimeSyncPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TimeSyncPolicySpec   `json:"spec,omitempty"`
	Status TimeSyncPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TimeSyncPolicyList contains a list of TimeSyncPolicy.
type TimeSyncPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TimeSyncPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TimeSyncPolicy{}, &TimeSyncPolicyList{})
}
