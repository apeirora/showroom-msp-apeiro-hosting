package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

func TestGardenerAdoptionPatchPreconditions(t *testing.T) {
	for _, conflict := range []bool{false, true} {
		calls := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			if r.Method != "PATCH" || r.URL.Path != "/apis/core.gardener.cloud/v1beta1/namespaces/garden-dev/shoots/old" {
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			}
			var patch map[string]any
			if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
				t.Error(err)
			}
			metadata := patch["metadata"].(map[string]any)
			if metadata["uid"] != "pinned" || metadata["resourceVersion"] != "7" {
				t.Error("missing identity preconditions")
			}
			annotations := metadata["annotations"].(map[string]any)
			if annotations[shootOwnerAnnotation] != "tenant/cr" || annotations[legacyShootOwnerAnnotation] != nil {
				t.Error("ownership patch differs")
			}
			w.Header().Set("Content-Type", "application/json")
			if conflict {
				w.WriteHeader(http.StatusConflict)
				_ = json.NewEncoder(w).Encode(metav1.Status{Status: "Failure", Reason: metav1.StatusReasonConflict, Code: 409})
				return
			}
			_, _ = w.Write([]byte(`{"apiVersion":"core.gardener.cloud/v1beta1","kind":"Shoot","metadata":{"name":"old"}}`))
		}))
		dyn, err := dynamic.NewForConfig(&rest.Config{Host: server.URL})
		if err != nil {
			t.Fatal(err)
		}
		g := &gardenerClient{dyn: dyn, namespace: "garden-dev"}
		live := &unstructured.Unstructured{}
		live.SetName("old")
		live.SetUID("pinned")
		live.SetResourceVersion("7")
		err = g.AdoptShoot(context.Background(), live, "tenant/cr")
		server.Close()
		if (err != nil) != conflict || calls != 1 {
			t.Fatalf("conflict=%v err=%v calls=%d", conflict, err, calls)
		}
	}
}

func TestGardenerAdoptedDeletionRejectsReplacement(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != "GET" {
			t.Errorf("replacement received %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"apiVersion":"core.gardener.cloud/v1beta1","kind":"Shoot","metadata":{"name":"old","uid":"replacement","resourceVersion":"2"}}`))
	}))
	defer server.Close()
	dyn, err := dynamic.NewForConfig(&rest.Config{Host: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	g := &gardenerClient{dyn: dyn, namespace: "garden-dev"}
	if _, err := g.DeleteShootWithUID(context.Background(), "old", "pinned"); err == nil {
		t.Fatal("accepted replacement")
	}
	if calls != 1 {
		t.Fatalf("requests=%d", calls)
	}
}
