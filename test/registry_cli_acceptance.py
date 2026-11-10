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


def empty_registry(version):
    fields = ["epics", "tasks", "worktree_operations"]
    if version >= 2:
        fields += ["task_results", "human_qa_records", "integration_authorities", "integration_intents", "integration_attempts", "integration_results"]
    if version >= 3:
        fields += ["task_spec_policies", "task_problem_revisions", "task_spec_revisions", "task_spec_assessments", "task_solution_selections", "task_result_spec_bindings", "task_content_publications"]
    if version >= 4:
        fields += ["epic_base_versions", "epic_base_updates", "task_queue_events", "task_preparations"]
    return f"format_version: {version}\n" + "".join(f"{key}: []\n" for key in fields)


def registration(ply, root, version, package, pin):
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
    registry = root / ".ply/work-items.yaml"
    # Seed historical formats only in this disposable fixture, before registrations.
    registry.write_text(empty_registry(version))
    run([ply, "workspace", "epic", "adopt", "greetings", "--title", "Greetings", "--project", "example", "--repo", "example", "--worktree", repo, "--ref", "refs/heads/main", "--expected-oid", oid], root, structured=False)
    create = [ply, "workspace", "task", "create", "greeting", "--title", "Named greeting", "--description", "Add an optional name", "--project", "example", "--repo", "example", "--epic", "greetings", "--format", "json"]
    if version < 3:
        before = registry.read_bytes()
        run(create, root, success=False)
        assert registry.read_bytes() == before
        create += ["--upgrade-store", digest(registry)]
    run(create, root)
    task = run([ply, "workspace", "task", "show", "greeting", "--format", "json"], root)
    doc = root / "spec.md"
    doc.write_text("# Named greeting\n\nKeep the default greeting and support an optional name.\n")
    recorder = dict(actor_claim="CB-01A isolated acceptance", control_surface="fixture only", recorded_at_utc="2026-10-09T00:00:00Z")
    draft = dict(kind="WorkspaceTaskSpecDraft@2", schema_version=2, format="json", format_version=1, canonicalization="RFC8785", publication_key="fixture/goal", task_id="greeting", recorder=recorder, registry_upgrade=None, spec_id="greeting-v1", contract_kind="goal", expected_previous=None, problem=task["problem"]["head"], title="Named greeting", objective="Support an optional name while keeping the default greeting.", documents=[dict(id="spec", source=dict(kind="file", locator=str(doc), sha256=digest(doc), size_bytes=doc.stat().st_size, media_type="text/markdown", git_provenance=None))], design=[dict(document_id="spec", section=None)], requirements=[dict(id="greeting", acceptance="The existing default and named greetings are verified.")], constraints=["No agent starts during this fixture."], executor=dict(provider=None, model=None, effort=None), change_reason="Register isolated acceptance input.")
    path = root / "goal.json"
    path.write_text(json.dumps(draft))
    spec = run([ply, "workspace", "task", "spec", "record", "greeting", "--file", path, "--format", "json"], root)["outcome_ref"]
    scope = ["--project", "example", "--repo", "example", "--epic", "greetings"]
    queue = run([ply, "workspace", "task", "queue", "list", *scope, "--format", "json"], root)
    q = dict(kind="WorkspaceTaskQueueDraft@2", schema_version=2, publication_key="fixture/queue", project_id="example", repo_id="example", epic_id="greetings", expected_revision=queue["revision"], entries=[dict(task_id="greeting", goal=dict(spec_id=spec["spec_id"], spec={key: spec[key] for key in ("revision", "manifest_sha256")}))], recorder=recorder, registry_upgrade=queue["registry"]["required_upgrade"])
    path = root / "queue.json"
    path.write_text(json.dumps(q))
    run([ply, "workspace", "task", "queue", "set", "--file", path, "--format", "json"], root)
    baseline = registry.read_bytes()
    before_queue = run([ply, "workspace", "task", "queue", "list", *scope, "--format", "json"], root)
    assert before_queue["pending"][0]["state"] == "ready"
    assert before_queue["registry"]["format_version"] == 4
    args = [ply, "workspace", "registry", "upgrade", "--package", package, "--sha256", pin, "--format", "json"]
    preview = run([*args, "--check"], root)
    assert registry.read_bytes() == baseline
    run([*args, "--apply"], root, success=False)
    assert registry.read_bytes() == baseline
    result = run([*args, "--apply", "--confirm", preview["confirmation"]], root)
    assert result["state"] == "upgraded"
    assert Path(result["backup"]).read_bytes() == baseline
    after = registry.read_bytes()
    assert run([*args, "--apply", "--confirm", preview["confirmation"]], root)["state"] == "already_upgraded"
    assert registry.read_bytes() == after
    after_queue = run([ply, "workspace", "task", "queue", "list", *scope, "--format", "json"], root)
    assert before_queue["pending"] == after_queue["pending"]
    assert before_queue["revision"] == after_queue["revision"]
    run([ply, "workspace", "task", "create", "second", "--title", "Unrelated task", "--description", "Retain first goal", *scope, "--format", "json"], root)
    assert run([ply, "workspace", "task", "queue", "list", *scope, "--format", "json"], root)["pending"] == before_queue["pending"]
    run([ply, "workspace", "task", "spec", "show", "greeting", "--spec", "greeting-v1", "--revision", "1", "--format", "json"], root)
    assert run(["git", "worktree", "list", "--porcelain"], repo, structured=False).stdout.count("worktree ") == 1
    assert run(["git", "rev-parse", "HEAD"], repo, structured=False).stdout.strip() == oid
    print(f"CBA-02/04: native registration format {version} -> 4 -> 5, history and retry PASS", flush=True)


def main():
    package = Path(sys.argv[1]).resolve(strict=True)
    pin = sys.argv[2]
    assert digest(package / "manifest.json") == pin
    manifest = json.loads((package / "manifest.json").read_text())
    with tempfile.TemporaryDirectory(prefix="ply-cb01a-acceptance-", dir="/tmp") as temporary:
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
        for version in range(1, 5):
            registration(dest, root / f"workspace-v{version}", version, package, pin)
        broken = root / "diagnostics"
        broken.mkdir()
        run([dest, "workspace", "init"], broken, structured=False)
        registry = broken / ".ply/work-items.yaml"
        for raw, code in ((empty_registry(1).replace("version: 1", "version: 99"), "registry_unsupported_format"), (empty_registry(1) + "authorization_override: true\n", "registry_unknown_field"), (empty_registry(1) + "broken: [", "registry_malformed_yaml"), (empty_registry(1).replace("epics: []", "epics: null"), "registry_invalid_content")):
            registry.write_text(raw)
            before = registry.read_bytes()
            result = run([dest, "workspace", "epic", "list", "--format", "json"], broken, success=False)
            diagnostic = json.loads(result.stderr)
            assert diagnostic["code"] == code and diagnostic["next_action"]
            assert registry.read_bytes() == before
        print("CBA-03: four distinct diagnostics with no writes PASS", flush=True)


if __name__ == "__main__":
    main()
