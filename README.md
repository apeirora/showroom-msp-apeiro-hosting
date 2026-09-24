# showroom-msp-sap-provider

SAP Hosting provider for simple-cluster, packaging the OpenStack and IronCore
Gardener flavors into one provider stack.

## Description

Exposes two namespaced CRs in group `sap.eu` — `OpenStackShoot` and
`IronCoreShoot` — via a single APIExport, sync-agent and operator. When a KCP
account workspace binds this provider (alongside simple-cluster), the operator's
NamespaceReconciler materialises both `sap-gardener-openstack` and
`sap-gardener-ironcore` SimpleClusterProviderConfigs (namespaced, owned by the
account namespace) so users can create SimpleClusters against either flavor. A
flavor CR's free-text `spec` (a Gardener Shoot manifest) overrides that flavor's
SCPC shoot template; an empty spec uses the operator's built-in default.

## Integrations & repository layout

This provider ships **three** integrations and contains **two** Go modules — the
provider bundle at the repo root and `ironcore-operator/` (the IronCore VM operator):

| Integration | Resource (Kind) | Operator + CRD | Location |
|---|---|---|---|
| OpenStackShoot | `Shoot` (`core.gardener.cloud`) | bundle operator (`sap.eu`) | this repo (root) |
| IronCoreShoot | `Shoot` (`core.gardener.cloud`, IronCore infra) | bundle operator (`sap.eu`) | this repo (root) |
| IronCore VM | `IronCoreVM` (`sap.eu`) | `ironcore-operator` | this repo (`ironcore-operator/`) |

All three integrations live in this repository — it is self-contained. The two
Gardener flavors talk directly to Gardener. The **IronCore VM** integration provides
both the operator (`ironcore-operator/`, image + chart `ironcore-operator`) and its
Platform Mesh surface: the `ironcorevms.sap.eu` APIResourceSchema
projected namespaced into the KCP workspace, the `publish-ironcore-ironcorevms`
PublishedResource in the sync-agent, and the VM CRUD UI in `ui/src/lib/vm.ts` +
`ui/src/components/Vm*`.

> The IronCore VM operator reconciles `IronCoreVM` CRs and resolves `SSHKey`
> references (`sshkeys.apeirora.eu`) via the CRD contract — the shared
> `sshkey-operator` (in the simple-cluster monorepo) must be deployed in the cluster
> at runtime.

> **Schema note:** `charts/sapprovider-pm-integration` carries the `IronCoreVM`
> APIResourceSchema used for the KCP projection; it is generated from
> `ironcore-operator/config/crd/bases`. Bump `vmSchemaRevision` and refresh it when
> that CRD changes.

## Getting Started

### Prerequisites
- go version v1.24.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/showroom-msp-sap-provider:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/showroom-msp-sap-provider:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/showroom-msp-sap-provider:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/showroom-msp-sap-provider/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.

## Contributing
Contributions go via pull requests with Conventional-Commit titles (release-please
derives the changelog and version bump from them). Before opening a PR, run
`make test` and `make lint`, and regenerate manifests with `make manifests generate`
if you changed API types.

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

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

---

<p align="center">
  <sub>Built with care by <a href="https://apeirora.eu/">ApeiroRA</a></sub>
</p>

<p align="center">
  <a href="https://apeirora.eu/">
    <img src="docs/assets/eu-funded.svg" alt="Funded by the European Union -- NextGenerationEU. Supported by the Federal Ministry for Economic Affairs and Energy on the basis of a decision by the German Bundestag." width="400" />
  </a>
</p>
