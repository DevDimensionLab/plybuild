#!/usr/bin/env python3
"""Register one disposable native Goal through public CLI; optionally start it.

This is a real-provider acceptance exercise, never part of ordinary unit tests.
It does not install Ply, change a real project, fabricate review, or create QA.
"""
import argparse
import datetime
import hashlib
import json
import os
from pathlib import Path
import subprocess
import tempfile


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--ply", required=True, type=Path)
    parser.add_argument("--start", action="store_true")
    args = parser.parse_args()
    binary = args.ply.resolve(strict=True)
    root = Path(tempfile.mkdtemp(prefix="ply-headless-journey-", dir="/private/tmp")).resolve()
    wrapper = root / "fixture"
    repo, epic = wrapper / "main", wrapper / "epic"
    repo.mkdir(parents=True)
    logs = root / "commands"
    logs.mkdir()
    count = 0

    def run(argv, cwd=root, structured=False):
        nonlocal count
        count += 1
        process = subprocess.run(list(map(str, argv)), cwd=cwd, text=True, capture_output=True)
        (logs / f"{count:03d}.json").write_text(json.dumps(dict(argv=list(map(str, argv)), cwd=str(cwd), exit=process.returncode, stdout=process.stdout, stderr=process.stderr), indent=2))
        if process.returncode:
            raise RuntimeError(f"Command failed; preserved at {logs / f'{count:03d}.json'}: {process.stderr}")
        return json.loads(process.stdout) if structured else process.stdout.strip()

    def ply(*argv, cwd=root):
        return run([binary, "workspace", *argv, "--format", "json"], cwd, True)

    def write(name, value):
        path = root / name
        path.write_text(json.dumps(value, indent=2) + "\n")
        return path

    def git(*argv, cwd=repo):
        return run(["git", "-c", "core.hooksPath=/dev/null", *argv], cwd)

    (repo / "README.md").write_text("# Disposable headless worker fixture\nImplement the registered Goal. No external effects.\n")
    (repo / "normalize.py").write_text("#!/usr/bin/env python3\nraise SystemExit('not implemented')\n")
    (repo / "normalize.py").chmod(0o755)
    (repo / "test").mkdir()
    (repo / "test" / "delivery_acceptance.sh").write_text('''#!/bin/sh
set -eu
python3 -I -B - <<'PY'
import json, os, subprocess
binary = os.environ['PLY_CANDIDATE_BINARY']
for data, expected in [('[3, 1, 3, -2]', [-2, 1, 3]), ('[]', []), (' [0, -1, 0] ', [-1, 0]), ('[9007199254740993, 1]', [1, 9007199254740993])]:
    p = subprocess.run([binary], input=data, text=True, capture_output=True)
    assert p.returncode == 0 and json.loads(p.stdout) == expected and not p.stderr, (data, p)
for data in ['[true]', '[1.5]', '{}', '["1"]', 'not json', '[1] trailing']:
    p = subprocess.run([binary], input=data, text=True, capture_output=True)
    assert p.returncode == 2 and p.stdout == '' and p.stderr.strip(), (data, p)
print('10 behavioral acceptance cases passed')
PY
''')
    git("init", "-b", "main")
    git("config", "user.email", "fixture@example.invalid")
    git("config", "user.name", "Ply worker fixture")
    git("add", "--all")
    git("commit", "-m", "Disposable worker fixture baseline")
    oid = git("rev-parse", "HEAD")
    git("worktree", "add", "-b", "epic", epic, oid)
    run([binary, "workspace", "init"])
    run([binary, "workspace", "project", "add", "fixture", "--name", "Headless fixture", "--wrapper", wrapper, "--repo", f"fixture={repo}"])
    run([binary, "workspace", "epic", "adopt", "epic", "--title", "Disposable worker Epic", "--project", "fixture", "--repo", "fixture", "--worktree", epic, "--ref", "refs/heads/epic", "--expected-oid", oid])
    ply("task", "create", "normalize", "--title", "Normalize integer arrays", "--description", "Implement one real small Goal for the authorized DW-01 worker acceptance exercise.", "--epic", "epic", "--project", "fixture", "--repo", "fixture")
    task = ply("task", "show", "normalize")
    spec = root / "spec.md"
    spec.write_text("""# Normalize integer arrays

Implement normalize.py as a standalone executable Python 3 program, using only
the standard library. Read all stdin as one JSON array of integers, remove
duplicates, sort ascending, and write a JSON array followed by one newline.
Accept negative integers, empty arrays, surrounding whitespace and arbitrary
Python integer precision. Boolean values, floats, strings, non-arrays, malformed
JSON and trailing non-whitespace input fail with exit 2, an explanatory stderr
message, and no stdout. Valid input exits 0 without stderr.

Keep and pass every existing acceptance case in test/delivery_acceptance.sh.
Add meaningful regression/unit tests and brief usage documentation. Build the
candidate by copying the self-contained executable to the controller's stated
output path. No dependencies, network, installation, remote effects or human QA.
Independent actual review and native automatic local Epic delivery are required.
""")
    now = datetime.datetime.now(datetime.timezone.utc).isoformat().replace("+00:00", "Z")
    recorder = dict(actor_claim="Human-requested DW-01 acceptance controller", control_surface="Disposable public-CLI worker exercise", recorded_at_utc=now)
    data = spec.read_bytes()
    goal = dict(kind="WorkspaceTaskSpecDraft@2", schema_version=2, format="json", format_version=1, canonicalization="RFC8785", publication_key="headless/normalize/v1", task_id="normalize", recorder=recorder, registry_upgrade=None, spec_id="normalize-v1", contract_kind="goal", expected_previous=None, problem=task["problem"]["head"], title="Normalize integer arrays", objective="Implement the specified executable and pass actual tests and independent review before native automatic local Epic delivery.", documents=[dict(id="spec", source=dict(kind="file", locator=str(spec), sha256="sha256:"+hashlib.sha256(data).hexdigest(), size_bytes=len(data), media_type="text/markdown", git_provenance=None))], design=[dict(document_id="spec", section=None)], requirements=[dict(id="behavior", acceptance="All normalization, invalid-input and exact output requirements in the Spec pass meaningful tests."), dict(id="delivery", acceptance="A distinct real Codex reviewer passes the exact candidate before automatic local integration and retained native closeout; HumanQA stays absent.")], constraints=["Keep this disposable Task local; no network, remote, installed binary or unrelated project changes."], executor=dict(provider="codex", model="gpt-6-astra", effort="max"), change_reason="Authorized real worker exercise, not a synthetic provider or reviewer.", delivery=dict(schema_version=3, mode="local_epic_integration", project_id="fixture", repo_id="fixture", epic_id="epic", target_ref="refs/heads/epic", target_worktree=str(epic), acceptance=dict(schema_version=1, mode="automatic", responsible_actor="Human-requested DW-01 worker", timeout_seconds=120, require_human_qa=False)))
    published = ply("task", "spec", "record", "normalize", "--file", write("goal.json", goal))
    target = ["--project", "fixture", "--repo", "fixture", "--epic", "epic"]
    queue = ply("task", "queue", "list", *target)
    ref = published["outcome_ref"]
    queue_draft = dict(kind="WorkspaceTaskQueueDraft@2", schema_version=2, publication_key="headless/queue/v1", project_id="fixture", repo_id="fixture", epic_id="epic", expected_revision=queue["revision"], entries=[dict(task_id="normalize", goal=dict(spec_id="normalize-v1", spec=dict(revision=ref["revision"], manifest_sha256=ref["manifest_sha256"])))], recorder=recorder, registry_upgrade=queue["registry"]["required_upgrade"])
    ply("task", "queue", "set", "--file", write("queue.json", queue_draft))
    command = [binary, "workflow", "execute", "--headless", "--spec", "normalize-v1", "--worker-timeout", "1800", "--format", "json"]
    preview = run([*command, "--check"], epic, True)
    if preview["state"] != "ready":
        raise RuntimeError(f"Worker preview is not ready: {preview}")
    result = dict(root=str(root), binary=str(binary), epic=str(epic), initial_oid=oid, preview=preview, execute_argv=list(map(str, command)))
    if args.start:
        result["start"] = run(command, epic, True)
    write("journey.json", result)
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
