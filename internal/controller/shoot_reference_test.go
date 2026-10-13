package controller

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	api "github.com/apeirora/showroom-msp-apeiro-hosting/api/v1alpha1"
)

func TestShootRejectsUnapprovedReferencesBeforeGardener(t *testing.T) {
	for _, kind := range []string{"OpenStackShoot", "IronCoreShoot"} {
		for _, spec := range []string{
			"spec:\n  credentialsBindingName: other\n",
			"spec:\n  secretBindingName: other\n",
			"spec:\n  dns:\n    providers:\n    - secretName: other\n",
		} {
			t.Run(kind+spec, func(t *testing.T) {
				scheme := runtime.NewScheme()
				_ = corev1.AddToScheme(scheme)
				_ = api.AddToScheme(scheme)
				var sh FlavorShoot
				newShoot := func() FlavorShoot { return &api.OpenStackShoot{} }
				if kind == "IronCoreShoot" {
					newShoot = func() FlavorShoot { return &api.IronCoreShoot{} }
				}
				sh = newShoot()
				sh.SetName("test")
				sh.SetNamespace("tenant")
				switch v := sh.(type) {
				case *api.OpenStackShoot:
					v.Spec = spec
				case *api.IronCoreShoot:
					v.Spec = spec
				}
				secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "garden", Namespace: "tenant"}, Data: map[string][]byte{"kubeconfig": []byte("config")}}
				c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(sh).WithObjects(sh, secret).Build()
				calls := 0
				r := &ShootReconciler{Client: c, Scheme: scheme, New: newShoot, KindName: kind, SCPCName: "config",
					KubeconfigSource: KubeconfigSourceGlobal, GlobalKubeconfigSecretName: "garden", GlobalKubeconfigSecretNamespace: "tenant",
					defaults:    &flavorDefaults{targetNamespace: "garden-dev"},
					NewGardener: func([]byte, string) (Gardener, error) { calls++; return newFakeGardener(), nil },
				}
				_, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: client.ObjectKeyFromObject(sh)})
				if err != nil {
					t.Fatal(err)
				}
				if calls != 0 {
					t.Fatalf("Gardener factory calls = %d, want 0", calls)
				}
				got := newShoot()
				if err := c.Get(context.Background(), client.ObjectKeyFromObject(sh), got); err != nil {
					t.Fatal(err)
				}
				if got.GetShootStatus().Phase != phaseFailed {
					t.Fatalf("phase = %s", got.GetShootStatus().Phase)
				}
			})
		}
	}
}
