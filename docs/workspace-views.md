# Workspace views for applications

Ply exposes versioned JSON projections for Home screens, navigation and activity.
The application chooses its presentation; Ply owns registration, progress and
evidence semantics. These views do not start agents, fetch repositories or repair
state. Existing Task list, ready queue and detailed Task show remain available.

## Start with one Home request

```sh
ply capabilities --format json
ply workspace overview --format json
ply workspace overview --project ply --format json
```

Discover the new operations in `PlyCapabilities@1.read_extensions`. The original
`operations` array and `coverage: "workspace-core-read"` remain unchanged. Combine
both arrays when discovering operations, but construct argv from commands your
application understands. An older binary without `read_extensions` does not
advertise these views. Do not silently fall back to parsing text after an advertised
JSON command fails.

`WorkspaceOverviewReadback@1` contains `tasks`, `epics` and `attention` derived from
the same captured registrations and process batch. It avoids one detailed Task
request per row and does not run live Git or Herdr checks. There is no cross-store
transaction or guarantee that separate command invocations see identical data.

For navigation, combine this with `workspace project list` and
`workspace project show <project-id>`. Project JSON now includes `companions`;
repository JSON includes nullable `wrapper`. Empty arrays and `null` are real values,
not fields to infer from a filesystem scan.

## Read operations

All commands below support `--format text|json` and default to text. Successful JSON
output is one object followed by a newline, with `kind` and `schema_version: 1`.
Use `--format json`; the global `--json` logging option is not a result selector.

| Command after `ply workspace` | JSON kind | Filters |
| --- | --- | --- |
| `overview` | `WorkspaceOverviewReadback@1` | `--project` |
| `task list --progress` | `WorkspaceTaskProgressListReadback@1` | `--project`, `--repo`, `--epic` |
| `task show <task-id> --progress` | `WorkspaceTaskProgressReadback@1` | none |
| `epic list` | `WorkspaceEpicListReadback@1` | `--project` |
| `attention` | `WorkspaceAttentionReadback@1` | `--project`, `--repo`, `--epic` |
| `journal recent` | `WorkspaceActivityReadback@1` | `--project`, `--since`, `--limit` |
| `run list` | `WorkspaceRunListReadback@1` | `--project`, `--active` |
| `worktree list` | `WorkspaceWorktreeListReadback@1` | `--project`, `--repo` |
| `status` | `WorkspaceStatusReadback@1` | none |
| `task lifecycle show <task-id>` | `WorkspaceWorkItemLifecycleReadback@1` | none |
| `epic lifecycle show <epic-id>` | `WorkspaceWorkItemLifecycleReadback@1` | none |
| `project metadata show <project-id>` | `WorkspaceProjectMetadataReadback@1` | none |

ID filters are exact, independently optional, and combined with AND. Unknown IDs
fail; known filters with no common entries return empty arrays. Results never
silently truncate. `journal recent --limit 0` means all events; a positive limit is
explicit. `--progress` and `--ready` are mutually exclusive. An ordinary Task list
keeps its existing shape and text; `--progress` opts into the richer contract.

## Progress and human attention

Each progress Task row carries the same `progress` block in list, show and overview.
Tasks sort by `task_id`; Epics sort by `epic_id`. Epics with no Tasks remain visible,
with zero counts. Epic bases are the current **registered** revisions, not newly
observed branch heads.

The progress block separates these facts:

- `has_progress`: execution, result, integration or process history exists. Merely
  registering a Task or writing its plan does not count.
- `state`: `nothing`, `in_progress`, `attention` or `integrated`. This is a registered
  summary, not live integration readiness or an agent liveness probe.
- `technical_gate` and `human_qa_outcome`: nullable native facts. Agent reports,
  transport completion and journal claims cannot manufacture a human QA result.
- `task_result_id`, `integration_result_id`, `result_count`: the evidence behind
  selected facts. Conflicting recorded results do not become a successful result
  by choosing a timestamp or discarding failed results.
- `next_actions`: explicit actions with `source`, `actor`, `kind`, `reason`,
  `severity`, `since_utc` and `evidence_ids`. Use these fields instead of interpreting
  a free-text status.
- `last_activity_utc`: the latest available recorded event time, or `null`. A read
  does not invent an activity timestamp.

`result_count` includes historical results. If the registered Problem, selected Spec
or assessment changes, those older results remain history rather than current QA
evidence: the selected result/QA fields are null, the reason is
`recorded_result_basis_stale`, and the agent action is `inspect_task_basis`. Merely
publishing an unselected Spec revision does not change the selected basis.

`attention` includes only actions whose actor is `human`, for Tasks with progress
and with both Task and parent Epic lifecycle `active`. Untouched Tasks belong in a
client's Ready next view, never Needs you. A reported wait remains distinguishable
from native integration or QA facts through its source and evidence IDs. Entries
sort by severity, then oldest known time, then stable subject/source/kind keys;
missing times sort last within the severity group.

For operational decisions, use the existing detailed Task show, readiness and
workflow operations. A Home projection does not grant authority to integrate.

## Lifecycle without changing execution behavior

Task and Epic lifecycle is explicit: `active`, `parked`, `frozen`, `archived`.
Older registrations default to `active` without creating or migrating any file.
Lifecycle organizes visibility and attention. It does not stop an agent, remove a
worktree, clear a queue or change existing workflow permissions. Ply never guesses
which old Tasks should be archived.

```sh
ply workspace task lifecycle show my-task --format json
ply workspace task lifecycle set my-task parked --actor human --format json
ply workspace epic lifecycle set my-epic archived --actor human --expected-revision 7 --format json
```

Set requires an actor claim and records a durable change event. Repeating the
current state is a no-op. `--expected-revision` checks the **shared workspace
lifecycle revision**, including changes to other Tasks and Epics, before writing.
After a conflict, read again and make a new intentional choice. Actor text records
who is claimed to have made the change; it is not authentication or human QA.
Mutation results use `WorkspaceWorkItemLifecycleMutation@1`.

Storage is an atomic `.ply/work-item-lifecycle.json` sidecar. Existing work-item
formats and integration authority hashes retain their meaning. Ordinary Task list
and Task/Epic show JSON also receive the additive `lifecycle` field.

## Places metadata

Register navigation relationships explicitly; these commands neither create a
directory nor register another repository:

```sh
ply workspace project companion add ply --role planning --path ./ply/planning --actor human --format json
ply workspace project companion add ply --role docs --path ./ply/main/docs --actor human --format json
ply workspace project repo-wrapper set ply --repo ply --wrapper ./ply --actor human --format json
ply workspace project metadata show ply --format json
ply workspace project companion remove ply --role docs --path ./ply/main/docs --actor human --format json
ply workspace project repo-wrapper clear ply --repo ply --actor human --format json
```

Companion roles are `planning`, `docs` and `other`. Add/set resolves physical paths,
validates directories and, for a repository wrapper, checks containment of its
registered repository locator. Subsequent reads retain missing paths as registered
metadata. Remove/clear can remove stale registrations. There is no automatic
directory classification or repair.

Metadata mutations require `--actor`, accept `--expected-revision`, preserve history,
and return `WorkspaceProjectMetadataMutation@1`. Their revision is shared by all
project metadata in this workspace, separate from lifecycle revisions. Repeating
an existing relationship is a no-op. The atomic sidecar is
`.ply/project-metadata.json`.

## Activity, runs and worktrees

```sh
ply workspace journal recent --since 24h --limit 100 --format json
ply workspace journal recent --since 2026-10-06T04:00:00Z --format json
ply workspace run list --active --format json
ply workspace worktree list --project ply --format json
```

Activity reuses Task journal event identities and includes lifecycle changes. The
`--since` bound is inclusive and accepts a UTC RFC3339 timestamp or a positive Go
duration such as `30m` or `24h`. Undated events are shown last without a time filter;
with a filter they are excluded and counted separately. Partial time coverage does
not mean a readable source is corrupt. Historic Herdr phases without event times
are available in the run view, not invented as timestamped journal events.

Run rows preserve provider (`codex` or `claude`), transport (`herdr`, `terminal` or
`headless`), available session/tab/panel bindings and recorded return facts. Missing
identities are `null`; corrupt records remain diagnosable without hiding valid
siblings. A cached `working` phase cannot prove that a process is alive. The public
state is `completed`, `returned`, `stopped` or `unknown`; `state_freshness` describes
that claim independently of the source read. `--active` means **unresolved, including
unknown**, not definitely running. Do not use it to automatically restart a Task.

Worktree inventory is an explicit, potentially more expensive Git observation. It
reports owner (`project_main`, `epic`, `task` or `none`), current local identity,
cleanliness and available comparisons with a registered base. Nullable
`merged_into_base` means ancestry against that recorded base; it does not authorize
deletion. Detached, locked, prunable, unmanaged and inaccessible entries remain
visible. Per-fact freshness distinguishes unknown identity, cleanliness and
comparison. The command disables optional locks and filesystem-monitor hooks and
does not fetch, prune, repair or clean up anything.

## Freshness and cache invalidation

`fresh`, `stale` and `unknown` describe observation quality, not business success.
`source_basis: "registered"` means recorded facts, even when they were captured
successfully now. Top-level progress/overview freshness describes the shared
register capture; each Task's `progress.freshness` additionally reflects its process
sources. Inspect both. The schema defines required nullable fields and closed
semantic enums; clients must tolerate additional object fields.

`workspace status` hashes the owners of registered facts: Projects, work items,
process journals and run stores. A stable `registry_sha256` lets an application skip
reloading unchanged **registered** data. It includes bytes and relevant directory
structure, not file modification timestamps. Two observations detect changing
sources; stale or unreadable data has a null digest and must not be cached as fresh.
Symbolic links and oversized source files also produce explicit unknown results.

The status digest is not an ETag for every possible display fact. It does not detect
external Git changes, live process changes or modifications to externally referenced
artifacts. Poll worktree observations separately and provide an explicit refresh
for external changes. A stable digest is not source validation, a security decision
or proof that a running agent still exists. Do not claim a transactional snapshot
across these independent files.

## Contracts and verification

The self-contained [workspace views schema](../schemas/read-contract/workspace-views-v1.schema.json)
covers the read and metadata mutation kinds. [Valid and invalid examples](../schemas/read-contract/workspace-views-examples/)
show empty, populated, nullable and unknown forms. Original core read schemas remain
documented in [Workspace read contract](read-contract.md).

`test/workspace_views_schema_test.py` checks schemas and serialized responses using
the existing pinned `test/read_contract_requirements.txt` dependencies.
`test/workspace_views_acceptance.py --ply /absolute/path/to/ply --baseline /absolute/path/to/old-ply --root /fresh/fixture/path`
builds a disposable registered workspace, compares legacy results, exercises
lifecycle and metadata, checks read-only effects and records command timings. It
preserves its fixture and report. The `--baseline` comparison is optional.

Performance measurements describe a specific fixture and binary; this contract
does not promise a universal subsecond bound for every workspace or Git inventory.
