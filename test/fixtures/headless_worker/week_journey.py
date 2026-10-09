#!/usr/bin/env python3
"""Prepare, start, and inspect a disposable real-provider AW-01 journey.

Model calls require --start. --service additionally installs the fixture's named
local service. --inspect only observes an existing fixture. No real project,
provider configuration, installed binary, or synthetic review is modified.
"""
import argparse
import datetime
import hashlib
import json
import os
from pathlib import Path
import signal
import subprocess
import sys
import tempfile
import time


ADDY_COMMIT = "1401c8b8030e023baeebb31781a6653fe8e93026"
ADDY_MANIFEST = "29466c59aa8bc8c3a6eed0abc3765d009cd0d1095df4785f6fb948831e703561"
TASKS = ("normalize", "normalize-descending")
MODEL, EFFORT = "gpt-6-astra", "max"
TERMINAL = {"completed", "completed_with_blocks", "blocked", "expired", "paused", "stopped"}

# Preserve the existing DW-01 ten-case acceptance block byte for byte. Both
# Goals must retain it; independent final checks also exercise the executable.
BASELINE_SCRIPT = '''#!/bin/sh
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
'''

SPEC_NORMALIZE = """# Normalize integer arrays

Implement normalize.py as a self-contained executable Python 3 program using
only the standard library. Read all stdin as one JSON array of integers, remove
duplicates, sort ascending, and write one JSON array followed by one newline.
Accept negative integers, empty arrays, surrounding whitespace and arbitrary
Python integer precision. Boolean values, floats, strings, non-arrays, malformed
JSON and trailing non-whitespace input fail with exit 2, explanatory stderr, and
no stdout. Valid input exits 0 without stderr.

Preserve the complete existing test/delivery_acceptance.sh byte for byte as the
prefix of that file. Additional checks may be appended after its final line.
Keep and pass all ten existing cases. Add meaningful regression/unit tests and
brief usage documentation. Build the candidate by copying the self-contained
executable to the controller's stated output path. No dependencies, network,
installation, remote effects or human QA. Actual independent review and native
automatic local Epic delivery are required. All files belong to this disposable
fixture. Do not implement any descending-order option in this first Goal.
"""

SPEC_DESCENDING = """# Add descending normalization

Extend the delivered normalize.py with the single optional flag --descending.
Without arguments, preserve every behavior of the ascending normalization Goal:
one JSON integer array from stdin, deduplication, arbitrary integer precision,
the same valid/invalid inputs, exit codes, stderr rules and final newline.
With exactly one --descending argument, return the unique integers in descending
order. Duplicate flags, unknown flags, positional arguments, and '--descending'
followed by any other argument fail with exit 2, explanatory stderr, no stdout.
Empty arrays and one-element arrays work in both orders. Input type validation
also applies in descending mode: bool, float, string, non-array, malformed JSON
and trailing non-whitespace input remain errors.

This Goal depends on the completed normalize Goal and must start from its updated
Epic base. Keep the original complete ten-case acceptance script byte for byte
as its prefix and append executable checks for descending behavior and argument
errors. Preserve existing regression tests and usage behavior; add focused tests
and document the new flag. Do not merely add documentation or loosen tests.
Build the self-contained executable at the controller's exact output path.
Actual separate review and automatic local Epic delivery are required. No
dependencies, network, installation, remote effects, or human QA.
"""


def require(condition, message):
    if not condition:
        raise RuntimeError(message)


def sha(raw):
    return "sha256:" + hashlib.sha256(raw).hexdigest()


def now():
    return datetime.datetime.now(datetime.timezone.utc).isoformat().replace("+00:00", "Z")


def read_json(path):
    return json.loads(Path(path).read_bytes())


def write_json(path, value):
    Path(path).write_text(json.dumps(value, indent=2, ensure_ascii=False) + "\n")
    return Path(path)


def inside(root, path):
    path = Path(path).resolve(strict=True)
    require(path != root and root in path.parents, f"Fixture evidence escaped {root}: {path}")
    return path


def bound(root, binding):
    path = inside(root, binding["locator"])
    raw = path.read_bytes()
    require(sha(raw) == binding["sha256"], f"Evidence digest mismatch: {path}")
    return raw


class Commands:
    def __init__(self, binary, root):
        self.binary, self.root = binary, root
        self.logs = root / "commands"
        self.logs.mkdir(exist_ok=True)
        self.count = max((int(p.stem) for p in self.logs.glob("*.json")), default=0)

    def run(self, argv, cwd=None, structured=False, raw=False):
        argv = list(map(str, argv))
        cwd = cwd or self.root
        self.count += 1
        log = self.logs / f"{self.count:05d}.json"
        started = now()
        try:
            p = subprocess.run(argv, cwd=cwd, text=not raw, capture_output=True, timeout=90)
        except subprocess.TimeoutExpired:
            write_json(log, dict(argv=argv, cwd=str(cwd), started_at=started,
                                finished_at=now(), error="CLI timed out after 90 seconds"))
            raise
        write_json(log, dict(argv=argv, cwd=str(cwd), started_at=started, finished_at=now(),
                             exit=p.returncode, stdout=p.stdout.decode(errors="replace") if raw else p.stdout,
                             stderr=p.stderr.decode(errors="replace") if raw else p.stderr))
        require(p.returncode == 0, f"Command failed; inspect {log}: {p.stderr}")
        return p.stdout if raw else (json.loads(p.stdout) if structured else p.stdout.strip())

    def ply(self, *argv, cwd=None):
        return self.run([self.binary, *argv, "--format", "json"], cwd, True)

    def workspace(self, *argv):
        return self.ply("workspace", *argv)

    def git(self, repo, *argv):
        return self.run(["git", "-c", "core.hooksPath=/dev/null", *argv], repo)


def executable_identity(commands, expected_commit=None, expected_sha=None):
    identity = commands.ply("version")
    actual = sha(commands.binary.read_bytes())
    require(identity["sha256"] == actual, "Ply self-reported digest differs from invoked executable")
    if expected_commit:
        require(identity["source_commit"] == expected_commit, "Unexpected Ply source commit")
    if expected_sha:
        require(actual == expected_sha, "Unexpected Ply executable SHA-256")
    return dict(path=str(commands.binary), **identity)


def goal_draft(task_id, problem, spec, epic, recorder):
    return dict(
        kind="WorkspaceTaskSpecDraft@2", schema_version=2, format="json", format_version=1,
        canonicalization="RFC8785", publication_key=f"aw01/{task_id}/v1", task_id=task_id,
        recorder=recorder, registry_upgrade=None, spec_id=task_id + "-v1", contract_kind="goal",
        expected_previous=None, problem=problem, title=task_id.replace("-", " ").capitalize(),
        objective="Implement the preserved behavior, pass actual tests and separate review, and deliver automatically to the disposable Epic.",
        documents=[dict(id="spec", source=dict(kind="file", locator=str(spec),
            sha256=sha(spec.read_bytes()), size_bytes=spec.stat().st_size,
            media_type="text/markdown", git_provenance=None))],
        design=[dict(document_id="spec", section=None)],
        requirements=[dict(id="behavior", acceptance="All preserved Spec behaviors pass, with the original ten-case acceptance block unchanged and meaningful added tests."),
                      dict(id="delivery", acceptance="A distinct actual Codex reviewer passes the exact tested candidate before native automatic local Epic delivery. HumanQA remains absent.")],
        constraints=["Only this disposable Task; no network, installation, push, notifications, or unrelated project effects."],
        executor=dict(provider="codex", model=MODEL, effort=EFFORT),
        change_reason="Human-authorized AW-01 real two-Spec acceptance fixture.",
        delivery=dict(schema_version=3, mode="local_epic_integration", project_id="fixture",
            repo_id="fixture", epic_id="epic", target_ref="refs/heads/epic", target_worktree=str(epic),
            acceptance=dict(schema_version=1, mode="automatic", responsible_actor="AW-01 disposable acceptance controller",
                            timeout_seconds=120, require_human_qa=False)))


def prepare(args):
    binary = args.ply.resolve(strict=True)
    root = Path(tempfile.mkdtemp(prefix="ply-aw01-week-", dir="/private/tmp")).resolve()
    print(f"Disposable fixture: {root}", file=sys.stderr, flush=True)
    commands = Commands(binary, root)
    report = dict(kind="PlyAW01WeekJourney@1", root=str(root), created_at=now(),
                  state="preparing", model_calls_requested=False, controls=[], interruptions=[], service_observations=[], errors=[])
    write_json(root / "journey.json", report)
    report["executable"] = executable_identity(commands, args.expected_source_commit, args.expected_sha256)
    commands.run([binary, "workflow", "week", "start", "--help"])
    wrapper, repo, epic = root / "fixture", root / "fixture/main", root / "fixture/epic"
    repo.mkdir(parents=True)
    (repo / "test").mkdir()
    (repo / "README.md").write_text("# Disposable AW-01 multi-Spec fixture\nImplement only the native Goal selected for this attempt.\n")
    (repo / "normalize.py").write_text("#!/usr/bin/env python3\nraise SystemExit('not implemented')\n")
    (repo / "normalize.py").chmod(0o755)
    (repo / "test/delivery_acceptance.sh").write_text(BASELINE_SCRIPT)
    commands.git(repo, "init", "-b", "main")
    commands.git(repo, "config", "user.email", "fixture@example.invalid")
    commands.git(repo, "config", "user.name", "Ply AW-01 fixture")
    commands.git(repo, "add", "--all")
    commands.git(repo, "commit", "-m", "Disposable AW-01 fixture baseline")
    oid = commands.git(repo, "rev-parse", "HEAD")
    commands.git(repo, "worktree", "add", "-b", "epic", epic, oid)
    commands.run([binary, "workspace", "init"])
    commands.run([binary, "workspace", "project", "add", "fixture", "--name", "AW-01 fixture",
                  "--wrapper", wrapper, "--repo", f"fixture={repo}"])
    commands.run([binary, "workspace", "epic", "adopt", "epic", "--title", "Disposable AW-01 Epic",
                  "--project", "fixture", "--repo", "fixture", "--worktree", epic,
                  "--ref", "refs/heads/epic", "--expected-oid", oid])
    recorder = dict(actor_claim="Human-requested AW-01 acceptance controller",
                    control_surface="Disposable public-CLI multi-Spec exercise", recorded_at_utc=now())
    selections = []
    for task_id, text in zip(TASKS, (SPEC_NORMALIZE, SPEC_DESCENDING)):
        commands.workspace("task", "create", task_id, "--title", task_id.replace("-", " "),
                           "--description", "Real developer and separate review in the disposable AW-01 journey.",
                           "--epic", "epic", "--project", "fixture", "--repo", "fixture")
        task = commands.workspace("task", "show", task_id)
        spec = root / f"{task_id}-spec.md"
        spec.write_text(text)
        draft = write_json(root / f"{task_id}-goal.json", goal_draft(task_id, task["problem"]["head"], spec, epic, recorder))
        published = commands.workspace("task", "spec", "record", task_id, "--file", draft)
        reference = published["outcome_ref"]
        selections.append(dict(task_id=task_id, goal=dict(spec_id=task_id + "-v1", spec=dict(
            revision=reference["revision"], manifest_sha256=reference["manifest_sha256"])),
            depends_on=[] if task_id == TASKS[0] else [TASKS[0]]))
    queue = commands.workspace("task", "queue", "list", "--project", "fixture", "--repo", "fixture", "--epic", "epic")
    queue_file = write_json(root / "queue.json", dict(kind="WorkspaceTaskQueueDraft@2", schema_version=2,
        publication_key="aw01/queue/v1", project_id="fixture", repo_id="fixture", epic_id="epic",
        expected_revision=queue["revision"], entries=[dict(task_id=s["task_id"], goal=s["goal"]) for s in selections],
        recorder=recorder, registry_upgrade=queue["registry"]["required_upgrade"]))
    commands.workspace("task", "queue", "set", "--file", queue_file)
    contract = dict(kind="ply.workflow.week-contract", schema_version=1, name="aw01-journey",
        workspace=str(root), return_worktree=str(epic), return_ref="refs/heads/epic", project_id="fixture",
        repo_id="fixture", epic_id="epic", selected=selections, model=MODEL, effort=EFFORT,
        duration_seconds=args.duration, max_attempts=3, agent_timeout_seconds=1800, max_provider_retries=20)
    contract_file = write_json(root / "week-contract.json", contract)
    preview = commands.ply("workflow", "week", "start", "--file", contract_file, "--check", cwd=epic)
    report.update(state="prepared", epic=str(epic), repo=str(repo), initial_oid=oid,
        baseline_sha256=sha(BASELINE_SCRIPT.encode()), contract=str(contract_file), preview=preview,
        service_requested=args.service, exercise_controls=args.exercise_controls,
        interruption_requested=args.interrupt_controller)
    write_json(root / "journey.json", report)
    if args.start:
        report["model_calls_requested"] = True
        write_json(root / "journey.json", report)
        flags = ["--service"] if args.service else []
        state = commands.ply("workflow", "week", "start", "--file", contract_file, *flags, cwd=epic)
        report.update(state="started", week_id=state["id"], start=state, started_at=now())
        if args.service:
            report["service_observations"].append(commands.ply("workflow", "week", "service", "status", state["id"], "--workspace", root))
        report["control_argv"] = {action: [str(binary), "workflow", "week", action, state["id"],
            "--workspace", str(root), "--format", "json"] for action in ("show", "pause", "stop", "resume")}
        report["service_argv"] = {action: [str(binary), "workflow", "week", "service", action, state["id"],
            "--workspace", str(root), "--format", "json"] for action in ("status", "stop", "start", "uninstall")}
        write_json(root / "journey.json", report)
    return commands, report


def verify_prompt(raw, role):
    text = raw.decode()
    require(text.startswith("# Pinned Ply worker role package\n"), "Prompt did not start with the preserved role package")
    require(f"Role: {role}\n" in text, f"Missing {role} prompt role")
    require(f"Source commit: {ADDY_COMMIT}\n" in text, "Prompt Addy source commit changed")
    require(f"Manifest SHA-256: {ADDY_MANIFEST}\n" in text, "Prompt manifest pin changed")
    marker = "## Embedded provenance manifest\n\n```json\n"
    manifest_bytes = text.split(marker, 1)[1].split("```\n\n## Upstream material", 1)[0].encode()
    require(sha(manifest_bytes) == "sha256:" + ADDY_MANIFEST, "Actual prompt manifest bytes differ")
    manifest = json.loads(manifest_bytes)
    require(manifest["source_commit"] == ADDY_COMMIT and manifest["license"] == "MIT", "Missing source/license provenance")
    selection = manifest["roles"][role]
    files = {item["path"]: item for item in manifest["files"]}
    for name in selection["adapters"] + selection["files"]:
        begin, end = f"--- BEGIN EMBEDDED FILE: {name} ---\n", f"\n--- END EMBEDDED FILE: {name} ---"
        require(text.count(begin) == 1, f"Missing or duplicated prompt file: {name}")
        content = text.split(begin, 1)[1].split(end, 1)[0].encode()
        require(len(content) == files[name]["bytes"] and sha(content) == "sha256:" + files[name]["sha256"],
                f"Actual prompt source bytes differ: {name}")
    return dict(role=role, source_commit=ADDY_COMMIT, manifest_sha256=ADDY_MANIFEST, prompt_sha256=sha(raw))


def verify_behavior(binary, descending=False):
    passed = 0
    # The candidate is a self-contained program that needs neither writes nor
    # network. These independent checks cannot alter candidate or native files.
    prefix = ["/usr/bin/sandbox-exec", "-p", "(version 1) (allow default) (deny network*) (deny file-write*)"] if sys.platform == "darwin" else []
    for argv in ([[], ["--descending"]] if descending else [[]]):
        for data, expected in [("[3, 1, 3, -2]", [-2, 1, 3]), ("[]", []),
                               (" [0, -1, 0] ", [-1, 0]), ("[9007199254740993, 1]", [1, 9007199254740993]),
                               ("[7,7]", [7])]:
            p = subprocess.run(prefix + [str(binary)] + argv, input=data, text=True, capture_output=True, timeout=10)
            want = list(reversed(expected)) if argv else expected
            require(p.returncode == 0 and not p.stderr and p.stdout.endswith("\n") and len(p.stdout.splitlines()) == 1
                    and json.loads(p.stdout) == want, f"Behavior failed: {argv}, {data}, {p}")
            passed += 1
        for data in ["[true]", "[1.5]", "{}", '["1"]', "not json", "[1] trailing", "[[]]", "null"]:
            p = subprocess.run(prefix + [str(binary)] + argv, input=data, text=True, capture_output=True, timeout=10)
            require(p.returncode == 2 and not p.stdout and p.stderr.strip(), f"Invalid-input behavior failed: {argv}, {data}, {p}")
            passed += 1
    if descending:
        for argv in [["--unknown"], ["--descending", "--descending"], ["extra"], ["--descending", "extra"]]:
            p = subprocess.run(prefix + [str(binary)] + argv, input="[1]", text=True, capture_output=True, timeout=10)
            require(p.returncode == 2 and not p.stdout and p.stderr.strip(), f"Argument behavior failed: {argv}, {p}")
            passed += 1
    return dict(passed_cases=passed, network_and_writes_denied=sys.platform == "darwin")


def verify_run(commands, report, entry, index):
    root, run_id = commands.root, entry["run_id"]
    run = commands.ply("workflow", "worker", "show", run_id, "--workspace", root)
    directory = inside(root, run["paths"]["run_root"])
    state = read_json(directory / "state.json")
    request, worker = state["request"], run["worker"]
    require(run["delivery"]["phase"] == "completed" and worker["phase"] == "completed", f"Run incomplete: {run_id}")
    require(request["workspace_root"] == str(root), "Run used another workspace")
    goal = json.loads(bound(root, request["delivery"]["goal"]))
    selected = read_json(report["contract"])["selected"][index]
    require(goal["task_id"] == entry["task_id"] and goal["spec_id"] == selected["goal"]["spec_id"]
            and request["delivery"]["goal"]["sha256"] == selected["goal"]["spec"]["manifest_sha256"],
            "Run was not bound to the selected immutable Goal")
    require(request["runtime"]["model"] == MODEL and request["delivery"]["reasoning_effort"] == EFFORT, "Explicit model choice changed")
    require(request["runtime"]["ply_executable"]["sha256"] == report["executable"]["sha256"], "Worker control executable differs")
    require(entry["base_oid"] == state["observed"]["epic"]["oid"], "Week/base observation differs from native Run")
    agreement = request["delivery"]["agreement"]
    require(agreement["mode"] == "local_epic_integration" and agreement["acceptance"]["mode"] == "automatic"
            and agreement["acceptance"]["require_human_qa"] is False, "Run constructed a human acceptance gate")
    require(run["delivery_status"]["human_judgment"] is None and run["delivery"].get("human_integration") is None,
            "Unexpected human acceptance/integration evidence")
    candidates = run["delivery"]["candidates"]
    require(candidates and all(c["human_qa"] is None for c in candidates), "Unexpected HumanQA record")
    candidate = candidates[-1]
    require(candidate["oid"] == worker["candidate_oid"] and candidate["tree"] == worker["candidate_tree"], "Final candidate identity differs")
    require(candidate["integration"] and run.get("closeout"), "Native integration or closeout receipt is missing")
    integration = json.loads(bound(root, candidate["integration"]))
    verification = json.loads(bound(root, candidate["verification"]))
    require(run["delivery_status"]["acceptance"]["outcome"] == "pass", "Automatic acceptance did not pass")
    review = json.loads(bound(root, worker["review"]))
    proof = json.loads(bound(root, worker["review_proof"]))
    require(proof["run_id"] == run_id and proof["candidate_oid"] == candidate["oid"]
            and proof["candidate_tree"] == candidate["tree"] and proof["review"] == worker["review"], "Review proof identity differs")
    reviewer_result = json.loads(bound(root, proof["result"]))
    require(reviewer_result["intent_sha256"] == proof["job"]["intent_sha256"], "Reviewer result belongs to another preserved intent")
    reviewer_session = reviewer_result["thread_id"]
    require(review["decision"] == "passed" and not review["findings"] and not review["open_actionable_findings"], "Open final review findings")
    require(review["candidate_oid"] == candidate["oid"] and review["candidate_tree"] == candidate["tree"], "Review assessed another candidate")
    require(reviewer_session == review["reviewer_session_id"] and reviewer_session not in worker["developer_sessions"], "Reviewer is not an independent actual session")
    attempts, developer_sessions, first_developer = [], [], True
    for intent_path in sorted((directory / "worker/attempts").glob("*/intent.json")):
        intent, raw_intent = read_json(intent_path), intent_path.read_bytes()
        result_path = intent_path.parent / "result.json"
        require(result_path.is_file(), f"Attempt has no terminal result: {intent_path.parent}")
        result = read_json(result_path)
        require(result["run_id"] == run_id and result["intent_sha256"] == sha(raw_intent), "Attempt/result binding differs")
        kind = intent["kind"]
        row = dict(kind=kind, job_id=intent["job_id"], exit=result["exit"], error=result.get("error", ""))
        if kind in ("developer", "reviewer"):
            prompt = bound(root, intent["prompt"])
            role = "reviewer" if kind == "reviewer" else ("correction" if b"Role: correction\n" in prompt else "developer")
            if kind == "developer" and first_developer:
                require(role == "developer", "Initial implementation received a correction role")
                first_developer = False
            row["skills"] = verify_prompt(prompt, role)
            argv = intent["argv"]
            require("--model" in argv and argv[argv.index("--model") + 1] == MODEL
                    and 'model_reasoning_effort="max"' in argv and "--ignore-user-config" in argv,
                    "Provider argv did not preserve model/config boundaries")
            stdout = bound(root, result["stdout"])
            bound(root, result["stderr"])
            if result["exit"] == 0 and not result.get("error"):
                require(result["thread_id"] and result["process"]["quiescence"] is True, "Provider completion is unproven")
                events = [json.loads(line) for line in stdout.splitlines() if line.strip()]
                require([e.get("thread_id") for e in events if e.get("type") == "thread.started"] == [result["thread_id"]]
                        and any(e.get("type") == "turn.completed" for e in events), "Actual provider completion is missing")
                answer = json.loads(bound(root, result["answer"]))
                row.update(session_id=result["thread_id"], answer=answer)
                if kind == "developer":
                    developer_sessions.append(result["thread_id"])
            policy = read_json(intent_path.parent / "tool-policy.json")
            require(all(not s["enabled"] for s in policy["mcp_servers"] or []), "MCP was enabled for the worker")
        attempts.append(row)
    require(developer_sessions and len(set(developer_sessions)) == len(developer_sessions), "Missing or reused developer session")
    require(set(worker["developer_sessions"]) == set(developer_sessions) and reviewer_session not in developer_sessions,
            "Native developer/reviewer session evidence differs")
    require(any(a.get("session_id") == reviewer_session and a.get("answer", {}).get("decision") == "passed" for a in attempts), "Final review has no matching successful actual attempt")
    probes = [read_json(p) for p in (directory / "worker/policy").glob("probe-*.json")]
    require(probes and all(p["passed"] is True for p in probes), "Actual runtime probe missing or failed")
    for interruption in report.get("interruptions", []):
        if interruption["run_id"] != run_id:
            continue
        resumed = read_json(inside(root, interruption["attempt_directory"]) / "result.json")
        require(resumed["intent_sha256"] == interruption["intent_sha256"] and resumed["exit"] == 0
                and worker["generation"] > interruption["generation"], "Interrupted attempt was not recovered under the same intent")
        interruption["recovered_result"] = resumed
        interruption["observed_generation_after"] = worker["generation"]
    target = inside(root, state["observed"]["target"]["worktree_locator"])
    original = commands.run(["git", "show", f"{candidate['oid']}:test/delivery_acceptance.sh"], target, raw=True)
    require(original.startswith(BASELINE_SCRIPT.encode()), "Original acceptance block was weakened or rewritten")
    candidate_binary = inside(root, worker["candidate_binary"])
    source = commands.run(["git", "show", f"{candidate['oid']}:normalize.py"], target, raw=True)
    require(candidate_binary.read_bytes() == source, "Built executable differs from committed implementation")
    behavior = verify_behavior(candidate_binary, descending=index == 1)
    return dict(task_id=entry["task_id"], run_id=run_id, base_oid=state["observed"]["epic"]["oid"],
                base_revision=entry["base_revision"], candidate_oid=candidate["oid"], candidate_tree=candidate["tree"],
                candidate_binary_sha256=sha(candidate_binary.read_bytes()), developer_sessions=developer_sessions,
                reviewer_session=reviewer_session, attempts=attempts, behavior=behavior,
                integration=integration, verification=verification, closeout=run["closeout"], runtime_probes=probes)


def inspect(commands, report, state):
    require(state["phase"] == "completed", f"Week is {state['phase']}, not completed")
    require([e["task_id"] for e in state["entries"]] == list(TASKS)
            and all(e["state"] == "completed" for e in state["entries"]), "Selected Tasks did not both complete")
    rows = [verify_run(commands, report, entry, i) for i, entry in enumerate(state["entries"])]
    require(rows[0]["run_id"] != rows[1]["run_id"], "Dependent Tasks reused one Run")
    require(rows[0]["base_oid"] == report["initial_oid"], "First Spec used an unexpected base")
    require(rows[1]["base_oid"] == rows[0]["candidate_oid"] and rows[1]["base_revision"] > rows[0]["base_revision"], "Second Spec did not start automatically from the updated Epic base")
    all_sessions = [s for row in rows for s in row["developer_sessions"] + [row["reviewer_session"]]]
    require(len(set(all_sessions)) == len(all_sessions), "Provider sessions were reused across Tasks")
    native_runs = list((commands.root / ".ply/workflow-runs/runs").glob("*/state.json"))
    require({p.parent.name for p in native_runs} == {r["run_id"] for r in rows}, "Unexpected or duplicate native Run")
    head = commands.git(report["epic"], "rev-parse", "HEAD")
    require(head == rows[-1]["candidate_oid"], "Final Epic head differs from delivered candidate")
    require(commands.git(report["epic"], "status", "--porcelain") == "", "Final disposable Epic is dirty")
    require(commands.git(report["repo"], "rev-parse", "HEAD") == report["initial_oid"], "Fixture main branch was modified")
    queue = commands.workspace("task", "queue", "list", "--project", "fixture", "--repo", "fixture", "--epic", "epic")
    require(not queue["pending"] and queue["current"] is None, "Native queue was not closed")
    base = commands.workspace("epic", "base", "show", "epic", "--repo", "fixture")
    require(base["freshness"] == "fresh" and base["versions"][-1]["oid"] == head
            and [v["oid"] for v in base["versions"]] == [report["initial_oid"]] + [r["candidate_oid"] for r in rows],
            "Native Epic base history differs from the two actual deliveries")
    binding = read_json(inside(commands.root, commands.root / ".ply/week-runs" / report["week_id"] / "binding.json"))
    contract = read_json(report["contract"])
    require(binding["contract"] == contract and binding["contract_sha256"] == state["contract_sha256"]
            and binding["control"]["sha256"] == report["executable"]["sha256"], "Frozen week execution contract differs")
    started, deadline = (datetime.datetime.fromisoformat(binding[k].replace("Z", "+00:00")) for k in ("started_at", "deadline"))
    require((deadline - started).total_seconds() == contract["duration_seconds"], "Visible deadline differs from approved duration")
    if report.get("exercise_controls"):
        require([c["action"] for c in report["controls"]] == ["pause", "resume", "stop", "resume"]
                and all("observed_boundary" in c for c in report["controls"] if c["action"] != "resume"), "Requested active controls remain unproven")
        for control in report["controls"]:
            if control.get("via_service"):
                require(control["service_stopped"]["state"] == "installed"
                        and state["generation"] > control["observed_boundary"]["generation"],
                        "Named service was not unloaded and restarted after the safe stop boundary")
    require(not report.get("interruption_requested") or len(report["interruptions"]) == 1, "Requested controller interruption remains unproven")
    if report["service_requested"]:
        services = report["service_observations"]
        require(any(s["state"] == "running" for s in services)
                and all(s["label"] == "com.ply.week." + state["id"] and s["control_sha256"] == report["executable"]["sha256"] for s in services),
                "Named service runtime or frozen executable is unproven")
    report.update(state="passed", observed_at=now(), week_status=state, runs=rows, final_oid=head,
                  final_queue=queue, final_base=base, actual_start=binding["started_at"], deadline=binding["deadline"],
                  observed_seconds=(datetime.datetime.now(datetime.timezone.utc) - started).total_seconds(),
                  executable_after=executable_identity(commands))
    require(report["executable_after"]["sha256"] == report["executable"]["sha256"], "Invoked binary changed while Runs were active")
    return report


def interrupt_controller(commands, report, run):
    """Kill only the observed native controller in this generated fixture.

    The provider and inherited-lock attempt runner are deliberately untouched.
    PID, full command, start identity, UID and current native generation/job are
    checked immediately before the signal. There is no arbitrary PID option.
    """
    root, worker = commands.root, run["worker"]
    require(root.name.startswith("ply-aw01-week-") and report["root"] == str(root), "Interruption requires a generated fixture")
    run_id, pid = run["run_id"], worker["owner_pid"]
    directory = inside(root, run["paths"]["run_root"])
    state = read_json(directory / "state.json")
    control = Path(state["request"]["runtime"]["ply_executable"]["path"])
    require(sha(control.read_bytes()) == report["executable"]["sha256"], "Controller binary changed")
    expected = " ".join(map(str, [control, "workflow", "worker", "_serve", run_id, "--workspace", root, "--format", "json"]))
    require(isinstance(pid, int) and pid > 1 and pid != os.getpid(), "Invalid fixture controller PID")
    command = commands.run(["/bin/ps", "-ww", "-p", pid, "-o", "command="])
    identity = commands.run(["/bin/ps", "-p", pid, "-o", "uid=", "-o", "lstart="])
    require(command == expected and identity.split()[0] == str(os.getuid()), "PID does not identify the exact fixture worker controller")
    job = worker["job"]
    attempt = inside(root, job["directory"])
    process = read_json(attempt / "process.json")
    require(process["state"] == "running" and not (attempt / "result.json").exists(), "Provider attempt is not actively running")
    current = read_json(directory / "state.json")["result"]["worker"]
    require(current["owner_pid"] == pid and current["generation"] == worker["generation"] and current["job"] == job,
            "Native controller identity changed before interruption")
    require(commands.run(["/bin/ps", "-ww", "-p", pid, "-o", "command="]) == command
            and commands.run(["/bin/ps", "-p", pid, "-o", "uid=", "-o", "lstart="]) == identity,
            "Process identity changed before interruption")
    record = dict(run_id=run_id, controller_pid=pid, controller_command=command, process_identity=identity,
                  generation=worker["generation"], job_id=job["id"], intent_sha256=job["intent_sha256"],
                  attempt_directory=str(attempt), intended_signal="SIGKILL", requested_at=now())
    report["interruptions"].append(record)
    write_json(root / "journey.json", report)
    os.kill(pid, signal.SIGKILL)
    record["signal_sent_at"] = now()
    return record


def observe(commands, report, seconds, exercise_controls=False, interrupt=False):
    week_id, root = report["week_id"], commands.root
    until = time.monotonic() + seconds
    control_step = 0
    last = None
    while True:
        state = commands.ply("workflow", "week", "show", week_id, "--workspace", root)
        report["week_status"] = state
        summary = (state["phase"], state.get("active_task_id"), tuple((e["task_id"], e["state"]) for e in state["entries"]))
        if summary != last:
            print(f"{now()} {summary}", file=sys.stderr, flush=True)
            last = summary
        if report["service_requested"] and state.get("active_task_id") and not any(s["state"] == "running" for s in report["service_observations"]):
            report["service_observations"].append(commands.ply("workflow", "week", "service", "status", week_id, "--workspace", root))
        if interrupt and not report["interruptions"]:
            task_id = TASKS[1] if exercise_controls else TASKS[0]
            active = next((e for e in state["entries"] if e["task_id"] == task_id == state.get("active_task_id") and e.get("run_id")), None)
            if active:
                run = commands.ply("workflow", "worker", "show", active["run_id"], "--workspace", root)
                worker = run["worker"]
                if worker.get("job") and worker["job"]["kind"] == "developer" and worker.get("owner_state") == "lock_held":
                    job = inside(root, worker["job"]["directory"])
                    if (job / "process.json").is_file() and not (job / "result.json").exists():
                        interrupt_controller(commands, report, run)
        if exercise_controls and control_step < 4:
            active = next((e for e in state["entries"] if e["task_id"] == state.get("active_task_id") and e.get("run_id")), None)
            if control_step in (0, 2) and active:
                run = commands.ply("workflow", "worker", "show", active["run_id"], "--workspace", root)
                wanted_stage = "developer" if control_step == 0 else "reviewer"
                if run["worker"].get("job") and run["worker"]["job"]["kind"] == wanted_stage and run["worker"]["phase"] == "running" and run["worker"].get("owner_state") == "lock_held":
                    job = inside(root, run["worker"]["job"]["directory"])
                    if (job / "process.json").is_file() and not (job / "result.json").exists():
                        process = read_json(job / "process.json")
                        if process["state"] == "running":
                            action, requested_at = ("pause" if control_step == 0 else "stop"), now()
                            if action == "stop" and report["service_requested"]:
                                value = commands.ply("workflow", "week", "service", "stop", week_id, "--workspace", root)
                            else:
                                value = commands.ply("workflow", "week", action, week_id, "--workspace", root)
                            report["controls"].append(dict(action=action, requested_at=requested_at, run_id=active["run_id"],
                                active_job=run["worker"]["job"], active_process=process, result=value,
                                via_service=action == "stop" and report["service_requested"]))
                            control_step += 1
            elif control_step in (1, 3) and state["phase"] == ("paused" if control_step == 1 else "stopped"):
                previous = report["controls"][-1]
                require(any(e.get("run_id") == previous["run_id"] for e in state["entries"]), "Control lost its native Run")
                native = commands.ply("workflow", "worker", "show", previous["run_id"], "--workspace", root)
                require(native["worker"]["phase"] == state["phase"], "Week parked before native worker reached its safe boundary")
                completed = read_json(inside(root, previous["active_job"]["directory"]) / "result.json")
                require(completed["process"]["quiescence"] is True and completed["intent_sha256"] == previous["active_job"]["intent_sha256"], "Active attempt was not completed and preserved before parking")
                requested = datetime.datetime.fromisoformat(previous["requested_at"].replace("Z", "+00:00"))
                finished = datetime.datetime.fromisoformat(completed["finished_at"].replace("Z", "+00:00"))
                require(finished >= requested, "Control missed its active provider attempt; probe remains unproven")
                previous["observed_boundary"] = state
                previous["native_boundary"] = native["worker"]
                previous["completed_attempt"] = completed
                if previous["via_service"]:
                    previous["service_stopped"] = commands.ply("workflow", "week", "service", "status", week_id, "--workspace", root)
                value = commands.ply("workflow", "week", "resume", week_id, "--workspace", root)
                report["controls"].append(dict(action="resume", requested_at=now(), run_id=previous["run_id"], result=value))
                if previous["via_service"]:
                    report["controls"][-1]["service_after_resume"] = commands.ply("workflow", "week", "service", "status", week_id, "--workspace", root)
                control_step += 1
        write_json(root / "journey.json", report)
        if state["phase"] == "completed":
            require(not exercise_controls or control_step == 4, "Journey completed before both active control probes; controls remain unproven")
            final_workers = [commands.ply("workflow", "worker", "show", e["run_id"], "--workspace", root) for e in state["entries"]]
            if all(r["worker"]["phase"] == "completed" for r in final_workers):
                return inspect(commands, report, state)
            report["waiting_for_native_completion"] = [r["run_id"] for r in final_workers if r["worker"]["phase"] != "completed"]
        elif state["phase"] in TERMINAL and not (exercise_controls and control_step in (1, 3, 2, 4) and state["phase"] in ("paused", "stopped")):
            report.update(state="incomplete", observed_at=now())
            return report
        if time.monotonic() >= until:
            report.update(state="observing", observed_at=now(), wait_expired=True)
            return report
        time.sleep(min(3, max(0.01, until - time.monotonic())))


def self_test():
    # Exercise the independent oracle with a tiny correct executable and a real
    # broken variant. This is harness-only evidence, never provider evidence.
    with tempfile.TemporaryDirectory(prefix="ply-aw01-harness-smoke-") as tmp:
        program = Path(tmp) / "normalize"
        program.write_text('''#!/usr/bin/env python3
import json, sys
try:
 if sys.argv[1:] not in ([], ['--descending']): raise ValueError('arguments')
 data=json.loads(sys.stdin.read())
 if not isinstance(data,list) or any(type(x) is not int for x in data): raise ValueError('integer array required')
 print(json.dumps(sorted(set(data),reverse=bool(sys.argv[1:]))))
except Exception as error:
 print(str(error),file=sys.stderr); sys.exit(2)
''')
        program.chmod(0o755)
        result = verify_behavior(program, True)
        program.write_text(program.read_text().replace("reverse=bool(sys.argv[1:])", "reverse=False"))
        try:
            verify_behavior(program, True)
        except RuntimeError:
            pass
        else:
            raise RuntimeError("Harness failed to reject the descending regression")
        assets = Path(__file__).resolve().parents[3] / "internal/workerskills"
        manifest_bytes = (assets / "manifest.json").read_bytes()
        manifest = json.loads(manifest_bytes)
        for role, selection in manifest["roles"].items():
            prompt = f"# Pinned Ply worker role package\n\nRole: {role}\nSource commit: {ADDY_COMMIT}\nManifest SHA-256: {ADDY_MANIFEST}\n"
            prompt += "## Embedded provenance manifest\n\n```json\n" + manifest_bytes.decode() + "```\n\n## Upstream material\n"
            for name in selection["adapters"] + selection["files"]:
                prompt += f"--- BEGIN EMBEDDED FILE: {name} ---\n" + (assets / name).read_text() + f"\n--- END EMBEDDED FILE: {name} ---\n"
            verify_prompt(prompt.encode(), role)
            corrupted = prompt.replace("Permission is hereby granted, free of charge", "Permission was altered", 1)
            try:
                verify_prompt(corrupted.encode(), role)
            except RuntimeError:
                pass
            else:
                raise RuntimeError("Prompt verifier accepted corrupted embedded license bytes")
        print(json.dumps(dict(kind="HarnessOnlySmoke", correct_candidate=result,
                             known_regression_rejected=True, source_prompt_parser_roles=4,
                             altered_prompt_license_rejected=True, provider_calls=0), indent=2))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--ply", type=Path)
    parser.add_argument("--start", action="store_true", help="explicitly authorize this fixture's real provider calls")
    parser.add_argument("--service", action="store_true", help="with --start, install and start this fixture's named launchd service")
    parser.add_argument("--exercise-controls", action="store_true", help="with --start and --wait, request pause and stop during active attempts, then resume")
    parser.add_argument("--interrupt-controller", action="store_true", help="with --start and --wait, SIGKILL the exactly identified fixture worker controller during development; preserve the provider attempt")
    parser.add_argument("--inspect", type=Path, help="observe an existing disposable journey; makes no start/control call")
    parser.add_argument("--wait", type=int, default=0, metavar="SECONDS", help="observe to completion or timeout; timeout never means pass")
    parser.add_argument("--duration", type=int, default=7200, help="native contract duration from actual start (1..604800 seconds)")
    parser.add_argument("--expected-source-commit")
    parser.add_argument("--expected-sha256")
    parser.add_argument("--self-test", action="store_true", help="offline oracle smoke only; no native workspace or provider")
    args = parser.parse_args()
    if args.self_test:
        require(not args.start and not args.service and not args.inspect and not args.interrupt_controller, "Self-test cannot start a journey")
        self_test()
        return 0
    require(args.ply is not None, "--ply is required")
    require(0 <= args.wait <= 604800 and 1 <= args.duration <= 604800, "Invalid wait/duration")
    require(not args.service or args.start, "--service requires explicit --start")
    require(not args.exercise_controls or (args.start and args.wait > 0), "Control probes require explicit --start and a positive --wait")
    require(not args.interrupt_controller or (args.start and args.wait > 0), "Controller interruption requires explicit --start and a positive --wait")
    require(not args.inspect or not (args.start or args.service or args.exercise_controls or args.interrupt_controller), "--inspect is read-only")
    if args.inspect:
        root = args.inspect.resolve(strict=True)
        report = read_json(root / "journey.json")
        require(root.name.startswith("ply-aw01-week-") and report["kind"] == "PlyAW01WeekJourney@1"
                and report["root"] == str(root), "Not an AW-01 disposable fixture")
        commands = Commands(args.ply.resolve(strict=True), root)
        current = executable_identity(commands, args.expected_source_commit, args.expected_sha256)
        require(current["sha256"] == report["executable"]["sha256"], "Inspect using the exact frozen Ply executable")
    else:
        commands, report = prepare(args)
    try:
        if report.get("week_id"):
            report = observe(commands, report, args.wait, args.exercise_controls, args.interrupt_controller)
            if report.get("service_requested"):
                report["service_status"] = commands.ply("workflow", "week", "service", "status", report["week_id"], "--workspace", commands.root)
        write_json(commands.root / "journey.json", report)
        print(json.dumps(report, indent=2))
        return 0 if report["state"] in ("prepared", "started", "passed") or (report["state"] == "observing" and args.wait == 0) else 1
    except Exception as error:
        report.update(state="failed", observed_at=now())
        report["errors"].append(str(error))
        write_json(commands.root / "journey.json", report)
        raise


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as error:
        print(f"AW-01 journey failed: {error}", file=sys.stderr)
        raise SystemExit(1)
