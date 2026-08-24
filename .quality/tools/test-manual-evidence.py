#!/usr/bin/env python3
"""Negative meta-tests for criterion-bound structured manual evidence."""

import copy
import hashlib
import json
import sys
import tempfile
from pathlib import Path


PARSER = Path(sys.argv[1])
namespace = {"__name__": "quality_manual_evidence_meta", "__file__": str(PARSER)}
exec(compile(PARSER.read_text(encoding="utf-8"), str(PARSER), "exec"), namespace)

COMMIT = "1" * 40
MODULE = "example.invalid/manual-evidence"
TREE = {
    "commit_tree": "2" * 40,
    "git_clean": True,
    "measurement_clean": True,
    "dirty_paths": [],
    "status_sha256": hashlib.sha256(b"\0").hexdigest(),
    "inventory_overlay_sha256": None,
}
INVENTORY = {
    "path": ".quality/inventory",
    "present": True,
    "sha256": "3" * 64,
    "overlay": False,
}
TOOL = {field: "instrument-" + field for field in namespace["BASELINE_TOOL_FIELDS"]}


def evidence_digest(value):
    payload = json.dumps(
        value, ensure_ascii=False, separators=(",", ":"), sort_keys=True,
    ).encode("utf-8")
    return hashlib.sha256(payload).hexdigest()


def receipt(criterion, kind, population):
    evidence = {"command": "fixture command for " + criterion, "population": population}
    return {
        "criterion": criterion,
        "kind": kind,
        "verdict": "PASS",
        "evidence": evidence,
        "evidence_sha256": evidence_digest(evidence),
    }


def valid_document():
    return {
        "schema_version": 2,
        "repository": {
            "module": MODULE,
            "commit": COMMIT,
            "tree": {
                "commit_tree": TREE["commit_tree"],
                "status_sha256": TREE["status_sha256"],
                "inventory_overlay_sha256": None,
            },
        },
        "inventory": INVENTORY,
        "instrument": TOOL,
        "criteria": [
            receipt("Q1.6", "test-double-contract", [{
                "subject": "fixture doubles", "dependencies": 2,
                "default_doubles": 2, "argument_recorders": 2,
                "dependency_struct_passed_whole": True,
            }]),
            receipt("Q1.7", "partial-failure-content", [{
                "subject": "fixture partial failures", "partial_failures": 2,
                "content_assertions": 2, "exit_only_assertions": 0,
            }]),
            receipt("Q1.9", "nonempty-iteration", [{
                "subject": "fixture collection", "iterated_collections": 2,
                "empty_population_assertions": 2,
            }]),
            receipt("Q2.4", "mutation-run", [{
                "subject": "scripts/mutate-fixture", "declared": 8,
                "killed": 8, "survived": 0, "unusable": 0,
            }]),
            receipt("Q2.8", "input-magnitude-control", [{
                "subject": "scripts/verify-fixture", "comparison": "magnitude",
                "control_magnitude": 10, "treatment_magnitude": 4,
                "expected_relation": "less",
            }]),
            receipt("Q2.9", "bad-input-read-only-control", [{
                "subject": "scripts/verify-fixture", "bad_input_exit": 2,
                "produced_artifacts": 0, "read_only_artifacts_changed": 0,
            }]),
        ],
        "findings": {
            "seam_test": [
                {"command": "seam 1", "mutation": "mutation 1", "measured": "1 killed", "unverified": "none 1"},
                {"command": "seam 2", "mutation": "mutation 2", "measured": "1 killed", "unverified": "none 2"},
                {"command": "seam 3", "mutation": "mutation 3", "measured": "1 killed", "unverified": "none 3"},
            ],
            "survivors": [
                {"location": "branch 1", "classification": "missing test", "measured": "1 survivor", "unverified": "none 1"},
                {"location": "branch 2", "classification": "redundant code", "measured": "1 survivor", "unverified": "none 2"},
                {"location": "branch 3", "classification": "fixture never reaches branch", "measured": "1 survivor", "unverified": "none 3"},
            ],
            "human_output": {"command": "fixture --help", "measured": "1 line", "unverified": "none"},
            "premise_check": {"premise": "fixture premise", "command": "inspect fixture", "measured": "1 fact", "held": True, "unverified": "none"},
        },
    }


def validate(document, tree=TREE):
    with tempfile.TemporaryDirectory(prefix="ply-manual-evidence-meta.") as temporary:
        path = Path(temporary) / "manual-evidence.json"
        path.write_text(json.dumps(document), encoding="utf-8")
        return namespace["validate_manual"](
            path, COMMIT, MODULE, tree, INVENTORY, TOOL,
        )


def assert_status(document, expected, label, tree=TREE):
    actual = validate(document, tree)
    if actual.get("status") != expected:
        raise AssertionError("{}: expected {}, got {}".format(label, expected, actual))


document = valid_document()
manual = validate(document)
if manual.get("status") != "valid" or manual.get("criterion_receipts") != 6:
    raise AssertionError("six valid criterion receipts were not accepted: " + repr(manual))
parsed = {
    "criteria": [
        {"id": criterion, "verdict": "UNMEASURABLE", "measured": "manual"}
        for criterion in namespace["MANUAL_CRITERION_KINDS"]
    ],
    "denominators": {
        "test_functions": 1, "mutation_harnesses": 1, "acceptance_scripts": 1,
    },
}
namespace["apply_manual_verdicts"](parsed, manual)
if any(item["verdict"] != "PASS" for item in parsed["criteria"]):
    raise AssertionError("valid receipts did not make every manual row reachable: " + repr(parsed))

stale = copy.deepcopy(document)
stale["repository"]["commit"] = "4" * 40
assert_status(stale, "stale", "stale commit")

dirty_tree = copy.deepcopy(TREE)
dirty_tree.update({"git_clean": False, "measurement_clean": False, "dirty_paths": ["value.go"]})
assert_status(document, "stale", "dirty tree", dirty_tree)

duplicate = copy.deepcopy(document)
duplicate["criteria"].append(copy.deepcopy(duplicate["criteria"][0]))
assert_status(duplicate, "invalid", "duplicate criterion")

empty = copy.deepcopy(document)
empty["criteria"][0]["evidence"] = {"command": "", "population": []}
empty["criteria"][0]["evidence_sha256"] = evidence_digest(empty["criteria"][0]["evidence"])
assert_status(empty, "invalid", "empty receipt")

wrong_kind = copy.deepcopy(document)
wrong_kind["criteria"][0]["kind"] = "mutation-run"
assert_status(wrong_kind, "invalid", "wrong receipt kind")

wrong_digest = copy.deepcopy(document)
wrong_digest["criteria"][0]["evidence_sha256"] = "0" * 64
assert_status(wrong_digest, "invalid", "wrong evidence digest")

false_pass = copy.deepcopy(document)
false_pass["criteria"][3]["evidence"]["population"][0]["killed"] = 7
false_pass["criteria"][3]["evidence_sha256"] = evidence_digest(false_pass["criteria"][3]["evidence"])
assert_status(false_pass, "invalid", "false PASS payload")

print("manual evidence rejects stale, dirty, duplicate, empty, wrong-kind, wrong-digest, and false-PASS receipts")
