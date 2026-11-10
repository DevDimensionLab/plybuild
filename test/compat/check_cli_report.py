#!/usr/bin/env python3

import argparse
import json
import sys
from pathlib import Path


def fail(message):
    print("cli compatibility: " + message, file=sys.stderr)
    raise SystemExit(1)


def load_json(path, label):
    try:
        return json.loads(Path(path).read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        fail("cannot read {}: {}".format(label, error))


def command_map(tree):
    if not isinstance(tree, dict) or tree.get("schema_version") != 1:
        fail("command-tree schema_version must be 1")
    commands = tree.get("commands")
    if not isinstance(commands, list):
        fail("command tree must contain a commands list")
    mapped = {}
    for command in commands:
        if not isinstance(command, dict) or not isinstance(command.get("path"), str):
            fail("command entry is invalid")
        path = command["path"]
        if path in mapped:
            fail("duplicate command path: " + path)
        mapped[path] = command
    return mapped


def flag_map(command):
    mapped = {}
    for flag in command.get("flags", []):
        if not isinstance(flag, dict):
            fail("flag entry is invalid")
        key = "{}:{}".format(flag.get("scope"), flag.get("name"))
        if key in mapped or "None" in key:
            fail("duplicate or invalid flag: " + key)
        mapped[key] = flag
    return mapped


def changes(before, after, kind, path, ignored=()):
    output = []
    for field in sorted(set(before) | set(after)):
        if field in ignored:
            continue
        if before.get(field) != after.get(field):
            output.append(
                {
                    "after": after.get(field),
                    "before": before.get(field),
                    "field": field,
                    "kind": kind,
                    "path": path,
                }
            )
    return output


def compare(base, current):
    deltas = []
    base_commands = command_map(base)
    current_commands = command_map(current)
    for path in sorted(set(base_commands) | set(current_commands)):
        if path not in current_commands:
            deltas.append({"before": base_commands[path], "kind": "command_removed", "path": path})
            continue
        if path not in base_commands:
            deltas.append({"after": current_commands[path], "kind": "command_added", "path": path})
            continue
        old_command = base_commands[path]
        new_command = current_commands[path]
        deltas.extend(changes(old_command, new_command, "command_changed", path, ("flags",)))
        old_flags = flag_map(old_command)
        new_flags = flag_map(new_command)
        for flag_key in sorted(set(old_flags) | set(new_flags)):
            flag_path = path + "#" + flag_key
            if flag_key not in new_flags:
                deltas.append({"before": old_flags[flag_key], "kind": "flag_removed", "path": flag_path})
            elif flag_key not in old_flags:
                deltas.append({"after": new_flags[flag_key], "kind": "flag_added", "path": flag_path})
            else:
                deltas.extend(changes(old_flags[flag_key], new_flags[flag_key], "flag_changed", flag_path))
    return deltas


def canonical(value):
    return json.dumps(value, sort_keys=True, separators=(",", ":"))


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--base", required=True)
    parser.add_argument("--current", required=True)
    parser.add_argument("--allowlist", required=True)
    parser.add_argument("--output", required=True)
    parser.add_argument("--base-ref", required=True)
    args = parser.parse_args()

    base = load_json(args.base, "base command tree")
    current = load_json(args.current, "current command tree")
    allowlist = load_json(args.allowlist, "CLI allowlist")
    if set(allowlist) != {"allowed_deltas", "base_ref", "schema_version"}:
        fail("CLI allowlist keys are invalid")
    if allowlist["schema_version"] != 1 or allowlist["base_ref"] != args.base_ref:
        fail("CLI allowlist metadata is stale")
    allowed = allowlist["allowed_deltas"]
    if not isinstance(allowed, list):
        fail("allowed_deltas must be a list")
    allowed_canonical = [canonical(value) for value in allowed]
    if allowed_canonical != sorted(allowed_canonical) or len(allowed_canonical) != len(set(allowed_canonical)):
        fail("allowed_deltas must be canonically sorted and unique")

    deltas = compare(base, current)
    delta_canonical = sorted(canonical(value) for value in deltas)
    if delta_canonical != allowed_canonical:
        fail("CLI deltas do not match allowlist: expected {!r}, got {!r}".format(allowed, deltas))

    report = {"base_ref": args.base_ref, "deltas": deltas, "schema_version": 1}
    output = Path(args.output)
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n", encoding="utf-8")


if __name__ == "__main__":
    main()
