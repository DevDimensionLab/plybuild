# NF04A interruption owner boundary v1 — T86252AB1

Decision: **outside-repository rig boundary; no proven minimal safe Ply patch**.
Keep the interruption remainder unresolved. Do not commission a developer from
this Spec. Root w9:pE owns the batch disposition; the rig/taskrun maintainers own
any separately selected cancellation contract and implementation.

Static source basis: clean Epic commit `0bfe0ae3df2ec36c4be55f0c00211c24ae8ad1d4`,
tree `132c0dd609a44fd2383eaf5a8e9d08de11a8b9b7`. Only this document changes.
No Go, rig, signal or provider tests, process inspection or denial replay ran.

## Historical evidence

The [retained backlog](/Users/perottochristensen/github/ply/planning/koe/ply-task-run-native-avslutning-og-qa-friksjon.md)
and [session 141 control](/Users/perottochristensen/github/ply/planning/design/evidence/69-factory-rig-human-qa-v3-result-control.json)
identify factory QA v3, `ply-factory-rig-human-qa-run-003`, interrupted on
2 October 2026. The control binds rig commit
`ceb6a30e0790a25b9d1906c3b4395abb8dfe60e6` and Ply commit
`15a0d82fbcad3a38d722306438401e449f570f07`.

Command 34's exit remained unknown. The helper report was written about ten
seconds after the outer blocked report; no completed provider turn, collected
typed delivery, full journey or human product verdict was established. The
control infers that terminating the direct Ply group bypassed its normal waiter
while Codex occupied a separate group. It establishes neither a normal-path
hang nor an effect outside the private fixture.

The historical denied process query was not retried. The human's report of no
rows for two main PIDs was not all-descendant closure. The control explicitly
says v4 did not fix cancellation supervision. Its preflight/banner/heartbeat
work is a separate historical preparation, not a cancellation repair or fresh
process observation. Old roots, unknown effects and ownership remain preserved.

## Current source reconciliation

| Boundary inspected statically | Current semantics and limit |
| --- | --- |
| External [riglib.py](/Users/perottochristensen/github/ply-factory-tests/riglib.py), `Commands.run`, lines 138–184 | Starts its direct command with `start_new_session=True`. On timeout/KeyboardInterrupt it sends TERM, then possibly KILL, to `child.pid`'s group and waits for that child. It records unknown descendant lifetime. The named current file and historical Git blob both hash to `53652aa7662a75333b443e712a98bd00670c3e8da857b8dbc01b894740980a2d`, matching the control. This source is available, unchanged at that path, and outside the Ply repo and this Task's implementation scope. |
| [service.go](../../internal/taskrun/service.go), `Apply`, schema-2 runner selection; [exec_runner_darwin.go](../../internal/taskrun/exec_runner_darwin.go), `execRunner.Run` / `runOwnedExec` | Schema 2 selects ExecRunner. The Darwin pipe runner creates a separate provider group with `Setpgid: true`. Its sole waiter serializes Wait4 and group signals; timeout/observation error drives TERM/KILL before reaping. It has no signal subscription for parent interruption. The historical runner has the same group separation and timeout ownership. An outer group kill still does not by itself reach that provider group. This is a source boundary, not a newly reproduced incident. |
| [process.go](../../internal/taskrun/process.go), `InteractiveRunner.Run`; [process_darwin.go](../../internal/taskrun/process_darwin.go), `darwinChild.Forward` / `Wait` | The interactive terminal path already forwards INT/TERM/HUP to a birth/group-bound child, and disables forwarding before reaping. That is a different runner from factory schema 2; copying its forwarding alone would not prove the rig-to-pipe-runner journey safe. |
| [delivery_automatic_verification.go](../../internal/taskrun/delivery_automatic_verification.go), `runDeliveryVerification`; [delivery_automatic_process_unix.go](../../internal/taskrun/delivery_automatic_process_unix.go), `configureAutomaticProcess` | Delivered in `aae8fd70d73d72991ecda62fe9168b808587a93a`: automatic acceptance uses NotifyContext plus timeout, a separate command group, group KILL on cancellation and WaitDelay. Signal handling remains installed until the caller preserves the result. This is existing automatic-verifier behavior, not a fix to the external rig or provider runner. |
| [worker_unix.go](../../internal/taskrun/worker_unix.go), `StartWorker` / `ControlWorker`; [worker_attempt_unix.go](../../internal/taskrun/worker_attempt_unix.go), `RunWorkerAttempt`; [worker_exec_darwin.go](../../internal/taskrun/worker_exec_darwin.go) | The newer detached native controller uses Setsid; pause/stop requests retain active-attempt ownership. Provider attempts reuse runOwnedExec with a different maximum budget. Native stop is not equivalent to signaling the historical rig group, and does not prove all writers absent. |
| [collect.go](../../internal/taskrun/collect.go), `collectPreview` | Schema-2 qualification still requires bound completion, exit zero and group quiescence, plus the native report/receipt/target gates. Preserve these checks; a later helper report cannot turn unknown completion into a qualifying result. |

The [existing pipe-runner tests](../../internal/taskrun/exec_runner_darwin_test.go)
cover normal/delayed completion, timeout, stream limits and reaping. The
[Darwin forwarding test](../../internal/taskrun/process_darwin_test.go) models
birth/group mismatch and post-wait rejection. The
[automatic interruption tests](../../internal/taskrun/delivery_automatic_execution_unix_test.go)
exercise verifier interruption and result preservation. These are source-read
test inventories, not passes in this Task; none establishes this external
rig/separate-provider cancellation journey. NF01 reservation isolation and NF02
receipt-fixture isolation already delivered on this base do not resolve it.

## Next owner and safe later boundary

The **external rig maintainer**, with the **Ply taskrun maintainer**, must first
select whether the rig should request a native cancellation and wait, or whether
the synchronous pipe runner should consume termination while keeping its sole
waiter alive through completion preservation. Changing only the rig's direct
kill target does not establish ownership of the separate provider group;
adding a signal handler does not establish the outer escalation/wait contract.
The inspected source does not prove a minimal patch that resolves both ends.
This Task authorizes neither an external-repo edit nor a lifecycle redesign.

If a later assignment proves a small correction, bind its exact source branch,
affected consuming functions and bounded cancellation semantics before execution.
Reuse the existing runner/terminal/command seams. A safe RED must reach the
actual interruption branch with an owned test supervisor and separate synthetic
child group, readiness handshake, explicit handles and bounded waits, without
a real provider or signals to the test host/old Run. GREEN must observe forwarding,
wait/reap and durable interrupted completion at that same branch, including
identity-change, early-exit and escalation failures. No enum-only assertion or
empty event stream counts as this coverage. Declare exact selected test names
and inner/outer limits in that later Spec; none is commissioned here.

Record outcomes per owned child/group. Missing identity, failed observation,
unreaped child or writer outside the owned group remains unknown and retains
resources. Never signal after ownership is lost/reaping has occurred, invent a
successful exit, or infer global writer absence from stopped, one client exit
or group emptiness. Preserve genuine denials and historical unknowns unchanged.

The immediate user benefit is an exact owner boundary and avoiding a misplaced
Ply patch or another unisolated QA run. Static reconciliation is bounded; patch
effort and safety remain uncertain until the cross-boundary contract is chosen.
Prefer that small explicit decision before proposing new supervisor/runtime
infrastructure, consistent with the supplied cost/benefit guidance.

## Acceptance and handoff

Bounded docs acceptance checks the decision/evidence/owner sections, source
references, this one-file base-to-candidate diff, a clean candidate and
`git diff --check`. A private make candidate serves native docs verification
only; it proves no interruption behavior and is never installed.

Root starts separate native Spec review. Hold verify until root releases it
after PASS, then the original owner verifies this exact candidate without
`--review`. Hold the single authorized local Epic integration for root's later
release and fresh target/acceptance check. Keep source cleanup pending for root;
no postdelivery follow/process inspection. Root owns learning and the aggregate
notification; this child sends no Slack.

Learning: group cancellation and completion evidence belong to their actual
consumer. Existing automatic-verifier cancellation cannot close an unchanged
external rig's separate-provider interruption finding.
