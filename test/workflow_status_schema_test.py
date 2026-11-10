#!/usr/bin/env python3
"""Validate WorkflowStatusReadback@1 examples and optional actual CLI responses.

Use the pinned test/read_contract_requirements.txt validator dependency. This
script only reads fixtures and captured responses; it never runs Ply or writes.
"""

import argparse
import collections
import copy
import json
from pathlib import Path
import unittest

from jsonschema import Draft202012Validator, FormatChecker

ROOT = Path(__file__).resolve().parents[1] / "schemas" / "read-contract"
SCHEMA = ROOT / "workflow-status-v1.schema.json"
EXAMPLES = ROOT / "workflow-status-examples"
HIDDEN_CATEGORIES = ("inactive", "completed", "backlog")
CATEGORIES = ("needs_you", "follow_up", "in_progress", "ready_next", *HIDDEN_CATEGORIES)


def unique_object(pairs):
    value = {}
    for key, item in pairs:
        if key in value:
            raise ValueError(f"duplicate JSON key: {key}")
        value[key] = item
    return value


def load(path):
    return json.loads(path.read_text(encoding="utf-8"), object_pairs_hook=unique_object)


def nodes(value, path=()):
    yield path, value
    if isinstance(value, dict):
        for key, item in value.items():
            yield from nodes(item, (*path, key))
    elif isinstance(value, list):
        for index, item in enumerate(value):
            yield from nodes(item, (*path, index))


def at_path(value, path):
    for key in path:
        value = value[key]
    return value


def schema_branches(root, schema):
    """Expand this contract's local refs, nullable unions and allOf rules."""
    yield schema
    if "$ref" in schema:
        reference = schema["$ref"]
        if not reference.startswith("#/"):
            raise AssertionError("workflow status schema must be self-contained")
        target = at_path(root, tuple(part.replace("~1", "/").replace("~0", "~") for part in reference[2:].split("/")))
        yield from schema_branches(root, target)
    for keyword in ("anyOf", "allOf"):
        for branch in schema.get(keyword, ()):
            yield from schema_branches(root, branch)


def required_fields_at_path(root, path):
    schemas = [root]
    for key in path:
        children = []
        for schema in schemas:
            for branch in schema_branches(root, schema):
                child = branch.get("items") if isinstance(key, int) else branch.get("properties", {}).get(key)
                if child is not None:
                    children.append(child)
        schemas = children
    return {key for schema in schemas for branch in schema_branches(root, schema) for key in branch.get("required", ())}


def check_counts(value):
    counts = value["counts"]
    assert counts["visible"] == len(value["items"]), "visible count differs from items"
    assert counts["total"] == counts["visible"] + sum(counts["hidden"][key] for key in HIDDEN_CATEGORIES), "hidden counts overlap or omit Tasks"
    categories = collections.Counter(item["category"] for item in value["items"])
    assert all(categories[key] == counts["categories"][key] for key in CATEGORIES), "category counts differ from visible rows"
    assert len({item["task_id"] for item in value["items"]}) == len(value["items"]), "duplicate Task row"
    if value["include_all"]:
        assert not any(counts["hidden"][key] for key in HIDDEN_CATEGORIES), "--all retains hidden Tasks"
    for item in value["items"]:
        if item["category"] in ("inactive", "completed"):
            assert not any(action["current"] for action in item["next_actions"]), "historical Task has current actions"
        if item["category"] == "needs_you":
            assert item["progress"]["has_progress"], "untouched Task became Needs you"
            assert any(action["current"] and action["actor"] == "human" for action in item["next_actions"]), "Needs you lacks a current human action"


class WorkflowStatusSchema(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.schema = load(SCHEMA)
        Draft202012Validator.check_schema(cls.schema)
        cls.validator = Draft202012Validator(cls.schema, format_checker=FormatChecker())
        cls.examples = {path.name: load(path) for path in sorted(EXAMPLES.glob("*.json"))}

    def test_examples_cover_empty_populated_unknown_and_history(self):
        for name in ("empty.valid.json", "populated.valid.json", "historical.valid.json", "unknown.valid.json"):
            self.assertIn(name, self.examples)
        for name, value in self.examples.items():
            with self.subTest(example=name):
                expected = name.endswith(".valid.json")
                self.assertTrue(expected or name.endswith(".invalid.json"))
                self.assertEqual(self.validator.is_valid(value), expected, list(self.validator.iter_errors(value)))
                if expected:
                    check_counts(value)

    def test_required_fields_include_explicit_nulls(self):
        for name, value in self.examples.items():
            if not name.endswith(".valid.json"):
                continue
            for path, node in nodes(value):
                if not isinstance(node, dict):
                    continue
                required = required_fields_at_path(self.schema, path)
                for key in node:
                    changed = copy.deepcopy(value)
                    del at_path(changed, path)[key]
                    with self.subTest(example=name, path=(*path, key)):
                        self.assertEqual(self.validator.is_valid(changed), key not in required)

    def test_additive_fields_are_safe_at_every_object_level(self):
        for name, value in self.examples.items():
            if not name.endswith(".valid.json"):
                continue
            changed = copy.deepcopy(value)
            for path, node in nodes(value):
                if isinstance(node, dict):
                    at_path(changed, path)["future_field"] = {"opaque": [None, True, 7]}
            with self.subTest(example=name):
                self.validator.validate(changed)
                check_counts(changed)

    def test_startup_recovery_fields_are_typed_without_runtime_claims(self):
        original = self.examples["startup-recovery.valid.json"]
        prefix = ("items", 0, "runs", 0, "startup_recovery")
        for suffix, wrong in (
            (("generation",), 0), (("generation",), True),
            (("binding", "sha256"), "unbound"),
            (("provider_executable", "path"), None),
            (("control_executable", "sha256"), "sha256:no"),
            (("transport_mode",), "restart_authorized"),
            (("original_transport", "observed_at"), None),
            (("original_transport", "launch_pending"), "live"),
            (("replacement_transport", "sha256"), "unknown"),
        ):
            changed = copy.deepcopy(original)
            path = (*prefix, *suffix)
            at_path(changed, path[:-1])[path[-1]] = wrong
            with self.subTest(path=path):
                self.assertFalse(self.validator.is_valid(changed))

    def test_enum_and_timestamp_errors_do_not_pass(self):
        cases = [
            (("schema_version",), 2), (("schema_version",), True),
            (("freshness",), "live"), (("source_basis",), "observed_git"),
            (("observed_at_utc",), "yesterday"),
            (("items", 0, "category"), "done"),
            (("items", 0, "task_lifecycle"), "deleted"),
            (("items", 0, "next_actions", 0, "actor"), None),
            (("items", 0, "last_activity_utc"), "today"),
            (("items", 0, "progress", "technical_gate"), "pass"),
            (("items", 0, "progress", "human_qa_outcome"), "approved"),
            (("items", 0, "runs", 0, "state"), "working"),
            (("items", 0, "runs", 0, "provider"), "model-name"),
            (("items", 3, "queues", 0, "readiness"), "start_authorized"),
            (("items", 3, "queues", 0, "rank"), 0),
        ]
        for path, wrong in cases:
            changed = copy.deepcopy(self.examples["populated.valid.json"])
            at_path(changed, path[:-1])[path[-1]] = wrong
            with self.subTest(path=path):
                self.assertFalse(self.validator.is_valid(changed))

    def test_action_codes_are_open_and_actor_is_explicit(self):
        value = copy.deepcopy(self.examples["populated.valid.json"])
        action = value["items"][0]["next_actions"][0]
        action.update(kind="future_owner_action", actor="unknown")
        self.validator.validate(value)

    def test_counts_detect_dropped_rows_and_overlapping_hidden_tasks(self):
        original = self.examples["populated.valid.json"]
        changed = copy.deepcopy(original)
        changed["items"].pop()
        with self.assertRaises(AssertionError):
            check_counts(changed)
        changed = copy.deepcopy(original)
        changed["counts"]["hidden"]["backlog"] += 1
        with self.assertRaises(AssertionError):
            check_counts(changed)


def validate_responses(directory):
    validator = Draft202012Validator(load(SCHEMA), format_checker=FormatChecker())
    checked = 0
    for path in sorted(directory.rglob("*.json")):
        value = load(path)
        if not isinstance(value, dict) or value.get("kind") != "WorkflowStatusReadback@1":
            continue
        if not path.read_bytes().endswith(b"\n"):
            raise AssertionError(f"missing final newline: {path}")
        validator.validate(value)
        check_counts(value)
        checked += 1
    if not checked:
        raise AssertionError("no actual WorkflowStatusReadback@1 responses found")
    print(f"Validated {checked} actual WorkflowStatusReadback@1 responses.")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--responses", type=Path)
    parser.add_argument("--verbose", action="store_true")
    args = parser.parse_args()
    result = unittest.TextTestRunner(verbosity=2 if args.verbose else 1).run(
        unittest.defaultTestLoader.loadTestsFromTestCase(WorkflowStatusSchema)
    )
    if not result.wasSuccessful():
        raise SystemExit(1)
    if args.responses:
        validate_responses(args.responses)
