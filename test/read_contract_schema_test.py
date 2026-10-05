#!/usr/bin/env python3
"""Validate the public read schemas, their examples and optional CLI responses.

Install test/read_contract_requirements.txt in a virtual environment first. Run
with --responses DIRECTORY to also validate captured JSON stdout from all four
new CLI response kinds. No Ply command or fixture mutation is performed here.
"""

import argparse
import copy
import itertools
import json
from pathlib import Path
import sys
import unittest

try:
    from jsonschema import Draft202012Validator
except ImportError:
    raise SystemExit(
        "Install the test-only validator: python -m pip install "
        "-r test/read_contract_requirements.txt"
    )


SCHEMA_ROOT = Path(__file__).resolve().parents[1] / "schemas" / "read-contract"
SCHEMAS = {
    "WorkspaceProjectListReadback@1": "workspace-project-list-readback-v1",
    "WorkspaceProjectReadback@1": "workspace-project-readback-v1",
    "WorkspaceTaskListReadback@1": "workspace-task-list-readback-v1",
    "PlyCapabilities@1": "ply-capabilities-v1",
}


def unique_object(pairs):
    result = {}
    for key, value in pairs:
        if key in result:
            raise ValueError(f"duplicate JSON key: {key}")
        result[key] = value
    return result


def load_json(path):
    return json.loads(path.read_text(encoding="utf-8"), object_pairs_hook=unique_object)


def nodes(value, path=()):
    """Yield every existing instance node, including array members."""
    yield path, value
    if isinstance(value, dict):
        for key, child in value.items():
            yield from nodes(child, (*path, key))
    elif isinstance(value, list):
        for index, child in enumerate(value):
            yield from nodes(child, (*path, index))


def at_path(value, path):
    for key in path:
        value = value[key]
    return value


def replace(value, path, replacement):
    if not path:
        return replacement
    result = copy.deepcopy(value)
    at_path(result, path[:-1])[path[-1]] = replacement
    return result


def wrong_type(value):
    if isinstance(value, list):
        return None
    if isinstance(value, dict):
        return []
    if isinstance(value, bool):
        return 1
    if isinstance(value, int):
        return "1"
    return []


class ReadContractSchemas(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.validators = {}
        for kind, name in SCHEMAS.items():
            schema = load_json(SCHEMA_ROOT / f"{name}.schema.json")
            Draft202012Validator.check_schema(schema)
            cls.validators[kind] = Draft202012Validator(schema)
        cls.valid_examples = {}
        cls.invalid_examples = {}
        for path in sorted((SCHEMA_ROOT / "examples").glob("*.json")):
            value = load_json(path)
            if path.name.endswith(".valid.json"):
                cls.valid_examples[path.name] = value
            elif path.name.endswith(".invalid.json"):
                cls.invalid_examples[path.name] = value
            else:
                raise AssertionError(f"example is not marked valid/invalid: {path}")

    def valid(self, value):
        self.validators[value["kind"]].validate(value)

    def test_examples_cover_each_schema_and_expected_validity(self):
        for examples, expect_valid in (
            (self.valid_examples, True), (self.invalid_examples, False)
        ):
            self.assertEqual({example["kind"] for example in examples.values()}, set(SCHEMAS))
            for name, example in examples.items():
                with self.subTest(example=name):
                    self.assertTrue(name.startswith(SCHEMAS[example["kind"]] + "."))
                    self.assertEqual(self.validators[example["kind"]].is_valid(example), expect_valid)

    def test_every_named_field_is_required_including_nullable_fields(self):
        for name, example in self.valid_examples.items():
            validator = self.validators[example["kind"]]
            for path, value in nodes(example):
                if not isinstance(value, dict):
                    continue
                for key in value:
                    with self.subTest(example=name, path=(*path, key)):
                        altered = copy.deepcopy(example)
                        del at_path(altered, path)[key]
                        self.assertFalse(validator.is_valid(altered), "missing field was accepted")

    def test_wrong_types_are_rejected_at_every_level(self):
        for name, example in self.valid_examples.items():
            validator = self.validators[example["kind"]]
            for path, value in nodes(example):
                with self.subTest(example=name, path=path):
                    self.assertFalse(validator.is_valid(replace(example, path, wrong_type(value))))

    def test_additive_fields_are_allowed_at_every_object_level(self):
        for name, example in self.valid_examples.items():
            with self.subTest(example=name):
                altered = copy.deepcopy(example)
                for path, value in nodes(example):
                    if isinstance(value, dict):
                        at_path(altered, path)["future_field"] = {"opaque": [None, True, 7]}
                self.valid(altered)

    def test_nulls_are_explicit_and_limited_to_declared_fields(self):
        for name, example in self.valid_examples.items():
            validator = self.validators[example["kind"]]
            for path, value in nodes(example):
                if not path:
                    continue
                nullable = value is None
                nullable |= path in (("build", "vcs_revision"), ("build", "vcs_modified"))
                nullable |= example["kind"] == "WorkspaceTaskListReadback@1" and path in (
                    ("scope", "project_id"), ("scope", "repo_id"), ("scope", "epic_id")
                )
                with self.subTest(example=name, path=path):
                    self.assertEqual(validator.is_valid(replace(example, path, None)), nullable)

    def test_empty_collections_are_arrays(self):
        cases = (
            ("workspace-project-list-readback-v1.populated.valid.json", "projects"),
            ("workspace-project-readback-v1.populated.valid.json", "repositories"),
            ("workspace-task-list-readback-v1.populated.valid.json", "tasks"),
        )
        for name, field in cases:
            example = self.valid_examples[name]
            with self.subTest(field=field):
                self.valid(replace(example, (field,), []))
                self.assertFalse(self.validators[example["kind"]].is_valid(replace(example, (field,), None)))

    def test_title_state_combinations_and_literal_unavailable_text(self):
        example = self.valid_examples["workspace-task-list-readback-v1.populated.valid.json"]
        validator = self.validators[example["kind"]]
        for title, source, status in itertools.product(
            (None, "Current title", "Problem content unavailable", ""),
            ("registration", "problem_revision"),
            ("available", "unavailable"),
        ):
            with self.subTest(title=title, source=source, status=status):
                altered = copy.deepcopy(example)
                altered["tasks"][0].update(title=title, title_source=source, title_status=status)
                expected = (
                    status == "available" and isinstance(title, str)
                ) or (
                    status == "unavailable" and source == "problem_revision" and title is None
                )
                self.assertEqual(validator.is_valid(altered), expected)

    def test_raw_worktree_state_remains_a_string_without_a_new_enum(self):
        example = self.valid_examples["workspace-task-list-readback-v1.populated.valid.json"]
        self.valid(replace(example, ("tasks", 0, "worktree_state"), "future_registered_state"))

    def test_version_and_repository_counts_are_integers(self):
        for example in self.valid_examples.values():
            validator = self.validators[example["kind"]]
            for version in (True, "1", 1.5, 2, None):
                with self.subTest(kind=example["kind"], version=version):
                    self.assertFalse(validator.is_valid(replace(example, ("schema_version",), version)))
        example = self.valid_examples["workspace-project-readback-v1.populated.valid.json"]
        validator = self.validators[example["kind"]]
        for count in (-1, 0.5, True, "2", None):
            with self.subTest(count=count):
                self.assertFalse(validator.is_valid(replace(example, ("project", "repo_count"), count)))

    def test_catalog_preserves_modes_selectors_and_legacy_schema_marker(self):
        example = self.valid_examples["ply-capabilities-v1.build-known.valid.json"]
        validator = self.validators[example["kind"]]
        cases = (
            (("operations",), example["operations"][::-1]),
            (("operations",), example["operations"][:-1]),
            (("operations", 2, "selectors"), []),
            (("operations", 2, "filter_policy"), "independent_and"),
            (("operations", 3, "filter_policy"), "all_or_none"),
            (("operations", 3, "command"), ["workspace", "task", "show"]),
            (("operations", 4, "result_schemas", 0, "schema_version"), 1),
            (("operations", 4, "result_schemas", 1, "schema_version"), None),
            (("operations", 4, "result_schemas", 2, "schema_version"), 2),
        )
        for path, value in cases:
            with self.subTest(path=path, value=value):
                self.assertFalse(validator.is_valid(replace(example, path, value)))


def validate_responses(directory):
    seen = set()
    files = sorted(directory.rglob("*.json"))
    for path in files:
        raw = path.read_bytes()
        if not raw.endswith(b"\n"):
            raise AssertionError(f"response lacks terminating newline: {path}")
        value = load_json(path)
        if not isinstance(value, dict) or value.get("kind") not in SCHEMAS:
            raise AssertionError(f"response has unsupported kind: {path}")
        kind = value["kind"]
        ReadContractSchemas.validators[kind].validate(value)
        seen.add(kind)
    if seen != set(SCHEMAS):
        raise AssertionError(f"missing response kinds: {sorted(set(SCHEMAS) - seen)}")
    print(f"Validated {len(files)} CLI responses covering all four schemas.")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--responses", type=Path, help="directory of captured CLI JSON responses")
    parser.add_argument("--verbose", "-v", action="store_true")
    args = parser.parse_args()
    suite = unittest.defaultTestLoader.loadTestsFromTestCase(ReadContractSchemas)
    result = unittest.TextTestRunner(verbosity=2 if args.verbose else 1).run(suite)
    if not result.wasSuccessful():
        return 1
    if args.responses is not None:
        validate_responses(args.responses)
    return 0


if __name__ == "__main__":
    sys.exit(main())
