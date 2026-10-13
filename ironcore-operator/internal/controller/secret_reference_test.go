package controller

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	api "github.com/apeirora/showroom-msp-apeiro-hosting/ironcore-operator/api/v1alpha1"
)

type referenceSpy struct {
	client.Client
	reads []client.ObjectKey
}

func (s *referenceSpy) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	if _, isVM := obj.(*api.IronCoreVM); !isVM {
		s.reads = append(s.reads, key)
	}
	return s.Client.Get(ctx, key, obj, opts...)
}

func TestVMRejectsForeignReferencesBeforeCredentialRead(t *testing.T) {
	for _, kind := range []string{"kubeconfig", "Secret", "SSHKey"} {
		for _, deleting := range []bool{false, true} {
			t.Run(kind+map[bool]string{true: "/delete", false: "/create"}[deleting], func(t *testing.T) {
				scheme := runtime.NewScheme()
				_ = corev1.AddToScheme(scheme)
				_ = api.AddToScheme(scheme)
				vm := &api.IronCoreVM{ObjectMeta: metav1.ObjectMeta{Name: "vm", Namespace: "tenant", Finalizers: []string{api.FinalizerName}}}
				vm.Spec.IronCoreKubeconfig = corev1.SecretReference{Name: "credentials"}
				if kind == "kubeconfig" {
					vm.Spec.IronCoreKubeconfig.Namespace = "other"
				} else {
					vm.Spec.SSHKeyRef = &corev1.TypedObjectReference{Name: "key", Kind: kind, Namespace: ptr("other")}
					if kind == "SSHKey" {
						vm.Spec.SSHKeyRef.APIGroup = ptr(api.SSHKeyAPIGroup)
					}
				}
				if deleting {
					now := metav1.Now()
					vm.DeletionTimestamp = &now
				}
				spy := &referenceSpy{Client: fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(vm).WithObjects(vm).Build()}
				r := &IronCoreVMReconciler{Client: spy, Scheme: scheme}
				_, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "tenant", Name: "vm"}})
				if err != nil {
					t.Fatal(err)
				}
				if len(spy.reads) != 0 {
					t.Fatalf("reference reads before rejection: %v", spy.reads)
				}
				got := &api.IronCoreVM{}
				if err := spy.Client.Get(context.Background(), client.ObjectKeyFromObject(vm), got); err != nil {
					t.Fatal(err)
				}
				if got.Status.Phase != api.PhaseFailed {
					t.Fatalf("phase = %s, want Failed", got.Status.Phase)
				}
			})
		}
	}
}

func TestVMDirectHelpersRejectForeignReferences(t *testing.T) {
	for _, helper := range []string{"auth", "keypair", "sshkey"} {
		t.Run(helper, func(t *testing.T) {
			vm := &api.IronCoreVM{ObjectMeta: metav1.ObjectMeta{Name: "vm", Namespace: "tenant"}}
			vm.Spec.IronCoreKubeconfig = corev1.SecretReference{Name: "config"}
			vm.Spec.SSHKeyRef = &corev1.TypedObjectReference{Name: "key", Kind: "SSHKey", APIGroup: ptr(api.SSHKeyAPIGroup), Namespace: ptr("other")}
			spy := &referenceSpy{}
			r := &IronCoreVMReconciler{Client: spy}
			var err error
			switch helper {
			case "auth":
				_, err = r.ironcoreClients(context.Background(), vm)
			case "keypair":
				_, _, err = r.ensureSSHKeypair(context.Background(), vm)
			case "sshkey":
				_, _, err = r.resolveSSHKeyResource(context.Background(), vm)
			}
			if err == nil {
				t.Fatal("expected denial")
			}
			if len(spy.reads) != 0 {
				t.Fatalf("reads = %v", spy.reads)
			}
		})
	}
}

func TestVMLocalSecretReferences(t *testing.T) {
	for _, ns := range []*string{nil, ptr(""), ptr("tenant")} {
		t.Run(map[bool]string{true: "omitted", false: "explicit"}[ns == nil], func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = corev1.AddToScheme(scheme)
			_ = api.AddToScheme(scheme)
			secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "key", Namespace: "tenant"}, Data: map[string][]byte{api.SSHSecretPublicKeyKey: []byte("public")}}
			spy := &referenceSpy{Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()}
			vm := &api.IronCoreVM{ObjectMeta: metav1.ObjectMeta{Name: "vm", Namespace: "tenant"}}
			vm.Spec.IronCoreKubeconfig = corev1.SecretReference{Name: "config"}
			vm.Spec.SSHKeyRef = &corev1.TypedObjectReference{Name: "key", Kind: "Secret", Namespace: ns}
			r := &IronCoreVMReconciler{Client: spy, Scheme: scheme}
			ref, pub, err := r.ensureSSHKeypair(context.Background(), vm)
			if err != nil {
				t.Fatal(err)
			}
			if pub != "public" || ref.Namespace != "tenant" || ref.Name != "key" {
				t.Fatalf("ref = %v, pub = %q", ref, pub)
			}
			if len(spy.reads) != 1 || spy.reads[0] != ref {
				t.Fatalf("reads = %v", spy.reads)
			}
		})
	}
}

func TestVMLocalKubeconfigNamespace(t *testing.T) {
	for _, ns := range []string{"", "tenant"} {
		t.Run(ns, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = corev1.AddToScheme(scheme)
			secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "config", Namespace: "tenant"}, Data: map[string][]byte{"kubeconfig": []byte("invalid")}}
			spy := &referenceSpy{Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()}
			vm := &api.IronCoreVM{ObjectMeta: metav1.ObjectMeta{Name: "vm", Namespace: "tenant"}}
			vm.Spec.IronCoreKubeconfig = corev1.SecretReference{Name: "config", Namespace: ns}
			r := &IronCoreVMReconciler{Client: spy}
			_, err := r.ironcoreClients(context.Background(), vm)
			if err == nil {
				t.Fatal("invalid kubeconfig must fail construction")
			}
			want := client.ObjectKey{Namespace: "tenant", Name: "config"}
			if len(spy.reads) != 1 || spy.reads[0] != want {
				t.Fatalf("reads = %v", spy.reads)
			}
		})
	}
}

func TestVMLocalGeneratedSecret(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = api.AddToScheme(scheme)
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	vm := &api.IronCoreVM{ObjectMeta: metav1.ObjectMeta{Name: "vm", Namespace: "tenant", UID: "vm-id"}}
	r := &IronCoreVMReconciler{Client: c, Scheme: scheme}
	ref, pub, err := r.ensureSSHKeypair(context.Background(), vm)
	if err != nil {
		t.Fatal(err)
	}
	secret := &corev1.Secret{}
	if err := c.Get(context.Background(), ref, secret); err != nil {
		t.Fatal(err)
	}
	if ref.Namespace != vm.Namespace || ref.Name != "vm-ssh" || len(pub) == 0 || len(secret.Data[api.SSHSecretPrivateKeyKey]) == 0 {
		t.Fatalf("invalid generated keypair at %v", ref)
	}
	if len(secret.OwnerReferences) != 1 || secret.OwnerReferences[0].UID != vm.UID {
		t.Fatal("missing VM owner")
	}
}
