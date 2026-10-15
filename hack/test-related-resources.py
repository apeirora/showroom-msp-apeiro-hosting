#!/usr/bin/env python3
# Copyright 2026 ApeiroRA contributors
# SPDX-License-Identifier: Apache-2.0
"""Test rendered resources against the api-syncagent v0.7.0 caller."""
import argparse
import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile

REVISION = "d9e14d998a374d81fb73254069b9be0f3d02cd9f"
ROOT = Path(__file__).resolve().parents[1]

def run(args, **kwargs):
    return subprocess.run(args, check=True, text=True, **kwargs)

parser = argparse.ArgumentParser()
parser.add_argument("--upstream", type=Path)
args = parser.parse_args()
with tempfile.TemporaryDirectory(prefix="hosted-related-") as temporary:
    temp = Path(temporary)
    upstream = args.upstream
    if upstream is None:
        upstream = temp / "api-syncagent"
        run(["git", "clone", "--quiet", "--depth", "1", "--branch", "v0.7.0", "https://github.com/kcp-dev/api-syncagent.git", str(upstream)])
    upstream = upstream.resolve()
    actual = run(["git", "rev-parse", "HEAD"], cwd=upstream, capture_output=True).stdout.strip()
    if actual != REVISION:
        raise SystemExit(f"Expected api-syncagent {REVISION}, got {actual}")
    if run(["git", "status", "--porcelain"], cwd=upstream, capture_output=True).stdout:
        raise SystemExit("The upstream source must be clean.")
    chart_source = next(ROOT.glob("charts/*provider-sync-agent"))
    chart = temp / "chart"
    shutil.copytree(chart_source, chart)
    run(["helm", "dependency", "build", str(chart)])
    fixture = ROOT / "hack/tests/related_namespace_test.go.txt"
    overlay = temp / "overlay.json"
    overlay.write_text(json.dumps({"Replace": {str(upstream / "internal/sync/hosted_related_namespace_test.go"): str(fixture)}}))
    env = dict(os.environ, HOSTED_PROVIDER_LINE="apeiro" if chart_source.name.startswith("apeiro") else "sap")
    failed = False
    for client_enabled in (False, True):
        manifest = temp / f"render-{client_enabled}.yaml"
        with manifest.open("w") as output:
            run(["helm", "template", "related-test", str(chart), "--namespace", "provider-system", "--set", f"publishedResources.clientKubeconfig.enabled={str(client_enabled).lower()}"], stdout=output)
        env["HOSTED_PUBLISHED_RESOURCES"] = str(manifest)
        env["HOSTED_CLIENT_SECRET"] = str(client_enabled).lower()
        result = subprocess.run(["go", "test", "-p=2", "-mod=readonly", "-overlay", str(overlay), "./internal/sync", "-run", "^TestHostedPublishedResourceRelatedNamespace$", "-count=1", "-v"], cwd=upstream, env=env)
        failed = failed or result.returncode != 0
    if run(["git", "status", "--porcelain"], cwd=upstream, capture_output=True).stdout:
        raise SystemExit("The test changed upstream source.")
    if failed:
        raise SystemExit(1)
