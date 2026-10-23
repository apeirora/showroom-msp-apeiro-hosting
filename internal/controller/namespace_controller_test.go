/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	apeirov1alpha1 "github.com/apeirora/showroom-msp-apeiro-hosting/api/v1alpha1"
)

var _ = Describe("Namespace Controller", func() {
	ctx := context.Background()

	const (
		globalSecretName      = "apeiro-gardener-kubeconfig"
		globalSecretNamespace = "apeiro-provider"
		bothExports           = DefaultSimpleClusterExportName + "," + DefaultApeiroExportName
	)

	nsCounter := 0
	newNamespace := func(exports string) string {
		nsCounter++
		name := fmt.Sprintf("acct-%d", nsCounter)
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
		if exports != "" {
			ns.Annotations = map[string]string{bindExportsAnnotation: exports}
		}
		Expect(k8sClient.Create(ctx, ns)).To(Succeed())
		return name
	}

	setExports := func(name, exports string) {
		ns := &corev1.Namespace{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name}, ns)).To(Succeed())
		if ns.Annotations == nil {
			ns.Annotations = map[string]string{}
		}
		ns.Annotations[bindExportsAnnotation] = exports
		Expect(k8sClient.Update(ctx, ns)).To(Succeed())
	}

	newReconciler := func() *NamespaceReconciler {
		manifest, err := ConfigureIronCoreManifest(gardenDevConfig(GinkgoT()))
		Expect(err).NotTo(HaveOccurred())
		return &NamespaceReconciler{
			IronCoreManifest:                manifest,
			Client:                          k8sClient,
			Scheme:                          k8sClient.Scheme(),
			KubeconfigSource:                KubeconfigSourceGlobal,
			GlobalKubeconfigSecretName:      globalSecretName,
			GlobalKubeconfigSecretNamespace: globalSecretNamespace,
			ClientKubeconfigSecretName:      "apeiro-gardener-kubeconfig",
			SimpleClusterExportName:         DefaultSimpleClusterExportName,
			ApeiroExportName:                DefaultApeiroExportName,
		}
	}

	getSCPC := func(namespace, name string) (*unstructured.Unstructured, error) {
		scpc := &unstructured.Unstructured{}
		scpc.SetGroupVersionKind(scpcGVK)
		err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, scpc)
		return scpc, err
	}

	scpcShootTemplate := func(scpc *unstructured.Unstructured) string {
		resources, _, _ := unstructured.NestedSlice(scpc.Object, "spec", "resources")
		Expect(resources).NotTo(BeEmpty())
		first := resources[0].(map[string]interface{})
		tpl, _ := first["template"].(string)
		return tpl
	}

	reconcileNS := func(name string) {
		r := newReconciler()
		_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name}})
		Expect(err).NotTo(HaveOccurred())
	}

	It("blocks SCPC writes on an old API contract and reports readiness after the API upgrade", func() {
		crd := &unstructured.Unstructured{}
		crd.SetAPIVersion("apiextensions.k8s.io/v1")
		crd.SetKind("CustomResourceDefinition")
		key := types.NamespacedName{Name: "simpleclusterproviderconfigs.showroom.apeirora.eu"}
		Expect(k8sClient.Get(ctx, key, crd)).To(Succeed())
		annotations := crd.GetAnnotations()
		original := annotations[accountContractAnnotation]
		delete(annotations, accountContractAnnotation)
		crd.SetAnnotations(annotations)
		Expect(k8sClient.Update(ctx, crd)).To(Succeed())
		DeferCleanup(func() {
			Expect(k8sClient.Get(ctx, key, crd)).To(Succeed())
			values := crd.GetAnnotations()
			values[accountContractAnnotation] = original
			crd.SetAnnotations(values)
			Expect(k8sClient.Update(ctx, crd)).To(Succeed())
		})
		name := newNamespace(bothExports)
		reconcileNS(name)
		for _, scpcName := range allSCPCNames() {
			_, err := getSCPC(name, scpcName)
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
		}
		ns := &corev1.Namespace{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name}, ns)).To(Succeed())
		found := false
		for _, condition := range ns.Status.Conditions {
			if condition.Type == accountContractCondition {
				found = true
				Expect(condition.Status).To(Equal(corev1.ConditionFalse))
				Expect(condition.Reason).To(Equal("AccountContractRequired"))
			}
		}
		Expect(found).To(BeTrue())
		Expect(k8sClient.Get(ctx, key, crd)).To(Succeed())
		annotations = crd.GetAnnotations()
		annotations[accountContractAnnotation] = "2"
		crd.SetAnnotations(annotations)
		Expect(k8sClient.Update(ctx, crd)).To(Succeed())
		reconcileNS(name)
		for _, scpcName := range allSCPCNames() {
			_, err := getSCPC(name, scpcName)
			Expect(err).NotTo(HaveOccurred())
		}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name}, ns)).To(Succeed())
		for _, condition := range ns.Status.Conditions {
			if condition.Type == accountContractCondition {
				Expect(condition.Status).To(Equal(corev1.ConditionTrue))
			}
		}
	})

	It("creates both namespace-owned SCPCs when both exports are bound", func() {
		name := newNamespace(bothExports)
		reconcileNS(name)

		for _, scpcName := range allSCPCNames() {
			scpc, err := getSCPC(name, scpcName)
			Expect(err).NotTo(HaveOccurred(), "SCPC %s should exist", scpcName)

			Expect(scpc.GetLabels()).To(HaveKeyWithValue(managedByLabelKey, managedByNamespaceBind))

			owners := scpc.GetOwnerReferences()
			Expect(owners).To(HaveLen(1))
			Expect(owners[0].Kind).To(Equal("Namespace"))
			Expect(owners[0].Name).To(Equal(name))

			// The flavor's provider CR is applied locally into the account
			// namespace; the operator reaches the garden itself, so the ensured
			// SCPC carries no garden kubeconfig ref.
			targetNS, _, _ := unstructured.NestedString(scpc.Object, "spec", "target", "namespace")
			Expect(targetNS).To(BeEmpty())
			_, hasRef, _ := unstructured.NestedMap(scpc.Object, "spec", "target", "kubeconfigSecretRef")
			Expect(hasRef).To(BeFalse())

			// The local provider CR (OpenStackShoot/IronCoreShoot) carries the
			// plain Gardener Shoot manifest in its spec.
			Expect(scpcShootTemplate(scpc)).To(ContainSubstring("apeiro.eu/v1alpha1"))
			Expect(scpcShootTemplate(scpc)).To(ContainSubstring("core.gardener.cloud/v1beta1"))
		}
	})

	It("does not create SCPCs when the Apeiro export is absent", func() {
		name := newNamespace(DefaultSimpleClusterExportName)
		reconcileNS(name)

		for _, scpcName := range allSCPCNames() {
			_, err := getSCPC(name, scpcName)
			Expect(err).To(HaveOccurred())
		}
	})

	It("removes the managed SCPCs when the Apeiro binding is dropped", func() {
		name := newNamespace(bothExports)
		reconcileNS(name)
		for _, scpcName := range allSCPCNames() {
			_, err := getSCPC(name, scpcName)
			Expect(err).NotTo(HaveOccurred())
		}

		// Unbind Apeiro: the annotation now lists only SimpleCluster.
		setExports(name, DefaultSimpleClusterExportName)
		reconcileNS(name)

		for _, scpcName := range allSCPCNames() {
			_, err := getSCPC(name, scpcName)
			Expect(err).To(HaveOccurred())
		}
	})

	It("keeps the embedded SCPC templates even when a flavor CR carries a spec", func() {
		name := newNamespace(bothExports)
		custom := "apiVersion: core.gardener.cloud/v1beta1\nkind: Shoot\nmetadata:\n  name: custom\nspec:\n  region: eu\n"
		Expect(k8sClient.Create(ctx, &apeirov1alpha1.OpenStackShoot{
			ObjectMeta: metav1.ObjectMeta{Name: "svc", Namespace: name},
			Spec:       custom,
		})).To(Succeed())

		reconcileNS(name)

		// The CR spec feeds the CR's own Gardener Shoot, never the SCPC.
		os, err := getSCPC(name, openstackSCPCName)
		Expect(err).NotTo(HaveOccurred())
		Expect(scpcShootTemplate(os)).NotTo(Equal(custom))
		Expect(scpcShootTemplate(os)).To(ContainSubstring("core.gardener.cloud/v1beta1"))

		ic, err := getSCPC(name, ironcoreSCPCName)
		Expect(err).NotTo(HaveOccurred())
		Expect(scpcShootTemplate(ic)).To(ContainSubstring("ironcore"))
		Expect(scpcShootTemplate(ic)).To(ContainSubstring("ironcore-infra"))
		Expect(scpcShootTemplate(ic)).NotTo(ContainSubstring("wdf4-ora-dev"))
	})
	It("creates account SCPCs for both flavors without a destination flag", func() {
		name := newNamespace(bothExports)
		r := newReconciler()
		_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: types.NamespacedName{Name: name}})
		Expect(err).NotTo(HaveOccurred())
		for _, flavor := range flavors {
			scpc, err := getSCPC(name, flavor.scpcName)
			Expect(err).NotTo(HaveOccurred())
			for _, field := range []string{"namespace", "kcpAccount", "kubeconfigSecretRef"} {
				_, found, err := unstructured.NestedFieldNoCopy(scpc.Object, "spec", "target", field)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeFalse())
			}
			namespace, _, err := unstructured.NestedString(scpc.Object, "spec", "kubeconfigSource", "secretRef", "namespace")
			Expect(err).NotTo(HaveOccurred())
			Expect(namespace).To(Equal("default"))
		}
	})
})
