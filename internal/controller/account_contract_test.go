package controller

import (
	"context"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAccountContractVersions(t *testing.T) {
	for _, version := range []string{"", "0", "-1", "invalid", "1", "2"} {
		obj := &unstructured.Unstructured{}
		obj.SetAnnotations(map[string]string{accountContractAnnotation: version})
		account, err := contractVersion(obj)
		if (err == nil) != (version == "1" || version == "2") || account != (err == nil) {
			t.Fatalf("contract=%q account=%v err=%v", version, account, err)
		}
	}
	if account, err := contractVersion(&unstructured.Unstructured{}); account || err != nil {
		t.Fatalf("absent contract account=%v err=%v", account, err)
	}
}

var (
	bindingGVK = schema.GroupVersionKind{Group: "apis.kcp.io", Version: "v1alpha2", Kind: "APIBinding"}
	schemaGVK  = schema.GroupVersionKind{Group: "apis.kcp.io", Version: "v1alpha1", Kind: "APIResourceSchema"}
	crdGVK     = schema.GroupVersionKind{Group: "apiextensions.k8s.io", Version: "v1", Kind: "CustomResourceDefinition"}
)

func boundSCPCBinding(name, schemaName, schemaUID string) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: map[string]any{"apiVersion": bindingGVK.GroupVersion().String(), "kind": bindingGVK.Kind, "metadata": map[string]any{"name": name, "uid": name + "-uid"}, "status": map[string]any{"phase": "Bound", "conditions": []any{map[string]any{"type": "InitialBindingCompleted", "status": "True"}, map[string]any{"type": "BindingUpToDate", "status": "True"}}, "boundResources": []any{map[string]any{"group": "showroom.apeirora.eu", "resource": accountContractResource, "schema": map[string]any{"name": schemaName, "UID": schemaUID, "identityHash": name + "-identity"}}}}}}
}

func boundSCPCSchema(annotations map[string]any) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: map[string]any{"apiVersion": schemaGVK.GroupVersion().String(), "kind": schemaGVK.Kind, "metadata": map[string]any{"name": "core-schema", "uid": "schema-uid", "annotations": annotations}, "spec": map[string]any{"group": "showroom.apeirora.eu", "names": map[string]any{"plural": accountContractResource}}}}
}

func kcpContractReconciler(t *testing.T, objs ...client.Object) *NamespaceReconciler {
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
	return &NamespaceReconciler{Client: fake.NewClientBuilder().WithScheme(scheme).WithRESTMapper(mapper).WithObjects(objs...).Build()}
}

func TestAccountContractUsesOnlyTheBoundSchemaUID(t *testing.T) {
	for _, matching := range []bool{true, false} {
		obj := boundSCPCSchema(map[string]any{accountContractAnnotation: "1"})
		if !matching {
			obj.SetUID("replacement")
		}
		r := kcpContractReconciler(t, boundSCPCBinding("core", "core-schema", "schema-uid"), obj)
		if account, err := r.accountContract(context.Background()); (err == nil) != matching || account != matching {
			t.Fatalf("matching UID=%v account=%v err=%v", matching, account, err)
		}
	}
}

func TestAccountContractLegacyBoundSchema(t *testing.T) {
	r := kcpContractReconciler(t, boundSCPCBinding("core", "core-schema", "schema-uid"), boundSCPCSchema(nil))
	if account, err := r.accountContract(context.Background()); account || err != nil {
		t.Fatalf("legacy schema account=%v err=%v", account, err)
	}
}

func TestAccountContractAmbiguousBoundSchema(t *testing.T) {
	r := kcpContractReconciler(t, boundSCPCBinding("core", "core-schema", "schema-uid"), boundSCPCBinding("other", "other-schema", "other-uid"), boundSCPCSchema(nil))
	if _, err := r.accountContract(context.Background()); err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("ambiguous schema err=%v", err)
	}
}
