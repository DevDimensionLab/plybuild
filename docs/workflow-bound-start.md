# Approved preview-bound Goal start (GS-01)

Version 1 supports an explicitly selected pending Goal with interactive Codex.
Discover support in `ply capabilities --format json`: require operation
`workflow.execute.bound-start`, mode `interactive_codex_explicit_goal_v1`, and
the version 1 request/preview/status schemas. Also require its `.apply` and
`.status` operations. An older or unsuitable binary must be rejected by the
consumer. Never fall back to ordinary `execute`.

Run preview and apply from the registered Epic return worktree:

```sh
ply workflow execute preview --spec selected-goal \
  --request-id 550e8400-e29b-41d4-a716-446655440000 \
  --herdr-workspace w9 --format json > /absolute/preview.json
jq '.request' /absolute/preview.json > /absolute/approved-request.json
# The consumer displays .basis and obtains its own fresh user authentication.
ply workflow execute apply --file /absolute/approved-request.json --format json
ply workflow execute status 550e8400-e29b-41d4-a716-446655440000 \
  --workspace /physical/workspace --format json
```

The consumer supplies and journals the request ID **before apply**. It must keep
the entire request object unchanged. IDs contain 8–128 ASCII letters, digits,
dots, underscores or hyphens, beginning with a letter or digit; UUIDs work.
Preview does not reserve the ID or create any files, native intent, reservation,
worktree, Run, trust grant, agent or tab. Status performs only historical reads.

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
the strict request and the additive preview/status envelopes. Unknown request
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

The opaque binding covers the shown basis and request except historical Ply executable
identity. Replacing or removing that Ply installation leaves a new preview valid;
provider, Herdr, mandate, policy, environment and delivery inputs remain bound. It is a
comparison token, not a user authentication credential or new Git/delivery
authorization. Plyapp retains WebAuthn, Slack identity and user confirmation.

The native request key and resulting Run ID bind both the consumer request ID
and its approval binding. A different request that stopped before effects cannot
adopt another request's Run merely because both selected the same Goal.

Apply first looks up the existing identity. Same ID and same request read the
same attempt even after its own queue/trust/Run state or installed binary
changes. Different request content or binding under the same ID returns
`bound_start_request_conflict`. There is no implicit retry/resume/restart.

For a new identity, apply rechecks the complete basis while holding the existing
native Project and WorkItem locks, then durably publishes the immutable intent
before execution effects. The guarded native composition verifies its expected
registry, project and runtime before content and worktree writes, advances the
expected registry only after its own publications, and checks again under the
native TaskRun → Project → WorkItem locks at Run reservation. Transport checks
the bound environment again before the tab and agent effects. A valid request
can therefore leave a retained attempt when a later conflict stops progress;
it never silently selects another revision or weaker policy.

The first Run preserves its historical Ply observation and exact launch environment. Bound
apply can operate from an app CLI process without `HERDR_ENV=1` because the
request names the approved Herdr workspace; there is no public bypass switch.
Ordinary interactive execute keeps its terminal requirement and history.
Headless contracts retain their authority and use the current Ply invocation.

Legacy previews made before PV-01 included Ply in the opaque digest but returned
no prior basis in the request and saved nothing on preview. If that basis was
never retained, a changed digest cannot be distinguished from a changed mandate
or policy. Such a request remains rejected without effects; arbitrary old
digests are never accepted. Existing retained intents, requests, contexts and
Run identities remain readable without rewriting history. See
[PV-01 evidence and consumer inventory](ply-upgrades.md).

Status `kind` is `ply.workflow.bound-start-status`, `schema_version` is 1:

| State | Execution effect | Meaning |
| --- | --- | --- |
| `not_recorded` | `not_confirmed` | No intent was found. Absence alone is **not proof of no effect** after a lost answer or damaged storage. |
| `retained_attempt` | `unknown` | A durable request exists; no readable native Run is confirmed yet. It may be in progress or interrupted. Never start another writer to recover it. |
| `existing_run` | `reserved` or `unknown` | The same native Run exists. `run` is a historical inventory projection with Run ID, Herdr/session identities, freshness and unresolved evidence. It is not a provider liveness observation. |
| `outcome_unknown` | `unknown` | Required preserved evidence cannot be read or validated. Inspect this same identity; no automatic repair or new start. |
| `rejected` (apply) | `none` | This invocation rejected its input/basis before publishing its new intent. This does not describe another earlier invocation. |

Stable error codes are `bound_start_invalid_arguments`,
`bound_start_invalid_request`, `bound_start_invalid_binding`,
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
