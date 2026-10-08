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

## Provider upgrades and current runtime observations

The provider executable in the request is **historical launch evidence** once
the session has accepted its mandate. Removing, relocating or replacing that
file does not identify the current process, change its accepted provider, or
authorize launching a different executable. The preview shows its original
binding and whether that file is present, missing or changed/unavailable.
It does not search for a replacement and never executes either launcher.

Herdr's supported observation supplies the exact registered provider, native
session, workspace/tab/pane/terminal, foreground Task cwd and known status
(`working`, `idle`, `done` or `blocked`), with no pending launch. It supplies no
effective-policy claim. Therefore, when the historical launcher is unavailable
or changed, the owner must also preserve a private current runtime observation
and pass `--runtime-evidence FILE` to `workflow execute continue`. The owner
reads its actual current runtime/session/policy; copying requested policy or
looking up an installed binary cannot produce that observation. Unknown or
changed permissions stop continuation. This is the same recipient-observation
trust boundary as initial acceptance, not OS attestation or renewed permission.

The strict JSON contract is `ply.workflow.delivery-runtime-observation@1`:

| Field | Required observation or binding |
| --- | --- |
| `kind`, `schema_version` | `ply.workflow.delivery-runtime-observation`, `1` |
| `run_id`, `request_sha256`, `session_id` | The original native run, immutable request and Ply session |
| `context`, `acceptance` | `{locator, sha256}` for the exact preserved context and acceptance |
| `runtime_claim` | Actual `runtime_id`, `model_id` (null if unknown), `profile_id` (nullable), `native_session_id`, `effective_policy_sha256` |
| `delivery_permission` | Actual `permission_confirmed`, bound `launch_contract_sha256`, `actual_policy_evidence`, and original `delivery_agreement_sha256`/`allowed_effects` where present |
| `observed_at_utc` | Actual RFC3339 UTC observation time, within ten minutes of a new transition, with at most thirty seconds of clock skew |

The current effective-policy digest, evidence bindings, session, original
agreement and effects must match accepted authority. Evidence uses the existing
locator-sorted `{locator, sha256, role}` contract. Preserve submitted files and
their evidence unchanged. New observations use new private files. A different
actual profile label alone is not a permission change. No automatic collector
pretends that Herdr can observe policy; the owning agent must make the claim
from its current runtime context. A later authority change must be reported and
stops dependent effects under that same recipient contract.

The continuation preserves this observation in its immutable proof. Subsequent
callbacks validate that proof, accepted policy and fresh live owner, allowing
the old launcher to remain absent. Every new launch and initial acceptance still
checks its applicable bound executable. Older continuation proofs remain
readable; they require this explicit transition before omitting launcher checks.

| Operation | Required runtime dependencies |
| --- | --- |
| Start, initial acceptance | Bound provider launcher, control, Herdr, launch contract and applicable acceptance checks |
| Upgrade continuation | Positive accepted authority, fresh owner policy observation, live Herdr owner, original/current controls, required frozen artifacts |
| Working report, verify/reuse, QA, authorized delivery after continuation | Preserved continuation/current-policy proof, accepted evidence, live owner, bound control and required artifacts; no historical launcher |
| Installed `report --incomplete` | Exact run/context/Task, frozen acceptance and controls/history, report evidence and live Herdr owner; no launcher, requested/current policy-file or trust-config dependency |

`report --incomplete` accepts only `stopped` or `needs_input`, including before
acceptance. It cannot qualify a candidate, record QA, switch controls or authorize
work. It preserves unresolved attempts and the existing publication predecessor.
Missing ownership observation or altered immutable identity artifacts still
reject reporting. A durable report retains its exact human question or stop
reason in readback; runtime drift remains a separate diagnostic and invalidates
technical eligibility. A local file or rejected command is not a recorded report.
After an uncertain reply, inspect the same event ID and bytes before retrying.

Continuation recovery first observes a preserved transition proof. A proof
published while its observation was fresh can be adopted after a lost reply,
without expiring that already-performed transition or overwriting its before
state. Current immutable evidence and live ownership are still rechecked, and
intervening reports remain in the current state. Concurrent/repeated calls bind
one generation and one control. An equivalent current observation reaffirms the
existing proof for the same control; it does not create another generation.
Requalification retains its original attempt
key, receipt and exact candidate conditions.

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
| ac-01 | `TestDeliveryProviderUpgradePreservesAcceptedContinuation` covers removed/relocated/changed launchers for both providers; `TestDeliveryProviderUpgradePreservesNativeLegacyContinuation` extends a real legacy proof. |
| ac-02 | Current-runtime observation, authority/live-owner and initial acceptance/new-launch rejection tests keep current authority separate from launcher history. |
| ac-03 | `TestContinuityNativeOldReaderCorrection` exercises working reports, receipt reuse, synthetic exact-candidate QA and local integration after launcher removal. |
| ac-04 | Incomplete-report tests cover pre/post-acceptance missing runtime evidence, exact questions, immutable identity guards, interrupted publication and retained unknown verifier attempts. |
| ac-05 | `TestDeliveryReceiptReuseRejectsChangedInputsWithoutRunningAcceptance`, correction status tests and native upgrade scenarios retain distinct verification, qualification and human judgment. |
| ac-06 | Missing-current-observation, changed policy, missing/wrong live session and required-artifact guards produce distinct native diagnostic codes and causes. |
| ac-07 | Concurrent continuation/report/requalification tests and recovery after expired evidence, intervening reports and partial proof publication preserve one transition and unrelated records. |
| ac-08 | `TestContinuityNativeProviderUpgradePreservesAcceptedOwner` preserves actual failures from pinned historical binaries, then verifies the new native commands with removed/replaced/relocated provider fixtures. |
| ac-09 | `scripts/continuity-journey`, its native driver/retry regressions and [the installed acceptance journey](../continuity-acceptance.md) exercise short commands without manual state edits. |

Automated fixture verdicts are labeled synthetic. This implementation's actual
human pass remains an external, exact-candidate gate after technical qualification.
