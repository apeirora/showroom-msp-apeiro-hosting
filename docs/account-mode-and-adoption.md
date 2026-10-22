# Account children and Shoot adoption

Both IronCore and OpenStack SCPCs create children in the account workspace.
SimpleCluster uses the account as its only destination. No SCPC destination flag is required.

This configuration requires a SimpleCluster release with the universal account destination.
The account must bind the SimpleCluster and Apeiro exports.
The sync agent projects the account child into the account namespace on the service cluster.
The kubeconfig Secret returns to namespace `default` in the account.
Existing SimpleClusters need a reviewed account handoff before the new core controller resumes.

For a migration, set `apeiro.eu/adopt-shoot-uid` on the new IronCoreShoot or OpenStackShoot.
Use the UID of the real Gardener Shoot, not the UID of either IronCoreShoot object.
Keep the same physical namespace and CR name.
The controller derives the Shoot name from the kind, physical namespace, and CR name.

The adoption path requires:

- The exact configured garden project, derived name, and supplied Shoot UID.
- No Shoot deletion timestamp.
- The `sap.eu/owner` or `apeiro.eu/owner` annotation with the physical CR namespace and name.
- No conflicting value in either ownership annotation.

The controller replaces `sap.eu/owner` with `apeiro.eu/owner` through a metadata patch.
The patch includes the Shoot UID and resourceVersion.
A conflict stops adoption.
The controller keeps the existing Shoot and spec.
The controller reports Failed when adoption fails.
An absent Shoot never causes creation while the adoption annotation remains.

Keep both provider controllers and sync paths stopped during the object handoff.
Transfer the kubeconfig Secret owner reference before the new controller starts.
Keep the old objects until the new path is Ready.
Before old object deletion, stop the old provider and sync path.
Remove only their known finalizers.
Then delete the old objects with orphan propagation and UID preconditions.
Orphan propagation alone does not suppress the provider finalizer.
Never delete the Gardener Shoot as part of the handoff.

For rollback, stop the new path first.
Restore the old provider selection, secret owner references, and Shoot ownership annotation.
Restore the old controller and sync path only after the new path stops.
Orphan-delete the new objects with their provider controllers stopped.
Rehearse the handoff and rollback on disposable resources before a live migration.
