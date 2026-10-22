package controller

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

// AdoptShootUIDAnnotation pins a migration to an existing Gardener Shoot.
// Keep this annotation to prevent creation if that Shoot disappears.
const AdoptShootUIDAnnotation = "apeiro.eu/adopt-shoot-uid"

const legacyShootOwnerAnnotation = "sap.eu/owner"

// shootAdopter changes metadata with UID and resourceVersion preconditions.
type shootAdopter interface {
	AdoptShoot(context.Context, *unstructured.Unstructured, string) error
}

func (r *ShootReconciler) validateAdoption(sh FlavorShoot, live *unstructured.Unstructured) error {
	expected := sh.GetAnnotations()[AdoptShootUIDAnnotation]
	if expected == "" {
		return nil
	}
	if r.KindName != "IronCoreShoot" {
		return fmt.Errorf("shoot adoption supports only IronCoreShoot")
	}
	if string(live.GetUID()) != expected || live.GetNamespace() != r.defaults.targetNamespace ||
		live.GetName() != gardenShootName(r.KindName, sh.GetNamespace(), sh.GetName()) ||
		!live.GetDeletionTimestamp().IsZero() {
		return fmt.Errorf("Gardener Shoot adoption identity mismatch or deletion in progress")
	}
	owner := sh.GetNamespace() + "/" + sh.GetName()
	legacy, current := live.GetAnnotations()[legacyShootOwnerAnnotation], live.GetAnnotations()[shootOwnerAnnotation]
	if (legacy != "" && legacy != owner) || (current != "" && current != owner) || (legacy == "" && current == "") {
		return fmt.Errorf("Gardener Shoot adoption ownership mismatch")
	}
	return nil
}

func (r *ShootReconciler) adoptShoot(ctx context.Context, g Gardener, sh FlavorShoot, live *unstructured.Unstructured) error {
	if sh.GetAnnotations()[AdoptShootUIDAnnotation] == "" {
		return nil
	}
	if err := r.validateAdoption(sh, live); err != nil {
		return err
	}
	owner := sh.GetNamespace() + "/" + sh.GetName()
	if live.GetAnnotations()[shootOwnerAnnotation] == owner && live.GetAnnotations()[legacyShootOwnerAnnotation] == "" {
		return nil
	}
	adopter, ok := g.(shootAdopter)
	if !ok {
		return fmt.Errorf("Gardener client does not support Shoot adoption")
	}
	return adopter.AdoptShoot(ctx, live, owner)
}

func (g *gardenerClient) AdoptShoot(ctx context.Context, live *unstructured.Unstructured, owner string) error {
	if live.GetUID() == "" || live.GetResourceVersion() == "" {
		return fmt.Errorf("shoot adoption requires UID and resourceVersion")
	}
	patch, err := json.Marshal(map[string]any{"metadata": map[string]any{
		"uid": string(live.GetUID()), "resourceVersion": live.GetResourceVersion(),
		"annotations": map[string]any{shootOwnerAnnotation: owner, legacyShootOwnerAnnotation: nil},
	}})
	if err != nil {
		return err
	}
	_, err = g.dyn.Resource(shootGVR).Namespace(g.namespace).Patch(ctx, live.GetName(), types.MergePatchType, patch, metav1.PatchOptions{})
	return err
}
