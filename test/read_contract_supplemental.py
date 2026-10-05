#!/usr/bin/env python3
"""RL-01 queue/filter acceptance supplement; never creates or changes seed data.

Use the three manifests created by TestReadContractAcceptanceFixture. The frozen
read_contract_acceptance.py remains the primary compatibility test. This supplement
covers actual format-4 queues, every filter combination, unavailable format-4 titles,
and unchanged ready/Task-show reads. Expected rows come exclusively from the seeds.

Optional --evidence preserves each subprocess's argv/cwd/exit/duration and separate
stdout/stderr with hashes. --responses preserves successful new JSON response kinds
for read_contract_schema_test.py. Both destinations must be new and outside fixtures.
"""

import argparse
from datetime import datetime, timezone
import hashlib
import itertools
import json
import os
from pathlib import Path
import shlex
import subprocess
import sys
import tempfile
import time

sys.dont_write_bytecode = True
from read_contract_acceptance import require, snapshot, subset, unique_object


TASK_FIELDS = {"project": "project_id", "repo": "repo_id", "epic": "parent_epic_id"}
NEW_KINDS = {"PlyCapabilities@1", "WorkspaceProjectListReadback@1",
             "WorkspaceProjectReadback@1", "WorkspaceTaskListReadback@1"}


def sha256(data):
    return hashlib.sha256(data).hexdigest()


def json_bytes(value):
    return (json.dumps(value, ensure_ascii=False, sort_keys=True, indent=2) + "\n").encode()


def main():
    parser = argparse.ArgumentParser(description=__doc__,
                                     formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--ply", type=Path, required=True)
    parser.add_argument("--baseline", type=Path, required=True)
    parser.add_argument("--fixtures", type=Path, required=True,
                        help="primary manifest.json, beside queue-healthy and unavailable-v4")
    parser.add_argument("--evidence", type=Path)
    parser.add_argument("--responses", type=Path)
    args = parser.parse_args()
    base = args.fixtures.resolve().parent
    manifests = [args.fixtures.resolve(), base / "queue-healthy/manifest.json",
                 base / "unavailable-v4/manifest.json"]
    fixtures = [json.loads(path.read_text(), object_pairs_hook=unique_object) for path in manifests]
    for f in fixtures:
        require(f["kind"] == "PlyReadContractFixture@1", "unknown fixture kind")
        require(not any(Path(f["home"]).iterdir()), "private home must begin empty")
        for key in ("root", "home", "outside_root", "empty_root"):
            require(str(Path(f[key]).resolve()) == f[key], "fixture root must be physical")
    for binary in (args.ply, args.baseline):
        require(binary.is_absolute() and binary.is_file(), "absolute binary required")
    for destination in (args.evidence, args.responses):
        if destination is not None:
            require(destination.is_absolute(), "artifact destination must be absolute")
            require(base not in destination.resolve().parents and destination.resolve() != base,
                    "artifact destination must be outside watched fixtures")
            destination.mkdir(parents=True)

    # Deny and record named external tools during only the new registration reads.
    # Old ready/show modes retain their existing Git observation dependencies.
    guard = (args.evidence / "no-process-tools" if args.evidence else
             Path(tempfile.mkdtemp(prefix="rl01-no-process-tools-", dir=base.parent)))
    if args.evidence:
        guard.mkdir()
    attempts = guard / "attempts.txt"
    for name in ("git", "gh", "codex", "claude", "curl", "wget", "ssh"):
        file = guard / name
        file.write_text("#!/bin/sh\nprintf '%s\\n' \"$0\" >> " + shlex.quote(str(attempts)) +
                        "\nexit 97\n")
        file.chmod(0o755)

    watch = [str(base), str(args.ply), str(args.baseline)]
    before = snapshot(watch)
    calls = []
    result = {"kind": "PlyReadContractSupplement@1", "result": "failed", "human_verdict": None,
              "candidate_sha256": sha256(args.ply.read_bytes()),
              "baseline_sha256": sha256(args.baseline.read_bytes()),
              "manifests": [{"path": str(p), "sha256": sha256(p.read_bytes())} for p in manifests]}

    def run(argv, fixture, *, baseline=False, cwd=None, guarded=False):
        executable = args.baseline if baseline else args.ply
        env = {key: os.environ[key] for key in ("PATH", "TMPDIR", "LANG", "LC_ALL", "SYSTEMROOT")
               if key in os.environ}
        env.update({"HOME": fixture["home"], "XDG_CONFIG_HOME": fixture["home"] + "/.config",
                    "GIT_OPTIONAL_LOCKS": "0"})
        if guarded:
            env["PATH"] = str(guard) + os.pathsep + env.get("PATH", "")
        command = [str(executable), *argv]
        directory = cwd or fixture["root"]
        start = time.monotonic()
        response = subprocess.run(command, cwd=directory, env=env, capture_output=True,
                                  timeout=20, check=False)
        record = {"argv": command, "cwd": directory, "exit": response.returncode,
                  "duration_seconds": time.monotonic() - start, "guarded": guarded,
                  "stdout_sha256": sha256(response.stdout), "stderr_sha256": sha256(response.stderr)}
        calls.append(record)
        if args.evidence:
            for stream in ("stdout", "stderr"):
                path = args.evidence / f"{len(calls):04d}.{stream}"
                path.write_bytes(getattr(response, stream))
                record[stream] = str(path)
        return response

    def document(argv, f, kind, expected, cwd=None):
        response = run(argv, f, cwd=cwd, guarded=True)
        require(response.returncode == 0 and not response.stderr,
                f"{argv}: exit={response.returncode}, stderr={response.stderr!r}")
        require(response.stdout.endswith(b"\n"), f"{argv}: missing newline")
        value = json.loads(response.stdout, object_pairs_hook=unique_object)
        subset(value, {"kind": kind, "schema_version": 1, **expected})
        if args.responses and kind in NEW_KINDS:
            (args.responses / f"{len(calls):04d}-{kind}.json").write_bytes(response.stdout)
        return value

    def filtered(f, filters):
        return [t for t in f["tasks"] if all(t[TASK_FIELDS[key]] == value
                                            for key, value in filters.items())]

    def filter_argv(filters):
        return [word for key, value in filters.items() for word in ("--" + key, value)]

    def task_list(f, filters, cwd=None):
        return document(["workspace", "task", "list", "--format", "json", *filter_argv(filters)],
                        f, "WorkspaceTaskListReadback@1", {"workspace": {"root": f["root"]},
                        "scope": {key + "_id": filters.get(key) for key in TASK_FIELDS},
                        "tasks": filtered(f, filters)}, cwd)

    def clean_error(argv, f, code, cwd=None):
        response = run(argv, f, cwd=cwd, guarded=True)
        require(response.returncode != 0 and not response.stdout and
                response.stderr.startswith(f"Error: {code}:".encode()),
                f"{argv}: expected clean {code} error, got {response!r}")

    def compatible(argv, f, cwd=None, *, observed_time=False):
        old, new = run(argv, f, baseline=True, cwd=cwd), run(argv, f, cwd=cwd)
        streams = [old.stdout, new.stdout]
        if observed_time and old.returncode == new.returncode == 0:
            for i, raw in enumerate(streams):
                value = json.loads(raw, object_pairs_hook=unique_object)
                if value.get("kind") == "WorkspaceTaskReadback@3":
                    stamp = value["observed_at_utc"]
                    require(isinstance(stamp, str) and stamp.endswith("Z"), "invalid Task-show observation time")
                    parsed = datetime.fromisoformat(stamp.replace("Z", "+00:00"))
                    require(parsed.tzinfo == timezone.utc, "Task-show observation time is not UTC")
                    value["observed_at_utc"] = "<per-call observation time>"
                streams[i] = json_bytes(value)
        require((old.returncode, streams[0], old.stderr) ==
                (new.returncode, streams[1], new.stderr), f"old read behavior changed: {argv}, cwd={cwd}")
        return new

    try:
        primary, healthy, unavailable = fixtures
        document(["capabilities", "--format", "json"], primary, "PlyCapabilities@1",
                 {"coverage": "workspace-core-read"}, primary["outside_root"])
        document(["workspace", "project", "list", "--format", "json"], primary,
                 "WorkspaceProjectListReadback@1", {"workspace": {"root": primary["root"]},
                 "scope": {"project_id": None}, "projects": primary["projects"]})
        for p in primary["projects"]:
            document(["workspace", "project", "show", "--format", "json", "--", p["project_id"]],
                     primary, "WorkspaceProjectReadback@1", {"workspace": {"root": primary["root"]},
                     "scope": {"project_id": p["project_id"]}, "project": p,
                     "repositories": primary["repositories"][p["project_id"]]})

        for f in fixtures:
            domains = {"project": [p["project_id"] for p in f["projects"]],
                       "repo": sorted({r["repo_id"] for rs in f["repositories"].values() for r in rs}),
                       "epic": sorted({t["parent_epic_id"] for t in f["tasks"]})}
            task_list(f, {})
            for size in (1, 2, 3):
                for keys in itertools.combinations(TASK_FIELDS, size):
                    for values in itertools.product(*(domains[key] for key in keys)):
                        task_list(f, dict(zip(keys, values)))
            # Ordinary reads must not inherit the current Epic or repository.
            for p in f["projects"]:
                task = next(t for t in f["tasks"] if t["project_id"] == p["project_id"])
                task_list(f, {}, str(Path(p["wrapper"]) / task["parent_epic_id"]))

        # New text scopes retain the exact registered rows and all supplied labels.
        selection = {"project": "alpha", "repo": "alpha-main", "epic": "epic-alpha"}
        for size in (1, 2, 3):
            for keys in itertools.combinations(selection, size):
                filters = {key: selection[key] for key in keys}
                argv = ["workspace", "task", "list", *filter_argv(filters)]
                response = run(argv, primary, guarded=True)
                explicit = run(argv + ["--format", "text"], primary, guarded=True)
                require(response.returncode == explicit.returncode == 0 and
                        not response.stderr and not explicit.stderr and response.stdout == explicit.stdout,
                        "default and explicit text differ")
                lines = response.stdout.decode().splitlines()
                require(all(value in lines[0] for value in filters.values()), "text scope lost a filter")
                expected = filtered(primary, filters)
                require([line.strip().split(":", 1)[0] for line in lines[1:]] ==
                        [t["task_id"] for t in expected], "text filter returned incorrect Task rows")
                for line, task in zip(lines[1:], expected):
                    fields = ("project_id", "repo_id", "worktree_state")
                    # The preserved Epic-only renderer names the Epic in its header.
                    if "project" in filters or "repo" in filters:
                        fields += ("parent_epic_id",)
                    require(all(task[key] in line for key in fields), "text row lost raw binding/state")
                    require((task["title"] or "Problem content unavailable") in line, "text row lost title")
        plain = run(["workspace", "task", "list"], primary, guarded=True)
        explicit = run(["workspace", "task", "list", "--format", "text"], primary, guarded=True)
        logged = run(["--json", "workspace", "task", "list"], primary, guarded=True)
        require(plain.returncode == explicit.returncode == logged.returncode == 0 and
                plain.stdout == explicit.stdout == logged.stdout, "text/default/global logging semantics differ")

        for key in TASK_FIELDS:
            for bad in ("", "invalid/id", "UPPERCASE"):
                clean_error(["workspace", "task", "list", "--" + key, bad, "--format", "json"], primary,
                            "workspace_work_invalid_arguments")
            # Known disjoint filters must not hide a missing third identifier.
            filters = {"project": "alpha", "repo": "zeta-main", "epic": "epic-zeta"}
            filters[key] = "rl01-missing-id"
            clean_error(["workspace", "task", "list", *filter_argv(filters), "--format", "json"], primary,
                        "workspace_work_not_found")
        for argv, code in ((["capabilities"], "capabilities_invalid_arguments"),
                           (["workspace", "project", "list"], "workspace_project_invalid_arguments"),
                           (["workspace", "project", "show", "alpha"], "workspace_project_invalid_arguments"),
                           (["workspace", "task", "list"], "workspace_work_invalid_arguments")):
            clean_error(argv + ["--format", "xml"], primary, code, primary["outside_root"])
            clean_error(argv + ["extra", "--format", "json"], primary, code)
        clean_error(["workspace", "project", "show", "--format", "json", "rl01-missing-id"], primary,
                    "workspace_project_not_found")

        target = ["--project", "alpha", "--repo", "alpha-main", "--epic", "epic-alpha"]
        queue_response = compatible(["workspace", "task", "queue", "list", *target, "--format", "json"], healthy)
        require(queue_response.returncode == 0, "healthy queue read failed")
        queue = json.loads(queue_response.stdout, object_pairs_hook=unique_object)
        require([entry["task_id"] for entry in queue["pending"]] == ["task-b-legacy", "task-a-current"],
                "canonical pending order differs from seed")
        registered = task_list(healthy, {"project": "alpha", "repo": "alpha-main", "epic": "epic-alpha"})
        require(len(registered["tasks"]) == 4, "ordinary scope lost Tasks outside the pending queue")
        require(registered["tasks"][0]["task_id"] == "task-a-current", "ordinary list followed queue priority")
        for fmt in ("text", "json"):
            compatible(["workspace", "task", "list", "--ready", *target, "--format", fmt], healthy)
            epic_cwd = str(Path(healthy["projects"][0]["wrapper"]) / "epic-alpha")
            compatible(["workspace", "task", "list", "--ready", "--format", fmt], healthy, epic_cwd)
            for size in (1, 2):
                for keys in itertools.combinations(selection, size):
                    filters = {key: selection[key] for key in keys}
                    response = compatible(["workspace", "task", "list", "--ready", *filter_argv(filters),
                                           "--format", fmt], healthy)
                    require(response.returncode != 0, "ready accepted a partial target")
        for f in (primary, healthy):
            for task in f["tasks"]:
                for fmt in ("text", "json"):
                    compatible(["workspace", "task", "show", "--format", fmt, "--", task["task_id"]], f,
                               observed_time=fmt == "json")
        # Corrupt v4 content remains strict in pre-existing ready/show modes.
        compatible(["workspace", "task", "list", "--ready", *target, "--format", "json"], unavailable)
        compatible(["workspace", "task", "show", "--format", "json", "task-f-corrupt"], unavailable)
        result["result"] = "technical_pass"
        result["compatibility_normalization"] = "Task-show @3 JSON compares all fields except its validated per-call observed_at_utc; raw outputs remain preserved. Other old reads compare bytes."
    except Exception as exc:
        result["error"] = str(exc)
        raise
    finally:
        after = snapshot(watch)
        delta = sorted(key for key in before.keys() | after.keys() if before.get(key) != after.get(key))
        result.update({"calls": len(calls), "commands": calls, "state_delta": delta,
                       "external_tool_attempts": attempts.read_text().splitlines() if attempts.exists() else [],
                       "before_sha256": sha256(json_bytes(before)), "after_sha256": sha256(json_bytes(after))})
        if delta or result["external_tool_attempts"]:
            result["result"] = "failed"
        if args.evidence:
            (args.evidence / "before.json").write_bytes(json_bytes(before))
            (args.evidence / "after.json").write_bytes(json_bytes(after))
            (args.evidence / "result.json").write_bytes(json_bytes(result))
        require(not delta, f"read phase changed fixture or binary state: {delta[:12]}")
        require(not result["external_tool_attempts"], "new read command invoked an external tool")
    print(json.dumps({key: value for key, value in result.items() if key != "commands"}))


if __name__ == "__main__":
    main()
