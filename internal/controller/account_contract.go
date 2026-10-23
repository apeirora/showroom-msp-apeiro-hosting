package controller

import (
	"context"
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const accountContractAnnotation = "showroom.apeirora.eu/account-resource-contract"
const accountContractCondition corev1.NamespaceConditionType = "HostingAccountContractReady"
const accountContractResource = "simpleclusterproviderconfigs"

// contractVersion reports whether the API supports account SCPCs.
// An API without the annotation predates the contract and keeps legacy local SCPCs.
func contractVersion(obj *unstructured.Unstructured) (bool, error) {
	value, ok := obj.GetAnnotations()[accountContractAnnotation]
	if !ok {
		return false, nil
	}
	version, err := strconv.Atoi(value)
	if err != nil || version < 1 {
		return false, fmt.Errorf("SimpleClusterProviderConfig API requires %s >= 1; deploy the universal account core first", accountContractAnnotation)
	}
	return true, nil
}

// accountContract reads metadata from the same API server that receives the SCPC.
// It does not depend on the location or Deployment of the SimpleCluster controller.
func (r *NamespaceReconciler) accountContract(ctx context.Context) (bool, error) {
	crd := &unstructured.Unstructured{}
	crd.SetGroupVersionKind(schema.GroupVersionKind{Group: "apiextensions.k8s.io", Version: "v1", Kind: "CustomResourceDefinition"})
	err := r.Get(ctx, client.ObjectKey{Name: accountContractResource + ".showroom.apeirora.eu"}, crd)
	if err == nil {
		return contractVersion(crd)
	}
	if !apierrors.IsNotFound(err) && !apimeta.IsNoMatchError(err) {
		return false, fmt.Errorf("read SimpleClusterProviderConfig API contract: %w", err)
	}
	// A bound KCP API has no local CRD. Resolve only the schema actually bound to this resource.
	bindings := &unstructured.UnstructuredList{}
	bindings.SetGroupVersionKind(schema.GroupVersionKind{Group: "apis.kcp.io", Version: "v1alpha2", Kind: "APIBindingList"})
	if err := r.List(ctx, bindings); err != nil {
		return false, fmt.Errorf("SimpleClusterProviderConfig API contract unavailable: %w", err)
	}
	var schemaName, schemaUID string
	for _, binding := range bindings.Items {
		phase, _, _ := unstructured.NestedString(binding.Object, "status", "phase")
		if phase != "Bound" || binding.GetUID() == "" || !binding.GetDeletionTimestamp().IsZero() {
			continue
		}
		initial, current := false, false
		conditions, _, _ := unstructured.NestedSlice(binding.Object, "status", "conditions")
		for _, entry := range conditions {
			condition, ok := entry.(map[string]any)
			if !ok {
				continue
			}
			if condition["type"] == "InitialBindingCompleted" {
				initial = condition["status"] == "True"
			}
			if condition["type"] == "BindingUpToDate" {
				current = condition["status"] == "True"
			}
		}
		if !initial || !current {
			continue
		}
		resources, _, _ := unstructured.NestedSlice(binding.Object, "status", "boundResources")
		for _, entry := range resources {
			resource, ok := entry.(map[string]any)
			if !ok || resource["group"] != "showroom.apeirora.eu" || resource["resource"] != accountContractResource {
				continue
			}
			name, _, _ := unstructured.NestedString(resource, "schema", "name")
			identity, _, _ := unstructured.NestedString(resource, "schema", "identityHash")
			uid, _, _ := unstructured.NestedString(resource, "schema", "UID")
			if name == "" || identity == "" || uid == "" || (schemaName != "" && (schemaName != name || schemaUID != uid)) {
				return false, fmt.Errorf("ambiguous or incomplete bound SimpleClusterProviderConfig API schema")
			}
			schemaName, schemaUID = name, uid
		}
	}
	if schemaName == "" {
		return false, fmt.Errorf("SimpleClusterProviderConfig API has no bound schema contract")
	}
	apiSchema := &unstructured.Unstructured{}
	apiSchema.SetGroupVersionKind(schema.GroupVersionKind{Group: "apis.kcp.io", Version: "v1alpha1", Kind: "APIResourceSchema"})
	if err := r.Get(ctx, client.ObjectKey{Name: schemaName}, apiSchema); err != nil {
		return false, fmt.Errorf("read bound SimpleClusterProviderConfig API schema contract: %w", err)
	}
	if string(apiSchema.GetUID()) != schemaUID {
		return false, fmt.Errorf("bound SimpleClusterProviderConfig schema UID differs")
	}
	group, _, _ := unstructured.NestedString(apiSchema.Object, "spec", "group")
	plural, _, _ := unstructured.NestedString(apiSchema.Object, "spec", "names", "plural")
	if group != "showroom.apeirora.eu" || plural != accountContractResource {
		return false, fmt.Errorf("bound SimpleClusterProviderConfig schema resource differs")
	}
	return contractVersion(apiSchema)
}

// accountContractStatus reports False only when SCPC writes are refused.
func (r *NamespaceReconciler) accountContractStatus(ctx context.Context, ns *corev1.Namespace, account bool, contractErr error) error {
	condition := corev1.NamespaceCondition{Type: accountContractCondition, Status: corev1.ConditionTrue, Reason: "AccountContractAvailable", Message: "SimpleClusterProviderConfig API supports account resource contract 1"}
	if !account {
		condition.Reason = "LegacyContract"
		condition.Message = "SimpleClusterProviderConfig API has no account resource contract; writing legacy local SCPCs"
	}
	if contractErr != nil {
		condition.Status = corev1.ConditionFalse
		condition.Reason = "AccountContractRequired"
		condition.Message = contractErr.Error()
	}
	patch := client.MergeFromWithOptions(ns.DeepCopy(), client.MergeFromWithOptimisticLock{})
	condition.LastTransitionTime = metav1.Now()
	found := false
	for i, existing := range ns.Status.Conditions {
		if existing.Type != accountContractCondition {
			continue
		}
		if existing.Status == condition.Status && existing.Reason == condition.Reason && existing.Message == condition.Message {
			return nil
		}
		if existing.Status == condition.Status {
			condition.LastTransitionTime = existing.LastTransitionTime
		}
		ns.Status.Conditions[i] = condition
		found = true
		break
	}
	if !found {
		ns.Status.Conditions = append(ns.Status.Conditions, condition)
	}
	return r.Status().Patch(ctx, ns, patch)
}
