// Copyright 2026 ApeiroRA contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestNamespaceIronCoreAccountTarget(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		name := "local"
		if enabled {
			name = "account"
		}
		t.Run(name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			if err := corev1.AddToScheme(scheme); err != nil {
				t.Fatal(err)
			}
			scheme.AddKnownTypeWithName(scpcGVK, &unstructured.Unstructured{})
			manifest, err := ConfigureIronCoreManifest(gardenDevConfig(t))
			if err != nil {
				t.Fatal(err)
			}
			ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "account-a", UID: "namespace-uid", Annotations: map[string]string{
				bindExportsAnnotation: DefaultSimpleClusterExportName + "," + DefaultSapExportName,
			}}}
			c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(ns).Build()
			r := &NamespaceReconciler{Client: c, Scheme: scheme, IronCoreManifest: manifest,
				IronCoreKCPAccount: enabled, SimpleClusterExportName: DefaultSimpleClusterExportName, SapExportName: DefaultSapExportName}
			request := ctrl.Request{NamespacedName: types.NamespacedName{Name: ns.Name}}
			for iteration := 0; iteration < 2; iteration++ {
				if _, err := r.Reconcile(context.Background(), request); err != nil {
					t.Fatal(err)
				}
				for _, f := range flavors {
					scpc := &unstructured.Unstructured{}
					scpc.SetGroupVersionKind(scpcGVK)
					if err := c.Get(context.Background(), types.NamespacedName{Namespace: ns.Name, Name: f.scpcName}, scpc); err != nil {
						t.Fatal(err)
					}
					accountMode := enabled && f.scpcName == ironcoreSCPCName
					target, _, _ := unstructured.NestedBool(scpc.Object, "spec", "target", "kcpAccount")
					if target != accountMode {
						t.Fatalf("%s account target = %v, want %v", f.scpcName, target, accountMode)
					}
					wantNamespace, wantSecretNamespace := "{{ .SC.namespace }}", ""
					if accountMode {
						wantNamespace, wantSecretNamespace = "", "default"
					}
					namespace, _, _ := unstructured.NestedString(scpc.Object, "spec", "target", "namespace")
					secretNamespace, _, _ := unstructured.NestedString(scpc.Object, "spec", "kubeconfigSource", "secretRef", "namespace")
					if namespace != wantNamespace || secretNamespace != wantSecretNamespace {
						t.Fatalf("%s namespace = %q, Secret namespace = %q", f.scpcName, namespace, secretNamespace)
					}
					if _, exists, _ := unstructured.NestedMap(scpc.Object, "spec", "target", "kubeconfigSecretRef"); exists {
						t.Fatal("target must reuse the operator provider credential")
					}
					owners := scpc.GetOwnerReferences()
					if len(owners) != 1 || owners[0].UID != ns.UID {
						t.Fatal("namespace must retain SCPC ownership")
					}
				}
			}
		})
	}
}
