# Active delivery continuity

## Decision and contract boundary

An active delivery is bound to its immutable request, mandate, native session,
accepted runtime authority, Task worktree and selected delivery agreement. An
upgrade does not change those facts. Publication bytes are content addressed;
newer unrelated metadata is not evidence that older bytes were corrupted.

Use two boundaries: operation-scoped semantic validation of Task content, and
an explicit native continuation of the same accepted delivery using a separately
preserved compatible control. Preserve the old control, request, context,
acceptance, receipts and candidate history byte for byte. Never restart the
provider, send another Task prompt or roll back another Task's publications.

Alternatives considered:

- Relax every decoder: rejected because required unknown semantics could grant
  effects the reader does not understand.
- Freeze or copy the whole workspace registry: rejected because this hides
  concurrent publications and loses updates.
- Replace the old bound executable: rejected because it breaks its integrity
  binding and provenance.
- A permanent compatibility daemon: unnecessary for this bounded transition;
  native commands and existing file locks can preserve its evidence.

## Required dependency set

The shared registry retains strict structural validation and the complete
records. A Task operation validates semantic content for that Task and its
explicit dependency closure: selected and historical Spec/problem/assessment/
selection publications, referenced goal origins and documents, result-bound
Specs, relevant Epic queue/base/preparation and integration authority. Current
contracts prohibit cross-Task content dependencies; a future required dependency
needs explicit reader support. Every required artifact retains physical-path,
digest, size, schema and chain checks. Shared registry identity/reference checks
remain global; unrelated publication bytes are preserved without semantic or
artifact traversal. Mutation reads and publishes the complete latest registry
under the existing lock; it never
writes a filtered snapshot. Global audit operations can still validate all
content.

Unknown required fields/versions are incompatibility, malformed bytes are
invalid content, mismatched digests are integrity failures, mismatched candidate
or predecessor is stale state, and a reserved effect without a conclusive
receipt is unknown. Diagnostics retain the cause, artifact, operation, reader
capability and supported next action. No error grants permission.

## Continuing an already-started run

The installed native continuation command observes the same accepted run in the
same physical Task cwd and the same live provider/session. It checks the original
control digest, immutable inputs, accepted policy evidence, agreement and native
start receipt. It proves that this reader supports all required contract
semantics and preserves that proof with old/new executable and context bindings.
It copies its executable to a new immutable control slot and binds a continuation
proof in the mutable run state, retaining the exact original context. The frozen
request, context, prior controls and accepted policy stay intact. Only then may
callbacks select the new slot; old controls cannot decode the transitioned state
and must use the installed continuation command to inspect its supported path.
This transition adds no delivery effects or authority and sends no provider input.
Unsupported contracts or changed authority stop with a precise supported action.
Repeated/concurrent requests must observe and return the same transition.

Routine compatible continuation adds zero human actions: the delivery owner uses
the supported command. A real authority change remains a separate explicit
bounded decision. Neither a directory, installed version nor preserved requested
launch profile establishes effective runtime authority.

## Verification and distinct outcomes

Preserve command completion independently from qualification. A known completed
verification whose qualification failed may be requalified without rerunning the
command only when the exact attempt, source commit/tree, original review,
acceptance script and snapshot, outputs, frozen inputs and relevant accepted
runtime evidence still match. Unknown execution without a receipt is observed,
never replayed. Candidate publication uses its original deterministic attempt key
and existing native idempotent artifacts. A changed source or input requires new
verification and reports why the old receipt cannot be reused.

Readback separately presents current source, the exact last verification and its
exit, qualified candidate history/currentness, human judgment and observed final
delivery. An older qualified candidate must not label a correction ready. Human
QA requires a real pass/fail/blocked answer for the exact qualified candidate;
only that pass plus the unchanged original authority permits integration.

## Validation and human journey

Use isolated workspaces and an actual pre-delivery-schema reader built from
commit f2de0c638f1d2b1c38417aff70b79ca7bb4be31a. Qualify an earlier legacy candidate,
publish unrelated newer Goal/Spec bytes including delivery, commit a correction,
and reproduce successful acceptance followed by old-reader qualification failure.
Continue through the installed native command, reuse only matching evidence,
then inspect all distinct outcomes. Exercise tampering, required incompatible
semantics, source/review/script mismatch, interrupted receipts, concurrent
publication/continuation and duplicate callbacks. Existing QA and delivery tests
remain authoritative for their separate gates.

The repeatable installed-product journey uses the same isolated fixture and
short commands rather than internal hash copying or manual .ply changes. Its
provider and historical QA fixtures are explicitly synthetic. The user's actual
judgment of this product candidate is recorded separately by its delivery owner.
The live passkey run remains read-only evidence owned by its original agent.

## Compatibility policy and limits

This implementation supports the known delivery-owner v2 run and its native
handoff/receipt contracts, with the pre-delivery-field schema as the historical
regression boundary. It does not promise arbitrary future run/registry schema
support. Required incompatible semantics need a supporting reader; unknown
execution effects require observation. It adds no unattended loop, remote
publication, permission changes, or human approval inferred from automation.

## Operating the supported transition

The delivery owner uses the installed CLI from the existing Task worktree:

```shell
ply workflow execute continue RUN_ID --context ORIGINAL_CONTEXT --check
ply workflow execute continue RUN_ID --context ORIGINAL_CONTEXT
```

The first command is optional readback. The second preserves the compatible
control and returns its exact path. Subsequent callbacks use that path and the
unchanged context. The old executable remains preserved. For the latest
successful verification whose qualification failed, use the returned control:

```shell
/returned/ply-control workflow execute verify RUN_ID \
  --context ORIGINAL_CONTEXT --review ORIGINAL_REVIEW --reuse ATTEMPT_ID
```

The native result supplies these identities; the owner does not ask the human
to repair them. A mismatched receipt returns a reason without running acceptance.
The ordinary verify command is a deliberate new execution after changed inputs.
Human QA and integration still use their original separate callback gates.

## Acceptance coverage

| Requirement | Executable evidence |
| --- | --- |
| ac-01, ac-08 | `TestContinuityNativeOldReaderCorrection` uses the pinned historical CLI and real newer publication bytes; `TestTaskScopedContentPreservesUnrelatedNewerPublications` covers future unrelated semantics. |
| ac-02 | Scoped required-artifact tests, continuation authority/session/control rejection, and `TestDeliveryReceiptReuseRejectsChangedInputsWithoutRunningAcceptance` preserve the blocking boundaries. |
| ac-03 | Native historical continuation and `TestDeliveryContinuationPreservesAcceptedSessionAndFrozenInputs` retain the original owner and frozen artifacts. |
| ac-04 | Receipt requalification and status tests preserve actual command completion independently of changed inputs or failed qualification. |
| ac-05 | Concurrent native continuation/publication, interrupted proof publication and concurrent requalification tests assert one effect and retained sibling records. |
| ac-06 | Native Task diagnostics and continuation diagnostics distinguish unsupported schemas, missing/tampered artifacts, stale sessions and unknown effects. |
| ac-07 | `TestDeliveryStatusSeparatesCorrectionFromEarlierCandidate` and existing exact-candidate QA/integration tests keep technical, human and final outcomes separate. |
| ac-09 | `scripts/continuity-journey`, its native driver regression and [the installed acceptance journey](../continuity-acceptance.md) exercise the same commands without manual state edits. |

Automated fixture verdicts are labeled synthetic. This implementation's actual
human pass remains an external, exact-candidate gate after technical qualification.
