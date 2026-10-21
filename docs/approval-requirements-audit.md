# Approval requirements audit — reconciled through AL11

Current source: `c6477d6094628d74eb6a1d04b41207174d9f93ea`, 21 October 2026,
on `ply_agentic_workflow_support`. AL12 reconciles documentation only. The
[historical AL02 baseline](#historical-al02-baseline--21-october-2026) below
preserves the original inventory, measured costs and proposed P1–P5 decisions
at `d335e70e`. Its line references describe that old source. The current
decisions and evidence below supersede its remaining-work claims.

The selected backlog is still open. Delivered repairs, an accepted report-only
pilot and unresolved product findings are different outcomes. This audit does
not deliver the unresolved findings or declare the whole backlog complete.

## Current delivered behavior and decisions

| Area | Current decision and delivered behavior | Current source and regression reference |
| --- | --- | --- |
| AL01 receipt, retained | **keep** Mandate receipt acknowledges the exact owner/cwd/mandate; unknown runtime facts stay unknown. A real refusal remains a refusal. | [report_startup.go](../internal/taskrun/report_startup.go), [execute receipt and stops](workflow-execute.md) |
| AL02 audit | Documentation delivered at `922961b`; no product change or ordinary CLI installation. | Historical baseline below; native `al02-reviewed.json`, `al02-closeout-follow.json`, `al02-task-final.json` |
| P1 / AL03 | **reconcile/remove, delivered** Optional Codex trust is `redundant` or `withheld` with a useful stderr warning; the ordinary start continues. Pre-launch drift withholds only the grant. | [herdr_codex_trust.go](../internal/taskrun/herdr_codex_trust.go), [trust tests](../internal/taskrun/herdr_codex_trust_test.go): `TestWorkflowCodexTrustDriftAfterReservationWithholdsGrant`, `TestWorkflowCodexTrustCandidateCLIJourney` |
| P2 / AL04 | **reconcile, delivered** After a report session is bound, provider/Ply executable or confirmed trust drift warns and reporting continues under the same report-only authority. Herdr, saved identity, frozen artifacts and contradictory trust bindings still fail validation. This did not implement the originally proposed immutable control copy/continuation for reports. | [herdr_store.go](../internal/taskrun/herdr_store.go): `workflowFreshRuntimeOperation`; [report_continuity_test.go](../internal/taskrun/report_continuity_test.go): `TestAL04ReportCLIJourneyAcrossInstall`, `TestAL04PartialPromptEffectAfterDrift` |
| AL05 input roles | **simplify, delivered** Ordinary reports, evidence, acceptance scripts and candidate reviews accept 0644, bind actual bytes and warn. Group/other writable inputs are refused with role/path/reason. Authority-bearing inputs and native state retain owner-only mode. Mode bits do not authenticate authorship. | [schema.go](../internal/taskrun/schema.go): `readFileAs`; [al05_file_inputs_test.go](../internal/taskrun/al05_file_inputs_test.go): `TestAL05OrdinaryRoundReportIsConsumed`, `TestAL05ChangedOrdinaryInputsDoNotQualify` |
| P3 / AL06B + AL06C2 | **simplify, technical change delivered and actual pilot accepted** Report requests freeze `manual` or `auto` and forward that mode to Claude. Actual read-only Claude auto pilot completed separately; see the P3 outcome for AL06/AL06C1 failures and the report-only boundary. | [schema.go](../internal/taskrun/schema.go), [herdr_transport.go](../internal/taskrun/herdr_transport.go), [report_claude_mode_test.go](../internal/taskrun/report_claude_mode_test.go): `TestAL06ClaudeAutoReportCLIJourney`, `TestAL06ClaudeReceiptPreservesUnknownOrDifferentObservation` |
| P4 / AL07 | **simplify, delivered** Already-authorized automatic closeout/installation recovery names the coordinator; unknown installation asks for observation. Completed/retired history exposes no fictitious reviewer. Human-owned and headless journeys retain their proper actors. | [read_closeout.go](../internal/taskrun/read_closeout.go), [progress.go](../internal/workspaceview/progress.go), [al07_read_actions_test.go](../internal/taskrun/al07_read_actions_test.go): `TestAL07ConsumedInventoryRecovery`, `TestAL07CandidateCLIReadActionsDoNotWrite` |
| P5 / AL08 | **implement, delivered and actually installed** Native `workflow execute install-cli` consumes the retained qualified executable after native delivery/closeout, with attributed coordinator operation, backup, receipt and recovery. It replaces per-job `install-qualified.py` copies for this self-hosted Epic. | [CLI installation](cli-installation.md#coordinator-installation-after-native-self-hosted-delivery), [delivery_cli_install.go](../internal/taskrun/delivery_cli_install.go), [qualified_test.go](../internal/clirelease/qualified_test.go): `TestQualifiedInstallCannotReplayAfterLegacyRelease` |
| AL09 follow | **reconcile, delivered** Both follow paths distinguish the observation window from the Herdr request budget. At most one read continuation handles a request timeout; persistent failure stays unknown with a next follow command. No prompt/start is replayed. A later tab-disappearance race is still open. | [follow_observation.go](../internal/taskrun/follow_observation.go), [al09_follow_observation_test.go](../internal/taskrun/al09_follow_observation_test.go): `TestAL09FollowWindowClipsHealthyRequest`, `TestAL09ReportFollowObservation` |
| AL10 unfinished reports | **implement, delivered** Blocked/unknown report-only returns preserve failed/not_run/unknown requirements with exact Spec/candidate/artifact binding. Only complete reports require all requirements passed. Warning includes Run/round/source; receipt/legacy and post-start guidance name the actual phase. No technical pass, TaskResult, QA or integration is created by an incomplete report. | [herdr_callback.go](../internal/taskrun/herdr_callback.go), [task_spec_guard.go](../internal/workspace/task_spec_guard.go), [report_incomplete_test.go](../internal/taskrun/report_incomplete_test.go): `TestAL10IncompleteReportCorrectsToCompleteAcceptedReturn`, `TestAL10ProviderStopAfterLegacyAcceptanceNamesRecordedStart` |
| S4 / AL11 | **simplify, delivered**, replacing historical **keep**. A new receipt report can bind `launch_contract_pending_runtime_acceptance`, selected profile, empty policy hash and empty evidence; a redundant coordinator policy document is unnecessary. Optional evidence remains hash-checked. Legacy positive-acceptance, terminal/factory and Delivery contracts retain their requirements. | [launch binding](workflow-run.md#choose-claude-or-codex), [schema.go](../internal/taskrun/schema.go): `parseRequestForSurface`; [report_launch_contract_test.go](../internal/taskrun/report_launch_contract_test.go): `TestAL11LaunchContractBoundaries`, `TestAL11ReportLaunchContractCLIJourney` |

## Evidence levels and native outcomes

Evidence directory (local raw records, not portable product prerequisites):
`/Users/perottochristensen/.local/share/ply/planning/backlog-sequential-20261021/`.
The following records were read alongside the source. Native results outrank
stale planning status; this table is a source index, not a second status system.

| Work | Delivered source | Native/evidence records in that directory |
| --- | --- | --- |
| AL02 | `922961b` | `al02-reviewed.json`, `al02-closeout-follow.json`, `al02-task-final.json` (retained native artifacts) |
| AL03 | `f5535ef` | `al03-reviewed.json`, `al03-closeout-follow.json`, `al03-installation.json` |
| AL04 | `7e049ed` | `al04-review2-readback.json`, `al04-closeout-follow.json`, `al04-installation.json` |
| AL05 | `c02e3ee` | `al05-review3-readback.json`, `al05-native-final-readback.json`, `al05-closeout-follow.json`, `al05-installation.json` |
| AL06B | `1040987` | `al06b-review1-readback.json`, `al06b-native-final-readback.json`, `al06b-closeout-follow.json`, `al06b-installation.json` |
| AL07 | `7fde35b` | `al07-review2-readback.json`, `al07-native-final-readback.json`, `al07-closeout-follow.json`, `al07-installation.json` |
| AL08 | `2941fcc` | `al08-review2-readback.json`, `al08-native-verify3.json`, `al08-follow-current.json`, `al08-installation.json`, `al08-native-install.json` |
| AL09 | `3e09cab` | `al09-review1-readback.json`, `al09-native-verify2.json`, `al09-closeout-follow.json`, `al09-installation.json` |
| AL10 | `4f822789` | `al10-review2-readback.json`, `al10-native-verify2-verification.json`, `al10-completed.json`, `al10-installation.json` |
| AL11 | `c6477d6` | `al11-review1-readback.json`, `al11-native-verify2-verification.json`, `al11-completed.json`, `al11-installation.json` |

These product Tasks have exact passed independent review, native verification
and local Epic delivery/closeout evidence; the installed smoke records are
separate. Completed historical readback may show `current=false` or
`inputs_match=false` after source removal: that is not a reversal of its durable
delivery receipt. AL06 and AL06C are deliberately not rows of delivered Tasks.

**Actually executed CLI/native:** Candidate and installed CLI journeys consume
real CLI processes and disposable native stores. AL03–AL07 have coordinator
installation receipts; AL08–AL11 use native `install-cli` receipts and installed
feature smoke. This proves the consumed CLI behavior, not live provider policy.

**Modeled provider/Herdr:** Those regression journeys use synthetic provider
and Herdr stand-ins. The P1/P2 after-columns below are fixture results, not
remeasurements of the original plyapp journey. No new after-timing or broad
question-count study was performed by AL12. AL09 reproduced a clipped-deadline
mechanism; missing timing in old failures prevents claiming it caused all of them.

**Actual Claude auto pilot:** `al06c-result.json` and its native report/review
records demonstrate actual reads, two corrections and an accepted report-only
return. They do not turn synthetic refusal tests into real runtime denials or
claim every possible provider action is authorized. See P3 outcome below.

**Known evidence limits:** AL05's earlier broad development run stopped
incomplete at `TestDeliveryCLILocalJourneys/local_branch_integration/master`
(44 observed passes, no observed failure); it is not green broad-suite evidence
(`al05-broad-development-incomplete.json`). AL10's broad development run was
also stopped/incomplete, with no actual Go exit observed; it is not green
(`al10-broad-development-observed.json`, `al10-owner-al10-broad-stopped.log`).
AL11's initial unit red exists in
the owner transcript, not a dedicated log. Its first green attempt failed in
the helper that mixed warning stderr with JSON stdout; that log was overwritten
by the passing rerun. `al11-owner-intermediate-attempts.txt` preserves the limit;
the later passing logs and old-CLI red are separate evidence, not reconstructed
missing logs. No Go suite was run for this documentation reconciliation.

## Retained boundaries and responsible actors

The human chooses product scope, target and required product judgment once.
The owner reuses that frozen mandate for authorized work/commits; the coordinator
owns routine review, already-authorized recovery and installation. Missing
scope or an unselected effect is a real decision for the human. Required human
QA remains human-owned; automatic acceptance does not fabricate that judgment.

Exact preview/apply binding, owner/session/candidate identity, separate review,
byte-bound acceptance and preservation of unknown/partial effects remain useful
product checks. They are agent-owned commands, not repeated human approvals.
Receipt is not runtime permission: Codex on-request/auto_review and Claude
manual/auto still control their own actions, trust and login dialogs. Ply does
not answer those dialogs. Actual denials stay preserved; source reading or
policy labels cannot clear them. Unknown earlier effects block dependent work
until evidence resolves them; they do not justify locking unrelated safe work.

The user's continuation principle still applies: remove a redundant condition,
reconcile ordinary state or give a useful warning with Run/path/cause/handling
through the existing logger. A warning is not an acknowledgement gate. Keep a
stop where the next effect lacks authority, identity or a known prior outcome,
and identify its supported next action and actor. No new process gate is added
by this audit. Document-only acceptance uses the unchanged qualified CLI and
source/reference checks, with no artificial rebuild or installation.

## Concrete unresolved work

These are existing selected backlog items for root to coordinate. AL12 changes
their documentation, not their product behavior. Each link identifies the
existing task to avoid duplicating findings in a new ledger.

| Open finding | Actual evidence and boundary | Existing backlog / next useful check |
| --- | --- | --- |
| Report closeout after start | AL06C1 ended the provider and released preparation after a sourced stop; it was not accepted/delivered/retired. AL10 corrected misleading wording only. `report_startup.go` still refuses pre-work retirement once native start exists. | [post-start closeout](/Users/perottochristensen/github/ply/planning/koe/ply-rapportstopp-etter-start-skal-kunne-avsluttes.md): consumed started-report stop, known absence, same-source closeout, lost reply and live/unknown-effect controls |
| Review-return text limit | `deliveryParseReviewReturn` still requires each evidence gap to be one line and 1–600 characters. AL05, AL07 and AL10 reviewers reformatted actual returns; this was input validation, not a provider refusal. | [review text](/Users/perottochristensen/github/ply/planning/koe/ply-reviewretur-skal-tale-vanlig-forklaringstekst.md): accept ordinary longer/multiline explanation within a useful total bound, preserve findings and exact retry |
| Final follow race | During AL10 closeout the finished owner tab disappeared between observations; `agent_not_found` returned unknown. Next follow closed the same delivery. AL09's budget repair remains delivered. | [follow race](/Users/perottochristensen/github/ply/planning/koe/ply-follow-skal-tale-forbigaende-observasjonsfeil.md#ny-avgrenset-race-under-al10-closeout): deterministic inventory/get disappearance; reconcile positive absence, preserve conflicting/unknown observations |
| Concurrent in-place file rewrite | Source-only gap: same inode/same size during `readFileAs` is not a reproduced product failure. AL05's post-binding fault proves later digest mismatch, not stable bytes during reading. | [file binding](/Users/perottochristensen/github/ply/planning/koe/ply-filbinding-samtidig-overskriving-under-lesing.md): deterministic read-boundary seam and behavior before selecting a fix |
| Useful build identity | AL09, AL10 and AL11 installed version output has empty `source_commit`; native commit/tree/hash receipts remain exact. The cause of absent embedded metadata is unproven. | [build identity](/Users/perottochristensen/github/ply/planning/koe/ply-kvalifisert-cli-skal-ha-nyttig-buildidentitet.md): inspect actual build flags/fallback; preserve unknown, avoid a new mandatory build gate |
| Simple authorized quota backup | Original AL06 owner/reviewer hit real session quota; AL10 owner hit actual weekly limit after successful integration (UI reset 25 October 17:00 Europe/Oslo). Standing Codex backup enabled continuation but is not a delivered native backup feature. | [quota backup](/Users/perottochristensen/github/ply/planning/koe/providerkvote-for-forste-handling-skal-gi-enkel-backup.md): preserve old work/stop, exact owner absence and fresh backup mandate; no quota probe, account/policy change or denied-action replay |

The older unfinished-requirements report finding overlaps AL10 and is delivered
at that scope; it must not become a duplicate of the separate post-start
closeout gap. The final follow race belongs to the existing follow task;
AL11's empty embedded commit overlaps the existing build-identity task. No new
product defect was reproduced by AL12's source reconciliation.

## Reusable difficult tests already delivered

These are existing behavior tests, not new prose-checking tests or tests rerun
by AL12. CLI tests below execute the selected binary with modeled provider/Herdr.

| Boundary | Permanent regression and what it proves |
| --- | --- |
| Grant drift / lost start reply | `TestWorkflowCodexTrustDriftAfterReservationWithholdsGrant`, `TestWorkflowCodexTrustLostStartReplyKeepsDecision`: withhold only the grant; no duplicate start |
| Report continuity / partial prompt | `TestAL04IdentityHistoryAndArtifactGuardsAfterDrift`, `TestAL04PartialPromptEffectAfterDrift`: preserve identity/history and unknown prompt state without replay |
| File drift / ordinary modes | `TestAL05InPlaceOverwriteAfterBindingDoesNotQualify`, `TestAL05CandidateCLIJourneyWithOrdinaryModes`: later changed bytes do not qualify; normal 0644 input works |
| Claude receipt / refusal | `TestAL06ClaudeAutoLostStartAndPreservedRefusal`, `TestAL06ClaudeReceiptPreservesUnknownOrDifferentObservation`: synthetic denial survives retry; first unknown/different observation remains immutable |
| Recovery actors / historical reads | `TestAL07ConsumedInventoryRecovery`, `TestAL07CandidateCLIReadActionsDoNotWrite`: correct coordinator action, observe unknown install, no write during readback |
| Installation interruption / later selection | `TestQualifiedCLIInstallationInterruptionAndKnownSmokeFailure`, `TestQualifiedInstallCannotReplayAfterLegacyRelease`: known failure evidence, guarded recovery and no old operation replacing a later selection |
| Follow budget / persistent failure | `TestAL09FollowWindowClipsHealthyRequest`, `TestAL09FollowRequestTimeoutContinuation`, `TestAL09ReportFollowObservation`: distinguish window expiry from request failure, one read continuation, unknown remains unknown |
| Incomplete report / correction | `TestAL10IncompleteReportAfterReceiptIsPreservedWithoutQualification`, `TestAL10IncompleteReportCorrectsToCompleteAcceptedReturn`: truthful unfinished requirements, no false qualification, same-session correction |
| Launch-only contract / legacy | `TestAL11LaunchContractBoundaries`, `TestAL11ReportLaunchContractCLIJourney`: no redundant policy file for new receipt reports; optional evidence drift, wrong identity, invalid mode and legacy boundaries remain checked |

The linked current source files above locate these tests. Installation tests
also live in [delivery_cli_install_test.go](../internal/taskrun/delivery_cli_install_test.go).
The unresolved table identifies difficult scenarios still to implement; none
is silently counted as tested by this inventory.

# Historical AL02 baseline — 21 October 2026

Base: `d335e70e` on `ply_agentic_workflow_support` (AL01 delivered). AL02 changed
no product code. The inventory, source line numbers, costs and
repair proposals in this historical section describe that base, before AL03–AL11.
They are retained as the original reasoning, not as current open requirements.
Outcome subsections explicitly identify later delivered behavior.

Two principles from [planning/goal.md](/Users/perottochristensen/github/ply/planning/goal.md)
guide the decisions. First, "Utviklingsløpet skal fortsette; stopp krever en god
grunn": a locked development process is a product defect. Second, root cause
comes before compensation. Each decision is therefore one of the following:

- **remove**: delete the requirement.
- **simplify**: keep its value at lower cost.
- **reconcile**: have Ply resolve the state itself and continue.
- **implement**: add the missing handling so the journey can continue.
- **reuse-mandate**: the human already chose this, so do not ask again.
- **keep**: the requirement stays, with a stated stop basis and a supported path forward.

Where work can safely continue past a doubt, the recommended signal is a `warn`
in Ply's existing logrus logger (historical `cmd/root.go:109`). It must not be a new receipt
gate or a human acknowledgement.

Evidence classes:

- **observed**: persisted evidence that already exists, such as a transcript or a run file.
- **source**: code or text read at the base.
- **modeled**: a derived journey that was not run.

No new timing measurement and no external provider start were made for this
audit. The only times quoted below come from existing transcript timestamps.

## Historical reported AGENTS/policy question

The human reported this question from the plyapp planner: "Plys automatiske
Codex-tillitskontroll krever at du inspiserer …/AGENTS.md og policyen før en ny
start-preview …". It is traced to the actual rule as follows.

| Step | Evidence |
| --- | --- |
| The planner added the optional `codex_project_trust` grant (`process-local`, root `/Users/perottochristensen/github/plyapp-server`) to its report start request. | `/Users/perottochristensen/github/plyapp/planning/reports/pwa-status-http-implementation-20261020/codex-start-request-extra-trust-rejected.json` |
| Preview walked the cwd and repository ancestry and found the tracked `AGENTS.md` in the Task worktree. This triggers the stop "project config, hooks, rules or instructions could be activated". | `internal/taskrun/herdr_codex_trust.go:215-240`, `internal/taskrun/herdr_codex_trust.go:231` |
| The stop text tells the caller to involve the human: "Have the human inspect this path and policy before requesting a new preview". This is the only product message that sends an agent-observable fact to the human. | `internal/taskrun/herdr_codex_trust.go:60` |
| Native output arrived at 20:28:33Z. The planner then asked the human at 20:29:04Z, and the human answered "for noe tull" at 20:29:39Z. | `/Users/perottochristensen/.codex/sessions/2026/10/20/rollout-2026-10-20T21-46-08-01a15ad9-8ab6-7ac0-a473-9a13b8025f83.jsonl:676`, `…01a15ad9-8ab6-7ac0-a473-9a13b8025f83.jsonl` lines 684 and 693 |
| The planner found that the root was already `trusted` in `~/.codex/config.toml`, removed the grant and withdrew the question. Preview without the grant succeeded at 20:30:09Z. | `/Users/perottochristensen/github/plyapp/planning/reports/pwa-status-http-implementation-20261020/trust-resolution.json`, transcript lines 714 (withdrawal) and 717 (successful preview) |

**Cause, verified:** the product's stop wording, combined with a redundant
optional grant. Standing agent policy already forbids routine confirmations
(`/Users/perottochristensen/.agents/AGENTS.md:11-12`, `/Users/perottochristensen/.agents/AGENTS.md:19-35`,
`/Users/perottochristensen/.agents/skills/ply-human-input/SKILL.md:22-26`). The
product text overrode that policy. Earlier grants for `ply/main` passed because
that repository has no `AGENTS.md`, which made the grant look like a safe pattern
to copy.

**Not proven:** why the planner added the grant. No skill or Ply document
prescribes it for ordinary starts; the likeliest source is the earlier Ply-repo
review requests. Also unverified is whether Codex activates a project
`AGENTS.md` differently depending on trust, which decides whether `AGENTS.md`
belongs in the block list at all. The second blocker in the same journey (the
profile/policy acceptance lock) was a separate defect and was fixed by AL01.

Cost observed in this journey: 1 human question, a 35 s human answer, and 96 s
from the stop to an accepted preview. AL01 did not touch this rule, so the
question still recurred at the AL02 base for a request carrying that optional
grant. AL03 subsequently removed this whole-start lock (P1 outcome below).

## Historical inventory at d335e70e

| ID | Journey | Layer | Requirement | Source | Value | Mandate | Actor | Cost | Evidence | Decision | Path forward |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| S1 | startup | product | Read-only preview, then `--apply --confirm <digest>` with `HERDR_ENV=1` | `docs/workflow-run.md:67-72` | Apply cannot run effects that differ from the checked preview | Start already requested | requesting agent | 1 extra command, 0 questions | source | **keep** The digest binds exact effects and needs no human | A fresh preview gives a new digest |
| S2 | startup | product | The `workflow execute` command is the human mandate (goal, scope, delivery target) | `docs/workflow-execute.md:186-190` | Prevents agents from choosing product scope | This is the mandate | human | 1 action per Task | observed | **keep** This is a real product decision, taken once | n/a, not a stop |
| S3 | startup | product | Report-run Claude starts are forced to `--permission-mode manual`; schema accepts only `manual` | `internal/taskrun/herdr_transport.go:335`, `internal/taskrun/schema.go:251`, `cmd/workflow_run.go:43` | None beyond the runtime itself. Delivery already allows `auto` (`internal/taskrun/delivery_service.go:115`) | Code/test mandate given | runtime | A native prompt for each non-allowlisted tool action | source, modeled | **simplify** Accept `auto` as Delivery does (P3) | Owner works under the runtime's own review |
| S4 | startup | product | A report request must carry effective-policy evidence | `internal/taskrun/schema.go:254-256` | Records which launch policy was requested | Requester-produced | requesting agent | 1 artifact, 0 questions | source | **keep** Agent-owned, and after AL01 it is not an acceptance gate | n/a |
| F1 | folder-trust | product | An optional Codex grant stops the whole start when `AGENTS.md`, `.codex` or `.agents` exists in the ancestry | `internal/taskrun/herdr_codex_trust.go:231` | Avoids silently trusting project config, hooks and rules | Ordinary start mandated; grant optional | requesting agent | 1 question, 96 s (observed) | observed, source | **reconcile** Withhold only the grant effect, `warn`, and continue the ordinary start; report a redundant grant when the root is already trusted (P1) | The Codex native trust policy still applies |
| F2 | folder-trust | product | Stop text: "Have the human inspect this path and policy" | `internal/taskrun/herdr_codex_trust.go:60` | None. The agent can read the path and the config | Agent policy forbids the question | requesting agent | Caused the reported question | observed | **remove** Address the agent and state the handling Ply chose (P1) | n/a |
| F3 | folder-trust | product | A confirmed grant rechecks every trust path hash on each later operation and stops with "changed since confirmation" | `internal/taskrun/herdr_codex_trust.go:413`, `internal/taskrun/herdr_store.go:230` | Valid only before launch. After launch the argv grant is already applied | Run accepted | ply | Any later edit to `~/.codex` blocks callbacks | source, modeled | **reconcile** Check before launch; after launch, `warn` and continue (P2) | Pre-launch drift withholds the grant |
| F4 | folder-trust | product | Claude folder trust is written automatically for execute; a failed write warns and the start continues | `docs/workflow-execute.md:160-174`, `cmd/workflow_execute.go:227` | Avoids the native trust dialog | Bound in launch intent | ply | 0 questions; at most 1 native dialog if the lock fails (`docs/workflow-execute.md:274-277`) | source | **keep** Already follows the warn-and-continue model | The dialog is answered in the same tab, then `execute resume` |
| C1 | code-test | product | Mandate receipt `check-in` before Task writes | `docs/workflow-execute.md:192-210` | Binds the exact mandate, session and cwd to the owner | Mandate | owner agent | 1 command, 0 questions | observed | **keep** Already a pure receipt, observed in this Run | Idempotent retry |
| C2 | code-test | runtime | Per-action runtime review (Claude auto classifier, Codex on-request with auto_review) | `docs/workflow-execute.md:176-184` | Real sandbox and approval enforcement | Provider policy | runtime | Prompts only when the provider asks | source | **keep** Runtime boundary; denials are preserved | `observe-stop`, then `resolve-stop` |
| C3 | code-test | instruction | Personal rules: no commit without explicit instruction; stop and ask when unsure | `/Users/perottochristensen/.claude/CLAUDE.md:68-70` | Stops unrequested Git effects | The mandate explicitly authorizes commits | owner agent | 0 questions when the mandate is read | source | **reuse-mandate** The frozen mandate is that explicit instruction; no edit to personal policy | n/a |
| C4 | code-test | instruction | The orchestrator approves routine subagent work and access itself | `/Users/perottochristensen/.agents/AGENTS.md:19-35` | Removes routine human approval | Standing rule | coordinator | Reduces cost | source | **keep** This is the rule F2 violated | n/a |
| R1 | review | product | Independent review round started by the coordinator on the same Run | `/Users/perottochristensen/.agents/skills/ply-delivery/SKILL.md:52-58` | A separate reviewer for the exact candidate | Agreement gate | coordinator | 0 questions (AL01 had 2 rounds) | observed | **keep** Real quality gate with no human | Findings, a new commit, then a new round |
| V1 | verification | product | Automatic acceptance needs the exact clean candidate, a built `--candidate-binary` and the frozen acceptance script | `cmd/workflow_execute.go:400` | Binds what was tested to what is delivered | Agreement | owner agent | 1 build per candidate, 0 questions | source | **keep** Exact candidate binding is agent-owned; AL02 used the unchanged CLI for document-only acceptance | A completed receipt is recovered without rerun |
| D1 | delivery | product | Automatic local Epic integration has no human gate when `require_human_qa=false` | `internal/workflowhandoff/delivery.go:376-377` | n/a | Agreement | owner agent | 0 questions | source | **keep** Already gate-free | n/a |
| D2 | delivery | product | Human QA when the agreement requires it | `docs/workflow-execute.md:557-566` | Real product judgment | Chosen in the agreement | human | 1 question | source | **keep** A genuine product decision | Fail or blocked is preserved and vetoes delivery |
| D3 | delivery | instruction | Ask for the delivery mode only when it is missing | `/Users/perottochristensen/.agents/skills/ply-delivery/SKILL.md:77` | A missing target is a real choice | Reuses an earlier choice | human | 0 when given | source | **keep** Asks only for a missing choice | n/a |
| I1 | installation | product | Native closeout `installation_pending` when the plan has an install profile | `internal/integration/execute.go:245` | Integration is not claimed complete before a required install | Agreement | ply | 0 questions | source | **keep** Resumable, with the source preserved | `resume_installation` |
| I2 | installation | instruction | Self-hosted Ply CLI install after delivery uses the same ad-hoc `install-qualified.py`, copied into 4 job folders (identical sha256 `02736ab7…`) | `/Users/perottochristensen/.local/share/ply/planning/report-startup-lock-20261020/install-qualified.py:1` | Hash-checked install with a backup | Selected per job by the human | coordinator | A per-job copy outside native receipts | observed | **implement** Investigate the native install profile for the Epic (P5) | Script backup and rollback |
| X1 | recovery | product | Report-run callbacks stop on changed Ply, provider or Herdr executable bytes; report Runs have no continuation path | `internal/taskrun/schema.go:350`, `internal/taskrun/herdr_store.go:241` | Prevents an unknown binary from acting as the bound control | Run accepted | ply | Observed lock: `task_run_conflict: executable bytes changed` after an ordinary CLI update (`/Users/perottochristensen/.local/share/ply/planning/report-startup-lock-20261020/affected-run.json`) | observed, source | **implement** Use Delivery's immutable control copy or continuation for report Runs (P2) | Today only AL01 retire, and only for pre-work owners |
| X2 | recovery | product | Unknown or partial earlier effects block dependent work until reconciled with exact evidence | `docs/workflow-execute.md:262-265` | Prevents replay of an uncertain effect | n/a | ply | Blocks only the dependent effect | source | **keep** Truthfully unknown, handled locally | Exact-evidence reconcile (AL01) or retire |
| X3 | recovery | product | Closeout `resume_closeout` and `resume_installation` are shown with actor `human` | `internal/workspaceview/progress.go:173`, `internal/taskrun/read_inventory.go:440` | None. Resuming an authorized operation is coordinator work (`workflow execute follow`) | Agreement | coordinator | 1 misleading human action | source, modeled | **simplify** Label the coordinator (P4) | n/a |
| X4 | recovery | product | A blocked, paused or stopped headless worker sets next actor `human` | `internal/taskrun/worker_unix.go:210` | Control of a headless worker | Pause or stop was chosen by a principal | human | 1 action | source | **keep** Change only if evidence shows a coordinator-resolvable cause | `worker resume` on the same Run |

## Historical question counts and modeled proposals

| Journey | Earlier contract | AL02 base (d335e70e) | Proposed, not measured after | Evidence |
| --- | --- | --- | --- | --- |
| Codex report start with the optional grant in a repo that has `AGENTS.md` | 1 question; stop 20:28:33Z, answer after 35 s, accepted preview after 96 s | 1 | 0 (P1) | observed |
| Same start without the grant | 0 | 0 | 0 | observed |
| Report start to acceptance with unknown profile/policy | Locked, no supported path | 0 (AL01) | 0 | source |
| Delivery start through execute and check-in (this Run) | 1 mandate plus a permission confirmation | 1 mandate | 1 mandate | observed |
| Review rounds | 0 | 0 | 0 | observed |
| Automatic local delivery | 0 | 0 | 0 | source |
| Delivery with required human QA | 1 | 1 | 1 (kept) | source |
| Recovery: report Run after a CLI update while the owner works | Locked | Locked unless pre-work (retire) | 0 (P2) | modeled |
| Recovery: resolve a stop | 1 permission confirmation | 0 | 0 | source |
| Recovery: interrupted closeout shown to the human | 1 | 1 | 0 (P4) | modeled |

## Historical hard-scenario requirements

These scenarios become product tests in the repair Specs, in the existing taskrun
fixtures. They are not AL02 tests.

| Scenario | Where it applies | The test must show |
| --- | --- | --- |
| stale lock/owner | P2: report owner after a CLI update; P4: closeout | The same Run continues, or is retired with history kept; no second writer |
| interrupted agent | P2: owner callback after a provider update mid-Run | Work and session identity are preserved; continuation without replay |
| lost receipt | P1/P2: preview/apply digest reply lost | A repeated start reads the same reservation; the grant is not applied twice |
| drift | P1: trust facts change between preview and apply; P2: executables change after launch | Before launch the grant is withheld and a warning logged; after launch callbacks warn and continue |
| partial effect | P2: tab created, prompt state unknown | The effect stays unknown and only the dependent step is blocked |

## Original repair proposals and subsequent outcomes

### P1 (high): Codex trust stop becomes agent-owned and local

- **Evidence:** the observed question above. `internal/taskrun/herdr_codex_trust.go:60`
  and `internal/taskrun/herdr_codex_trust.go:231`.
- **Scope:** `internal/taskrun/herdr_codex_trust.go` and `internal/taskrun/herdr_service.go:184-188`,
  the trust section of `docs/workflow-run.md:131-145`, and the start help text in
  `cmd/workflow_run.go:43`.
  - If the observed config already trusts the root, the grant is redundant: the
    preview reports it as a no-op.
  - Otherwise, an unconfirmable grant removes only the trust effect from the
    preview, logs a logrus `warn` (path, reason, handling), and continues the
    ordinary start. The stderr pattern is `cmd/workflow_status.go:46`.
  - Malformed grants keep exit 2. No change to persistent config, sandbox or
    runtime approvals.
- **Stop basis:** none remains for the whole start. Explicit `untrusted` still
  withholds the grant, and Codex native policy applies.
- **Regression journeys:**
  - A fixture linked worktree with a tracked `AGENTS.md` and a trusted root:
    preview with the grant gets a confirmation, zero human-directed text, and a
    report-mode start.
  - The same fixture with an untrusted root: start without the grant effect,
    plus a warning.
  - Explicit `untrusted`.
  - Lost apply reply.
  - The code/test forbidden effects (commit, push, install) are unchanged.
- **Done when:** the plyapp journey needs 0 human questions; no product message
  routes an agent-observable fact to the human; all existing trust tests still pass.
- **Depends on:** none. Before removing `AGENTS.md` from the block list, check the
  documented Codex behaviour for it.

### P1 outcome (AL03)

| Journey | Before (observed, AL02) | After (modeled: fixture and candidate CLI) |
| --- | --- | --- |
| Grant plus `AGENTS.md`, root already trusted | Preview stopped with exit 4; 1 human question; 96 s to a usable preview | Preview reports `redundant`. Apply makes 1 tab, 1 agent start without `trust_level`, 1 prompt. 0 questions |
| Grant plus `AGENTS.md`, root not trusted; explicit `untrusted`; unknown level | Preview stopped with exit 4 | Preview reports `withheld` and logs a `warn` on stderr. The ordinary start continues under native policy. 0 questions |
| Trust drift after reservation, before launch | Reserved start stopped mid-launch | Grant withheld and recorded; later callbacks are not blocked by trust |

The plyapp journey was not repeated, since this Task may not restart the
planner's Run. The "after" column comes from the regression fixtures and
the candidate CLI journey, not from a live provider. Cost: one more read of the
Codex config files during preview. Benefit: the human question disappears, and
no trust effect is applied that Ply cannot confirm.

### P2 (high): report Runs survive drift after launch

- **Evidence:** the observed `executable bytes changed` lock
  (`/Users/perottochristensen/.local/share/ply/planning/report-startup-lock-20261020/affected-run.json`),
  plus `internal/taskrun/herdr_store.go:241` and `internal/taskrun/schema.go:350`.
  For trust drift: `internal/taskrun/herdr_store.go:230`.
- **Scope:** reuse the Delivery mechanisms for report Runs: the immutable
  per-attempt control executable and the existing continuation. Check trust
  freshness only before launch; after launch, log a `warn`. No new receipt gate.
- **Stop basis:** keep the stop for an unbound callback executable (identity
  mismatch) and for unknown prior effects.
- **Regression journeys:** CLI install mid-Run; provider update mid-Run; a
  `~/.codex/config.toml` edit after launch; a partial tab/prompt effect;
  legacy frozen Runs keep their contract.
- **Done when:** a working report owner completes its round after each kind of
  drift, with history and effect identity preserved.
- **Depends on:** P1 for the trust part; it touches the same file.

### P2 outcome (AL04)

| Control | What it protected | Now |
| --- | --- | --- |
| Provider executable bytes on every report operation | That the launch used the bound provider | Checked until the native session is bound. After that, a change or removal is logged with `warn` |
| Ply executable bytes on every report callback | That the bound control runs the callback | The bound path, context, cwd, run, request and session stay required. Changed bytes from an install at that path are logged with `warn` |
| Codex trust facts on every operation | That the grant matched the confirmed facts at launch | Checked before launch, as in AL03. After launch, drift of a confirmed grant is logged with `warn`. A contradictory saved binding stays an error |

This uses no new control copy or continuation for report Runs: report authority
is code, test and report only, and these operations launch nothing. Delivery,
headless workers, Herdr, frozen inputs, policy evidence, stops and retirement are
unchanged. `workflow execute` commands send warnings to stderr, as `workflow run`
already does, so stdout keeps a single JSON value.

| Journey | Before (observed source/fixture) | After (fixture and candidate CLI) |
| --- | --- | --- |
| Report owner after a Ply install at the bound path | Locked: `executable bytes changed` | Same round reported and reviewed; 0 questions; `warn` |
| Provider changed or removed, Codex config and `.git/config` edited after launch | Locked | Same round reported; no new tab or provider start |
| Lost report reply | n/a (locked) | Same report record; no second dispatch |
| Prompt reply lost, then install | Show reported drift | Prompt state stays unknown and is not resent; the owner checks in and reports |

The fixture Ply and provider bytes are synthetic. The CLI journey starts with a
real build of the AL03 base (`f5535ef`) as the installed Ply and installs the
candidate at the same path. No live plyapp Run was repeated.

## AL05: file requirements for ordinary inputs

Source of the friction: in AL02, verify stopped on a 0644 `acceptance.sh` with
the generic error `expected bounded private regular file`. The owner first
suspected the 0755 candidate binary. In AL04, a progress report was refused
because a shell-redirected evidence log was 0644. A single check in
`internal/taskrun/schema.go` (`readFile`) mixed type, size and private mode, and
its error named neither the path nor the role.

| Role | Source (base `7e049ed`) | Value of owner-only mode | Decision |
| --- | --- | --- | --- |
| Delivery report (`execute report --file`) | `workflowDeliveryReport` used `private=true` | None. Ply checks the content, copies the bytes into its private event store, and the report grants no authority | **remove**. Group/other write access stays refused |
| Round report (`workflow run report --file`) | `WorkflowReportRound` used `private=true` (review r001-f01) | None. Its run, request, session and round control are checked, and it is canonicalized and kept in Ply's immutable report slot. Later reads use that copy | **remove**, same as the delivery report |
| Coordinator review (`workflow run review --file`) | `WorkflowReviewRun` | It is the coordinator's control decision for the round, not owner evidence | **keep** owner-only, named `coordinator review` |
| Report evidence and verifier evidence | `validateDeliveryReport` used `workflowBound` (private). Also re-read by historical qualification | None. The owner's digest binds the bytes. They are not copied, and every read re-reads the original path against that digest | **remove**, both at publication and on re-read |
| Acceptance entrypoint | `deliveryReadBinding` at reservation, `workflowBound` after the run, in live revalidation, receipt reuse and status | None. The command runs a private snapshot, and the bytes are re-checked by hash | **remove** |
| Candidate review (`--review`) | Same readers as the entrypoint | None. Its candidate, reviewer and findings are validated, and its digest is bound | **remove** |
| Candidate binary | `prepareAutomaticVerification` already used `private=false` and a 0700 snapshot | n/a | **keep** (baseline, proved by a test) |
| Runtime acceptance, notification route | `WorkflowAccept`, `deliveryBegin` | Carry authority claims or a notification authorization | **keep**, with the role named in the error |
| Native run state, stores, controls | `readValue`, `PrivateRead`, `workflowBound` | Ply creates them 0600. Any other mode means a change outside Ply | **keep**. Run state is named; other callers report `private file` |

After the change, a refusal reads like this example: `delivery report <path>: must not be writable by group
or other users, who could change it at any time before Ply reads it; observed regular file, mode 0664,
412 bytes; remove that write access (chmod go-w) and submit the same path
again`. Accepted shared-readable inputs log `warn` with role, path and mode,
once per callback, after Ply has consumed the input (report published or
verification reserved).

Mode bits do not authenticate an owner (review r001-f03). Before AL05, owner-only
mode incidentally excluded files owned by another non-root user, because those
were unreadable at 0600. AL05 does not add a UID check. The caller chooses which
path to submit, and a readable file owned by someone else is consumed like any
other ordinary input. A 0644 mode is never treated as proof of who wrote the
file. What actually binds the input is the digest of the bytes Ply read and
content validation. Reports and the acceptance entrypoint are also copied
privately (event store, report slot, executed attempt snapshot). Report
evidence and the review are not copied at submission and are re-read from
their original path against the digest; the review's bytes become a candidate
artifact when a candidate is created. Group/other write access is refused
because it lets others change the file at any time. Missing or unopenable
ordinary inputs name their role and path and keep the actual errno
(r001-f02). Native readers keep raw errors, so `os.IsNotExist` recovery is
unchanged.

Still protected: absolute physical paths without symlinks, regular files, size
limits, inode replacement while the file is opened or read, the owner's digest
for evidence, and re-reading by digest after binding. Executed bytes are always
the bound snapshot.

**Known limit, not changed:** `readFile` detects only inode replacement
(`os.SameFile`) and growth past the limit. An in-place rewrite of the same inode
with the same size during the read is not detected there. No read seam exists
to test it. What actually protects the ordinary roles today: evidence must match
the digest the owner submitted. For the acceptance entrypoint and review, the
bytes that were read are the bytes bound, snapshotted and executed, and any
later difference fails the post-run digest check. The fault-boundary test
`TestAL05InPlaceOverwriteAfterBindingDoesNotQualify` proves that last part
after binding, not during the read. The remaining work is in the backlog as
`planning/koe/ply-filbinding-samtidig-overskriving-under-lesing.md`: a read
seam in `readFileAs`, a test that rewrites the same inode with the same size
during the read, and a decision on whether a changed mtime or size refuses.

Not changed by AL05: the observe-stop/resolve-stop evidence and runtime
evidence paths. They keep
owner-only mode, and their errors now name the path, requirement and observed
mode under the generic role `private file`.

| Journey | Before | After (fixture and candidate CLI) |
| --- | --- | --- |
| 0644 report with a 0644 evidence log | exit 2, generic error, then manual chmod | Published; one `warn` per file; 0 questions |
| 0644 `acceptance.sh`, 0644 review, 0755 binary | Verify refused; owner debugged the binary | Qualified and delivered through `workflow delivery execute` |
| Lost report reply | n/a | Same event; no second report |
| 0664/0666 input, directory, symlink, oversized, changed evidence | Generic error | Refused; the error names the role, path, requirement and observation; no event |

### P3 (medium, delivered technical change and actual pilot): Claude `auto` for report Runs

- **Evidence:** `internal/taskrun/herdr_transport.go:335` and
  `internal/taskrun/schema.go:251`, while Delivery accepts `auto`
  (`internal/taskrun/delivery_service.go:115`).
- **Scope:** accept `manual` or `auto` in report requests and forward the
  requested mode; update the help text and docs. Existing digests are unchanged.
- **Stop basis:** unknown modes are still rejected.
- **Regression journeys:** an `auto` report start launches with
  `--permission-mode auto`; a `manual` request is byte-identical to today;
  frozen requests are unchanged.
- **Done when:** a Claude report Run has no Ply-forced per-tool prompts.
- **Depends on:** after P1, because both edit `cmd/workflow_run.go:43` and
  `docs/workflow-run.md`.

### P3 outcome (AL06B and AL06C2)

AL06B (`1040987`) accepts and forwards frozen `manual`/`auto` report choices;
unsupported modes fail before effects. The original AL06 Claude owner committed
a dry-run candidate, then owner continuation and the separate review hit an
actual session quota. That Run was retired, not delivered; its dry-run evidence
was preserved. Fresh Codex AL06B ownership did not inherit its acceptance.
AL06B's fixture/CLI tests alone were not an actual Claude pilot.

The first actual pilot, AL06C1 (`wfr_476ec245…`), encountered the coordinator's
incorrect recipient Spec after actual reads. Preserve the rejected reports and
sourced stop: provider ended, source retained, preparation advanced, Run neither
accepted, delivered nor retired. Queue release is not Run closeout.

Revised AL06C2 (`wfr_d94e2170…`) used actual Claude, with UI showing Opus 5.5
and auto mode/classifier for real read actions. The native report was accepted
after two correction rounds; each read verifier exited 0. The coordinator
recorded 0 human mode changes and 0 routine human questions, ended the provider
and closed the preparation choice. Unknown runtime model/profile/policy/sandbox
facts in the receipt stayed unknown. Report acceptance is not TaskResult,
full Task delivery or full Task closeout; the clean source was retained.
`al06c-result.json`, `al06c2-round0-review.json`, `al06c2-accepted.stderr` and
`al06c2-provider-ended.json` preserve the outcome and the earlier acceptance
rejection. Pretty and canonical reports parsed equally but had different byte
hashes; the native predecessor hash was used. P3's technical and actual-provider
checks are satisfied at these distinct evidence levels. Post-start report
closeout remains a separate unresolved finding.

### P4 (low, delivered AL07): resume steps belong to the coordinator

- **Evidence:** `internal/workspaceview/progress.go:173` and
  `internal/taskrun/read_inventory.go:440`.
- **Scope:** pending `resume_closeout` and `resume_installation` show the
  coordinator or agent actor. Human-owned integration keeps the human
  (`internal/taskrun/read_delivery.go:94`).
- **Stop basis:** n/a; this is labelling only.
- **Regression journeys:** an interrupted closeout (stale owner); install unknown.
- **Done when:** status and inventory name no human for resuming an already
  authorized operation.
- **Depends on:** none.

### P4 outcome (AL07)

`6198dd9` corrected startup protocol guidance and coordinator-owned recovery
labels; `7fde35b` repaired the headless regression found by native review1.
The passed exact review2, verification, closeout and installed CLI/history reads
are recorded in `al07-native-final-readback.json`, `al07-closeout-follow.json`
and `al07-installation.json`. Completed/retired historical readback no longer
invents an available reviewer or a human action. Human-owned integration keeps
its human actor; headless recovery keeps the recipient. Unknown installation
means observe first, not reinstall. These read paths do not perform effects.

### P5 (low, delivered AL08): native installation for the self-hosted Epic

- **Evidence:** the identical copied `install-qualified.py` (inventory row I2),
  and `internal/integration/execute.go:245`.
- **Scope:** an investigation Spec. Can the Epic's install profile replace the
  per-job script, with backup and rollback?
- **Stop basis:** keep `install_unknown` until the installed result is observed.
- **Regression journeys:** an interrupted install; a lost install receipt.
- **Done when:** a root job installs through a native receipt, without a
  script copy.
- **Depends on:** after P1–P4.

AL08 implementation outcome for P5: `workflow execute install-cli` now reads
completed self-hosted Ply delivery and its retained qualified executable, with a
separate attributed coordinator operation after closeout. It reuses the atomic
backup/receipt/recovery engine in `internal/clirelease/install.go`; the native
reader and installed smoke live in `internal/taskrun/delivery_cli_install.go`.
The generic integration install-profile adapter remains available for other
products, but generating another command profile would retain the per-job
installation scripting that I2 identified. Usage is in
[CLI installation](cli-installation.md#coordinator-installation-after-native-self-hosted-delivery).

Implementation regressions use isolated targets: `TestQualifiedCLIInstallationPublicConsumer`
consumes native verify/integrate/closeout and then runs the retained candidate
controller after source removal; `TestQualifiedCLIInstallationInterruptionAndKnownSmokeFailure`
preserves actual failed-process evidence and observes changed bytes without
retrying it. `internal/clirelease/qualified_test.go` covers backup/replace
interruptions, reply recovery without a second copy, changed inputs and guarded
rollback across later operations. Subsequently, separate exact-candidate review2
passed, native verification and
integration/closeout completed, and root installed through native `install-cli`
and ran the installed public-consumer scenario (18.820 s). The recorded receipt
is `al08-installation.json` in the evidence directory above. This completes P5,
not the broader audit.

Native AL08 review round 001 found `r001-f01`: an incoming legacy release
failed to archive the outgoing qualified receipt, so an old qualified retry
could replace that later selection. The permanent regression
`TestQualifiedInstallCannotReplayAfterLegacyRelease` reproduced this before the
fix. The shared history boundary now follows the outgoing qualified operation,
regardless of the incoming install path; explicit legacy rollback/reinstall
remains supported. The finding and failed review remain native history; the
corrected `2941fcc` candidate subsequently passed review2, verification,
closeout and the actual root installation recorded above.
