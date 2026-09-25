#!/usr/bin/env python3

import argparse
import json
import sys
from pathlib import Path


EXPECTED_KEYS = {
    "base_ref",
    "compatible_changes",
    "incompatible_changes",
    "module",
    "schema_version",
    "tool_module",
    "tool_version",
}


def fail(message):
    print("api compatibility: " + message, file=sys.stderr)
    raise SystemExit(1)


def load_allowlist(path):
    try:
        value = json.loads(Path(path).read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        fail("cannot read allowlist: {}".format(error))
    if not isinstance(value, dict) or set(value) != EXPECTED_KEYS:
        fail("allowlist keys are invalid")
    if value["schema_version"] != 1:
        fail("allowlist schema_version must be 1")
    for key in ("base_ref", "module", "tool_module", "tool_version"):
        if not isinstance(value[key], str) or not value[key]:
            fail("allowlist {} must be a non-empty string".format(key))
    for key in ("compatible_changes", "incompatible_changes"):
        changes = value[key]
        if (
            not isinstance(changes, list)
            or any(not isinstance(change, str) or not change for change in changes)
            or changes != sorted(changes)
            or len(changes) != len(set(changes))
        ):
            fail("allowlist {} must be sorted and unique".format(key))
    return value


def parse_report(path):
    sections = {"Compatible changes:": [], "Incompatible changes:": []}
    current = None
    try:
        lines = Path(path).read_text(encoding="utf-8").splitlines()
    except OSError as error:
        fail("cannot read apidiff report: {}".format(error))
    for line in lines:
        if not line:
            continue
        if line in sections:
            current = line
            continue
        if current is None or not line.startswith("- ") or len(line) == 2:
            fail("unrecognized apidiff output: {}".format(line))
        sections[current].append(line[2:])
    for changes in sections.values():
        if len(changes) != len(set(changes)):
            fail("apidiff report contains duplicate changes")
        changes.sort()
    return sections


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--allowlist", required=True)
    parser.add_argument("--report", required=True)
    parser.add_argument("--output", required=True)
    parser.add_argument("--base-ref", required=True)
    parser.add_argument("--module", required=True)
    parser.add_argument("--tool-version", required=True)
    args = parser.parse_args()

    allowed = load_allowlist(args.allowlist)
    expected_metadata = {
        "base_ref": args.base_ref,
        "module": args.module,
        "tool_module": "golang.org/x/exp",
        "tool_version": args.tool_version,
    }
    for key, expected in expected_metadata.items():
        if allowed[key] != expected:
            fail("allowlist {} does not match {}".format(key, expected))

    sections = parse_report(args.report)
    actual = dict(expected_metadata)
    actual.update(
        {
            "schema_version": 1,
            "compatible_changes": sections["Compatible changes:"],
            "incompatible_changes": sections["Incompatible changes:"],
        }
    )
    mismatches = []
    for key in ("compatible_changes", "incompatible_changes"):
        if actual[key] != allowed[key]:
            mismatches.append(
                "{} expected {!r}, got {!r}".format(key, allowed[key], actual[key])
            )
    if mismatches:
        fail("; ".join(mismatches))

    output = Path(args.output)
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(json.dumps(actual, indent=2, sort_keys=True) + "\n", encoding="utf-8")


if __name__ == "__main__":
    main()
