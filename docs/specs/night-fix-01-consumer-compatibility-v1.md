# Bounded native consumer compatibility corrections — v1

T4A5CC805 · Task `ply-nf01a-compatibility-spec-20261024` · NF01A.
English implementation Spec, source-checked on `7404cc7dec9af9ff1afd822f24dc00214c798b36`.
This Task authors this file only; root w9:pE owns the later implementation Task.

## Problem and boundary

The finite compatibility projection rejects some already terminal mandates before
reading their native evidence, and loses a captured headless preparation binding.
Correct those four cases within the existing PAR-A01–A04 contract. Historical
reconciliation identified two retired reports, five completed headless deliveries,
one completed human integration, and MR05's misleading preparation diagnostic.
Those counts describe historical evidence, not current blockers or promised reductions.
The frozen NF01A assignment and batch mandate are normative; historical assessment
and reference tables are contextual. Do not inspect or operate their live Runs.

Keep registry format4 unchanged. No migration, actual activation check/apply,
week-control migration, new writer, stop-clear, rebind, retire, history patch,
cleanup requirement or inferred provider/process termination. Completed week state
alone remains insufficient. No new status database, caller-supplied compatibility
facts, global control change or upgrade/install authority follows from this Spec.

## Current source

| Seam, verified on this base | Existing behavior and bounded change |
| --- | --- |
| [registry_parallel_activation.go](../../internal/taskrun/registry_parallel_activation.go): `parallelActivationConsumers`, `parallelConsumerClass` | Inventory freshness is checked first, then non-Herdr transport is rejected before `workflowRead` and `owned[PreparationID]`. The class function rejects report-only contracts before retirement and uses only `automaticIntegratedResult` for completed local integration. Move terminal fact evaluation and validated preparation binding ahead of transport rejection; retain the fallback gates. |
| [read_inventory.go](../../internal/taskrun/read_inventory.go): `ReadInventory`, `workflowInventoryReturn`; [herdr_store.go](../../internal/taskrun/herdr_store.go): `workflowRead`, `workflowBound`, `workflowReservations` | Freshness verifies preserved request/state and report-slot hashes, report identity and completion review/terminal binding; it is not liveness. `workflowRead` also validates request/root/Run/session and headless identity. Reservations distinguish retired-report release and exact TaskResult qualification from a bare reviewed return. Keep damaged bindings blocking. |
| [report_startup.go](../../internal/taskrun/report_startup.go): `workflowReportRetirement`, `workflowReportQueueReleased` | Typed retirement verifies kind/version, Run, request SHA and native session; bound files pass `workflowBound`. Queue validation resolves the preparation's actual queue in a validated snapshot and requires release/advance naming that exact preparation. Retirement alone does not discharge its queue reservation. Reuse both. |
| [delivery_retire.go](../../internal/taskrun/delivery_retire.go): `deliveryRetirement`; [human_integration.go](../../internal/taskrun/human_integration.go): `deliveryRelease`, `preservedHumanIntegration`, `validateHumanIntegrationDecision` | Delivery retirement and ownership release validate native event/binding identity. Release binds Task/result, OID/tree, source ref/locator, owner/session and provenance. Human authorization additionally binds the actual pass/QA, plan/decision digest, exact release and native decision history. Reuse these readers; call no retirement/release/decision writers. |
| [delivery_automatic_closeout.go](../../internal/taskrun/delivery_automatic_closeout.go): `automaticIntegratedResult`; [read_closeout.go](../../internal/taskrun/read_closeout.go): `workflowTaskCloseout`; [task_closeout.go](../../internal/workspace/task_closeout.go): `ReadTaskCloseoutAt` | Automatic integration validates envelope, event, exact candidate/TaskResult, completion/recovery and parent or historical completion; deliberately excludes human takeover. Closeout validates candidate/result/source identity and store receipt integrity. Preserve these semantics. |
| [delivery_authority.go](../../internal/taskrun/delivery_authority.go): `completeLocalDelivery`; [read_trace_delivery.go](../../internal/taskrun/read_trace_delivery.go): `deliveryIntegration`, `traceWireObject`, `traceWireJSON`; [task_integration.go](../../internal/workspace/task_integration.go) | Human and automatic local delivery write the same `PlyDeliveryIntegration@1` envelope. Trace already binds its native result to registered IntegrationResult/authority/QA. Nested `Result.Integration.Readback.Value` embeds `canonicaljson.Object` as Go members; it is not a flat `integration_result` JSON property. Reuse that decoding/binding pattern for a small read-only human terminal helper, not the automatic-policy predicate or an effectful integration check. |
| [registry_parallel_activation.go](../../internal/workspace/registry_parallel_activation.go): `ParallelConsumerCompatibility`, `ActivateParallelRegistry`; [workspace_registry.go](../../cmd/workspace_registry.go) | Existing ID/class/effects/reasons/digest wire shape and public read-only `--check` are sufficient. No schema or new command is required. |

Production changes should remain in the classifier and, if needed, a small
read-only terminal integration helper beside the existing readers. Extract shared
pure binding checks from trace only if needed; preserve its readback semantics.
Do not rewrite the inventory, queue fold, lifecycle, permissions or automatic
integration policy. Tests extend the fixture seams below.

## Classification rules

1. **Integrity before compatibility.** Retain inventory-level unknown blockers and
   per-row damaged/stale/unknown rejection. Only a fresh, identified workflow Run
   may reach `workflowRead`; non-workflow native journals keep their conservative
   unknown fallback. A read/validation error produces an explicit unknown consumer,
   never disappearance or terminal success. Validate request/Run/preparation/Task
   against the registered preparation, source and queue before recording `owned`.
   Record this binding before the transport gate. A damaged row cannot suppress
   its separate unresolved preparation. Preserve deterministic sorting and IDs.
2. **Read terminal evidence before transport and owner-control compatibility.**
   A helper should distinguish proven terminal, absent/incomplete proof, and
   invalid proof. Terminal history needs no current old executable bytes or UI
   availability. Retain its control SHA as history, not live authority. Include
   the exact evidence read in `EvidenceSHA256` (retirement plus queue release, or
   integration/release/closeout plus the relevant validated registry facts).
   Missing proof continues conservatively; contradictory evidence is unknown.
3. **Retired report terminal proof.** Before the report-only rejection, require
   `workflowReportRetirement` and `workflowReportQueueReleased` for this Run's
   validated preparation. Bind the report retirement's SHA, Run/request/native
   session and preparation→Task/source/queue; a release/advance in another queue,
   for another preparation or before an unrelated reservation is insufficient.
   Use the validated append-only registry and its existing queue reader; do not
   invent a second fold. Validate every referenced immutable report/terminal slot
   through the existing inventory checks; a report may legitimately be absent
   for a pre-start retirement. If proof is complete, emit PAR-A02, compatible=true,
   remaining_effects=[]: this mandate is ended, not successfully delivered.
   A retirement file without exact queue release remains blocked. Retained
   refusals/unknown original effects remain history and authorize no replay.
   An accepted, blocked, pending or sealed report, even with a matching TaskResult,
   is not a new terminal rule: keep unclassified remaining report/review/correction
   or separate QA/integration authority blocked. TaskResult matching handoff,
   start and terminal hashes can release a target reservation; it alone does not
   discharge every effect of that mandate.
4. **Completed local delivery terminal proof.** Reject any unresolved native
   delivery attempt before success. For the latest exact candidate require a
   valid `deliveryRelease`, recorded native integration completion and
   `workflowTaskCloseout` with state=complete and lifecycle_completed=true,
   plus delivery.phase=completed. Use `automaticIntegratedResult` unchanged for
   automatic delivery, including headless. For accepted human takeover use
   `preservedHumanIntegration` and the human integration binding below. Either
   complete route yields PAR-A02 with no remaining effects. Phase/completed,
   a closed queue, accepted=true, a lone result or closeout alone is insufficient.
5. **Human integration binding.** Read candidate.Integration via `workflowBound`;
   require a matching native `integration` event and `PlyDeliveryIntegration@1`
   envelope with this Run/request SHA/candidate key, no error and Result.Completed.
   Decode its actual `DeliveryIntegrationResult` wire form, requiring complete
   recovery and its exact TaskResult, OID/tree and source. Decode the embedded
   canonical object using the existing `traceWireObject`/`traceWireJSON` pattern;
   compare the full IntegrationResult with the validated registry record, then
   bind its AuthorityID to this Task/result/candidate, selected target/ref and
   the frozen delivery authorization (Run/request/mandate/preparation). Require
   successful native outcome (exact_effect/already_integrated, recovery complete),
   the exact accepted human plan/decision/release and registered pass QA; reject
   a later applicable fail/blocked answer or a different latest candidate.
   Bind authority's HumanQARecordID and human authorization to those same records.
   Cross-check closeout.NativeIntegrationID and Plan ownership/target against
   this integration/release/candidate. ParentOID may be the candidate or a valid
   historical completion already established by the native readback; do not
   require today's parent to remain unchanged. Reuse validated registry/native
   authority checks and trace binding logic, never infer success from transport,
   a Boolean, map shape or arbitrary receipt. Reading these historical facts
   grants no human judgment, execution, requalification or cleanup authority.
6. **Fallbacks preserve remaining effects.** With no terminal proof, retain the
   non-Herdr unknown blocker; a captured headless MR05 Run remains unknown with
   its Run/preparation binding reflected in the Run consumer identity/reason.
   Do not emit the false separate preparation/no-complete-binding blocker for
   a correctly bound captured Run. Its week blocker remains; no release or
   clearance follows. SL01's uncaptured/unknown startup keeps its preparation
   blocker. For Herdr keep existing PAR-A01 owner work, PAR-A03 automatic
   cleanup-only with proven qualified coordinator, and PAR-A04 human integration/
   restorable correction. Do not broaden PAR-A03 to human or active headless
   controls. Preserve existing delivery-retirement classification/validator and
   all A01/A04 control-byte compatibility gates; add no new lifecycle claim.

## RED/GREEN acceptance plan

These are planned product tests, **not executed** by this authoring Task. Add
focused failures for each of the four deviations, observe RED for its intended
reason, implement the minimum reader change, then observe GREEN and the strict
negative controls. Fixture records must be generated by native paths or existing
controlled seams, not fabricated historical success/accepted flags. Tampering is
only a negative test, with the expected hashes/slots deliberately controlled.

Reuse [registry_parallel_activation_test.go](../../internal/taskrun/registry_parallel_activation_test.go)
(`TestParallelActivationLegacyCloseoutClasses`),
[registry_parallel_retained_test.go](../../internal/taskrun/registry_parallel_retained_test.go)
(`retainedActivationFixture`, `activationNativeState`, no-lock/write wrappers),
[parallel_epic_return_test.go](../../internal/taskrun/parallel_epic_return_test.go)
(`parallelRuntime`, `parallelQualify`),
[report_startup_test.go](../../internal/taskrun/report_startup_test.go)
(`al01AbsentReportFixture`, retirement durable-retry tests),
[report_retirement_poststart_test.go](../../internal/taskrun/report_retirement_poststart_test.go),
[human_integration_test.go](../../internal/taskrun/human_integration_test.go)
(`humanIntegrationCandidate`, preserved decision/release/closeout tests), and
[worker_delivery_test.go](../../internal/taskrun/worker_delivery_test.go)
(`newWorkerDeliveryFixture`, native review/receipt fixtures). Inject simulated
Herdr/provider, LocalProcesses, StartupProcessTree and SourceWriters at fixture
construction; forbid fallback to actual OS/provider inspection. Existing helper
names do not by themselves guarantee isolation. Use bounded owned testchildren
only if a test needs a real process; wait/reap them. No foreign process probes.

| Case | Required observations |
| --- | --- |
| AT01 — retired reports | Model pre-start and stopped post-start report retirement with exact native queue release/advance: classifier and consumer collection yield PAR-A02/empty effects. RED today is report-only rejection. Retry/history stay immutable. Remove release, change queue/preparation, damage retirement hash/Run/request/session or report-slot/terminal hash: unknown/blocked, never A02. No report required for genuine pre-start case. |
| AT02 — headless terminal | Generate headless automatic native integration, recorded attempt, exact ownership release and complete closeout. RED today is early transport rejection; GREEN is A02. Missing/partial/unknown integration, wrong event kind/Run/request, wrong candidate/result, changed bytes, unreleased owner, missing/mismatched closeout, cleanup_pending or phase-only completed stays blocked. Live/working/unknown headless stays unknown; no generic headless-reader compatibility claim. |
| AT03 — human terminal | Use native human release→preserved pass/plan→local integration→complete closeout, including legacy explicit takeover. RED today is A04 because automatic reader excludes it; GREEN is A02. Reject accepted-only, decision-only, missing release, other result/OID/tree/source/plan/target, unregistered integration, authority/QA mismatch, negative/superseding answer, unresolved attempt and incomplete closeout. Test actual nested canonical-object wire form and valid historical parent completion. |
| AT04 — MR05 diagnostic | Fresh headless request binds registered preparation/Task/source; queue reserve still active, no release/advance and no positive worker completion. RED today includes an erroneous separate preparation blocker. GREEN removes only that duplicate, identifies the bound Run in diagnostics and remains state=blocked through the Run and week consumers. Unknown/uncaptured SL01 start and damaged binding retain preparation blockers. |
| AT05 — authority boundaries | Accepted/blocked/sealed/pending reports, qualifying TaskResult alone, unfinished review/correction, delivery phase/completed and released queue alone never become A02. Retain genuine incompatible A01/A04 and missing-reader A03 blockers; existing compatible A01/A03 behavior remains. Weeks including completed/stop remain unknown. Wrong event/request/Run bindings remain unknown even if the queue is released. |
| AT06 — read-only stability | Repeat collection and public check: deterministic IDs/classes/reasons/digests for unchanged inputs; changed bound evidence changes the digest or blocks. Compare full native file and directory snapshots, registry bytes/format4, refs and source/control/history bytes before/after; zero writes, lock acquisition, process/provider/cleanup calls. Preserve unknown siblings and inventory-level reasons. |

Suggested new focused names: `TestNF01RetiredReportCompatibility`,
`TestNF01HeadlessTerminalCompatibility`, `TestNF01HumanTerminalCompatibility`,
`TestNF01PreparationBindingDiagnostics`, `TestNF01StrictCompatibilityNegatives`,
`TestNF01PublicCompatibilityCheck`. Capture actual failing/passing names, exit
codes and raw logs. Run only these and directly affected existing guard regressions
with `go test ./internal/taskrun -run '<exact focused expression>' -count=1`
in the later implementation Task. No full-suite/coverage run belongs to NF01A.

## Public CLI check demo

The later implementor extends the isolated qualified-reader fixture above, seeded
from the exact clean private make-built candidate (`PLY_DELIVERY_TEST_BINARY`),
with modeled OS/provider seams and native report/headless/human records. Never use
historical live IDs. From that fixture's parent checkout invoke the public CLI:

```sh
"$PLY_DELIVERY_TEST_BINARY" workspace registry activate-parallel \
  --qualified-run "$FIXTURE_QUALIFIED_RUN" --coordinator 'synthetic NF01 coordinator' \
  --allow-retained-source --check --format json
```

This is a **future isolated check demo**, not permission to run an actual
workspace activation check now. Use only `--check`, never `--apply`. First a
terminal-only fixture shows the three proven terminal kinds as A02/empty effects
and state=preview. Add MR05-like bound active headless and week fixtures: output
must identify the Run, omit only the erroneous separate preparation diagnostic,
retain unknown remaining effects and state=blocked. Include one damaged binding
and incompatible A01/A04 variant. Assert JSON state and exact consumers/reasons,
not exit=0 alone (blocked checks can exit 0). Repeat and compare full native
snapshots/format4; no activation/backup/cleanup or new control files. Model
reader qualification; do not bypass the public native qualified-reader gate.
Fixtures prove modeled compatibility decisions, not real provider termination,
installed behavior, human approval or activation readiness of the user's workspace.

## Delivery boundary

**NF01A-01:** current-source table and classification rules cover all four gaps
and the existing typed validators. **NF01A-02:** AT01–AT06 and the isolated public
CLI demo specify focused RED/GREEN and strict negatives. **NF01A-03:** this
Task commits only this Spec. External docs acceptance checks completeness,
source references, sole-file scope and `git diff --check`; its RED/GREEN provides
document evidence only. No product/test code, Go test suite or live probes here.

Return clean candidate OID/tree and actual docs acceptance logs to root w9:pE,
with one short authoring lesson. Root starts/seals separate native review. Hold
owner verify until root authorizes it after sealed PASS; then build the exact
clean private candidate with `make build PLY_BUILD_OUTPUT=/absolute/private/ply`,
check embedded version/build identity and verify without `--review`. Hold
integrate for root release. Only the frozen local_epic_integration target
`refs/heads/ply_agentic_workflow_support` at
`/Users/perottochristensen/github/ply/ply_agentic_workflow_support` is authorized.
No installation, main/master/remote/push/PR, own agents, direct Herdr start, Slack
or postdelivery follow/process inspection. Parent owns notification and cleanup
follow; preserve source/cleanup pending honestly. Future product implementation
is a separate root-registered Task with focused tests, independent native review,
verify/delivery and its separately authorized qualified installed CLI smoke.
