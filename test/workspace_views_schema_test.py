#!/usr/bin/env python3
"""Validate the versioned workspace views and captured real CLI responses.

Uses the same pinned, test-only jsonschema dependency as read_contract_schema_test.
No Ply invocation or fixture writes occur here.
"""

import argparse
import copy
import json
from pathlib import Path
import unittest

from jsonschema import Draft202012Validator, FormatChecker

ROOT = Path(__file__).resolve().parents[1] / "schemas" / "read-contract"
SCHEMA = ROOT / "workspace-views-v1.schema.json"
EXAMPLES = ROOT / "workspace-views-examples"


def unique_object(pairs):
    out = {}
    for key, value in pairs:
        if key in out:
            raise ValueError(f"duplicate key: {key}")
        out[key] = value
    return out


def load(path):
    return json.loads(path.read_text(encoding="utf-8"), object_pairs_hook=unique_object)


class WorkspaceViewsSchema(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.schema = load(SCHEMA)
        Draft202012Validator.check_schema(cls.schema)
        cls.validator = Draft202012Validator(cls.schema, format_checker=FormatChecker())
        cls.examples = {p.name: load(p) for p in sorted(EXAMPLES.glob("*.json"))}

    def test_all_kinds_have_valid_and_invalid_examples(self):
        kinds = {entry["properties"]["kind"]["const"] for entry in self.schema["oneOf"]}
        valid, invalid = set(), set()
        for name, example in self.examples.items():
            expected = name.endswith(".valid.json")
            self.assertTrue(expected or name.endswith(".invalid.json"), name)
            with self.subTest(example=name):
                self.assertEqual(self.validator.is_valid(example), expected)
            (valid if expected else invalid).add(example["kind"])
        self.assertEqual(valid, kinds)
        self.assertEqual(invalid, kinds)

    def test_additive_fields_are_safe(self):
        for name, value in self.examples.items():
            if not name.endswith(".valid.json"):
                continue
            extended = copy.deepcopy(value)
            extended["future_field"] = {"unknown": [1, None]}
            with self.subTest(example=name):
                self.validator.validate(extended)

    def test_stale_result_basis_has_no_current_success_claim(self):
        value = copy.deepcopy(self.examples["progress-list.populated.valid.json"])
        progress = value["tasks"][0]["progress"]
        progress.update(state="in_progress", integration_classification="unknown",
                        has_progress=True, result_count=1, technical_gate=None,
                        human_qa_outcome=None, task_result_id=None,
                        integration_result_id=None,
                        reasons=["recorded_result_basis_stale"], next_actions=[{
                            "source": "content", "kind": "inspect_task_basis",
                            "actor": "agent", "severity": "next_step",
                            "reason": "Historical results do not cover the current Task basis.",
                            "since_utc": None, "evidence_ids": [],
                        }])
        self.validator.validate(value)

    def test_fundamental_enums_and_missing_values_fail(self):
        cases = [
            ("progress-list.populated.valid.json", ("tasks", 0, "lifecycle"), "deleted"),
            ("progress-list.populated.valid.json", ("tasks", 0, "progress", "state"), "done"),
            ("progress-list.populated.valid.json", ("tasks", 0, "progress", "technical_gate"), "pass"),
            ("progress-list.populated.valid.json", ("tasks", 0, "progress", "human_qa_outcome"), "approved"),
            ("runs.populated.valid.json", ("runs", 0, "provider"), "opus"),
            ("runs.populated.valid.json", ("runs", 0, "transport"), "browser"),
            ("runs.populated.valid.json", ("runs", 0, "state"), "running"),
        ]
        for name, path, wrong in cases:
            value = copy.deepcopy(self.examples[name])
            target = value
            for key in path[:-1]:
                target = target[key]
            target[path[-1]] = wrong
            with self.subTest(example=name, path=path):
                self.assertFalse(self.validator.is_valid(value))
        for name, value in self.examples.items():
            if name.endswith(".valid.json"):
                for field in ("kind", "schema_version", "workspace"):
                    changed = copy.deepcopy(value)
                    del changed[field]
                    with self.subTest(example=name, missing=field):
                        self.assertFalse(self.validator.is_valid(changed))


def validate_responses(directory):
    schema = load(SCHEMA)
    validator = Draft202012Validator(schema, format_checker=FormatChecker())
    supported = {item["properties"]["kind"]["const"] for item in schema["oneOf"]}
    seen = set()
    for path in sorted(directory.rglob("*.json")):
        value = load(path)
        if not isinstance(value, dict) or value.get("kind") not in supported:
            continue
        if not path.read_bytes().endswith(b"\n"):
            raise AssertionError(f"missing final newline: {path}")
        validator.validate(value)
        seen.add(value["kind"])
    if not seen:
        raise AssertionError("no supported workspace view responses found")
    print(f"Validated actual responses for {len(seen)} workspace view kinds.")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--responses", type=Path)
    parser.add_argument("--verbose", action="store_true")
    args = parser.parse_args()
    result = unittest.TextTestRunner(verbosity=2 if args.verbose else 1).run(
        unittest.defaultTestLoader.loadTestsFromTestCase(WorkspaceViewsSchema)
    )
    if not result.wasSuccessful():
        raise SystemExit(1)
    if args.responses:
        validate_responses(args.responses)
