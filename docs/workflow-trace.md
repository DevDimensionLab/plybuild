# Workflow trace for one Task

`ply workflow trace <task-id>` reads the recorded process for one registered Task
across its runs. Use it to investigate repeated checks, questions, waiting,
candidate changes, human QA and explicit handoffs. Untouched, inactive and
completed Tasks are valid selections. Run from any directory inside the physical
Ply workspace; an Epic worktree is not required.

```sh
ply workflow trace task-id
ply workflow trace task-id --format json
ply --json workflow trace task-id
```

The text starts with evidence-linked analysis and a source-position map using
Human, Agent, Ply and Unknown roles. Full declarations, event payloads, candidate
facts, relations, sources, coverage and diagnostics follow. No history is silently
truncated. Strings use JSON escapes to preserve exact whitespace and safely show
control characters and Unicode in an ASCII terminal. For example, `\n` represents
an original newline in a question; `null` is missing information, not zero.
The leading summary explains that coverage concerns selected records, and that
reported QA/wait and verifier elapsed time are not active effort or time savings.

## What is recorded

The projection keeps these layers separate:

| Layer | Where to inspect it | Meaning |
| --- | --- | --- |
| Current declarations | `current_goals` | Latest readable registered goal for each Spec. `status: unknown` means a newer revision could not be safely classified. |
| Frozen contract | `runs[].frozen_goal`, `frozen_basis`, `declared_process` | The exact goal, execution basis and obligations that the original run used. A later goal edit does not rewrite a frozen run. |
| Recorded history | `events`, `chains`, `relations` | Preserved native records and reported process contributions. A declaration is not evidence that its steps occurred. |
| Candidate facts | `candidates[].axes` | Independent reported, technical, human QA and integration facts for an exact repository/commit/tree pair. |
| Analysis | `analysis` | Recorded counts, compatible elapsed intervals and inspectable questions, with definitions, evidence and coverage. |
| Unknowns | `coverage`, `diagnostics`, null fields | Missing, invalid, changing or ambiguous evidence and limits of the recorded history. |

Supported run families are `legacy_task_run` (`trn_`), `workflow_run_v1` (older
`wfr_` with preserved review rounds), and `goal_execute_v2` (goal execution).
Each retains its own identity and semantics. Missing declared process remains
unknown. An explanatory guide is not a registered workflow definition or an
executed activity. `history_state: empty` means there are no recorded runs;
Task declarations or other registered events may still be present. Text says
"No execution recorded" and explicitly labels absent process declarations
"Declared process unknown".

An event's `role` is `human`, `agent`, `ply` or `unknown`. `actor_claim` records
what the source actually claims, while `recorder` identifies the recording
surface. Provider names and recorder names do not authenticate a human.
`evidence_class` distinguishes `reported`, `human_attestation`, `controlled`,
`registered` and `unknown` evidence. Even controlled verification is separate
from a candidate-bound human judgment.

A report saying "done", a passed technical gate, human pass/fail, and local
integration are independent facts. Candidate axes retain all outcomes and flag
conflicts; they do not choose the latest favorable answer. A human pass for C1
does not qualify C2. A changed candidate does not itself prove a correction or
explain its cause.

## Reading the process

Begin with the analysis questions and their `evidence_ids`. Find those IDs in the
history, inspect each event's `data` and `source_ids`, then follow the source
locator and digest. Candidate and result references keep a failure or pass tied
to the actual candidate. Exact questions and explicit waiting fields (`reason`,
`dependency`, `next_actor`) remain in event data.

A `needs_input` report followed by a `working` report preserves both statements.
It does not invent a human answer or response interval. Explicit `responds_to`,
`corrects`, `verifies`, `supersedes` and `handoff` relations retain their source
meaning. Adjacent roles and nearby timestamps do not create a handoff or causal
link. No waiting report is not evidence of no waiting.

For example, save a readback and follow an analysis question without parsing text:

```sh
ply workflow trace task-id --format json > trace.json
jq '.analysis[] | select(.question != null) |
    {id, question, definition, coverage, evidence_ids, unknowns}' trace.json
```

Use an actual analysis ID from that output:

```sh
jq --arg analysis_id 'repeat_verification' '
  . as $trace |
  .analysis[] | select(.id == $analysis_id) as $analysis |
  {analysis: $analysis,
   events: [$trace.events[] |
     select(.id as $id | $analysis.evidence_ids | index($id))],
   sources: [$trace.sources[] |
     select(.id as $id | $analysis.source_ids | index($id))]}' trace.json
```

`repeat_verification` exists only when the recorded evidence supports it. If the
same candidate and verification inputs are not exactly bound, multiple attempts
stay distinct without being called the same verification basis. Even an exact
repeated basis does not prove redundancy or equal external conditions. An
analysis question identifies evidence to inspect; its answer and any improvement
remain unmeasured.

## Order, identity and time

`events`, `runs`, `sources`, `chains`, `candidates` and `analysis` use deterministic
ID sorting. Current goals sort by Spec ID. Relations sort by source event, type
and destination; diagnostics sort by code and detail. These are presentation
orders, not a global chronology. Text's source map follows each chain's own
`event_ids`; full event details use the deterministic event-ID order.

`events[].positions` gives the authoritative source chain and original sequence.
Each `chains[].event_ids` list follows that chain's sequence, with tied positions
sorted by event ID solely for presentation. Equal positions assert no relative
order. An event may have several positions when exact identity is reached through
several source paths. Independent chains with missing or overlapping clocks do
not establish a global event order. Diagram spacing does not represent duration.

The same exact native QA, result axis or integration record reached through
multiple adapters is represented once with its provenance retained. A result's
reported and technical axes remain separate facts. Equal text or different,
ambiguously related IDs are not silently treated as one action, or as proof of
different real actions.

`occurred_at_utc`, `reported_at_utc`, `registered_at_utc` and top-level
`observed_at_utc` have different meanings. `time_basis` retains clock, precision
and uncertainty. Untimed records keep null times; filesystem modification time
never supplies a missing event time. Observation freshness does not prove agent
liveness. `live_state` is always `unknown`.

Only compatible endpoints support an interval. Native verifier start/end can
yield verifier elapsed time; waiting and QA intervals remain reported intervals.
An explicitly linked request and answer can yield their recorded elapsed time.
Missing, reversed, uncertain or incomparable endpoints yield an unknown value
and explanatory evidence or diagnostics. None of these intervals measures active
work, human attention, cost, waiting cause, efficiency or time savings.

## JSON contract and failures

The self-contained [schema](../schemas/read-contract/workflow-trace-v1.schema.json)
defines `WorkflowTraceReadback@1`, `schema_version: 1`.
[Examples](../schemas/read-contract/workflow-trace-examples/) cover empty history,
all three run families with candidate-bound QA and questions, partial history,
and invalid actor/authority claims. These examples are synthetic technical
fixtures, not actual human QA.

All projection fields are required and collections are arrays, including `[]`.
Missing scalar evidence uses explicit `null`. The projection envelope, roles,
versions, clocks, source references and analysis fields are strict. Safe native
payloads in `data`, `declaration`, `frozen_goal`, `declared_process` and
`frozen_basis` retain their native contract shape and can carry native additions.
Consumers should use the projection's typed fields and the payload's own version
when interpreting its details.

| Reference or status | Contract |
| --- | --- |
| `source_ids` | IDs in `sources`; each source has its locator, nullable SHA-256 digest normalized as `sha256:<64 lowercase hex digits>`, and observation status. |
| `evidence_ids`, `event_ids` | IDs in `events`. |
| `run_ids`, `candidate_id` | IDs in `runs` and `candidates`. |
| `result_id`, `result_ids`, `native_id` | Original native identities retained for inspection; they are not another count of actions. |
| `positions[].chain_id` | An ID in `chains`. Sequence belongs to that chain only. |
| `analysis[].coverage` | `exact` for its defined recorded evidence, `lower_bound` when missing records may add to a count, or `unknown`. Every value names its definition and evidence. |
| `coverage.state` | `complete` or `partial` coverage of selected registered sources, never completeness of real work. |
| `coverage.read` | Source families consulted, not event/source references or proof that every record existed. |
| `freshness` | `fresh` or `unknown` source observation, independent of live activity. |
| `history_state` | `empty`, `recorded` or `partial`; empty/partial success is different from a fatal read failure. |
| `next_transition_authorized` | Always `false`. A trace grants no operational authority. |

`run_ids` references the safely bound outer runs in this result. Original internal
handoff aliases remain in native records and `events[].data.native_run_ids` when
present; an alias is not an additional delivery run.

For trace, `--json` and `--format json` select the same result, including when
`--json` precedes `workflow`. `--json=false` selects no format. A true `--json`
with explicit `--format text` is an error. A successful JSON read emits exactly
one object and a newline to stdout. Logging stays on stderr, including with
`--debug`; diagnostics are also retained as structured facts in the result.
Other commands keep their existing global logging behavior.

Unknown Tasks, missing workspaces and fatal authoritative registry failures exit
nonzero with empty stdout and an explanation on stderr. Corrupt, missing,
foreign, orphaned, hash-mismatched or changing optional history produces
identifiable diagnostics and partial/unknown coverage while independent valid
history remains visible. Re-reading an old file does not prove live activity.

Discover the additive operation in `ply capabilities --format json` under
`read_extensions`, with `id: "workflow.trace"` and
`command: ["workflow", "trace"]`. The existing core catalog, status, journal,
execute and write contracts remain unchanged.

## Observation boundary and validation

Trace reads registered local stores without Git, provider, Herdr, Slack,
network or process probes. It does not initialize a workspace/profile or create
directories, locks, caches, receipts or journal entries. It does not refresh,
repair, run acceptance/verifier commands, start/retry a run, record QA, integrate
or change a workflow. Concurrent filesystem changes are diagnosed where observed;
the read is not a transactional filesystem snapshot.

The schema tests also check reference integrity, source sequence and candidate
isolation. Validate preserved actual CLI JSON responses with the test-only
validator from `test/read_contract_requirements.txt`:

```sh
python3 test/workflow_trace_schema_test.py --responses /path/to/captured-responses
```

The response directory must contain at least one actual trace response. Runtime
Ply has no Python dependency. Acceptance records read cost and source coverage for
its tested workspace; those measurements are not a claim of process improvement.
