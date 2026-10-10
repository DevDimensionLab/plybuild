#!/bin/sh
# Draft Task-scoped driver. See resume-next-action-goal-draft.md before binding.
# Deliberately no heredoc: macOS /bin/sh may use /tmp despite private TMPDIR.
set -eu
exec python3 -I -B -c '
import hashlib
import json
import os
from pathlib import Path
import re
import signal
import stat
import subprocess
import sys
import tempfile
import time

EXPECTED_REF = "refs/heads/ply-task/aw-01-current-resume-guidance-038968865f8f"
PACKAGE = "github.com/devdimensionlab/plybuild/internal/taskrun"
GO = Path("/Users/perottochristensen/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.7.darwin-arm64/bin/go")
MODULE_CACHE = Path("/Users/perottochristensen/go/pkg/mod")
REQUIRED = (
    "TestAW01GuidanceLegacyRunningReadback",
    "TestAW01GuidanceResumeTransition",
    "TestAW01GuidanceWaitingResume",
    "TestAW01GuidancePauseStopAndBlocked",
    "TestAW01GuidancePendingControls",
    "TestAW01GuidanceCloseoutAndDriftPriority",
    "TestAW01GuidanceReadDoesNotWrite",
    "TestAW01GuidanceInteractiveUnchanged",
)

def require(ok, message):
    if not ok:
        raise RuntimeError(message)

def sha(path):
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for block in iter(lambda: stream.read(1048576), b""):
            digest.update(block)
    return "sha256:" + digest.hexdigest()

def physical_env(name):
    value = os.environ.get(name, "")
    path = Path(value)
    require(value and path.is_absolute() and path.resolve(strict=True) == path, name + " must be an existing physical absolute path")
    return path

def main():
    repo = Path.cwd().resolve(strict=True)
    require((repo / "go.mod").is_file() and (repo / "internal/taskrun").is_dir(), "run from the candidate repository root")
    private_home = physical_env("HOME")
    private_tmp = physical_env("TMPDIR")
    require(private_home.is_dir() and private_tmp.is_dir(), "native private HOME/TMPDIR are required")
    candidate = physical_env("PLY_CANDIDATE_BINARY")
    require(stat.S_ISREG(candidate.stat().st_mode) and os.access(candidate, os.X_OK), "native candidate must be a regular executable")
    oid, tree = (os.environ.get(name, "") for name in ("PLY_CANDIDATE_OID", "PLY_CANDIDATE_TREE"))
    require(all(re.fullmatch(r"(?:[0-9a-f]{40}|[0-9a-f]{64})", value) for value in (oid, tree)), "native candidate OID/tree are required")
    require(GO.is_file() and os.access(GO, os.X_OK) and MODULE_CACHE.is_dir(), "the recorded local Go toolchain/cache must already exist; no download fallback")
    env = {"HOME": str(private_home), "TMPDIR": str(private_tmp), "PATH": os.environ.get("PATH", "/usr/bin:/bin"), "LC_ALL": "C", "GIT_OPTIONAL_LOCKS": "0"}

    def git(*args):
        result = subprocess.run(["/usr/bin/git", "-c", "core.fsmonitor=false", "-c", "core.hooksPath=/dev/null", *args], cwd=repo, env=env, capture_output=True, text=True, timeout=10)
        require(result.returncode == 0, "read-only Git observation failed: " + " ".join(args) + ": " + result.stderr.strip())
        return result.stdout.strip()

    def source_binding():
        require(git("symbolic-ref", "--quiet", "HEAD") == EXPECTED_REF, "focused acceptance is bound only to the exact selected Task ref")
        require(git("rev-parse", "HEAD") == oid and git("rev-parse", "HEAD^{tree}") == tree, "native candidate OID/tree do not match this checkout")
        require(git("status", "--porcelain=v1", "--untracked-files=all") == "", "candidate source must remain clean")

    source_binding()
    candidate_before = sha(candidate)
    scratch = Path(tempfile.mkdtemp(prefix="aw01-guidance-", dir=private_tmp))
    os.chmod(scratch, 0o700)
    for name in ("tmp", "go-cache", "go-tmp"):
        (scratch / name).mkdir(mode=0o700)
    env.update({
        "TMPDIR": str(scratch / "tmp"),
        "GOCACHE": str(scratch / "go-cache"),
        "GOTMPDIR": str(scratch / "go-tmp"),
        "GOMODCACHE": str(MODULE_CACHE),
        "GOROOT": str(GO.parent.parent),
        "GOTOOLCHAIN": "local",
        "GOENV": "off",
        "GOWORK": "off",
        "GOPROXY": "off",
        "GOSUMDB": "off",
        "PLY_CANDIDATE_BINARY": str(candidate),
        "PLY_CANDIDATE_OID": oid,
        "PLY_CANDIDATE_TREE": tree,
    })
    identity = {
        "kind": "AW01CurrentGuidanceAcceptanceIdentity@1",
        "source_ref": EXPECTED_REF, "candidate_oid": oid, "candidate_tree": tree,
        "candidate_binary": {"path": str(candidate), "sha256": candidate_before, "coverage": "identity only; this driver does not execute the CLI"},
        "go": {"path": str(GO), "sha256": sha(GO)},
        "driver_sha256": sha(Path(sys.argv[1]).resolve(strict=True)),
        "required_tests": list(REQUIRED), "artifact_directory": str(scratch),
        "evidence_kind": "deterministic public-function regressions compiled from the bound candidate source",
    }
    (scratch / "identity.json").write_text(json.dumps(identity, indent=2) + "\n")
    command = [str(GO), "test", "-json", "-mod=readonly", "-count=1", "-p=1", "-parallel=1", "-timeout=120s", "-run", "^(" + "|".join(REQUIRED) + ")$", "./internal/taskrun"]
    started = time.monotonic()
    timed_out = False
    interrupted = {"signal": None}
    runner_error = None
    cleanup_errors = []
    child, code = None, None

    def request_interruption(number, frame):
        # Record rather than raise: even a signal during Popen must allow the
        # returned child identity to be captured before our finally cleanup.
        if interrupted["signal"] is None:
            interrupted["signal"] = number

    previous_handlers = {number: signal.signal(number, request_interruption) for number in (signal.SIGTERM, signal.SIGHUP, signal.SIGINT)}
    with (scratch / "go-test.jsonl").open("wb") as stdout, (scratch / "go-test.stderr").open("wb") as stderr:
        try:
            if interrupted["signal"] is None:
                child = subprocess.Popen(command, cwd=repo, env=env, stdout=stdout, stderr=stderr, start_new_session=True)
            while child is not None and interrupted["signal"] is None:
                remaining = 300 - (time.monotonic() - started)
                if remaining <= 0:
                    timed_out = True
                    break
                try:
                    code = child.wait(timeout=min(0.25, remaining))
                    break
                except subprocess.TimeoutExpired:
                    continue
        except BaseException as error:
            runner_error = type(error).__name__ + ": " + str(error)
        finally:
            if child is not None and (code is None or interrupted["signal"] is not None or timed_out or runner_error is not None):
                # Only this newly-created test process group, never a native
                # worker or an independently existing process.
                try:
                    os.killpg(child.pid, signal.SIGKILL)
                except ProcessLookupError:
                    pass
                except OSError as error:
                    cleanup_errors.append("cannot stop own test group: " + str(error))
                try:
                    code = child.wait(timeout=10)
                except BaseException as error:
                    cleanup_errors.append("own test child exit remains unresolved: " + type(error).__name__ + ": " + str(error))
            for number, handler in previous_handlers.items():
                signal.signal(number, handler)
    observed = {name: {"run": 0, "pass": 0} for name in REQUIRED}
    bad_events = []
    package_pass = 0
    for number, line in enumerate((scratch / "go-test.jsonl").read_text().splitlines(), 1):
        try:
            event = json.loads(line)
        except json.JSONDecodeError:
            bad_events.append({"line": number, "reason": "non-JSON go-test output"})
            continue
        if event.get("Package") != PACKAGE:
            continue
        action, name = event.get("Action"), event.get("Test")
        if action in ("fail", "skip"):
            bad_events.append({"test": name, "action": action})
        if name in observed and action in observed[name]:
            observed[name][action] += 1
        if name is None and action == "pass":
            package_pass += 1
    errors = []
    if code != 0 or timed_out:
        errors.append("the bounded Go test process did not finish successfully")
    if interrupted["signal"] is not None:
        errors.append("acceptance interrupted by signal " + str(interrupted["signal"]))
    if runner_error is not None:
        errors.append(runner_error)
    errors.extend(cleanup_errors)
    if bad_events or package_pass != 1:
        errors.append("a failed/skipped test, malformed output or missing package pass was observed")
    for name, counts in observed.items():
        if counts != {"run": 1, "pass": 1}:
            errors.append("mandatory regression did not run and pass exactly once: " + name)
    try:
        source_binding()
        require(sha(candidate) == candidate_before, "candidate executable changed during acceptance")
    except Exception as error:
        errors.append(str(error))
    summary = {
        "kind": "AW01CurrentGuidanceAcceptance@1", "outcome": "pass" if not errors else "fail",
        "argv": command, "exit_code": code, "timed_out": timed_out,
        "interrupted_signal": interrupted["signal"], "runner_error": runner_error, "cleanup_errors": cleanup_errors,
        "elapsed_wall_seconds": round(time.monotonic() - started, 3),
        "required_tests": observed, "bad_events": bad_events, "errors": errors,
        "identity": identity, "scope": "focused deterministic package cases; no fullsuite or live provider/service qualification",
    }
    (scratch / "summary.json").write_text(json.dumps(summary, indent=2) + "\n")
    print(json.dumps(summary, sort_keys=True))
    return 0 if not errors else 1

try:
    raise SystemExit(main())
except Exception as error:
    print(json.dumps({"kind": "AW01CurrentGuidanceAcceptance@1", "outcome": "fail", "error": str(error)}))
    raise SystemExit(1)
' "$0"
