package controller

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestAccountSCPCRendering(t *testing.T) {
	for _, f := range flavors {
		scpc, err := renderSCPC(f, "account", true)
		if err != nil {
			t.Fatal(err)
		}
		for _, field := range []string{"namespace", "kubeconfigSecretRef", "kcpAccount"} {
			if _, found, _ := unstructured.NestedFieldNoCopy(scpc.Object, "spec", "target", field); found {
				t.Fatalf("%s has destination selector %s", f.scpcName, field)
			}
		}
		secretNS, _, _ := unstructured.NestedString(scpc.Object, "spec", "kubeconfigSource", "secretRef", "namespace")
		if secretNS != "default" {
			t.Fatal("exported kubeconfig Secret is not in the account default namespace")
		}
		resources, _, _ := unstructured.NestedSlice(scpc.Object, "spec", "resources")
		if resources[0].(map[string]any)["apiVersion"] != "apeiro.eu/v1alpha1" {
			t.Fatal("wrong child group")
		}
	}
}

func TestLegacySCPCRendering(t *testing.T) {
	for _, f := range flavors {
		scpc, err := renderSCPC(f, "account", false)
		if err != nil {
			t.Fatal(err)
		}
		target, _, _ := unstructured.NestedString(scpc.Object, "spec", "target", "namespace")
		if target != "{{ .SC.namespace }}" {
			t.Fatalf("%s target namespace %q", f.scpcName, target)
		}
		if _, found, _ := unstructured.NestedFieldNoCopy(scpc.Object, "spec", "target", "kubeconfigSecretRef"); found {
			t.Fatalf("%s keeps the garden kubeconfig ref", f.scpcName)
		}
		if _, found, _ := unstructured.NestedFieldNoCopy(scpc.Object, "spec", "kubeconfigSource", "secretRef", "namespace"); found {
			t.Fatalf("%s moves the kubeconfig Secret out of the account namespace", f.scpcName)
		}
	}
}
