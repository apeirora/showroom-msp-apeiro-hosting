package controller

import (
	"context"
	"fmt"
	"testing"

	api "github.com/apeirora/showroom-msp-apeiro-hosting/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type adoptionGardener struct {
	*fakeGardener
	creates, deletes, adopts int
	failAdopt                bool
}

func (g *adoptionGardener) CreateShoot(ctx context.Context, shoot *unstructured.Unstructured) error {
	g.creates++
	return g.fakeGardener.CreateShoot(ctx, shoot)
}
func (g *adoptionGardener) DeleteShoot(ctx context.Context, name string) (bool, error) {
	g.deletes++
	return g.fakeGardener.DeleteShoot(ctx, name)
}
func (g *adoptionGardener) AdoptShoot(_ context.Context, live *unstructured.Unstructured, owner string) error {
	g.adopts++
	if g.failAdopt {
		return fmt.Errorf("adoption conflict")
	}
	annotations := live.GetAnnotations()
	annotations[shootOwnerAnnotation] = owner
	delete(annotations, legacyShootOwnerAnnotation)
	live.SetAnnotations(annotations)
	return nil
}

func TestAdoptionThroughReconcileNeverCreatesOrDeletes(t *testing.T) {
	for _, kind := range []string{"IronCoreShoot", "OpenStackShoot"} {
		for _, scenario := range []string{"legacy", "already-adopted", "absent", "uid-mismatch", "namespace-mismatch", "name-mismatch", "deleting", "foreign-legacy-owner", "foreign-new-owner", "missing-owner", "patch-conflict"} {
			t.Run(kind+"/"+scenario, func(t *testing.T) {
				scheme := runtime.NewScheme()
				if err := corev1.AddToScheme(scheme); err != nil {
					t.Fatal(err)
				}
				if err := api.AddToScheme(scheme); err != nil {
					t.Fatal(err)
				}
				scheme.AddKnownTypeWithName(scpcGVK, &unstructured.Unstructured{})
				meta := metav1.ObjectMeta{Name: "test", Namespace: "tenant", UID: "new-cr", Annotations: map[string]string{AdoptShootUIDAnnotation: "shoot-uid"}}
				newShoot := func() FlavorShoot {
					if kind == "OpenStackShoot" {
						return &api.OpenStackShoot{}
					}
					return &api.IronCoreShoot{}
				}
				sh := newShoot()
				sh.SetName(meta.Name)
				sh.SetNamespace(meta.Namespace)
				sh.SetUID(meta.UID)
				sh.SetAnnotations(meta.Annotations)
				secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "garden", Namespace: "tenant"}, Data: map[string][]byte{"kubeconfig": []byte("config")}}
				c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(sh).WithObjects(sh, secret).Build()
				g := &adoptionGardener{fakeGardener: newFakeGardener()}
				name := gardenShootName(kind, sh.GetNamespace(), sh.GetName())
				live := &unstructured.Unstructured{Object: map[string]any{"apiVersion": "core.gardener.cloud/v1beta1", "kind": "Shoot", "metadata": map[string]any{"name": name, "namespace": testGardenNamespace, "uid": "shoot-uid", "resourceVersion": "1", "annotations": map[string]any{legacyShootOwnerAnnotation: "tenant/test"}}, "status": map[string]any{"lastOperation": map[string]any{"state": "Succeeded"}}}}
				switch scenario {
				case "already-adopted":
					live.SetAnnotations(map[string]string{shootOwnerAnnotation: "tenant/test"})
				case "uid-mismatch":
					live.SetUID("replacement")
				case "namespace-mismatch":
					live.SetNamespace("foreign")
				case "name-mismatch":
					live.SetName("foreign")
				case "deleting":
					now := metav1.Now()
					live.SetDeletionTimestamp(&now)
				case "foreign-legacy-owner":
					live.SetAnnotations(map[string]string{legacyShootOwnerAnnotation: "other/cr"})
				case "foreign-new-owner":
					live.SetAnnotations(map[string]string{legacyShootOwnerAnnotation: "tenant/test", shootOwnerAnnotation: "other/cr"})
				case "missing-owner":
					live.SetAnnotations(nil)
				case "patch-conflict":
					g.failAdopt = true
				}
				if scenario != "absent" {
					g.shoots[name] = live
				}
				r := &ShootReconciler{Client: c, Scheme: scheme, KindName: kind, SCPCName: "config", New: newShoot, KubeconfigSource: KubeconfigSourceGlobal, GlobalKubeconfigSecretName: "garden", GlobalKubeconfigSecretNamespace: "tenant", NewGardener: func([]byte, string) (Gardener, error) { return g, nil }, defaults: &flavorDefaults{targetNamespace: testGardenNamespace}}
				success := scenario == "legacy" || scenario == "already-adopted"
				for i := 0; i < 2; i++ {
					_, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: client.ObjectKeyFromObject(sh)})
					if (err == nil) != success {
						t.Fatalf("success=%v, err=%v", success, err)
					}
				}
				if g.creates != 0 || g.deletes != 0 {
					t.Fatalf("creates=%d deletes=%d", g.creates, g.deletes)
				}
				if success {
					if string(live.GetUID()) != "shoot-uid" || live.GetAnnotations()[shootOwnerAnnotation] != "tenant/test" || live.GetAnnotations()[legacyShootOwnerAnnotation] != "" {
						t.Fatal("Shoot identity or owner differs")
					}
					if err := c.Get(context.Background(), client.ObjectKeyFromObject(sh), sh); err != nil {
						t.Fatal(err)
					}
					if sh.GetShootStatus().Phase != phaseReady || sh.GetShootStatus().KubeconfigSecretRef == nil {
						t.Fatalf("adopted CR not ready: %+v", sh.GetShootStatus())
					}
					if scenario == "legacy" && g.adopts != 1 {
						t.Fatalf("adoption count=%d", g.adopts)
					}
				}
			})
		}
	}
}
