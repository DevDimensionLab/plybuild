#!/usr/bin/env python3
"""Run all Go packages, partitioning only taskrun into isolated processes."""

import os
from pathlib import Path
import re
import signal
import subprocess
import sys
import time


def output(argv, cwd):
    return subprocess.check_output(argv, cwd=cwd, text=True).splitlines()


def main():
    repo, temporary = (Path(arg).resolve(strict=True) for arg in sys.argv[1:])
    package_dir = repo / "internal/taskrun"
    binary = str(temporary / "taskrun.test")
    # Exporters and subprocess helpers are opt-in entrypoints, not suite modes.
    for name in (
        "PLY_WORKFLOW_FIXTURE_ROOT", "PLY_CLAUDE_ENV_FIXTURE",
        "PLY_JOURNAL_FIXTURE_ROOT", "PLY_DELIVERY_JOURNEY_ROOT",
        "PLY_SYNTHETIC_PROCESS", "PLY_SYNTHETIC_RESERVATION_CHILD",
    ):
        if os.environ.get(name):
            raise RuntimeError(f"unset {name} before running acceptance")

    names = output([binary, "-test.list=^(Test|Example|Fuzz)"], package_dir)
    if not names or len(names) != len(set(names)) or any(
        not re.fullmatch(r"(?:Test|Example|Fuzz)\w*", name) for name in names
    ):
        raise RuntimeError("test discovery returned empty, duplicate or invalid names")
    names.sort()
    shards = [names[index::4] for index in range(4)]
    if sorted(name for shard in shards for name in shard) != names:
        raise RuntimeError("taskrun partitions do not cover every discovered test once")

    packages = output(["go", "list", "./..."], repo)
    taskrun = output(["go", "list", "./internal/taskrun"], repo)
    if len(taskrun) != 1 or packages.count(taskrun[0]) != 1:
        raise RuntimeError("cannot uniquely identify taskrun in the full package list")
    other_packages = [package for package in packages if package != taskrun[0]]
    jobs = [("other-packages", ["go", "test", *other_packages,
                                "-count=1", "-timeout=30m"], repo, None)]
    for index, shard in enumerate(shards, 1):
        if shard:
            expression = "^(" + "|".join(re.escape(name) for name in shard) + ")$"
            jobs.append((f"taskrun-{index}", [binary, f"-test.run={expression}",
                         "-test.count=1", "-test.timeout=30m", "-test.v",
                         "-test.paniconexit0"], package_dir, shard))
    print(f"Discovered {len(packages)} Go packages and {len(names)} taskrun tests; "
          "each test is assigned to exactly one process.", flush=True)

    running = []
    failed = False
    try:
        for label, argv, cwd, selected in jobs:
            log = temporary / f"{label}.log"
            with log.open("wb") as stream:
                process = subprocess.Popen(argv, cwd=cwd, stdout=stream,
                                           stderr=subprocess.STDOUT,
                                           start_new_session=True)
            running.append((label, process, log, selected))
            print(f"Started {label}: {len(selected) if selected is not None else len(other_packages)} "
                  f"{'tests' if selected is not None else 'packages'}", flush=True)
        while running:
            for job in list(running):
                label, process, log, selected = job
                code = process.poll()
                if code is None:
                    continue
                transcript = log.read_text(errors="replace")
                print(f"\n--- {label}: actual exit {code} ---\n{transcript}", flush=True)
                if selected is not None and code == 0:
                    observed = re.findall(r"^=== RUN   ([^\s/]+)$", transcript, re.MULTILINE)
                    if sorted(observed) != sorted(selected):
                        print(f"{label}: observed top-level tests differ from discovery", flush=True)
                        failed = True
                failed = failed or code != 0
                running.remove(job)
            if running:
                time.sleep(0.2)
    finally:
        # Stop only owned groups. Some synthetic providers create their own
        # groups, so interruption cannot attest to all descendant quiescence.
        # The shell preserves temporary evidence on any unsuccessful exit.
        if running:
            print("Acceptance interrupted: stopping owned worker groups; descendant "
                  "quiescence is unverified. Preserve temporary evidence.", flush=True)
        for _, process, _, _ in running:
            try:
                os.killpg(process.pid, signal.SIGTERM)
            except ProcessLookupError:
                pass
        for label, process, log, _ in running:
            try:
                code = process.wait(timeout=5)
            except subprocess.TimeoutExpired:
                try:
                    os.killpg(process.pid, signal.SIGKILL)
                except ProcessLookupError:
                    pass
                code = process.wait()
            transcript = log.read_text(errors="replace")
            print(f"\n--- {label}: interrupted, actual exit {code} ---\n{transcript}", flush=True)
    return 1 if failed else 0


if __name__ == "__main__":
    def interrupted(signum, frame):
        raise KeyboardInterrupt(f"acceptance interrupted by signal {signum}")

    signal.signal(signal.SIGTERM, interrupted)
    signal.signal(signal.SIGHUP, interrupted)
    sys.exit(main())
