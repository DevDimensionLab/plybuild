# TC9E836AD — Ply module seams for parallel agent work

Analysis only, 23 October 2026. Task `ply-modularity-analysis-20261023`. This report recommends future scopes; it authorizes no implementation, new Task/Spec, queue mutation, installation or lifecycle change. The machine-readable companion is [modularity-candidates.json](modularity-candidates.json).

## Basis and recommendation

Analyzed clean source commit `31f3ca5073c7615933ec0e77bb78e22d9f69bed4`, tree `c721fced22323b2676c42fe60c90b05d6a08dd6b`, on `refs/heads/ply-task/analyze-ply-module-seams-for-parallel-ag-614f6f3349b0` in the assigned Task worktree. The two report files are subsequent analysis changes; production source is unchanged. This is the delivered A/B source base, not a successful coverage baseline.

Prioritize MOD-01 event evidence/replay, MOD-02 rejected-QA recovery, then MOD-03 trace capture/projection. MOD-01 and MOD-03 offer a practical two-worktree option after freezing existing record types; integrate MOD-01 first. MOD-02 has the highest effect risk and should keep one implementation owner. MOD-04 and MOD-05 are evaluated alternatives whose existing seams should be retained. Priority expresses useful responsibility separation, not a CRAP score or commitment to do all five.

The codebase-design vocabulary is used here: a module presents an interface (including ordering, error and authorization rules); a seam is where that interface lives. Depth means useful behavior behind a small interface. Extracting a helper or adding a Go interface alone is not depth. Plain functions and concrete structs are sufficient for these proposals; adapters are justified only where behavior actually varies.

Only static source/test inspection and local Git history informed product conclusions. All named existing tests below were inspected, not executed in this analysis. Some journey fixtures start helper processes; listing them does not select them for execution here. Documentation acceptance validates scope, references, consistency and whitespace only. The mandated private candidate build/version inspection is delivery metadata, not a product-test or runtime-fixture qualification.

Historical planning inputs: [CRAP backlog](../../../planning/reports/crap-backlog-20261022/backlog-plan.md), [failed 23 October coverage](../../../planning/reports/agentic-crap-20261023/report.md), and [ROB-01](../../../planning/specs/rob-01-avgrensede-specs-for-robust-kodefundament-v1.md). These are sibling planning-worktree inputs, not files delivered by this Task. The 22 October master analysis used a timed-out/unqualified profile; the original 23 October fulltest exited 1 with taskrun/workspace timeout. Full coverage remains **UNKNOWN/BLOCKED**. No zero profile counter proves absent tests. The retained B `/bin/ps` denial is not retried or reinterpreted. No new baseline or numeric before/after claim is made.

Local Git history corroborates already completed work: `e5c96d7` introduced within-observation verified-content reuse; `82365d9` accepted ordinary delivery inputs with specific refused-file diagnostics; `aae8fd7` introduced automatic acceptance/local Epic delivery. The source head `31f3ca5` and predecessor `a618d33` concern coverage runtime reproducibility; their presence does not establish a successful full run.

## Candidate comparison

| ID | Priority | Candidate | Risk | Relative effort | Decision |
| --- | --- | --- | --- | --- | --- |
| MOD-01 | 1 | Separate event evidence reads from journal transitions | high | large | recommend |
| MOD-02 | 2 | Make rejected-QA recovery proof distinct from effect publication | very_high | large | recommend |
| MOD-03 | 3 | Separate trace evidence capture from deterministic projection | medium | medium | recommend |
| MOD-04 | 4 | Preserve phase-scoped content observation as the validation seam | high | small_if_needed | retain_existing_seam |
| MOD-05 | 5 | Keep notification outcome classification shared and effect ownership separate | medium | small_if_needed | retain_existing_seam |

Effort is relative: large crosses persisted evidence and historical replay; medium reorganizes existing projection; small_if_needed retains a working seam and limits future work to a demonstrated need. No time estimates or conflict-free merge guarantees are implied.

## MOD-01 — Separate event evidence reads from journal transitions

Priority **1**; risk **high**; relative effort **large**; decision `recommend`. Hard prerequisite candidates: none for the stated first slice.

**Current responsibility and dependency direction.** readJournal reads and verifies the chain, replays fold, projects unobserved process state and checks a cached prefix. fold both loads immutable claim/report/status bytes and validates transitions; foldProviderCompletion also loads launch and provider streams. readJournal replays fold again for cache validation, so knowledge of bound storage and historical transitions is interleaved. Show/collection callers -> readJournal -> fold -> bound claim/report/status files and foldProviderCompletion; appendEvent also invokes the same fold. workflowRead is the separate WorkflowRun persisted-state reader, not a caller of fold.

- [internal/taskrun/store.go](../../internal/taskrun/store.go) — `readJournal`, `fold`, `appendEvent`.
- [internal/taskrun/provider_completion.go](../../internal/taskrun/provider_completion.go) — `foldProviderCompletion`.
- [internal/taskrun/herdr_store.go](../../internal/taskrun/herdr_store.go) — `workflowRead`.
- [internal/taskrun/result_cache.go](../../internal/taskrun/result_cache.go) — `resultCacheValue`.

**Proposed module/interface.** Keep a package-private function: replayEvent(j *journal, e Event, load func(FileBinding, int) ([]byte, error)) error. The loader owns physical read and digest checks; replay owns event-specific path/identity checks and transition order. This is illustrative, not a final signature.

```go
// Existing journal/Event/Result remain authoritative.
err := replayEvent(&j, event, boundFileReader)
// Production and an in-memory corrupt/missing-byte reader cross this seam.
```

**Preserve existing behavior.**

- Exact chain sequence, predecessor digests and event UTC validation stay in readJournal.
- Schema-1 negative claim/cache compatibility keeps original bytes and actual claim provenance.
- A durable process start is not liveness; unknown setup must not relaunch.
- Collection qualification still requires all current process, acceptance, delivery and budget facts; provider completion is not inferred from report prose.
- Preserve error category, precedence, read limits and per-replay evidence freshness; do not cache across observations.

These are responsibility refactors only. Broader retry eligibility, new terminal states, new authority, a new store/schema or different human-QA interpretation would require a separately selected product change.

**Regression protection (proposed, not run).** Through replayEvent, provide identical acceptance events with valid, missing, changed-digest and legacy negative claim bytes; assert complete existing Result and error category. Keep Show-level prefix/corruption and historical-return tests as caller-visible guards. Inventory fixture setup before selecting runnable tests; named journey tests may start helper processes.

- [internal/taskrun/recovery_test.go](../../internal/taskrun/recovery_test.go) — `TestR9CachePrefixAndCorruption`.
- [internal/taskrun/return_repair_history_test.go](../../internal/taskrun/return_repair_history_test.go) — `TestReturnRepairLegacyNegativeCacheKeepsClaimProvenance`.
- [internal/taskrun/scenarios_test.go](../../internal/taskrun/scenarios_test.go) — `TestR7ProcessEvidenceCannotBeInferredFromReport`.

**First slice:** Characterize acceptance_received replay, then separate its bound claim loading from its transition in the existing taskrun package.

**Migration:**

1. Freeze one acceptance family with an in-memory evidence reader and characterization of existing failures.
2. Move only acceptance evidence acquisition behind the loader, preserving transition checks and read order.
3. Review the complete caller plus helper, then consider report/status families separately; provider stream validation is a later slice.

A change to evidence acquisition becomes local without redistributing transition rules across callers. The cost is threading a reader through both replay paths and preserving subtle historical error ordering; moving switch cases alone buys little depth.

**Two-agent ownership:**

- Agent A: Own replay implementation, signatures, journal/cache compatibility and any new private fold file. Files: `internal/taskrun/store.go`, `internal/taskrun/provider_completion.go`, `internal/taskrun/result_cache.go`.
- Agent B: Own a separate characterization test file using the agreed seam; coordinate before changing this existing historical fixture. Files: `internal/taskrun/return_repair_history_test.go`.
- One owner for shared contracts/fixtures: `internal/taskrun/types.go`, `internal/taskrun/fixture_test.go`, `internal/taskrun/recovery_test.go`. A listed existing test file is not jointly writable: allocate new scenario files or explicitly transfer it.

Text conflict if both edit store.go or shared fixtures. Even disjoint tests can accidentally bless a different negative-state projection or weaken completion evidence. Freeze Result/Event and expected errors before parallel writing.

**Safe integration and retest:** Integrate characterization first when it uses existing callers; integrate the extraction next. Tests that require the new replay signature wait for that slice. Recheck Show/cache and historical negative-return behavior on the updated Epic before the next delivery.

**Backlog/ROB-01 reconciliation:** Reuse CRAP-01 and ply-task-run-sjeldne-crash-kombinasjoner.md; overlaps the first responsibility family ROB-01 may select. This narrows that existing candidate; it does not publish a competing Spec or require a full crash matrix.

## MOD-02 — Make rejected-QA recovery proof distinct from effect publication

Priority **2**; risk **very_high**; relative effort **large**; decision `recommend`. Hard prerequisite candidates: none for the stated first slice.

**Current responsibility and dependency direction.** deliveryQARecovery combines owner/runtime observations, immutable request and staged-draft proof, current candidate checks, registry absence checks, receipt creation, report append and retirement of one attempted reservation. Preview and apply deliberately share that procedure; completed recovery is an early read-only receipt lookup. CLI check/confirm adapter -> WorkflowPreviewDeliveryQARecovery or WorkflowRecoverDeliveryQA -> deliveryQARecovery -> workflowhandoff rejection/evidence validation and workspace locked registry snapshot -> workflowKeep/deliveryAppendEvent/workflowSave.

- [internal/taskrun/delivery_qa_recovery.go](../../internal/taskrun/delivery_qa_recovery.go) — `WorkflowPreviewDeliveryQARecovery`, `WorkflowRecoverDeliveryQA`, `deliveryQARecovery`, `deliveryQARecoveryOwner`.
- [internal/workflowhandoff/delivery_qa_rejection.go](../../internal/workflowhandoff/delivery_qa_rejection.go) — `ValidateRejectedDeliveryHumanQA`.
- [cmd/workflow_execute_qa_recovery.go](../../cmd/workflow_execute_qa_recovery.go) — `workflowExecuteQARecoveryCommand`.
- [internal/taskrun/herdr_store.go](../../internal/taskrun/herdr_store.go) — `workflowRead`, `workflowSave`.

**Proposed module/interface.** Prefer a private function proveRejectedQA(existing captured inputs) (deliveryQARejectionReceipt, error), using the existing receipt type. Keep preview/apply orchestration and current locked absence checks together. An abstract generic RecoveryEngine would hide the wrong facts.

```go
// Pseudocode: keep the existing lock nesting and live rechecks.
receipt, err := proveRejectedQA(state, request, stagedProof)
// Only the existing locked executor may persist receipt/report/state.
```

**Preserve existing behavior.**

- Only the known pre-publication +00:00 timestamp schema rejection qualifies; arbitrary failed/unknown effects do not.
- Same live owner, original control/context, current runtime observation and unchanged qualified candidate remain required.
- Keep TaskRun -> project -> work-items lock order; absence check and state save remain inside those locks.
- Completed replay returns original receipt and leaves later QA/integration reservations intact.
- Preserve receipt -> event -> state fault points, prior-receipt recovery and exactly one recover event.
- No QA verdict, provider restart, changed callback control or integration authority is created.

These are responsibility refactors only. Broader retry eligibility, new terminal states, new authority, a new store/schema or different human-QA interpretation would require a separately selected product change.

**Regression protection (proposed, not run).** Cross the existing preview/recover interface with injected Fault at after_receipt, after_event and after_state; retry concurrently and read workflow state to assert one event, no invented QA, and no retirement of a later attempt. Reuse the existing rejection table for changed owner/evidence/control/candidate and published QA. A pure proof table supplements rather than replaces these effect guards.

- [internal/taskrun/delivery_qa_recovery_test.go](../../internal/taskrun/delivery_qa_recovery_test.go) — `TestDeliveryQARecoveryInterruptedAndConcurrent`, `TestDeliveryQARecoveryRejectsUnprovenEffectsAndAuthority`, `TestDeliveryQARecoveryPreservesEvidenceAndAllowsOrdinaryQA`.

**First slice:** Characterize the current unpublished-attempt eligibility decisions with captured evidence, then extract a private rejection-proof function without moving locks or writes.

**Migration:**

1. Name and characterize only eligibility proof, preserving the first failing diagnostic.
2. Extract the private proof computation while all read timing, locks and writes stay in deliveryQARecovery.
3. Only after separate review consider grouping publication steps; keep original fault stages and replay observations.

Reviewers can distinguish why a rejected attempt is safe to retire from how its recovery is durably recorded. The cost is avoiding stale proof between preview and apply; a reusable proof must never become authority outside the locked recheck.

**Two-agent ownership:**

- Agent A: Own the proof extraction, orchestration and all lock/publication changes. Files: `internal/taskrun/delivery_qa_recovery.go`, `internal/workflowhandoff/delivery_qa_rejection.go`.
- Agent B: Own behavior characterization in a dedicated new test file; existing shared rejectedQAFixture remains with A unless explicitly transferred. Files: `internal/taskrun/delivery_qa_recovery_test.go`.
- One owner for shared contracts/fixtures: `internal/taskrun/delivery_types.go`, `internal/taskrun/delivery_fixture_test.go`, `internal/taskrun/delivery_qa_recovery_test.go`, `internal/workflowhandoff/delivery_qa_rejection.go`. A listed existing test file is not jointly writable: allocate new scenario files or explicitly transfer it.

Both agents editing the fixture/test table is a text conflict. Disjoint proof and executor edits can accept stale absence evidence or retire another attempt. Do not divide those semantic responsibilities between simultaneous implementation owners.

**Safe integration and retest:** Land caller-level characterization first, then the proof extraction; requalify interruption/concurrent-retry and ordinary-QA preservation checks on the new Epic base. No automatic dependency on MOD-01 because these are different stored run families.

**Backlog/ROB-01 reconciliation:** Reuse CRAP-02 and ply-task-run-native-avslutning-og-qa-friksjon.md; overlaps ROB-01 first-family selection and its second regression Spec. Startup/inventory recovery (CRAP-03/04) is explicitly separate: readInventoryRecovery only projects recorded metadata and never establishes live authority.

## MOD-03 — Separate trace evidence capture from deterministic projection

Priority **3**; risk **medium**; relative effort **medium**; decision `recommend`. Hard prerequisite candidates: none for the stated first slice.

**Current responsibility and dependency direction.** Read captures registered sources, journal and run history, links native records, rechecks source bytes and registry freshness, then calls finish. builder owns event/source/candidate deduplication, provenance, chains and analysis. These responsibilities already form a useful internal module; expose less capture orchestration to that implementation rather than creating a second trace model. cmd trace consumers -> workflowtrace.Read -> workspace/taskjournal/taskrun captured histories -> builder journal/history/linkNative -> finish/analyze -> Result. Trace reads state; it never authorizes a workflow action.

- [internal/workflowtrace/read.go](../../internal/workflowtrace/read.go) — `Read`, `read`, `directoryDigest`.
- [internal/workflowtrace/project.go](../../internal/workflowtrace/project.go) — `builder`, `event`, `journal`, `finish`.
- [internal/workflowtrace/analysis.go](../../internal/workflowtrace/analysis.go) — `measure`, `analyze`.

**Proposed module/interface.** A private projectTrace function may accept existing captured snapshot/history/registry types and return the existing Result. Keep builder private. Capture passes final source validity and diagnostics after recheck; no new exported projector interface or alternative Event/Task-status schema.

```go
// Pseudocode; use existing record types, not a new wire format.
result := projectTrace(task, journalSnapshot, traceHistory, registry, sources)
err := ValidateReferences(result)
```

**Preserve existing behavior.**

- Deduplicate exact native identity while retaining every source projection and actor provenance.
- Independent chains have no asserted global chronology; sorting is presentation only.
- Reported outcome, technical result, human QA and integration remain separate candidate axes.
- Changed/missing sources make freshness or numeric evidence unknown; partial counts remain lower bounds.
- Current goals remain separate from frozen runs; equal question text does not authenticate or merge actions.

These are responsibility refactors only. Broader retry eligibility, new terminal states, new authority, a new store/schema or different human-QA interpretation would require a separately selected product change.

**Regression protection (proposed, not run).** Use the same projectTrace/Result seam as Read with synthetic existing records: duplicate native event with two chains, conflicting QA, missing source, and changed waiting evidence. Preserve Read-level source recheck and concurrent-append tests; pure projection tests cannot establish physical-source freshness.

- [internal/workflowtrace/analysis_test.go](../../internal/workflowtrace/analysis_test.go) — `TestTraceNativeDedupRetainsAttemptsCandidatesAndSeparateJudgments`, `TestTraceConflictsAndPartialSourcesDoNotChooseWinners`, `TestTraceSourceSequenceAndExactNativeDedupPreserveProjections`.
- [internal/workflowtrace/read_test.go](../../internal/workflowtrace/read_test.go) — `TestTraceSourceRecheckRejectsChangedAndSymlinkedEvidence`, `TestTraceConcurrentJournalAppendMakesCapturedPrefixPartial`.
- [internal/workflowtrace/review_regression_test.go](../../internal/workflowtrace/review_regression_test.go) — `TestReviewChangedWaitingEvidenceDoesNotKeepNumericInterval`.

**First slice:** Characterize event identity/deduplication and per-chain order through the existing builder.finish result before moving projection orchestration.

**Migration:**

1. Characterize current builder.finish outputs and ValidateReferences with existing records.
2. Group the already private projection orchestration behind one function without changing Result or capture recheck order.
3. Migrate one Read call site; retain physical-source tests and remove any duplicate adapter layer.

Projection changes can be reviewed with small deterministic inputs while source acquisition keeps physical freshness checks local. This already has considerable depth, so the improvement must remove orchestration knowledge from Read; wrapping finish alone is not sufficient.

**Two-agent ownership:**

- Agent A: Own projection and deterministic fixtures. Files: `internal/workflowtrace/project.go`, `internal/workflowtrace/analysis.go`, `internal/workflowtrace/analysis_test.go`.
- Agent B: Own captured-source/recheck integration after the private input contract is agreed. Files: `internal/workflowtrace/read.go`, `internal/workflowtrace/read_test.go`.
- One owner for shared contracts/fixtures: `internal/workflowtrace/types.go`, `internal/taskrun/read_trace.go`, `internal/taskjournal/types.go`, `cmd/workflow_trace.go`. A listed existing test file is not jointly writable: allocate new scenario files or explicitly transfer it.

project.go/analysis tests can be disjoint from Read files, but changing source validity or identity conventions can silently change numeric claims. One owner approves the input contract and unchanged Result schema before coding.

**Safe integration and retest:** Land deterministic characterization, then the projection implementation, then adapt capture. Re-run partial-source, dedup, candidate-axis and changed-interval checks against the rebased Epic; do not reuse old candidate evidence.

**Backlog/ROB-01 reconciliation:** Related to CRAP-03 process-journal/observer work but narrower than startup recovery. Reuse that concern rather than duplicating ROB-01 lifecycle regressions; no trace feature or new metrics are selected.

## MOD-04 — Preserve phase-scoped content observation as the validation seam

Priority **4**; risk **high**; relative effort **small_if_needed**; decision `retain_existing_seam`. Hard prerequisite candidates: none for the stated first slice.

**Current responsibility and dependency direction.** Content and queue closure validators share verified owned bytes within one phase. Work-item snapshot, locked publication and goal validation create fresh observations. Task scope limits which content is relevant; phase lifetime is as significant as read arguments. WorkItems Snapshot/Publish or goal operation -> observeTaskContent -> TaskContentStorage.Read -> content/queue validators. Publication still belongs to the store/session, never the observing reader.

- [internal/workspace/content_observation.go](../../internal/workspace/content_observation.go) — `observeTaskContent`.
- [internal/workspace/work_item_store.go](../../internal/workspace/work_item_store.go) — `validateClosure`, `publish`.
- [internal/workspace/goal_operation.go](../../internal/workspace/goal_operation.go) — `validateGoalOperationContent`.
- [internal/workspace/task_content_scope.go](../../internal/workspace/task_content_scope.go) — `WithTaskContentScope`, `validateTaskScopeTransition`.

**Proposed module/interface.** The existing observeTaskContent(*TaskContentStorage) *TaskContentStorage is small and deep enough. Prefer this concrete function/struct. A future validateClosure helper must preserve each caller’s validation/error order; do not introduce a repository/cache interface or a long-lived memoizer.

```go
view := observeTaskContent(source) // fresh for this phase only
// Both existing validators consume view; it cannot publish.
// New snapshot, lock acquisition or publication => a new view.
```

**Preserve existing behavior.**

- Successful reads only are cached, keyed by root/kind/digest, with owned copies.
- Failed reads are retried; identical digest in another root does not imply existence.
- No observation survives a registry reread, lock acquisition or proposed publication.
- Scope and queue/content closure checks both remain enforced; snapshot and publish error ordering is preserved.

These are responsibility refactors only. Broader retry eligibility, new terminal states, new authority, a new store/schema or different human-QA interpretation would require a separately selected product change.

**Regression protection (proposed, not run).** Exercise WorkItemStore Snapshot and session Publish, the real caller seam: count unique managed reads within one phase, corrupt content before the next snapshot/relock/publication, and require rejection. Keep direct view tests for copy ownership, missing-file retry and cross-root separation.

- [internal/workspace/content_observation_test.go](../../internal/workspace/content_observation_test.go) — `TestWorkItemSnapshotReusesVerifiedContentWithinObservation`, `TestContentObservationOwnsBytesAndRetriesMissingSources`, `TestWorkItemObservationRejectsCorruptionOnNextSnapshotAndPublication`.

**First slice:** Keep observeTaskContent intact; only a demonstrated caller problem justifies extracting combined closure validation, with same-phase and next-phase tests first.

**Migration:**

1. Retain current implementation and document lifetime at any new caller.
2. If a later concrete duplication warrants it, characterize caller error order and factor only that repeated closure responsibility.
3. Review all phase creation sites together; never broaden observation lifetime for speed.

Already concentrates verified-byte reuse without weakening publication. A new abstraction now would add interface cost with little leverage; future work should address a specific caller defect, not rename the existing module.

**Two-agent ownership:**

- Agent A: Single owner for observation lifetime and closure orchestration if selected later. Files: `internal/workspace/content_observation.go`, `internal/workspace/work_item_store.go`, `internal/workspace/goal_operation.go`.
- Agent B: Independently characterize mutation/corruption scenarios after lifetime contract is frozen. Files: `internal/workspace/content_observation_test.go`.
- One owner for shared contracts/fixtures: `internal/workspace/task_content_scope.go`, `internal/workspace/work_item_store.go`, `internal/workspace/work_item_store_test.go`. A listed existing test file is not jointly writable: allocate new scenario files or explicitly transfer it.

Shared workspace store/scope files create text conflicts with registry/queue refactors. Disjoint caches and recovery edits can still reuse stale evidence across locked absence checks. Implementation should wait when those invariants are changing.

**Safe integration and retest:** No refactor is currently recommended. If selected later, integrate focused characterization then a single lifetime-preserving slice; after other workspace delivery recheck scoped closure and corruption-on-publication before continuing.

**Backlog/ROB-01 reconciliation:** Related to CRAP-08 registry invariant families and ROB-01 prerequisite stability, not a replacement for queue validation. Git history e5c96d7 already introduced this verified-content reuse; preserve its work instead of proposing it as missing.

## MOD-05 — Keep notification outcome classification shared and effect ownership separate

Priority **5**; risk **medium**; relative effort **small_if_needed**; decision `retain_existing_seam`. Hard prerequisite candidates: none for the stated first slice.

**Current responsibility and dependency direction.** Both notification execution and integration notification sending call the same classify function after an injected transport observation. HTTP adapter bounds reads and disables redirects/replay. Separate effect owners reserve durable attempts and persist classification; unknown dispatch is deliberately distinct from proven non-delivery. Execute or sendIntegration -> Dependencies.Transport -> postHTTP Observation -> classify -> caller-owned durable Attempt. ObserveIntegration is a read-only state observation and does not consult credentials or send.

- [internal/workflownotification/transport.go](../../internal/workflownotification/transport.go) — `sendHTTP`, `postHTTP`, `classify`.
- [internal/workflownotification/service.go](../../internal/workflownotification/service.go) — `Execute`.
- [internal/workflownotification/integration_transport.go](../../internal/workflownotification/integration_transport.go) — `ObserveIntegration`, `SendIntegration`, `RetryIntegration`, `sendIntegration`.
- [internal/workflownotification/types.go](../../internal/workflownotification/types.go) — `Dependencies`, `Observation`, `Attempt`.

**Proposed module/interface.** The existing function classify(Observation, time.Time) Attempt is sufficient; the Transport function already has real HTTP and injected-test adapters. Do not introduce another Go interface, generic notifier or unified notification ledger.

```go
observation := deps.Transport(ctx, endpoint, payload)
attempt := classify(observation, deps.Now())
// The caller preserves its reserved attempt identity before dispatch.
```

**Preserve existing behavior.**

- Only bounded exact 200/ok acknowledges transport; it never proves human reading or product approval.
- Truncated/oversize or ambiguous dispatch remains unknown; unknown and acknowledged attempts cannot be retried.
- Only recognized rejection or valid single Retry-After establishes a retryable result; retain explicit confirmation rules.
- No hidden HTTP redirect, proxy lookup or replay; preserve separate route/event/credential bindings and durable reservation order.

These are responsibility refactors only. Broader retry eligibility, new terminal states, new authority, a new store/schema or different human-QA interpretation would require a separately selected product change.

**Regression protection (proposed, not run).** At postHTTP/classify, use an injected RoundTripper for redirects, oversized/truncated body, timeout before versus after possible dispatch, and duplicate/invalid Retry-After. Keep owner-specific reservation/receipt-failure tests to prove shared classification does not create resend authority.

- [internal/workflownotification/transport_test.go](../../internal/workflownotification/transport_test.go) — `TestHTTPClassificationAndBoundedAdapter`, `TestTruncatedBodyIsUnknown`.
- [internal/workflownotification/notification_test.go](../../internal/workflownotification/notification_test.go) — `TestReservationAndReceiptFilesystemFailure`, `TestCorruptPriorStateNeverClaimsNotAttempted`.

**First slice:** Retain classify(Observation, time.Time) Attempt and the injected Transport function; use a focused classification table only when an observed outcome needs a change.

**Migration:**

1. Keep classifier and injected transport unchanged unless a real defect is selected.
2. For a selected outcome change, add a focused caller-visible table case and review effect-owner retry consequences.
3. Avoid merging the two durable schemas; adjust callers only if behavior demands it.

This is a positive example of depth: callers learn one observation/result contract and reuse precise outcome rules. Additional wrappers would add cost; the existing effect ledgers differ for valid ownership reasons.

**Two-agent ownership:**

- Agent A: Own shared classification and HTTP adapter behavior. Files: `internal/workflownotification/transport.go`, `internal/workflownotification/transport_test.go`.
- Agent B: Own an independently justified integration receipt change with separate tests after classifier contract is fixed. Files: `internal/workflownotification/integration_transport.go`.
- One owner for shared contracts/fixtures: `internal/workflownotification/types.go`, `internal/workflownotification/service.go`, `internal/workflownotification/store.go`. A listed existing test file is not jointly writable: allocate new scenario files or explicitly transfer it.

Classifier and ledger files need not overlap textually, but a changed unknown/rejected distinction changes retry eligibility. One owner must approve any such behavior change; that would be a separately selected product change, not this refactor.

**Safe integration and retest:** No extraction is currently recommended. If independent adapter/receipt maintenance is later selected, serialize Epic returns and rerun shared classification plus the second owner’s unknown-effect/retry guards on the updated base.

**Backlog/ROB-01 reconciliation:** No direct CRAP-01/02 or ROB-01 implementation duplicate. Existing notification/continuity work already owns retries and receipts; reuse its tests. Do not attach this to ROB-01 merely to increase candidate count.

## Candidate-pair ownership matrix

All pairs permit independent analysis. Implementation independence is limited to the stated slices and frozen existing contracts. Every local Epic integration is serial; none of these rows authorizes overlapping integration effects. Retained candidates MOD-04/05 are conditional future work, not parallel assignments.

| Pair | Implementation | Text conflict / ownership | Semantic conflict and serial return order |
| --- | --- | --- | --- |
| MOD-01 + MOD-02 | parallel_after_contract_freeze | Separate store/fold and delivery_qa_recovery files; taskrun types/fixtures need one owner. | Different run families, but shared TaskResult/evidence meanings must stay fixed. MOD-01 then MOD-02; recheck recovery evidence and historical Result contracts. |
| MOD-01 + MOD-03 | parallel_after_contract_freeze | Disjoint packages if read_trace and journal types remain unchanged. | Trace must not infer completion or change identity when fold historical projection changes. MOD-01 then MOD-03; compare historical negative state, dedup and candidate axes. |
| MOD-01 + MOD-04 | conditional | Fold files disjoint; workspace types/store become shared if extraction expands. | Evidence freshness and publication-phase assumptions may drift. MOD-01 first; MOD-04 only if selected later, then validate fresh scoped closure. |
| MOD-01 + MOD-05 | independent_scoped | Disjoint taskrun and notification implementation files. | Preserve shared meaning of unknown; neither may invent completion from a notice. MOD-01 then selected MOD-05 maintenance; rerun both affected guards. |
| MOD-02 + MOD-03 | parallel_after_contract_freeze | Recovery and trace files disjoint; trace input record types frozen. | Recovery report is not human QA or an integration; trace must preserve that distinction. MOD-02 then MOD-03; recheck trace candidate axes after recovery fixture changes. |
| MOD-02 + MOD-04 | serialize_implementation | Potential overlap in workspace store/scope/evidence dependencies. | Moving absence checks or reusing bytes across locks can retire an unknown effect. Selected MOD-04 contract change first, then rebase MOD-02 and revalidate locked absence/replay; otherwise defer MOD-04. |
| MOD-02 + MOD-05 | independent_scoped | Recovery and notification files disjoint. | Unknown recovery and unknown Slack delivery have separate receipts; sharing classification must not grant recovery authority. MOD-02 then selected MOD-05 maintenance; verify neither side changes the other effect ledger. |
| MOD-03 + MOD-04 | conditional | Trace Read and workspace store disjoint unless ReadViewBasis/types change. | Observation lifetime and source validity affect trace freshness. Selected MOD-04 first, then rebase MOD-03 and recheck changed-source/partial-history behavior; otherwise MOD-03 alone. |
| MOD-03 + MOD-05 | independent_scoped | Disjoint trace and notification files. | A transport acknowledgement must not become a human-judgment trace record. MOD-03 then selected MOD-05 maintenance; retain distinct evidence classes. |
| MOD-04 + MOD-05 | independent_scoped | Workspace content and notification transport files disjoint. | Keep separate observation lifetimes and receipt stores; no shared cache or new generic store. Selected MOD-04 then selected MOD-05; recheck closure and notification replay independently. |

## Practical worktree and integration plan

1. The coordinator chooses one coherent future slice under the existing backlog/ROB-01 selection process. It freezes the relevant existing record types, errors, file ownership and fixture ownership before starting implementation. No general event bus, plugin framework, second Task-status representation or universal recovery engine is proposed.
2. A useful optional pair is Agent A on MOD-01 replay and Agent B on MOD-03 projection, each in its own native Task worktree from a qualified Epic base. A owns taskrun replay files; B owns workflowtrace files. `internal/taskrun/read_trace.go`, `internal/taskjournal/types.go` and exported Result/Event shapes stay fixed under one contract owner. If either needs those files, stop parallel implementation of that seam and sequence the contract change first.
3. Within MOD-02, split independent characterization from implementation, not proof from the effect executor. Only one implementation owner touches lock placement, reservation retirement and receipt/event/state publication. Existing fixture owners approve new scenario files; a second agent must not rewrite the same fixture concurrently.
4. Integrate MOD-01 before MOD-03; MOD-02 can be implemented independently once contracts are fixed, but its Epic return is also serial. If ROB-01 chooses MOD-02 first, that choice supersedes this suggested ordering; this analysis does not impose a lifecycle dependency between separate run families. MOD-04/05 stay deferred unless separately justified. Matrix orders involving them are conditional, not a mandate to implement them.
5. After each integration, observe the new Epic ref/base and queue through the native coordinator route. Rebase or recreate the later authorized worktree against that observed base using supported workflow. Review text conflicts and semantic assumptions even if Git merges cleanly. Rerun the affected bounded caller-level guards identified above, plus checks for changed shared types/fixtures. A widened diff requires a fresh review scope.
6. Preserve old OID/tree/script/binary/review/test bindings as historical evidence. New source means a new exact candidate, fresh bounded acceptance and separate native review; do not edit old evidence to make it current. Coordinator-owned final fulltest waits for a stable integrated base and its separate readiness/authority. This Task neither runs it nor clears the retained denial.

## Verification and delivery limits

The frozen external acceptance entrypoint checks these two files, structured candidate relationships, source/test symbols, Markdown links and Git whitespace/scope against the recorded analysis base. This is document integrity evidence, not proof that the proposed refactors preserve runtime behavior. Review must judge whether the suggested interfaces actually remove caller knowledge, since a script cannot prove module depth.

This analysis Task changes only `docs/robust-codebase/modularity-analysis.md` and `docs/robust-codebase/modularity-candidates.json`. Acceptance code and logs are private external delivery artifacts. Native managed task-requirements evidence, exact review and automatic acceptance belong to the candidate verification flow; this report is not a second requirements/status database.

Authorized return is one native local integration from this Task source to `refs/heads/ply_agentic_workflow_support` at `/Users/perottochristensen/github/ply/ply_agentic_workflow_support`. Clean documentation is held unchanged for coordinator-started separate native review. The same owner verifies without a handwritten review, then integrates after sealed pass. `owner_handoff` exit 5 means integration is observed and source cleanup awaits coordinator follow. This static report does not itself assert that those later steps passed. Parent backlog-coordinator w9:pE owns review/follow/closeout and aggregate notification; this owner sends no separate Slack.
