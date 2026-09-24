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

- **168** distinct third-party dependencies across **3** published components.
- Licences in use: `Apache-2.0` (68), `MIT` (38), `BSD-3-Clause` (35), `Apache-2.0 AND BSD-3-Clause` (16), `Apache-2.0 AND MIT` (4), `Apache-2.0 AND BSD-3-Clause AND MIT` (3), `ISC` (2), `BSD-2-Clause` (1), `Unlicense OR Apache-2.0` (1).
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
| `github.com/ironcore-dev/ironcore` v0.3.0 | `Apache-2.0` |
| `github.com/prometheus/client_golang` v1.22.0 | `Apache-2.0` |
| `github.com/prometheus/client_golang` v1.23.2 | `Apache-2.0` |
| `github.com/prometheus/client_model` v0.6.1 | `Apache-2.0` |
| `github.com/prometheus/client_model` v0.6.2 | `Apache-2.0` |
| `github.com/prometheus/common` v0.62.0 | `Apache-2.0` |
| `github.com/prometheus/common` v0.67.1 | `Apache-2.0` |
| `github.com/prometheus/procfs` v0.15.1 | `Apache-2.0` |
| `github.com/prometheus/procfs` v0.19.1 | `Apache-2.0` |
| `k8s.io/kube-openapi` v0.0.0-20250710124328-f3f2b991d03b | `Apache-2.0` |
| `k8s.io/kube-openapi` v0.0.0-20250910181357-589584f1c912 | `Apache-2.0` |
| `sigs.k8s.io/randfill` v1.0.0 | `Apache-2.0` |

## Full inventory

### Apache-2.0 (68)

| Dependency | Version | Used by |
|---|---|---|
| `@luigi-project/client` | `2.31.0` | ui |
| `@sap-theming/theming-base-content` | `11.36.4` | ui |
| `@ui5/webcomponents` | `2.24.0` | ui |
| `@ui5/webcomponents-base` | `2.24.0` | ui |
| `@ui5/webcomponents-fiori` | `2.24.0` | ui |
| `@ui5/webcomponents-icons` | `2.24.0` | ui |
| `@ui5/webcomponents-icons-business-suite` | `2.24.0` | ui |
| `@ui5/webcomponents-icons-tnt` | `2.24.0` | ui |
| `@ui5/webcomponents-localization` | `2.24.0` | ui |
| `@ui5/webcomponents-react` | `2.24.1` | ui |
| `@ui5/webcomponents-react-base` | `2.24.1` | ui |
| `@ui5/webcomponents-theming` | `2.24.0` | ui |
| `cel.dev/expr` | `v0.25.1` | sap-provider-operator |
| `github.com/go-logr/logr` | `v1.4.3` | ironcore-operator, sap-provider-operator |
| `github.com/go-logr/stdr` | `v1.2.2` | sap-provider-operator |
| `github.com/go-logr/zapr` | `v1.3.0` | ironcore-operator, sap-provider-operator |
| `github.com/go-openapi/jsonpointer` | `v0.21.0` | sap-provider-operator |
| `github.com/go-openapi/jsonpointer` | `v0.22.1` | ironcore-operator |
| `github.com/go-openapi/jsonreference` | `v0.20.2` | sap-provider-operator |
| `github.com/go-openapi/jsonreference` | `v0.21.2` | ironcore-operator |
| `github.com/go-openapi/swag` | `v0.23.0` | sap-provider-operator |
| `github.com/go-openapi/swag` | `v0.25.1` | ironcore-operator |
| `github.com/go-openapi/swag/cmdutils` | `v0.25.1` | ironcore-operator |
| `github.com/go-openapi/swag/conv` | `v0.25.1` | ironcore-operator |
| `github.com/go-openapi/swag/fileutils` | `v0.25.1` | ironcore-operator |
| `github.com/go-openapi/swag/jsonname` | `v0.25.1` | ironcore-operator |
| `github.com/go-openapi/swag/jsonutils` | `v0.25.1` | ironcore-operator |
| `github.com/go-openapi/swag/loading` | `v0.25.1` | ironcore-operator |
| `github.com/go-openapi/swag/mangling` | `v0.25.1` | ironcore-operator |
| `github.com/go-openapi/swag/netutils` | `v0.25.1` | ironcore-operator |
| `github.com/go-openapi/swag/stringutils` | `v0.25.1` | ironcore-operator |
| `github.com/go-openapi/swag/typeutils` | `v0.25.1` | ironcore-operator |
| `github.com/go-openapi/swag/yamlutils` | `v0.25.1` | ironcore-operator |
| `github.com/google/btree` | `v1.1.3` | ironcore-operator, sap-provider-operator |
| `github.com/google/gnostic-models` | `v0.7.0` | ironcore-operator, sap-provider-operator |
| `github.com/inconshreveable/mousetrap` | `v1.1.0` | sap-provider-operator |
| `github.com/ironcore-dev/ironcore` | `v0.3.0` | ironcore-operator |
| `github.com/modern-go/concurrent` | `v0.0.0-20180306012644-bacd9c7ef1dd` | ironcore-operator, sap-provider-operator |
| `github.com/modern-go/reflect2` | `v1.0.3-0.20250322232337-35a7c28c31ee` | ironcore-operator, sap-provider-operator |
| `github.com/prometheus/client_model` | `v0.6.1` | sap-provider-operator |
| `github.com/prometheus/client_model` | `v0.6.2` | ironcore-operator |
| `github.com/prometheus/common` | `v0.62.0` | sap-provider-operator |
| `github.com/prometheus/common` | `v0.67.1` | ironcore-operator |
| `github.com/prometheus/procfs` | `v0.15.1` | sap-provider-operator |
| `github.com/prometheus/procfs` | `v0.19.1` | ironcore-operator |
| `github.com/spf13/cobra` | `v1.9.1` | sap-provider-operator |
| `go.opentelemetry.io/auto/sdk` | `v1.2.1` | sap-provider-operator |
| `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp` | `v0.58.0` | sap-provider-operator |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace` | `v1.34.0` | sap-provider-operator |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc` | `v1.34.0` | sap-provider-operator |
| `go.opentelemetry.io/proto/otlp` | `v1.5.0` | sap-provider-operator |
| `gomodules.xyz/jsonpatch/v2` | `v2.4.0` | sap-provider-operator |
| `gomodules.xyz/jsonpatch/v2` | `v2.5.0` | ironcore-operator |
| `google.golang.org/genproto/googleapis/api` | `v0.0.0-20260414002931-afd174a4e478` | sap-provider-operator |
| `google.golang.org/genproto/googleapis/rpc` | `v0.0.0-20260414002931-afd174a4e478` | sap-provider-operator |
| `google.golang.org/grpc` | `v1.82.1` | sap-provider-operator |
| `k8s.io/api` | `v0.34.0` | sap-provider-operator |
| `k8s.io/api` | `v0.34.1` | ironcore-operator |
| `k8s.io/apiextensions-apiserver` | `v0.34.0` | sap-provider-operator |
| `k8s.io/apiextensions-apiserver` | `v0.34.1` | ironcore-operator |
| `k8s.io/apiserver` | `v0.34.0` | sap-provider-operator |
| `k8s.io/component-base` | `v0.34.0` | sap-provider-operator |
| `k8s.io/klog/v2` | `v2.130.1` | ironcore-operator, sap-provider-operator |
| `sigs.k8s.io/apiserver-network-proxy/konnectivity-client` | `v0.31.2` | sap-provider-operator |
| `sigs.k8s.io/controller-runtime` | `v0.22.1` | sap-provider-operator |
| `sigs.k8s.io/controller-runtime` | `v0.22.3` | ironcore-operator |
| `sigs.k8s.io/randfill` | `v1.0.0` | ironcore-operator, sap-provider-operator |
| `sigs.k8s.io/structured-merge-diff/v6` | `v6.3.0` | ironcore-operator, sap-provider-operator |

### MIT (38)

| Dependency | Version | Used by |
|---|---|---|
| `@tanstack/react-virtual` | `3.14.5` | ui |
| `@tanstack/virtual-core` | `3.17.3` | ui |
| `@types/jquery` | `3.5.34` | ui |
| `@types/openui5` | `1.150.0` | ui |
| `@types/prop-types` | `15.7.15` | ui |
| `@types/qunit` | `2.19.14` | ui |
| `@types/react` | `18.3.31` | ui |
| `@types/react-dom` | `18.3.7` | ui |
| `@types/sizzle` | `2.3.10` | ui |
| `@types/trusted-types` | `2.0.7` | ui |
| `@zxing/library` | `0.21.3` | ui |
| `clsx` | `2.1.1` | ui |
| `csstype` | `3.2.3` | ui |
| `github.com/beorn7/perks` | `v1.0.1` | ironcore-operator, sap-provider-operator |
| `github.com/blang/semver/v4` | `v4.0.0` | sap-provider-operator |
| `github.com/cenkalti/backoff/v4` | `v4.3.0` | sap-provider-operator |
| `github.com/cespare/xxhash/v2` | `v2.3.0` | ironcore-operator, sap-provider-operator |
| `github.com/emicklei/go-restful/v3` | `v3.12.2` | sap-provider-operator |
| `github.com/emicklei/go-restful/v3` | `v3.13.0` | ironcore-operator |
| `github.com/felixge/httpsnoop` | `v1.0.4` | sap-provider-operator |
| `github.com/fxamacker/cbor/v2` | `v2.9.0` | ironcore-operator, sap-provider-operator |
| `github.com/go-task/slim-sprig/v3` | `v3.0.0` | sap-provider-operator |
| `github.com/josharian/intern` | `v1.0.0` | sap-provider-operator |
| `github.com/json-iterator/go` | `v1.1.12` | ironcore-operator, sap-provider-operator |
| `github.com/mailru/easyjson` | `v0.7.7` | sap-provider-operator |
| `github.com/onsi/ginkgo/v2` | `v2.22.0` | sap-provider-operator |
| `github.com/onsi/gomega` | `v1.36.1` | sap-provider-operator |
| `github.com/x448/float16` | `v0.8.4` | ironcore-operator, sap-provider-operator |
| `go.uber.org/multierr` | `v1.11.0` | ironcore-operator, sap-provider-operator |
| `go.uber.org/zap` | `v1.27.0` | ironcore-operator, sap-provider-operator |
| `graphql` | `16.14.2` | ui |
| `graphql-sse` | `2.6.0` | ui |
| `js-tokens` | `4.0.0` | ui |
| `loose-envify` | `1.4.0` | ui |
| `react` | `18.3.1` | ui |
| `react-dom` | `18.3.1` | ui |
| `scheduler` | `0.23.2` | ui |
| `ts-custom-error` | `3.3.1` | ui |

### BSD-3-Clause (35)

| Dependency | Version | Used by |
|---|---|---|
| `@lit-labs/ssr-dom-shim` | `1.6.0` | ui |
| `github.com/antlr4-go/antlr/v4` | `v4.13.1` | sap-provider-operator |
| `github.com/evanphx/json-patch/v5` | `v5.9.11` | ironcore-operator, sap-provider-operator |
| `github.com/fsnotify/fsnotify` | `v1.9.0` | ironcore-operator, sap-provider-operator |
| `github.com/gogo/protobuf` | `v1.3.2` | ironcore-operator, sap-provider-operator |
| `github.com/google/go-cmp` | `v0.7.0` | ironcore-operator, sap-provider-operator |
| `github.com/google/uuid` | `v1.6.0` | ironcore-operator, sap-provider-operator |
| `github.com/grpc-ecosystem/grpc-gateway/v2` | `v2.26.3` | sap-provider-operator |
| `github.com/munnerz/goautoneg` | `v0.0.0-20191010083416-a7dc8b61c822` | ironcore-operator, sap-provider-operator |
| `github.com/pmezard/go-difflib` | `v1.0.0` | sap-provider-operator |
| `github.com/pmezard/go-difflib` | `v1.0.1-0.20181226105442-5d4384ee4fb2` | ironcore-operator |
| `github.com/spf13/pflag` | `v1.0.10` | ironcore-operator |
| `github.com/spf13/pflag` | `v1.0.6` | sap-provider-operator |
| `go4.org/netipx` | `v0.0.0-20231129151722-fdeea329fbba` | ironcore-operator |
| `golang.org/x/crypto` | `v0.51.0` | ironcore-operator |
| `golang.org/x/exp` | `v0.0.0-20240823005443-9b4947da3948` | sap-provider-operator |
| `golang.org/x/net` | `v0.53.0` | sap-provider-operator |
| `golang.org/x/net` | `v0.55.0` | ironcore-operator |
| `golang.org/x/oauth2` | `v0.34.0` | ironcore-operator |
| `golang.org/x/oauth2` | `v0.36.0` | sap-provider-operator |
| `golang.org/x/sync` | `v0.20.0` | ironcore-operator, sap-provider-operator |
| `golang.org/x/sys` | `v0.43.0` | sap-provider-operator |
| `golang.org/x/sys` | `v0.45.0` | ironcore-operator |
| `golang.org/x/term` | `v0.42.0` | sap-provider-operator |
| `golang.org/x/term` | `v0.43.0` | ironcore-operator |
| `golang.org/x/text` | `v0.36.0` | sap-provider-operator |
| `golang.org/x/text` | `v0.37.0` | ironcore-operator |
| `golang.org/x/time` | `v0.14.0` | ironcore-operator |
| `golang.org/x/time` | `v0.9.0` | sap-provider-operator |
| `golang.org/x/tools` | `v0.43.0` | sap-provider-operator |
| `google.golang.org/protobuf` | `v1.36.11` | ironcore-operator, sap-provider-operator |
| `gopkg.in/evanphx/json-patch.v4` | `v4.12.0` | sap-provider-operator |
| `gopkg.in/evanphx/json-patch.v4` | `v4.13.0` | ironcore-operator |
| `gopkg.in/inf.v0` | `v0.9.1` | ironcore-operator, sap-provider-operator |
| `lit-html` | `2.8.0` | ui |

### Apache-2.0 AND BSD-3-Clause (16)

| Dependency | Version | Used by |
|---|---|---|
| `github.com/google/cel-go` | `v0.29.0` | sap-provider-operator |
| `github.com/google/pprof` | `v0.0.0-20241029153458-d1b30febd7db` | sap-provider-operator |
| `github.com/prometheus/client_golang` | `v1.22.0` | sap-provider-operator |
| `github.com/prometheus/client_golang` | `v1.23.2` | ironcore-operator |
| `go.opentelemetry.io/otel` | `v1.43.0` | sap-provider-operator |
| `go.opentelemetry.io/otel/metric` | `v1.43.0` | sap-provider-operator |
| `go.opentelemetry.io/otel/sdk` | `v1.43.0` | sap-provider-operator |
| `go.opentelemetry.io/otel/trace` | `v1.43.0` | sap-provider-operator |
| `k8s.io/apimachinery` | `v0.34.0` | sap-provider-operator |
| `k8s.io/apimachinery` | `v0.34.1` | ironcore-operator |
| `k8s.io/client-go` | `v0.34.0` | sap-provider-operator |
| `k8s.io/client-go` | `v0.34.1` | ironcore-operator |
| `k8s.io/utils` | `v0.0.0-20250604170112-4c0f3b243397` | sap-provider-operator |
| `k8s.io/utils` | `v0.0.0-20251002143259-bc988d571ff4` | ironcore-operator |
| `sigs.k8s.io/json` | `v0.0.0-20241014173422-cfa47c3a1cc8` | sap-provider-operator |
| `sigs.k8s.io/json` | `v0.0.0-20250730193827-2d320260d730` | ironcore-operator |

### Apache-2.0 AND MIT (4)

| Dependency | Version | Used by |
|---|---|---|
| `go.yaml.in/yaml/v2` | `v2.4.2` | sap-provider-operator |
| `go.yaml.in/yaml/v2` | `v2.4.3` | ironcore-operator |
| `go.yaml.in/yaml/v3` | `v3.0.4` | ironcore-operator, sap-provider-operator |
| `gopkg.in/yaml.v3` | `v3.0.1` | sap-provider-operator |

### Apache-2.0 AND BSD-3-Clause AND MIT (3)

| Dependency | Version | Used by |
|---|---|---|
| `k8s.io/kube-openapi` | `v0.0.0-20250710124328-f3f2b991d03b` | sap-provider-operator |
| `k8s.io/kube-openapi` | `v0.0.0-20250910181357-589584f1c912` | ironcore-operator |
| `sigs.k8s.io/yaml` | `v1.6.0` | ironcore-operator, sap-provider-operator |

### ISC (2)

| Dependency | Version | Used by |
|---|---|---|
| `github.com/davecgh/go-spew` | `v1.1.1` | sap-provider-operator |
| `github.com/davecgh/go-spew` | `v1.1.2-0.20180830191138-d8f796af33cc` | ironcore-operator |

### BSD-2-Clause (1)

| Dependency | Version | Used by |
|---|---|---|
| `github.com/pkg/errors` | `v0.9.1` | sap-provider-operator |

### Unlicense OR Apache-2.0 (1)

| Dependency | Version | Used by |
|---|---|---|
| `@zxing/text-encoding` | `0.9.0` | ui |
