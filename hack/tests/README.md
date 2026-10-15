# Related Secret contract

Run `python3 hack/test-related-resources.py` from the repository root.
Use Python 3, Helm 3 and Go 1.26.3 or later.

The test renders the sync agent chart with client kubeconfig sync disabled and enabled.
It runs the actual api-syncagent v0.7.0 resolver at commit `d9e14d998a374d81fb73254069b9be0f3d02cd9f`.
A Go overlay adds the test file without changes to the upstream checkout.
The test uses fake Kubernetes clients with Secret metadata only.

For Shoot resources, the service namespace is the physical primary namespace.
The Secret namespace in each KCP account is `default`.
The reverse client kubeconfig map uses the same two namespaces.
Two accounts use the same Secret name to check namespace isolation.
A failed Shoot with no kubeconfig reference must not block create or delete.
The namespaced VM SSH Secret keeps its existing namespace map.

To reuse a clean checkout at that exact revision, pass `--upstream /path/to/api-syncagent`.
