# Prepared Task runs in Herdr

`ply workflow run` starts one interactive Claude or Codex session for an existing, prepared
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
| `runtime` | Interactive runtime: selected provider, requested model, nullable config profile, physical hash-bound agent/Ply executables and full permission evidence |
| `agreement` | Native A: 3 corrections, 5400 active seconds, 2 environment measures |
| `human_authority` | Nonempty `actor_claim`, `start_surface: "human_authorized_herdr"`, `authorized: true` |
| `herdr` | Physical hash-bound `executable`, `workspace_id`, and `tab_label` of 1–80 characters |
| `coordinator` | Named `actor_claim` and `may_request_changes: true` |
| `return_mode` | `"reviewed_report_only"` |

### Choose Claude or Codex

Set `runtime.provider` in the request file to `"claude"` or `"codex"`.
Omitting this field defaults to Codex before the request is hashed. An explicit
null, empty string or unknown provider is rejected. Existing explicit Codex
requests and saved runs keep their bindings. There is no provider fallback.
The request is the single source for the choice; no separate CLI flag overrides it.
`runtime.model` selects the model within that provider and is always required.
Prefer the full model identifier when known: a CLI alias can resolve to a different
identifier, and acceptance still rejects a known literal model mismatch.
Preview and run readback expose `provider` in JSON and `Agent` in text.

| Runtime field | Codex | Claude |
| --- | --- | --- |
| `provider` | `"codex"` or omitted | `"claude"` |
| `mode` | `"interactive"` | `"interactive"` |
| `model` | An installed Codex model name | A Claude model name or supported alias, such as `"sonnet"` |
| `executable` | Physical path and SHA-256 of the selected Codex binary | Physical path and SHA-256 of the selected Claude binary |
| `config_profile` | Null or a Codex config profile | Null |
| `permission_binding.profile_id` | Bound Codex managed permission profile | `"manual"`, Claude's native permission mode |

Both providers still require the full effective-policy evidence and hash-bound
Ply executable. A mode name alone is not permission evidence. Claude's `manual`
mode retains its normal approval prompts. Other Claude permission modes are not
part of this contract. The separate native terminal and factory starts remain
Codex-only.

Digests use `sha256:` plus 64 lowercase hexadecimal characters. Apply requires
`HERDR_ENV=1`, the exact fresh preview confirmation, and the selected local
`claude` or `codex` resolving to the bound executable. A missing or different
binary stops before reservation or tab creation. The new tab receives the caller's
PATH and exact physical Task worktree. The recipient reports its actual runtime
before making Task writes.

Codex retains the requested model, optional config profile, managed permission
profile, `on-request`, and `approvals_reviewer="auto_review"`. Claude receives
`--model MODEL --permission-mode manual` and inherits the tab's cwd; fresh Herdr
observations must also report that exact foreground cwd. Codex flags and trust
grants are never sent to Claude. Ply never answers native login or permission dialogs.
Persistent Claude folder trust requires the separate grant below.

### Optional Claude folder trust

The Claude-only top-level `claude_project_trust` field has this shape:

```json
{"mode":"configuration","worktree_root":"/absolute/physical/task-worktree"}
```

Omission preserves existing requests and runs. Null, unknown modes, a different
provider or a path other than the prepared physical Task worktree are rejected.
The read-only preview shows the exact project key and configuration location.
After reserving the single start, immediately before launching Claude, Ply merges
only `projects[project_key].hasTrustDialogAccepted=true`. The key is the physical
Task worktree in Unicode NFC, supported by Claude Code 2.1.285's cwd trust fallback.
No main-checkout or other-worktree entry is added. The trust persists and allows
that folder's native project settings/hooks to become active; it is separate from
tool permissions, model selection and runtime acceptance.

The usual location is `~/.claude.json`, or `CLAUDE_CONFIG_DIR/.claude.json` with an
explicit absolute override. An existing `.config.json` inside the Claude config
directory has native legacy priority. Empty/relative overrides and custom OAuth
configuration are skipped with a reason. The tab receives the current home and
any explicit supported config directory. Before writing trust, Ply sends one
reserved shell setup in the new pane, using Bash/zsh builtins to remove inherited
overrides and confirm the intended environment. A private nonce acknowledgement
binds completion; terminal echo alone is insufficient. Unsupported shells, readonly
variables or an unconfirmed setup leave the same tab preserved without a trust
write or provider start. Ply never retries uncertain shell input. Shell hooks that
change environment again after this acknowledgement remain outside this binding.

The configuration writer preserves other JSON values and file permissions, keeps
a private unique backup, and atomically replaces the file in its own directory.
An already-true entry is a file no-op. A missing file is left to native onboarding.
Invalid JSON, duplicate keys, unexpected structures and unsafe paths are not repaired.
The normal native `.lock` directory coordinates concurrent writers; Ply never
steals an existing lock. Native fallback/exit writers can ignore it, so complete
exclusion from another Claude process is not guaranteed.

Run readback reports `claude_project_trust.state` as `written`, `already`, `skipped`,
`failed` or `unknown`, with paths, available hashes and a reason. It never returns
configuration contents or claims native readiness. A failed or uncertain trust
write produces a warning and continues the same ordinary Claude start. Any remaining
native dialog must be inspected in that tab. There is no automatic restart, rollback
or trust retry; check, show, follow and repeated start never write trust. Later normal
config changes do not invalidate the preserved run or its callbacks.

New Claude goals launched through [`workflow execute`](workflow-execute.md) bind
this grant in their launch intent. Previously saved intents are never upgraded.

### Optional Codex process trust

The optional, Codex-only top-level `codex_project_trust` field has exactly this shape:

```json
{"mode":"process-local","repository_root":"/absolute/physical/repository"}
```

Omit it to preserve existing request bytes, digests and launch behavior. A supplied
null, another mode, a Claude provider or extra field is rejected. For a linked worktree, the trust
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
terminal, starts the selected agent once, then binds the actual native session before sending
the Task. Codex versions that defer their `SessionStart` hook until the first turn
need a short readiness exchange. When Herdr explicitly reports an interactive,
settled selected agent without a session, Ply sends one fixed message asking for `PLY_READY`
without tools or file changes. This message contains no Task instructions. The
attempt is recorded before input in `startup-bootstrap-attempt.json`; it is never
replayed. Ply then waits for the real hook session and a fresh settled observation.
It neither invents a session ID nor parses terminal text to obtain one.

Agent start, the optional readiness exchange and fresh observations share one
90-second window. If
an ungranted legacy request encounters native onboarding, the attempt remains
subject to the same deadline. A confirmed process-local grant handles project trust
before launch. Ply addresses only the exact reserved workspace/tab/pane/terminal
and named agent, checking both Herdr's provider and the native session's provider;
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
be known for positive acceptance. Before the first positive acceptance, Ply also
checks the fresh Herdr provider, session and location; the agent may be working
while its callback runs. Negative claims remain visible without an
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
tests use subprocess stand-ins and do not start real Herdr/Claude/Codex agents.
