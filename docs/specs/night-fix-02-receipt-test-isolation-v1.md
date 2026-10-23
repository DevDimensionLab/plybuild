# NF02 receipt journey isolation Spec v1 — T2768B475

Specify a test-only correction for the receipt journey on source `6a3a4f01de5214f7a0aa799c45d86e4ad6cb136c`
(`6a3a4f0`). This authoring Task changes only this document. All test changes
and commands below are requirements for a later developer, not executed results.

## Source finding

[TestDeliveryReceiptCLIJourney](../../internal/taskrun/delivery_receipt_test.go)
creates a receipt fixture, binds the candidate executable, patches the modeled
owner pane, then calls `workflowTestStart`. It exercises invalid receipt bindings,
idempotent check-in, unknown sandbox coverage, rejection of legacy begin, and
candidate verification through subprocess CLI. It ends in
[assertAutomaticCLIClosure](../../internal/taskrun/delivery_automatic_cli_test.go).

That shared helper registers/checks/executes public Delivery, verifies exact
integration and source retention, sets only `agent_status=idle`, then calls
`cleanupCLIAt(binary, f.Parent, "execute", "follow", ...)`.
[cleanupCLIAt](../../internal/taskrun/delivery_default_cleanup_test.go) creates a
fresh `exec.Command`; it cannot inherit Go dependency closures. Idle agent status
also does not establish that the bound pane has left. `ownerPaneStandin` supports
explicit `cleanupOwnerLeaves`/`closed_transport`, which this path does not use.

The remaining OS fallthrough is statically concrete:
[ResumeAutomaticLocalDelivery](../../internal/taskrun/delivery_automatic_closeout.go)
calls `automaticOwnerLeft`, then
[taskWriters.cleanup](../../internal/taskrun/delivery_task_writers.go) composes
Herdr absence with local liveness. Once Herdr absence is established,
[observePathWriterLiveness](../../internal/taskrun/delivery_writer_liveness.go)
uses `systemLocalProcesses` when `LocalProcesses` is nil: bounded `lsof` and,
for a session, `/bin/ps`. Setting `f.D.LocalProcesses` cannot isolate the child CLI.
This is test isolation debt, not a demonstrated production guard bug. No fresh
journey, process inspection or historical denial was replayed for this Spec.

AL18 already fixed the FS01 pane-get stand-in and modeled the handoff, explicit
owner exit and in-process follow. Its
[TestFS01ReviewRoundsFromExistingTaskToDelivery](../../internal/taskrun/delivery_review_round_test.go)
uses a coordinator `LocalProcesses` table and the shared unknown-inventory
assertion. Preserve that fix. AL19's
[TestDeliveryStartupCLI](../../internal/taskrun/delivery_startup_cli_test.go)
already uses cause/action `resolve-stop`, preserving the original denial and
`PermissionConfirmed=false`. Preserve it too. Historical receipt `/bin/ps`
denial is context, not evidence that either delivered fix regressed.

## Minimal test-only change

Limit implementation to `internal/taskrun/delivery_receipt_test.go` and, only
where extracting assertions helps, `internal/taskrun/delivery_automatic_cli_test.go`
or `internal/taskrun/delivery_default_cleanup_test.go`. Add a receipt-local
isolated closure helper; retain the public closure helper's behavior for its
other callers. Do not silently convert every existing CLI test to simulation.

1. Before startup, install explicit fixture dependencies, then invoke
   `workflowTestStart`. Reuse the shapes in
   [parallelRuntime](../../internal/taskrun/parallel_epic_return_test.go):
   `Dependencies.LocalProcesses` returns a valid nonempty table containing the
   observing PID; `StartupProcessTree` returns the controlled launch tree.
   Use `ownerPaneStandin` before launch, retaining its updated executable hash
   and request. Use receipt-local setup, without changing blanket fixture defaults.
   Audit all startup/verification paths before running; stop at any unmodeled
   boundary, never try the real OS to discover whether isolation works.
2. Explicitly replace `Workspace.SourceWriters` with a counted simulated observer
   appropriate to this disposable fixture. Rebind `Workflow.TaskWorkspace` to
   the updated workspace dependencies. Merely copying `f.D` is insufficient:
   [SystemDependencies](../../internal/taskrun/service.go) creates callbacks that
   capture earlier dependency values. Do not assume changing `LocalProcesses`
   changes a captured `SourceWriters` observer. Reuse the initialization pattern
   in [nf01Runtime](../../internal/taskrun/registry_nf01_test.go), without importing
   its unrelated registry behavior or relaxing production defaults.
3. Keep the existing receipt subprocess checks and exact candidate verify.
   Preserve wrong context/provider/session rejection, repeated receipt identity,
   unknown sandbox/per-action facts and legacy protocol rejection. Keep the
   candidate path/hash in `Runtime.PlyExecutable`, `D.Executable`, verify input
   and the acceptance script aligned. Retain the native generated bound control;
   do not substitute a fake successful executable or manufacture terminal receipts.
4. Replace this journey's call to `assertAutomaticCLIClosure` with an explicitly
   named simulated-observer closure helper. Use `WorkflowDeliveryIntegrate` with
   the same Run/context and injected owner dependencies. Assert
   `cleanupOwnerHandoff`: `owner_handoff`, phase `closing`, delivered integration
   and retained source. Snapshot target OID/reflog, delivery events, queue closure
   and Epic base at this point. Use `cleanupOwnerLeaves` to model actual exit,
   not merely idle status. Copy coordinator dependencies with `CWD` returning
   `f.Parent`; execute `WorkflowFollow` from outside the source with the existing
   modeled Herdr identity conventions. Keep owner callback cwd at the Task source.
5. Reuse `assertCleanupUnknownInspectionPreservesDeliveredIntegration` through
   its callback seam, as AL18 does, and reuse `assertCleanupRemoved` for positive
   completion. Supplement their assertions with queue/base stability and absence
   of fabricated human QA. Keep the target worktree, archive candidate and source
   branch/removal assertions. Repeat completion/readback with the same Run and
   result and prove the target reflog/integration event did not change.

Public Delivery registration/check/execute is no longer covered by this receipt
journey's closure section. Existing public Delivery tests retain that duty.
Do not add a CLI flag, environment setting, wire field, PATH-based absence tool,
new process service or production injection mechanism to preserve that label.
If investigation shows any retained subprocess step requires an OS dependency
that cannot be controlled through existing seams, stop that step and return the
precise boundary to root; do not claim complete isolation or broaden product scope.

## Controlled regression

Tie the regression oracle to the actual receipt journey's startup and closure
callables/dispatch. Add a small receipt-local spy/tripwire at those boundaries,
shared by `TestDeliveryReceiptCLIJourney` and its regression case. It must observe
what that orchestration actually calls, rather than validate a separately supplied
route or a list of non-nil adapters. Keep all instrumentation test-only and local;
no generic route platform or new production seam is needed.

Controlled RED must exercise that same orchestration with deliberately wrong
wiring: a stale workspace binding carrying a failing `SourceWriters` spy, or a
closure dispatch that attempts the subprocess boundary. Install fail-fast test
tripwires before startup and before invoking any action that could reach the OS.
The boundary tripwire must fail before `exec.Command`/real inspection is attempted;
the stale-binding spy must fail at its actual production-function call boundary.
For an early boundary failure assert zero downstream launch/inspection/integration
effects and unchanged native state; for a later wiring failure compare the state
at entry to that phase. Never execute the old OS journey to obtain RED. A negative
test of the tripwire alone does not qualify: changing the real journey's dispatch
back to `cleanupCLIAt`, or rebinding `TaskWorkspace` to the stale workspace, must
make this regression fail safely. If that protection cannot be demonstrated with
the local instrumentation, return the exact boundary to root before execution.

GREEN must show phase-specific behavioral evidence through that same call path:

- At the startup/recovery boundary consuming `StartupProcessTree`, record the
  spy's actual shell PID and call count. If ordinary startup does not enter that
  branch, assert zero calls there and exercise the consuming branch using the
  same receipt startup callable with a controlled preserved-shell fixture. Do
  not call the adapter directly and report that as startup coverage.
- At the source-writer boundary, prove that production code reaches the rebound
  `Workspace.SourceWriters` through `Workflow.TaskWorkspace`: record the source
  path and phase, while the earlier captured/stale observer is a failing spy.
  Require the intended spy to run and the stale spy to remain unused.
- After modeled owner exit, require `WorkflowFollow` to reach the injected
  `LocalProcesses` spy with the bound session/transcript and observing PID table.
  Assert its call count and the resulting native closeout, not adapter presence.

Keep all observers initialized before startup. Preserve genuine native Git/journal
effects in the disposable fixture, without altering their oracle files. Report
only branches actually exercised; a declaration of wiring is not coverage.

Add receipt-local negative closeout subcases with fresh fixtures: modeled Herdr
inventory denied (`retire_inventory_denied` via the existing helper), local
`os.ErrPermission`, and incomplete/empty local inventory. Each must return unknown
(`owner_unknown` where exposed), keep phase `closing` and final delivery
`delivered`, preserve source directory/ref, target OID/reflog, native integration
binding/events, closed queue and advanced Epic base. No negative observation may
produce an absence receipt. Count the denied managed or local inspection exactly
once per explicit follow attempt; no automatic retry or reintegration is allowed.
After the test explicitly replaces the simulated failure with a complete table,
continue the same closeout and assert normal removal. Do not retry actual denial.
Also preserve a positive external Task/session writer case that blocks removal.

## Coverage and acceptance

| Boundary | Evidence the later implementation may claim |
| --- | --- |
| Receipt check-in, invalid bindings, legacy rejection, candidate verify | Exact supplied subprocess CLI with modeled provider/Herdr; retain argv, cwd, stdout/stderr and exits. |
| Startup and writer absence/denial; integration and closeout | In-process production functions with simulated observer coverage and real disposable native Git/queue/base records. |
| Public Delivery execute/follow after the split | Not exercised by the new receipt closure; no subprocess closeout claim. |
| Installed CLI, live provider, actual OS inventory or sandbox enforcement | Not established by this test or this docs-only Task. |

For the later developer, run only the receipt-call-path regression first; preserve
its controlled RED, then GREEN. Run the corrected
`TestDeliveryReceiptCLIJourney` and its new guard subcases with an absolute clean
make-built `PLY_DELIVERY_TEST_BINARY`, `-count=1` and a bounded Go timeout. Include
`TestAR03ComposedWriterEvidence` as the focused guard matrix only after auditing
its setup for pre-start isolation; its cases cover denied, empty, missing-self,
duplicate identity, external writer and session elsewhere. No real-OS writer test,
old broad CLI cleanup test or full suite is required. Existing AL18/AL19 assertions
remain intact; static preservation is not an executed regression pass. Record
actual test selections, limits, candidate identity and outcomes without claiming
unrun coverage. Format touched Go files and run `git diff --check`.

This authoring Task's bounded acceptance checks this Spec's required design
boundaries, local source links, single-file scope and `git diff --check`; it runs
no Go tests or receipt CLI journey. Root w9:pE starts separate native review of
the clean docs commit. Owner verify waits for sealed Spec-reviewPASS and uses
the private make-built candidate without `--review`. Integration waits for root
release and returns only to the frozen local Epic target; no installation.
Source cleanup remains pending for root, with no postdelivery follow or OS probe.
