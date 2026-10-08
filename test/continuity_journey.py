#!/usr/bin/env python3
"""Operate the disposable installed-product journey without copying internal hashes."""

import argparse
from datetime import datetime, timezone
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
        self.binary = str(Path(m["binary"]).resolve(strict=True))
        self.cwd, self.run_id, self.context = m["task_worktree"], m["run_id"], m["context"]
        self.logs = self.root / "journey-observations"
        self.logs.mkdir(mode=0o700, exist_ok=True)

    def call(self, binary, *args, json_output=True):
        argv = [binary, *args]
        if json_output:
            argv += ["--format", "json"]
        env = dict(os.environ, HERDR_ENV="1")
        result = subprocess.run(argv, cwd=self.cwd, env=env, capture_output=True, text=True)
        record = {"argv": argv, "cwd": self.cwd, "exit": result.returncode,
                  "stdout": result.stdout, "stderr": result.stderr}
        path = self.logs / f"{time.time_ns()}.json"
        path.write_text(json.dumps(record, indent=2) + "\n")
        if result.returncode:
            if result.stdout:
                print(result.stdout)
            raise RuntimeError(result.stderr.strip() or f"native command exited {result.returncode}")
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
        count = Path(self.manifest["counter"]).read_text().count("acceptance executed\n")
        calls = [json.loads(line)["argv"][:2] for line in
                 Path(self.manifest["provider_calls"]).read_text().splitlines()]
        print(f"Acceptance executions: {count} (two original candidate checks; reuse adds zero)")
        print(f"Provider starts: {calls.count(['agent', 'start'])}; "
              f"Task tabs created: {calls.count(['tab', 'create'])}")
        print("Original control, original context and unrelated publication hashes: unchanged")

    def continue_delivery(self):
        result = self.call(self.binary, "workflow", "execute", "continue", self.run_id,
                           "--context", self.context)
        if result.get("state") not in ("continued", "existing"):
            raise RuntimeError("continuation did not reach its supported native state")
        (self.root / "continued.json").write_text(json.dumps(result, indent=2) + "\n")
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
    parser.add_argument("action", choices=["inspect", "continue", "verify", "qa", "integrate"])
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
        getattr(journey, args.action)()


if __name__ == "__main__":
    try:
        main()
    except (OSError, ValueError, KeyError, RuntimeError) as error:
        print("Journey stopped:", error, file=sys.stderr)
        raise SystemExit(1)
