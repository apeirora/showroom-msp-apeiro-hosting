package controller

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestAccountSCPCRendering(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		for _, f := range flavors {
			scpc, err := renderSCPC(f, "account", enabled)
			if err != nil {
				t.Fatal(err)
			}
			account, _, _ := unstructured.NestedBool(scpc.Object, "spec", "target", "kcpAccount")
			want := enabled && f.scpcName == ironcoreSCPCName
			if account != want {
				t.Fatalf("%s enabled=%v account=%v", f.scpcName, enabled, account)
			}
			ns, _, _ := unstructured.NestedString(scpc.Object, "spec", "target", "namespace")
			if want {
				if ns != "" {
					t.Fatal("account target has namespace")
				}
				secretNS, _, _ := unstructured.NestedString(scpc.Object, "spec", "kubeconfigSource", "secretRef", "namespace")
				if secretNS != "default" {
					t.Fatal("account secret is not in default")
				}
				resources, _, _ := unstructured.NestedSlice(scpc.Object, "spec", "resources")
				if resources[0].(map[string]any)["apiVersion"] != "apeiro.eu/v1alpha1" {
					t.Fatal("account child has wrong group")
				}
			} else if ns != "{{ .SC.namespace }}" {
				t.Fatal("local target changed")
			}
			if _, found, _ := unstructured.NestedFieldNoCopy(scpc.Object, "spec", "target", "kubeconfigSecretRef"); found {
				t.Fatal("unexpected target kubeconfig")
			}
		}
	}
}
