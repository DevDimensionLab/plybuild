"""Isolated two-Task roundtrip for Goal Specs and queue prioritization."""

import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys

root = Path(sys.argv[1])
workspace = root / "workspace"
repo = workspace / "ply/main"
epic = workspace / "ply/epic"
repo.mkdir(parents=True)
home = root / "home"
home.mkdir()
env = dict(os.environ, HOME=str(home), GIT_CONFIG_NOSYSTEM="1")


def command(*argv):
    result = subprocess.run(argv, cwd=workspace, env=env, capture_output=True, text=True)
    assert result.returncode == 0, (argv, result.returncode, result.stdout, result.stderr)
    return result.stdout


def ply(*args):
    return command(root / "ply", *args)


def read_json(*args):
    return json.loads(ply(*args, "--format", "json"))


def git(*args):
    return command("git", "-C", repo, "-c", "user.name=Ply tests",
                   "-c", "user.email=tests@example.invalid", *args).strip()


def state_bytes():
    return {
        path.relative_to(workspace): hashlib.sha256(path.read_bytes()).hexdigest()
        for directory in (workspace / ".ply", workspace / "ply/data/work")
        for path in directory.rglob("*") if path.is_file()
    }


def git_state():
    return git("show-ref"), git("worktree", "list", "--porcelain")


def goal_entry(task_id):
    spec_id = f"{task_id}-goal"
    objective = f"Deliver {task_id}."
    args = ("workspace", "task", "goal", "register", task_id,
            "--spec", spec_id, "--title", "Same title", "--objective", objective,
            "--criterion", "The exact Goal is queued.", *target)
    preview = read_json(*args)
    applied = read_json(*args, "--apply", "--confirmation", preview["confirmation"])
    assert applied["applied"], applied
    spec = read_json("workspace", "task", "spec", "show", task_id, "--spec", spec_id,
                     "--revision", str(applied["goal"]["spec"]["revision"]))
    assert spec["kind"] == "WorkspaceTaskSpecReadback@2", spec
    assert spec["goal"]["objective"] == objective and spec["integrity"] == "valid", spec
    assert spec["goal"]["goal"] == applied["goal"], spec
    return {"task_id": task_id, "goal": applied["goal"]}


ply("workspace", "init")
git("init", "-b", "main")
git("commit", "--allow-empty", "-m", "base")
base = git("rev-parse", "HEAD")
git("worktree", "add", "-b", "epic", epic, base)
ply("workspace", "project", "add", "ply", "--name", "Ply", "--wrapper", workspace / "ply",
    "--repo", f"ply={repo}")
ply("workspace", "epic", "adopt", "epic", "--title", "Epic", "--project", "ply",
    "--repo", "ply", "--worktree", epic, "--ref", "refs/heads/epic", "--expected-oid", base)

target = ("--project", "ply", "--repo", "ply", "--epic", "epic")
queue_list = ("workspace", "task", "queue", "list", *target)
before = state_bytes()
git_before = git_state()
empty = read_json(*queue_list)
assert empty["revision"] == 0 and empty["pending"] == [], empty
assert state_bytes() == before, "empty queue read wrote state"

entries = [goal_entry("task"), goal_entry("task-b")]
queued = read_json(*queue_list)
assert [row["task_id"] for row in queued["pending"]] == ["task", "task-b"], queued
entries.reverse()
draft = {
    "kind": "WorkspaceTaskQueueDraft@2",
    "schema_version": 2,
    "publication_key": "queue/reorder",
    "project_id": "ply",
    "repo_id": "ply",
    "epic_id": "epic",
    "expected_revision": queued["revision"],
    "entries": entries,
    "registry_upgrade": None,
    "recorder": {
        "actor_claim": "synthetic fixture",
        "control_surface": "local CLI roundtrip",
        "recorded_at_utc": "2026-09-30T00:00:00Z",
    },
}
draft_path = workspace / "queue.json"
draft_path.write_text(json.dumps(draft))
queue_set = ("workspace", "task", "queue", "set", "--file", draft_path)
prioritized = read_json(*queue_set)
assert prioritized["revision"] == queued["revision"] + 1, prioritized
assert [row["state"] for row in prioritized["pending"]] == ["ready", "ready"], prioritized
assert [{key: row[key] for key in ("task_id", "goal")} for row in prioritized["pending"]] == entries, prioritized

after = state_bytes()
assert read_json(*queue_list) == prioritized, "queue readback lost the new order"
assert read_json(*queue_set) == prioritized, "replaying the same draft changed the readback"
assert state_bytes() == after, "list or replay changed stored state"
assert git_state() == git_before, "queue or Goal registration changed Git"
assert not (home / ".ply").exists(), "roundtrip wrote to HOME"
