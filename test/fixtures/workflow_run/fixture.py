#!/usr/bin/env python3
"""Expose the native Go fixture builder to the plan-owned black-box harness."""
import argparse
import json
import os
import subprocess
import sys

p = argparse.ArgumentParser(description=__doc__)
p.add_argument("--builder", required=True)
p.add_argument("--root", required=True)
p.add_argument("--ply", required=True)
a = p.parse_args()
env = dict(os.environ, PLY_WORKFLOW_FIXTURE_ROOT=a.root,
           PLY_WORKFLOW_FIXTURE_BINARY=a.ply, PLY_WORKFLOW_PYTHON=sys.executable)
r = subprocess.run([a.builder, "-test.run=^TestWorkflowFixtureExport$", "-test.v"],
                   capture_output=True, text=True, env=env)
if r.returncode:
    sys.stderr.write(r.stdout + r.stderr)
    raise SystemExit(r.returncode)
rows = [line.removeprefix("WH01_FIXTURE:") for line in r.stdout.splitlines()
        if line.startswith("WH01_FIXTURE:")]
if len(rows) != 1:
    raise SystemExit("Native fixture did not return exactly one manifest: " + r.stdout)
print(json.dumps(json.loads(rows[0])))
