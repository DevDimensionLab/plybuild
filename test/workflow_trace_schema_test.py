#!/usr/bin/env python3
"""Validate WorkflowTraceReadback@1 fixtures and captured actual CLI responses.

Use test/read_contract_requirements.txt. This validator only reads existing
files; it does not execute Ply or change a workspace. --responses must contain
at least one actual trace response and validates its references as well as shape.
"""

import argparse
import copy
import datetime as dt
import json
from pathlib import Path
import unittest

from jsonschema import Draft202012Validator, FormatChecker

ROOT = Path(__file__).resolve().parents[1] / "schemas" / "read-contract"
SCHEMA = ROOT / "workflow-trace-v1.schema.json"
EXAMPLES = ROOT / "workflow-trace-examples"
NATIVE_PAYLOADS = {"data", "declaration", "frozen_goal", "declared_process", "frozen_basis"}


def validator_for(schema):
    checker = FormatChecker()

    # jsonschema's optional RFC3339 package is not in the pinned dependency
    # set. Require actual calendar/time validity even without that extra.
    @checker.checks("date-time", raises=ValueError)
    def valid_datetime(value):
        return not isinstance(value, str) or dt.datetime.fromisoformat(value.replace("Z", "+00:00")).tzinfo is not None

    return Draft202012Validator(schema, format_checker=checker)


def unique_object(pairs):
    result = {}
    for key, value in pairs:
        if key in result:
            raise ValueError(f"duplicate JSON key: {key}")
        result[key] = value
    return result


def load(path):
    return json.loads(path.read_text(encoding="utf-8"), object_pairs_hook=unique_object)


def nodes(value, path=()):
    """Projection nodes only; native versioned payloads own their shape."""
    yield path, value
    if isinstance(value, dict):
        for key, item in value.items():
            if key not in NATIVE_PAYLOADS:
                yield from nodes(item, (*path, key))
    elif isinstance(value, list):
        for index, item in enumerate(value):
            yield from nodes(item, (*path, index))


def at_path(value, path):
    for key in path:
        value = value[key]
    return value


def check_references(value):
    tables = {}
    for name in ("events", "sources", "runs", "candidates", "chains", "analysis"):
        table = {row["id"]: row for row in value[name]}
        assert len(table) == len(value[name]), f"duplicate {name} identity"
        tables[name] = table
    references = {"source_ids": "sources", "evidence_ids": "events", "event_ids": "events", "run_ids": "runs"}
    for path, node in nodes(value):
        if not isinstance(node, dict):
            continue
        for key, table in references.items():
            for reference in node.get(key, []):
                assert reference in tables[table], f"unresolved {key} at {path}: {reference}"
    for event in value["events"]:
        if event["candidate_id"] is not None:
            assert event["candidate_id"] in tables["candidates"], "unresolved candidate"
        for position in event["positions"]:
            assert position["chain_id"] in tables["chains"], "unresolved source chain"
            assert event["id"] in tables["chains"][position["chain_id"]]["event_ids"], "position omitted from source chain"
    for relation in value["relations"]:
        assert relation["from"] in tables["events"] and relation["to"] in tables["events"], "unresolved relation endpoint"
    for chain in value["chains"]:
        positions = [(p["sequence"], event["id"])
                     for event in value["events"] for p in event["positions"]
                     if p["chain_id"] == chain["id"]]
        # One exact event can be reached at several preserved positions.
        # Compare the entire source-position list instead of choosing a time
        # or first occurrence for that event.
        assert chain["event_ids"] == [eid for _, eid in sorted(positions)], "source sequence or tied-position presentation was changed"
    for candidate in value["candidates"]:
        for eid in candidate["event_ids"]:
            assert tables["events"][eid]["candidate_id"] == candidate["id"], "candidate has foreign evidence"
        for axis in candidate["axes"].values():
            assert set(axis["evidence_ids"]).issubset(candidate["event_ids"]), "axis has foreign evidence"
            outcomes = set(axis["outcomes"])
            if axis["state"] == "not_recorded":
                assert not outcomes and not axis["evidence_ids"], "absence axis contains a claim"
            else:
                assert axis["evidence_ids"] and outcomes, "recorded axis lacks evidence"
                assert (len(outcomes) > 1) == (axis["state"] == "conflicting"), "axis conflict was flattened"
    for analysis in value["analysis"]:
        if analysis["value"] is not None and analysis["value"] > 0:
            assert analysis["evidence_ids"], "analysis value lacks evidence"
        if analysis["coverage"] == "unknown":
            assert analysis["value"] is None and analysis["unknowns"], "unknown analysis manufactures a value"


class WorkflowTraceSchema(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.schema = load(SCHEMA)
        Draft202012Validator.check_schema(cls.schema)
        cls.validator = validator_for(cls.schema)
        cls.examples = {path.name: load(path) for path in sorted(EXAMPLES.glob("*.json"))}

    def test_empty_multifamily_and_partial_examples(self):
        for name in ("empty.valid.json", "multi-family.valid.json", "partial.valid.json"):
            self.assertIn(name, self.examples)
        for name, value in self.examples.items():
            with self.subTest(example=name):
                expected = name.endswith(".valid.json")
                self.assertTrue(expected or name.endswith(".invalid.json"))
                self.assertEqual(self.validator.is_valid(value), expected, list(self.validator.iter_errors(value)))
                if expected:
                    check_references(value)

    def test_projection_objects_reject_unknown_or_missing_fields(self):
        original = self.examples["multi-family.valid.json"]
        for path, node in nodes(original):
            if not isinstance(node, dict):
                continue
            changed = copy.deepcopy(original)
            at_path(changed, path)["invented_authority"] = True
            with self.subTest(path=path, operation="add"):
                self.assertFalse(self.validator.is_valid(changed))
            for key in node:
                changed = copy.deepcopy(original)
                del at_path(changed, path)[key]
                with self.subTest(path=(*path, key), operation="omit"):
                    self.assertFalse(self.validator.is_valid(changed))

    def test_versions_roles_times_and_authority_are_not_free_text(self):
        for path, wrong in (
            (("kind",), "WorkflowTraceReadback@2"), (("schema_version",), 2),
            (("schema_version",), True), (("live_state",), "alive"),
            (("next_transition_authorized",), True), (("freshness",), "live"),
            (("observed_at_utc",), "2026-99-07T10:00:00Z"),
            (("events", 0, "role"), "provider"),
            (("events", 0, "occurred_at_utc"), "yesterday"),
            (("events", 0, "reported_at_utc"), 0),
            (("events", 0, "registered_at_utc"), ""),
            (("events", 0, "evidence_class"), "authenticated_human"),
            (("events", 0, "time_basis", "kind"), "mtime"),
            (("events", 0, "positions", 0, "sequence"), -1),
            (("sources", 0, "sha256"), "unbound"),
            (("analysis", 0, "coverage"), "measured_improvement"),
            (("analysis", 0, "value"), -1),
        ):
            changed = copy.deepcopy(self.examples["multi-family.valid.json"])
            at_path(changed, path[:-1])[path[-1]] = wrong
            with self.subTest(path=path):
                self.assertFalse(self.validator.is_valid(changed))

    def test_native_payloads_remain_governed_by_native_contracts(self):
        changed = copy.deepcopy(self.examples["multi-family.valid.json"])
        changed["events"][0]["data"]["future_native_evidence"] = {"opaque": [None, 7, True]}
        changed["runs"][2]["frozen_goal"]["future_native_obligation"] = ["a", "b"]
        self.validator.validate(changed)

    def test_source_uncertainty_preserves_offsets_separately_from_utc_fields(self):
        changed = copy.deepcopy(self.examples["multi-family.valid.json"])
        changed["events"][0]["time_basis"] = dict(
            kind="reported", clock="source clock", precision="second",
            uncertainty=dict(earliest="2026-10-07T10:00:00+02:00", latest="2026-10-07T11:00:00+02:00"),
        )
        self.validator.validate(changed)
        changed["events"][0]["time_basis"]["uncertainty"]["latest"] = "not-a-time"
        self.assertFalse(self.validator.is_valid(changed))

    def test_broken_links_duplicate_ids_and_candidate_contamination_fail(self):
        original = self.examples["multi-family.valid.json"]
        cases = []
        changed = copy.deepcopy(original); changed["analysis"][0]["evidence_ids"] = ["missing"]; cases.append(changed)
        changed = copy.deepcopy(original); changed["events"].append(changed["events"][0]); cases.append(changed)
        changed = copy.deepcopy(original); changed["relations"][0]["to"] = "missing"; cases.append(changed)
        changed = copy.deepcopy(original); changed["candidates"][0]["event_ids"].append("qa-c3-pass"); cases.append(changed)
        changed = copy.deepcopy(original); changed["chains"][2]["event_ids"].reverse(); cases.append(changed)
        changed = copy.deepcopy(original); changed["analysis"][1]["value"] = 0; cases.append(changed)
        for changed in cases:
            with self.assertRaises(AssertionError):
                check_references(changed)

    def test_duplicate_json_keys_are_not_accepted(self):
        with self.assertRaises(ValueError):
            json.loads('{"role":"human","role":"agent"}', object_pairs_hook=unique_object)

    def test_one_event_can_preserve_multiple_source_positions(self):
        changed = copy.deepcopy(self.examples["multi-family.valid.json"])
        changed["events"][0]["positions"].append(dict(chain_id="legacy-events", sequence=3))
        changed["events"][1]["positions"].append(dict(chain_id="legacy-events", sequence=2))
        changed["chains"][0]["event_ids"] = ["legacy-report", "older-review", "legacy-report"]
        self.validator.validate(changed)
        check_references(changed)


def validate_responses(directory):
    validator = validator_for(load(SCHEMA))
    checked = 0
    for path in sorted(directory.rglob("*.json")):
        value = load(path)
        if not isinstance(value, dict) or value.get("kind") != "WorkflowTraceReadback@1":
            continue
        data = path.read_bytes()
        assert data.endswith(b"\n") and data.count(b"\n") == 1, f"not one JSON object plus newline: {path}"
        validator.validate(value)
        check_references(value)
        checked += 1
    if not checked:
        raise AssertionError("no actual WorkflowTraceReadback@1 responses found")
    print(f"Validated {checked} actual WorkflowTraceReadback@1 responses.")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--responses", type=Path)
    parser.add_argument("--verbose", action="store_true")
    args = parser.parse_args()
    result = unittest.TextTestRunner(verbosity=2 if args.verbose else 1).run(
        unittest.defaultTestLoader.loadTestsFromTestCase(WorkflowTraceSchema)
    )
    if not result.wasSuccessful():
        raise SystemExit(1)
    if args.responses:
        validate_responses(args.responses)
