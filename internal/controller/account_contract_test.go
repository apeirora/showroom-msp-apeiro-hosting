package controller

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAccountContractVersions(t *testing.T) {
	for _, version := range []string{"", "0", "-1", "invalid", "1", "2"} {
		obj := &unstructured.Unstructured{}
		obj.SetAnnotations(map[string]string{accountContractAnnotation: version})
		err := contractVersion(obj)
		if (err == nil) != (version == "1" || version == "2") {
			t.Fatalf("contract=%q err=%v", version, err)
		}
	}
}

func TestAccountContractUsesOnlyTheBoundSchemaUID(t *testing.T) {
	bindingGVK := schema.GroupVersionKind{Group: "apis.kcp.io", Version: "v1alpha2", Kind: "APIBinding"}
	schemaGVK := schema.GroupVersionKind{Group: "apis.kcp.io", Version: "v1alpha1", Kind: "APIResourceSchema"}
	crdGVK := schema.GroupVersionKind{Group: "apiextensions.k8s.io", Version: "v1", Kind: "CustomResourceDefinition"}
	for _, matching := range []bool{true, false} {
		binding := &unstructured.Unstructured{Object: map[string]any{"apiVersion": bindingGVK.GroupVersion().String(), "kind": bindingGVK.Kind, "metadata": map[string]any{"name": "core", "uid": "binding-uid"}, "status": map[string]any{"phase": "Bound", "conditions": []any{map[string]any{"type": "InitialBindingCompleted", "status": "True"}, map[string]any{"type": "BindingUpToDate", "status": "True"}}, "boundResources": []any{map[string]any{"group": "showroom.apeirora.eu", "resource": accountContractResource, "schema": map[string]any{"name": "core-schema", "UID": "schema-uid", "identityHash": "core-identity"}}}}}}
		obj := &unstructured.Unstructured{Object: map[string]any{"apiVersion": schemaGVK.GroupVersion().String(), "kind": schemaGVK.Kind, "metadata": map[string]any{"name": "core-schema", "uid": "schema-uid", "annotations": map[string]any{accountContractAnnotation: "1"}}, "spec": map[string]any{"group": "showroom.apeirora.eu", "names": map[string]any{"plural": accountContractResource}}}}
		if !matching {
			obj.SetUID("replacement")
		}
		scheme := runtime.NewScheme()
		if err := corev1.AddToScheme(scheme); err != nil {
			t.Fatal(err)
		}
		scheme.AddKnownTypeWithName(bindingGVK, &unstructured.Unstructured{})
		scheme.AddKnownTypeWithName(bindingGVK.GroupVersion().WithKind("APIBindingList"), &unstructured.UnstructuredList{})
		scheme.AddKnownTypeWithName(schemaGVK, &unstructured.Unstructured{})
		scheme.AddKnownTypeWithName(crdGVK, &unstructured.Unstructured{})
		mapper := apimeta.NewDefaultRESTMapper([]schema.GroupVersion{bindingGVK.GroupVersion(), schemaGVK.GroupVersion(), crdGVK.GroupVersion()})
		for _, gvk := range []schema.GroupVersionKind{bindingGVK, schemaGVK, crdGVK} {
			mapper.Add(gvk, apimeta.RESTScopeRoot)
		}
		c := fake.NewClientBuilder().WithScheme(scheme).WithRESTMapper(mapper).WithObjects(binding, obj).Build()
		r := &NamespaceReconciler{Client: c}
		if err := r.accountContract(context.Background()); (err == nil) != matching {
			t.Fatalf("matching UID=%v err=%v", matching, err)
		}
	}
}
