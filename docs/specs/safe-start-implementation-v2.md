# Safe startup producer and consumer implementation — v2

TC0C1831A · Task `ply-three01-startup-implementation-spec-20261024`.
This contract specifies the next two sequential implementations. This docs delivery does not close the startup bug.

## Scope and source binding

Ply base: `b7a7e435bb22c83f31e2600df69933a15e4e39b9`.
Herdr source and prepared dependency return base: `5da0a01e1eedda054db0c81dd3a780000c40d9f0`.
Both checkouts were clean during inspection. Installed CLI and default server report Herdr `0.8.2`, protocol `20`;
installed Codex reports `0.164.0`. These observations do not identify the installed Herdr source or detection overrides.

| Owner | Existing source and concrete gap |
| --- | --- |
| Ply launch/input | [herdr_transport.go](../../internal/taskrun/herdr_transport.go): `workflowLaunch`, `workflowStartupCall`, `workflowPromptGenerationUntil`; Task submission checks settlement, without current positive prompt evidence. Preserve `workflowStartArgv` provider/model/effort/profile/policy. |
| Ply readiness/bootstrap | [herdr_readiness.go](../../internal/taskrun/herdr_readiness.go): `workflowAwaitReadinessGeneration` accepts session-present plus settled; sessionless bootstrap uses explicit `interactive_ready`. Neither closes the submission race. |
| Ply continuation | [delivery_resume.go](../../internal/taskrun/delivery_resume.go), [delivery_start_recovery.go](../../internal/taskrun/delivery_start_recovery.go), [startup_process.go](../../internal/taskrun/startup_process.go): retain generation, pre-Task, exact exit/no-child and replacement guards. |
| Herdr classification | [Codex manifest](/Users/perottochristensen/github/herdr/src/detect/manifests/codex.toml), [agent_detection.rs](/Users/perottochristensen/github/herdr/src/pane/agent_detection.rs), [terminal/state.rs](/Users/perottochristensen/github/herdr/src/terminal/state.rs): existing update rule recognizes an older layout; `codex_prompt_ready` is startup-oriented and managed Active can outlive current prompt evidence. |
| Herdr API/input | [api/agents.rs](/Users/perottochristensen/github/herdr/src/app/api/agents.rs): `queue_agent_prompt` checks blocked/pending/foreground before queuing text and delayed Enter; `handle_agent_send_keys` has no expected-menu or one-shot guard. [runtime.rs](/Users/perottochristensen/github/herdr/src/terminal/runtime.rs), [Unix actor](/Users/perottochristensen/github/herdr/src/pty/actor/unix.rs), [Windows actor](/Users/perottochristensen/github/herdr/src/pty/actor.rs) own consumption. |
| Herdr compatibility | [schema.rs](/Users/perottochristensen/github/herdr/src/api/schema.rs), [server capabilities](/Users/perottochristensen/github/herdr/src/api/schema/server.rs), [server dispatch](/Users/perottochristensen/github/herdr/src/api/server.rs), [CLI](/Users/perottochristensen/github/herdr/src/cli/agent.rs), [repository rules](/Users/perottochristensen/github/herdr/AGENTS.md). |

The named [incident](/Users/perottochristensen/github/ply/planning/koe/ply-ma-handtere-codex-oppdatering-for-oppdragsprompt.md)
and its `approved-update-terminal.txt` show readiness text during brew and then shell, not a captured skip menu.
Codex tag `rust-v0.164.0` resolves to commit `6aaac694803004a63c6f393d153c96d22521a6b9`.
Its [picker source](https://raw.githubusercontent.com/openai/codex/6aaac694803004a63c6f393d153c96d22521a6b9/codex-rs/tui/src/update_prompt.rs)
maps key `2` directly to `NotNow` and Continue; key `3` persists dismissal and Enter selects the current highlight.
The [upstream render snapshot](https://raw.githubusercontent.com/openai/codex/6aaac694803004a63c6f393d153c96d22521a6b9/codex-rs/tui/src/snapshots/codex_tui__update_prompt__tests__update_prompt_modal.snap)
grounds one proposed layout: version header, release-notes line, ordered Update/Skip/persistent-skip options and footer.
This is authoritative versioned source evidence, not a live terminal capture or proof of installed detection.

## Herdr producer contract (proposed)

Add only a startup inspection and guarded submission for managed Codex before its first Task prompt.
Advertise optional JSON ping capability `safe_startup_v1: true` only when both methods and writer guards exist.
Missing/false capability disables this startup path with a concrete dependency stop; CLI help or bundled schema alone is insufficient.
Installed `herdr agent --help` has neither proposed operation. Add CLI `agent startup inspect --file PATH`
and `agent startup submit --file PATH`, where files contain the corresponding params below; preserve normal JSON result/error envelopes.
Expose the capability in `status server --json`. Keep old `agent.prompt`/`agent.send_keys` shapes and behavior compatible.
Do not attach ignored safety fields to an old method. Add new JSON method types, dispatch, CLI parsing and generated schema;
preserve endpoint generation-1 codecs, method-shape fixtures and existing-value digests without rewriting their expectations.

Inspection request; binding is mandatory and checked, not a hint:

```json
{"id":"inspect-1","method":"agent.startup.inspect","params":{"target":"w9:pX","binding":{"workspace_id":"w9","tab_id":"w9:tX","terminal_id":"term_X","name":"ply-example","cwd":"/absolute/task","provider":"codex","workflow_run_id":"wfr_example","generation":"original"},"attempt_id":null}}
```

Example result, with illustrative opaque values and timestamp:

```json
{"id":"inspect-1","result":{"startup":{"contract":"safe_startup_v1","server_epoch":"epoch-X","managed_launch_id":"launch-X","evidence_token":"opaque-X","observed_at_unix_ms":1792839000000,"surface":"update_skip_once","prompt_ready":false,"launch_pending":true,"native_session":null,"menu":{"layout_id":"codex-0.164-picker-v1","selected_index":0,"skip_action":"skip_once"},"attempt":null}}}
```

`surface` is `update_skip_once`, `provider_prompt`, `blocked` or `unknown`; unknown enum values fail closed.
`prompt_ready=true` requires a current empty editable provider composer, correct foreground identity/cwd, no pending launch,
and no live blocker/working/update/onboarding signal. A native session, idle OSC title, historical Active or absent field never supplies this proof.
Inspection reads current detection-buffer content and cursor, not the scrollable viewport or cached lifecycle classification.
Bind tokens to server epoch, managed launch, exact binding, content/cursor revision, foreground identity and native session (including null).
Expire tokens after at most two seconds; reinspection may refresh evidence, never reset an input reservation.
Echo the validated binding in real results; omitted required fields or identity drift is an error, never readiness.

For `codex-0.164-picker-v1`, require the complete source-grounded layout, all three ordered options, current selection,
and the matching installed provider version. Reject truncation, conflicting surfaces and unknown variants.
Select skip with the single character `2`, with no Enter, arrows, Escape, update action or persistent dismissal.
Herdr owns this mapping; callers cannot provide arbitrary keys. Selection is still token-bound even though key `2` is direct.
The old manifest's recognition alone does not qualify this layout. Follow Herdr's evidence-based detection rules before changing manifests.

Submission request for that observation:

```json
{"id":"submit-1","method":"agent.startup.submit","params":{"target":"w9:pX","binding":{"workspace_id":"w9","tab_id":"w9:tX","terminal_id":"term_X","name":"ply-example","cwd":"/absolute/task","provider":"codex","workflow_run_id":"wfr_example","generation":"original"},"expected":{"server_epoch":"epoch-X","managed_launch_id":"launch-X","evidence_token":"opaque-X"},"attempt_id":"startup-example-skip","kind":"skip_once","text":null,"deadline_unix_ms":1792839001000}}
```

`kind` is `skip_once`, `bootstrap` or `task`. Bootstrap accepts only the existing tool-free readiness text;
Task requires a nonempty text and an exact native session bound by inspection. Skip requires null text and the recognized menu.
Bootstrap requires a sessionless provider prompt; Task requires a session-present provider prompt.
Return `result.submission` with attempt ID, kind, binding, payload digest, outcome, bytes-written status and reason.
Outcomes are `not_sent`, `submitted` or `unknown`; `submitted` means bytes including required submission were written, not model completion.
Inspection with `attempt_id` returns the volatile outcome within the same server epoch, without sending bytes; after restart it is unavailable/unknown.
Validation errors report `not_sent` only with positive zero-byte proof. Disconnect/timeout/partial write never implies zero bytes.

## Atomic submission and durable outcomes

Implement at existing API/runtime/PTY seams, with narrow server-owned data; no generic dialog service or global settings.
Ply alone owns durable consumed reservations for skip/bootstrap/Task; no second durable Herdr journal or restore protocol is required.
Herdr atomically consumes an epoch-bound volatile one-shot token and `(managed_launch_id, kind)` slot before queueing,
binding Run/generation, attempt ID, request digest and deadline. Another ID, payload or generation cannot reuse an occupied slot.
Same-epoch duplicate submissions return the existing outcome without enqueueing; Ply uses read-only lookup, never transport retries.
Restart creates a new server epoch and invalidates all old tokens/outcomes; Ply retains every consumed reservation and treats missing outcomes as unknown.
Fresh tokens, reconnects, restored terminals or new IDs cannot authorize resending a consumed/unknown Ply attempt.
Crash after reservation but before a known write remains conservatively consumed; positive zero-byte evidence does not release the slot.

Use phase-aware guards at the existing PTY writer, with one serialized input transaction and current parsed detection evidence:

- **first-byte admission:** after queueing/delay, require the exact token's fresh content/cursor, binding/epoch/launch/session,
  foreground provider and deadline. Drain/process pending detection output before validation; token expiry sends zero bytes.
  Activate the consumed token's transaction lease at admission; its two-second inspection TTL does not expire an admitted long paste.
- **text-stream:** retain existing paste encoding/boundaries and Unix write offset; own prefix echo, redraw and cursor changes
  do not require the old empty-composer revision or restart validation. Check a cheap revocable identity/input/deadline guard before each write continuation.
  Parsed unrelated blockers/shell/updater/session changes revoke the lease. Serialize competing input; revoke on externally admitted input.
  Windows uses an explicit guarded write loop in place of blind `write_all` for this path; advertise capability only where these guards work.
- **echo-wait:** after all text bytes are acknowledged, process output and wait up to two seconds, clamped to the startup deadline,
  for the same provider's editable composer showing this transaction's text or its version-supported paste-summary representation.
  Delayed echo may succeed within that bound. Do not compare raw full text to a collapsed paste widget; an unrecognized representation,
  unavailable echo or ambiguous attribution times out as partial/unknown, with no Enter. Own echo is attributed only while input remains exclusively leased.
- **Enter admission:** retain the existing submit delay, then require the fresh compatible composer/own-paste evidence,
  unchanged binding/foreground/session, unrevoked lease and remaining deadline immediately before writing Enter on both platforms.

No new write syscall or Enter is admitted after deadline/revocation; already in-flight bytes cannot be retracted and remain possible partial input.
Keep writer/parsed-state admission serialized; an API-layer check followed by an unguarded queued write is insufficient.
Never send cleanup keys or replay partial text. A later idle/prompt appearance does not prove that text or Task submission was absent.
After a lost skip reply, same-launch fresh provider-prompt evidence may establish dismissal and permit observation to continue;
the skip attempt remains consumed and its uncertain write history remains recorded. A remaining/changed menu gets no second key.
Bootstrap uncertainty resolves only through a fresh matching native session and settled positive provider prompt; Task uncertainty never permits resend.

## Ply consumer and sequencing

Use the delivered producer's actual published schema/CLI, then extend only the cited startup/readiness/prompt/resume seams and fixtures.
Check the selected server capability before any startup text; keep unsupported servers usable for unrelated actions.
Hold consumer activation and subsequent real tasks until the named-endpoint gate below passes; keep the existing working controller meanwhile.
Keep Run/Task/request, provider executable and permission controls unchanged; a receipt is not blanket runtime permission.

1. Read preserved state first. Active denial, possible Task input, accepted work, missing reservation evidence or unresolved prior effect stops dependent sends.
2. Inspect the same reserved location/foreground cwd and generation. Recognized menu permits one durable `startup-skip-attempt.json`
   under `workflowStartupPath`; reserve with `workflowUpdate`/generation fence before calling the new submission method.
3. Observe afresh after skip. In both session-present and sessionless branches require explicit positive prompt evidence and settled state.
   Sessionless alone may reserve `startup-bootstrap-attempt.json` and submit the existing readiness exchange once.
   Then bind the fresh stable native session; literal historical `PLY_READY` and skip success are insufficient.
4. Before first Task input, rerun freshness/identity/Task-base/denial checks and inspect again. Atomically reserve `prompt-attempt.json`
   with its evidence/attempt binding, then use guarded submission. Herdr revalidates at consumption; never fall back to `agent prompt` or send-keys.
5. All startup sends use the existing generation send fence, including Task input. Concurrent owners/resume share the durable reservations.
   Reservation with lost response permits read-only lookup, never a new attempt ID. Preserve `prompt_unknown` and prior refusals.

Keep the existing shared startup deadline, normally 90 seconds, covering launch, skip, bootstrap, session bind and Task send.
Clamp calls and producer deadlines to remaining time; no input or delayed Enter after expiry. Skip/reinspection never extends it.
A supported resume gets its own bounded observation window but retains old reservations and their expired send deadlines.
It creates no tab/agent and cannot replay an expired attempt. Existing guarded recovery alone can establish a new generation.

Unknown layout/missing Skip, updater, shell, authentication/trust/permission, optional-model or security surfaces send zero model text and guessed keys.
Report exact surface, phase, evidence and `next_actor=coordinator` for inspection. Required upgrade or app-safety confirmation uses `next_actor=user`.
Neither timeout nor routine autonomy approves upgrade/install/app safety. Human-approved completed upgrade may use existing `recover-start --check`
and guarded apply, preserving history and requiring fresh Task/base/input, exact transport exit and no-child evidence.
A shell, missing launcher or failed read alone proves none of these. Recovery guards remain unchanged.

## RED/GREEN product acceptance

These product tests are planned, not executed by this docs Task. Add each regression first, observe RED at the real consuming branch,
then implement the focused fix and record GREEN. Existing safe rejections should stay green. Use synthetic state/PTY writers and Ply's
[stand-in](../../internal/taskrun/herdr_fixture_test.go); inspect ordered bytes/calls, reservations and outcomes, with zero real upgrades.
Herdr engine tests use synthetic manifests/minimal strings; do not freeze invented provider screens against bundled rules, per its repository rules.

| Case | Proposed focused test and owner | Required assertion |
| --- | --- | --- |
| AT01 skip | Herdr `startup_skip_direct_key_is_one_shot`; Ply `TestWorkflowStartupSkipsOnceBeforeReadiness` | Source-bound layout metadata selects only byte `2`; concurrent/duplicate requests and lost replies write at most once; skip→fresh prompt→optional bootstrap→Task ordering. |
| AT02 blocked | Herdr `startup_unknown_surface_writes_nothing`; Ply `TestWorkflowStartupRejectsBlockedSurface` | Unknown/truncated/missing-skip/updater/shell/security/model-offer surfaces write zero text/keys despite idle/session metadata; exact next actor. |
| AT03 readiness | Ply `TestWorkflowStartupRequiresPositivePromptInBothBranches`; extend [bootstrap tests](../../internal/taskrun/herdr_bootstrap_test.go) | Both session branches reject absent/false evidence, stale Active and identity drift; sessionless bootstrap once, then fresh matching session/prompt. |
| AT04 consumption | Herdr `startup_guard_rechecks_at_writer_and_enter`, `startup_own_echo_allows_partial_paste` | Require positive long/partial paste with own redraw, delayed-echo success and supported paste-summary success. Unrelated surface/identity/input changes or unrecognized echo cancel safely before text/Enter; no self-invalidation on prefix echo. |
| AT05 durability | Herdr `startup_token_is_consumed_within_epoch`; Ply `TestWorkflowStartupUnknownEffectsNeverReplay` | Crash before/after enqueue, disconnect, changed payload/new ID, concurrent resume and server epoch replacement: volatile slots reject duplicates; Ply durable consumed reservations prevent all replay even when producer lookup is unavailable; refusal retained. |
| AT06 deadline | Herdr `startup_expiry_cancels_queued_input`; extend [readiness tests](../../internal/taskrun/herdr_readiness_test.go) | One total launch deadline, expired token/queue/Enter sends nothing further, resume does not revive expired reservations. |
| AT07 recovery | Extend `TestDeliveryStartRecoveryRejectsUnsafeStartEvidence` in [recovery tests](../../internal/taskrun/delivery_start_recovery_test.go) and [generation tests](../../internal/taskrun/delivery_recovery_generation_test.go) | Human upgrade/security remain separate; possible Task input, dirty/drifted base, active launch, remaining child or unavailable exit proof rejects; old-generation Task sender cannot reserve or send. |
| AT08 compatibility | Herdr `startup_capability_requires_guarded_methods`; Ply `TestWorkflowStartupRequiresProducerCapability`, `TestWorkflowStartupBindsQualifiedNamedEndpoint` | Old server/missing fields/unknown contract fails closed locally, no unsafe fallback; codec fixtures stay compatible. Mismatched CLI/socket/caller pane or unqualified endpoint blocks activation; real named-shell execute/review routes to the capable server. |

## Sequential delivery and runtime smoke

1. Root registers Herdr producer implementation first from the prepared local custom dependency return:
   project `ply`, repo `herdr`, Epic `ply-safe-start-herdr`, `/Users/perottochristensen/github/herdr-worktrees/ply-startup-support`,
   `refs/heads/codex/ply-startup-support`. No upstream issue/PR/push/release. Publish the concrete schema, CLI help, capability,
   source OID/tree, binary hash and focused test evidence after separate review and authorized native local return.
2. Use [justfile](/Users/perottochristensen/github/herdr/justfile) recipes: `just test-one FILTER` for RED/GREEN, then `just check` before commit;
   unavailable toolchain/SDK is an explicit evidence gap, not a pass. Build in the isolated implementation worktree.
   For isolated tests clear inherited `HERDR_SOCKET_PATH`/`HERDR_CLIENT_SOCKET_PATH` and use an explicit named session with the private binary.
   Verify its socket/session/capability before input. Debug Codex disables the update flow; synthetic tests cannot substitute for release-layout qualification.
3. Before modifying a bundled detection rule, obtain the required detection-buffer text/ANSI and explain evidence in that isolated session,
   using installed provider/version and the repository's throwaway detection workflow. Observe a naturally available menu; do not provoke upgrade/security.
   Keep preexisting overrides intact and restore temporary ones. The source-grounded modeled contract can be built/reviewed before capture;
   live menu compatibility remains unqualified and no live skip pass is claimed. Do not change bundled manifests without their required capture.
   Producer reports must state this limit separately from modeled tests; an untriggered menu proves no skip behavior. No real upgrade occurs in acceptance.
4. Root deploys the qualified producer binary to an isolated named Herdr session for **real Task routing**, not merely smoke testing.
   Keep the active default server, protected PROD/PLAN and existing working launcher intact. Replacing CLI bytes alone does not upgrade a server.
   Required before consumer activation: observe the actual named session/socket, running producer capability/epoch, binary path/hash and positive/negative
   guarded-submission smoke. Preserve these endpoint facts and results; if any check fails, retain the old controller and hold activation/next tasks.
5. Root owns a visible coordination shell created by that named server. [Pane environment](/Users/perottochristensen/github/herdr/src/integration/env.rs)
   supplies genuine `HERDR_BIN_PATH`, `HERDR_SOCKET_PATH` and workspace/tab/pane values; read and verify them, never copy old `w9:pE` identity.
   Native execute inherits the actual socket through its environment: [workflowCallUntil](../../internal/taskrun/herdr_transport.go) inherits the caller environment,
   and [Herdr socket selection](/Users/perottochristensen/github/herdr/src/session.rs) honors `HERDR_SOCKET_PATH` without an explicit conflicting session flag.
   [Runtime selection](../../internal/taskexecute/runtime.go) resolves `herdr` from PATH; use command-scoped PATH selecting the verified binary's directory,
   verify resolution/hash in execute `--check`, and use the same context for start/review/returns. `HERDR_BIN_PATH` alone does not select Ply's CLI.
   [Native execute](../../cmd/workflow_execute.go) accepts the actual `--herdr-workspace`; live `deliveryCallerPane` checks must match shell ancestry on this endpoint.
   Root may drive that real shell via supported pane commands with real handles/return files while retaining responsibility in this conversation.
6. Root starts Ply consumer from the then-delivered Ply Epic, binds the delivered producer contract and qualified endpoint, and runs bounded AT01–AT08,
   relevant Go/producer tests, separate reviews and native candidate qualification. Gate activation on a successful native preview/start/review from that
   genuine coordination shell and its observed destination; then route the next cancellation and coverage tasks through the same proven endpoint.
   Recheck capability/epoch and routing context before each start; drift blocks sends. No global launcher/config change or new orchestration service is needed.

## This documentation candidate

Requirement mapping: ac-001 → source binding, proposed producer and consumer sequencing;
ac-002 → atomic/durable rules and AT01–AT08; ac-003 → single-file docs checks, clean commit, separate native review and exact qualification.
Commit only this file. External acceptance checks structure, source links, JSON consistency, scope and `git diff --check`;
its RED/GREEN establishes docs evidence only. Root w9:pE owns orchestration, review, later implementation and aggregate notice; no child Slack/internal agents.
Return clean OID/tree and evidence, hold verify until root releases after review PASS, then verify without `--review` using the clean private `make build` binary.
Hold integrate for root's separate release and fresh target check. The only authorized return is Ply Epic `ply_agentic_workflow_support`.
After observed integration/`owner_handoff`, stop source writes and leave cleanup to root; no postdelivery follow/process probes or installation.
