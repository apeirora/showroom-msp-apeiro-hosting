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
	"github.com/apeirora/showroom-msp-apeiro-hosting/internal/secretpolicy"
)

const testGardenNamespace = "garden-dev"

const shootManifestHeader = "apiVersion: core.gardener.cloud/v1beta1\nkind: Shoot\n"

func TestShootRejectsUnapprovedReferencesBeforeGardener(t *testing.T) {
	for _, kind := range []string{"OpenStackShoot", "IronCoreShoot"} {
		for _, spec := range []string{
			"spec:\n  credentialsBindingName: other\n",
			"spec:\n  secretBindingName: other\n",
			"spec:\n  dns:\n    providers:\n    - secretName: other\n",
			"spec:\n  dns:\n    providers:\n    - credentialsRef:\n        apiVersion: v1\n        kind: Secret\n        name: other\n",
			"spec:\n  credentialsBindingName: 42\n",
			"spec:\n  dns:\n    providers:\n    - secretName: 42\n",
			"spec:\n  dns:\n    providers: invalid\n",
			"spec:\n  dns:\n    providers:\n    - invalid\n",
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
					v.Spec = shootManifestHeader + spec
				case *api.IronCoreShoot:
					v.Spec = shootManifestHeader + spec
				}
				secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "garden", Namespace: "tenant"}, Data: map[string][]byte{"kubeconfig": []byte("config")}}
				c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(sh).WithObjects(sh, secret).Build()
				spy := &shootSecretSpy{Client: c}
				calls := 0
				r := &ShootReconciler{Client: spy, Scheme: scheme, New: newShoot, KindName: kind, SCPCName: "config",
					KubeconfigSource: KubeconfigSourceGlobal, GlobalKubeconfigSecretName: "garden", GlobalKubeconfigSecretNamespace: "tenant",
					defaults:    &flavorDefaults{targetNamespace: testGardenNamespace},
					NewGardener: func([]byte, string) (Gardener, error) { calls++; return newFakeGardener(), nil },
				}
				_, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: client.ObjectKeyFromObject(sh)})
				if err != nil {
					t.Fatal(err)
				}
				if spy.reads != 0 {
					t.Fatalf("credential reads = %d", spy.reads)
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

type shootSecretSpy struct {
	client.Client
	reads int
}

func (s *shootSecretSpy) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	if _, ok := obj.(*corev1.Secret); ok {
		s.reads++
	}
	return s.Client.Get(ctx, key, obj, opts...)
}

func TestShootReferencePolicyAtCreateBoundary(t *testing.T) {
	for _, tc := range []struct {
		name, spec string
		entries    []secretpolicy.Entry
		wantOK     bool
	}{
		{name: "approved", spec: "spec:\n  credentialsBindingName: infra\n  dns:\n    providers:\n    - secretName: dns\n", entries: []secretpolicy.Entry{
			{Namespace: testGardenNamespace, Name: "infra", Purpose: secretpolicy.GardenerCredentialBinding},
			{Namespace: testGardenNamespace, Name: "dns", Purpose: secretpolicy.GardenerDNS},
		}, wantOK: true},
		{name: "legacy approved", spec: "spec:\n  secretBindingName: infra\n", entries: []secretpolicy.Entry{{Namespace: testGardenNamespace, Name: "infra", Purpose: secretpolicy.GardenerSecretBinding}}, wantOK: true},
		{name: "wrong namespace", spec: "spec:\n  credentialsBindingName: infra\n", entries: []secretpolicy.Entry{{Namespace: "tenant", Name: "infra", Purpose: secretpolicy.GardenerCredentialBinding}}},
		{name: "wrong purpose", spec: "spec:\n  secretBindingName: infra\n", entries: []secretpolicy.Entry{{Namespace: testGardenNamespace, Name: "infra", Purpose: secretpolicy.GardenerCredentialBinding}}},
		{name: "second DNS denied", spec: "spec:\n  dns:\n    providers:\n    - secretName: dns\n    - secretName: other\n", entries: []secretpolicy.Entry{{Namespace: testGardenNamespace, Name: "dns", Purpose: secretpolicy.GardenerDNS}}},
		{name: "empty binding denied", spec: "spec:\n  credentialsBindingName: ''\n"},
		{name: "no refs", spec: "spec:\n  region: custom\n", wantOK: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sh := &api.OpenStackShoot{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "tenant"}, Spec: shootManifestHeader + tc.spec}
			r := &ShootReconciler{SecretPolicy: secretpolicy.Policy{Allowlist: tc.entries}, defaults: &flavorDefaults{targetNamespace: testGardenNamespace}}
			g := newFakeGardener()
			err := r.createShoot(context.Background(), g, sh, "remote")
			if (err == nil) != tc.wantOK {
				t.Fatalf("error = %v, wantOK = %v", err, tc.wantOK)
			}
			if tc.wantOK {
				created := g.shoots["remote"]
				if created == nil || created.GetNamespace() != testGardenNamespace {
					t.Fatal("expected remote Shoot in configured namespace")
				}
			} else if len(g.shoots) != 0 {
				t.Fatal("Gardener received denied manifest")
			}
		})
	}
}

func TestShootApprovedReferencesThroughReconcile(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = api.AddToScheme(scheme)
	sh := &api.IronCoreShoot{ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "tenant"}, Spec: shootManifestHeader + "metadata:\n  namespace: ignored\nspec:\n  credentialsBindingName: ironcore-infra\n  dns:\n    providers:\n    - credentialsRef:\n        apiVersion: v1\n        kind: Secret\n        name: designate-dns\n"}
	secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "garden", Namespace: "tenant"}, Data: map[string][]byte{"kubeconfig": []byte("config")}}
	c := fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(sh).WithObjects(sh, secret).Build()
	g := newFakeGardener()
	calls := 0
	r := &ShootReconciler{Client: c, Scheme: scheme, KindName: "IronCoreShoot", SCPCName: "config", New: func() FlavorShoot { return &api.IronCoreShoot{} },
		KubeconfigSource: KubeconfigSourceGlobal, GlobalKubeconfigSecretName: "garden", GlobalKubeconfigSecretNamespace: "tenant",
		NewGardener: func([]byte, string) (Gardener, error) { calls++; return g, nil }, defaults: &flavorDefaults{targetNamespace: testGardenNamespace},
		SecretPolicy: secretpolicy.Policy{Allowlist: []secretpolicy.Entry{
			{Namespace: testGardenNamespace, Name: "ironcore-infra", Purpose: secretpolicy.GardenerCredentialBinding},
			{Namespace: testGardenNamespace, Name: "designate-dns", Purpose: secretpolicy.GardenerDNS},
		}},
	}
	if _, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: client.ObjectKeyFromObject(sh)}); err != nil {
		t.Fatal(err)
	}
	if calls != 1 || len(g.shoots) != 1 {
		t.Fatalf("factory calls = %d, shoots = %d", calls, len(g.shoots))
	}
	if err := c.Get(context.Background(), client.ObjectKeyFromObject(sh), sh); err != nil {
		t.Fatal(err)
	}
	if sh.Status.Phase != phaseProvisioning {
		t.Fatalf("phase = %s", sh.Status.Phase)
	}
}

func TestDNSCredentialsRefPolicy(t *testing.T) {
	r := &ShootReconciler{defaults: &flavorDefaults{targetNamespace: testGardenNamespace}, SecretPolicy: secretpolicy.Policy{Allowlist: []secretpolicy.Entry{
		{Namespace: testGardenNamespace, Name: "dns", Purpose: secretpolicy.GardenerDNS},
		{Namespace: "approved-dns", Name: "dns", Purpose: secretpolicy.GardenerDNS},
	}}}
	for _, tc := range []struct {
		name, fields string
		wantOK       bool
	}{
		{"default namespace", "", true},
		{"explicit target", "        namespace: garden-dev\n", true},
		{"explicit approved namespace", "        namespace: approved-dns\n", true},
		{"foreign namespace", "        namespace: other\n", false},
		{"invalid namespace type", "        namespace: 42\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sh := &api.OpenStackShoot{Spec: shootManifestHeader + "spec:\n  dns:\n    providers:\n    - credentialsRef:\n        apiVersion: v1\n        kind: Secret\n        name: dns\n" + tc.fields}
			err := r.validateRawShoot(sh)
			if (err == nil) != tc.wantOK {
				t.Fatalf("err = %v, wantOK = %v", err, tc.wantOK)
			}
		})
	}
	for _, spec := range []string{
		"spec:\n  dns:\n    providers:\n    - credentialsRef: invalid\n",
		"spec:\n  dns:\n    providers:\n    - credentialsRef:\n        apiVersion: other/v1\n        kind: Secret\n        name: dns\n",
		"spec:\n  dns:\n    providers:\n    - credentialsRef:\n        apiVersion: v1\n        kind: Other\n        name: dns\n",
		"spec:\n  dns:\n    providers:\n    - credentialsRef:\n        apiVersion: v1\n        kind: Secret\n",
		"spec:\n  dns:\n    providers:\n    - credentialsRef:\n        apiVersion: v1\n        kind: Secret\n        name: dns\n      secretName: other\n",
	} {
		if err := r.validateRawShoot(&api.OpenStackShoot{Spec: shootManifestHeader + spec}); err == nil {
			t.Fatalf("accepted invalid DNS reference: %s", spec)
		}
	}
}
