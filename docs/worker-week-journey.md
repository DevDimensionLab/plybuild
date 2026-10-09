# Disposable real multi-Spec journey

`test/fixtures/headless_worker/week_journey.py` creates a new workspace under
`/private/tmp/ply-aw01-week-*`. It registers two real native Goals, a local test
Epic, the queue, and an exact week contract. The second Goal depends on completion
of the first. Neither Goal can touch an existing project, remote, installation,
notification route, or provider configuration.

The first Goal implements a small Python integer-array normalizer. The second
adds `--descending` and argument validation while retaining all original behavior.
Both must preserve the ten-case DW-01 acceptance block byte for byte, add meaningful
tests, pass actual native acceptance and independent actual review, and integrate
automatically into the disposable Epic. The explicit executor remains
`gpt-6-astra` with `max` effort and the existing Codex account login.

## Prepare without model calls

```sh
python3 -B test/fixtures/headless_worker/week_journey.py --self-test
python3 -B test/fixtures/headless_worker/week_journey.py --ply /absolute/qualified/ply
```

The first command checks only the harness oracle and prompt parser. It runs a
small correct executable, rejects a real descending regression, and rejects
altered embedded license bytes. These are harness tests, never worker/provider
evidence. On macOS, the executable checks deny writes and network access.

The second command creates a disposable fixture and uses public CLI commands to
register both Goals and preview `workflow week start --check`. It makes no model
call and installs no service. `prepared` means registration and preview succeeded;
it does not mean AW-01 acceptance passed.

## Start the real journey

Use the qualified installed executable for the installed-CLI acceptance journey:

```sh
python3 -B test/fixtures/headless_worker/week_journey.py \
  --ply /absolute/installed/ply --start --wait 5400 \
  --expected-source-commit COMMIT --expected-sha256 sha256:DIGEST
```

`--start` explicitly enables model calls. The harness invokes the week start once;
the week controller selects Spec 2 automatically. Default contract duration is
7200 seconds from actual start; `--duration` supports 1–604800 seconds. A bounded
`--wait` observes the independent controller and does not change its deadline.
Waiting timeout returns nonzero with `observing`, never a pass. Without a wait,
start returns while the controller continues.

Optional exercised boundaries:

```sh
python3 -B test/fixtures/headless_worker/week_journey.py \
  --ply /absolute/installed/ply --start --service \
  --exercise-controls --interrupt-controller --wait 5400
```

- `--service` installs and starts only the fixture's named launchd job,
  `com.ply.week.WEEK_ID`. The report preserves its label, plist path, frozen
  executable identity, observed runtime, and stop/uninstall commands.
- `--exercise-controls` requests pause during an actual developer process, waits
  for its preserved terminal result and native paused boundary, then resumes the
  same Run. It requests stop during an actual reviewer process and similarly
  verifies completion before resume. With `--service`, stop unloads the named
  fixture service and resume restarts it; final validation requires a newer
  scheduler generation. These are service/control tests, not a host reboot test.
- `--interrupt-controller` kills only the exactly identified fixture worker
  `_serve` process during development. It rechecks the full command, UID, start
  identity, native generation/job, and frozen binary immediately before SIGKILL.
  The provider and inherited-lock attempt runner are untouched. With both options,
  the interruption targets Spec 2; otherwise it targets Spec 1. Validation requires
  completion of the same preserved attempt and a later controller generation.

These options require explicit `--start` and a positive `--wait`. Missing the
intended active stage leaves that probe unproven and fails qualification. The
harness neither manufactures a failed provider answer nor claims that incidental
corrections prove an intentional regression exercise; AW-05 needs its separate
failure/correction evidence.

## Observe and validate preserved evidence

```sh
python3 -B test/fixtures/headless_worker/week_journey.py \
  --ply /absolute/frozen/ply --inspect /private/tmp/ply-aw01-week-EXAMPLE --wait 300
```

`--inspect` makes no start, resume, or control call. It can write its own command
logs and report while observing native state. Use the same executable bytes as
the original invocation. A changed input binary is rejected, while the running
controller retains its separately frozen native executable.

Every command's argv, cwd, timestamps, exit code, stdout, and stderr is preserved
under `commands/`; the consolidated report is `journey.json`. Preparation errors
also leave the disposable directory and command logs for inspection. Nothing is
automatically deleted or uninstalled. Use the exact `control_argv` and
`service_argv` recorded for this fixture when stopping or removing its service.
Do not remove its files while a controller or worker is active.

Only `state: passed` means the full inspection passed. It checks:

- Two distinct completed native Runs, an empty queue, and exact Epic history:
  baseline → Spec 1 candidate → Spec 2 candidate, with a higher native base revision.
- Actual developer completions and distinct actual reviewer sessions, preserved
  JSONL completion events, hashes, successful review, acceptance and native
  integration/closeout evidence; no HumanQA or human integration record.
- Exact pinned Addy commit, manifest, adapters, license and selected role-file
  bytes in the actual hashed prompts, including correction roles when present.
- Model/config boundaries and actual runtime probes; unchanged executable bytes
  throughout the journey. `executable.path` identifies the invoked file; supplying
  an installed path is what makes this an installed-CLI exercise.
- Byte-preserved original acceptance tests, a built executable identical to its
  committed source, and independent final behavior checks (13 for Spec 1 and
  30 for Spec 2) under denied writes/network on macOS.
- Requested active controls and interruption, when enabled, with actual start,
  deadline, last progress, and observed duration. A short run is never described
  as week-proven.

Blocked, stopped, expired, missing, changed, or incomplete evidence is retained
and reported. A terminal phase alone cannot produce a pass. This journey does
not replace deterministic dependency/backoff or Git-effect recovery tests.
