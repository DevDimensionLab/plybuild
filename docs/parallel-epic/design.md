# TE4B29D43 — Concurrent Tasks on one Epic: design

Design proposal, 23 October 2026. Native Task: `ply-parallel-epic-spec-20261023`.
Presentation code TE4B29D43 changes no Task, Run, provider session or authority.
The implementable contract is [spec.md](spec.md). This authoring candidate changes
only these two documents. Product implementation and synthetic product tests were
not executed. Every new command, type and format below is **proposed**.

## Source basis

Inspected source: commit `f534ee4d3ec064d54f7db3c74d61813c9e5e6fe0`, tree
`1a2e6350996a9880bf77e39e9097be86a05fd8e1`. Parent:
`refs/heads/ply_agentic_workflow_support`, registered checkout
`/Users/perottochristensen/github/ply/ply_agentic_workflow_support`,
Project/repo `ply`, Epic `ply-agentic-workflow-support`.

Keep three version facts separate:

| Fact | Evidence and implication |
| --- | --- |
| Actual shared persisted format **4** | Coordinator's native readbacks around queue revisions 508/509, supplied during this Task; not a new independent live-store inspection here. Activation must start from the actual fresh format/digest. |
| Highest source-supported format **7** | [registry_compatibility.go](../../internal/workspace/registry_compatibility.go), `CurrentRegistryFormat = 7`; support for a newer format does not establish use of it. |
| Public queue schemas | [task_queue_types.go](../../internal/workspace/task_queue_types.go), `WorkspaceTaskQueueReadback`, and [task_queue_registered.go](../../internal/workspace/task_queue_registered.go) have their own kind/schema contracts. They are not registry format numbers. |

The frozen execution bound base revision 114 and queue revision 509. These are
observed start facts, not enduring mutation preconditions. Queue revisions must
always be freshly compared. No live queue is a fixture for this design.

[Earlier modularity analysis](../robust-codebase/modularity-analysis.md) and its
[structured companion](../robust-codebase/modularity-candidates.json) identify
useful seams and the existing single-current restriction. They are historical
analysis, not another Task-status store or an execution approval. Full coverage
remains **UNKNOWN/BLOCKED**. Preserve the known `/bin/ps` denial; this authoring
Task runs no process, browser, permission or fullsuite probes.

The present restriction is a contract: `applyQueueEvent` stores one current ID;
`validateQueueRegistry` rejects a second reservation, validates closures against
that ID, and requires every preparation to have its reservation. Removing the
preview guard cannot create supported concurrent execution.

## Alternatives and decision

Both alternatives retain one append-only native queue journal, existing Task
preparations and per-Run ownership. Neither adds a scheduler, a side queue, a
sentinel current, another Epic or an independent status database.

| Alternative | Representation and interface | Cost and decision |
| --- | --- | --- |
| A: explicit membership transitions in the existing fold | `reserve_parallel` adds one preparation ID to a fold's active set; `advance`/`release` remove only their named ID. Serial `reserve` remains exclusive. `set` changes pending only. The active set is derived, never separately persisted. | **Selected.** Membership, ordering, admission mode and closure have one CAS chain. Existing reservation/operation publication remains the transaction seam. Smallest complete extension of the current journal contract. |
| B: infer membership from all preparations minus terminal events | Scan preparations and close history to derive all nonterminal preparations. Admission mode must still be recorded in each preparation, and its queue reservation must still be validated for atomic worktree creation. | Viable with format changes, but “not terminal” cannot distinguish an unreserved/partial publication from authorized membership. A safe implementation ends up reconstructing A's reservation fold anyway. More scattered recovery reasoning for no reduced public interface. |

Both change persisted semantics. Format 4 through 7 readers reject concurrent
reservation history, not merely an unfamiliar JSON field. Relabeling the new
semantics as format 4 or 7 would break strict validation. Choose **proposed
format 8**, behind the explicit activation gate below. This number is the next
producer version in the inspected source, not the observed live format.

A separate target-local journal would require an authority handover that old
readers understand before ignoring the old queue. The existing source has no
such handover. Keeping a fake legacy current while using a second authoritative
queue would introduce two truths and a bypass. That alternative is excluded.

The cohesive responsibility is **Task participation and effective delivery
basis**: who may act independently, which source generation its proof covers,
and when that proof may return to the shared parent. Membership without
supported stale-parent reconciliation is an incomplete feature.

## Membership module

Put a concrete module in proposed `internal/workspace/task_participation.go`,
using existing registry records and canonical digests. Its interface includes
these operations (names illustrative, not existing exports):

```go
snapshot := foldParticipation(registry, queueID) // pending, active, terminal, revision
plan, err := planAdmission(snapshot, exactGoal, mode, resourceClaims)
err = commitParticipation(session, plan, expectedRevision)
basis, err := effectiveDeliveryBasis(registry, preparationID, reconciliationID)
```

`commitParticipation` runs inside existing Project/work-item locks and publishes
through `WorkItemStoreSession.Publish`. It validates the expected revision,
exact Task/Goal/selection, unique preparation/operation/source claims, and named
transition. It does not start providers, run Git or infer liveness. Concrete
functions suffice; no new Go interface for a single implementation. Existing
Git/OS/Herdr adapters remain the seams for real variation and fault injection.

Active membership is a set keyed by preparation ID, ordered for readback by
reservation revision then ID. A Task has at most one active preparation, and
its old terminal preparation remains terminal. Pending order is orthogonal:
reserve removes only its exact Task; close removes only its exact active ID;
set may replace pending only and cannot requeue any active or terminal Task.
`preparationDisposition` consults this same fold. New multi-active callers use
`active`; legacy singleton projections may keep `current` as specified below.
No caller uses a projected `current == nil` as an admission test.

Dependency direction:
`cmd` -> `taskexecute` / `taskrun` -> `workflowhandoff` -> `workspace`
participation/basis -> existing store locks, immutable content and Git adapters.
Readback in `workspaceview`/`taskjournal` consumes native projections. Workspace
never imports taskrun to discover owners or start reviewers. Run authorization
is passed as an already validated native binding and checked again at effects.

Admission and mutation remain separate. A stale unconsumed preview fails CAS
with no effect and names a fresh preview. Once an intent/reservation exists,
retry recovers that exact identity before ranking anything else. Unrelated
queue edits can be rebased into a new publication of the same intent only after
rechecking its exact selected entry; they cannot redirect an intent to a new
Goal. An unknown start effect requires observation of its reservation, never a
second provider start.

## Singleton inventory

This inventory covers direct queue `.Current` reads and indirect “current,
prepared, fresh” assumptions in production source at the inspected base. Line
numbers are deliberately omitted so the filenames and symbols remain usable.
Other fields named `Current` for candidate freshness, PR reconsideration or a
Task-journal snapshot are not queue membership; preserve their meaning.

| Journey | Existing source and assumption | Required treatment |
| --- | --- | --- |
| Persist, replay, validate, disposition | [task_queue_store.go](../../internal/workspace/task_queue_store.go): `queueFold`, `applyQueueEvent`, `foldQueue`, `preparationDisposition`, `validateQueueRegistry`, `validateQueueTransition`; `set` excludes current, reserve requires empty, close clears current; base-update plan validates a single current reference. | One derived active set; replay old serial histories with unchanged validity and digest. New parallel events only in format 8. Named close never empties all members. Validate every preparation-to-reservation link and append-only transition. |
| Goal preview and queue binding | [task_goal_execution.go](../../internal/workspace/task_goal_execution.go): `buildGoalExecutionPreviewInTarget`, `bindGoalExecutionQueue`, `PrepareTaskGoalExecution`, stored/resume helpers. Another current blocks before selection and during binding. | Explicit parallel admission only for named Goal; exact existing execution recovery precedes selection. Merge only the selected entry's binding into latest pending CAS; preserve unrelated active/pending. Serial default still requires exclusive participation. |
| Recovery candidate discovery | [goal_operation.go](../../internal/workspace/goal_operation.go): `goalExecutionCandidates` scans pending and one current plus preserved acceptance-path history. | Include all active preparations in target, but scope explicit recovery to exact Task/intent; never adopt another active owner. Keep historical acceptance-path recovery. |
| Preparation and initial start | [task_preparation.go](../../internal/workspace/task_preparation.go): `buildPreparePlan`, reservation closure in `PrepareTask`, `preparationReadback`, `TaskRunPreparation`; [task_start_guard.go](../../internal/workspace/task_start_guard.go): `PreviewInitialTaskGoalExecution`, `WithTaskStartGuard`. | Same-Task resume remains exact; legacy next retains exclusivity. Parallel plan records admission mode and commits preparation/reservation/worktree intent atomically. Start guard checks selected member and relevant claims, not global emptiness. |
| Bound initial start | [taskexecute/bound_start.go](../../internal/taskexecute/bound_start.go): `buildBoundBasis`, `approvedStartGuard.Check`, `ApplyBoundStart`; [goal_operation.go](../../internal/workspace/goal_operation.go): `TaskStartupDependencySHA256` retains shared lifecycle/queue/base records, omitting only foreign external content registrations. | Preserve this dependency protection. Parallel contention requires fresh exact no-effect preview or same-intent recovery, not removing shared ownership from the digest. Update bound request schema/input for explicit mode; never treat a registry digest mismatch as permission to relaunch. |
| Task-scoped preservation | [task_content_scope.go](../../internal/workspace/task_content_scope.go): `outsideTaskContentScope`, `validateTaskScopeTransition` preserve foreign records and keep queue set outside Task scope. | Add reconciliation records to this preservation projection and classify new events by exact preparation. Cross-Task admission/set remains a dedicated locked operation; effective-basis updates cannot mutate foreign records. |
| Execute launch and start CAS | [taskexecute/service.go](../../internal/taskexecute/service.go): `Execute`, `launchIntent`, `startRequest`; [taskrun/herdr_service.go](../../internal/taskrun/herdr_service.go), [herdr_store.go](../../internal/taskrun/herdr_store.go); [delivery_start_recovery.go](../../internal/taskrun/delivery_start_recovery.go). | Freeze parallel/reconciliation authority before launch. Keep per-Run reservation, provider/session, control and retry identities. Metadata contention is a no-effect refresh, unknown transport is not retry eligibility. |
| Legacy TaskRun start / handoff | [task_run_bridge.go](../../internal/workflowhandoff/task_run_bridge.go): preparation must be current, prepared, fresh; [taskrun/service.go](../../internal/taskrun/service.go): native start/collection. | Keep old TaskRun route exclusive; new interactive parallel delivery goes through a versioned request and membership-aware bridge. Do not silently broaden legacy return authority. |
| Report / receipt / stop | [report_startup.go](../../internal/taskrun/report_startup.go): `workflowReportPreparation`, `workflowReleaseReportQueue`; [delivery_service.go](../../internal/taskrun/delivery_service.go), [delivery_stop.go](../../internal/taskrun/delivery_stop.go), [delivery_runtime_observation.go](../../internal/taskrun/delivery_runtime_observation.go). | Report/receipt remain bound to exact preparation hash/owner; release tests membership of that ID. A stop or question in A affects A; shared-target unknown effects still block dependent target effects. Receipt never clears refusal. |
| Review | [delivery_review_round.go](../../internal/taskrun/delivery_review_round.go): `WorkflowDeliveryRequestReview` loads `TaskRunPreparation`, candidate/session binding and exact return checks. | Read A's membership and effective basis, preserve unique reviewer session/control per Run. Candidate movement invalidates its pass; B's membership cannot select A's reviewer. |
| Verify / candidate / authority | [task_goal_candidate.go](../../internal/workspace/task_goal_candidate.go): `ValidateDeliveryCandidateTarget` requires frozen parent still current; [task_spec_guard.go](../../internal/workspace/task_spec_guard.go): basis/parent/delivery digest and relevance; [delivery_authority.go](../../internal/taskrun/delivery_authority.go): `deliveryAuthorityState`, native candidate authorization; [delivery_callback.go](../../internal/taskrun/delivery_callback.go) / [delivery_automatic_verification.go](../../internal/taskrun/delivery_automatic_verification.go). | Validate immutable original basis plus a native effective-basis receipt for new Runs; never weaken legacy parent checks. Bind new review, verification, candidate and TaskResult to one reconciliation generation. Keep human veto and stop checks. |
| Local integrate / queue / base | [delivery_return.go](../../internal/workflowhandoff/delivery_return.go): `IntegrateDeliveryCandidate`, `CloseDeliveryTaskQueue`; [task_integration.go](../../internal/workspace/task_integration.go): `buildIntegrationPlan`, `ApplyTaskIntegration`, persisted recovery; [delivery_completion.go](../../internal/taskrun/delivery_completion.go), [delivery_authority.go](../../internal/taskrun/delivery_authority.go). | Exact member closure, effective parent, same registered target. Preserve native locked target reservation/CAS and one-attempt ledger. Recover native result before re-evaluating a newly moved parent. Reservation spans Git, queue, base and Run publications. |
| Epic lookup / base update / unresolved effects | [epic_base.go](../../internal/workspace/epic_base.go): `queueHasGoalSpec`, `queueUnresolved`, `buildBasePlan`, `UpdateEpicBase`; plan `CurrentPreparationID`. | Inspect all members; proposed base plan schema 2 binds sorted active IDs. Distinguish source-local unresolved operations from shared-target unresolved effects, with no broad ignore of unknown effects. Base publication never rewrites other Tasks' start basis. |
| Owner retire / retain | [delivery_retire.go](../../internal/taskrun/delivery_retire.go): `deliveryRetireReleaseQueue`; [delivery_retain_effects.go](../../internal/taskrun/delivery_retain_effects.go): exact current check before release. | Release A by ID, preserve B, retain A's real absence/refusal/effect evidence. Unknown writer absence remains unknown; no queue membership change proves process exit. |
| Backup / replacement | [task_backup_goal.go](../../internal/workspace/task_backup_goal.go): selected-retirement/current preview; [taskexecute/backup_goal.go](../../internal/taskexecute/backup_goal.go): foreign-current guard and queued detection; [taskrun/backup_goal.go](../../internal/taskrun/backup_goal.go), [backup_retirement.go](../../internal/taskrun/backup_retirement.go). | Use exact source/replacement membership; no arbitrary current adoption. Keep replacement route serial/exclusive in first implementation. Retirement is not delivery and cannot inherit acceptance. |
| Adoption | [task_adoption.go](../../internal/workspace/task_adoption.go): `observeAdoption` refuses current activity and unresolved effects. | Keep adoption exclusive; don't make existing-source adoption a second parallel launch path. Native fresh-worktree parallel delivery is the minimum complete slice. |
| Closeout / resource retirement | [task_closeout.go](../../internal/workspace/task_closeout.go), [task_closeout_reservation.go](../../internal/workspace/task_closeout_reservation.go), [task_closeout_retention.go](../../internal/workspace/task_closeout_retention.go), [delivery_automatic_closeout.go](../../internal/taskrun/delivery_automatic_closeout.go). | These rely on named result/integration/ownership, not direct q.Current. Verify their upstream membership assumptions; remove only A's source after evidence retention and actual writer rules. Preserve archive and owner_handoff. |
| Native queue readback / text | [task_queue.go](../../internal/workspace/task_queue.go): `buildQueueReadback`, `SetTaskQueue`, `closeTaskQueue`, `TaskQueueText`; [task_queue_types.go](../../internal/workspace/task_queue_types.go). | Schema 2 lists all active, per-mode readiness, reasons and IDs; schema 1 must refuse an unrepresentable multi-active view. `--ready` never hides active members. |
| Registered index / public workflow view | [task_queue_registered.go](../../internal/workspace/task_queue_registered.go): `ReadRegisteredTaskQueues`, `indexRegisteredQueues`, `applyContentChanges`; [workspaceview/workflow_status.go](../../internal/workspaceview/workflow_status.go): `addQueues`; [workflow_status_types.go](../../internal/workspaceview/workflow_status_types.go). | Project every active Task, diagnostic and provenance independently. Pending readiness distinguishes serial conflict from parallel eligibility. Update view types/schema and rendering; no first-member truncation. |
| Trace / Task journal / recorded results | [taskrun/read_trace_native.go](../../internal/taskrun/read_trace_native.go), [read_trace_workflow.go](../../internal/taskrun/read_trace_workflow.go), [taskjournal/native.go](../../internal/taskjournal/native.go), [workspace/task_journal.go](../../internal/workspace/task_journal.go), [task_recorded_progress.go](../../internal/workspace/task_recorded_progress.go). | Native validators and current-result relevance must understand membership and effective basis. Historical events remain exact; journal snapshot `Current` means current facts of one Task, not sole Epic owner. No new persisted status projection. |
| Headless worker | [worker_queue_release.go](../../internal/taskrun/worker_queue_release.go): `workerPreparationCurrent`, `WithParkedWorker`; [worker_pipeline.go](../../internal/taskrun/worker_pipeline.go), [worker_retry.go](../../internal/taskrun/worker_retry.go). | Remain exclusive on this target, including pause/resume/retry. Preserve the exact advanced integration-recovery exception. Reject parallel selection before provider effects; never treat no singleton projection as a free queue. |
| Week orchestration | [weekrun/manager.go](../../internal/weekrun/manager.go): Start, initial-state recovery, next selection; [control.go](../../internal/weekrun/control.go): registration recovery/dispatch; [native.go](../../internal/weekrun/native.go): parked-worker release. | Remain exclusive, check full membership at all seven q.Current sites and native release seam. No concurrent week, worker or `--next` expansion. |
| CLI/help and external readers | [workflow_execute.go](../../cmd/workflow_execute.go), [workspace_task_prepare.go](../../cmd/workspace_task_prepare.go), [workspace_task_queue.go](../../cmd/workspace_task_queue.go), [workflow_execute_review.go](../../cmd/workflow_execute_review.go). | Add explicit parallel selection/status diagnostics, describe serial exclusions and required reconciliation; retain existing command meaning. Audit JSON consumers via schema tests before activation. |

The inventory is an implementation checklist, not a request to mechanically
rename every `Current`. Re-run a bounded source search at implementation base,
classify additional hits and review consumers before enabling admission.

## Persisted and public compatibility

**Install without activation.** New reader/writer supports existing formats
1–7 unchanged, including this workspace's actual format 4. Normal serial
execute/register/report/read operations must preserve the current format and
all existing boundaries. No automatic 4→7 or 4→8 conversion, including on
`--check`. Existing `registry_upgrade` behavior is retained for older commands.

**Minimal producer boundary.** Extend the native explicit registry upgrade seam
in [registry_compatibility.go](../../internal/workspace/registry_compatibility.go)
and manifest checks in [clirelease/manifest.go](../../internal/clirelease/manifest.go).
A proposed format-8 activation plan binds actual from-version (4, 5, 6 or 7),
registry SHA256, selected target, release package/manifest/executable digests and
reader capabilities. Direct 4→8 preserves all old records and pre-existing
format/adoption/Epic-create boundaries; it must not fabricate intermediate
upgrade history. Add one distinct immutable parallel-format boundary, rather
than overwrite earlier boundaries. Format-8 shape/validation must select the
predecessor capabilities from this boundary: a 4→8 registry need not have
format-5/6/7 feature boundaries or feature records, whereas a 7→8 registry
retains and validates all three. Existing `validateRegistryBoundary` uses
`>=5/6/7` tests and must be changed explicitly for this case; simply raising
`CurrentRegistryFormat` is insufficient. Extend
[work_item_store.go](../../internal/workspace/work_item_store.go) YAML shape,
wire encoding and transition checks for this predecessor-dependent shape.
Empty new reconciliation arrays are explicit.
Older formats 1–3 use their existing supported preparation/content path first;
parallel activation from them is excluded from the first complete slice.

The present `RegistryFormatBoundary` checks availability of a new reader, not
compatibility of every existing frozen control. Existing explicit upgrade
implementation specializes 4→5; it does **not** already support this plan.
Extending it requires tests and explicit authority, not an assumed flag.

**Known-consumer guard.** Extend the finite persisted projection in
[read_inventory.go](../../internal/taskrun/read_inventory.go), `ReadInventory`,
which reads native and Herdr attempts with explicit incomplete/damaged evidence
and performs no process probe. Its current fields do not certify executable
format compatibility; add that specific observation from exact retained
control/release artifacts, plus native week/bound-start bindings.
Registry format is workspace-wide. A target-local
feature cannot promise format transparency to foreign frozen readers. The
activation preview must read native persisted Run/control/worker/week bindings
that can still mutate or resume against this registry and compare their exact
frozen executable capability to format 8. Use a native read-only compatibility
projection with IDs, digests, authority state and reasons, not OS process scans.
An unretired paused/stopped Run with a frozen older reader is a blocker even
when no process is running. Unknown/missing executable capability or incomplete
native inventory is a blocker, not absence. Repeat this check under the
activation's native locks and bind the inventory digest to its confirmation.
Concurrent new starts/registrations must share that lock or conflict by digest.

The minimum migration rule is to finish such a consumer through its existing
authorized route, or wait; unchanged compatible controls can continue. There
is no supported in-place rebinding of an incompatible frozen control in this
slice. Do not retire another Task, replace its binary, rewrite its mandate,
replay denied actions, or modify a stopped PWA owner. Human authorization for
this feature does not grant those actions. An actual incompatible foreign
consumer therefore prevents **live activation**, while isolated qualification,
source delivery and installation preserving format 4 can still complete.
Historical terminal evidence stays byte-identical and readable with the new
reader; old executable readback may explicitly report unsupported format after
activation. Source cleanup proves none of these reader-compatibility facts.

This is a finite extension to native format publication and control inventory,
not a global liveness registry, process watcher, migration service or general
scheduler. If native bindings cannot provide a complete compatibility guard,
return `activation_consumer_compatibility_unknown` with exact missing facts.
Deliver the qualified source/capability and propose the narrowly missing native
inventory seam as a separate scope decision; do not enable a weaker guard or
build a general migration framework. No foreign PWA action is permitted.

**Public contracts.** Proposed queue readback schema 2 has `active[]` (Task,
preparation, mode, reservation revision, source path/ref, per-member freshness,
effective parent/reconciliation, reasons/evidence), existing ranked `pending[]`,
and explicit serial/parallel readiness. Internal consumers use that projection.
Schema 1 remains available for zero/one active member with its exact singleton
shape and legacy `current` disposition; this is a read-only projection, not
an exclusive-write permission. Parallel native callers use schema 2 even
when only one member exists. A schema-1 request with more than one active
member returns a structured incompatibility error and next action to request
schema 2. It must not return `current:null`
with success. Text output prints each active member. Registered-queue and
workflow-status envelopes need their own schema bumps for arrays; candidate
freshness booleans named `current` retain their independent meaning. The new
registered-queue and workflow-status versions both use schema 2; schema-1
callers receive the same explicit incompatibility rule when multiple members
cannot be represented. Per-preparation readback schema 2 reports `active`
for new plans. Old preparation plans retain their schema-1 readback and
`current` disposition; requesting schema 1 for a new schema-2 preparation
explicitly refuses instead of relabeling its nested plan.

**New execution contracts.** Proposed execute/preparation plans schema 2 freeze
`admission_mode=parallel_explicit` and `source_reconciliation=merge_parent_v1`.
Legacy plans stay schema 1. Parallel workflow request/Run envelopes use schema
3 (current delivery is schema 2); their handoff uses schema 4 (current task-spec
handoff is schema 3). Candidate schema 4 binds a reconciliation reference and
effective parent. Workspace Goal/Spec revision schema 2 remains unchanged:
its immutable original implementation basis is not edited. Extend new workflow
and handoff validators only for explicitly versioned requests and authority;
old receipt/begin contracts and schemas remain untouched. Versioned callback
contexts and immutable controls must support all new records before launch.

## Parent reconciliation design

Today `ValidateDeliveryCandidateTarget` and `validateSpecImplementationBasis`
require the original registered Task parent, while `buildIntegrationPlan`
requires Task parent == current Epic base == observed clean parent. The handoff
bridge also checks original preparation parent against owner start evidence.
Ancestry alone is insufficient. [delivery_return.go](../../internal/workflowhandoff/delivery_return.go)
explicitly calls refreshing a later feature. `UpdateEpicBase` changes registered
Epic metadata; it neither merges source nor changes frozen Task authority.

Choose owner-controlled **merge of the current parent into the existing Task
source**, preserving original ancestry, original preparation, Task, Run,
session and frozen Goal/Spec. No rebase or replacement Run is needed. The
proposed native `reconcile-source` operation owns this Git effect; raw Git or
hand-edited preparations cannot establish its receipt. A clean Git merge still
needs semantic review and new tests.

Format 8 adds append-only reconciliation records and a validated effective-basis
projection. Generation zero is the original frozen parent/start basis with a
null reconciliation reference, so the first participant can qualify without
performing a redundant merge. Every later generation requires its exact
committed native receipt and predecessor chain. Each record binds Task/preparation/Run/request/mandate, preceding
reconciliation, exact source before/after OID/tree, source ref/worktree/common
dir, original start, target ref/worktree, old/new registered base revision and
OID/tree, merge intent/result and evidence. Original worktree binding and
preparation plan remain immutable. Candidate validation uses original identity
plus this receipt's effective parent; it does not overwrite `Task.Worktree.ParentOID`
or relax old `validateSpecImplementationBasis` checks globally. New native
validators, integration intents and evidence readers must agree on the same
projection; a receipt is not an unchecked override parameter.

The owner first commits its current work, then the coordinator seals or closes
any pending review of that source without forging a pass. A reconciliation
intent blocks new candidate/review/verify/integration effects for that Run until
it resolves. Known stops, unknown prior effects and released owners retain
existing gates. No process scan is needed to prove that membership alone does
not grant these effects; use the existing actual owner/writer checks where an
effect requires them and preserve unknown observations.

The [spec protocol](spec.md#source-reconciliation-protocol) defines exact CAS,
conflict completion and crash observation. On success, any earlier candidate,
review, verification and QA stay historical and cannot qualify the new source.
Human fail/blocked judgment remains a veto requiring actual resolution; automatic
acceptance cannot erase it. Every effective-base change needs a new independent
review and exact native verification, including an already-contained parent
whose merge produces the same source OID. Bind the effective-basis generation,
not only OID/tree, in qualification keys.

This makes repeated parent movement ordinary: if C integrates while B is being
reconciled/reviewed, B receives a concrete new stale-base diagnosis and repeats
against C's observed base. No integration-order scheduler or fairness promise
is added. A text/semantic conflict in B does not stop A's independent source
work, but an unresolved shared-target integration remains a target gate.
