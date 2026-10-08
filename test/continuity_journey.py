#!/usr/bin/env python3
"""Operate the disposable installed-product journey without copying internal hashes."""

import argparse
from datetime import datetime, timedelta, timezone
import fcntl
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import time

sys.dont_write_bytecode = True


def digest(path):
    return "sha256:" + hashlib.sha256(Path(path).read_bytes()).hexdigest()


class NativeCommandError(RuntimeError):
    def __init__(self, record):
        super().__init__(record["stderr"].strip() or f"native command exited {record['exit']}")
        try:
            self.payload = json.loads(record["stdout"])
        except ValueError:
            self.payload = {}


class Journey:
    def __init__(self, manifest_path):
        self.path = manifest_path.resolve(strict=True)
        self.root = self.path.parent
        self.manifest = json.loads(self.path.read_text())
        m = self.manifest
        if m.get("kind") != "PlyContinuityJourney@1" or m.get("schema_version") != 1:
            raise RuntimeError("this is not a disposable continuity journey")
        if Path(m["workspace"]).resolve(strict=True) != self.root:
            raise RuntimeError("the fixture workspace has moved; preserved bindings cannot be rewritten")
        for key in ("task_worktree", "target_worktree", "old_control", "context", "review", "counter"):
            if not Path(m[key]).resolve(strict=True).is_relative_to(self.root):
                raise RuntimeError(f"{key} escapes the isolated fixture")
        if digest(m["old_control"]) != m["old_control_sha256"] or digest(m["context"]) != m["context_sha256"]:
            raise RuntimeError("the preserved old control or context has changed")
        for path, expected in m["unrelated_publications"].items():
            if digest(path) != expected:
                raise RuntimeError("unrelated publication bytes have changed")
        self.upgrade = m.get("provider_upgrade")
        if self.upgrade:
            for key in ("historical_provider", "installed_provider", "before_control"):
                binding = self.upgrade[key]
                if not Path(binding["path"]).resolve(strict=False).is_relative_to(self.root):
                    raise RuntimeError(f"{key} escapes the isolated fixture")
                if key != "historical_provider" and digest(binding["path"]) != binding["sha256"]:
                    raise RuntimeError(f"{key} changed after preparation")
            for binding in self.upgrade["failure_evidence"]:
                path = Path(binding["locator"]).resolve(strict=True)
                if not path.is_relative_to(self.root) or digest(path) != binding["sha256"]:
                    raise RuntimeError("preserved pre-upgrade failure evidence changed")
            report = Path(self.upgrade["incomplete_report"]).resolve(strict=True)
            if not report.is_relative_to(self.root):
                raise RuntimeError("incomplete report escapes the isolated fixture")
            observation = self.upgrade["runtime_observation"]
            path = Path(observation["locator"]).resolve(strict=True)
            if not path.is_relative_to(self.root) or digest(path) != observation["sha256"]:
                raise RuntimeError("synthetic runtime observation template changed")
        self.binary = str(Path(m["binary"]).resolve(strict=True))
        self.cwd, self.run_id, self.context = m["task_worktree"], m["run_id"], m["context"]
        self.logs = self.root / "journey-observations"
        self.logs.mkdir(mode=0o700, exist_ok=True)

    def call(self, binary, *args, json_output=True, show_failure=True):
        argv = [binary, *args]
        if json_output:
            argv += ["--format", "json"]
        env = dict(os.environ, HERDR_ENV="1",
                   PATH=str(self.root / "bin") + os.pathsep + os.environ.get("PATH", ""))
        result = subprocess.run(argv, cwd=self.cwd, env=env, capture_output=True, text=True)
        record = {"argv": argv, "cwd": self.cwd, "exit": result.returncode,
                  "stdout": result.stdout, "stderr": result.stderr}
        path = self.logs / f"{time.time_ns()}.json"
        path.write_text(json.dumps(record, indent=2) + "\n")
        if result.returncode:
            if result.stdout and show_failure:
                print(result.stdout)
            raise NativeCommandError(record)
        if json_output:
            return json.loads(result.stdout)
        print(result.stdout, end="")
        return None

    def show(self):
        return self.call(self.binary, "workflow", "execute", "show", self.run_id)

    def control(self):
        path = self.root / "continued.json"
        if not path.is_file():
            raise RuntimeError("run 'journey continue' first")
        continued = json.loads(path.read_text())
        control = continued["preview"]["control_executable"]
        if not Path(control["path"]).resolve(strict=True).is_relative_to(self.root):
            raise RuntimeError("continuation control escapes the isolated fixture")
        if digest(control["path"]) != control["sha256"]:
            raise RuntimeError("preserved continuation control changed")
        return control["path"]

    def inspect(self):
        print(self.manifest["fixture_notice"])
        self.call(self.binary, "workflow", "execute", "show", self.run_id, json_output=False)
        run = self.show()
        native_session = run["transport"]["agent_session_id"]
        if native_session != self.manifest["native_session_id"]:
            raise RuntimeError("the accepted native session changed")
        print(f"Native session: {native_session} (unchanged)")
        if self.upgrade:
            historical = self.upgrade["historical_provider"]
            current = self.upgrade["installed_provider"]
            original = Path(historical["path"])
            observed = "missing" if not original.exists() else (
                "unchanged" if digest(original) == historical["sha256"] else "different bytes")
            print(f"Historical synthetic provider: {observed} at {original}")
            print(f"Current synthetic installation: {current['path']}")
            print("The new installation is not evidence about the accepted session's provider or authority.")
            for binding in self.upgrade["failure_evidence"]:
                receipt = json.loads(Path(binding["locator"]).read_text())
                operation = receipt["argv"][3]
                print(f"Preserved pre-fix {operation}: actual exit {receipt['exit']} (no report or continuation recorded)")
        count = Path(self.manifest["counter"]).read_text().count("acceptance executed\n")
        calls = [json.loads(line)["argv"][:2] for line in
                 Path(self.manifest["provider_calls"]).read_text().splitlines()]
        print(f"Acceptance executions: {count} (two original candidate checks; reuse adds zero)")
        print(f"Provider starts: {calls.count(['agent', 'start'])}; "
              f"Task tabs created: {calls.count(['tab', 'create'])}")
        print("Original control, original context and unrelated publication hashes: unchanged")

    def report_incomplete(self):
        if not self.upgrade:
            raise RuntimeError("this older fixture has no provider cleanup report")
        path = self.upgrade["incomplete_report"]
        report = json.loads(Path(path).read_text())
        run = self.call(self.binary, "workflow", "execute", "report", self.run_id,
                        "--context", self.context, "--file", path, "--incomplete")
        events = [event for event in run["delivery"]["events"] if event["id"] == report["event_id"]]
        if len(events) != 1:
            raise RuntimeError("the bounded incomplete report was not durably recorded once")
        print("The original owner's incomplete outcome is durably recorded once.")
        print("This report supplies no verification, human judgment, control replacement or delivery authority.")

    def continue_delivery(self):
        # Serialize this wrapper's evidence selection as well as the native
        # command. Native run locking independently protects the actual effect.
        with (self.root / "journey-continuation.lock").open("a+b") as lock:
            os.chmod(lock.name, 0o600)
            fcntl.flock(lock.fileno(), fcntl.LOCK_EX)
            self._continue_delivery()

    def _save_private_json(self, path, value):
        temporary = path.with_name(path.name + f".{time.time_ns()}.pending")
        with temporary.open("x") as stream:
            os.fchmod(stream.fileno(), 0o600)
            json.dump(value, stream, indent=2)
            stream.write("\n")
            stream.flush()
            os.fsync(stream.fileno())
        os.replace(temporary, path)
        directory = os.open(path.parent, os.O_RDONLY)
        try:
            os.fsync(directory)
        finally:
            os.close(directory)

    def _runtime_selection(self):
        path = self.root / "continuation-runtime.json"
        if not path.exists():
            return None
        binding = json.loads(path.read_text())
        observation = Path(binding["locator"]).resolve(strict=True)
        if not observation.is_relative_to(self.root) or digest(observation) != binding["sha256"]:
            raise RuntimeError("the selected immutable runtime observation changed")
        return binding

    def _select_current_runtime(self):
        # This generated claim is only for the labeled isolated fixture. Keep
        # every source file immutable, including observations from prior tries.
        observation = json.loads(Path(self.upgrade["runtime_observation"]["locator"]).read_text())
        observation["observed_at_utc"] = datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")
        path = self.logs / f"synthetic-runtime-{time.time_ns()}.json"
        self._save_private_json(path, observation)
        binding = {"locator": str(path), "sha256": digest(path)}
        self._save_private_json(self.root / "continuation-runtime.json", binding)
        print("Preserved an explicitly synthetic current session and authority observation.")
        return binding

    def _continuation_preview(self, binding):
        evidence = ["--runtime-evidence", binding["locator"]] if binding else []
        return self.call(self.binary, "workflow", "execute", "continue", self.run_id,
                         "--context", self.context, *evidence, "--check", show_failure=False)

    def _can_refresh_runtime(self, error, binding):
        if error.payload.get("diagnostic", {}).get("code") != "delivery_runtime_observation":
            return False
        observation = json.loads(Path(binding["locator"]).read_text())
        observed = datetime.fromisoformat(observation["observed_at_utc"].replace("Z", "+00:00"))
        if observed.tzinfo is None or observed >= datetime.now(timezone.utc) - timedelta(minutes=10):
            return False
        preview = error.payload.get("preview", {})
        if (preview.get("runtime_observation") != binding or preview.get("proof")
                or not preview.get("directory") or preview.get("generation", 0) < 1):
            return False
        directory = Path(preview["directory"])
        expected_parent = Path(self.context).parents[2] / "delivery" / "continuations"
        if directory.resolve(strict=False).parent != expected_parent.resolve(strict=False):
            raise RuntimeError("the native continuation directory is outside the original run")
        if preview.get("runtime_refresh_allowed") is not True:
            raise RuntimeError("The native preview has not confirmed that this expired observation can be refreshed. "
                               "Keep its exact evidence and inspect the preserved transition; no new observation was generated.")
        # Native validation found no published proof and checked any retained
        # control/before-state as the same generation. Keep those bytes intact.
        return True

    def _continue_delivery(self):
        binding = self._runtime_selection()
        try:
            preview = self._continuation_preview(binding)
        except NativeCommandError as error:
            code = error.payload.get("diagnostic", {}).get("code")
            if self.upgrade and binding is None and code == "delivery_runtime_observation_required":
                directory = Path(self.context).parents[2] / "delivery" / "continuations"
                if directory.exists() and any(directory.iterdir()):
                    raise RuntimeError("A continuation artifact already exists without this journey's selected observation. "
                                       "Inspect and recover that original transition before generating another observation.")
            elif not (self.upgrade and binding and self._can_refresh_runtime(error, binding)):
                raise
            binding = self._select_current_runtime()
            preview = self._continuation_preview(binding)
        if preview.get("state") == "existing":
            result = preview
        elif preview.get("state") == "ready":
            evidence = ["--runtime-evidence", binding["locator"]] if binding else []
            result = self.call(self.binary, "workflow", "execute", "continue", self.run_id,
                               "--context", self.context, *evidence)
        else:
            raise RuntimeError("read-only continuation observation did not permit an existing or ready transition")
        if result.get("state") not in ("continued", "existing"):
            raise RuntimeError("continuation did not reach its supported native state")
        self._save_private_json(self.root / "continued.json", result)
        print("Continuation:", result["state"])
        print("The same delivery is ready. Run 'journey verify' to reuse its successful check.")
        print("No provider restart, new Task, new terminal or renewed delivery permission.")

    def verify(self):
        before = Path(self.manifest["counter"]).read_bytes()
        run = self.call(self.control(), "workflow", "execute", "verify", self.run_id,
                        "--context", self.context, "--review", self.manifest["review"],
                        "--reuse", self.manifest["attempt"])
        if Path(self.manifest["counter"]).read_bytes() != before:
            raise RuntimeError("receipt reuse unexpectedly reran acceptance")
        candidates = run["delivery"]["candidates"]
        if len(candidates) != 2 or candidates[-1]["oid"] != self.manifest["corrected_candidate"]:
            raise RuntimeError("the exact correction did not qualify")
        print("The corrected candidate is technically qualified; acceptance was not rerun.")
        print("Human judgment and final delivery remain separate. Run 'journey inspect'.")

    def qa(self, outcome, answer):
        run = self.show()
        candidates = run["delivery"]["candidates"]
        if not candidates or candidates[-1]["oid"] != self.manifest["corrected_candidate"]:
            raise RuntimeError("qualify the correction with 'journey verify' before judging it")
        candidate = candidates[-1]
        # Native answer is the explicit pass/fail/blocked choice. Preserve any
        # additional exact feedback reversibly in observation; JSON quoting
        # also keeps multiline feedback within the native plain-text contract.
        observation = ("Explicit judgment of this disposable corrected candidate; "
                       "not acceptance of the parent implementation delivery.")
        if answer != outcome:
            observation += " Exact feedback: " + json.dumps(answer, ensure_ascii=False)
        if len(observation) > 2000:
            raise RuntimeError("feedback exceeds the native 2000-character observation limit; no QA was submitted")
        actor = self.manifest.get("qa_actor_claim", "Human operating the disposable product journey")
        path = self.root / f"human-answer-{outcome}.json"
        if path.exists():
            record = json.loads(path.read_text())
            if (record["answer"] != outcome or record["observation"] != observation
                    or record["actor_claim"] != actor or record["result_oid"] != candidate["oid"]):
                raise RuntimeError("this outcome already preserves a different answer; inspect its native record")
        else:
            now = datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")
            record = {"kind": "DeliveryHumanAttestation@1", "schema_version": 1,
                      "task_id": "task", "task_result_id": candidate["task_result"]["id"],
                      "result_oid": candidate["oid"], "result_tree": candidate["tree"],
                      "outcome": outcome, "actor_claim": actor,
                      "start_surface": "installed disposable native CLI journey",
                      "started_at_utc": now, "completed_at_utc": now, "answer": outcome,
                      "observation": observation}
            with path.open("x") as stream:
                json.dump(record, stream)
                stream.write("\n")
            path.chmod(0o600)
        self.call(self.control(), "workflow", "execute", "qa", self.run_id, outcome,
                  "--context", self.context, "--evidence", str(path))
        print(f"Preserved disposable-candidate judgment: {outcome}")
        if answer != outcome:
            print("Exact feedback:", answer)

    def integrate(self):
        run = self.call(self.control(), "workflow", "execute", "integrate", self.run_id,
                        "--context", self.context)
        if run["delivery"]["phase"] != "completed":
            raise RuntimeError("native integration did not report observed completion")
        print("The disposable candidate was locally integrated after its matching human pass.")
        print("Run 'journey inspect' to see verification, qualification, judgment and delivery separately.")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--fixture", type=Path,
                        default=Path(__file__).resolve().with_name("human-journey.json"))
    parser.add_argument("action", choices=["inspect", "report-incomplete", "continue", "verify", "qa", "integrate"])
    parser.add_argument("outcome", nargs="?", choices=["pass", "fail", "blocked"])
    parser.add_argument("--answer", help="additional exact feedback to preserve with the explicit outcome")
    args = parser.parse_args()
    if args.action == "qa" and args.outcome is None:
        parser.error("qa requires your explicit pass, fail or blocked judgment")
    if args.action != "qa" and (args.outcome is not None or args.answer is not None):
        parser.error("outcome and --answer apply only to qa")
    journey = Journey(args.fixture)
    if args.action == "qa":
        journey.qa(args.outcome, args.answer if args.answer is not None else args.outcome)
    elif args.action == "continue":
        journey.continue_delivery()
    else:
        getattr(journey, args.action.replace("-", "_"))()


if __name__ == "__main__":
    try:
        main()
    except (OSError, ValueError, KeyError, RuntimeError) as error:
        print("Journey stopped:", error, file=sys.stderr)
        raise SystemExit(1)
