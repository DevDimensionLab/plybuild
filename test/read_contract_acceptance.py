#!/usr/bin/env python3
"""RL-01 read-phase acceptance. Does not create fixtures or build/install Ply.

Run only after implementation is selected, with two built binaries and disposable
fixtures prepared by the delivery's canonical Go fixture helpers. Seed all expected
values independently of candidate output. JSON fixture manifest:

  kind: "PlyReadContractFixture@1"
  root, empty_root, outside_root, corrupt_projects_root, corrupt_tasks_root, home:
    distinct physical absolute directories (home is an empty private test home)
  projects: sorted expected Project objects from Spec section 3
  repositories: object mapping each project_id to sorted expected Repository objects
  tasks: sorted expected Task objects from Spec section 3
  protected_paths: extra fixture Git/worktree paths outside the directories above

Populated fixture: at least two projects, repositories and Epics, Unicode and
punctuation in names/paths/titles, a missing checkout, Tasks outside the queue,
all four known worktree states, current/legacy/unavailable titles, and a legitimate
available title "Problem content unavailable". Corrupt fixtures contain a valid
workspace marker and an invalid respective registry. The empty root is initialized
but has no registrations. The outside root has no containing workspace.

Subprocesses get the private home; this script never changes the parent environment.
This is technical acceptance, never a human product verdict. Product-side tests own
fixture validation, old ready/show variants, injected faults and JSON Schema checks.
"""

import argparse
import hashlib
import itertools
import json
import os
from pathlib import Path
import stat
import subprocess


def require(condition, message):
    if not condition:
        raise AssertionError(message)


def subset(actual, expected, where="response"):
    """Required fields/types and exact list membership; tolerate additive fields."""
    require(type(actual) is type(expected), f"{where}: wrong type")
    if isinstance(expected, dict):
        for key, value in expected.items():
            require(key in actual, f"{where}: missing {key}")
            subset(actual[key], value, f"{where}.{key}")
    elif isinstance(expected, list):
        require(len(actual) == len(expected), f"{where}: wrong list length")
        for index, (left, right) in enumerate(zip(actual, expected)):
            subset(left, right, f"{where}[{index}]")
    else:
        require(actual == expected, f"{where}: expected {expected!r}, got {actual!r}")


def unique_object(pairs):
    result = {}
    for key, value in pairs:
        require(key not in result, f"duplicate JSON key: {key}")
        result[key] = value
    return result


def snapshot(paths):
    result = {}

    def visit(path):
        try:
            info = path.lstat()
        except FileNotFoundError:
            result[str(path)] = ["absent"]
            return
        mode = stat.S_IMODE(info.st_mode)
        if stat.S_ISLNK(info.st_mode):
            result[str(path)] = ["symlink", mode, os.readlink(path)]
        elif stat.S_ISDIR(info.st_mode):
            result[str(path)] = ["directory", mode]
            for child in sorted(path.iterdir()):
                visit(child)
        elif stat.S_ISREG(info.st_mode):
            result[str(path)] = ["file", mode, hashlib.sha256(path.read_bytes()).hexdigest()]
        else:
            raise AssertionError(f"fixture contains special file: {path}")

    for value in paths:
        visit(Path(value))
    return result


def catalog_entry(op, mode, kind, selectors=(), filters=(), policy="none", schemas=None):
    return {
        "id": op, "command": ["workspace", *op.split(".")], "mode": mode,
        "selectors": list(selectors), "formats": ["text", "json"],
        "filters": list(filters), "filter_policy": policy, "effect": "read",
        "result_schemas": schemas or [{"kind": kind, "schema_version": 1}],
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__,
                                     formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--ply", type=Path, required=True, help="absolute candidate binary")
    parser.add_argument("--baseline", type=Path, required=True, help="binary from the bound old base")
    parser.add_argument("--fixtures", type=Path, required=True, help="prepared fixture manifest")
    args = parser.parse_args()
    fixture = json.loads(args.fixtures.read_text(), object_pairs_hook=unique_object)
    require(fixture["kind"] == "PlyReadContractFixture@1", "unknown fixture manifest")
    root_keys = ("root", "empty_root", "outside_root", "corrupt_projects_root",
                 "corrupt_tasks_root", "home")
    roots = [fixture[key] for key in root_keys]
    require(len(set(roots)) == len(roots), "fixture directories must be distinct")
    for value in roots:
        require(Path(value).is_dir() and str(Path(value).resolve()) == value,
                "fixture paths must be existing physical absolute directories")
    require(not any(Path(fixture["home"]).iterdir()), "private test home must start empty")
    for binary in (args.ply, args.baseline):
        require(binary.is_absolute() and binary.is_file(), "absolute binary required")
    projects, repos, tasks = fixture["projects"], fixture["repositories"], fixture["tasks"]
    require(len(projects) >= 2 and len({t["repo_id"] for t in tasks}) >= 2,
            "fixture needs multiple projects/repositories")
    require(len({t["parent_epic_id"] for t in tasks}) >= 2, "fixture needs multiple Epics")
    require({t["title_source"] for t in tasks} == {"registration", "problem_revision"},
            "fixture needs legacy and current titles")
    require(any(t["title"] is None and t["title_status"] == "unavailable" for t in tasks),
            "fixture needs unavailable problem content")
    require(any(t["title"] == "Problem content unavailable" and
                t["title_status"] == "available" for t in tasks), "fixture needs literal title")
    require({t["worktree_state"] for t in tasks} >= {
        "unbound", "creating", "worktree_ready", "reconciliation_required"},
        "fixture needs all four worktree states")
    watch = roots + fixture["protected_paths"] + [str(args.ply), str(args.baseline)]
    before = snapshot(watch)
    env = {key: os.environ[key] for key in ("PATH", "TMPDIR", "LANG", "LC_ALL", "SYSTEMROOT")
           if key in os.environ}
    env.update({"HOME": fixture["home"], "XDG_CONFIG_HOME": fixture["home"] + "/.config",
                "GIT_OPTIONAL_LOCKS": "0"})
    calls = 0

    def run(argv, cwd=None, binary=None):
        nonlocal calls
        calls += 1
        return subprocess.run([str(binary or args.ply), *argv], cwd=cwd or fixture["root"],
                              env=env, capture_output=True, timeout=20, check=False)

    def document(argv, kind, expected, cwd=None):
        response = run(argv, cwd)
        require(response.returncode == 0 and not response.stderr,
                f"{argv}: exit={response.returncode}, stderr={response.stderr!r}")
        require(response.stdout.endswith(b"\n"), f"{argv}: missing final newline")
        value = json.loads(response.stdout, object_pairs_hook=unique_object)
        subset(value, {"kind": kind, "schema_version": 1, **expected})
        return value

    def error(argv, code=None, cwd=None):
        response = run(argv, cwd)
        require(response.returncode != 0 and not response.stdout and response.stderr,
                f"{argv}: expected clean failure")
        if code:
            require(response.stderr.startswith(f"Error: {code}:".encode()),
                    f"{argv}: wrong error class: {response.stderr!r}")

    def task_list(filters):
        scope = {key + "_id": filters.get(key) for key in ("project", "repo", "epic")}
        expected = [t for t in tasks if all(t[{"project": "project_id", "repo": "repo_id",
                    "epic": "parent_epic_id"}[key]] == value for key, value in filters.items())]
        argv = ["workspace", "task", "list", "--format", "json"]
        for key, value in filters.items():
            argv += ["--" + key, value]
        document(argv, "WorkspaceTaskListReadback@1", {
            "workspace": {"root": fixture["root"]}, "scope": scope, "tasks": expected})

    try:
        entries = [
            catalog_entry("project.list", "default", "WorkspaceProjectListReadback@1"),
            catalog_entry("project.show", "default", "WorkspaceProjectReadback@1"),
            catalog_entry("task.list", "ready", "WorkspaceTaskQueueReadback@1", ["--ready"],
                          ["project", "repo", "epic"], "all_or_none"),
            catalog_entry("task.list", "registered", "WorkspaceTaskListReadback@1", (),
                          ["project", "repo", "epic"], "independent_and"),
            catalog_entry("task.show", "default", None, schemas=[
                {"kind": "WorkspaceTaskReadback@1", "schema_version": None},
                {"kind": "WorkspaceTaskIntegrationReadback@1", "schema_version": 1},
                {"kind": "WorkspaceTaskReadback@3", "schema_version": 3}]),
        ]
        cap = document(["capabilities", "--format", "json"], "PlyCapabilities@1",
                       {"coverage": "workspace-core-read", "operations": entries},
                       fixture["outside_root"])
        require(isinstance(cap["build"]["version"], str) and cap["build"]["version"],
                "build version missing")
        require(cap["build"]["vcs_revision"] is None or
                isinstance(cap["build"]["vcs_revision"], str), "invalid build revision")
        require(cap["build"]["vcs_modified"] is None or
                type(cap["build"]["vcs_modified"]) is bool, "invalid modified metadata")
        repeat = run(["capabilities", "--format", "json"], fixture["outside_root"])
        subset(json.loads(repeat.stdout), cap)
        for where in (fixture["root"], fixture["empty_root"]):
            populated = where == fixture["root"]
            document(["workspace", "project", "list", "--format", "json"],
                     "WorkspaceProjectListReadback@1", {"workspace": {"root": where},
                     "scope": {"project_id": None}, "projects": projects if populated else []}, where)
        for project in projects:
            document(["workspace", "project", "show", "--format", "json", "--", project["project_id"]],
                     "WorkspaceProjectReadback@1", {"workspace": {"root": fixture["root"]},
                     "scope": {"project_id": project["project_id"]}, "project": project,
                     "repositories": repos[project["project_id"]]})
        task_list({})
        domains = {"project": [p["project_id"] for p in projects],
                   "repo": sorted({r["repo_id"] for rs in repos.values() for r in rs}),
                   "epic": sorted({t["parent_epic_id"] for t in tasks})}
        for key, values in domains.items():
            for value in values:
                task_list({key: value})
            require("rl01-missing-id" not in values, "reserved missing ID is present")
            error(["workspace", "task", "list", "--" + key, "rl01-missing-id", "--format", "json"],
                  "workspace_work_not_found")
            error(["workspace", "task", "list", "--" + key, "", "--format", "json"],
                  "workspace_work_invalid_arguments")
        for project, repo in itertools.product(domains["project"], domains["repo"]):
            task_list({"project": project, "repo": repo})
        task_list({"project": tasks[0]["project_id"], "repo": tasks[0]["repo_id"],
                   "epic": tasks[0]["parent_epic_id"]})
        document(["workspace", "task", "list", "--format", "json"], "WorkspaceTaskListReadback@1",
                 {"workspace": {"root": fixture["empty_root"]}, "scope": {
                     "project_id": None, "repo_id": None, "epic_id": None}, "tasks": []},
                 fixture["empty_root"])
        error(["workspace", "project", "list", "--format", "json"],
              "workspace_not_found", fixture["outside_root"])
        error(["workspace", "project", "list", "--format", "json"], cwd=fixture["corrupt_projects_root"])
        error(["workspace", "task", "list", "--format", "json"], cwd=fixture["corrupt_tasks_root"])
        for argv, code in ((["capabilities"], "capabilities_invalid_arguments"),
                           (["workspace", "project", "list"], "workspace_project_invalid_arguments"),
                           (["workspace", "task", "list"], "workspace_work_invalid_arguments")):
            error(argv + ["--format", "xml"], code)
        legacy = [["workspace", "project", "list"], ["workspace", "task", "list"]]
        legacy += [["workspace", "project", "show", "--", p["project_id"]] for p in projects]
        legacy += [["workspace", "task", "list", "--epic", e] for e in domains["epic"]]
        for argv in legacy:
            old, new = run(argv, binary=args.baseline), run(argv)
            require(old.returncode == new.returncode == 0 and
                    old.stdout == new.stdout and old.stderr == new.stderr,
                    f"legacy output changed: {argv}")
    finally:
        after = snapshot(watch)
        delta = sorted(key for key in before.keys() | after.keys() if before.get(key) != after.get(key))
        require(not delta, f"read phase changed fixture/binary state: {delta[:12]}")
    print(json.dumps({"result": "technical_pass", "calls": calls,
                      "candidate_sha256": hashlib.sha256(args.ply.read_bytes()).hexdigest(),
                      "baseline_sha256": hashlib.sha256(args.baseline.read_bytes()).hexdigest(),
                      "human_verdict": None}))


if __name__ == "__main__":
    main()
