# Prepared Task runs in Herdr

`ply workflow run` starts one interactive Codex session for an existing, prepared
Ply Task. It binds the frozen Task Spec, native handoff, runtime, Herdr identity,
reports and coordinator decisions. Its result is a reviewed report. It does not
publish TaskResult, attest provider inactivity, install a candidate, perform human
QA or integrate changes.

Run the coordinator commands inside the bound Ply workspace. Run recipient
callbacks from the exact Task worktree using the hash-bound Ply executable.

```sh
ply workflow run start --file /absolute/prepared-task-request.json --check
ply workflow run start --file /absolute/prepared-task-request.json --apply --confirm sha256:<preview-digest>
ply workflow run show wfr_<digest> --format json
ply workflow run follow wfr_<digest> --timeout 60
ply workflow run review wfr_<digest> --file /absolute/review.json
```

All six leaves support `--format text|json`; text is the default. Preview and
`show` do not write, create locks or contact Herdr. `show` labels saved transport
observations as cached and checks artifact freshness even after acceptance.
`follow` makes fresh observations; its timeout stops observation, not the agent.

The request is strict JSON with kind `ply.workflow.herdr-run-request` and
`schema_version: 1`. Supply these required fields:

| Field | Value |
| --- | --- |
| `request_key` | Stable native request key; different bytes under the same key conflict |
| `workspace_root` | Absolute physical containing workspace |
| `preparation_id`, `preparation_sha256` | Exact current native Task preparation |
| `handoff_draft` | Native Handoff@2 bound to that preparation and frozen Task Spec, four total rounds |
| `runtime` | Native TaskRun Runtime: interactive Codex, requested model, nullable config profile, physical hash-bound Codex/Ply executables and full permission evidence |
| `agreement` | Native A: 3 corrections, 5400 active seconds, 2 environment measures |
| `human_authority` | Nonempty `actor_claim`, `start_surface: "human_authorized_herdr"`, `authorized: true` |
| `herdr` | Physical hash-bound `executable`, `workspace_id`, and `tab_label` of 1–80 characters |
| `coordinator` | Named `actor_claim` and `may_request_changes: true` |
| `return_mode` | `"reviewed_report_only"` |

Digests use `sha256:` plus 64 lowercase hexadecimal characters. Apply requires
`HERDR_ENV=1`, the exact fresh preview confirmation, and local `codex` resolving
to the requested executable. The new tab receives the caller's validated PATH;
the recipient still reports its actual runtime before making Task writes. Starts
retain the requested model, optional config profile, managed permission profile,
`on-request`, and `approvals_reviewer="auto_review"`. No global config is changed.

The optional top-level `codex_project_trust` field has exactly this shape:

```json
{"mode":"process-local","repository_root":"/absolute/physical/repository"}
```

Omit it to preserve existing request bytes, digests and launch behavior. A supplied
null, another mode or extra field is rejected. For a linked worktree, the trust
root is the main repository owning the shared `.git` directory, which may differ
from the Task worktree and its cwd. Preview verifies both Git links and the physical
identities before offering a confirmation. The same confirmation covers this one
root and the process lifetime; apply rechecks the repository, executable and policy
bindings before effects. Changing trust under an existing request key conflicts.

The grant uses Codex CLI configuration before startup and makes no persistent trust
change. It preserves the requested model, configuration and permission profiles,
sandbox, `on-request` and `auto_review`. Explicit distrust, incompatible managed
policy, unknown search roots, symlinks or project config/hooks/rules/instructions
that could become active prevent a new grant before any tab is created. Inspect
the reported path and effective policy before requesting a new preview.

Ply reserves the target before creating a background tab, binds its pane and
terminal, starts Codex once, then binds the actual native session before prompting.
Agent start and its fresh readiness observations share one 35-second window. If
an ungranted legacy request encounters native onboarding, the attempt remains
subject to the same deadline. A confirmed process-local grant handles project trust
before launch. Ply only
observes the exact reserved workspace/tab/pane/terminal and named Codex agent;
it never answers trust dialogs. An early nonzero start reply can be followed by
these same-attempt observations within the remaining time. Missing `launch_pending`
is not proof that the provider exited. Only a full fresh,
settled native session with no pending launch permits the first Task prompt.

If readiness remains blocked, disappears, changes identity or reaches the deadline,
the coordinator inspects the preserved attempt and last observation. An unbound
recipient is never asked to accept. Safe diagnostics retain the phase, exit or
timeout class, stream sizes and known provider error code; unknown details stay
unknown. Raw child output, terminal contents and capabilities are not published.
There is no automatic restart or later readiness recovery after this window.

Every possibly submitted prompt has a durable attempt. A lost reply, timeout or
crash never permits automatic replay. Repeat `start` to read the same reservation,
or use `show` and `follow` to inspect it. Never delete state to manufacture a new
start. Native and Herdr runs share target exclusion. Even an accepted report keeps
the target reserved until a later qualified native TaskResult binds its terminal.

The generated instructions name the private context and redacted native handoff.
Acceptance is the recipient's first action:

```sh
ply workflow run accept wfr_<digest> --context /absolute/current/context.json --file /absolute/acceptance.json
ply workflow run report wfr_<digest> --context /absolute/current/context.json --file /absolute/round-report.json
```

Acceptance uses native Acceptance@2 fields with kind
`ply.workflow.run-acceptance`, schema version 1. A known different model is rejected;
an unknown actual model may be null. Necessary runtime and permission facts must
be known for positive acceptance. Negative claims remain visible without an
invented start receipt.

A round report uses all native TaskRun Report fields with kind
`ply.workflow.round-report`, schema version 1, plus `round`, `control_id`, and
`previous_report_sha256`. Round zero has null control and predecessor. Include
the frozen Spec's managed `task-requirements` artifact, actual candidate binding,
verifier results and cumulative `budget_usage`. Write each round's evidence to
new paths in the native reply staging directory. Previously reported artifacts
must remain unchanged. Reporting is the last target write of that round.

After two fresh idle/done observations at most one second apart, `follow` makes
the report ready for review. A `finished` tab label expresses transport status.
It never means technical approval, provider inactivity or a human pass.

Review uses kind `ply.workflow.run-review`, schema version 1, and all fields:
`review_id`, `run_id`, `request_sha256`, `handoff_sha256`, `round`, `report_sha256`,
`reviewer`, `decision`, `findings`. Reviewer must equal the named coordinator.
Each finding has `observed`, `expected`, and `acceptance` text. `changes_requested`
and `blocked` require findings; `accepted` may use an empty array.

`changes_requested` reserves another correction before one prompt to the same
session, including when internal corrections were already reported. The next
report must use the new context, control/predecessor binding and cumulative floor.
Unknown/decreasing usage or any reached A limit prevents another correction.
`accepted` requires a complete report, native technical gate `passed`, unchanged
evidence and a clean candidate. It seals the native workflow terminal once.
`blocked` preserves the reported outcome and coordinator stop without fabricating
a terminal. Identical reviews are idempotent; conflicting decisions are rejected.

Exit codes are 0 for successful operations/readback, 2 for invalid input, 4 for
conflicts or rejected transitions, and 5 for awaited reports or uncertain effects.
Other local I/O errors use 1. A nonzero status does not undo an attempted effect;
inspect the run ID and preserved state. The coordinator owns result control and
preparation of the later human QA and integration gates.

The executable stand-in and native fixture exporter live in test scope. To run
the plan-owned acceptance harness, build `./cmd/ply` and
`go test -c ./internal/taskrun` into a private artifact directory, then supply a
JSON argv file containing `python3`, the absolute
`test/fixtures/workflow_run/fixture.py` path, `--builder`, and the test executable.
Pass that file as `--fixture-command-file` to `wh01_acceptance.py`, along with
`--ply` and a private `--artifact-root`. Case directories are retained. These
tests use subprocess stand-ins and do not start real Herdr/Codex agents.
