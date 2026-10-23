# Agent work and follow-up v1 (T7CF442AE)

Early concrete producer contract, subject to the exact reviewed schema. Commands:

```
ply workflow agents --format json
ply workflow followup register --file request.json --format json
ply workflow followup list [--due] --format json
ply workflow followup show FOLLOWUP --format json
ply workflow followup request FOLLOWUP --file mutation.json --format json
ply workflow followup dispatch FOLLOWUP --format json
ply workflow followup recheck FOLLOWUP --file mutation.json --format json
ply workflow followup cancel FOLLOWUP --file mutation.json --format json
ply workflow followup retarget FOLLOWUP --file mutation.json --format json
```

All operations accept an explicit physical `--workspace`. Reads create nothing.
`AgentWorkReadback@1` contains `observed_at_utc`, `sources`, `agents`, and
`followups`. Each agent has an opaque `agent_ref`, safe provider/status/tab/workspace
metadata, nullable verified `work`, binding status, and operation capabilities
(`enabled`, `disabled`, `unknown`, each with a reason). Raw terminal/session IDs,
paths, prompts and source errors are never copied into public output.

`AgentFollowup@1` has `id`, `revision`, `generation`, `owner_ref`, `mandate_ref`,
`registrar` (kind `human_requested` or `agent_self_registered`, safe actor reference),
`target_ref`, `wait_kind`, `wait_ref`, nullable verified `work`, `next_check_at_utc`,
`state`, `cycle`, `transport`, and nullable `recheck`. The server authenticates the
human and authorizes calls on its private owner socket or Slack service; CLI
registrar provenance is an explicit caller declaration, never inferred human QA.

Register request uses `publication_key`, `owner_ref`, `mandate_ref`, `registrar`,
`target_ref`, `wait_kind`, `wait_ref`, `next_check_at_utc`. Stable ID derives from
owner/mandate/wait identity, excluding aliases and target generations. Retrying the
same bytes reads back; changed bytes conflict. `list` recovers a lost response.
Mutations use `request_id`, `owner_ref`, `expected_revision`, `generation`,
and the operation's arguments. Recheck additionally requires `target_ref` and
`outcome` (`waiting`, `fulfilled`, `cancelled`); waiting supplies a new explicit
`next_check_at_utc`. Recipient revalidation uses `show` before reporting.

Dispatch claims one canonical cycle durably before any input. A consumed cycle
cannot be replayed with a new request ID, alias, restart or target generation.
Sent means transport acknowledged/queued, not read. Lost ACK or crash after claim
means unknown/possibly sent; only a recipient's explicit waiting recheck creates
a new cycle. A request makes an unconsumed cycle due, without replaying a consumed
cycle. Retarget requires a new generation and preserves its consumed state.

An opt-in server loop calls `list --due`, then `dispatch ID` once per listed item.
Each invocation is bounded and starts no background worker, model or session.
Working defers; missing identity/readiness, pending launch, known native
question/pause/refusal, unknown veto coverage, cancellation or supersession prevent
input. Herdr get → prompt is best effort, not atomic. A recipient must revalidate
ID/generation/target/state before status work. Fixed message grants no authority.

Status reads include unbound own Herdr agents. Observation is fresh only from a
new bounded Herdr read; native work and next actor come from registered sources.
Agent done never establishes Task delivery. Source outcomes distinguish available,
empty, partial, unavailable and unknown. Public metadata is bounded to 4096 agents
and 256 UTF-16 units per field. Native registry format 4 is unchanged.

Consumer migration keeps HTTP status-v1, Tailscale identity, Firebase CORS and
4096/256 limits. The server uses these CLI contracts without importing Go packages
or private store files. TUI owner-socket and Slack share its service. No HTTP/PWA
actions or application activation are part of this producer Task.

Representative draft fixtures: `schemas/agent-work-v1/examples/overview.json`.
