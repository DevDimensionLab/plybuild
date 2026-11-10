#!/usr/bin/env python3
"""Exercise a built status CLI against preserved multi-Project local fixtures.

This is automated technical acceptance, not human product QA. The supplied
--root must be new; all fixtures, responses, commands and costs remain there.
"""

import argparse
import datetime as dt
import json
import os
from pathlib import Path
import shlex
import statistics
import sys

sys.dont_write_bytecode = True

from workspace_views_acceptance import Journey, contains, sha, snapshot
from workflow_status_schema_test import validate_responses


class StatusJourney(Journey):
    def status(self, *flags, cwd=None, expected=0):
        before = len(self.commands)
        text, elapsed = self.command([self.ply, "workflow", "status", *flags], cwd=cwd, expected=expected)
        result = self.commands[before]
        if expected:
            assert text == "" and result["stderr"], result
            return None, elapsed
        assert text.endswith("\n") and text.count("\n") == 1, result
        value = json.loads(text)
        assert value["kind"] == "WorkflowStatusReadback@1", value
        return value, elapsed

    def setup_question(self):
        source = self.root / "inputs" / "question.md"
        source.write_text("Synthetic registered agent question. No actual human QA.\n")
        event = dict(
            kind="ply.workspace.task-journal-event-input", schema_version=1,
            publication_key="workflow-status:question", task_id="t0-0",
            actor=dict(id="fixture", role="tool", session_id=None),
            activity_id="workflow-status-fixture", run_binding=None,
            occurred_at=dt.datetime.now(dt.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            time_basis=dict(kind="observed", clock="fixture UTC", precision="second", uncertainty=None),
            sources=[dict(locator=str(source), sha256=sha(source))],
            type="step_started", title="Synthetic agent question", detail="",
            step=dict(id="waiting", kind="waiting", parent_step_id=None),
            outcome=None, candidate=None, relations=[],
            waiting=dict(reason="Which output should be the default?", dependency="fixture-answer", next_actor="human"),
        )
        draft = self.root / "inputs" / "question.json"
        draft.write_text(json.dumps(event) + "\n")
        self.ply_json("workspace", "task", "journal", "record", "t0-0", "--file", draft)

    def guard_live_commands(self):
        trap = self.root / "forbidden-live-command.log"
        directory = self.root / "blocked-executables"
        directory.mkdir()
        for name in ("git", "herdr", "codex", "claude"):
            script = directory / name
            script.write_text("#!/bin/sh\nprintf '%s\\n' " + shlex.quote(name) + " >> " + shlex.quote(str(trap)) + "\nexit 93\n")
            script.chmod(0o700)
        self.env["PATH"] = str(directory) + os.pathsep + self.env.get("PATH", "")
        return trap

    def check_status(self):
        empty, _ = self.status("--json")
        assert len(empty["projects"]) == 5 and len(empty["epics"]) == 6, empty
        assert empty["items"] == [] and empty["counts"]["hidden"]["backlog"] == 50, empty
        self.setup_question()
        active, _ = self.status("--json")
        assert [row["task_id"] for row in active["items"]] == ["t0-0"], active
        question = active["items"][0]
        assert question["category"] == "needs_you", question
        assert question["next_actions"][0]["actor"] == "human", question
        assert question["next_actions"][0]["reason"] == "Which output should be the default?", question
        assert question["progress"]["technical_gate"] is None and question["progress"]["human_qa_outcome"] is None
        self.ply_json("workspace", "epic", "lifecycle", "set", "e0", "parked", "--actor", "technical fixture")
        parked, _ = self.status("--json")
        assert parked["items"] == [] and parked["counts"]["hidden"] == dict(inactive=10, completed=0, backlog=40), parked
        history, _ = self.status("--all", "--json")
        assert len(history["items"]) == 50 and not any(history["counts"]["hidden"].values()), history
        old = next(row for row in history["items"] if row["task_id"] == "t0-0")
        assert old["category"] == "inactive" and all(not a["current"] for a in old["next_actions"]), old
        self.ply_json("workspace", "epic", "lifecycle", "set", "e0", "active", "--actor", "technical fixture")

        nested = self.workspace / "arbitrary" / "nested"
        nested.mkdir(parents=True)
        outside = self.root / "outside"
        outside.mkdir()
        trap = self.guard_live_commands()
        before = snapshot(self.workspace)
        for argv in (["--json", "workflow", "status"], ["workflow", "--json", "status"], ["workflow", "status", "--json"], ["--debug", "workflow", "status", "--format", "json"]):
            text, _ = self.command([self.ply, *argv], cwd=nested)
            assert text.count("\n") == 1 and json.loads(text)["workspace"]["root"] == str(self.workspace), text
        full, _ = self.status("--all", "--format", "json", cwd=nested)
        assert full["counts"]["total"] == 50 and len(full["items"]) == 50
        for flags in (["--project", "p1"], ["--repo", "r1"], ["--epic", "e1"], ["--project", "p1", "--repo", "r1", "--epic", "e1"]):
            selected, _ = self.status("--all", "--json", *flags)
            assert len(selected["items"]) == 10 and all(row["project_id"] == "p1" for row in selected["items"]), selected
        disjoint, _ = self.status("--all", "--json", "--project", "p0", "--repo", "r1")
        assert disjoint["projects"] == disjoint["epics"] == disjoint["items"] == [], disjoint
        empty_epic, _ = self.status("--json", "--epic", "empty")
        assert len(empty_epic["epics"]) == 1 and not empty_epic["items"], empty_epic
        for flags in (["--project", "missing"], ["--project", "p0", "--epic", "missing"], ["--json", "--format", "text"], ["--format", "xml"], ["--unexpected"]):
            self.status("--json", *flags, expected=1)
        self.status("--json", cwd=outside, expected=1)
        assert list(outside.iterdir()) == [], "status initialized an absent workspace"
        text, _ = self.command([self.ply, "workflow", "status"], cwd=nested)
        for part in ("Needs you (1)", "Which output should be the default?", "Registered facts", "Hidden: 49 backlog"):
            assert part in text, text
        for argv in (["--help"], ["workflow", "--help"], ["workflow", "status", "--help"]):
            help_text, _ = self.command([self.ply, *argv])
            assert "status" in help_text
        times = [self.status("--json", cwd=nested)[1] for _ in range(5)]
        self.measurements["workflow_status_50_tasks_5_projects"] = dict(samples_seconds=times, median_seconds=statistics.median(times), max_seconds=max(times))
        assert not trap.exists(), "status executed a live Git/provider/Herdr command"
        assert snapshot(self.workspace) == before, "status changed registered sources, history, Git metadata or timestamps"

        if self.baseline:
            compatible = [["workspace", "overview"], ["workspace", "status"], ["workspace", "task", "list"], ["workspace", "task", "list", "--progress"], ["workspace", "task", "show", "t0-0", "--progress"], ["workspace", "run", "list"]]
            for argv in compatible:
                old, _ = self.command([self.baseline, *argv, "--format", "json"], capture=False)
                current, _ = self.ply_json(*argv)
                previous = json.loads(old)
                for stamp in ("as_of", "observed_at_utc"):
                    if stamp in previous: previous[stamp] = current[stamp]
                contains(current, previous)
            old, _ = self.command([self.baseline, "capabilities", "--format", "json"], capture=False)
            current, _ = self.ply_json("capabilities")
            previous = json.loads(old)
            assert current["operations"] == previous["operations"] and current["coverage"] == previous["coverage"]
            assert any(op["command"] == ["workflow", "status"] for op in current["read_extensions"])
            for extension in previous["read_extensions"]: assert extension in current["read_extensions"]
        assert not trap.exists() and snapshot(self.workspace) == before
        validate_responses(self.capture)

    def report(self, outcome, error=None):
        value = dict(kind="PlyWorkflowStatusAcceptance@1", outcome=outcome, human_qa="not_performed", ply=str(self.ply), ply_sha256=sha(self.ply), baseline=str(self.baseline) if self.baseline else None, workspace=str(self.workspace), task_count=50, project_count=5, measurements=self.measurements, error=error, commands=self.commands)
        (self.root / "report.json").write_text(json.dumps(value, indent=2) + "\n")
        print(json.dumps({key: value[key] for key in ("kind", "outcome", "workspace", "measurements", "error")}, indent=2))


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--ply", type=Path, required=True)
    parser.add_argument("--baseline", type=Path)
    parser.add_argument("--root", type=Path, required=True)
    args = parser.parse_args()
    journey = StatusJourney(args)
    try:
        journey.setup()
        journey.check_status()
    except Exception as error:
        journey.report("fail", str(error))
        raise
    else:
        journey.report("pass")
