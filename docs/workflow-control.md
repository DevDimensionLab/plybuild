# Request-bound runtime control

`workflow control` lets a client inspect, approve and follow one exact pause,
stop or resume request. It supports opted-in native Codex headless workers on
macOS, either standalone or owned by an opted-in week. Interactive Herdr runs,
contracts without control opt-in and unknown ownership return no actions. Installing a
new CLI does not change an existing Run's frozen contract.

## Opt in when starting new work

For a new standalone Goal, inspect the complete start agreement first:

```sh
ply workflow execute --headless --bound-control --spec SPEC_ID --check --format json
ply workflow execute --headless --bound-control --spec SPEC_ID --format json
```

For a new week, use the [week contract](workflow-week.md) with
`"schema_version": 2` and the additional field
`"runtime_control": {"schema_version": 1, "mode": "bound_control_v1"}`.
`workflow week start --file /absolute/week.json --check --format json` shows the
control agreement before start. Version 1 and starts without this opt-in retain
their existing behavior. The native `control_contract` binds parent kind,
contract digest, supported schemas and, for a child, the exact week ID and
binding. It is separate from the worker's execution-policy version.

## Inspect and approve the exact scope

```sh
ply capabilities --format json
ply workflow control show RUN_ID --workspace /physical/workspace --format json
ply workflow control show --week WEEK_ID --workspace /physical/workspace --format json
ply workflow control preview RUN_ID --action pause --request-id mobile-control-001 \
  --workspace /physical/workspace --format json
ply workflow control preview --week WEEK_ID --action resume --request-id mobile-control-002 \
  --workspace /physical/workspace --format json
ply workflow control apply --file /absolute/approved-control-request.json --format json
ply workflow control status mobile-control-001 --workspace /physical/workspace --format json
```

Capabilities advertise separate `codex_headless_standalone_v1` and
`codex_headless_week_v1` modes for show, preview, apply and status. Show, preview
and status are read-only: they do not register work, create locks, recover
controllers or wake providers. Apply is `approved_runtime_control`.

A Run locator owned by a week resolves to that exact **whole week and selected
workset**. Present its dependencies, absolute deadline, active Run and each
action's boundary, effect and caveats before approval. `--week` also works
between Goals; Task and Run identities are then `null`. Run and week locators
are mutually exclusive. There is no implicit next item or alternate scope.

Save the preview's complete `request` unchanged, obtain the client's approval,
then pass that request object to apply. A rejected preview has `request: null`
and cannot be approved. Binding is a comparison token, not authentication:
the client still owns user identity, access, revocation and one-use approval.
Never substitute a start or answer receipt, shell command or legacy worker
command when this capability is unavailable.

Requests use strict JSON: duplicate or unknown fields, wrong types, unsupported
versions and noncanonical targets are rejected. IDs contain 8–128 ASCII letters,
digits, dots, underscores or hyphens, starting with a letter or digit. Preview
validity is at most five minutes. Apply revalidates the frozen identities,
semantic state, queue authority and validity under native publication locks.
Timestamps and log growth do not renew approval; changes of owner, work or
permission require a fresh preview. Returning to an earlier phase does not
reuse an earlier revision.

## Read request outcome separately from runtime observation

| Readback | Meaning and next step |
| --- | --- |
| `accepted` / `pending_boundary` | The exact request is durably reserved. Keep observing its ID; the boundary is not yet confirmed. |
| `accepted` / `paused` or `stopped` | Positive native evidence confirms parking. Work and identity remain available for a separately approved resume. |
| `accepted` / `resumed` | The native owner consumed the continuation permission. This does not prove provider liveness or completed work. |
| `accepted` / `completed_before_boundary` | Standalone integration finished first; the receipt does not claim a pause. |
| `accepted` / `superseded` | The original request was superseded by the absolute deadline's stop sequence; inspect the preserved reason. |
| `not_recorded` / `not_confirmed` | No retained request was found. This alone cannot prove that an earlier uncertain invocation had no effect. |
| `outcome_unknown` / `unknown` | Preserve the same request and inspect its retained evidence locally. Do not issue another ID to retry an uncertain effect. |
| `rejected` / `none` | This invocation did not reserve a control. Resolve its stable error before obtaining new approval. |

Identical replay returns the retained request even after expiry or owner changes;
it never wakes again. Reusing its ID with changed content gives
`control_request_conflict`. Concurrent requests cannot overwrite a pending
request. Status is historical readback, with explicit source/freshness; a receipt
or saved PID is not proof that a provider is alive. Public output contains opaque
revisions and identities, not private callbacks, streams or authentication data.

`runtime_observation.state` is one of `running`, `waiting`, `paused`, `stopped`,
`completed` or `unknown`. Internal `phase` stays separate: `ready` and `blocked`
project to `unknown`, `pausing` and `stopping` remain `running`, and `expired`
projects to `stopped` only at the parent's confirmed parking boundary. An absent
week or unreadable state returns a nonzero exit and safe JSON `error.code` /
`error.detail`, no actions, and `state_revision: null`.

Pause and stop wait for the active or reserved attempt's positive result and
process boundary. They send no signals or terminal input. An already reserved
integration can finish. A week pause that wins before the next selection still
parks the week after that completion, with a null active Run and no next Goal.
An explicit resume preserves the Run, logical worker session, agreement,
selection, dependencies, provider wait and original absolute deadline. Starting
the next selected Goal after resume creates a new Run for that Goal.

Blocked/authentication states, missing attempt evidence, runtime drift, released
queue ownership and unsupported contracts expose no resume action. Follow
`unavailable_reason` and `next_action` locally. Missing processes do not establish
standalone ownership. A child cannot clear its parent's pause or extend time.
Expiry prevents dispatch immediately; an accepted pending resume receives a
`deadline_expired` supersession receipt referencing the native stop sequence.

## Native publication and recovery

The lock order is **parent week control/dispatch guard → TaskRun store → Project
registry → WorkItem registry** when those resources are involved. Standalone
control starts at TaskRun. Apply rechecks native registry revisions inside the
publication transaction. Worker dispatch takes the parent guard before TaskRun;
no code takes the parent guard while holding TaskRun or registry locks.

The week owner/epic locks and worker owner/supervisor locks are long-lived process
ownership capabilities. Parent observation probes child locks without blocking;
the parent guard never waits for a child owner. Provider starts, controller wakes
and launchd waits happen after releasing the short publication guard. A durable
dispatch reservation counts as active even before spawn, so a winning pause
cannot fall into a read-before-spawn gap. Reconciliation reloads guarded state.
Owner-generation and failure publication use the same guard as apply.
Intent publication and Job reservation share the winning dispatch guard. A pause
that wins first leaves no intent to obstruct a later same-Run successor. A
partially published or reserved attempt remains unknown and is never overwritten.
An existing supervisor observes bound pause/stop during provider backoff without
changing `RetryNotBefore` or starting another provider attempt.

Each immutable request/digest/sequence is stored before wake. The original owner
may consume it, or one successor recorded against that reservation may do so;
the receipt includes `predecessor_generation` for that successor. Unlinked owner
changes preserve an unknown outcome. A completed native state write retains the
receipt for recovery if immutable receipt publication was interrupted. Public
replay and status perform no recovery. Corrupt readback after a known durable
reservation cannot turn a possible effect into `none`.
Consumers revalidate the exact child/session/job, executable bytes and required
positive result/stream evidence before opening a resume gate or confirming a
boundary. Terminal consumption rejection is retained against the request and
stops automatic owner retries. Later evidence does not reverse that rejection.
Recovery of an already saved positive receipt is separate and preserves its
original generation and observation, even when a later owner publishes it.

The week marks operational execution before consuming controls, including before
the first Goal. Missing mutable state after that point requires local recovery.
Only an unstarted registration gap can reconstruct its initial gate/state from
the immutable binding, retaining the original deadline.

See [JSON Schema](../schemas/workflow-control-v1.schema.json) for required fields
and explicit nullability, and [examples](../schemas/workflow-control-examples)
for valid and invalid documents. Output permits additive fields within version
1. The focused `TestMR05` tests cover native reservation, evidence loss, ownership,
dispatch races and CLI/schema contracts in disposable fixtures. Actual provider,
service and consumer journeys are separate evidence.
