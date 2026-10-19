# Current Run, Worker and Week interpretation

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
