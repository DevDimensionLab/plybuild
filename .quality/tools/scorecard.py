#!/usr/bin/env python3
"""Convert the pinned Markdown audit into deterministic, project-aware JSON."""

import argparse
import hashlib
import json
import os
import re
import shlex
import shutil
import stat
import subprocess
import sys
import tempfile
from pathlib import Path, PurePosixPath


SCHEMA_VERSION = 2
WRAPPER_VERSION = 1
EXPECTED_LEVEL_COUNTS = {"L0": 8, "L1": 9, "L2": 10, "L3": 9}
EXPECTED_IDS = (
    {"Q0." + str(value) for value in range(1, 9)} |
    {"Q1." + str(value) for value in range(1, 10)} |
    {"Q2." + str(value) for value in range(1, 11)} |
    {"Q3." + str(value) for value in range(1, 10)}
)
RATCHET_IDS = {"Q0.6", "Q0.8", "Q1.1", "Q1.2", "Q1.3", "Q1.4", "Q2.1", "Q3.4"}
INVENTORY_SECTIONS = ("subjects", "seams", "features", "adapters")
TEST_SUPPORT_PATHS = ("internal/testutil",)
EXCLUDED_REPOSITORY_ROOTS = (".git", ".quality", "target", "vendor")
SURVIVOR_CLASSES = {"missing test", "redundant code", "fixture never reaches branch"}
MANUAL_CRITERION_KINDS = {
    "Q1.6": "test-double-contract",
    "Q1.7": "partial-failure-content",
    "Q1.9": "nonempty-iteration",
    "Q2.4": "mutation-run",
    "Q2.8": "input-magnitude-control",
    "Q2.9": "bad-input-read-only-control",
}
MANUAL_CRITERION_POPULATIONS = {
    "Q1.6": "test_functions",
    "Q1.7": "test_functions",
    "Q1.9": "test_functions",
    "Q2.4": "mutation_harnesses",
    "Q2.8": "acceptance_scripts",
    "Q2.9": "acceptance_scripts",
}
BASELINE_TOOL_FIELDS = (
    "wrapper_version", "upstream_version", "upstream_sha256", "parser_sha256",
    "wrapper_sha256", "call_scanner_sha256", "q06_contract_sha256",
    "report_template_sha256",
)
BUILD_SELECTOR_ENV = (
    "GO386", "GOAMD64", "GOARM", "GOARM64", "GOMIPS", "GOMIPS64", "GOPPC64",
    "GORISCV64", "GOWASM", "GOEXPERIMENT",
)
BASELINE_BUILD_FIELDS = (
    "version", "goos", "goarch", "cgo_enabled", "goenv", "goflags", "gowork",
) + tuple(name.lower() for name in BUILD_SELECTOR_ENV)


class AuditBroken(Exception):
    pass


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("--report", required=True)
    parser.add_argument("--output", required=True)
    parser.add_argument("--display-output")
    parser.add_argument("--repo", required=True)
    parser.add_argument("--commit", required=True)
    parser.add_argument("--upstream", required=True)
    parser.add_argument("--upstream-exit", required=True, type=int)
    parser.add_argument("--baseline")
    parser.add_argument("--manual-evidence")
    parser.add_argument("--inventory-overlay-sha256")
    parser.add_argument("--partial", action="store_true")
    return parser.parse_args()


def sha256_bytes(value):
    return hashlib.sha256(value).hexdigest()


def sha256_file(path):
    return sha256_bytes(Path(path).read_bytes())


def module_name(repo):
    go_mod = Path(repo) / "go.mod"
    for line in go_mod.read_text(encoding="utf-8").splitlines():
        if line.startswith("module "):
            return line.split(None, 1)[1].strip()
    raise AuditBroken("go.mod has no module declaration")


def normalize_inventory_path(value, label):
    if not value or "\\" in value or value.startswith("/"):
        raise AuditBroken("{} must be a non-empty POSIX repository-relative path: {}".format(label, value))
    parsed = PurePosixPath(value)
    if value != parsed.as_posix() or value in (".", "..") or ".." in parsed.parts:
        raise AuditBroken("{} is not a normalized repository-relative path: {}".format(label, value))
    return value


def parse_inventory_bytes(raw):
    result = {name: [] for name in INVENTORY_SECTIONS}
    section = None
    try:
        lines = raw.decode("utf-8").splitlines()
    except UnicodeDecodeError as error:
        raise AuditBroken(".quality/inventory is not UTF-8: " + str(error))
    for number, raw_line in enumerate(lines, 1):
        line = raw_line.strip()
        match = re.fullmatch(r"\[([a-z]+)\]", line)
        if match:
            section = match.group(1) if match.group(1) in result else None
        elif section and line and not line.startswith("#"):
            result[section].append({"line": number, "value": line})
    return result


def inventory_counts(inventory):
    return {"declared_" + name: len(inventory[name]) for name in INVENTORY_SECTIONS}


def inventory_snapshot(repo, overlay_sha256):
    path = Path(repo) / ".quality" / "inventory"
    if not path.is_file():
        if os.path.lexists(path):
            raise AuditBroken(".quality/inventory must be a readable regular file")
        if overlay_sha256:
            raise AuditBroken("inventory overlay is absent")
        return (
            {name: [] for name in INVENTORY_SECTIONS},
            {"path": ".quality/inventory", "present": False, "sha256": None, "overlay": False},
        )
    try:
        raw = path.read_bytes()
    except OSError as error:
        raise AuditBroken("could not read .quality/inventory: " + str(error))
    digest = sha256_bytes(raw)
    if overlay_sha256 and digest != overlay_sha256:
        raise AuditBroken("inventory overlay checksum does not match --inventory-overlay-sha256")
    return parse_inventory_bytes(raw), {
        "path": ".quality/inventory", "present": True, "sha256": digest,
        "overlay": bool(overlay_sha256),
    }


def verify_inventory_identity(repo, expected):
    overlay_sha256 = expected["sha256"] if expected["overlay"] else None
    _, observed = inventory_snapshot(repo, overlay_sha256)
    if observed != expected:
        raise AuditBroken(".quality/inventory changed during audit")


def metric(criterion_id, measured):
    def count(name, pattern, direction="max"):
        match = re.search(pattern, measured)
        if not match:
            return None
        return {"name": name, "value": int(match.group(1)), "direction": direction}

    if criterion_id == "Q0.6":
        skipped = count("skipped_tests", r"(\d+) skipped tests(?: of \d+)?")
        guard = re.search(r"leak guard defined (\d+) times, called (\d+) times", measured)
        if skipped and guard:
            skipped["precondition_met"] = int(guard.group(1)) > 0 and int(guard.group(2)) > 0
            skipped["guard_call_sites"] = int(guard.group(2))
        used = re.search(r"leak guard used at (\d+) call sites", measured)
        if skipped and used:
            skipped["precondition_met"] = int(used.group(1)) > 0
            skipped["guard_call_sites"] = int(used.group(1))
        return skipped
    if criterion_id == "Q0.8":
        value = count("scripts_without_meta_test", r"(\d+) of (\d+) scripts")
        match = re.search(r"(\d+) of (\d+) scripts", measured)
        if value and match:
            value["population"] = int(match.group(2))
        passed = re.search(r"(\d+) scripts, all with a test- counterpart", measured)
        if not value and passed:
            value = {
                "name": "scripts_without_meta_test", "value": 0, "direction": "max",
                "population": int(passed.group(1)),
            }
        return value
    if criterion_id == "Q1.1":
        value = count("packages_without_tests", r"(\d+) of (\d+) packages")
        match = re.search(r"(\d+) of (\d+) packages", measured)
        if value and match:
            value["population"] = int(match.group(2))
        return value
    if criterion_id == "Q1.2":
        return count("process_exiting_calls", r"(\d+) process-exiting calls")
    if criterion_id == "Q1.3":
        value = count("direct_external_files", r"(\d+) files make direct")
        passed = re.search(r"(\d+) direct calls outside (\d+) declared adapters", measured)
        if not value and passed:
            value = {
                "name": "direct_external_files", "value": int(passed.group(1)),
                "direction": "max", "declared_adapters": int(passed.group(2)),
            }
        return value
    if criterion_id == "Q1.4":
        match = re.search(r"(\d+) of (\d+) declared seams", measured)
        if match:
            return {
                "name": "covered_seams", "numerator": int(match.group(1)),
                "denominator": int(match.group(2)), "direction": "min_ratio",
            }
    if criterion_id == "Q2.1":
        match = re.search(r"(\d+) harnesses for (\d+) declared subjects", measured)
        if match:
            return {
                "name": "subjects_with_harnesses", "numerator": int(match.group(1)),
                "denominator": int(match.group(2)), "direction": "min_ratio",
            }
    if criterion_id == "Q3.4":
        return count("state_claim_phrases", r"(\d+) state-claim phrases")
    return None


def parse_report(path):
    text = Path(path).read_text(encoding="utf-8")
    header = re.search(r"^commit ([0-9a-f]+) . tool version (\S+) . (.+ mode)$", text, re.MULTILINE)
    if not header:
        raise AuditBroken("could not parse the commit/tool header")
    denominators = {}
    patterns = [
        (r"packages (\d+) . with tests (\d+) . go files (\d+) . test functions (\d+) \((\d+) skipped\)",
         ("packages", "packages_with_tests", "go_files", "test_functions", "skipped_tests")),
        (r"table-driven tests (\d+) . scripts (\d+) . mutation harnesses (\d+) . acceptance scripts (\d+)",
         ("table_driven_tests", "scripts", "mutation_harnesses", "acceptance_scripts")),
        (r"enabled CI workflows (\d+) . inventory (.+)", ("enabled_ci_workflows", "inventory")),
    ]
    for pattern, names in patterns:
        match = re.search(pattern, text)
        if not match:
            raise AuditBroken("could not parse a denominator line: " + pattern)
        for name, value in zip(names, match.groups()):
            denominators[name] = int(value) if value.isdigit() else value
    criteria = []
    current = None
    for line in text.splitlines():
        found = re.match(r"^  (PASS|FAIL|UNMEASURABLE|N-A)\s+(Q\d+\.\d+)\s+(.+)$", line)
        if found:
            current = {
                "id": found.group(2), "level": "L" + found.group(2)[1],
                "verdict": found.group(1), "title": found.group(3).strip(), "measured": "",
            }
            criteria.append(current)
            continue
        if current and line.startswith("               measured: "):
            current["measured"] = line.split("measured: ", 1)[1]
        elif current and line.startswith("               lift:     "):
            current["lift"] = line.split("lift:     ", 1)[1]
    if not criteria or any(not item["measured"] for item in criteria):
        raise AuditBroken("the report has no complete criterion population")
    for item in criteria:
        if (item["id"] == "Q0.5" and item["verdict"] == "PASS" and
                not re.search(r"checksum unchanged \([0-9a-f]{12}\)", item["measured"])):
            raise AuditBroken("Q0.5 reported PASS without a non-empty tree checksum")
        found_metric = metric(item["id"], item["measured"])
        if found_metric:
            item["metric"] = found_metric
    return {
        "reported_commit": header.group(1), "upstream_version": header.group(2),
        "mode": header.group(3), "denominators": denominators, "criteria": criteria,
    }


def force_verdict(item, verdict, reason):
    if item["verdict"] != verdict:
        item.setdefault("upstream_verdict", item["verdict"])
        item["verdict"] = verdict
    item.setdefault("wrapper_findings", []).append(reason)


def replace_measurement(item, measured, found_metric, verdict, reason):
    item.setdefault("upstream_measured", item["measured"])
    item["measured"] = measured
    item["metric"] = found_metric
    force_verdict(item, verdict, reason)


def exact_path_member(relative, roots):
    return any(relative == root or relative.startswith(root + "/") for root in roots)


def isolated_go_environment(temporary):
    environment = os.environ.copy()
    environment["GOWORK"] = "off"
    build_cache = Path(temporary) / "gocache"
    build_cache.mkdir(exist_ok=True)
    environment["GOCACHE"] = str(build_cache)
    if not environment.get("GOMODCACHE"):
        probe = subprocess.run(
            ["go", "env", "GOMODCACHE"], stdout=subprocess.PIPE, stderr=subprocess.PIPE,
            text=True, check=False,
        )
        if probe.returncode != 0 or not probe.stdout.strip():
            raise AuditBroken("go env GOMODCACHE returned no path")
        environment["GOMODCACHE"] = probe.stdout.strip()
    for name in ("HOME", "GOTMPDIR"):
        location = Path(temporary) / name.lower()
        location.mkdir(exist_ok=True)
        environment[name] = str(location)
    return environment


def decode_json_stream(value):
    decoder = json.JSONDecoder()
    offset = 0
    result = []
    while offset < len(value):
        while offset < len(value) and value[offset].isspace():
            offset += 1
        if offset == len(value):
            break
        item, offset = decoder.raw_decode(value, offset)
        result.append(item)
    return result


def selected_file_digest(repo, paths):
    digest = hashlib.sha256()
    for relative in sorted(paths):
        digest.update(relative.encode("utf-8"))
        digest.update(b"\0")
        digest.update((Path(repo) / relative).read_bytes())
        digest.update(b"\0")
    return digest.hexdigest()


def go_build_context(repo):
    root = Path(repo).resolve()
    with tempfile.TemporaryDirectory(prefix="quality-go-context-") as temporary:
        environment = isolated_go_environment(temporary)
        version = subprocess.run(
            ["go", "version"], cwd=root, env=environment, stdout=subprocess.PIPE,
            stderr=subprocess.PIPE, text=True, check=False,
        )
        if version.returncode != 0 or not version.stdout.strip():
            raise AuditBroken("go version failed: " + version.stderr.strip())
        environment_run = subprocess.run(
            ["go", "env", "-json", "GOOS", "GOARCH", "CGO_ENABLED", "GOENV", "GOFLAGS", "GOWORK",
             *BUILD_SELECTOR_ENV],
            cwd=root, env=environment, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
            text=True, check=False,
        )
        if environment_run.returncode != 0:
            raise AuditBroken("go env for build context failed: " + environment_run.stderr.strip())
        try:
            selected_environment = json.loads(environment_run.stdout)
        except json.JSONDecodeError as error:
            raise AuditBroken("go env returned invalid JSON: " + str(error))
        goflags = selected_environment.get("GOFLAGS", "")
        try:
            flag_tokens = shlex.split(goflags)
        except ValueError as error:
            raise AuditBroken("GOFLAGS is malformed: " + str(error))
        forbidden = {
            "-C", "-modfile", "-overlay", "-toolexec",
            "-args", "-count", "-exec", "-failfast", "-list", "-n", "-run",
            "-short", "-shuffle", "-skip",
        }
        rejected = None
        for token in flag_tokens:
            normalized = "-" + token[2:] if token.startswith("--") else token
            name = normalized.split("=", 1)[0]
            if name in forbidden or name.startswith("-test."):
                rejected = token
                break
        if rejected:
            raise AuditBroken("GOFLAGS may bypass or desynchronize the audit: " + rejected)
        listed = subprocess.run(
            ["go", "list", "-e", "-json", "./..."], cwd=root, env=environment,
            stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=False,
        )
        if listed.returncode != 0:
            raise AuditBroken("go list could not select the build population: " + listed.stderr.strip())
    try:
        packages = decode_json_stream(listed.stdout)
    except json.JSONDecodeError as error:
        raise AuditBroken("go list returned invalid JSON: " + str(error))
    production = set()
    tests = set()
    for package in packages:
        directory = Path(package.get("Dir", ""))
        try:
            relative_directory = directory.resolve(strict=True).relative_to(root)
        except (FileNotFoundError, ValueError, OSError):
            continue
        for field, destination in (
                ("GoFiles", production), ("CgoFiles", production),
                ("TestGoFiles", tests), ("XTestGoFiles", tests)):
            for name in package.get(field) or []:
                candidate = directory / name
                try:
                    resolved = candidate.resolve(strict=True)
                    resolved.relative_to(root)
                except (FileNotFoundError, ValueError, OSError) as error:
                    raise AuditBroken("selected Go file is missing or outside the repository: {} ({})".format(
                        candidate, error,
                    ))
                relative = (relative_directory / name).as_posix()
                if exact_path_member(relative, EXCLUDED_REPOSITORY_ROOTS):
                    continue
                if candidate.is_symlink() or not candidate.is_file():
                    raise AuditBroken("selected Go file is not a regular non-symlink file: " + relative)
                destination.add(relative)
    overlap = production & tests
    if overlap:
        raise AuditBroken("go list selected files as both production and test: " + ", ".join(sorted(overlap)))
    manifest = {"production": sorted(production), "tests": sorted(tests)}
    metadata = {
        "version": version.stdout.strip(),
        "goos": selected_environment.get("GOOS"), "goarch": selected_environment.get("GOARCH"),
        "cgo_enabled": selected_environment.get("CGO_ENABLED"),
        "goenv": environment.get("GOENV", ""),
        "goflags": goflags,
        "gowork": selected_environment.get("GOWORK"),
        "production_files": {
            "count": len(production), "sha256": selected_file_digest(root, production),
        },
        "test_files": {"count": len(tests), "sha256": selected_file_digest(root, tests)},
    }
    metadata.update({name.lower(): selected_environment.get(name, "") for name in BUILD_SELECTOR_ENV})
    return {"manifest": manifest, "metadata": metadata}


def verify_selected_file_identity(repo, build_context):
    for population, metadata_name in (("production", "production_files"), ("tests", "test_files")):
        paths = build_context["manifest"][population]
        try:
            digest = selected_file_digest(repo, paths)
        except OSError as error:
            raise AuditBroken("selected Go file changed or disappeared during audit: " + str(error))
        expected = build_context["metadata"][metadata_name]
        if len(paths) != expected["count"] or digest != expected["sha256"]:
            raise AuditBroken("selected {} Go population changed during audit".format(population))


def verify_build_context_stable(repo, build_context):
    verify_selected_file_identity(repo, build_context)
    refreshed = go_build_context(repo)
    if refreshed != build_context:
        raise AuditBroken("Go build context or selected file population changed during audit")


def scan_go(repo, mode, build_context):
    scanner_source = Path(__file__).with_name("go-callscan.go.src")
    if not scanner_source.is_file():
        raise AuditBroken("missing project-local Go call scanner")
    with tempfile.TemporaryDirectory(prefix="quality-callscan-") as temporary:
        scanner = Path(temporary) / "go-callscan.go"
        manifest = Path(temporary) / "build-manifest.json"
        scanner.write_bytes(scanner_source.read_bytes())
        manifest.write_text(json.dumps(build_context["manifest"], sort_keys=True), encoding="utf-8")
        run = subprocess.run(
            ["go", "run", str(scanner), "--repo", str(Path(repo).resolve()), "--mode", mode,
             "--manifest", str(manifest)],
            stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True,
            env=isolated_go_environment(temporary), check=False,
        )
    if run.returncode != 0:
        detail = run.stderr.strip().splitlines()
        raise AuditBroken("Go call scanner failed: " + (detail[0] if detail else "no diagnostic"))
    try:
        findings = json.loads(run.stdout)
    except json.JSONDecodeError as error:
        raise AuditBroken("Go call scanner returned invalid JSON: " + str(error))
    if not isinstance(findings, list):
        raise AuditBroken("Go call scanner returned a non-list population")
    return findings


def run_q06_contract(repo):
    contract_source = Path(__file__).with_name("q06-contract_test.go.src")
    support = Path(repo) / "internal" / "testutil"
    if not contract_source.is_file() or not support.is_dir():
        return False
    with tempfile.TemporaryDirectory(prefix="quality-q06-contract-") as temporary:
        workspace = Path(temporary)
        root = workspace / "repository"
        package = root / "testutil"
        package.mkdir(parents=True)
        (root / "go.mod").write_text("module example.invalid/qualityq06\n\ngo 1.18\n", encoding="utf-8")
        for source in support.glob("*.go"):
            if not source.name.endswith("_test.go"):
                shutil.copy2(source, package / source.name)
        shutil.copy2(contract_source, package / "quality_audit_contract_test.go")
        run = subprocess.run(
            ["go", "test", "./testutil", "-count=1", "-run", "^TestQualityAuditSafeWriterContract$"],
            cwd=root, env=isolated_go_environment(temporary), stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT, text=True, check=False,
        )
        return run.returncode == 0


def run_q06_repository_controls(repo, names):
    expression = "^(?:{})$".format("|".join(re.escape(value) for value in names))
    with tempfile.TemporaryDirectory(prefix="quality-q06-repository-controls-") as temporary:
        run = subprocess.run(
            ["go", "test", "-json", "./internal/testutil", "-count=1", "-run", expression],
            cwd=repo, env=isolated_go_environment(temporary), stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT, text=True, check=False,
        )
    if run.returncode != 0:
        return []
    started = set()
    passed = set()
    for line in run.stdout.splitlines():
        try:
            event = json.loads(line)
        except json.JSONDecodeError:
            continue
        test_name = event.get("Test")
        if test_name not in names:
            continue
        if event.get("Action") == "run":
            started.add(test_name)
        elif event.get("Action") == "pass":
            passed.add(test_name)
    return sorted(started & passed)


def supplement_q06(parsed, repo, build_context):
    item = next((value for value in parsed["criteria"] if value["id"] == "Q0.6"), None)
    if not item:
        return
    structure = scan_go(repo, "q06-structure", build_context)
    structure_counts = {
        name: sum(value["kind"] == name for value in structure)
        for name in (
            "central_guard_definition", "central_guard_call", "safe_write_definition",
            "safe_copy_definition", "guarded_safe_write", "guarded_safe_copy",
            "repository_guard_control", "repository_writer_control",
        )
    }
    test_writes = scan_go(repo, "test-writes", build_context)
    write_calls = sum(value["kind"] == "safe_write" for value in test_writes)
    copy_calls = sum(value["kind"] == "safe_copy" for value in test_writes)
    unsafe_calls = [value for value in test_writes if value["kind"] == "unsafe_test_write"]
    repository_controls_present = (
        structure_counts["repository_guard_control"] + structure_counts["repository_writer_control"]
    )
    control_names = (
        "TestWouldLeakIntoRepositoryRejectsWorkingTreeAndAcceptsTempDir",
        "TestSafeFixtureWritersRefuseRepositoryPaths",
    )
    passed_controls = run_q06_repository_controls(repo, control_names) if repository_controls_present == 2 else []
    repository_controls_passed = len(passed_controls)
    contract_passed = int(run_q06_contract(repo))
    guard_definitions = structure_counts["central_guard_definition"]
    guard_calls = structure_counts["central_guard_call"]
    writer_definitions = structure_counts["safe_write_definition"] + structure_counts["safe_copy_definition"]
    guarded_writers = int(structure_counts["guarded_safe_write"] > 0) + int(
        structure_counts["guarded_safe_copy"] > 0
    )
    test_population = parsed["denominators"]["test_functions"]
    precondition = all((
        test_population > 0, guard_definitions == 1, guard_calls > 0,
        writer_definitions == 2, guarded_writers == 2, write_calls > 0, copy_calls > 0,
        not unsafe_calls, repository_controls_passed == 2, contract_passed == 1,
    ))
    found_metric = {
        "name": "skipped_tests", "value": parsed["denominators"]["skipped_tests"],
        "direction": "max", "population": test_population, "precondition_met": precondition,
        "central_guard_definitions": guard_definitions, "central_guard_call_sites": guard_calls,
        "safe_writer_definitions": writer_definitions, "guarded_safe_writers": guarded_writers,
        "safe_writer_call_sites": write_calls + copy_calls, "write_safe_call_sites": write_calls,
        "copy_safe_call_sites": copy_calls, "unsafe_direct_test_writes": len(unsafe_calls),
        "unsafe_direct_test_write_sites": unsafe_calls,
        "repository_control_tests_expected": 2,
        "repository_control_tests_passed": repository_controls_passed,
        "repository_control_tests_run_and_passed": passed_controls,
        "contract_tests_expected": 1, "contract_tests_passed": contract_passed,
    }
    verdict = "PASS" if precondition and found_metric["value"] == 0 else "FAIL"
    measured = (
        "central guard {}/1; guarded safe writers {}/2; safe-writer call sites {} "
        "(write {}, copy {}); repository controls {}/2; audit contract {}/1; "
        "unsafe direct test writes {}; {} skipped tests of {}"
    ).format(
        guard_definitions, guarded_writers, write_calls + copy_calls, write_calls, copy_calls,
        repository_controls_passed, contract_passed, len(unsafe_calls), found_metric["value"],
        test_population,
    )
    replace_measurement(item, measured, found_metric, verdict, "project-local semantic leak-guard contract")


def supplement_q13(parsed, repo, inventory, build_context):
    item = next((value for value in parsed["criteria"] if value["id"] == "Q1.3"), None)
    if not item:
        return
    adapters = [
        normalize_inventory_path(entry["value"], "adapter at inventory line {}".format(entry["line"]))
        for entry in inventory["adapters"]
    ]
    if len(adapters) != len(set(adapters)):
        raise AuditBroken("inventory contains duplicate adapter paths")
    root = Path(repo).resolve()
    selected_package_directories = {
        PurePosixPath(value).parent.as_posix()
        for value in build_context["manifest"]["production"]
    }
    adapter_receipts = []
    for value in adapters:
        candidate, issue = inspect_path_components(root, value)
        if issue:
            status = issue
        elif not candidate.is_dir():
            status = "not_directory"
        elif value not in selected_package_directories:
            status = "no_selected_production_package"
        else:
            status = "valid"
        adapter_receipts.append({"path": value, "status": status})
    calls = scan_go(repo, "external", build_context)
    violations = [
        call for call in calls
        if PurePosixPath(call["path"]).parent.as_posix() not in adapters
    ]
    invalid_adapters = [value for value in adapter_receipts if value["status"] != "valid"]
    missing_adapters = [value["path"] for value in invalid_adapters if value["status"] == "missing"]
    found_metric = {
        "name": "direct_external_call_sites", "value": len(violations), "direction": "max",
        "population": len(calls), "declared_adapters": len(adapters), "adapter_paths": adapters,
        "adapter_receipts": adapter_receipts, "invalid_adapter_paths": invalid_adapters,
        "missing_adapter_paths": missing_adapters, "test_support_paths": list(TEST_SUPPORT_PATHS),
        "violations": violations,
    }
    if not adapters:
        verdict = "UNMEASURABLE"
    elif invalid_adapters:
        verdict = "FAIL"
    elif not calls:
        verdict = "UNMEASURABLE"
    else:
        verdict = "PASS" if not violations else "FAIL"
    measured = (
        "{} direct external call sites outside {} exact adapters of {} production call sites; "
        "{} adapter paths invalid"
    ).format(
        len(violations), len(adapters), len(calls), len(invalid_adapters),
    )
    replace_measurement(item, measured, found_metric, verdict, "project-local import-aware exact-path scan")


def parse_subjects(inventory):
    subjects = []
    names = set()
    for entry in inventory["subjects"]:
        match = re.fullmatch(r"([^=]+)=([^:]+):(.+)", entry["value"])
        if not match:
            raise AuditBroken("malformed subject inventory entry at line {}".format(entry["line"]))
        name = match.group(1).strip()
        sources = [value.strip() for value in match.group(2).split(",")]
        harness = match.group(3).strip()
        if not name:
            raise AuditBroken("empty subject name at inventory line {}".format(entry["line"]))
        if name in names:
            raise AuditBroken("duplicate subject name in inventory: " + name)
        names.add(name)
        subjects.append({"name": name, "sources": sources, "harness": harness, "line": entry["line"]})
    return subjects


def normalized_subject_path(raw, allow_dot):
    if allow_dot and raw == ".":
        return raw
    try:
        return normalize_inventory_path(raw, "subject path")
    except AuditBroken:
        return None


def inspect_path_components(root, relative):
    current = root
    for part in PurePosixPath(relative).parts:
        try:
            names = {entry.name for entry in current.iterdir()}
        except (FileNotFoundError, NotADirectoryError, PermissionError, OSError):
            return None, "missing"
        if part not in names:
            candidate = current / part
            return None, "case_mismatch" if candidate.exists() else "missing"
        current = current / part
        try:
            if current.is_symlink():
                return None, "symlink"
        except OSError:
            return None, "missing"
    try:
        resolved = current.resolve(strict=True)
        resolved.relative_to(root)
    except FileNotFoundError:
        return None, "missing"
    except (ValueError, OSError):
        return None, "outside_repository"
    return current, None


def production_go_files(source):
    if source.is_file():
        return 1 if source.suffix == ".go" and not source.name.endswith("_test.go") else 0
    if not source.is_dir():
        return 0
    count = 0
    excluded = {".git", ".quality", "target", "vendor"}
    for directory, names, files in os.walk(source, followlinks=False):
        directory_path = Path(directory)
        names[:] = sorted(
            name for name in names
            if name not in excluded and not (directory_path / name).is_symlink()
        )
        for name in sorted(files):
            candidate = directory_path / name
            if name.endswith(".go") and not name.endswith("_test.go") and not candidate.is_symlink() and candidate.is_file():
                count += 1
    return count


def inspect_subject_source(root, raw, duplicate):
    receipt = {"path": raw, "kind": None, "production_go_files": 0, "status": "invalid_path"}
    relative = normalized_subject_path(raw, allow_dot=True)
    if relative is None:
        return receipt
    receipt["path"] = relative
    if duplicate:
        receipt["status"] = "duplicate"
        return receipt
    if relative != "." and exact_path_member(relative, EXCLUDED_REPOSITORY_ROOTS):
        receipt["status"] = "excluded_root"
        return receipt
    candidate, issue = inspect_path_components(root, relative)
    if issue:
        receipt["status"] = issue
        return receipt
    receipt["kind"] = "file" if candidate.is_file() else ("directory" if candidate.is_dir() else None)
    receipt["production_go_files"] = production_go_files(candidate)
    receipt["status"] = "valid" if receipt["production_go_files"] > 0 else "no_production_go"
    return receipt


def inspect_subject_harness(root, raw):
    receipt = {"path": raw, "status": "invalid_path"}
    relative = normalized_subject_path(raw, allow_dot=False)
    if relative is None:
        return receipt
    receipt["path"] = relative
    candidate, issue = inspect_path_components(root, relative)
    if issue:
        receipt["status"] = issue
        return receipt
    try:
        details = candidate.lstat()
    except OSError:
        receipt["status"] = "missing"
        return receipt
    if not stat.S_ISREG(details.st_mode):
        receipt["status"] = "not_regular"
    elif not os.access(candidate, os.X_OK):
        receipt["status"] = "not_executable"
    else:
        receipt["status"] = "valid"
    return receipt


def supplement_q21(parsed, repo, inventory):
    item = next((value for value in parsed["criteria"] if value["id"] == "Q2.1"), None)
    if not item:
        return
    root = Path(repo).resolve()
    subjects = parse_subjects(inventory)
    receipts = []
    for subject in subjects:
        counts = {value: subject["sources"].count(value) for value in subject["sources"]}
        sources = [inspect_subject_source(root, value, counts[value] > 1) for value in subject["sources"]]
        harness = inspect_subject_harness(root, subject["harness"])
        issues = [
            {"code": value["status"], "path": value["path"], "kind": "source"}
            for value in sources if value["status"] != "valid"
        ]
        if harness["status"] != "valid":
            issues.append({"code": harness["status"], "path": harness["path"], "kind": "harness"})
        receipts.append({
            "name": subject["name"], "declared_sources": subject["sources"],
            "source_roots": sources, "harness": harness,
            "status": "covered" if not issues else "uncovered", "issues": issues,
        })
    covered = sum(value["status"] == "covered" for value in receipts)
    invalid = [value for value in receipts if value["status"] != "covered"]
    source_precondition = all(
        source["status"] == "valid"
        for subject in receipts for source in subject["source_roots"]
    )
    found_metric = {
        "name": "valid_subject_harness_bindings", "numerator": covered,
        "denominator": len(subjects), "direction": "min_ratio", "subjects": receipts,
        "invalid_subject_bindings": invalid, "source_precondition_met": source_precondition,
    }
    verdict = "UNMEASURABLE" if not subjects else ("PASS" if covered == len(subjects) else "FAIL")
    measured = "{} of {} declared subjects have real production roots and a matching executable harness".format(
        covered, len(subjects),
    )
    replace_measurement(item, measured, found_metric, verdict, "project-local subject-to-harness validation")


def markdown_population(repo):
    return sum(
        not exact_path_member(path.relative_to(repo).as_posix(), (".git", "target", "vendor"))
        for path in Path(repo).rglob("*.md")
    )


def enrich_metrics(parsed):
    den = parsed["denominators"]
    for item in parsed["criteria"]:
        found_metric = item.get("metric")
        if not found_metric:
            continue
        if item["id"] == "Q1.2":
            found_metric["population"] = den["go_files"]
        elif item["id"] == "Q3.4":
            found_metric["population"] = den["markdown_files"]


def enforce_population_guards(parsed):
    den = parsed["denominators"]
    by_id = {item["id"]: item for item in parsed["criteria"]}
    zero_denominators = {
        "Q0.1": min(den["packages"], den["go_files"]),
        "Q0.2": den["test_functions"], "Q0.3": den["go_files"],
        "Q0.4": den["test_functions"], "Q0.5": den["test_functions"],
        "Q0.6": den["test_functions"], "Q0.8": den["scripts"],
        "Q1.1": den["packages"], "Q1.2": den["go_files"], "Q1.8": den["test_functions"],
        "Q2.2": den["mutation_harnesses"], "Q2.3": den["mutation_harnesses"],
        "Q2.4": den["mutation_harnesses"], "Q2.5": den["declared_features"],
        "Q2.6": den["acceptance_scripts"], "Q2.7": den["acceptance_scripts"],
        "Q2.8": den["acceptance_scripts"], "Q2.9": den["acceptance_scripts"],
        "Q2.10": den["acceptance_scripts"], "Q3.4": den["markdown_files"],
    }
    for criterion_id, population in zero_denominators.items():
        item = by_id.get(criterion_id)
        if item and population == 0:
            force_verdict(item, "UNMEASURABLE", "population is 0, so PASS is impossible")
    for criterion_id in ("Q1.4", "Q2.1"):
        item = by_id.get(criterion_id)
        found_metric = item.get("metric") if item else None
        if not found_metric:
            continue
        total = found_metric.get("denominator")
        covered = found_metric.get("numerator")
        if total == 0:
            force_verdict(item, "UNMEASURABLE", "declared population is 0, so PASS is impossible")
        elif covered < total:
            force_verdict(item, "FAIL", "declared population is not fully covered")


def metric_precondition(criterion_id, found_metric, current):
    if criterion_id == "Q0.6":
        population_valid = found_metric.get("population", 0) > 0
        return population_valid and (not current or found_metric.get("precondition_met") is True)
    if criterion_id == "Q2.1":
        return (
            found_metric.get("denominator", 0) > 0 and
            (not current or found_metric.get("source_precondition_met") is True)
        )
    if found_metric.get("direction") == "min_ratio":
        return found_metric.get("denominator", 0) > 0
    if criterion_id == "Q1.3":
        return (
            found_metric.get("population", 0) > 0 and
            found_metric.get("declared_adapters", 0) > 0 and
            (not current or not found_metric.get("invalid_adapter_paths"))
        )
    if "population" in found_metric:
        return found_metric["population"] > 0
    return isinstance(found_metric.get("value"), int)


def compare_metric(current, baseline):
    if current.get("name") != baseline.get("name") or current.get("direction") != baseline.get("direction"):
        raise AuditBroken("baseline and current ratchet metrics are incompatible")
    if current["direction"] == "max":
        if current["value"] < baseline["value"]:
            return "improved"
        if current["value"] == baseline["value"]:
            return "held"
        return "regressed"
    if current["direction"] == "min_ratio":
        if current["denominator"] == 0 or baseline.get("denominator", 0) == 0:
            return None
        current_scaled = current["numerator"] * baseline["denominator"]
        baseline_scaled = baseline["numerator"] * current["denominator"]
        if current_scaled > baseline_scaled:
            return "improved"
        if current_scaled == baseline_scaled:
            return "held"
        return "regressed"
    raise AuditBroken("unknown ratchet direction")


def apply_baseline(parsed, baseline_path, module, current_tool, current_inventory):
    if not baseline_path:
        return {"status": "not_supplied", "compared": 0, "improved": 0, "held": 0, "regressed": 0}
    try:
        baseline_raw = Path(baseline_path).read_bytes()
        baseline = json.loads(baseline_raw)
    except (OSError, json.JSONDecodeError) as error:
        raise AuditBroken("cannot read baseline JSON: " + str(error))
    if baseline.get("schema_version") != SCHEMA_VERSION:
        raise AuditBroken("baseline schema version does not match")
    if baseline.get("repository", {}).get("module") != module:
        raise AuditBroken("baseline belongs to a different Go module")
    baseline_inventory = baseline.get("inventory", {})
    for field in ("present", "sha256"):
        if baseline_inventory.get(field) != current_inventory.get(field):
            raise AuditBroken("baseline inventory identity mismatch: " + field)
    baseline_tree = baseline.get("repository", {}).get("tree", {})
    if baseline_tree.get("measurement_clean") is not True:
        raise AuditBroken("baseline was not measured from an accepted clean tree")
    baseline_tool = baseline.get("tool", {})
    if baseline_tool.get("partial") is not False:
        raise AuditBroken("baseline must be an authoritative full scorecard")
    for field in BASELINE_TOOL_FIELDS:
        if field not in baseline_tool or baseline_tool[field] != current_tool.get(field):
            raise AuditBroken("baseline tool identity mismatch: " + field)
    baseline_build = baseline_tool.get("go_build", {})
    current_build = current_tool.get("go_build", {})
    for field in BASELINE_BUILD_FIELDS:
        if field not in baseline_build or baseline_build[field] != current_build.get(field):
            raise AuditBroken("baseline Go build context mismatch: " + field)
    baseline_criteria = baseline.get("criteria", [])
    if not isinstance(baseline_criteria, list) or any(not isinstance(item, dict) for item in baseline_criteria):
        raise AuditBroken("baseline criteria must be a list of objects")
    baseline_ids = [item.get("id") for item in baseline_criteria]
    if len(baseline_ids) != len(set(baseline_ids)):
        raise AuditBroken("baseline contains duplicate criterion IDs")
    if set(baseline_ids) != EXPECTED_IDS:
        missing = sorted(EXPECTED_IDS - set(baseline_ids))
        extra = sorted(set(baseline_ids) - EXPECTED_IDS)
        raise AuditBroken("full baseline criterion mismatch; missing={} extra={}".format(missing, extra))
    old = {item["id"]: item for item in baseline_criteria}
    baseline_sha256 = sha256_bytes(baseline_raw)
    counts = {"improved": 0, "held": 0, "regressed": 0}
    compared = 0
    selected = [item for item in parsed["criteria"] if item["id"] in RATCHET_IDS]
    if not selected:
        return {
            "status": "not_applicable", "baseline_commit": baseline.get("repository", {}).get("commit", "unknown"),
            "baseline_sha256": baseline_sha256, "selected": 0, "current_not_comparable": 0,
            "compared": 0, **counts,
        }
    current_not_comparable = 0
    for item in selected:
        if "metric" not in item:
            raise AuditBroken("selected ratchet criterion has no structured metric: " + item["id"])
        previous = old.get(item["id"], {})
        if "metric" not in previous:
            raise AuditBroken("baseline has no metric for selected ratchet criterion: " + item["id"])
        if not metric_precondition(item["id"], previous["metric"], current=False):
            raise AuditBroken("baseline ratchet metric has an empty or invalid population: " + item["id"])
        if not metric_precondition(item["id"], item["metric"], current=True):
            current_not_comparable += 1
            continue
        outcome = compare_metric(item["metric"], previous["metric"])
        if outcome is None:
            current_not_comparable += 1
            continue
        compared += 1
        counts[outcome] += 1
        item["ratchet"] = {"baseline": previous["metric"], "current": item["metric"], "outcome": outcome}
        if outcome == "regressed":
            force_verdict(item, "FAIL", "ratchet metric regressed from the stored baseline")
        else:
            force_verdict(item, "PASS", "ratchet metric {} against the stored baseline".format(outcome))
    if compared == 0:
        return {
            "status": "not_comparable", "reason": "selected current ratchet population is empty or invalid",
            "baseline_commit": baseline.get("repository", {}).get("commit", "unknown"),
            "baseline_sha256": baseline_sha256,
            "selected": len(selected), "current_not_comparable": current_not_comparable,
            "compared": 0, **counts,
        }
    return {
        "status": "compared", "baseline_commit": baseline.get("repository", {}).get("commit", "unknown"),
        "baseline_sha256": baseline_sha256, "selected": len(selected),
        "current_not_comparable": current_not_comparable, "compared": compared, **counts,
    }


def git_output(repo, *arguments):
    run = subprocess.run(
        ["git", "-C", str(repo), *arguments], stdout=subprocess.PIPE, stderr=subprocess.PIPE,
        text=True, check=False,
    )
    if run.returncode != 0:
        raise AuditBroken("git {} failed: {}".format(" ".join(arguments), run.stderr.strip()))
    return run.stdout


def measurement_dirty_paths(repo):
    root = Path(repo).resolve()
    listed = subprocess.run(
        ["git", "-C", str(root), "ls-tree", "-r", "-z", "HEAD"],
        stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=False,
    )
    if listed.returncode != 0:
        raise AuditBroken("git ls-tree failed: " + listed.stderr.decode("utf-8", "replace").strip())
    tracked = {}
    for record in listed.stdout.split(b"\0"):
        if not record:
            continue
        metadata, raw_path = record.split(b"\t", 1)
        mode, object_type, object_id = metadata.decode("ascii").split()
        path = os.fsdecode(raw_path)
        if object_type == "blob" and not measurement_output_path(path):
            tracked[path] = (mode, object_id)

    working = {}
    for directory, names, filenames in os.walk(root, followlinks=False):
        directory_path = Path(directory)
        retained = []
        for name in names:
            path = directory_path / name
            relative = path.relative_to(root).as_posix()
            if measurement_output_path(relative):
                continue
            if path.is_symlink():
                working[relative] = path
            else:
                retained.append(name)
        names[:] = retained
        for name in filenames:
            path = directory_path / name
            relative = path.relative_to(root).as_posix()
            if not measurement_output_path(relative):
                working[relative] = path

    dirty = set(tracked) ^ set(working)
    for relative in set(tracked) & set(working):
        path = working[relative]
        mode, object_id = tracked[relative]
        expected_symlink = mode == "120000"
        if path.is_symlink() != expected_symlink:
            dirty.add(relative)
            continue
        try:
            contents = os.fsencode(os.readlink(path)) if expected_symlink else path.read_bytes()
        except OSError as error:
            raise AuditBroken("cannot read measured worktree path {}: {}".format(relative, error))
        stored = subprocess.run(
            ["git", "-C", str(root), "cat-file", "blob", object_id],
            stdout=subprocess.PIPE, stderr=subprocess.PIPE, check=False,
        )
        if stored.returncode != 0:
            raise AuditBroken("git cat-file failed for {}".format(relative))
        executable = bool(path.stat().st_mode & 0o111) if not expected_symlink else False
        expected_executable = mode == "100755"
        if contents != stored.stdout or executable != expected_executable:
            dirty.add(relative)
    return sorted(dirty)


def measurement_output_path(relative):
    return exact_path_member(relative, (".git", "target/quality-audit"))


def tree_identity(repo, commit, overlay_sha256):
    head = git_output(repo, "rev-parse", "HEAD").strip()
    if head != commit:
        raise AuditBroken("--commit does not equal repository HEAD")
    commit_tree = git_output(repo, "rev-parse", "HEAD^{tree}").strip()
    dirty_paths = measurement_dirty_paths(repo)
    git_clean = not dirty_paths
    accepted_overlay = False
    if overlay_sha256:
        if not re.fullmatch(r"[0-9a-f]{64}", overlay_sha256):
            raise AuditBroken("--inventory-overlay-sha256 must be a lowercase SHA-256")
        accepted_overlay = dirty_paths == [".quality/inventory"]
        if not accepted_overlay:
            raise AuditBroken("inventory overlay mode permits only untracked .quality/inventory")
    return {
        "commit_tree": commit_tree, "git_clean": git_clean, "measurement_clean": git_clean or accepted_overlay,
        "dirty_paths": dirty_paths,
        "status_sha256": sha256_bytes(("\0".join(dirty_paths) + "\0").encode("utf-8")),
        "inventory_overlay_sha256": overlay_sha256 if accepted_overlay else None,
    }


def has_text(value):
    return isinstance(value, str) and bool(value.strip())


def has_number(value):
    return has_text(value) and bool(re.search(r"\d", value))


def normalize_classification(value):
    return re.sub(r"\s+", " ", value.strip().lower().replace("_", " ").replace("-", " "))


def canonical_json_sha256(value):
    raw = json.dumps(
        value, ensure_ascii=False, separators=(",", ":"), sort_keys=True,
    ).encode("utf-8")
    return sha256_bytes(raw)


def is_integer(value, minimum=0):
    return isinstance(value, int) and not isinstance(value, bool) and value >= minimum


def validate_receipt_population(criterion, population):
    if not isinstance(population, list) or not population:
        return "evidence population must be a non-empty list"
    if any(not isinstance(value, dict) for value in population):
        return "evidence population entries must be objects"
    subjects = [value.get("subject") for value in population]
    if any(not has_text(value) for value in subjects):
        return "every evidence population entry needs a subject"
    if len(subjects) != len(set(value.strip() for value in subjects)):
        return "evidence population subjects must be unique"

    required = {
        "Q1.6": {
            "subject", "dependencies", "default_doubles", "argument_recorders",
            "dependency_struct_passed_whole",
        },
        "Q1.7": {
            "subject", "partial_failures", "content_assertions", "exit_only_assertions",
        },
        "Q1.9": {
            "subject", "iterated_collections", "empty_population_assertions",
        },
        "Q2.4": {"subject", "declared", "killed", "survived", "unusable"},
        "Q2.8": {
            "subject", "comparison", "control_magnitude", "treatment_magnitude",
            "expected_relation",
        },
        "Q2.9": {
            "subject", "bad_input_exit", "produced_artifacts", "read_only_artifacts_changed",
        },
    }[criterion]
    for value in population:
        if set(value) != required:
            return "{} evidence fields do not match its criterion contract".format(criterion)
        if criterion == "Q1.6":
            dependencies = value["dependencies"]
            if (not is_integer(dependencies, 1) or
                    value["default_doubles"] != dependencies or
                    value["argument_recorders"] != dependencies or
                    value["dependency_struct_passed_whole"] is not True):
                return "Q1.6 PASS requires defaults and argument recorders for every dependency and whole-struct passing"
        elif criterion == "Q1.7":
            failures = value["partial_failures"]
            if (not is_integer(failures, 1) or
                    not is_integer(value["content_assertions"], failures) or
                    value["exit_only_assertions"] != 0):
                return "Q1.7 PASS requires content assertions for every partial failure and no exit-only assertions"
        elif criterion == "Q1.9":
            collections = value["iterated_collections"]
            if (not is_integer(collections, 1) or
                    not is_integer(value["empty_population_assertions"], collections)):
                return "Q1.9 PASS requires an empty-population assertion for every iterated collection"
        elif criterion == "Q2.4":
            declared = value["declared"]
            if (not is_integer(declared, 1) or value["killed"] != declared or
                    value["survived"] != 0 or value["unusable"] != 0):
                return "Q2.4 PASS requires declared == killed with zero survived and unusable mutations"
        elif criterion == "Q2.8":
            control = value["control_magnitude"]
            treatment = value["treatment_magnitude"]
            relation = value["expected_relation"]

            def numeric(candidate):
                return isinstance(candidate, (int, float)) and not isinstance(candidate, bool)

            relation_holds = (
                relation == "less" and numeric(control) and numeric(treatment) and treatment < control or
                relation == "greater" and numeric(control) and numeric(treatment) and treatment > control or
                relation == "different" and numeric(control) and numeric(treatment) and treatment != control
            )
            if value["comparison"] != "magnitude" or not relation_holds:
                return "Q2.8 PASS requires a true magnitude comparison whose treatment has the declared relation"
        elif criterion == "Q2.9":
            if (not is_integer(value["bad_input_exit"]) or value["bad_input_exit"] == 0 or
                    value["produced_artifacts"] != 0 or value["read_only_artifacts_changed"] != 0):
                return "Q2.9 PASS requires non-zero bad-input exit and zero produced or changed artifacts"
    return None


def validate_criterion_receipts(evidence):
    receipts = evidence.get("criteria")
    if not isinstance(receipts, list):
        return None, "criteria must be a list"
    criteria = [value.get("criterion") for value in receipts if isinstance(value, dict)]
    if len(criteria) != len(receipts):
        return None, "criterion receipts must be objects"
    if len(criteria) != len(set(criteria)):
        return None, "criterion receipts must be unique"
    summaries = {}
    evidence_digests = []
    required_fields = {"criterion", "kind", "verdict", "evidence", "evidence_sha256"}
    for receipt in receipts:
        criterion = receipt["criterion"]
        if criterion not in MANUAL_CRITERION_KINDS:
            return None, "receipt criterion is not manually resolvable: " + str(criterion)
        if set(receipt) != required_fields:
            return None, "{} receipt fields do not match the schema".format(criterion)
        if receipt["kind"] != MANUAL_CRITERION_KINDS[criterion]:
            return None, "{} receipt kind is wrong".format(criterion)
        if receipt["verdict"] != "PASS":
            return None, "{} receipt verdict must be PASS".format(criterion)
        claim = receipt["evidence"]
        if (not isinstance(claim, dict) or set(claim) != {"command", "population"} or
                not has_text(claim.get("command"))):
            return None, "{} receipt evidence is empty or malformed".format(criterion)
        digest = receipt["evidence_sha256"]
        if not isinstance(digest, str) or not re.fullmatch(r"[0-9a-f]{64}", digest):
            return None, "{} evidence digest is not a lowercase SHA-256".format(criterion)
        if canonical_json_sha256(claim) != digest:
            return None, "{} evidence digest does not match its payload".format(criterion)
        population_error = validate_receipt_population(criterion, claim["population"])
        if population_error:
            return None, population_error
        evidence_digests.append(digest)
        summaries[criterion] = {
            "kind": receipt["kind"], "evidence_sha256": digest,
            "population": len(claim["population"]),
        }
    if len(evidence_digests) != len(set(evidence_digests)):
        return None, "criterion evidence digests must be unique"
    return summaries, None


def validate_manual(path, commit, module, tree, inventory, tool):
    if not path:
        return {"status": "not_supplied"}
    try:
        raw = Path(path).read_bytes()
    except OSError as error:
        return {"status": "invalid", "reason": str(error)}
    digest = sha256_bytes(raw)
    try:
        evidence = json.loads(raw)
    except json.JSONDecodeError as error:
        return {"status": "invalid", "reason": str(error), "sha256": digest}
    if evidence.get("schema_version") != 2:
        return {"status": "invalid", "reason": "schema_version must be 2", "sha256": digest}
    repository = evidence.get("repository", {})
    if repository.get("module") != module:
        return {"status": "invalid", "reason": "module does not match", "sha256": digest}
    if repository.get("commit") != commit:
        return {"status": "stale", "evidence_commit": repository.get("commit", "missing"), "sha256": digest}
    if not tree["measurement_clean"]:
        return {"status": "stale", "reason": "measured worktree is dirty", "sha256": digest}
    evidence_tree = repository.get("tree", {})
    if evidence_tree.get("commit_tree") != tree["commit_tree"]:
        return {"status": "stale", "reason": "commit tree does not match", "sha256": digest}
    if evidence_tree.get("status_sha256") != tree["status_sha256"]:
        return {"status": "stale", "reason": "measured tree status does not match", "sha256": digest}
    if evidence_tree.get("inventory_overlay_sha256") != tree.get("inventory_overlay_sha256"):
        return {"status": "stale", "reason": "inventory overlay does not match", "sha256": digest}
    if evidence.get("inventory") != inventory:
        return {"status": "stale", "reason": "inventory identity does not match", "sha256": digest}
    expected_instrument = {field: tool.get(field) for field in BASELINE_TOOL_FIELDS}
    if evidence.get("instrument") != expected_instrument:
        return {"status": "stale", "reason": "instrument identity does not match", "sha256": digest}
    criterion_receipts, receipt_error = validate_criterion_receipts(evidence)
    if receipt_error:
        return {"status": "invalid", "reason": receipt_error, "sha256": digest}
    findings = evidence.get("findings", {})
    seams = findings.get("seam_test")
    survivors = findings.get("survivors")
    human = findings.get("human_output")
    premise = findings.get("premise_check")
    if not isinstance(seams, list) or len(seams) != 3:
        return {"status": "invalid", "reason": "seam_test must contain exactly 3 receipts", "sha256": digest}
    if not isinstance(survivors, list) or len(survivors) != 3:
        return {"status": "invalid", "reason": "survivors must contain exactly 3 classifications", "sha256": digest}
    if len({json.dumps(value, sort_keys=True) for value in seams}) != 3:
        return {"status": "invalid", "reason": "seam receipts must be unique", "sha256": digest}
    if len({json.dumps(value, sort_keys=True) for value in survivors}) != 3:
        return {"status": "invalid", "reason": "survivor receipts must be unique", "sha256": digest}
    for receipt in seams:
        if (not isinstance(receipt, dict) or
                not all(has_text(receipt.get(key)) for key in ("command", "mutation", "unverified")) or
                not has_number(receipt.get("measured"))):
            return {"status": "invalid", "reason": "a seam receipt is empty or has no measured number", "sha256": digest}
    if len({receipt["mutation"].strip() for receipt in seams}) != 3:
        return {"status": "invalid", "reason": "seam mutations must be unique", "sha256": digest}
    for receipt in survivors:
        if (not isinstance(receipt, dict) or
                not all(has_text(receipt.get(key)) for key in ("location", "classification", "unverified")) or
                not has_number(receipt.get("measured"))):
            return {"status": "invalid", "reason": "a survivor receipt is empty or has no measured number", "sha256": digest}
        if normalize_classification(receipt["classification"]) not in SURVIVOR_CLASSES:
            return {"status": "invalid", "reason": "survivor classification is outside the documented vocabulary", "sha256": digest}
    if len({receipt["location"].strip() for receipt in survivors}) != 3:
        return {"status": "invalid", "reason": "survivor locations must be unique", "sha256": digest}
    if (not isinstance(human, dict) or
            not all(has_text(human.get(key)) for key in ("command", "unverified")) or
            not has_number(human.get("measured"))):
        return {"status": "invalid", "reason": "human_output is empty or has no measured number", "sha256": digest}
    if (not isinstance(premise, dict) or
            not all(has_text(premise.get(key)) for key in ("premise", "command", "unverified")) or
            not has_number(premise.get("measured")) or not isinstance(premise.get("held"), bool)):
        return {"status": "invalid", "reason": "premise_check is incomplete", "sha256": digest}
    return {
        "status": "valid", "commit": commit, "commit_tree": tree["commit_tree"], "sha256": digest,
        "seam_receipts": 3, "survivor_receipts": 3,
        "survivor_classifications": sorted({normalize_classification(value["classification"]) for value in survivors}),
        "criterion_receipts": len(criterion_receipts), "receipts": criterion_receipts,
        "inventory_sha256": inventory["sha256"], "instrument": expected_instrument,
    }


def apply_q39_manual_verdict(parsed, manual):
    item = next((value for value in parsed["criteria"] if value["id"] == "Q3.9"), None)
    if not item:
        return
    if item["verdict"] == "FAIL":
        if manual["status"] == "valid":
            reason = "manual evidence is valid, but the audit's automated failure takes precedence"
        elif manual["status"] in ("invalid", "stale"):
            reason = "the automated audit failed and manual evidence is " + manual["status"]
        else:
            reason = "the automated audit failed; manual evidence was not supplied"
        force_verdict(item, "FAIL", reason)
        return
    if manual["status"] == "valid":
        force_verdict(item, "PASS", "manual evidence is valid and bound to the measured clean tree")
    elif manual["status"] in ("invalid", "stale"):
        force_verdict(item, "FAIL", "manual evidence is " + manual["status"])
    else:
        force_verdict(item, "UNMEASURABLE", "manual evidence was not supplied")


def apply_manual_verdicts(parsed, manual):
    apply_q39_manual_verdict(parsed, manual)
    by_id = {value["id"]: value for value in parsed["criteria"]}
    for criterion, population_name in MANUAL_CRITERION_POPULATIONS.items():
        item = by_id.get(criterion)
        if not item:
            continue
        if manual["status"] in ("invalid", "stale"):
            if item["verdict"] != "FAIL":
                force_verdict(item, "FAIL", "manual evidence is " + manual["status"])
            continue
        if manual["status"] != "valid" or criterion not in manual["receipts"]:
            continue
        receipt = manual["receipts"][criterion]
        if item["verdict"] == "FAIL":
            force_verdict(item, "FAIL", "upstream failure takes precedence over a valid manual receipt")
            continue
        if item["verdict"] != "UNMEASURABLE":
            continue
        population = parsed["denominators"].get(population_name, 0)
        if population <= 0:
            force_verdict(
                item, "UNMEASURABLE",
                "valid manual receipt cannot resolve an empty {} population".format(population_name),
            )
            continue
        force_verdict(
            item, "PASS",
            "valid {} receipt covers {} evidence subject(s), bound by {}".format(
                receipt["kind"], receipt["population"], receipt["evidence_sha256"],
            ),
        )


def summarize(criteria):
    levels = {}
    for level, expected in EXPECTED_LEVEL_COUNTS.items():
        selected = [item for item in criteria if item["level"] == level]
        counts = {
            name: sum(item["verdict"] == name for item in selected)
            for name in ("PASS", "FAIL", "UNMEASURABLE", "N-A")
        }
        levels[level] = {"expected": expected, "measured": len(selected), **counts}
    attained = "none"
    for level in ("L0", "L1", "L2", "L3"):
        value = levels[level]
        if value["measured"] == value["expected"] and value["PASS"] + value["N-A"] == value["expected"]:
            attained = level
        else:
            break
    return levels, attained


def write_json(path, value):
    destination = Path(path)
    destination.parent.mkdir(parents=True, exist_ok=True)
    payload = json.dumps(value, indent=2, sort_keys=True) + "\n"
    descriptor, temporary = tempfile.mkstemp(prefix=".scorecard.", dir=str(destination.parent))
    try:
        with os.fdopen(descriptor, "w", encoding="utf-8") as stream:
            stream.write(payload)
        os.replace(temporary, destination)
    finally:
        if os.path.exists(temporary):
            os.unlink(temporary)


def instrument_identity(args):
    upstream = Path(args.upstream)
    directory = Path(__file__).parent
    result = {
        "upstream_sha256": sha256_file(upstream), "parser_sha256": sha256_file(__file__),
        "wrapper_sha256": sha256_file(directory / "quality-audit.sh"),
        "call_scanner_sha256": sha256_file(directory / "go-callscan.go.src"),
        "q06_contract_sha256": sha256_file(directory / "q06-contract_test.go.src"),
    }
    report_template = upstream.parent / "audit-report.template.md"
    result["report_template_sha256"] = sha256_file(report_template) if report_template.is_file() else None
    for field, digest in result.items():
        expected = os.environ.get("QUALITY_AUDIT_EXPECTED_" + field.upper())
        if expected is not None and expected != digest:
            raise AuditBroken("instrument changed before parser execution: " + field)
    return result


def tool_metadata(args, parsed, build_context, initial_identity):
    if instrument_identity(args) != initial_identity:
        raise AuditBroken("quality instrument changed during audit")
    return {
        "wrapper_version": WRAPPER_VERSION, "upstream_version": parsed["upstream_version"],
        **initial_identity,
        "mode": parsed["mode"], "partial": args.partial, "upstream_exit": args.upstream_exit,
        "go_build": build_context["metadata"],
    }


def main():
    args = parse_args()
    try:
        if args.upstream_exit not in (0, 1):
            raise AuditBroken("--upstream-exit must be 0 or 1")
        initial_instrument = instrument_identity(args)
        parsed = parse_report(args.report)
        ids = [item["id"] for item in parsed["criteria"]]
        if len(ids) != len(set(ids)):
            raise AuditBroken("the report contains duplicate criterion IDs")
        if not args.partial and set(ids) != EXPECTED_IDS:
            missing = sorted(EXPECTED_IDS - set(ids))
            extra = sorted(set(ids) - EXPECTED_IDS)
            raise AuditBroken("full report criterion mismatch; missing={} extra={}".format(missing, extra))
        if not args.commit.startswith(parsed["reported_commit"]):
            raise AuditBroken("report commit does not match the audited repository")
        module = module_name(args.repo)
        build_context = go_build_context(args.repo)
        inventory, inventory_info = inventory_snapshot(args.repo, args.inventory_overlay_sha256)
        parsed["denominators"].update(inventory_counts(inventory))
        parsed["denominators"]["markdown_files"] = markdown_population(args.repo)
        supplement_q06(parsed, args.repo, build_context)
        supplement_q13(parsed, args.repo, inventory, build_context)
        supplement_q21(parsed, args.repo, inventory)
        verify_build_context_stable(args.repo, build_context)
        enrich_metrics(parsed)
        enforce_population_guards(parsed)
        tree = tree_identity(args.repo, args.commit, args.inventory_overlay_sha256)
        current_tool = tool_metadata(args, parsed, build_context, initial_instrument)
        ratchet = apply_baseline(parsed, args.baseline, module, current_tool, inventory_info)
        manual = validate_manual(
            args.manual_evidence, args.commit, module, tree, inventory_info, current_tool,
        )
        apply_manual_verdicts(parsed, manual)
        levels, attained = summarize(parsed["criteria"])
        result = {
            "schema_version": SCHEMA_VERSION,
            "repository": {"module": module, "commit": args.commit, "tree": tree},
            "inventory": inventory_info, "tool": current_tool,
            "denominators": parsed["denominators"], "criteria": parsed["criteria"],
            "levels": levels, "attained_level": attained, "ratchet": ratchet,
            "manual_evidence": manual,
        }
        verify_inventory_identity(args.repo, inventory_info)
        write_json(args.output, result)
    except Exception as error:
        print("AUDIT BROKEN: " + str(error), file=sys.stderr)
        return 2
    failures = sum(item["verdict"] in ("FAIL", "UNMEASURABLE") for item in parsed["criteria"])
    print("structured scorecard: {} ({} non-passing criteria, {} ratchet regressions)".format(
        args.display_output or args.output, failures, ratchet.get("regressed", 0)))
    return 0 if failures == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
