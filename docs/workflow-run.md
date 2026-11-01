# Prepared Task runs in Herdr

`ply workflow run` starts one interactive Claude or Codex session for an existing, prepared
Ply Task. It binds the frozen Task Spec, native handoff, runtime, Herdr identity,
reports and coordinator decisions. Its result is a reviewed report. It does not
publish TaskResult, attest provider inactivity, install a candidate, perform human
QA or integrate changes. A later commit on the Epic does not invalidate a report-only
start from an unchanged prepared Task source. The exact Task source, input and
ownership checks still apply; delivery retains its parent revision check.

Run the coordinator commands inside the bound Ply workspace. Run recipient
callbacks from the exact Task worktree using ordinary installed `ply` and the
original context/session. Historical Ply path/hash observations are informative.

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
Its exit-5 observation result includes the same Run's next `follow` command.
The overall observation window is distinct from the Herdr request budget:
expiry of the window preserves the round, while a failed read leaves current
transport unknown. An independent `agent get` timeout permits at most one read
continuation within the original window, with a sourced warning through the
existing logger. The continuation rechecks the same session and sends no provider
input. Preserved stops return before polling, and session drift remains an error.

The request is strict JSON with kind `ply.workflow.herdr-run-request` and
`schema_version: 1`. Supply these required fields:

| Field | Value |
| --- | --- |
| `request_key` | Stable native request key; different bytes under the same key conflict |
| `workspace_root` | Absolute physical containing workspace |
| `preparation_id`, `preparation_sha256` | Exact current native Task preparation |
| `handoff_draft` | Native Handoff@2 bound to that preparation and frozen Task Spec, four total rounds |
| `runtime` | Interactive runtime: selected provider, requested model, nullable config profile, physical hash-bound agent/Ply executables and the selected launch permission profile |
| `agreement` | Explicit request limits for corrections, active time and environment measures |
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
| `permission_binding.profile_id` | Bound Codex managed permission profile | `"manual"` or `"auto"`, Claude's native permission mode |

New receipt reports bind the selected launch directly. Use this permission
binding with the provider's profile from the table above (this example selects
Claude auto); no separate effective-policy document is required:

```json
{
  "authority_kind": "launch_contract_pending_runtime_acceptance",
  "profile_id": "auto",
  "effective_policy_sha256": "",
  "evidence": []
}
```

The hash-bound agent/Ply executables, Task Spec, human authority and Herdr
session still apply. Optional evidence keeps its typed, sorted locator/hash
bindings and is checked for changes. A nonempty optional policy hash must match
supplied evidence. These are launch inputs, not observed runtime policy.
Unknown model, profile, policy and sandbox facts remain unknown at check-in;
receipt grants no blanket permission. Frozen legacy positive-acceptance report
contracts still require their effective-policy evidence. The separate terminal
request remains available; factory and execute-owned Delivery starts have been removed.

Claude's `manual`
mode retains its normal approval prompts. In `auto` mode Claude's own runtime
still evaluates each action; an actual refusal stays preserved and a later
receipt, label or actor cannot clear it. Preview, launch argv and readback show
the same frozen selection. Any other Claude permission mode is rejected before
reservation or launch; select `manual` or `auto` in a new request. Existing
report Runs keep their frozen mode. The separate native terminal start remains Codex-only.

Digests use `sha256:` plus 64 lowercase hexadecimal characters. Apply requires
`HERDR_ENV=1`, the exact fresh preview confirmation, and the selected local
`claude` or `codex` resolving to the bound executable. A missing or different
binary stops before reservation or tab creation. The new tab receives the caller's
PATH and exact physical Task worktree. The recipient reports its actual runtime
before making Task writes.

Codex retains the explicitly requested model, optional config profile and permission
profile. Native approval policy and reviewer come from Codex configuration. Claude receives
`--model MODEL --permission-mode MODE` with the selected `manual` or `auto` and inherits the tab's cwd; fresh Herdr
observations must also report that exact foreground cwd. Ply does not grant project
trust, write provider configuration or answer native login or permission dialogs.

### Claude auto report pilot

For a separate actual-provider check, prepare a disposable native Task with a
read-only inspection goal and `reviewed_report_only` return mode. Bind the
current installed Ply and Claude executables and the launch binding above in a new
request; set provider `claude`, mode `interactive`, null `config_profile`, and
`permission_binding.profile_id: "auto"`. Use the preview and confirmed start
commands above from the coordinator's Herdr workspace. Keep the request and
confirmation unchanged when observing a lost reply so it reuses the same Run.

In the launched Claude session, use the generated context and its exact bound
Ply executable for `workflow run check-in RUN_ID --context CONTEXT --provider
claude --native-session SESSION`. Optional `--profile` records only an actually
observed runtime label. An unknown or different first observation is preserved
without changing the requested launch mode; an identical retry is idempotent,
and a later receipt cannot rewrite that observation. Do the read-only work,
submit the generated round report through `workflow run report`, and let the
coordinator `follow` and `review` it with the generated review format. If changes
are requested, correct and report in the same session; an accepted review seals
the report return without publishing a TaskResult or integrating changes.
Preserve any actual refusal or unknown effect and stop the dependent action.

Synthetic CLI regressions qualify the implementation, not an actual Claude
pilot. AL06B delivered the technical change; the separate revised AL06C2 pilot
subsequently completed actual Claude auto reads, two correction rounds and an
accepted report-only return. Its provider ended and preparation choice closed;
no TaskResult or full Task delivery/closeout was claimed. The first AL06C1 Run
remains stopped with source/history preserved, not accepted or retired. See the
[reconciled audit](approval-requirements-audit.md#p3-outcome-al06b-and-al06c2)
for the distinct native outcomes and evidence limits.

### Launch facts after launch

Once the report owner's native session is bound, its launch has used the provider
executable and any process-local trust grant. Later report operations (check-in,
accept, report, stop, show, follow, review) never launch a provider. An ordinary
Ply or provider install, a removed provider launcher, or a Codex config or
`.git/config` edit after launch then logs a `warn` on stderr (`fact`, `bound`,
`observed`, `handling`) and the operation continues. The request keeps the original
launch bindings; nothing is rewritten and no new authority is inferred.

These checks still apply after launch: the callback runs in the exact Task
worktree, with the exact context, run, request and native
session; frozen inputs, policy evidence, receipts, reports and the Herdr
executable keep their bytes; preserved stops, negative acceptance and retirement
stay in force. Before the session is bound, provider executable identity is checked.
The native runtime continues to enforce its own permissions.

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
subject to the same deadline. An applied process-local grant handles project trust
before launch. Ply addresses only the exact reserved workspace/tab/pane/terminal
and named agent, checking both Herdr's provider and the native session's provider;
it never answers trust dialogs. An early nonzero start reply can be followed by
these same-attempt observations within the remaining time. Missing `launch_pending`
is not proof that the provider exited. Only a full fresh,
settled native session with no pending launch permits the first Task prompt.

Owner startup, independent review and new-terminal owner resume use public
Herdr `agent start`, `agent get` and `agent prompt` operations. They require no
private startup capability or forced Codex status-line profile. Public readiness
observations do not atomically prove that the terminal still contains an empty
provider prompt when input arrives. A known update or security dialog is not
permission to send Task text, approve an upgrade or answer the dialog; preserve
the attempt and let the coordinator inspect it. Ply does not automatically skip
update menus. Fixture tests establish Ply behavior only, not live-provider safety.

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
New report requests default to `mandate_receipt_v1`; their request and Handoff@2
freeze that choice. Before Task writes, acknowledge the exact mandate, Task,
provider and native session:

```sh
ply workflow run check-in wfr_<digest> --context /absolute/current/context.json --provider codex --native-session <bound-session>
ply workflow run report wfr_<digest> --context /absolute/current/context.json --file /absolute/round-report.json
```

Omitted model, profile, policy and sandbox coverage remain unknown. Optional
observations never establish sandbox coverage or blanket future permission.
Ply reuses native identity, target and input validation; runtime reviews concrete
actions. The launch profile is a request, not observed effective policy. Receipt
adds no commit, push, installation, human QA or delivery authority.

Persisted legacy report contracts retain `workflow run accept` and their original
positive runtime-attestation requirements. Repeating the original start request
never converts an existing reservation to the new protocol. An explicit
`startup_protocol: "legacy_acceptance_v1"` remains available for legacy callers.

A `blocked` or `stopped` pre-start round report can be preserved without positive
acceptance or a native start receipt, including after negative acceptance. It
must report unknown technical qualification and no work effects, tests or managed
artifacts. The report is an immutable stop observation, not a completed round or
native terminal. The coordinator can also preserve refused callback transport
with `workflow run observe-stop`, naming the actual provider/session, action,
evidence and known effect status. A preserved stop blocks work; receipt cannot
clear it or replay the refused action. `follow` returns that stop immediately,
even if the owner has exited, without polling transport or turning the
observation into a working-round report.

When the original report effects were unknown, read `report_startup.source_sha256`
from `workflow run show`, then append the coordinator's actual action-result or
effect-inspection evidence:

```sh
ply workflow run resolve-stop wfr_<digest> --stop sha256:<original-stop> --source sha256:<current-source> --effect-status not_started --basis action_result --observer 'Actual coordinator' --reason 'Concrete observed result' --evidence /absolute/native-result.json
```

The observation binds the original refusal, Run, request, session, exact source
and a frozen copy of the actual evidence. Writer absence alone is not evidence
of past effects. For pre-work retirement, only observed `not_started` or `none_observed` effects
qualify; known or observed `partial` effects stay blocked in that path. Post-start
resource retirement preserves unknown/partial effects without requiring resolution. An
unchanged-source retry reuses the immutable observation, including recovery of
a lost publication reply. The refusal still blocks receipt and working callbacks;
this reconciliation never continues or impersonates the owner.

For a selected replacement or resource closure of a stopped report owner:

```sh
ply workflow run retire wfr_<digest> --observer 'Actual coordinator' --reason 'Actual replacement reason'
```

Retirement observes an already absent pane and agent, complete Herdr inventories
and local Task/session processes. Pre-work source must remain clean at its original
commit, with no native start/terminal or unresolved effects. After native start
(receipt or legacy acceptance), a preserved coordinator stop is required. The
exact worktree/ref/repository and preparation are checked; dirty work and unique
commits are retained. The retirement records actual commit, tree, clean status
and source digest, plus original effects and any prior source-bound resolutions.
Unknown and partial effects remain unresolved; resource closure does not
reclassify them, replay an action or qualify work. Accepted native terminal
history cannot be retired. The first recorded retirement with original unknown/partial
effects emits one sourced warning through the existing logger on stderr, naming
the Run, original effects, source digest and retained record. JSON stdout stays
clean, and identical retries do not repeat the warning. An installed CLI
update does not require the historical launcher bytes or an owner callback for
this coordinator-only exit. Wrong identity, unavailable inspection, a remaining
or unobservable writer still blocks retirement; unknown effects also block the
pre-work path. No source is deleted, reset or committed by retirement. The original request,
acceptance and denial evidence remain history; callbacks close and only the exact
preparation/queue reservation is released. An already supported `advanced` or
`released` historical preparation is reconciled without touching the current
queue owner. This creates no TaskResult, accepted success, human QA or integration.
Identical retries recover durable retirement after interruption. Prepare a fresh Task copy in the same authorized
scope for replacement; it receives its own mandate and acceptance.

A round report uses all native TaskRun Report fields with kind
`ply.workflow.round-report`, schema version 1, plus `round`, `control_id`, and
`previous_report_sha256`. Round zero has null control and predecessor. Include
the frozen Spec's managed `task-requirements` artifact, actual candidate binding,
verifier results and cumulative `budget_usage`. Write each round's evidence to
new paths in the native reply staging directory. Previously reported artifacts
must remain unchanged. Reporting is the last target write of that round.

After check-in, a round that did not finish is still reported truthfully. A
`blocked`, `conflict`, `budget_exhausted` or `unknown` report may give each
requirement its actual `failed`, `not_run` or `unknown` outcome, with the actual
verifier exit and reason. Ply still checks the exact Task, Spec, session,
candidate (including dirty authorized work), requirement ID set, cited verifiers
and artifacts, budget floor and round binding. The report is preserved,
immutable and idempotent on identical retry. It is not a technical pass,
TaskResult, human QA or integration. Only a `complete` report requires every
requirement `passed`; the same unfinished requirements still refuse `complete`.
Ply logs a warning on stderr naming the Run, round and source file. JSON stdout
stays unchanged.

`evidence_gaps` lists only missing evidence that a qualification actually
requires. Any entry blocks `accepted`. Do not use it for facts the protocol
allows to stay unknown, such as model, profile, policy or sandbox coverage
omitted at check-in, or a declared way of measuring time. Report those through
the existing fields, without attesting them:

```json
"evidence_gaps": [],
"process_observation": "Model and effective policy were not observable and remain unknown from check-in; active_seconds is wall-clock time measured by the recipient."
```

Do not hide a missing verifier result, artifact or required observation to get
a pass. Keep it as an evidence gap or as a `not_run`/`unknown` requirement.

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
`accepted` records the coordinator's decision about the exact report and seals
the native workflow terminal once. Reported outcomes, technical assessment,
evidence gaps, findings and budget usage remain unchanged. Acceptance does not
require a passed assessment, an unexceeded budget or a clean candidate, and does
not attest technical pass or publish a TaskResult. Report, session, artifact and
candidate freshness checks still apply; dirty work must match what was reported.
`blocked` preserves the reported outcome and coordinator stop without fabricating
a terminal. Identical reviews are idempotent; conflicting decisions are rejected.
After either seal the coordinator closes the provider session and settles the
Task preparation. A report-only Run does not order human QA or integration;
each needs its own selected gate.

A coordinator-observed provider stop after the recorded native start (check-in,
or positive acceptance under a legacy contract) is not a pre-start stop.
Owner work may exist. The coordinator can use `run retire` after positive writer
absence and exact source/preparation checks, retaining all work and effects.
`resolve-stop` remains available for actual effect evidence, but resource closure
does not require converting uncertainty to a no-effect claim. An unresolved
effect still authorizes no replay or successful delivery.

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
