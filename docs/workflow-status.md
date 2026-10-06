# Workspace workflow status

`ply workflow status` shows registered work across all projects and Epics in the
discovered Ply workspace, with an explicit next actor and evidence sources. Run
it from any directory inside the workspace; an Epic worktree is not required.

```sh
ply workflow status
ply workflow status --project ply --repo ply --epic cli
ply workflow status --all --json
ply --json workflow status
ply workflow status --format json
```

Project, repository and Epic filters are exact, independently optional IDs combined
with AND. Unknown IDs are errors. Known filters with no common work return an empty
result. The workspace identity is its physical discovered root; paths and titles
are not Task identities. All matching registered projects and Epics, including empty
ones, remain in the JSON response. No rows are silently truncated.

## Presentation and history

Each Task has one category, in this order:

| JSON category | English heading | Meaning |
| --- | --- | --- |
| `needs_you` | Needs you | Active Task and Epic, recorded progress, and a current explicit human question or qualified product QA need. |
| `follow_up` | Follow-up | Agent work, a recorded blocker, conflicting evidence, or a relevant unknown that needs inspection. The actor may be unknown. |
| `in_progress` | In progress (recorded) | Recorded work without a higher priority action. This does not prove a live agent exists. |
| `ready_next` | Ready next (registered) | A queued goal without a known registered blocker. Actual start still requires execute preflight. |
| `inactive` | Inactive | Task or parent Epic is parked, frozen or archived. |
| `completed` | Completed (registered) | Native evidence confirms the result was integrated, with the required result binding. |
| `backlog` | Backlog | Untouched, unqueued work; active lifecycle or a Spec alone does not mean execution has started. |

The first four categories appear by default. Hidden counts are mutually exclusive:
inactive lifecycle takes precedence, then completed work, then untouched backlog.
`--all` includes those rows with their lifecycle, native progress and historical
actions, and reports zero hidden counts. An inactive or completed Task's actions
have `current: false`; an old question does not become a present human obligation.
Source problems remain visible in diagnostics even when a related Task is hidden.

All actions are preserved. A Task with multiple needs appears in the highest
applicable category. Native result/QA conflicts remain conflicts; timestamps never
select the most favorable outcome. A report saying work is done, or a terminated
transport, cannot establish technical pass, human pass or integration.

Queue rank is local to its Project/Repo/Epic queue. A current preparation,
unresolved effects, an invalid goal revision or changed Problem can block another
goal. Older solution Specs that require live validation have unknown readiness.
The response preserves the exact queue target and original `goal` Spec ID,
revision and manifest digest when known. A current queue entry's `spec` may be the
execution Spec, while `goal` retains the original chosen outcome. The nullable
`selection` preserves the selected-solution decision for legacy queues.
Corrupt or orphaned records without a safe Task binding produce
diagnostics instead of invented Tasks.

## JSON, flags and discovery

For this command, `--json` and `--format json` select the same result object.
`--json` works before `workflow`, between `workflow` and `status`, or after `status`.
`--json=false` selects no format. A true `--json` with explicit `--format text` is
an error in either order; using both JSON selectors is valid. Other commands keep
their existing global `--json` logging behavior.

A successful JSON read writes exactly one object and a newline to stdout.
Diagnostics and logging, including with `--debug`, use stderr when emitted outside
the result. Readable partial results may exit 0 with explicit `unknown` or `stale`
facts and structured diagnostics. Invalid arguments, missing workspaces and fatal
registry errors exit nonzero with an explanation on stderr and empty stdout. Reads
never initialize a workspace or a user profile.

```sh
ply workflow status --json |
  jq '.items[] | select(.category == "needs_you") | {task_id, next_actions}'
```

Discover the operation in `ply capabilities --format json`, under
`read_extensions` with `id: "workflow.status"` and
`command: ["workflow", "status"]`. The original `operations` array and
`coverage: "workspace-core-read"` remain unchanged. An advertised JSON operation
that fails should surface its error; clients should not fall back to parsing text.

## Version 1 contract

The self-contained [JSON schema](../schemas/read-contract/workflow-status-v1.schema.json)
defines `WorkflowStatusReadback@1`, `schema_version: 1`.
[Examples](../schemas/read-contract/workflow-status-examples/) cover empty,
populated, historical, unknown and startup recovery results, plus invalid actor
and live-state claims.

| Field | Contract |
| --- | --- |
| `workspace`, `filters`, `include_all` | Physical root and exact requested scope. Absent filters are explicit `null`. |
| `observed_at_utc` | UTC time the reader observed sources, never an agent activity timestamp. |
| `source_basis` | Always `registered`; these are durable facts, not live probes. |
| `freshness` | `fresh`, `stale` or `unknown` source observation. |
| `projects`, `epics`, `items`, `diagnostics` | Required arrays; an empty collection is `[]`. |
| `counts` | `total = visible + hidden.inactive + hidden.completed + hidden.backlog`. Category counts count visible Task rows exactly once. |
| Item identity and title | Required Project, Repo, Epic and Task IDs. `title: null` with `title_status: "unavailable"` explicitly represents an unreadable current Problem title. |
| `progress` | The existing native progress facts: technical gate, actual human QA, integration, evidence IDs and registered process facts remain separate. |
| `next_actions` | Explicit `actor` (`human`, `agent`, `ply`, `unknown`), source, stable kind, reason, severity, real timestamp or `null`, evidence IDs, nullable run/queue IDs and `current`. |
| `runs` | All bound attempts with provider, transport, available session/Herdr bindings, source bindings and independent `state_freshness`. Missing values are null. |
| Run `delivery` | Optional for runs with a delivery contract; preserves goal/basis, reports, event IDs, report bindings, reported phase and next action. It is omitted for other runs. |
| Run `startup_recovery` | Optional typed recovery history: generation, record binding, original and replacement executable bindings, and available original/replacement transport bindings. Executable paths and digests are recorded metadata, with no current runtime or permission claim. |
| `queues` | Bound queue entries with local rank or null, exact goal when known, current preparation or null, registered readiness and reasons. |
| `diagnostics` | Stable code, source, explicit actor, severity, reason, time or null, nullable Task/run/queue identity and evidence IDs. |

Rows sort by category priority. Action-bearing categories sort by existing severity
(`attention`, then `next_step`), oldest known action time, and stable identities;
unknown times sort last. Ready rows sort by Project/Repo/Epic, rank within that queue,
and Task identity, without inventing a global priority among unrelated queues.
Projects sort by ID; Epics by Project and Epic IDs. Runs sort by run ID, queue bindings
by queue ID and state. Actions sort by `current: true` first, severity, oldest known
time, then source, kind, actor, run ID, queue ID and reason. Diagnostics sort by
source, Task ID, run ID, queue ID, code and reason. Nullable identity fields sort as
empty strings. The array descriptions in the versioned schema specify these keys.

Objects allow additive fields. Clients must tolerate unknown object fields and
owner action/diagnostic codes. Source, actor, lifecycle, category, freshness,
readiness and verdict enum values are documented in the schema. Delivery phase,
meaning, recorded actor and current-preparation state preserve their owner
vocabulary rather than claiming a new native result. Changes to field meaning,
required nullability, closed enum semantics or ordering require a new contract
version.

## Observation limits and verification

The reader captures and reuses registered sources in one request and projects
Tasks in batches. It does not call detailed `task show` or `workflow execute show`
for each row. It performs no Git status or fetch, provider/Herdr calls, installation,
locks, journal writes, repair or restart. An unreadable optional source preserves
valid siblings and adds diagnostics. A detected source change produces bounded
unknown/stale knowledge; the stores do not offer a transactional snapshot.

`last_activity_utc` and `since_utc` use real recorded timestamps. Current delivery
reports and events do not contain event times, so their timestamps stay null and
text shows “time unknown.” File modification time and observation time never fill
that gap. A freshly read old file is not fresh run state; cached working never
proves liveness. The first contract promises no whole-response cache tag, including
for external Git state, process state or referenced artifacts.

Known startup recovery metadata is read through a separate historical adapter.
Its record, before-state, context and replacement transport remain bound to their
captured source digests and original run. The adapter does not inspect executable
bytes, contact a provider, attest permissions or authorize another recovery.
Original transport state and observation fields retain their recorded meaning;
they do not establish present liveness. Recovery does not supply a missing report
timestamp. The mutation and callback schemas keep their existing validation rules.

Run technical acceptance with the pinned schema-validator environment:

```sh
PLY_TEST_PYTHON=/path/to/venv/bin/python sh test/workflow_status_acceptance.sh
/path/to/venv/bin/python test/workflow_status_schema_test.py --verbose
/path/to/venv/bin/python test/workflow_status_schema_test.py --responses /path/to/captured/json
```

The schema tests use the existing pinned `test/read_contract_requirements.txt`
dependency and can validate captured actual CLI output. Acceptance runs three
bounded Go groups covering CLI/projection, registered queue/content, and preserved
run/report behavior, with a five-minute timeout per Go package rather than a wall
time limit for an entire group. It preserves the complete discovered test
inventory, actual group output and exits, `incomplete` for started tests interrupted
before a result, and `not_run` for tests outside the executed groups in
`go-test-coverage.json`. This is scoped verification, not a claim that
the full repository test suite passed. The local CLI journey and schema checks
also preserve their evidence; automated acceptance does not supply human QA.
Performance measurements describe the tested workspace and binary; no universal
latency bound is promised.

A representative fixture with 4 projects, 12 Epics and 300 planned Tasks measured
89,129,764 ns/op (about 89 ms), 61,793,013 B/op and 1,520,602 allocations/op on an
Apple M4 Max. This is a three-iteration observation of the registered reader, not
a general latency guarantee or a measured improvement over earlier commands:

```sh
go test ./internal/workspaceview -run '^$' -bench BenchmarkWorkflowStatusMultiProject -benchtime=3x -count=1
```
