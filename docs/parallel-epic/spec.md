# TE4B29D43 — Specify concurrent native Tasks on one Epic with serial exact delivery

Implementation specification, 23 October 2026. Native author Task:
`ply-parallel-epic-spec-20261023`; Spec `ply-parallel-epic-spec-v1`.
TE4B29D43 is a presentation label only. Source basis:
`f534ee4d3ec064d54f7db3c74d61813c9e5e6fe0`, tree
`1a2e6350996a9880bf77e39e9097be86a05fd8e1`.
Read [design.md](design.md) for source evidence, alternatives, every singleton
consumer category, version distinctions and the selected membership module.
All proposed commands/formats/types and all PAR-T tests below are future work.
Product implementation and synthetic product tests were not executed here.

## Scope and invariants

Deliver the human-selected complete capability: two distinct chosen Goals on
one existing Project/repo/Epic/parent launch native interactive owners in separate
visible Herdr tabs and isolated Task worktrees, act independently, and return
serially to that same parent with exact qualification after any parent movement.
The target used for this authoring assignment is Project/repo `ply`, Epic
`ply-agentic-workflow-support`, `refs/heads/ply_agentic_workflow_support` at
`/Users/perottochristensen/github/ply/ply_agentic_workflow_support`. Later feature
fixtures use isolated equivalents; they must not mutate this live queue.

The implementation owner may make routine technical choices within the selected
model. New product scope or foreign-owner mutation is not authorized by this
Spec. No new Epic workaround, private `.ply` patch, generic scheduler/eventbus,
background watcher, server/PWA operation, global permission change, main/master,
remote push, PR or release effect belongs to the feature delivery. Source
reconciliation and serial exact return are required, not deferred after merely
allowing a second start.

The first complete slice supports fresh-worktree, explicit named interactive
Goal execution with local Epic delivery. Headless workers, week control, legacy
TaskRun start, adopted-source/backup starts and implicit `--next` remain exclusive.
This is an explicit compatibility restriction, not a claim those paths become
parallel. All their reads and recovery paths must safely understand new
membership or reject before effects. Other local branch/PR delivery modes
retain existing contracts; broadening them is outside this slice.

Maintain these invariants through every transition:

| ID | Invariant |
| --- | --- |
| PAR-I01 | One authoritative participation fold over native queue history; pending, active membership and terminal disposition have no competing writable projection. |
| PAR-I02 | Each active Task has unique preparation/operation/worktree/ref/Run/control/provider-session and its own candidate/reviewer/acceptance. Common Epic identity is shared; effect ownership is not. |
| PAR-I03 | Serial admission is exclusive. Parallel admission requires an explicit named Goal and compatible active parallel members; no implicit adoption or retroactive parallel permission for a serial owner. |
| PAR-I04 | Registration, ranking, reserve, close, retire and cleanup for A preserve B's membership and all unrelated pending references/order, including `ply-consumer-agent-status` and `ply-registered-cross-repo-plan`. |
| PAR-I05 | A stop, necessary question or retained source affects its Run. Independent lawful source work on B continues. Shared-target unknown effects still block dependent base/integration effects. |
| PAR-I06 | Delivery requires a clean exact candidate, independent native review, executed exact verification and the frozen acceptance/human gates for the current effective parent generation. No old proof follows changed source/basis. |
| PAR-I07 | Integration to a shared Git common-dir/ref is serial and CAS-bound. Durable attempt identity spans Git, queue closure, base publication and Run result. Retry observes/resumes; it never integrates or closes twice. |
| PAR-I08 | Original Task/Run/session/Spec/start basis, human judgments, refusals and qualified receipts remain immutable history. Reconciliation adds evidence, not replacement authority. |
| PAR-I09 | Installation preserves actual persisted format, including format 4. Workspace format activation is a separate explicit qualified effect and cannot disrupt incompatible foreign frozen consumers. |

## Observable contract

Command sketches in this section are **proposed**, not current support. Extend
existing command help and native JSON validators with the least new choices:

| Trigger | Proposed input | Observable result and next action |
| --- | --- | --- |
| Choose parallel work | `ply workflow execute --spec <exact-spec-id> --parallel --check`, then existing execute/apply surface with the exact plan | Versioned plan freezes named Task/Goal, mode, source reconciliation policy and same target. Explicit `--parallel` cannot combine with `--next`, omitted spec, headless, restart/replacement or unsupported delivery mode. No provider starts on check. |
| Default serial execution | Existing `execute`/`--next`, legacy `prepare` and worker/week routes | Existing behavior on serial history; new checks inspect full active membership. They refuse another participant rather than adopting or ignoring it. Serial first owner cannot be retroactively converted by a second `--parallel` call. |
| Concurrent admission | Explicit parallel A and B; same source base, distinct named Goals | Each freezes `admission_mode=parallel_explicit`, `source_reconciliation=merge_parent_v1` and its own effect scope before launch. A remaining serial member blocks B with its ID and supported close/continue action. |
| Same Goal twice | Identical intent/confirmation/acceptance path | Return the existing preparation/Run/transport reservation, even after pending reordering. Changed bytes or conflicting owner claims fail with no new owner. Unknown start outcome returns observation guidance. |
| Inspect participation | Proposed `queue list --schema-version 2` and updated workflow status/readback | All active members, sorted deterministically, with exact identity, mode, effective parent, freshness, reasons and next actor; pending ranks retain order. Schema 1 succeeds only when representable; multi-active schema 1 explicitly refuses and points to schema 2. |
| Mutate pending | Existing native register/set with expected revision | Add/reorder only pending. Active and terminal Tasks cannot be inserted. CAS conflict returns no effect; retry retains current unrelated entries. |
| Review/verify A | Existing native separate-review/verify route through new versioned bindings | Only A's preparation/basis/candidate/Run qualifies. Reviewer is a distinct native session in its own pane; B's controls, reports or proof cannot substitute. No handwritten `--review` for feature delivery. |
| A's parent moved | Source/basis inspection or integration check | `reconciliation_required` names old/current base and source, preserves A's prior results, and returns owner-operated `reconcile-source` preview argv. A remains active. |
| A stopped or waiting | Existing receipt/report/stop callbacks | The actual refusal/question remains visible for A. B remains independently eligible unless a concrete shared resource or effect is involved. No clearing a stop by receipt, queue release or alternate actor. |
| Close/retire A | Existing exact owner/coordinator path, ID/revision and current real evidence | One terminal transition for A. B and protected pending persist. Retire/release is not integrated/delivered; unknown writer absence retains its diagnosis. |

Keep execution metadata locks short: do not hold a workspace registry lock while
an owner implements, a reviewer works or verification runs. Existing per-Run
locks and effect reservations handle those journeys. Use Project/work-item
locks for atomic publication and the existing target integration authority
ledger for durable serialization. `ApplyTaskIntegration` currently holds
Project -> WorkItem locks through its short Git effect; the unresolved target
reservation is keyed by common-dir/ref, not a separate generic scheduler lock.
Retain that real mutual exclusion and add effective-basis/source/parent CAS
checks at this seam. Git ref locks and exact post-effect/reflog observation
remain necessary; raw external Git is not coordinated by a Ply lock. Preserve
Run -> Project -> WorkItem acquisition order where locks nest; never acquire
a Run lock from inside a workspace store callback. Revalidate cross-store
bindings at each publication.
Concurrent independent native metadata operations either commit at a new CAS
revision or return a no-effect conflict; none replace each other's snapshot.

An unresolved source creation/reconciliation affects its Task and claimed
source path/ref. The participation module must classify operation dependencies:
shared common-dir administration is serialized briefly; an unknown operation
whose possible resource scope overlaps a new action blocks that action. A
known isolated source conflict alone does not block B's lawful source work.
Unresolved parent integration/base publication locks that target, including
other Epics sharing its ref. Keep repository identity based on common-dir;
identically named refs in different repositories are independent.

Coordinator follow-up uses existing native Run IDs and `execute show/follow`
for each owner, with explicit return/next-actor facts. It requests independent
review per candidate and resumes each Run's own operation. A poll/read is an
observation, not a background worker. No second mutable multi-Run status store,
new scheduler or hidden watchdog. Existing durable states remain the truth.

## Source reconciliation protocol

Implement proposed `ply workflow execute reconcile-source RUN_ID` with the
existing bound owner context and `--check` / `--apply --confirm <digest>`
conventions. Proposed `--continue` and `--abort` refer to the same recorded
operation. The operation is owner-controlled and source-only: it never merges
into the Epic or changes the frozen delivery target. Reuse the selected owner,
provider/session and exact control; new workflow schema 3 explicitly authorizes
`merge_parent_v1` at initial launch. Historical schema-1/2 Runs cannot acquire
that authority, and cannot use this command to upgrade themselves.

1. **Preview exact inputs.** Require active named membership, accepted owner,
   unchanged Goal/Spec/mandate, no active stop, no unknown prior delivery effect,
   clean committed source with original start ancestry, registered source
   identity, and clean current registered/observed parent. Parent must descend
   from the preceding effective base. Bind Task/preparation/Run/request,
   prior reconciliation digest, before OID/tree, original start, current base
   revision/OID/tree, target/common-dir and merge policy into confirmation.
   Refuse dirty/uncommitted work with action to commit through existing owner
   authority. Seal/close pending review before apply; do not infer reviewer exit
   or invent a result. Bind the review state so a concurrent new review conflicts.
2. **Reserve before Git.** Under native locks, re-read the selected member,
   source/parent identities, owner and stops, then persist a write-once
   reconciliation intent and per-source operation claim. Invalidate eligibility
   for new review/verification/integration on that Run while it is pending.
   A stale preview fails before effect. Reserve the exact parent commit, not
   whichever ref tip happens to exist when Git executes.
3. **Merge through native adapter.** Merge that exact parent into the existing
   source using the adapter's controlled Git environment (hooks and signing
   cannot execute arbitrary configuration), with expected source ref/OID and
   worktree preconditions. Preserve a deterministic operation marker and
   before/after observation. Do not hold a registry lock through manual conflict
   resolution. No rebase/reset/force move, new Task or editable parent binding.
4. **Handle conflict without pretending success.** Record merge exit, MERGE_HEAD,
   index/worktree conflict paths and known ref state under this operation.
   Return `source_conflict` with the exact operation and next owner action to
   resolve files in this source, then invoke native `--continue`. Continue
   validates operation ownership, parent commit, original source first-parent
   ancestry, resolved index and resulting clean merge commit; native adapter
   performs the final commit and records it. An unresolved semantic issue also
   requires owner corrections and review. `--abort` is offered only for this
   known pending native merge, preserving edits/evidence and refusing an unsafe
   rollback; it must observe the exact pre-state before recording no effect.
   Unknown abort/effect stays blocked. Other participants keep their own work.
5. **Publish effective basis.** Append receipt with old/new base and source,
   generation, effect proof and original owner/Spec linkage, then make its
   committed generation visible through the participation/basis module.
   Preserve the original Task worktree parent, preparation, frozen Spec and
   prior candidate records; validators resolve the new effective parent only
   for new versioned Runs. Current source must include both original source
   and selected parent. If already contained, record the no-Git-effect receipt
   explicitly; a changed base generation still invalidates prior qualification.
6. **Requalify exactly.** Owner completes semantic changes and commits. Coordinator
   requests a new separate native review on exact OID/tree/effective-basis digest;
   owner builds new private bytes from that clean commit and native verify
   executes bounded acceptance without handwritten `--review`. New candidate,
   verification receipt and TaskResult all bind reconciliation generation.
   Review/verify success is technical qualification, not a human answer. Human
   gate requires a new actual candidate-bound answer where selected. A prior
   human fail/blocked remains a veto until the supported real judgment resolves
   it; automatic mode cannot launder it by reconciling.
7. **Return only on current base.** Acquire existing native Project/work-item target lock and reservation/CAS,
   revalidate current parent == this effective parent, current clean candidate,
   review/verification/QA/stops/owner, then execute the native ff-only return.
   If parent moved again, no new integration effect: return to step 1 with the
   new base. Preserve any prior reserved integration attempt and reconcile it
   before a new candidate or reconciliation operation is allowed.

Concrete two-Task journey:

```text
P0 ── A1          A starts at P0, reviews/verifies A1, returns P0 -> A1
 └── B1          B starts at P0, retains separate owner and source B1
      └── B2     native reconcile merges A1 into B1; receipt effective_base=A1
                 B2 has both B1 and A1 ancestry; fresh review + verify on B2
                 serial native return A1 -> B2; close only B; base becomes B2
```

B's old pass on B1/P0 is historical after A's return; it cannot authorize B2/A1.
Even if Git reports an automatic merge, semantic conflicts remain possible.
If A1 and B1 conflict, B stays active in source-conflict state with preserved
intent. A's completed delivery and unrelated pending stay intact. If C advances
the parent during B review, B requires another receipt/generation and review.
No command claims conflict-free integration from path-disjoint diffs alone.

## Publication and crash recovery

Native operations need explicit fault seams and a retry result at each row.
Continue existing write-once keys/digests and one-attempt effect ledgers; ordinary
no-effect contention does not justify a fresh identity after an unknown effect.
A new terminal event is never inferred from an empty queue or missing source.

| Publication boundary | Durable fact / recovery observation | Allowed retry behavior |
| --- | --- | --- |
| Goal execution request, execution Spec, assessment, selection, queue binding | Existing content publications keyed by request; selected Goal remains exact while unrelated pending may move. | Resume missing publication using same input/intent; preserve latest unrelated entries. Before any reservation, a changed selected Goal conflicts. |
| Preparation + worktree operation + reservation | One native registry transaction contains all three; worktree create effect may not yet have happened. | Recover exact preparation and observe ref/worktree operation. Same Task never receives another owner because B reserved later. |
| Worktree/ref creation before ready outcome | Existing operation evidence/observation identifies no effect, prepared or partial/unknown. | Record known prepared outcome once; known no-effect can follow supported retry; unknown requires observation without another create. |
| Run request/control reservation before Herdr/session/start prompt | Existing per-Run launch state records each effect. | Return existing Run and actual transport. Unknown pane/session/prompt outcome is not replayed. No actor swap. |
| Reconciliation intent before merge | Before/parent OIDs, source claim and operation marker retained. | If exact no-effect state observed, resume the same reserved operation only. |
| Merge commit or conflict before receipt | Ref, MERGE_HEAD/index, first-parent/parent ancestry and operation-specific evidence distinguish known success, conflict and unknown. | Adopt observed exact result into same receipt; do not execute another merge. If marker/ref evidence cannot disambiguate, return `reconciliation_effect_unknown`. |
| Reconciliation receipt before effective-basis publication | Append-only receipt links generation and exact pre-state. | Idempotently publish this generation if selected predecessor still matches; no parent/Spec mutation. A concurrent newer generation conflicts. |
| Review/verify start before result | Existing sealed review/acceptance attempts plus new basis generation bound at reservation. | Recover result when completed; unknown execution blocks replay. Source/generation drift marks proof stale. |
| Integration authority/intent/attempt before Git | Existing one-attempt native ledger binds result and expected parent plus effective-basis receipt. | Recheck and resume exact attempt; no replacement Delivery/Run. |
| Target ref moved before integration result | Repository lock, exact source/ref observations and attempt-specific reflog prove before/after. | Record success once if exact effect proven; no second merge. Unknown/foreign movement remains unresolved shared-target gate. |
| Integration result before queue advance | Named durable integration result is authoritative. | Close only its preparation with fresh queue CAS, preserving B and pending. Recover from exact advance event, never queue emptiness. |
| Queue advance before Epic-base publication | Same integration reservation and registered old/new base versions. | Complete the same base publication once; other integration cannot pass pending target closure. Source work on B may continue. |
| Epic-base publication before Run completion | Observe exact integration, named advance and base result. | Seal completion once even if another Task later advances parent; recover historical own effect rather than require parent still equals own result. |
| Run completion before retained evidence/resource cleanup | Existing closeout receipt/archive/evidence set, actual Task writer checks and source identity. | Resume only A closeout. `owner_handoff` exit 5 is delivered with coordinator cleanup after owner leaves; no source writes/retry. Unknown absence preserves resources and diagnosis. |

Durable integration reservation key remains canonical Git common-dir plus
parent ref; native Project/work-item locks serialize publication/effect entry.
The current tests explicitly keep another ref in the same repository
independent after an unresolved reservation on the first ref. Preserve this
behavior; do not introduce a repository-wide unresolved-effect veto. Alias
worktrees of the same repository/ref cannot bypass the reservation. A second native integration must re-observe the target after the
first one's complete queue/base closure. No long-lived registry lock, implicit
background resumer or unlocked metadata patch may fill publication gaps.

## Synthetic regression plan

These are **planned synthetic tests**, not executed results or claims of real
two-agent capability. Use isolated temporary workspaces and real Git repositories
for ancestry, clean states, merge conflicts, ref CAS and recovery. Use existing
fake clock, injected Git fault hooks, fake Herdr/provider/OS adapters for transport
and writer observations. Simulated session/liveness must be labeled in evidence.
Never use active root/PWA/server/CRAP owners, live queue or real denied `/bin/ps`
journeys as fixtures. Start with a failing regression through the public/native
seam for each behavior before its production implementation.

| ID | Public/native scenario and decisive assertions | Owning seam / dependency |
| --- | --- | --- |
| PAR-T01 | Serial format-4 fixture starts A, rejects B; installed new binary serial check/read/register/execute preserve format 4 and pending. Repeat supported legacy formats with appropriate existing rules. | `taskexecute.Execute`, registry codec; before admission change. |
| PAR-T02 | Named `--parallel` A/B at same P0 both start; distinct preparations, worktrees/refs, Runs, fake visible Herdr panes, sessions and controls. Mutating A context for B fails. | Execute -> native prepare -> fake Herdr; after schema activation fixture. |
| PAR-T03 | Concurrent identical A previews/applies/retries create exactly one preparation/operation/Run/launch. Different exact input conflicts. Crash after each launch reservation returns unknown/readback without duplicate prompt/start. | Start guard, launch intent, transport recovery. |
| PAR-T04 | Concurrent A/B previews at qN: one wins CAS, other stale unused preview has no effect, refreshed named B succeeds. Concurrent planner insertion/reorder and recovery preserve selected exact reference and all unrelated order. | Queue publication/Project+WorkItem locks. |
| PAR-T05 | Seed protected pending `ply-consumer-agent-status`, `ply-registered-cross-repo-plan` and unrelated active C on another target. Register/set/reserve/advance/release A; assert full foreign records, references and relative order unchanged. | Native queue fold and transition validators. |
| PAR-T06 | A receipt/working/stopped/question and B progress remain distinct. A known refusal blocks its dependent effects; B's source work remains allowed. Shared-target unknown integration blocks dependent return for both. | Bound report/stop and dependency classification. |
| PAR-T07 | Separate native review per A/B, wrong reviewer session/context rejected; wrong candidate or reconciliation generation cannot inherit pass. Same source OID with changed base generation needs fresh round and acceptance. | Review/verify/evidence readers. |
| PAR-T08 | Real Git A1 and B1 from P0; return A, attempt B with old proof fails before target effect. Native reconcile creates B2 including A1, fresh exact review/verify then ff-only return. Assert two unique integrations, base versions and named closures; pending unchanged. | End-to-end native local return, no raw-Git receipt substitute. |
| PAR-T09 | Real text conflict leaves B's source active with recorded paths/merge intent and explicit continue/abort actions. Resolve through native continue, then semantic correction/review. Abort preserves unexpected edits or refuses; A/unrelated Tasks remain untouched. | Reconciliation adapter and receipt generation. |
| PAR-T10 | Inject C parent movement between reconcile preview/apply, during review, and just before integrate. Stale apply has no effect; completed old-base receipt becomes stale; integrate requires new reconcile/review/verify. | Source CAS and target CAS. |
| PAR-T11 | Fail/crash after content Spec/assessment/selection/set, atomic preparation reservation, Git worktree/ref create, ready publication, Run reservation. Recover exact identity each time without lost B/pending or duplicate create/start. | Existing native publication faults plus new parallel interleavings. |
| PAR-T12 | Fail/crash after reconcile intent, merge success/conflict, receipt and basis publication. Recover success/conflict from operation evidence; deliberately ambiguous evidence stays unknown without rerun. | Reconciliation fault seams. |
| PAR-T13 | Two competing local returns with aliases of same repository/ref: one effect at a time; crash after target effect/result/advance/base/Run publication and repeat. Exactly one own integration, queue terminal and base version; distinct repository same ref remains independent. | Repository lock, reflog/ledger, exact delivery closure. |
| PAR-T14 | Real human fail/blocked/required answer vetoes the affected candidate after reconciliation; automatic mode does not create QA or clear veto. Retire/release/retained-source A with known absence preserves B; unknown absence retains A; cleanup owner_handoff doesn't prove other readers absent. | Acceptance, retire, closeout; fake OS observations labeled. |
| PAR-T15 | New binary reads 4–7, old serial histories byte-preserved; direct explicit 4→8 activation preview has no writes, stale digest/unknown or incompatible foreign frozen control blocks apply. Supported isolated fixture activates exactly once with boundary/backup and no intermediate 7. Installation alone never activates. | Format boundary, release capabilities, finite native consumer guard. |
| PAR-T16 | Queue schema 2 and registered/workflow view show A/B accurately; schema 1 multi-active refuses. Every worker/week/legacy/backup/adoption/next entry and recover/resume path refuses parallel mode before effects, even when singleton projection is absent. Historical Run/control/candidate fixtures retain old meaning. | Consumer inventory, versioned CLI/help fixtures. |

Extend rather than discard existing tests inspected in source:
[task_goal_execution_test.go](../../internal/workspace/task_goal_execution_test.go)
(`TestGoalExecutionRecoversEachPublishedContentStep`,
`TestGoalExecutionSamePublicInputRetriesWithoutDuplicateEffects`,
`TestDeliveryCandidateAcceptsOwnedDescendantButNotDirtyOrMovedParent`);
[task_preparation_fault_test.go](../../internal/workspace/task_preparation_fault_test.go)
(`TestTaskPreparationP5FaultBoundariesAndRecovery`);
[task_queue_removed_source_test.go](../../internal/workspace/task_queue_removed_source_test.go)
(`TestTaskQueueRemovedSourceConcurrentReleaseAndStart`);
[task_integration_repository_lock_test.go](../../internal/workspace/task_integration_repository_lock_test.go)
(`TestTaskIntegrationSameRepositoryLockFollowsTargetRef`);
[delivery_review_round_test.go](../../internal/taskrun/delivery_review_round_test.go)
(`TestFS01ChangedCandidateCannotInheritPass`). Retain serial expectations and add
explicit parallel variants; merely changing old assertions would lose coverage.

Proposed new focused files: `internal/workspace/task_participation_test.go`,
`task_source_reconciliation_test.go` and `registry_parallel_activation_test.go`;
`internal/taskexecute/parallel_epic_test.go`;
`internal/taskrun/parallel_epic_return_test.go`;
`cmd/workflow_execute_parallel_test.go`. Names are planned, not existing sources.
Tests must assert native durable facts/observed Git effects and authority
rejection, not private helper call counts or mirrored implementation branches.

Run named package/test selectors for these bounded fixtures after each phase;
owner records exact commands/output, simulated facts and limitations. Full
coverage remains **UNKNOWN/BLOCKED**. No mechanical CC/CRAP threshold or new
full heavy suite is a gate until runtime prerequisites are separately resolved.
An installed CLI smoke on an isolated fixture proves distribution, schemas,
Git composition and callback routing; fake Herdr does not prove two real owners.
A later explicitly authorized real two-agent same-Epic journey must observe both
native sessions/receipts/worktrees simultaneously, separate review/verify,
moved-base reconciliation and both serial return receipts. Do not claim this
from source tests, installation or idle panes alone.

## Implementation order and delivery

Keep one implementation owner for the shared membership/basis contracts, with
separate native reviewer. These phases are bounded implementation order, not
new Tasks registered by this author. Shared queue/schema files must be changed
sequentially. Each phase requires its named red regressions before production
changes. The capability is complete only after phase 4 passes; partial support
must remain disabled for parallel launch.

| Phase | Files/contracts owned | Entry / exit and dependencies |
| --- | --- | --- |
| 1. Membership and versioned contracts | workspace queue types/store/preparation/start guard, proposed participation module; registry codec/compatibility and clirelease capabilities; taskexecute plan/input; public queue/view types. | Freeze event/schema/effective-basis contract from design. Red PAR-T01/04/05/15/16; green isolated codec/fold/CAS/compatibility checks. New admission remains gated until downstream phases qualify. No live format write. |
| 2. Explicit native launch and consumers | task_goal_execution/goal_operation, taskexecute launch/store, workflowhandoff bridge, taskrun report/recovery/review/retire/worker plus week and workspaceview inventory, cmd help. | Depends on phase 1 types; PAR-T02/03/06/11/16 through native seams. Two isolated fake-Herdr owners with separate identity. Legacy exclusions enforced throughout recovery, not only at CLI parsing. |
| 3. Source reconciliation and exact basis | proposed workspace reconciliation module/records + Git adapter; task_goal_candidate/task_spec_guard; new workflow/handoff validators, candidate/review/verify/authority and readback. | Depends on participation and ownership. PAR-T07/09/10/12/14. Native merge/continue/abort and fresh generation proofs; original records remain unchanged. No launch-only feature considered delivered. |
| 4. Serial exact return and activation guard | workspace integration/base/closeout, workflowhandoff delivery_return, taskrun delivery authority/completion, bounded CLI fixtures; explicit activation plan and known native consumer inventory. | Depends on all above. PAR-T08/13/15 and all invariants. Real-Git two-return fixture, fault recovery, old-format serial use and incompatible-reader refusal qualify complete source capability. |
| 5. Review, exact verify, source delivery; optional later install/activation | Same clean Task candidate and private executable; native review/verify/local Epic return. Coordinator owns installation and any real runtime smoke/activation under their distinct authority. | Full focused evidence, independent review pass, exact native verify without handwritten --review, one authorized local Epic integration. Installation preserves actual format; activation requires exact separate native plan and compatible consumers. |

Source changes may be reviewed in phases, but no parallel mutation is enabled
until reconciliation/qualification/return/consumer safety are complete. A future
implementation Spec can use these file contracts and criterion IDs directly;
no new product approval is needed for the already selected feature. If a
missing compatibility inventory seam proves broader than the finite native
guard described in design, preserve source-delivered capability and report
exact unavailable evidence/scope before expanding infrastructure. Do not act on
foreign frozen owners to manufacture an activation pass.

For the implementation delivery, use the frozen native source/target agreement,
clean committed candidate, meaningful focused acceptance and actual independent
native review. Coordinator starts reviewer; owner corrects findings, commits and
holds for a fresh round. Build private candidate via `make build` with embedded
commit/dirty identity, inspect `version --format json` and `go version -m`, then
native verify **without handwritten --review** on the exact reviewed candidate.
Keep qualified bytes immutable; observe local return, named queue closure and
Epic-base receipt. No human judgment is inferred from technical pass. One
Task's delivery never declares the whole Epic complete.

Coordinator may later install only the exact qualified feature binary with its
old backup/history preserved and run a relevant installed smoke under the
separate selected installation mandate. Source-delivered, installed,
capability-qualified, live-format-activated and actually-two-agent-observed are
separate facts. This docs-only author Task has **no installation authority**.
Live activation cannot be inferred from a feature binary's supported formats.

For this author Task, run only the frozen acceptance entrypoint's bounded
document/reference/consistency/whitespace checks, then commit only
`docs/parallel-epic/spec.md` and `docs/parallel-epic/design.md`. Build/inspect
private metadata only because the frozen automatic native agreement requires
it. Preserve all raw evidence outside temporary cleanup. HOLD the clean OID/tree
for coordinator `backlog-coordinator w9:pE` to start separate native review;
no internal agents, direct Herdr start, follow-on Task registration or Slack.
Parent owns follow/closeout/total notification. After sealed review pass, same
owner performs native verify without `--review` and the single authorized local
Epic return. `owner_handoff` exit 5 is normal delivery completion: stop source
writes and return to coordinator, with no integration retry. Known process
absence limitations remain limitations during cleanup.

## Requirement traceability

The author criteria are document completeness claims, not executed product
behavior. Native task-requirements evidence is produced by verification of this
exact docs candidate, not a second hand-authored status artifact.

| Frozen criterion | Resolved contract | Future observable evidence |
| --- | --- | --- |
| ac-001 | Design compares event-fold membership with preparation-minus-terminal derivation and selects event fold; inventories source seams and separates actual format 4, maximum supported 7, proposed 8 and public schemas. | Author docs/reference acceptance; future PAR-T01/04/15/16. |
| ac-002 | PAR-I01–09; complete launch/identity/pending/stop contract; source reconciliation steps 1–7; named crash/publication recovery; exact review/verify and serial return. | Future PAR-T02–14 plus consumer tests, raw native receipts and labeled runtime limits. |
| ac-003 | Exactly the two prescribed docs; bounded docs acceptance and diff --check; no production changes, product/fullsuite/process/browser probes, installation or live registry migration. | Author clean commit diff, preserved frozen-entrypoint output, private build/version identity and separate review/native verification. |

The future implementation is acceptable only when every PAR-I invariant and all
sixteen planned scenarios have meaningful native-seam evidence or a concrete
unresolved blocker; a blocked required behavior cannot be called complete.
Actual live activation and real two-agent observation are reported separately
when foreign-reader/runtime prerequisites prevent them. Protected pending,
old human gates and foreign history must remain intact at every boundary.
