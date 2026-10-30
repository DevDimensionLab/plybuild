# Current Run, Worker and Week interpretation

> Historisk dokument etter fjerningen av den kompilerte prosesskjeden.
> Den beskrevne current-state-/resume-oppskriften og dens testwrapper er fjernet.
> Se [gjeldende direkte operasjoner](../README.md#human-owned-process).

`internal/currentstate` owns the interpretation of current worker, product-answer
and Week facts. Its `RunAction`, `WeekEntry` and `WeekQueue` functions are pure:
they cannot launch, consume controls or answers, repair stores, take ownership,
write history, release a queue or deliver work.

Previously `WorkflowShow` chained `workerCurrentAction` and
`productCurrentAction`. Their eligibility depended on the saved guidance actor
and, for resume, the complete human message. Week Show selected observations
using `answerGuidance`, while `projectAnswerWait` recognized a recoverable
answer block by the `Native Run next actor ` prefix. Tick repeated phase rules
around that projection. Changing display wording could leave stale guidance or
prevent recovery of an answer-derived block.

Now `WorkflowShow` gathers native evidence through `currentRunAction` and
`productCurrentFacts`, then calls `currentstate.RunAction`. Worker obstacles and
Run diagnostics retain precedence; fresh worker controls remain unconsumed;
parking, provider waits, answer receipt, reservation, observed continuation,
confirmed use and completed delivery select their own actor. Presentation
messages do not select those transitions. Provider waits written by the worker
carry a `reason_code`; changing their diagnostic wording preserves the wait and
fresh control guidance.

Week Show observes the active native Run even before the next Tick. Both reads
and bound execution use `projectCurrentRun` → `currentstate.WeekEntry`. New Week
records carry `guidance_origin` and entry `block_origin`. Only
`answer_continuation` blocks can be cleared by continuation progress; ordinary
blocks, including uncertain release, retain their state even if the display
contains a native-actor prefix. Completion comes from the native receipt, not a
message. A completed entry never becomes running work from stale observations.
Native dispatch, ownership checks and queue-release policy remain in the
existing execution paths and adapters.

The deletion test: removing this module would put worker/control/answer
precedence back into Run callers and put wait, uncertainty, block recovery and
completion interpretation back into both Week Show and Tick. The adapters keep
evidence loading and execution; they no longer need the interpretation order or
the vocabulary used to render messages. NativeAdapter and fakeNative still
exercise the same Week interface; no new execution adapter was introduced.

`currentstate/legacy.go` isolates compatibility reads for historical provider
wait text and Week guidance without origins. New initial records have explicit
origins. Reading older records projects their meaning in memory and preserves
their stored bytes. Initial-registration recovery also accepts the exact older
initial shape without widening its recovery conditions. Frozen contracts are
unchanged.

## Concrete Worker outcomes and queue release

Worker records `outcome` when it interprets a completed developer block, a
blocked review or exhausted correction attempts. Other blockers record an
explicit `unknown`, so even diagnostic text that resembles a concrete failure
cannot change their meaning. Clearing or rewriting the diagnostic preserves
this fact; starting a continuation or another stage clears it. Provider waits
and AR01's `reason_code` keep their existing precedence.

`WorkerStatus.QueueReleaseAdvice` owns concrete outcome classification and the
parked phase, stage, reservation and retry conditions shared by Week observation
and native release. Week's NativeAdapter no longer knows Worker diagnostic
prefixes. `WithParkedWorker` reads the current state under the owner lock and
still validates preparation identity/disposition, controls, bound attempt and
process completion, provider evidence and the exact native queue before release.
Advice observed earlier supplies no authority for that effect. Schema-1 Week
continues independent work after release; bound Week still requires recovery of
its retained Run. Stopped concrete failures retain native recovery support but
do not become automatic Week release advice.
Completed releases retain advice for reconciliation after a lost reply; the
native preparation disposition prevents another effect or a new writer.

This gives the Worker module depth: callers need one advisory decision instead
of knowing the concrete failure vocabulary and its lifecycle conditions. It
gives locality because outcome production, historical compatibility and native
effect checks stay with Worker. The deletion test would restore failure lists
and phase/stage/retry knowledge across both observation and release. No new
adapter was introduced; NativeAdapter and fakeNative remain the existing seam.
`legacyWorkerOutcome` alone interprets the historical native release prefixes
when `outcome` is absent. It neither migrates old records nor makes the old
Week-only `reviewer blocked:` prefix a supported native release outcome.

AR02 regressions exercise actual Worker interpretation and a disposable native
queue through NativeAdapter. They cover changed and empty diagnostics, explicit
unknown effects, fresh ownership/control/process/provider failures after advice,
one queue effect across repetition and a lost reply, retained source and resume
rejection after release. Historical reads preserve state/request bytes.
`test/worker-outcome-acceptance.sh` adds exact-candidate public Worker CLI checks
and runs the AR01, Worker, Week and queue regressions with modeled provider facts.

## Verification

`TestAR01ConsumedResumeIgnoresDisplayWording` and
`TestAR01AnswerBlockRecoveryIgnoresDisplayWording` were run red before production
changes. They failed through WorkflowShow/ShowWorker and Week Show respectively.
Additional behavioral cases cover renamed provider diagnostics and an ordinary
block with native-looking display text. Week recovery checks unchanged record
bytes and absence of native effects before Tick agrees with Show.

All AW01 and NA01 regressions remain. Product-answer tests now cross the same
`currentRunAction` seam as WorkflowShow instead of calling the removed partial
product-only interpretation. Historical Week fixtures explicitly omit the new
guidance origin. No test was deleted or replaced by a list of implementation
names. SL01 and MR05 coverage remains in the acceptance selection.

The unchanged baseline also failed the native single-attempt journey at its
old immediate-completion assertion. CL02 now first returns released delivery
with cleanup pending. The journey checks that `closing` boundary and its durable
delivery/release, then takes the existing next worker recovery step and still
requires completed delivery. This changes only the test, with no relaxation of
writer absence or cleanup policy.

`test/current-state-acceptance.sh` uses the exact private candidate for the
existing native SL01 public CLI journey: worker show, Week show, question and
answer status across receipt, continuation, review and completion. Its provider
is an owned stand-in, not evidence of live model, mobile or passkey behavior.
It runs focused behavior/control tests and relevant vet with readonly modules,
offline dependencies and physical disposable directories outside the workspace.
Native review, exact-candidate acceptance and observed integration are separate
delivery facts.
