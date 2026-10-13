# Secret reference configuration

The operator checks each Secret reference before it uses the Secret.
An omitted reference namespace means the request namespace.
A reference to another namespace requires an exact operator exception.

Set `secretReferenceAllowlist` in the operator chart values:

```yaml
secretReferenceAllowlist:
  - namespace: credential-system
    name: shared-key
    purpose: use
```

The chart passes the list through `--secret-reference-allowlist` as JSON.
Each entry requires an explicit namespace, name, and purpose.
Wildcards and empty entry namespaces are invalid.

| Purpose | Allowed operation |
| --- | --- |
| `use` | Use the named credential |
| `copy` | Copy data from the named Secret |
| `write` | Create or change the named Secret |

An exception for one purpose does not permit another purpose.
For example, `use` does not permit a copy of the credential.
The default list is empty.
The controller checks still apply if Kubernetes RBAC permits broader access.

The configured global Gardener Secret is an exact `use` exception.
The operator derives its name and namespace from its startup flags.

The SAP chart derives remote exceptions from `scpc.ironcore.projectNamespace`
and `scpc.ironcore.shootSpec`.
These fields are operator configuration.
Tenant Shoot text cannot add exceptions.

| Remote purpose | Configured field |
| --- | --- |
| `gardener-credential-binding` | `credentialsBindingName` |
| `gardener-secret-binding` | `secretBindingName` |
| `gardener-dns` | DNS provider `secretName` or `credentialsRef.name` |

A DNS `credentialsRef.namespace` defaults to the configured project namespace.
An explicit namespace remains part of the exact exception.
