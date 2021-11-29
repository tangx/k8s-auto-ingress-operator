/*
Copyright 2021.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AutoIngressSpec defines the desired state of AutoIngress
type AutoIngressSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of AutoIngress. Edit autoingress_types.go to remove/update
	// Foo string `json:"foo,omitempty"`
	RootDomain string `json:"rootDomain,omitempty"`
	// TlsConfig  string `json:"tlsConfig"`
}

// AutoIngressStatus defines the observed state of AutoIngress
type AutoIngressStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AutoIngress is the Schema for the autoingresses API
type AutoIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutoIngressSpec   `json:"spec,omitempty"`
	Status AutoIngressStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AutoIngressList contains a list of AutoIngress
type AutoIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AutoIngress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AutoIngress{}, &AutoIngressList{})
}
