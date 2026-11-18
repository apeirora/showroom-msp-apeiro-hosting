/*
Copyright 2026.

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

// Package v1alpha1 contains API Schema definitions for the apeiro.eu v1alpha1 API group.
// +kubebuilder:object:generate=true
// +groupName=apeiro.eu
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// GroupVersion is group version used to register these objects.
	GroupVersion = schema.GroupVersion{Group: "apeiro.eu", Version: "v1alpha1"}

	// SchemeBuilder collects the registration functions for this group-version.
	// It uses the apimachinery builder because controller-runtime's
	// scheme.Builder is deprecated (SA1019).
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// addKnownTypes registers the API types of this group-version with the scheme.
func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(GroupVersion,
		&OpenStackShoot{}, &OpenStackShootList{},
		&IronCoreShoot{}, &IronCoreShootList{},
	)
	metav1.AddToGroupVersion(s, GroupVersion)
	return nil
}
