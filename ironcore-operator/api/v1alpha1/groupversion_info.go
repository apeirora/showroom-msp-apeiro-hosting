/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
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
	GroupVersion  = schema.GroupVersion{Group: "apeiro.eu", Version: "v1alpha1"}
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// addKnownTypes registers the API types of this group-version with the scheme.
// It uses the apimachinery builder because controller-runtime's
// scheme.Builder is deprecated (SA1019).
func addKnownTypes(s *runtime.Scheme) error {
	s.AddKnownTypes(GroupVersion,
		&IronCoreVM{}, &IronCoreVMList{},
	)
	metav1.AddToGroupVersion(s, GroupVersion)
	return nil
}
