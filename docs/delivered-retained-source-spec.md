# AL48: Delivered work with retained source

Decision-ready implementation Spec. This documentation deliverable implements no
product transition and performs no MP-01 action. The primary preserved AL48 Spec
is normative; the consumer finding supplies evidence, not producer authority.
Source inspection is against `60e99247283472083577ee8e1985141a7d499cd0`.

## User outcome and bounded decision

Select one explicit, metadata-only completion disposition for an interactive,
automatically accepted local Epic delivery: **delivery completed; source retained;
cleanup obligation open under a named coordinator**. The Run becomes terminal
`completed`, with resource disposition `retained_cleanup_open`. This is distinct
from physical retirement, legacy intentional keep and stopped-owner retirement.
It never means that no OS writer exists or that deletion was attempted.

The first implementation includes recording, recovery and public discovery of
this disposition. Physical cleanup remains unavailable through this new mode;
the later supported transition is specified below as a separate implementation
boundary. A user can record the disposition, inspect the obligation and retain
the source safely. They cannot use ordinary follow as a cleanup retry afterward.
There is no unresolved product choice needed for this first implementation.

Keep existing Task lifecycle/queue facts and workflow history authoritative. Add
no parallel Task status store, general recovery service or historical migration.
Only work depending on this exact resource is constrained. Other Tasks can use
the observed updated Epic base without waiting for this cleanup obligation.

## Concrete evidence and its limits

The preserved consumer finding identifies MP-01 Run
`wfr_b4a514c2e864655ff54c779d346781dd18519ce73c19ef0132ac1361b5fdc2a5`,
Task `plyapp-monitor-scenario-harness`, result
`trs_cf3452c7793f09eb7b79c0ce44e4ac04`, candidate
`c07251d96f6d8164ee44cbd548e31c5165f5686d`. It reports native review3 and
qualified verify00000007, exact local Epic integration, queue closure and Epic
base update; integration returned exit 5 `owner_handoff`. The exact owner
Herdr `w7:p3S/t3Q` then ended and its pane absence was checked. Clean source
remains under root; OS writer liveness remains unknown.

The finding binds an earlier actual denial, `fork/exec /bin/ps: operation not
permitted`, through its cleanup report/native receipt. That historical denial
is not replaced by the later Herdr absence. These are supplied consumer facts,
not a live producer re-observation or proof that every process ended. AL48 reads
no consumer Run state and repeats no follow, ps/lsof, sandbox journey, cleanup,
retirement or integration. No other provider or human shell is a relay.

The preserved primary Spec content digest is
`sha256:091df06e899575aed36dfef4f8130e2e0d7e73e752d45391dc8c87140a0e7a09`;
the consumer-finding digest is
`sha256:903605dc3fc668d842fc4160fd0dc136d32aa5854bf81b15f2e40b3d28565a28`.
Its original locator is
`plyapp/planning/reports/orchestrator-monitor-execution-20261022/native-cleanup-consumer-finding.md`.
Original integration/result names in that finding are evidence locators, not
new permission or assertions that AL48 has independently read those receipts.

## Actual support and gap

Paths below are relative to this repository. Symbols are inspected source
references; tests listed here were read, not executed as live cleanup proof.

| Existing support | Inspected path / symbol | Consequence |
| --- | --- | --- |
| CL02 default cleanup and frozen keep | [delivery_automatic_closeout.go](../internal/taskrun/delivery_automatic_closeout.go): `automaticCloseoutKeeps`, `ResumeAutomaticLocalDelivery`, `automaticIntegratedResult` | Exact observed automatic integration releases ownership. Cleanup-enabled Runs wait for writer absence and complete only with retired source. Older mandates without cleanup intentionally keep it. No operator-selected retained cleanup obligation exists here. |
| Positive writer evidence | [delivery_task_writers.go](../internal/taskrun/delivery_task_writers.go): `taskWriters.cleanup`, `taskWriters.herdr`, `taskWriters.local`; [delivery_writer_liveness.go](../internal/taskrun/delivery_writer_liveness.go): `observeWriterLiveness`, `systemLocalProcesses` | Herdr inventory and local process observation compose. Inspection failure is unknown, never absence. Production uses lsof and /bin/ps; its bounded scope cannot prove arbitrary external writer absence. |
| AL21 bounded observation | [al21_closeout_observation_test.go](../internal/taskrun/al21_closeout_observation_test.go): `TestAL21CloseoutObservation` | Modeled disappearance/transient recovery can finish the same closeout. Persistent/unknown inventory, local writers and unknown processes preserve closing plus delivered. This does not remove deletion guards or cure a real denial. |
| AL40 retain effects | [delivery_retain_effects.go](../internal/taskrun/delivery_retain_effects.go): `deliveryRetentionGuard`, `deliveryRetentionObserve`, `deliveryRetainedEffects`; [delivery_resource_status.go](../internal/taskrun/delivery_resource_status.go): `deliveryResourceStatus` | Requires active stop and rejects completed delivery, ownership release or human integration. Observation still calls cleanup writer checks. Retires stopped resources with retained effects, not already delivered work; no supported overlap for MP-01. |
| Native closeout plan and reservation | [task_closeout.go](../internal/workspace/task_closeout.go): `previewTaskCloseout`, `applyTaskCloseoutGuarded`; [task_closeout_reservation.go](../internal/workspace/task_closeout_reservation.go): `validateCloseoutTargetAvailable` | Exact plan digest, operation and source reservation already exist. Changing Keep conflicts with an existing plan. Current complete receipts require retired or kept resources: a new disposition needs explicit validation, not a forged Keep plan. |
| Evidence retention | [task_closeout_retention.go](../internal/workspace/task_closeout_retention.go): `retainCloseoutEvidence`, `verifyCloseoutRetention`, `ResolveTaskRetainedEvidence`; [task_closeout_automatic_retention.go](../internal/workspace/task_closeout_automatic_retention.go): `automaticCloseoutSnapshots` | Reuse durable evidence retention and automatic verification snapshots. Do not change historical origin paths/digests or delete source before retained evidence verifies. |
| Public projection | [delivery_status.go](../internal/taskrun/delivery_status.go): `deliveryStatus`; [read_closeout.go](../internal/taskrun/read_closeout.go): `workflowTaskCloseout`, `workflowCloseoutAction`; [read_inventory.go](../internal/taskrun/read_inventory.go): `InventoryRun`; [read_delivery.go](../internal/taskrun/read_delivery.go): `readInventoryDelivery` | Final delivery already distinguishes delivered from Run closing. Existing closeout guidance assumes eventual cleanup or keep. Extend these authoritative projections to expose the obligation rather than treating completed as no remaining action. |
| Physical protection | [task_closeout_git.go](../internal/workspace/task_closeout_git.go): `closeoutObserveSource`, `closeoutCheckAncestor`, `closeoutRemoveWorktree`, `closeoutRemoveBranch`; [task_closeout.go](../internal/workspace/task_closeout.go): `closeoutProtectResources` | Preserve source identity, evidence, ancestry, clean linked worktree, protected resources and no-force/CAS removal checks. |

Thus current support gives honest delivery readback and safe retention while
closing, or intentional frozen keep. Neither is the selected terminal outcome
with a discoverable assigned cleanup obligation. Do not repurpose AL40 or edit
its guards to manufacture overlap.

## Proposed public surface (not available today)

One new effect command, run by the actual coordinator from outside the source:

```sh
ply workflow execute retain-source RUN_ID \
  --cleanup-owner 'backlog-coordinator' \
  --reason 'Exact delivery observed; physical cleanup cannot qualify' \
  --format json
```

Use a stable accountable owner label (1..256 codepoints, trimmed, no controls)
and a bounded nonempty reason following native coordinator input conventions.
They identify responsibility and explanation; neither attests permission or
absence. CLI resolves exact candidate, operation, source and evidence from Run
history. No pasted OID/session, new Task registration, force flag, replacement
provider or permission-confirmation argument. Report-Run aliases reject the mode.

Read it through existing `workflow execute show RUN_ID --format json` and
`workflow run list --format json` inventory. Add a `source_disposition` projection
to `delivery_status` and the corresponding inventory delivery row, backed by the
same validated closeout/disposition records. Both expose:

- `state=retained_cleanup_open`, `run_id`, `task_id`, `task_result_id`, candidate
  OID/tree, exact source path/ref/common-dir and closeout operation/plan digest;
- `obligation_id`, `cleanup_owner`, reason, disposition binding and current
  obligation state `open`; retained evidence manifest/archive bindings;
- last recorded writer evidence with timestamp, method/scope and limitation:
  `managed_owner=absent`, `process_liveness=unknown` where applicable;
- historical denial/stop bindings and original partial/unknown effect claims,
  separately from current delivery evidence; no stop-resolution fabrication;
- `next_actor=cleanup_owner`, `next_action.kind=retain_and_inspect`, and text:
  “Delivery is complete. Source cleanup remains open under OWNER. Retain this
  source; physical cleanup requires a separately supported qualified transition.”

Human-readable show gives the same delivery/resource/owner distinction. Inventory
must include these completed Runs under existing inclusive/history selection;
default actionable inventory must retain rows with open obligations even though
the delivery Run is terminal. A broken binding yields an unknown diagnostic,
not an empty obligation or claim of removal. Existing final-delivery evidence
remains the integration receipt. AL40 `resource_retirement` keeps its meaning.

## State, identity and authority

Reuse the native closeout operation; if none exists, reserve it deterministically
for this exact integrated TaskResult. Stable `obligation_id` is derived from Run,
TaskResult and that operation, once. Add one optional disposition record to the
closeout receipt, plus append-only native workflow event
`delivered_source_retained`. It binds request/mandate digest, Run/provider/native
session, Task/worktree/result IDs, candidate OID/tree, integration/release
bindings, operation/plan digest, retained manifest, source observation and all
known denial/uncertainty evidence. Store owner/reason/time in that record.
No invented new session, repeated integration or separate registry is needed.

The closeout receipt gains an explicit validated terminal lifecycle/resource
combination: `state=complete`, `lifecycle_completed=true`,
`resource_state=retained_cleanup_open`, removal booleans false, verified retention
manifest and open obligation. The immutable original plan keeps `Keep=false`.
Do not relax the normal complete invariant: this combination is legal only with
the bound disposition record and exact delivery proof. Task lifecycle completion
reuses the existing lifecycle transition; WorktreeState must remain retained,
never retired. Source reservation survives completion and excludes adoption,
new writers and another closeout over that exact path/ref. It does not reserve
the entire Epic, registry or other Task's source.

Authority is coordinator metadata closeout of an already authorized delivery,
not renewed source writing/deletion authority. The implementation must extend
the native permitted coordinator transition checks explicitly. A new frozen
contract can advertise this disposition separately from cleanup permission;
never infer permission to delete from retention. For already-started interactive
automatic local Epic Runs with original cleanup authority, allow this safer
metadata-only disposition as a versioned supported coordinator transition when
all guards below pass. Bind the unchanged old mandate digest and actual caller;
append the new decision under current runtime review. No rewriting original
Goal, mandate, release, acceptance or receipt bytes. Runs lacking cleanup authority
stay on legacy keep in this first cut. No historic batch conversion or consumer
Run action is authorized by this document.

Preserve actual stop facts even if cleanup-specific stops remain unresolved.
Only an identified cleanup observation/removal blocker may coexist with this
transition. Unknown integration, installation, acceptance or other delivery
effects cannot be relabeled cleanup. Retention closes source-owner callbacks;
it neither resolves an active stop nor permits dependent replay. Existing
provider controls, denied runtime actions and task authority remain unchanged.

## Transition guards and order

1. Under native workflow serialization, require interactive automatic local
   Epic agreement; exact successful qualified candidate, actual review and
   acceptance, observed integration receipt, queue closure and Epic base update,
   and released writing ownership. Reuse `automaticIntegratedResult` binding
   checks. The target may have advanced through legitimate later work; verify
   candidate ancestry and preserved integration evidence, not a perpetual exact
   current-head requirement. Required installation must not be pending/unknown.
2. Reject stopped/replaced owner retirement, non-delivered candidate, PR/remote
   uncertainty, legacy keep, foreign ownership, malformed history, live managed
   owner and known local writer. Require fresh complete Herdr pane/agent absence
   through the read-only managed-inventory adapter. Failed inventory is unknown
   and blocks this first cut. Do not call the local process adapter: retain its
   last unknown/denied state; managed absence is explicitly not OS absence.
   Unresolved prior destructive cleanup reservation/unknown removal effect
   blocks disposition until legitimately observed; it cannot become “retained”.
3. Check physical registered linked source, unchanged direct branch/OID/tree,
   exact clean tracked/untracked/ignored/index/submodule state and protected
   resource boundaries using existing read-only closeout guards. Re-observe
   around retention to detect drift. A known local writer remains a blocker;
   elapsed time or Herdr absence never clears that positive observation.
4. Under the established workflow/work-item lock order, bind/reserve the exact
   closeout operation without changing its plan. Preserve archive and verified
   retention snapshots outside source (including known denial evidence); reject
   missing/corrupt evidence. Add only a disposition-specific coordinator guard:
   no reuse of a deletion guard with unknown rewritten as absent.
5. Persist the bound disposition/obligation decision before lifecycle completion.
   Complete lifecycle and append its workflow event using recoverable native
   event/index publication. Mark delivery/round/final-return `completed` only
   after the receipt, lifecycle and obligation agree. Removal booleans stay false.
   No physical deletion, process signal, pane close, owner start, queue successor
   advancement or second Git integration occurs anywhere in this command.

For an already reserved cleanup plan, preserve operation and plan digest and
append this explicit disposition; never toggle Keep to avoid `plan_changed`.
If there is no reserved plan, create one for the same native closeout identity
before publishing the disposition. Keep new metadata internal to native stores,
using their serialization and existing CAS/publication patterns.

## Idempotence, concurrency and recovery

Same Run, owner and reason return the same obligation and recorded result after
completion, without fresh process inspection or another event. Changed inputs
conflict; owner reassignment is outside this first implementation. Lost replies
are read back using the same command/identity; unknown effects are never retried
as a new operation. A crash after decision/retention, before lifecycle or event
indexing, resumes only metadata reconciliation. Source changes during a pending
recovery block it; no history or observed unknown is overwritten.

Concurrent `follow` and retain-source serialize on the same Run and source
reservation. If follow has already reserved destructive removal, retention
refuses rather than racing or guessing its effect. If retention wins, follow
only reads/completes metadata recovery and returns completed with the obligation;
it must never auto-remove this source. Repeated show/list remain read-only.
A newer Task on the same Epic can advance the queue/base; retention checks its
own historic queue closure and never closes or releases that successor.

No automatic Slack event is created on every read/retry. Reuse native event
identity and existing notification subscription behavior, with disposition and
cleanup owner visible in the completion projection. A completion notice is not
a deletion receipt or human judgment. This AL48 Task sends no independent notice;
parent backlog-coordinator `w9:pE` owns whole-job notifications.

## Later physical cleanup: separate supported boundary

First implementation deliberately leaves `obligation.state=open`; current
`follow`, `integrate`, AL40 retire and generic keep/reconcile must not bypass it.
Users can inspect show/list, keep the source and proceed with independent work.
There is no presently qualified MP-01 deletion action to recommend.

A later implementation must expose an explicit obligation-targeted cleanup
transition for the same Run/operation/obligation, without reopening delivery or
creating a new identity. Its public command must read the obligation's native
bindings rather than accept copied paths/OIDs as authority. It must first inspect
preserved refusals and unknown attempts without invoking the refused observer.
Permission restoration or a new provider/shell name alone is not evidence.

Legitimate enabling evidence is a separately documented, actual runtime change
that permits a new bounded observation, or a separately implemented supported
observation method with positive coverage of the same Task/session writers.
Its provenance, changed action, method, scope, time and actual result must be
recorded under native control and runtime review; the old refusal stays intact.
An operator's assertion of “allowed”, restored-looking Herdr, timeout, absent
pane or a different command spelling cannot satisfy this gate. If no such
supported path is available, obligation stays open indefinitely. This document
authorizes neither building an alternate observer nor trying one on MP-01.

Only after that gate may the explicit cleanup transition obtain fresh positive
writer absence (managed inventory plus bounded OS Task/session coverage), check
source/reservation/ownership/evidence afresh and revalidate immediately before
destructive phases. Known writers, denied or unknown observations block deletion.
Reuse unchanged `closeoutProtectResources`, `closeoutObserveSource`, candidate
archive/evidence/ancestry and no-force worktree removal plus exact-ref CAS branch
removal guards. Dirty/new source, nested repo, ignored files, index flags,
symlink, foreign source/branch, locks or changed candidate cannot be discarded.

Append cleanup reservation and effect observations to the same obligation.
A lost reply after worktree/branch effect requires receipt-aware observation of
that reserved effect before proceeding; no blind repeat, new obligation or force
removal. Only both removals positively observed can produce
`obligation.state=fulfilled`, resource `retired`. Original retained disposition,
denial, review, verification and integration remain available through history
and retained evidence resolution. Delivery stays completed throughout. No Run
migration, broad recovery platform or general authority revision is needed.

## Later regression plan and AL48 qualification

Use injected `Dependencies.LocalProcesses` and modeled Herdr fixtures; never
invoke actual denied helpers, active owners, consumer Runs or provider launch as
test targets. Write a failing focused test before later product changes. Inspect
transitive fixture helpers before running a bounded test selection.

| Scenario | Required assertion |
| --- | --- |
| R01 Exact delivery, released owner, managed absent, modeled OS denial/unknown | Completion with one named open obligation; source/ref and denial bytes preserved; process adapter invocation count zero; show/list discover it. Unknown OS is not reported absent. |
| R02 Guard matrix | Non-delivery, PR unknown, legacy keep, retired/replaced owner, live managed/known local writer, dirty/foreign source, missing evidence, unknown removal or installation all refuse without changing resources/history. |
| R03 Lost reply and crash stages | Same input recovers one obligation, operation and event before/after lifecycle/index publication. Changed inputs conflict. Tampered bound evidence blocks. No duplicate integration or notification. |
| R04 Concurrent follow/retention and successor | Either the guarded cleanup reservation wins or retention wins; never both effects. Retained terminal follow cannot delete; successor queue/base untouched; source reservation prevents new writers/adoption. |
| R05 Later cleanup boundary (future implementation) | Without sourced enabling change, process adapter not invoked. Modeled legitimate new observation with a writer/unknown still blocks. Positive bounded absence plus all original guards permits same-obligation removal; interrupted removal recovers by observation and preserves old denial. |

For this documentation candidate, private acceptance checks the single-file
scope, whitespace, local path/symbol references and contract coverage. Separate
native semantic review must evaluate the whole contract, especially the boundary
between metadata completion and positive deletion evidence. These docs checks
qualify neither future product behavior nor MP-01 cleanup. Build unchanged
product code using `make build PLY_BUILD_OUTPUT=...` on the clean exact docs
commit; preserve source identity, SHA and `go version -m` outside the Task. Parent
starts separate review; owner holds candidate/binary unchanged until full sealed
exact pass, then uses native verify **without --review**, and only an exit-zero
qualified/current/inputs-match result permits one authorized integrate. Delivered
`owner_handoff` exit 5 ends the owner normally; parent owns follow, source closeout,
qualified installation, smoke and notification. No local install in AL48.
