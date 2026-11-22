#!/usr/bin/env python3
"""Black-box CB-01A acceptance against the exact prebuilt consumer package."""
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile


def digest(path):
    return "sha256:" + hashlib.sha256(path.read_bytes()).hexdigest()


def run(argv, cwd, env=None, success=True, structured=True):
    p = subprocess.run([str(a) for a in argv], cwd=cwd, env=env, capture_output=True, text=True)
    if (p.returncode == 0) != success:
        raise AssertionError(f"{p.args}: exit {p.returncode}\n{p.stdout}\n{p.stderr}")
    if structured and success:
        return json.loads(p.stdout)
    return p


REGISTRY_FORMAT = 11
REGISTRY_FIELDS = ["epics", "tasks", "worktree_operations", "task_results", "human_qa_records", "integration_authorities", "integration_intents", "integration_attempts", "integration_results", "task_problem_revisions", "task_spec_revisions", "task_spec_assessments", "task_solution_selections", "task_result_spec_bindings", "task_content_publications", "epic_base_versions", "epic_base_updates", "task_queue_events", "task_preparations", "epic_create_operations"]


def empty_registry(version=REGISTRY_FORMAT):
    return f"format_version: {version}\nregistry:\n" + "".join(f"  {key}: []\n" for key in REGISTRY_FIELDS)


def registration(ply, root):
    root.mkdir()
    run([ply, "workspace", "init"], root, structured=False)
    repo = root / "example/main"
    repo.mkdir(parents=True)
    run(["git", "init", "-b", "main"], repo, structured=False)
    (repo / "README.md").write_text("Isolated acceptance fixture.\n")
    run(["git", "add", "README.md"], repo, structured=False)
    run(["git", "-c", "user.name=CB01A fixture", "-c", "user.email=fixture@example.invalid", "-c", "commit.gpgsign=false", "commit", "-m", "Fixture baseline"], repo, structured=False)
    oid = run(["git", "rev-parse", "HEAD"], repo, structured=False).stdout.strip()
    run([ply, "workspace", "project", "add", "example", "--name", "Example", "--wrapper", repo.parent, "--repo", f"example={repo}"], root, structured=False)
    registry = repo.parent / "data/work/work-items.yaml"  # the Project's registry, under its wrapper
    assert not registry.exists()
    run([ply, "workspace", "epic", "adopt", "greetings", "--title", "Greetings", "--project", "example", "--repo", "example", "--worktree", repo, "--ref", "refs/heads/main", "--expected-oid", oid], root, structured=False)
    create = [ply, "workspace", "task", "create", "greeting", "--title", "Named greeting", "--description", "Add an optional name", "--project", "example", "--repo", "example", "--epic", "greetings", "--format", "json"]
    run(create, root)
    # A fresh registry is written in the single supported format.
    assert registry.read_text().startswith(f"format_version: {REGISTRY_FORMAT}\nregistry:\n")
    recorder = dict(actor_claim="CB-01A isolated acceptance", control_surface="fixture only", recorded_at_utc="2026-10-09T00:00:00Z")
    problem_doc = root / "problem.md"
    problem_doc.write_text("# Greeting\n\nThe greeting cannot name the reader.\n")
    problem_draft = dict(kind="WorkspaceTaskProblemDraft@1", schema_version=1, format="json", format_version=1, canonicalization="RFC8785", publication_key="fixture/problem", task_id="greeting", registry_upgrade=None, expected_previous=None, origin=dict(kind="authored"), title="Named greeting", summary="Support an optional name", problem_document_id="problem", documents=[dict(id="problem", source=dict(kind="file", locator=str(problem_doc), sha256=digest(problem_doc), size_bytes=problem_doc.stat().st_size, media_type="text/markdown", git_provenance=None))], sources=[], claims=[], deadline=None, change_reason="Record the fixture problem", recorder=recorder)
    path = root / "problem.json"
    path.write_text(json.dumps(problem_draft))
    problem = run([ply, "workspace", "task", "problem", "record", "greeting", "--file", path, "--format", "json"], root)["outcome_ref"]
    doc = root / "spec.md"
    doc.write_text("# Named greeting\n\nKeep the default greeting and support an optional name.\n")
    draft = dict(kind="WorkspaceTaskSpecDraft@2", schema_version=2, format="json", format_version=1, canonicalization="RFC8785", publication_key="fixture/goal", task_id="greeting", recorder=recorder, registry_upgrade=None, spec_id="greeting-v1", contract_kind="goal", expected_previous=None, problem={key: problem[key] for key in ("revision", "manifest_sha256")}, title="Named greeting", objective="Support an optional name while keeping the default greeting.", documents=[dict(id="spec", source=dict(kind="file", locator=str(doc), sha256=digest(doc), size_bytes=doc.stat().st_size, media_type="text/markdown", git_provenance=None))], design=[dict(document_id="spec", section=None)], requirements=[dict(id="greeting", acceptance="The existing default and named greetings are verified.")], constraints=["No agent starts during this fixture."], executor=dict(provider=None, model=None, effort=None), change_reason="Register isolated acceptance input.")
    path = root / "goal.json"
    path.write_text(json.dumps(draft))
    spec = run([ply, "workspace", "task", "spec", "record", "greeting", "--file", path, "--format", "json"], root)["outcome_ref"]
    scope = ["--project", "example", "--repo", "example", "--epic", "greetings"]
    queue = run([ply, "workspace", "task", "queue", "list", *scope, "--format", "json"], root)
    assert queue["registry"]["format_version"] == REGISTRY_FORMAT and queue["registry"]["required_upgrade"] is None
    q = dict(kind="WorkspaceTaskQueueDraft@2", schema_version=2, publication_key="fixture/queue", project_id="example", repo_id="example", epic_id="greetings", expected_revision=queue["revision"], entries=[dict(task_id="greeting", goal=dict(spec_id=spec["spec_id"], spec={key: spec[key] for key in ("revision", "manifest_sha256")}))], recorder=recorder, registry_upgrade=None)
    path = root / "queue.json"
    path.write_text(json.dumps(q))
    run([ply, "workspace", "task", "queue", "set", "--file", path, "--format", "json"], root)
    before_queue = run([ply, "workspace", "task", "queue", "list", *scope, "--format", "json"], root)
    assert before_queue["pending"][0]["state"] == "ready"
    run([ply, "workspace", "task", "create", "second", "--title", "Unrelated task", "--description", "Retain first goal", *scope, "--format", "json"], root)
    assert run([ply, "workspace", "task", "queue", "list", *scope, "--format", "json"], root)["pending"] == before_queue["pending"]
    run([ply, "workspace", "task", "spec", "show", "greeting", "--spec", "greeting-v1", "--revision", "1", "--format", "json"], root)
    assert run(["git", "worktree", "list", "--porcelain"], repo, structured=False).stdout.count("worktree ") == 1
    assert run(["git", "rev-parse", "HEAD"], repo, structured=False).stdout.strip() == oid
    print(f"CBA-02/04: native registration in registry format {REGISTRY_FORMAT} without upgrade PASS", flush=True)


def main():
    package = Path(sys.argv[1]).resolve(strict=True)
    pin = sys.argv[2]
    assert digest(package / "manifest.json") == pin
    manifest = json.loads((package / "manifest.json").read_text())
    with tempfile.TemporaryDirectory(prefix="ply-cb01a-acceptance-", dir=os.environ["TMPDIR"]) as temporary:
        root = Path(temporary).resolve()
        consumer = root / "consumer"
        consumer.mkdir()
        (consumer / "bin").mkdir()
        dest = consumer / "bin/ply"
        # This obsolete-reader fixture deliberately cannot inspect any registry.
        old = b"#!/bin/sh\nprintf 'unsupported old reader\\n' >&2\nexit 1\n"
        dest.write_bytes(old)
        dest.chmod(0o755)
        env = dict(PATH="/usr/bin:/bin", LC_ALL="C", LANG="C", HOME=str(consumer))
        assert shutil.which("go", path=env["PATH"]) is None
        protected = consumer / ".ply"
        protected.mkdir()
        for name in ("work-items.yaml", "workspace.yaml", "queue", "agent-state", "git-state"):
            (protected / name).write_text("consumer state must not change\n")
        before = {p.name: digest(p) for p in protected.iterdir()}
        install = [package / "ply-install", "--package", package, "--destination", dest]
        run([*install, "--sha256", "sha256:" + "0" * 64], consumer, env, success=False)
        assert dest.read_bytes() == old
        bad = root / "wrong-platform"
        bad.mkdir()
        bad_manifest = dict(manifest, arch="amd64")
        (bad / "manifest.json").write_text(json.dumps(bad_manifest))
        run([package / "ply-install", "--package", bad, "--sha256", digest(bad / "manifest.json"), "--destination", dest], consumer, env, success=False)
        assert dest.read_bytes() == old
        run([package / "ply-install", "--package", package, "--sha256", pin, "--destination", consumer / "missing/ply"], consumer, env, success=False)
        assert dest.read_bytes() == old
        receipt = run([*install, "--sha256", pin], consumer, env)
        assert receipt["state"] == "installed" and digest(dest) == manifest["sha256"]
        assert Path(receipt["backup"]).read_bytes() == old
        identity = run([dest, "version", "--format", "json"], consumer, env)
        for key in ("version", "source_commit", "os", "arch", "sha256", "registry_formats"):
            assert identity[key] == manifest[key], key
        assert not identity["dirty"]
        dest.chmod(0o644)
        run([*install, "--sha256", pin], consumer, env)
        assert run([dest, "version", "--format", "json"], consumer, env) == identity
        run([package / "ply-install", "--rollback", "--sha256", manifest["sha256"], "--destination", dest], consumer, env)
        assert dest.read_bytes() == old
        run([*install, "--sha256", pin], consumer, env)
        assert {p.name: digest(p) for p in protected.iterdir()} == before
        print("CBA-05/06: no-Go consumer install, identity, failures, mode repair, rollback, no state writes PASS", flush=True)
        registration(dest, root / "workspace")
        broken = root / "diagnostics"
        broken.mkdir()
        run([dest, "workspace", "init"], broken, structured=False)
        # One Project; reading its registry observes no Git.
        (broken / "project").mkdir()
        (broken / "ply-workspace.json").write_text(json.dumps({"format_version": 1, "projects": [{"id": "project", "name": "Project", "wrapper": "project", "repos": [{"id": "repo", "path": "project/repo"}]}]}))
        registry = broken / "project/data/work/work-items.yaml"
        registry.parent.mkdir(parents=True)
        for raw, code in ((empty_registry(10), "registry_unsupported_format"), (empty_registry() + "  authorization_override: true\n", "registry_unknown_field"), (empty_registry() + "broken: [", "registry_malformed_yaml"), (empty_registry().replace("epics: []", "epics: null"), "registry_invalid_content")):
            registry.write_text(raw)
            before = registry.read_bytes()
            result = run([dest, "workspace", "epic", "list", "--format", "json"], broken, success=False)
            diagnostic = json.loads(result.stderr)
            assert diagnostic["code"] == code and diagnostic["next_action"]
            assert registry.read_bytes() == before
        print("CBA-03: four distinct diagnostics with no writes PASS", flush=True)


if __name__ == "__main__":
    main()
