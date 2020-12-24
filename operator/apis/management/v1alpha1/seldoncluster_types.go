/*
Copyright 2020.

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

// SeldonCoreSpec defines the desired state of the Seldon Core installation
type SeldonCoreSpec struct {
	Version string `json:"version,omitempty"`
}

// SeldonDeploySpec defines the desired state of the Seldon Deploy installation
type SeldonDeploySpec struct {
	Version string `json:"version,omitempty"`
}

// SeldonClusterSpec defines the desired state of SeldonCluster
type SeldonClusterSpec struct {
	SeldonCore *SeldonCoreSpec `json:"seldonCore,omitempty"`
}

type StatusState string

// CRD Status values
const (
	StatusStateAvailable StatusState = "Available"
	StatusStateCreating  StatusState = "Creating"
	StatusStateFailed    StatusState = "Failed"
)

// SeldonClusterStatus defines the observed state of SeldonCluster
type SeldonClusterStatus struct {
    State StatusState `json:"state,omitempty"`
    Description string `json:"description,omitempty"`
}

// SeldonCluster is the Schema for the seldonclusters API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type SeldonCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SeldonClusterSpec   `json:"spec,omitempty"`
	Status SeldonClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SeldonClusterList contains a list of SeldonCluster
type SeldonClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SeldonCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SeldonCluster{}, &SeldonClusterList{})
}
