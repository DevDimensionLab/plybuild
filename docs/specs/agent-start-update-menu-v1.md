# Safe agent startup across Codex update menus — v1

> Historisk dokument etter fjerningen av den kompilerte prosesskjeden.
> Execute-, worker/week-, delivery- og integrasjonsorkestreringen som omtales
> her er ikke en kjørbar oppskrift i dagens kilde. Fjernede kildefiler er
> bevart som historiske stier. Se [gjeldende direkte operasjoner](../../README.md#human-owned-process).

T2F4406D3 · Task `ply-startup-update-spec-20261023` · US01 authoring deliverable.
This is an English implementation Spec, not implementation authorization.

## Problem and authority

A delegated start must use the installed provider and reach its actual input
prompt before receiving the mandate. A supported **skip-once** of a recognized
Codex update menu is routine coordinator work. Upgrade/install and confirmation
that a new app is safe to run belong separately to the human. Neither standing
autonomy nor a startup timeout authorizes those choices, global settings changes,
blind Enter, or readiness/Task text into an updater or shell.

The frozen [US01 assignment](/Users/perottochristensen/github/.ply/task-content/v1/objects/sha256/7c8a1070da1f801387e1b7ee5d49e582d5f1f810f62d82001b63dd96deddddfe)
and [user rule](/Users/perottochristensen/github/.ply/task-content/v1/objects/sha256/eefec72e813be334b08bb477942fcf6a3165f3bba0a5f2e6630367369bce0931)
are normative. The frozen [incident](/Users/perottochristensen/github/.ply/task-content/v1/objects/sha256/784e8e6ef75ecbb69a47661c112f365512b3a4b83d8c46f186e118759840585c)
records human-approved DEP02 Codex 0.163→0.164, readiness messages during brew,
removal of the old launcher, and later successful native recovery on the same
Run before Task input. That is recorded history, not fresh observation by this
author or proof that automated menu handling exists. Related optional-model
onboarding is context; generic dialog suppression is outside this slice.

## Current ownership and seams

Ply references below are against Task base
`93766487d28f91c9b6cbdd8f95da3ff955796041`. External Herdr source was read statically
from a clean checkout at `5da0a01e1eedda054db0c81dd3a780000c40d9f0` under
`/Users/perottochristensen/github/herdr`; this does not identify installed bytes
or active detection manifests. Absolute external links identify that evidence
checkout; revalidate its revision before later implementation.

| Owner | Exact existing seam and limit |
| --- | --- |
| Ply launch and input | [`workflowLaunch`](../../internal/taskrun/herdr_transport.go) sends the reserved tab/start and calls readiness before Task input. [`workflowStartArgv`](../../internal/taskrun/herdr_transport.go) preserves provider/model/effort/profile/policy. [`workflowStartupCall`](../../internal/taskrun/herdr_transport.go) fences startup sends by generation. Transport uses existing `tab create`, `agent start`, `agent get`, `agent prompt`; no update-menu command is implemented here. |
| Ply readiness | [`workflowAwaitReadinessGeneration`](../../internal/taskrun/herdr_readiness.go) checks location, provider/session stability and settlement in one bounded window. The sessionless branch requires explicit `interactive_ready=true`; the session-present branch currently uses session plus settlement. [`workflowSettled`](../../internal/taskrun/herdr_transport.go) accepts idle/done with absent/false `launch_pending`. These are transport observations, not proof of a recognized live Codex prompt. |
| Ply bootstrap and mandate | [`workflowBootstrap`](../../internal/taskrun/herdr_readiness.go) records `startup-bootstrap-attempt.json` before one tool-free `PLY_READY` request. It then observes session and settlement; it does not parse a literal reply as independent readiness proof. [`workflowPromptGenerationUntil`](../../internal/taskrun/herdr_transport.go) rechecks target/session and records `prompt-attempt.json` before Task input; lost replies become `prompt_unknown`. |
| Herdr detection and readiness | The bundled [Codex manifest](/Users/perottochristensen/github/herdr/src/detect/manifests/codex.toml) has `startup_update`, a blocked rule for one exact text shape including `Update available!`, `Update now`, and skip-until-next-version wording. It neither selects Skip nor proves other layouts recognized. [`reconcile_managed_agent_at`](/Users/perottochristensen/github/herdr/src/terminal/state.rs) owns managed lifecycle; [`managed_agent_interactive_ready`](/Users/perottochristensen/github/herdr/src/terminal/state.rs) reports Active, which can persist beyond the current prompt observation. |
| Herdr input | [`queue_agent_prompt`](/Users/perottochristensen/github/herdr/src/app/api/agents.rs) rejects blocked/pending/non-provider foreground and queues text plus delayed submission. [`handle_agent_read`](/Users/perottochristensen/github/herdr/src/app/api/agents.rs) and [`handle_agent_explain`](/Users/perottochristensen/github/herdr/src/app/api/agents.rs) expose screen/detection context. [`handle_agent_send_keys`](/Users/perottochristensen/github/herdr/src/app/api/agents.rs) checks the named foreground provider, but provides no expected-menu/cursor fence or idempotent skip identity. Read-then-send is not an atomic safe skip. |
| Existing continuation | [`WorkflowResumeDeliveryStart`](../../internal/taskrun/delivery_resume.go) follows the same pre-Task launch without creating a tab/agent or replaying bootstrap. [`deliveryStartupAction`](../../internal/taskrun/delivery_resume.go) supplies phase-specific guidance. [`deliveryRecoveryPreTask`](../../internal/taskrun/delivery_start_recovery.go), [`deliveryRecoveryPreview`](../../internal/taskrun/delivery_start_recovery.go) and [`deliveryRecoveryExit`](../../internal/taskrun/delivery_start_recovery.go) guard executable replacement and positive exit evidence. [`deliveryRecoveryProcessTree`](../../internal/taskrun/startup_process.go) rejects remaining children or unavailable observations. |

## Bounded startup flow

These are decision states, not new public API names. Apply only to the selected
native Codex startup before any possible Task input. Preserve Task, Run, request,
mandate, provider, model, effort, policy, terminal binding and startup generation.
Keep historical unknowns and actual refusals; never reinterpret a failed read as
provider absence or a mandate receipt as runtime permission.

1. **Inspect the preserved attempt.** A Task prompt reservation, possible partial
   Task input, start receipt, accepted work, active denial or unresolved prior
   effect prevents automatic skip/restart/replay. Return concrete phase, known
   effect status and `next_actor=coordinator` for supported inspection/resolution.
2. **Classify fresh startup evidence.** Require the exact reserved
   workspace/tab/pane/terminal, name, cwd, provider and generation. Herdr must
   distinguish a live supported update menu, a live provider prompt, and an
   unrecognized/blocked surface. A cached idle label, session ID, shell prompt,
   missing flag or historical `PLY_READY` is insufficient.
3. **Known update menu with supported once-only Skip.** Require the full live
   layout, explicit single-launch Skip option, current selection and supported
   provider navigation semantics. Reserve one durable skip attempt for this
   generation before input. Under the Herdr guard described below, select only
   that option; never `Update now` or persistent `Skip until next version`.
   Freshly observe the same startup afterward. Menu acknowledgement is not
   readiness; menu persistence, changed layout or uncertain selection stops
   further keys. Repeated reads/resume cannot reserve another skip.
4. **Live provider prompt.** Before either bootstrap or Task text, require fresh
   positive prompt classification with no update/onboarding/security surface,
   correct foreground provider and no pending launch. Apply this to both current
   session-present and sessionless branches. If no session is bound, retain the
   existing single reserved tool-free readiness exchange; then require the fresh
   matching native session and settled prompt. With an already valid session,
   no extra bootstrap is needed. Recheck at submission, then reserve/send the
   first Task prompt once. Mandate check-in remains the recipient's later act.
5. **Blocked surface.** Unknown layout, missing Skip, optional-model onboarding,
   authentication/trust/permission dialog, updater or shell sends zero model
   input and zero guessed keys. Record the precise supported stop and
   `next_actor=coordinator` for classification or supported existing recovery.
   Upgrade-required and OS/app-safety confirmation use `next_actor=user`, naming
   the separate human upgrade/security action. Do not automatically approve
   trust/permissions or relabel an actual denial as onboarding.
6. **Lost reply or partial input.** Skip, bootstrap and Task sends have separate
   durable reservations and outcomes. Unknown delivery permits fresh read-only
   observation of the same attempt, never replay or a new attempt ID. A lost skip
   reply can progress only if the same provider's fresh prompt proves dismissal;
   remaining/partial menu navigation stays unknown. Partial bootstrap stays
   reserved until supported same-attempt observations resolve it. Possible Task
   input forbids restarting or resending even when the pane later looks idle.
   A known refusal remains refused; supported resolution must name the changed
   action and preserve its evidence.

Use the existing shared startup deadline (normally 90 seconds); skip does not
reset it or create an unbounded loop. Expiry preserves the last observation and
next responsible actor, with no input after expiry. A later supported `resume`
has its own bounded observation window and retains all reservations.

After an actual human-approved completed upgrade and required human app-safety
confirmation, the coordinator may use existing `workflow execute recover-start
RUN_ID --check` and guarded apply only when current guards permit. Preserve the
original request and history; recovery may bind a new native session/generation
and executable identities. It requires fresh Task/base/input checks, no possible
Task input, and positive exact-transport exit/no-child evidence (or the existing
guarded closed-terminal path). A shell alone or missing launcher is insufficient.
This authoring Task runs none of those observations or recovery commands.

## Minimum implementation boundary

**Herdr dependency:** the external repository above owns live detection,
foreground validation and terminal input. Extend its Codex detector/managed
readiness and input guard at the cited manifest, `src/terminal/state.rs` and
`src/app/api/agents.rs`. The minimum required behavior is a positive current
startup classification and a guarded one-shot Skip submission, bound to the
exact terminal/provider, generation-specific attempt key and expected live
menu/selection. Validate again where bytes are submitted; changed evidence sends
no bytes. Preserve attempted/unknown input across lost replies, reconnects and
partial navigation. The readiness/Task submission guard must likewise reject a
surface that changed after observation. This is a required upstream capability,
not an assertion that existing `agent send-keys` provides it. No new CLI name,
wire schema or provider protocol is presumed. The later implementor must bind
the supported published Herdr contract before wiring Ply; until available,
Ply reports the concrete dependency stop rather than faking safe skip.

In Ply, keep changes inside the cited startup/readiness/prompt/resume seams and
their durable state/diagnostics. Add only the per-generation skip reservation and
the minimum evidence binding needed by that supported transport. Reuse existing
identity, deadline, effect-unknown and generation fences. Recovery guards remain
unchanged. Do not turn this slice into a generic dialog platform, provider
upgrade service, new launcher, policy probe or global configuration mechanism.
Claude handling, later-session model offers and OS interaction remain outside
this Codex slice. A Herdr change needs its own authorized repository scope;
this Spec grants none and starts/registers no implementation.

## RED/GREEN acceptance plan

The following product tests are **planned, not executed** here. First add focused
failing cases at the owning layer; verify RED because the expected safe behavior
is missing, then implement the smallest fix and verify GREEN. Existing guard
regressions should stay green. Do not manufacture RED by weakening a guard.

Ply's [Herdr stand-in](../../internal/taskrun/herdr_fixture_test.go), the existing
`workflowStandin` constant, provides observation sequences, failure/lost-reply
injection and `herdr-calls.jsonl`. Extend that fixture with the eventual supported Herdr
contract, not a fabricated production API. Count ordered effects and inspect
durable reservations. Herdr's existing Rust tests can model terminal states,
screens and queued bytes without starting a real provider.

| Case | Required modeled assertion and existing regression seam |
| --- | --- |
| AT01 — recognized menu | Herdr fixture: full supported layout and cursor selects Skip exactly once; no upgrade/install/config action. Alter cursor/layout between read and submission: zero bytes. Ply fixture: one skip, then fresh readiness, then at most one bootstrap and one Task prompt; installed model/policy unchanged. Extend [`managed_agent_readiness_tracks_detection_state`](/Users/perottochristensen/github/herdr/src/terminal/state.rs). |
| AT02 — unknown/security | Unknown menu, missing Skip, OS safety, trust/permission and optional-model surfaces: zero bootstrap/Task text and zero guessed keys, concrete reason/next_actor. Contradictory idle/session metadata cannot override a live blocker. Extend [`TestWorkflowBootstrapRequiresExplicitInteractiveReadiness`](../../internal/taskrun/herdr_bootstrap_test.go) and Herdr detector/guard fixtures. |
| AT03 — updater/shell | Simulated updater output, provider exit or stale Active flag plus shell receives zero readiness/Task bytes. Include session-present idle with no live provider prompt to expose the current shortcut. Extend [`TestWorkflowReadinessWaitsWithinOneStart`](../../internal/taskrun/herdr_readiness_test.go) and [`queue_agent_prompt`](/Users/perottochristensen/github/herdr/src/app/api/agents.rs). |
| AT04 — readiness ordering | Skip success alone, absent readiness, changed identity/session/generation and deadline exhaustion never authorize Task input. Fresh prompt→optional bootstrap→bound settled session→first Task input is ordered. Retain [`TestWorkflowBootstrapDiscoversSessionBeforeTask`](../../internal/taskrun/herdr_bootstrap_test.go) and [`TestWorkflowReadinessUsesOneDeadline`](../../internal/taskrun/herdr_readiness_test.go). |
| AT05 — uncertainty | Inject lost/partial skip, bootstrap and Task replies, crash after reservation, concurrent resume and stale-generation senders. At most one attempt each; unknown remains unknown unless same-attempt evidence resolves it; refusal never replays. Extend [`TestWorkflowBootstrapLostReplyUsesFreshSessionWithoutResend`](../../internal/taskrun/herdr_bootstrap_test.go), [`TestDeliveryResumeRefusesUncertainPromptAndIdentityDrift`](../../internal/taskrun/delivery_resume_test.go) and [`TestDeliveryRecoveryGenerationOldReadyOwnerCannotReserveTaskPrompt`](../../internal/taskrun/delivery_recovery_generation_test.go). |
| AT06 — human recovery | Model a separately human-approved completed upgrade with app-safety confirmation; existing guarded recovery preserves Task/Run/request/history and can rebind executable/session. Dirty/drifted Task, possible Task input, active managed launch, remaining child or unavailable exit evidence rejects with no start. Retain [`TestDeliveryStartRecoveryPreservesHistoryAndUsesNewCallbacks`](../../internal/taskrun/delivery_start_recovery_test.go), [`TestDeliveryStartRecoveryRejectsUnsafeStartEvidence`](../../internal/taskrun/delivery_start_recovery_test.go) and [`TestDeliveryRecoveryManagedPendingAgentIsNotAnExitedProvider`](../../internal/taskrun/delivery_recovery_generation_test.go). |

Fixtures establish modeled decisions and byte ordering, never actual provider
layout support, OS policy, human approval or installed behavior. A later permitted
installed smoke must record those separately; an untriggered menu proves nothing.
No full-suite/coverage target or live upgrade/browser/process/quota/permission
probe is required by this Spec.

## Delivery boundary

Requirement mapping: **US01-01** ownership/source table; **US01-02** bounded flow
and uncertainty rules; **US01-03** AT01–AT06 and human boundaries; **US01-04** this
authoring delivery and its exact candidate evidence.

For T2F4406D3, commit only this file. The frozen external acceptance entrypoint
checks document completeness, source links/declarations, sole-file scope and
`git diff --check`; its RED/GREEN is docs evidence only. Return the clean
OID/tree and logs to coordinator w9:pE for separate native review. Hold owner
verify/integrate until review is sealed and the coordinator signals continuation.
Then use the exact candidate's native verify without `--review`, followed by the
single authorized `local_epic_integration` into `ply_agentic_workflow_support`.
No binary installation, push, PR, main/master changes, extra Slack, internal
agents or direct Herdr start. Parent owns final notice/follow; retain honest
source/cleanup pending facts and perform no postdelivery follow/process probe.

A later product implementation requires its own selected mandate, authorized
Herdr dependency, separate native review, candidate-bound relevant tests/native
verify and agreed local Epic integration. Product delivery also requires a
qualified installation and recorded installed smoke under that later authority;
human-only upgrade/app-safety choices still apply. This docs delivery neither
authorizes installation nor certifies the future product behavior.
