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
