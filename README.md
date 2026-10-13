<p align="center">
  <img src="charts/sapprovider-ui/files/icon.svg" alt="SAP Hosting" width="96" />
</p>

<h1 align="center">SAP Hosting Provider</h1>

<p align="center">
  <strong>Gardener clusters and IronCore VMs on SAP infrastructure, delivered as a Platform Mesh provider.</strong>
</p>

<p align="center">
  <a href="https://github.com/apeirora/showroom-msp-sap-hosting/releases"><img src="https://img.shields.io/github/v/release/apeirora/showroom-msp-sap-hosting?style=flat-square&color=blue" alt="Release"></a>
  <a href="go.mod"><img src="https://img.shields.io/github/go-mod/go-version/apeirora/showroom-msp-sap-hosting?style=flat-square" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache%202.0-green?style=flat-square" alt="License"></a>
  <a href="https://github.com/apeirora/showroom-msp-sap-hosting/actions"><img src="https://img.shields.io/github/actions/workflow/status/apeirora/showroom-msp-sap-hosting/test.yml?style=flat-square&label=CI" alt="CI"></a>
  <a href="https://api.reuse.software/info/github.com/apeirora/showroom-msp-sap-hosting"><img src="https://api.reuse.software/badge/github.com/apeirora/showroom-msp-sap-hosting" alt="REUSE"></a>
</p>

<p align="center">
  Part of the <a href="https://apeirora.eu/">ApeiroRA</a> Platform Mesh ecosystem
</p>

---

## What is the SAP Hosting provider?

This repository is the SAP hosting provider for the Showroom MSP: it lets a
Platform Mesh account order **Gardener Kubernetes
clusters** (OpenStack or IronCore flavor) and **IronCore virtual machines** on
SAP infrastructure, and hands the resulting cluster's admin kubeconfig back as
a Secret next to the request.

It plugs into the
[simple-cluster](https://github.com/apeirora/showroom-msp-simple-cluster) stack:
when an account binds this provider, the operator materialises ready-to-use
`SimpleClusterProviderConfig`s in the account's namespace, so a plain
`SimpleCluster` request can be fulfilled on SAP infrastructure. Independently of
that, the account can create the provider's own CRs directly for a standalone
Gardener cluster or a single VM.

## The resources

All user-facing CRs are namespaced, in API group `sap.eu`:

| Kind | What it provisions | Reconciled by |
|---|---|---|
| `OpenStackShoot` (`osshoot`) | A Gardener `Shoot` on OpenStack | bundle operator (repo root) |
| `IronCoreShoot` (`icshoot`) | A Gardener `Shoot` on IronCore | bundle operator (repo root) |
| `IronCoreVM` | A single IronCore `Machine` (VM) | `ironcore-operator/` |

A flavor CR's free-text `spec` is a plain Gardener Shoot manifest used
verbatim; an empty `spec` uses the configured IronCore template or the
built-in OpenStack template:

```yaml
apiVersion: sap.eu/v1alpha1
kind: OpenStackShoot
metadata:
  name: demo
```

The operator creates the Shoot in the configured Gardener project, tracks its
lifecycle in `status.phase` (`Provisioning` → `Ready`), and maintains an admin
kubeconfig in the Secret named by `status.kubeconfigSecretRef` — minted via the
`adminkubeconfig` subresource and renewed automatically before expiry.

An `IronCoreVM` describes one machine on a target IronCore cluster (reached
through a kubeconfig Secret): the operator ensures the `Network` and IP
`Prefix` exist, picks a `MachineClass` either by name (`cloudProfile`) or as
the smallest one satisfying `resources`, boots the given image and can attach
a public `VirtualIP`. SSH access is wired through `SSHKey` references
(`sshkeys.apeirora.eu`), served at runtime by the shared `sshkey-operator`
from the simple-cluster stack.

## How an account gets the provider

```
 Platform Mesh (KCP)                        MSP cluster                        Gardener / IronCore
┌──────────────────────┐   bind   ┌───────────────────────────────┐   ┌───────────────────────────┐
│ account workspace    │─────────►│ bind-watcher annotates the    │   │                           │
│ binds sap.eu (+      │          │ account namespace             │   │  ┌─────────────────────┐  │
│ simple-cluster)      │          │        │                      │   │  │ Shoot / Machine     │  │
└──────────┬───────────┘          │        ▼                      │   │  └──────────▲──────────┘  │
           │                      │ NamespaceReconciler creates   │   │             │             │
           │  sync-agent mirrors  │ sap-gardener-openstack and    │   │             │             │
           ├─────────────────────►│ sap-gardener-ironcore SCPCs   │   │             │             │
           │  sap.eu resources    │        │                      │   │             │             │
           │                      │        ▼                      │   │             │             │
           │                      │ SimpleCluster or *Shoot CR ───┼───┼─────────────┘             │
           │                      │ … kubeconfig Secret written   │   │                           │
           └──────────────────────┤ back next to the request      │   └───────────────────────────┘
                                  └───────────────────────────────┘
```

The `NamespaceReconciler` is the sole creator of the SCPCs: it waits until the
account namespace is annotated with **both** the simple-cluster and the sap.eu
APIExport names, then materialises both flavors and wires in the Gardener
kubeconfig — either one shared admin-managed Secret (`global` mode, the
default) or a per-client kubeconfig delivered from the account's KCP workspace
(`kcp` mode).

## Repository layout

Two Go modules, one UI, five Helm charts:

| Path | What it is |
|---|---|
| `/` (root module) | The provider bundle operator: both Shoot reconcilers + the `NamespaceReconciler` |
| `ironcore-operator/` | The IronCore VM operator (`IronCoreVM` → IronCore `Machine`) |
| `ui/` | "SAP Shoots" portal micro-frontend (Luigi + UI5): CRUD for both Shoot flavors and VMs through the GraphQL gateway |
| `charts/sapprovider-operator` | Deploys the bundle operator |
| `charts/ironcore-operator` | Deploys the IronCore VM operator |
| `charts/sapprovider-ui` | Deploys the portal micro-frontend |
| `charts/sapprovider-pm-integration` | Registers the provider in Platform Mesh: APIExport, APIResourceSchemas, marketplace tile and portal content |
| `charts/sapprovider-sync-agent` | Mirrors `sap.eu` resources between the KCP workspace and the MSP cluster |

> **Schema note:** `charts/sapprovider-pm-integration` carries the `IronCoreVM`
> APIResourceSchema used for the KCP projection; it is generated from
> `ironcore-operator/config/crd/bases`. Bump `vmSchemaRevision` and refresh it
> when that CRD changes.

## IronCore landscape values

The operator requires `scpc.ironcore` in the `sapprovider-operator` chart.
There is no built-in IronCore landscape. Missing settings stop the chart render
or the operator start with a field-specific error.

The chart mounts these values from a ConfigMap and sets `--ironcore-config`.
A checksum restarts the operator when the configuration changes.
For a direct operator start, supply the same configuration file with that flag.

| Value | Purpose |
|---|---|
| `scpc.ironcore.projectNamespace` | Garden project namespace, such as `garden-dev` |
| `scpc.ironcore.autoscaling` | SimpleCluster autoscaling support; neutral default `true` |
| `scpc.ironcore.profiles` | Machine profiles, ordered smallest first; each requires string `name`, `cpu`, and `memory` |
| `scpc.ironcore.shootSpec` | Complete Gardener Shoot spec |
| `shootSpec.cloudProfile` | CloudProfile name and kind |
| `shootSpec.credentialsBindingName` | CredentialsBinding in the Garden project |
| `shootSpec.kubernetes.version` | Kubernetes version |
| `shootSpec.region` | Cloud region |
| `shootSpec.provider.workers[0].machine` | Machine type, architecture, image name and version |
| `shootSpec.provider.workers[0].volume` | Disk size and volume type (storage class) |
| `shootSpec.provider.workers[0].zones` | Worker zones |
| `shootSpec.networking` | Network type, CIDRs, IP families and extension configuration |

All other Shoot settings, including addons, DNS, extensions, maintenance,
provider configuration and worker settings, pass through `shootSpec`.
CloudProfile, region, Kubernetes version and image name remain SimpleCluster inputs.
A single configured zone remains the `zone` input. Multiple zones pass through as a list.
`GARDEN_PROJECT_NAMESPACE` applies only to OpenStack. IronCore uses `projectNamespace`.
The operator supports one worker pool. SimpleCluster controls its minimum and
maximum worker counts. A standalone Shoot with an empty spec uses one worker.
With profiles, SimpleCluster selects the machine type from the requested CPU
and memory. Without profiles, set `machine.type` explicitly.
A standalone Shoot uses the first profile. An explicit provider CR spec still
overrides the configured Shoot spec. The configured Garden project remains the target.
OpenStack configuration is unchanged.

Use the [garden-dev example](charts/sapprovider-operator/examples/garden-dev.yaml)
with Kubernetes 1.33.13, gardenlinux 2150.1.0, a 50Gi disk and CPU profiles 1/2/4/8.
Check the CloudProfile and CredentialsBinding before use. Quota and available
capacity are separate requirements.

```sh
helm template sap-provider charts/sapprovider-operator \
  --namespace sap-provider \
  -f charts/sapprovider-operator/examples/garden-dev.yaml
```

## Installation

Images and charts are published to GHCR at one common version — a chart
`X.Y.Z` always pulls images `X.Y.Z`:

```sh
helm install sap-provider oci://ghcr.io/apeirora/charts/sapprovider-operator --version <X.Y.Z>
```

Every release also publishes an [OCM](https://ocm.software/) component
(`showroom.apeirora.eu/sap-hosting`) at `ghcr.io/apeirora/ocm`, referencing all
images and charts of that version.

Deployment-specific configuration worth knowing:

- **`GARDEN_PROJECT_NAMESPACE`** (env, via the chart's `extraEnv`) — the
  Gardener project namespace standalone Shoots are created in. The embedded
  defaults carry a placeholder; every deployment must set this.
- **`--kubeconfig-source`** — `global` (one shared Gardener kubeconfig Secret,
  see `scpc.kubeconfig.global.*` chart values) or `kcp` (per-client Secret in
  the account namespace).
- The pm-integration and sync-agent charts need the KCP workspace kubeconfig;
  see `charts/sapprovider-pm-integration/values.yaml`.

## Development

```sh
make build vet test          # root module (envtest-backed)
cd ironcore-operator && make build vet
cd ui && npm ci && npm run typecheck && npm run build
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for the pull-request workflow, the SBOM
regeneration step after dependency changes, and what CI requires.

## License

Licensed under the [Apache License 2.0](LICENSE). Copyright 2026 SAP SE or an
SAP affiliate company and showroom-msp-sap-hosting contributors. The repository
is [REUSE](https://reuse.software/) compliant; see [REUSE.toml](REUSE.toml) and
[docs/third-party-licenses.md](docs/third-party-licenses.md).

---

<p align="center">
  <sub>Built with care by <a href="https://apeirora.eu/">ApeiroRA</a></sub>
</p>

<p align="center">
  <a href="https://apeirora.eu/">
    <img src="docs/assets/eu-funded.svg" alt="Funded by the European Union -- NextGenerationEU. Supported by the Federal Ministry for Economic Affairs and Energy on the basis of a decision by the German Bundestag." width="400" />
  </a>
</p>
