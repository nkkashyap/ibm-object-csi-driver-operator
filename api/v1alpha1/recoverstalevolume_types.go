/*
Copyright 2024.

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

// Package v1alpha1 ...
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RecoverStaleVolumeSpec defines the desired state of RecoverStaleVolume
type RecoverStaleVolumeSpec struct {
	LogHistory int64                      `json:"logHistory,omitempty"`
	Data       []NamespacedDeploymentData `json:"data"`
}

// NamespacedDeploymentData ...
type NamespacedDeploymentData struct {
	Namespace   string   `json:"namespace"`
	Deployments []string `json:"deployments,omitempty"`
}

// RecoverStaleVolumeStatus defines the observed state of RecoverStaleVolume
type RecoverStaleVolumeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RecoverStaleVolume is the Schema for the recoverstalevolumes API
type RecoverStaleVolume struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RecoverStaleVolumeSpec   `json:"spec,omitempty"`
	Status RecoverStaleVolumeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RecoverStaleVolumeList contains a list of RecoverStaleVolume
type RecoverStaleVolumeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RecoverStaleVolume `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RecoverStaleVolume{}, &RecoverStaleVolumeList{})
}
