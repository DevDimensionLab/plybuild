# Bound product answers (SL01, schema 2)

SL01 allows one fixed product choice to continue an explicitly opted-in,
week-owned Codex developer on the same Task, reservation and provider session.
Discover all four operations in `ply capabilities --format json`:

| Capability | Command | Effect |
| --- | --- | --- |
| `workflow.execute.product-question` | `workflow execute question` | Read |
| `workflow.execute.bound-answer` | `workflow execute answer preview` | Read |
| `workflow.execute.bound-answer.apply` | `workflow execute answer apply` | Approved product answer |
| `workflow.execute.bound-answer.status` | `workflow execute answer status` | Read |

Require mode `week_codex_same_owner_turn_v1` and schema version **2** on every
operation. An advertised operation describes the binary's contract; the Run's
question and preview determine its eligibility. Old Runs are not upgraded.
Schema 1 and `week_codex_same_thread_v1` are not compatible aliases.

Before starting a new schema 2 week, its complete JSON contract must explicitly
include both opt-ins:

```json
{
  "runtime_control": {"schema_version": 1, "mode": "bound_control_v1"},
  "product_questions": {
    "schema_version": 1,
    "mode": "week_codex_same_owner_turn_v1",
    "max_continuation_turns": 1
  }
}
```

This is a fragment, not a complete start contract. Validate the complete
selected contract with `ply workflow week start --file /absolute/week.json
--check --format json` before the separately authorized start. The native
week binds the corresponding schema 2 worker contract before its first Run.
No standalone opt-in or upgrade of an existing week is supported by SL01.

For a new schema 2 week with `product_questions`, `max_provider_retries: 0`
means zero additional Ply-managed provider retries; the corresponding schema 2
worker preserves this limit. `max_attempts: 1`, `max_provider_retries: 0` and
`max_continuation_turns: 1` allow the initial developer turn, one approved
same-owner answer continuation and a separate reviewer turn. The latter two
are not extra developer attempts. The legacy schema 1 worker's zero-to-20
retry default remains unchanged. These limits do not attest Codex-internal
transport attempts.

## Read, preview, approve and apply

Use the physical workspace, not a Task worktree or a cwd-derived path. Read
operations create no intent, lock file, wake, provider turn or reservation.

```sh
ply workflow execute question RUN_ID \
  --workspace /physical/workspace --format json
ply workflow execute answer preview RUN_ID \
  --workspace /physical/workspace \
  --question-id QUESTION_ID --question-revision 1 \
  --option-id task_first --request-id answer-choice-0001 --format json
```

The consumer journals its request ID before approval, displays the exact
question, fixed choice, consequence, scope and continuation, and obtains its
own fresh user authentication. Preserve the complete returned `request`
object unchanged. The binding does not authenticate the user. The app owns
Slack identity and passkey verification; the native transaction independently
rechecks the approved basis.

```sh
ply workflow execute answer apply \
  --file /absolute/approved-answer.json --format json
ply workflow execute answer status answer-choice-0001 \
  --workspace /physical/workspace --format json
```

Only a ready preview has a non-null request. Request validity is explicit:
`validity.issued_at_utc <= now < validity.expires_at_utc`, with at most five
minutes and never beyond `original_deadline_utc`. Observation time is not an
expiry. Replaying an already retained exact request remains a read even after
expiry. A changed request under the same ID is a conflict; a lost reply never
authorizes a new ID, another apply or another developer session.

The request digest is lowercase SHA-256 prefixed by `sha256:` over canonical
JSON of the entire request, including mode, validity and binding. Property
order and insignificant whitespace do not change it. Question digests omit
only the question's own `question_sha256` field. Use the published canonical
vectors rather than hashing a pretty-printed file.

## Interpret the result without inventing progress

Apply and status use `ply.workflow.bound-answer-status`, schema 2. A semantic
rejection can return a valid JSON status and a nonzero exit. Preserve that
status. Argument, file and parser errors have stable
`bound_answer_invalid_arguments` or `bound_answer_invalid_request` codes on
stderr and no invented stdout status. Exit status alone does not prove that
an earlier invocation had no effect.

| State | Meaning |
| --- | --- |
| `not_recorded` | No readable request; unknown effect, not proof of no effect. |
| `answer_recorded` | The answer won durably; no continuation is reserved yet. |
| `continuation_deferred` | The retained answer waits before reservation for a named parent or runtime condition. |
| `continuation_reserved` | One exact continuation intent exists; start is unconfirmed. Parent pause after reservation preserves this state and identity. |
| `continuation_observed` | A new process attempt and the same provider session were positively observed. Answer use is not yet confirmed. |
| `continuation_completed` | The completed turn names the exact consumed answer ID and digest. This does not attest integration. |
| `outcome_unknown` | Evidence is incomplete, contradictory or unreadable; retain the same identity. |
| `rejected` | This invocation published no new intent. This says nothing about older invocations. |

The status basis is the preserved approved scope, worker and parent, not a
new liveness observation. Current known blockers appear in `reasons` and
evidence. An old `owner_state=confirmed` is not proof that an owner is alive.
Unknown identity fields and effect booleans are explicitly `null`; consumers
must preserve these nulls rather than translating them to false.

`continuation_observed` only permits “Work has started again. Answer use is
not yet confirmed.” “Your answer was used” additionally requires
`continuation_completed`, matching consumed request ID/digest, same source and
continuation session, and bound terminal result evidence. Delivery of the
selected content requires separate test, review and integration evidence for
the exact source revision. App-main integration is an app-owner receipt, not
an answer-status claim.

`continuation_started_at_utc` preserves the launch-boundary marker unchanged
in observed and completed receipts; it appears only after positive process
and session observation and is not an exact OS process-start timestamp.

## Same owner and parent control

`awaiting_product_answer` preserves the same Task, worker authority,
reservation, developer iteration, session, counters and original deadline.
It starts no model turns while waiting. It means a positively completed model
turn with the same owner; it does **not** mean global pause, all tool writers
quiescent, release eligibility or full MR04 v1 qualification.

The native transaction uses the existing parent control guard. A ready
preview under pause explicitly binds
`deferred_until_explicit_parent_resume`; answering does not resume the parent.
Pause after reservation keeps the exact intent and explains the blocker.
Stop, deadline, unavailable authority or an unknown previous effect prevents
a new dispatch. Status never repairs or wakes the controller. Each accepted
answer permits at most one continuation dispatch; no retry, new session or
budget reset is implied by a process restart.

The parent hashes can be cross-checked through public week readback:
`parent_week.contract_sha256` equals `week.control_contract.contract_digest`
(the child control contract's digest). `parent_week.workset_sha256` equals
`week.contract_sha256` and `week.control_contract.week_binding`: the hash of
the whole immutable week contract, including `selected`. It is not a separate
workset authority or a hash of only the selected list.

## Published wire and limits

The [schema](../schemas/workflow-bound-answer-v2/workflow-bound-answer-v2.schema.json)
has four closed roots: `question_read`, `preview`, `request`, and `status`.
The [examples](../schemas/workflow-bound-answer-v2/examples.json) contain valid,
schema-invalid, raw-invalid and semantic-invalid vectors. The
[locked wire rules](../schemas/workflow-bound-answer-v2/README.md) define exact
nested fields, nulls, time ordering and equality checks.

Requests are at most 64 KiB UTF-8; individual responses are at most 256 KiB.
Unknown or missing fields, duplicate keys, BOM, invalid UTF-8, lone surrogates,
multiple JSON values, floats, exponents and negative zero are rejected.
Paths must be absolute and physically normalized, with no symlink or dot
segments. Request strings are data, never shell text. No app-supplied prompt
or resume argv is accepted.

The wire directory preserves the exact reviewed design bytes, including its
historical draft labels. The frozen SHA-256 values are:

- README: `8b972f65ce0b4ea72299d802309cc3c5ba254874d2460007af8b5bdfcc996a41`
- Schema: `bba59f011f8d1217cd183304ecdb14a9311801289e3102e69aa91f8feb977075`
- Examples: `6f994001df8bac2f054c0014a4d2b68dbe32a08fc3de9835f4e99784065d3dd0`

Those bytes and synthetic examples do not attest a live Run or installation.
An ordinary installation and a separately bound week/app qualification must
identify the actual source, binary, tests, review and observed runtime effects.
