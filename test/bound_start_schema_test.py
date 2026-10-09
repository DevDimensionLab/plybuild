#!/usr/bin/env python3
"""Validate actual public bound-start CLI responses, plus strict request mutations."""
import argparse
import copy
import json
from pathlib import Path

from jsonschema import Draft202012Validator

parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument("--responses", type=Path, required=True)
args = parser.parse_args()
root = Path(__file__).resolve().parents[1]
schema = json.loads((root / "schemas/workflow-bound-start-v1.schema.json").read_text())
cap_schema = json.loads((root / "schemas/read-contract/ply-capabilities-v1.schema.json").read_text())
Draft202012Validator.check_schema(schema)
Draft202012Validator.check_schema(cap_schema)
bound = Draft202012Validator(schema)
capability = Draft202012Validator(cap_schema)
count = 0
for path in sorted(args.responses.glob("*.json")):
    value = json.loads(path.read_text())
    if value.get("kind", "").startswith("ply.workflow.bound-start-"):
        bound.validate(value)
        count += 1
    elif value.get("kind") == "PlyCapabilities@1":
        capability.validate(value)
        count += 1
assert count >= 8, f"missing actual public CLI evidence: {count}"
request = json.loads((args.responses / "request.json").read_text())
for path in [(), ("input",), ("input", "target")]:
    node = request
    for key in path:
        node = node[key]
    for field in node:
        changed = copy.deepcopy(request)
        at = changed
        for key in path:
            at = at[key]
        del at[field]
        assert not bound.is_valid(changed), (path, field, "missing field accepted")
        at[field] = None
        assert not bound.is_valid(changed), (path, field, "null accepted")
    changed = copy.deepcopy(request)
    at = changed
    for key in path:
        at = at[key]
    at["unexpected"] = True
    assert not bound.is_valid(changed), (path, "unknown field accepted")
print(f"Validated {count} actual CLI JSON responses and every request field's missing/null/unknown rejection.")
