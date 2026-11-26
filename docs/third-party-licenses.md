# Third-party licences

Every direct and transitive dependency of every component published from this
repository, with the licence it is distributed under.

This document is generated — do not edit it by hand:

```sh
./hack/generate-sbom.sh                     # refresh sbom/*.json
python3 hack/generate-third-party-licenses.py
```

The source of truth is the CycloneDX SBOM set under [`sbom/`](../sbom). Go
dependencies come from each module's `go.mod`; UI dependencies come from
`ui/package-lock.json`, production scope only, because dev-only packages are
build-time tooling and are not part of any shipped artefact.

## Summary

- **144** distinct third-party dependencies across **3** published components.
- Licences in use: `Apache-2.0` (57), `MIT` (33), `BSD-3-Clause` (33), `Apache-2.0 AND BSD-3-Clause` (14), `Apache-2.0 AND BSD-3-Clause AND MIT` (3), `Apache-2.0 AND MIT` (2), `ISC` (1), `Unlicense OR Apache-2.0` (1).
- **No copyleft dependencies.** Every licence above is permissive (Apache-2.0, BSD, MIT, ISC), so redistribution under Apache-2.0 carries no source-disclosure obligation.

## Obligations

| Licence | Obligation when redistributing | How this repository satisfies it |
|---|---|---|
| Apache-2.0 | Retain copyright, licence and NOTICE; state modifications | Dependencies are consumed unmodified as Go modules / npm packages; their licence headers ship inside the module and inside the container image layers. This repository's own licence is [Apache-2.0](../LICENSE). |
| BSD-2-Clause, BSD-3-Clause | Retain copyright notice and disclaimer; do not use the author's name to endorse | Notices ship unmodified with each dependency; no endorsement is claimed. |
| MIT, ISC | Retain copyright notice and permission notice | Notices ship unmodified with each dependency. |
| Unlicense OR Apache-2.0 | Dual-licensed; this project takes the Apache-2.0 branch | Treated as Apache-2.0 above. |

## Unmatched licence texts

syft could not match these projects' `LICENSE` files to an SPDX template
verbatim (added preamble, reformatting). Each was resolved by reading the file
in the Go module cache and is recorded in `OVERRIDES` in
`hack/generate-third-party-licenses.py`:

| Dependency | Resolved licence |
|---|---|
| `github.com/ironcore-dev/ironcore` v0.9.0 | `Apache-2.0` |
| `github.com/prometheus/client_golang` v1.24.0 | `Apache-2.0` |
| `github.com/prometheus/client_model` v0.6.2 | `Apache-2.0` |
| `github.com/prometheus/common` v0.70.0 | `Apache-2.0` |
| `github.com/prometheus/procfs` v0.21.1 | `Apache-2.0` |
| `k8s.io/kube-openapi` v0.0.0-20260721132016-d427ff9ee9ad | `Apache-2.0` |
| `sigs.k8s.io/randfill` v1.0.0 | `Apache-2.0` |

## Full inventory

### Apache-2.0 (57)

| Dependency | Version | Used by |
|---|---|---|
| `@luigi-project/client` | `2.33.0` | ui |
| `@sap-theming/theming-base-content` | `11.36.4` | ui |
| `@ui5/webcomponents` | `2.29.0` | ui |
| `@ui5/webcomponents-base` | `2.29.0` | ui |
| `@ui5/webcomponents-fiori` | `2.29.0` | ui |
| `@ui5/webcomponents-icons` | `2.29.0` | ui |
| `@ui5/webcomponents-icons-business-suite` | `2.29.0` | ui |
| `@ui5/webcomponents-icons-tnt` | `2.29.0` | ui |
| `@ui5/webcomponents-localization` | `2.29.0` | ui |
| `@ui5/webcomponents-react` | `2.29.0` | ui |
| `@ui5/webcomponents-react-base` | `2.29.0` | ui |
| `@ui5/webcomponents-theming` | `2.29.0` | ui |
| `cel.dev/expr` | `v0.25.2` | apeiro-provider-operator |
| `github.com/go-logr/logr` | `v1.4.4` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-logr/stdr` | `v1.2.2` | apeiro-provider-operator |
| `github.com/go-logr/zapr` | `v1.3.0` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/jsonpointer` | `v1.0.1` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/jsonreference` | `v1.0.0` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag/cmdutils` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag/conv` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag/fileutils` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag/jsonutils` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag/loading` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag/mangling` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag/netutils` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag/pools` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag/stringutils` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag/typeutils` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-openapi/swag/yamlutils` | `v0.29.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/google/btree` | `v1.1.3` | apeiro-provider-operator |
| `github.com/google/gnostic-models` | `v0.7.0` | apeiro-provider-operator, ironcore-operator |
| `github.com/inconshreveable/mousetrap` | `v1.1.0` | apeiro-provider-operator |
| `github.com/ironcore-dev/ironcore` | `v0.9.0` | ironcore-operator |
| `github.com/modern-go/concurrent` | `v0.0.0-20180306012644-bacd9c7ef1dd` | apeiro-provider-operator, ironcore-operator |
| `github.com/modern-go/reflect2` | `v1.0.3-0.20250322232337-35a7c28c31ee` | apeiro-provider-operator, ironcore-operator |
| `github.com/prometheus/client_model` | `v0.6.2` | apeiro-provider-operator, ironcore-operator |
| `github.com/prometheus/common` | `v0.70.0` | apeiro-provider-operator, ironcore-operator |
| `github.com/prometheus/procfs` | `v0.21.1` | apeiro-provider-operator, ironcore-operator |
| `github.com/spf13/cobra` | `v1.10.2` | apeiro-provider-operator |
| `go.opentelemetry.io/auto/sdk` | `v1.2.1` | apeiro-provider-operator |
| `go.opentelemetry.io/proto/otlp` | `v1.11.0` | apeiro-provider-operator |
| `gomodules.xyz/jsonpatch/v2` | `v2.4.0` | apeiro-provider-operator |
| `gomodules.xyz/jsonpatch/v2` | `v2.5.0` | ironcore-operator |
| `google.golang.org/genproto/googleapis/api` | `v0.0.0-20260803160001-6ac0973c030d` | apeiro-provider-operator |
| `google.golang.org/genproto/googleapis/rpc` | `v0.0.0-20260803160001-6ac0973c030d` | apeiro-provider-operator |
| `google.golang.org/grpc` | `v1.83.1` | apeiro-provider-operator |
| `k8s.io/api` | `v0.37.2` | apeiro-provider-operator, ironcore-operator |
| `k8s.io/apiextensions-apiserver` | `v0.37.0` | apeiro-provider-operator, ironcore-operator |
| `k8s.io/apiserver` | `v0.37.0` | apeiro-provider-operator |
| `k8s.io/component-base` | `v0.37.0` | apeiro-provider-operator |
| `k8s.io/klog/v2` | `v2.140.0` | apeiro-provider-operator, ironcore-operator |
| `k8s.io/streaming` | `v0.37.2` | apeiro-provider-operator |
| `sigs.k8s.io/apiserver-network-proxy/konnectivity-client` | `v0.36.0` | apeiro-provider-operator |
| `sigs.k8s.io/controller-runtime` | `v0.25.3` | apeiro-provider-operator, ironcore-operator |
| `sigs.k8s.io/randfill` | `v1.0.0` | apeiro-provider-operator, ironcore-operator |
| `sigs.k8s.io/structured-merge-diff/v6` | `v6.4.2` | apeiro-provider-operator, ironcore-operator |

### BSD-3-Clause (33)

| Dependency | Version | Used by |
|---|---|---|
| `@lit-labs/ssr-dom-shim` | `1.6.0` | ui |
| `github.com/antlr4-go/antlr/v4` | `v4.13.1` | apeiro-provider-operator |
| `github.com/evanphx/json-patch/v5` | `v5.9.11` | apeiro-provider-operator, ironcore-operator |
| `github.com/fsnotify/fsnotify` | `v1.9.0` | apeiro-provider-operator, ironcore-operator |
| `github.com/gogo/protobuf` | `v1.3.2` | apeiro-provider-operator |
| `github.com/google/go-cmp` | `v0.7.0` | apeiro-provider-operator |
| `github.com/google/uuid` | `v1.6.0` | apeiro-provider-operator, ironcore-operator |
| `github.com/grpc-ecosystem/grpc-gateway/v2` | `v2.29.0` | apeiro-provider-operator |
| `github.com/munnerz/goautoneg` | `v0.0.0-20191010083416-a7dc8b61c822` | apeiro-provider-operator, ironcore-operator |
| `github.com/pmezard/go-difflib` | `v1.0.1-0.20181226105442-5d4384ee4fb2` | apeiro-provider-operator, ironcore-operator |
| `github.com/spf13/pflag` | `v1.0.10` | apeiro-provider-operator, ironcore-operator |
| `go4.org/netipx` | `v0.0.0-20231129151722-fdeea329fbba` | ironcore-operator |
| `golang.org/x/crypto` | `v0.59.0` | ironcore-operator |
| `golang.org/x/exp` | `v0.0.0-20260410095643-746e56fc9e2f` | apeiro-provider-operator |
| `golang.org/x/mod` | `v0.41.0` | apeiro-provider-operator |
| `golang.org/x/net` | `v0.60.0` | apeiro-provider-operator |
| `golang.org/x/net` | `v0.62.0` | ironcore-operator |
| `golang.org/x/oauth2` | `v0.36.0` | apeiro-provider-operator, ironcore-operator |
| `golang.org/x/sync` | `v0.23.0` | apeiro-provider-operator |
| `golang.org/x/sync` | `v0.25.0` | ironcore-operator |
| `golang.org/x/sys` | `v0.48.0` | apeiro-provider-operator |
| `golang.org/x/sys` | `v0.50.0` | ironcore-operator |
| `golang.org/x/term` | `v0.46.0` | apeiro-provider-operator |
| `golang.org/x/term` | `v0.48.0` | ironcore-operator |
| `golang.org/x/text` | `v0.42.0` | apeiro-provider-operator |
| `golang.org/x/text` | `v0.44.0` | ironcore-operator |
| `golang.org/x/time` | `v0.15.0` | apeiro-provider-operator, ironcore-operator |
| `golang.org/x/tools` | `v0.49.0` | apeiro-provider-operator |
| `google.golang.org/protobuf` | `v1.36.12` | ironcore-operator |
| `google.golang.org/protobuf` | `v1.36.12-0.20260120151049-f2248ac996af` | apeiro-provider-operator |
| `gopkg.in/evanphx/json-patch.v4` | `v4.13.0` | apeiro-provider-operator, ironcore-operator |
| `gopkg.in/inf.v0` | `v0.9.1` | apeiro-provider-operator, ironcore-operator |
| `lit-html` | `2.8.0` | ui |

### MIT (33)

| Dependency | Version | Used by |
|---|---|---|
| `@tanstack/react-virtual` | `3.14.26` | ui |
| `@tanstack/virtual-core` | `3.19.11` | ui |
| `@types/jquery` | `3.5.34` | ui |
| `@types/openui5` | `1.154.0` | ui |
| `@types/qunit` | `2.19.14` | ui |
| `@types/react` | `19.3.3` | ui |
| `@types/react-dom` | `19.3.0` | ui |
| `@types/sizzle` | `2.3.10` | ui |
| `@types/trusted-types` | `2.0.7` | ui |
| `@zxing/library` | `0.21.3` | ui |
| `clsx` | `2.1.1` | ui |
| `csstype` | `3.2.3` | ui |
| `github.com/Masterminds/semver/v3` | `v3.4.0` | apeiro-provider-operator |
| `github.com/beorn7/perks` | `v1.0.1` | apeiro-provider-operator, ironcore-operator |
| `github.com/blang/semver/v4` | `v4.0.0` | apeiro-provider-operator |
| `github.com/cenkalti/backoff/v5` | `v5.0.3` | apeiro-provider-operator |
| `github.com/cespare/xxhash/v2` | `v2.3.0` | apeiro-provider-operator, ironcore-operator |
| `github.com/emicklei/go-restful/v3` | `v3.13.0` | apeiro-provider-operator, ironcore-operator |
| `github.com/felixge/httpsnoop` | `v1.0.4` | apeiro-provider-operator |
| `github.com/fxamacker/cbor/v2` | `v2.9.1` | apeiro-provider-operator, ironcore-operator |
| `github.com/go-task/slim-sprig/v3` | `v3.0.0` | apeiro-provider-operator |
| `github.com/json-iterator/go` | `v1.1.12` | apeiro-provider-operator, ironcore-operator |
| `github.com/onsi/ginkgo/v2` | `v2.34.0` | apeiro-provider-operator |
| `github.com/onsi/gomega` | `v1.44.0` | apeiro-provider-operator |
| `github.com/x448/float16` | `v0.8.4` | apeiro-provider-operator, ironcore-operator |
| `go.uber.org/multierr` | `v1.11.0` | apeiro-provider-operator, ironcore-operator |
| `go.uber.org/zap` | `v1.27.1` | apeiro-provider-operator, ironcore-operator |
| `graphql` | `16.14.2` | ui |
| `graphql-sse` | `2.6.4` | ui |
| `react` | `19.3.1` | ui |
| `react-dom` | `19.3.1` | ui |
| `scheduler` | `0.28.0` | ui |
| `ts-custom-error` | `3.3.1` | ui |

### Apache-2.0 AND BSD-3-Clause (14)

| Dependency | Version | Used by |
|---|---|---|
| `github.com/google/cel-go` | `v0.29.2` | apeiro-provider-operator |
| `github.com/google/pprof` | `v0.0.0-20260402051712-545e8a4df936` | apeiro-provider-operator |
| `github.com/prometheus/client_golang` | `v1.24.0` | apeiro-provider-operator, ironcore-operator |
| `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` | `v0.69.0` | apeiro-provider-operator |
| `go.opentelemetry.io/otel` | `v1.47.0` | apeiro-provider-operator |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace` | `v1.45.0` | apeiro-provider-operator |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc` | `v1.45.0` | apeiro-provider-operator |
| `go.opentelemetry.io/otel/log` | `v1.47.0` | apeiro-provider-operator |
| `go.opentelemetry.io/otel/metric` | `v1.47.0` | apeiro-provider-operator |
| `go.opentelemetry.io/otel/sdk` | `v1.47.0` | apeiro-provider-operator |
| `go.opentelemetry.io/otel/trace` | `v1.47.0` | apeiro-provider-operator |
| `k8s.io/apimachinery` | `v0.37.2` | apeiro-provider-operator, ironcore-operator |
| `k8s.io/utils` | `v0.0.0-20260626114624-be93311217bd` | apeiro-provider-operator, ironcore-operator |
| `sigs.k8s.io/json` | `v0.0.0-20250730193827-2d320260d730` | apeiro-provider-operator, ironcore-operator |

### Apache-2.0 AND BSD-3-Clause AND MIT (3)

| Dependency | Version | Used by |
|---|---|---|
| `k8s.io/client-go` | `v0.37.2` | apeiro-provider-operator, ironcore-operator |
| `k8s.io/kube-openapi` | `v0.0.0-20260721132016-d427ff9ee9ad` | apeiro-provider-operator, ironcore-operator |
| `sigs.k8s.io/yaml` | `v1.6.0` | apeiro-provider-operator, ironcore-operator |

### Apache-2.0 AND MIT (2)

| Dependency | Version | Used by |
|---|---|---|
| `go.yaml.in/yaml/v2` | `v2.4.4` | apeiro-provider-operator, ironcore-operator |
| `go.yaml.in/yaml/v3` | `v3.0.5` | apeiro-provider-operator, ironcore-operator |

### ISC (1)

| Dependency | Version | Used by |
|---|---|---|
| `github.com/davecgh/go-spew` | `v1.1.2-0.20180830191138-d8f796af33cc` | apeiro-provider-operator, ironcore-operator |

### Unlicense OR Apache-2.0 (1)

| Dependency | Version | Used by |
|---|---|---|
| `@zxing/text-encoding` | `0.9.0` | ui |
