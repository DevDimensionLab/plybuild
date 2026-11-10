#!/usr/bin/env python3
"""Installed CLI acceptance for workspace views, using a preserved local fixture.

Creates five local repositories, five populated Epics, one empty Epic and fifty
Tasks. It never uses a provider, network, real user QA, or external workspace write.
All input, output, measurements and fixture state remain available after the run.
"""

import argparse
import datetime as dt
import hashlib
import json
import os
from pathlib import Path
import statistics
import subprocess
import time


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def snapshot(root):
    """Include Git metadata, registry files, directory entries and timestamps."""
    out = {}
    for path in sorted(root.rglob("*")):
        rel = str(path.relative_to(root))
        st = path.lstat()
        if path.is_symlink():
            out[rel] = ["symlink", os.readlink(path), st.st_mtime_ns]
        elif path.is_dir():
            out[rel] = ["directory", st.st_mode, st.st_mtime_ns]
        else:
            out[rel] = ["file", st.st_mode, st.st_mtime_ns, sha(path)]
    return out


def contains(actual, legacy):
    """An additive field must not change a legacy field or ordered collection."""
    if isinstance(legacy, dict):
        assert isinstance(actual, dict)
        for key, value in legacy.items():
            assert key in actual, key
            contains(actual[key], value)
    elif isinstance(legacy, list):
        assert isinstance(actual, list) and len(actual) == len(legacy)
        for new, old in zip(actual, legacy):
            contains(new, old)
    else:
        assert actual == legacy, (actual, legacy)


class Journey:
    def __init__(self, args):
        self.ply = args.ply.resolve(strict=True)
        self.baseline = args.baseline.resolve(strict=True) if args.baseline else None
        self.root = args.root.absolute()
        assert not self.root.exists(), "Preserve existing fixtures; choose a new --root."
        self.root.mkdir(mode=0o700, parents=True)
        self.root = self.root.resolve()
        self.workspace = self.root / "workspace"
        self.workspace.mkdir()
        self.capture = self.root / "responses"
        self.capture.mkdir()
        (self.root / "inputs").mkdir()
        self.commands = []
        self.measurements = {}
        # Caller Git routing/config must never redirect fixture writes outside root.
        self.env = {key: value for key, value in os.environ.items()
                    if not key.startswith("GIT_")}
        self.env.update(GIT_OPTIONAL_LOCKS="0", GIT_TERMINAL_PROMPT="0",
                        GIT_CONFIG_NOSYSTEM="1", GIT_CONFIG_GLOBAL=os.devnull)

    def command(self, argv, cwd=None, expected=0, capture=True):
        start = time.perf_counter()
        process = subprocess.run(
            [str(v) for v in argv], cwd=cwd or self.workspace, env=self.env,
            capture_output=True, text=True, timeout=60,
        )
        elapsed = time.perf_counter() - start
        record = dict(argv=[str(v) for v in argv], cwd=str(cwd or self.workspace),
                      exit_code=process.returncode, elapsed_seconds=elapsed,
                      stdout=process.stdout, stderr=process.stderr)
        self.commands.append(record)
        index = len(self.commands)
        if capture and process.stdout.startswith("{"):
            (self.capture / f"{index:03}.json").write_text(process.stdout)
        assert (process.returncode == 0) if expected == 0 else (process.returncode != 0), record
        if expected == 0:
            assert process.stderr == "", record
        return process.stdout, elapsed

    def ply_json(self, *args, expected=0):
        out, elapsed = self.command([self.ply, *args, "--format", "json"], expected=expected)
        return (json.loads(out) if expected == 0 else out), elapsed

    def git(self, cwd, *args):
        # Explicit local identities and disabled optional hooks/filters keep setup
        # inside this fixture even when the caller has personal Git defaults.
        command = ["git", "-c", "core.hooksPath=/dev/null", "-c", "core.fsmonitor=false",
                   "-c", "commit.gpgsign=false", "-c", "user.name=Ply fixture",
                   "-c", "user.email=fixture@invalid", *args]
        process = subprocess.run(command, cwd=cwd, env=self.env, text=True,
                                 capture_output=True, timeout=30)
        assert process.returncode == 0, process.stderr
        return process.stdout.strip()

    def setup(self):
        self.command([self.ply, "workspace", "init"])
        for n in range(5):
            wrapper = self.workspace / f"project-{n}"
            main = wrapper / "main"
            main.mkdir(parents=True)
            self.git(main, "init", "--initial-branch=main")
            self.git(main, "commit", "--allow-empty", "-m", "Fixture base")
            oid = self.git(main, "rev-parse", "HEAD")
            self.command([self.ply, "workspace", "project", "add", f"p{n}",
                          "--name", f"Project {n}", "--wrapper", wrapper,
                          "--repo", f"r{n}={main}"])
            epic = wrapper / "epic"
            self.git(main, "worktree", "add", "-b", "epic", str(epic))
            self.command([self.ply, "workspace", "epic", "adopt", f"e{n}",
                          "--title", f"Epic {n}", "--project", f"p{n}", "--repo", f"r{n}",
                          "--worktree", epic, "--ref", "refs/heads/epic", "--expected-oid", oid])
            for task in range(10):
                self.ply_json("workspace", "task", "create", f"t{n}-{task}",
                              "--title", f"Task {n}-{task}", "--description", "Synthetic view acceptance",
                              "--project", f"p{n}", "--repo", f"r{n}", "--epic", f"e{n}")
            if n == 0:
                empty = wrapper / "empty"
                self.git(main, "worktree", "add", "-b", "empty", str(empty))
                self.command([self.ply, "workspace", "epic", "adopt", "empty",
                              "--title", "Empty Epic", "--project", "p0", "--repo", "r0",
                              "--worktree", empty, "--ref", "refs/heads/empty", "--expected-oid", oid])

    def check(self):
        tasks, _ = self.ply_json("workspace", "task", "list", "--progress")
        assert len(tasks["tasks"]) == 50
        assert all(row["progress"]["state"] == "nothing" for row in tasks["tasks"])
        attention, _ = self.ply_json("workspace", "attention")
        assert attention["items"] == [], "Unstarted work must not create Needs You."
        epics, _ = self.ply_json("workspace", "epic", "list")
        assert len(epics["epics"]) == 6
        assert sum(row["task_count"] for row in epics["epics"]) == 50
        assert next(row for row in epics["epics"] if row["epic_id"] == "empty")["task_count"] == 0

        source = self.root / "inputs" / "synthetic-wait.md"
        source.write_text("Synthetic process waiting observation; no human or technical QA claim.\n")
        now = dt.datetime.now(dt.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        event = dict(kind="ply.workspace.task-journal-event-input", schema_version=1,
                     publication_key="workspace-views:synthetic-wait", task_id="t0-0",
                     actor=dict(id="acceptance-fixture", role="tool", session_id=None),
                     activity_id="synthetic-view-acceptance", run_binding=None, occurred_at=now,
                     time_basis=dict(kind="observed", clock="fixture UTC", precision="second", uncertainty=None),
                     sources=[dict(locator=str(source), sha256=sha(source))], type="step_started",
                     title="Synthetic human design question", detail="",
                     step=dict(id="synthetic-wait", kind="waiting", parent_step_id=None),
                     outcome=None, candidate=None, relations=[],
                     waiting=dict(reason="Synthetic question for acceptance", dependency="fixture-answer", next_actor="human"))
        draft = self.root / "inputs" / "wait.json"
        draft.write_text(json.dumps(event) + "\n")
        self.ply_json("workspace", "task", "journal", "record", "t0-0", "--file", str(draft))
        attention, _ = self.ply_json("workspace", "attention")
        assert [row["subject_id"] for row in attention["items"]] == ["t0-0"]
        self.ply_json("workspace", "task", "lifecycle", "set", "t0-0", "frozen", "--actor", "acceptance fixture")
        assert self.ply_json("workspace", "attention")[0]["items"] == []
        self.ply_json("workspace", "task", "lifecycle", "set", "t0-0", "active", "--actor", "acceptance fixture")
        self.ply_json("workspace", "epic", "lifecycle", "set", "e0", "archived", "--actor", "acceptance fixture")
        assert self.ply_json("workspace", "attention")[0]["items"] == []
        self.ply_json("workspace", "epic", "lifecycle", "set", "e0", "active", "--actor", "acceptance fixture")
        life, _ = self.ply_json("workspace", "task", "lifecycle", "show", "t0-0")
        before = snapshot(self.workspace)
        noop, _ = self.ply_json("workspace", "task", "lifecycle", "set", "t0-0", "active", "--actor", "acceptance fixture", "--expected-revision", str(life["revision"]))
        assert noop["changed"] is False
        assert snapshot(self.workspace) == before, "Lifecycle no-op wrote state."
        self.ply_json("workspace", "task", "lifecycle", "set", "t0-0", "deleted", "--actor", "acceptance fixture", expected=1)
        self.ply_json("workspace", "task", "lifecycle", "set", "t0-0", "frozen", "--actor", "acceptance fixture", "--expected-revision", "0", expected=1)
        assert snapshot(self.workspace) == before, "Rejected lifecycle changed state."

        planning = self.workspace / "project-0" / "planning"
        planning.mkdir()
        self.ply_json("workspace", "project", "companion", "add", "p0", "--role", "planning", "--path", str(planning), "--actor", "acceptance fixture")
        self.ply_json("workspace", "project", "repo-wrapper", "set", "p0", "--repo", "r0", "--wrapper", str(planning.parent), "--actor", "acceptance fixture")
        project, _ = self.ply_json("workspace", "project", "show", "p0")
        assert project["project"]["companions"] == [dict(role="planning", locator=str(planning))]
        assert project["repositories"][0]["wrapper"] == str(planning.parent)

        readonly = [
            ["capabilities"], ["workspace", "project", "list"], ["workspace", "project", "show", "p0"],
            ["workspace", "task", "list"], ["workspace", "task", "list", "--progress"],
            ["workspace", "task", "show", "t0-0", "--progress"], ["workspace", "epic", "list"],
            ["workspace", "attention"], ["workspace", "overview"], ["workspace", "journal", "recent"],
            ["workspace", "journal", "recent", "--since", now], ["workspace", "run", "list"],
            ["workspace", "run", "list", "--active"], ["workspace", "worktree", "list"],
            ["workspace", "status"], ["workspace", "project", "metadata", "show", "p0"],
            ["workspace", "task", "lifecycle", "show", "t0-0"], ["workspace", "epic", "lifecycle", "show", "e0"],
        ]
        before = snapshot(self.workspace)
        for argv in readonly:
            response, _ = self.ply_json(*argv)
            assert isinstance(response, dict)
        assert snapshot(self.workspace) == before, "Read-only calls changed files or Git metadata."

        bulk, _ = self.ply_json("workspace", "task", "list", "--progress")
        single, _ = self.ply_json("workspace", "task", "show", "t0-0", "--progress")
        overview, _ = self.ply_json("workspace", "overview")
        assert single["task"] == next(row for row in bulk["tasks"] if row["task_id"] == "t0-0")
        assert overview["tasks"] == bulk["tasks"]
        assert all(row["progress"]["human_qa_outcome"] is None for row in bulk["tasks"])

        if self.baseline:
            for argv in readonly[1:4]:
                old, _ = self.command([self.baseline, *argv, "--format", "json"], capture=False)
                current, _ = self.ply_json(*argv)
                contains(current, json.loads(old))
                old_text, _ = self.command([self.baseline, *argv], capture=False)
                new_text, _ = self.command([self.ply, *argv], capture=False)
                assert old_text == new_text, (argv, old_text, new_text)

        for label, argv in {
            "overview_50_tasks_5_projects": ["workspace", "overview"],
            "status_50_tasks_5_projects": ["workspace", "status"],
            "inventory_5_repositories": ["workspace", "worktree", "list"],
        }.items():
            times = [self.ply_json(*argv)[1] for _ in range(5)]
            self.measurements[label] = dict(samples_seconds=times, median_seconds=statistics.median(times), max_seconds=max(times))
        assert self.measurements["overview_50_tasks_5_projects"]["max_seconds"] < 1, "Home overview exceeded the one-second target."

    def report(self, outcome, error=None):
        value = dict(kind="PlyWorkspaceViewsAcceptance@1", outcome=outcome,
                     human_qa="not_performed", ply=str(self.ply), ply_sha256=sha(self.ply),
                     baseline=str(self.baseline) if self.baseline else None,
                     workspace=str(self.workspace), task_count=50, project_count=5,
                     measurements=self.measurements, error=error, commands=self.commands)
        (self.root / "report.json").write_text(json.dumps(value, indent=2) + "\n")
        print(json.dumps({key: value[key] for key in ("outcome", "workspace", "measurements", "error")}, indent=2))


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--ply", type=Path, required=True)
    parser.add_argument("--baseline", type=Path)
    parser.add_argument("--root", type=Path, required=True)
    args = parser.parse_args()
    journey = Journey(args)
    try:
        journey.setup()
        journey.check()
    except Exception as exc:
        journey.report("fail", str(exc))
        raise
    else:
        journey.report("pass")
