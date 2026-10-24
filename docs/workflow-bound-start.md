# Named Goal startup with stable request identity (GS-01)

Version 1 supports an explicitly selected pending Goal with interactive Codex.
Discover support in `ply capabilities --format json`: require operation
`workflow.execute.bound-start`, mode `interactive_codex_explicit_goal_v1`, and
the version 1 request/preview/status schemas. The existing `.apply` and `.status` operations use the same request shape.
Ply upgrades require no compatibility step or historical preview check.

Run preview and apply from the registered Epic return worktree:

```sh
ply workflow execute preview --spec selected-goal \
  --request-id 550e8400-e29b-41d4-a716-446655440000 \
  --herdr-workspace w9 --format json > /absolute/preview.json
jq '.request' /absolute/preview.json > /absolute/approved-request.json
# The consumer can display .basis; existing authentication and authority still apply.
ply workflow execute apply --file /absolute/approved-request.json --format json
ply workflow execute status 550e8400-e29b-41d4-a716-446655440000 \
  --workspace /physical/workspace --format json
```

The consumer journals the request ID before apply. The selected inputs identify
the named Goal and return target. IDs contain 8–128 ASCII letters, digits, dots,
underscores or hyphens, beginning with a letter or digit. Preview is read-only;
it reserves no identity or effect. A historical preview checksum is optional
information and never a startup prerequisite.

The request shape is:

```json
{
  "kind": "ply.workflow.bound-start-request",
  "schema_version": 1,
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "workspace": "/physical/workspace",
  "input": {
    "spec_id": "selected-goal",
    "target": {"project_id": "", "repo_id": "", "epic_id": ""},
    "start_mode": "interactive",
    "herdr_workspace": "w9",
    "permission_profile": ""
  },
  "binding": "sha256:<64 lowercase hex digits from preview>"
}
```

The public [JSON Schema](../schemas/workflow-bound-start-v1.schema.json) defines
the strict request and the additive preview/status envelopes. The optional `binding` string is historical information. Unknown request
fields, duplicate keys, unsupported versions and modes are rejected. Empty
target fields select the current registered return worktree; all three explicit
target fields must be supplied together. Empty permission_profile selects the
existing provider configuration; the resulting profile is displayed and bound.

The preview `basis` contains the physical workspace and exact project/repo/Epic/
Task, Goal revision/hash and input files, queue/base revision, ref/OID/tree,
physical return and planned Task worktree/branch, model/effort, provider/Herdr/Ply
executable identities, requested launch policy, permission definition and policy
source hashes, HOME/CODEX_HOME/PATH, native trust behavior and unchanged delivery
and acceptance agreement. `runtime_acceptance` is explicitly `requested_only`.
The recipient still confirms its actual runtime and permissions before Task
writes. No persistent Codex configuration or trust is changed by this command.

Apply uses the currently registered named Goal/Spec and normal native
current-state authority. Changes since a historical preview do not require a
new preview or approval. This applies equally to old opaque requests with no
retained basis and new requests without a checksum. No installation history,
old binary or checksum reconstruction is needed.

The stable request identity is looked up first. Same ID and selected inputs
read the original attempt and Run, including after an informational checksum
change. Different selected inputs under that ID return
`bound_start_request_conflict`. Original stored request bytes remain unchanged.
Unknown or interrupted effects never cause a retry, new provider or repeated
Task prompt. A request that stopped before effects cannot adopt another
request's Run merely because both named the same Goal.

For a new identity, apply captures current state under native locks, then
publishes the intent before execution effects. The existing guard protects this
same startup's registry, project, source, runtime and launch environment. It
advances only after its own publications and rechecks at reservation and
transport boundaries. It compares against the current startup snapshot, not a
historical preview. Invalid current Goal/Task/return/authority or inconsistent
changes during startup still prevent the dependent effect.

Apply can operate from an app CLI process without `HERDR_ENV=1` because the
request names the Herdr workspace; existing authentication and native authority
remain required. The original schema and operation names remain supported; no
new schema or compatibility step is introduced. See
[PV-01 evidence and consumer inventory](ply-upgrades.md).

Status `kind` is `ply.workflow.bound-start-status`, `schema_version` is 1:

| State | Execution effect | Meaning |
| --- | --- | --- |
| `not_recorded` | `not_confirmed` | No intent was found. Absence alone is **not proof of no effect** after a lost answer or damaged storage. |
| `retained_attempt` | `unknown` | A durable request exists; no readable native Run is confirmed yet. It may be in progress or interrupted. Never start another writer to recover it. |
| `existing_run` | `reserved` or `unknown` | The same native Run exists. `run` is a historical inventory projection with Run ID, Herdr/session identities, freshness and unresolved evidence. It is not a provider liveness observation. |
| `outcome_unknown` | `unknown` | Required preserved evidence cannot be read or validated. Inspect this same identity; no automatic repair or new start. |
| `rejected` (apply) | `none` | This invocation rejected invalid current inputs before publishing its new intent. This does not describe another earlier invocation. |

Stable error codes are `bound_start_invalid_arguments`,
`bound_start_invalid_request`,
`bound_start_unsupported`, `bound_start_stale`, `bound_start_request_conflict`,
`bound_start_existing_attempt`, `bound_start_unavailable`, and `bound_start_outcome_unknown`. JSON operation
results carry `error.code` and `error.detail`; argument/file parsing errors also
use the stable code prefix on stderr and exit nonzero. Do not interpret exit
status alone as evidence that no start occurred: always journal the request and
use status after losing an apply result.

Initial bound start never adopts an earlier execution, startup recovery, or
unprepared launch intent under a new approval. Existing recovery commands remain
separate explicit operations, outside the original consumer approval.

Tests use disposable native Git/workspace fixtures and recorded provider/Herdr
stand-ins. Silver Bullet is only a read-only installed-binary preview/status
target for this delivery. Its real Slack start belongs to the app's separately
authenticated journey.
