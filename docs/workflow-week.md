# Run a selected queue without an open terminal

`ply workflow week` supervises registered native Goals through the existing
headless worker. Each Goal retains its native Task, preparation, Run, evidence,
review, acceptance and local Epic receipt. The week files only bind the selected
work and remember control/observation state. They are not a second Task database.

Each selected Goal must already be in the same native Epic queue and explicitly
authorize automatic acceptance and local Epic integration, without a human gate.
Register the Tasks and Goals first. Obtain each exact Goal revision and manifest
hash from `workspace task queue list --format json`.

From the receiving Epic worktree, create an absolute JSON file such as:

```json
{
  "kind": "ply.workflow.week-contract",
  "schema_version": 1,
  "name": "selected-maintenance",
  "workspace": "/absolute/workspace",
  "return_worktree": "/absolute/workspace/epic",
  "return_ref": "refs/heads/epic",
  "project_id": "example",
  "repo_id": "example",
  "epic_id": "maintenance",
  "selected": [
    {
      "task_id": "first-task",
      "goal": {
        "spec_id": "first-goal",
        "spec": {"revision": 1, "manifest_sha256": "sha256:<64 hex digits>"}
      },
      "depends_on": []
    },
    {
      "task_id": "second-task",
      "goal": {
        "spec_id": "second-goal",
        "spec": {"revision": 1, "manifest_sha256": "sha256:<64 hex digits>"}
      },
      "depends_on": ["first-task"]
    }
  ],
  "model": "gpt-6-astra",
  "effort": "max",
  "duration_seconds": 604800,
  "max_attempts": 3,
  "agent_timeout_seconds": 1800,
  "max_provider_retries": 20
}
```

The example hashes must be replaced with actual native Goal hashes. Model and
effort must match each Goal's immutable executor assignment; omitted native
choices retain the worker defaults shown above. No model switch, new payment
route, plugin installation or global Codex setting change is performed. The
existing account login is used.

To opt a new week into [request-bound runtime control](workflow-control.md), set
`schema_version` to `2` and add
`"runtime_control": {"schema_version": 1, "mode": "bound_control_v1"}`. All other
selection and execution bounds remain required. The check displays the frozen
control agreement. Version 1 retains its existing controls and does not gain
mobile actions when a new CLI is installed.

```sh
ply workflow week start --file /absolute/week.json --check --format json
ply workflow week start --file /absolute/week.json --service --format json
```

The check is read-only. Start records informative Ply identity and freezes the exact
selection and dependencies, the return target, the model, attempt bounds, actual
start time and deadline. Duration is 1–604800 seconds from this start, including
sleep and downtime. Repeating the same name/contract reads the same week;
changing its contract under that name is rejected. Empty selection completes
without starting work. Later queue additions never expand this selection.

Without `--service`, a detached controller survives closure of its launching
terminal. Detached recovery uses the same captured login environment as launchd;
a later shell cannot silently select a different HOME, CODEX_HOME or PATH. `--service` installs and starts the named per-user launchd service
immediately, resolving ordinary installed `ply` through the saved PATH. Neither mode publishes or installs a
worker's product candidate; native local Epic return is the authorized effect.

# Observe, pause and resume

Use the returned `week_<digest>` ID and physical workspace root:

```sh
ply workflow week show WEEK_ID --workspace /absolute/workspace --format json
ply workflow week pause WEEK_ID --workspace /absolute/workspace
ply workflow week stop WEEK_ID --workspace /absolute/workspace
ply workflow week resume WEEK_ID --workspace /absolute/workspace
```

Show reports actual start/deadline, last progress, active Task, native Run IDs,
base OIDs/revisions and current controller lock observation. A saved PID is not
proof that a controller lives. Files are under
`WORKSPACE/.ply/week-runs/WEEK_ID/`; `binding.json` preserves the immutable
contract/control, while `state.json` and `control.json` preserve observations and
control requests. Native `workflow worker show RUN_ID` gives detailed attempt
and reset evidence. If the service has unloaded, show reads the active worker
to report its confirmed pause/stop boundary without writing an observation.

An interrupted first registration can recover an immutable binding whose initial
state was not yet published, preserving its original deadline. An activity marker
is written before operational execution; missing state after that point is an
integrity error requiring inspection, never permission to recreate Task/Run work.

Pause and stop publish a request to the active native worker. An ongoing attempt
must finish with positive process and result evidence before the week reports
`paused` or `stopped`. `pausing` and `stopping` mean this boundary is still pending.
Both actions preserve the same Run, and resume cannot extend the original
deadline. Expiry requests stop and prevents any next Goal from starting. A
resume request after expiry wakes the preserved controller to enforce stop on
any surviving native Run; it grants no further work or time.

A completed native local return permits the next selected dependency to start
from the freshly observed Epic base. A concrete, positively completed developer
blocker or exhausted correction limit may release its preparation using the
native API, retaining its Task checkout, Run and evidence. This is irreversible
parking of that Task; its dependents are skipped, while independent selected
work continues. `completed_with_blocks` lists this preserved unfinished work.
No queue removal or Task recreation is used to disguise a failure.

Authentication errors and uncertain provider/Git effects block the controller
without releasing the Task. A recognized usage limit stays on the same active
Run in `waiting`, with its native persisted reset/backoff. A live native
supervisor is observed through its process lock; polling does not start another
supervisor. Explicit resume is necessary for a parked non-transient blocker.
The native worker still refuses any unsafe replay of an unknown attempt.

# Named login service

```sh
ply workflow week service install WEEK_ID --workspace /absolute/workspace
ply workflow week service status WEEK_ID --workspace /absolute/workspace --format json
ply workflow week service stop WEEK_ID --workspace /absolute/workspace
ply workflow week service start WEEK_ID --workspace /absolute/workspace
ply workflow week resume WEEK_ID --workspace /absolute/workspace
ply workflow week service uninstall WEEK_ID --workspace /absolute/workspace
```

The label is `com.ply.week.WEEK_ID`; its plist lives in the invoking account's
`~/Library/LaunchAgents`. It uses `gui/UID`, ordinary installed `ply` resolved through the saved PATH, the
registered return directory and the original HOME/PATH/TMPDIR plus an explicitly
set CODEX_HOME. It preserves existing authentication files. It needs the user's
available login environment; it is not a system-wide daemon.

Installation can take over an existing detached controller at a native API
boundary. The service waits for the same owner lock, while the detached
controller yields it between effects. The active native worker keeps its own
lock and Run. Launchd restarts a failed service from the preserved state; an
ordinary completed/paused/stopped/blocked exit stays parked. `service start`
loads/kickstarts supervision and respects a saved pause/stop; use `week resume`
to explicitly resume that work. `week resume` also loads an installed but stopped
service. No kickstart command uses `-k` to kill active work.

Service stop first persists the week stop and waits for the controller to
acknowledge that native worker control is published, then unloads launchd. The
native worker can still be finishing its current attempt. If acknowledgment is
not observed within 25 seconds, the service stays loaded and the CLI reports the
preserved pending request. Uninstall uses the same stop boundary and removes
only this exact managed plist; every wake and stop validates its frozen contents
before launchctl effects. Native history, product worktrees, frozen
controls and evidence remain. Service status reports whether the plist is
installed, loaded or running and its observed launchd PID.

A service restart is not a reboot test. Sleep pauses execution while the absolute
deadline continues. No reboot or week-long reliability claim follows from these
short tests.

# Verification

`go test ./internal/weekrun` uses explicit synthetic native adapters to exercise
selected dependencies, fresh-base succession, skipped dependents and independent
continuation, empty queues, lost replies before Run observation and after release,
persisted backoff, authenticated/unknown blockers, controls, time limits,
immutable bindings and process locks. The service tests redirect the binding's
HOME into a disposable directory and inject launchctl; they do not load a real
service or call a provider. On macOS, the generated plist is also checked by the
actual `plutil` parser.

A real provider/multi-Spec/service journey must be run separately in a disposable
native workspace. Inspect both developer and distinct reviewer sessions, native
acceptance/integration receipts and the second Goal's updated base before calling
that journey qualified. Preserve its actual start and observed elapsed duration.
