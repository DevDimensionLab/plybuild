# Workspace read contract

Use the registered Project and Task lists to inspect workspace registrations. These
lists include Tasks without queue positions, available worktrees or readable problem
content. The ready queue remains a separate view.

```sh
ply capabilities
ply capabilities --format json
ply workspace project list
ply workspace project list --format json
ply workspace project show --format json -- alpha
ply workspace task list --project alpha
ply workspace task list --repo core --epic alpha-epic --format text
ply workspace task list --project alpha --repo core --epic alpha-epic --format json
```

Each command defaults to `--format text`; only `text` and `json` are accepted. The
global `--json` flag controls logging and does not select the result format. Existing
ID syntax and `--` handling are unchanged. Project list, Task list and capabilities
accept no positional arguments; Project show takes exactly one Project ID.

## Registered Task filters

Ordinary `workspace task list` accepts `--project`, `--repo` and `--epic` independently
or in any combination. Every supplied ID is checked, and the resulting list contains
Tasks matching **all** supplied IDs exactly. The current directory selects the
containing workspace; it does not implicitly select a Project, repository or Epic.

Known filters with no common Tasks succeed with an empty list, including a Project
and repository that are not registered together. Unknown IDs fail with
`workspace_work_not_found`. Invalid IDs, explicitly empty filters, invalid formats
and invalid arguments fail with `workspace_work_invalid_arguments`. Format validation
happens before workspace data is read.

Text output shows Task IDs, titles and their registered bindings and worktree states.
When Project or repository filters are supplied, the heading identifies every supplied
filter. All previously valid text calls without those new filters preserve their
output, including Epic-only lists, empty lists and unavailable titles. Explicit
`--format text` is also accepted for ordinary lists.

`workspace task list --ready` keeps its existing target rules: supply all three target
flags, or omit all three from the appropriate Epic working directory. Its dataset,
ordering, errors, formats and `WorkspaceTaskQueueReadback@1` response are unchanged.
The independent filter rules above apply only to the ordinary registered-Task list.

For example, `ply workspace task list --project alpha --repo core --epic alpha-epic`
can display:

```text
Tasks for Project alpha, repository core, Epic alpha-epic in Ply workspace /work:
  beta-legacy: Legacy registration (Epic alpha-epic; alpha / core; unbound)
```

With a known repository `other-repo` and no matching Tasks,
`ply workspace task list --project alpha --repo other-repo` displays:

```text
No Tasks for Project alpha, repository other-repo are registered in Ply workspace /work.
```

## JSON responses

Successful JSON reads emit one UTF-8 JSON object followed by a newline. No logs or
progress text are written to stdout. Each new response has a `kind` and integer
`schema_version: 1`. All documented fields are required, including fields whose value
is `null`. Consumers must tolerate unknown additive fields. Incompatible shapes or
meanings require a new schema version; storage versions are independent.

Workspace responses contain `workspace.root`, the physical absolute root from the
validated workspace marker. This path is the workspace identity; no additional UUID
is introduced. Collections are always arrays and use `[]` when empty. There is no
pagination, hidden truncation, timestamp, canonical serialization or hash guarantee,
and no atomic snapshot guarantee across registers or commands.

| Command | `kind` | Payload |
| --- | --- | --- |
| `workspace project list` | `WorkspaceProjectListReadback@1` | `scope.project_id: null`, `projects: [Project]` |
| `workspace project show` | `WorkspaceProjectReadback@1` | `scope.project_id: string`, `project: Project`, `repositories: [Repository]` |
| `workspace task list` | `WorkspaceTaskListReadback@1` | `scope` with nullable `project_id`, `repo_id`, `epic_id`; `tasks: [Task]` |
| `capabilities` | `PlyCapabilities@1` | `build`, `coverage`, `operations`; no workspace is required |

Projects are sorted by `project_id`, repositories by `repo_id`, and Tasks by `task_id`,
all ascending. A Project contains string `project_id`, `name` and `wrapper`, plus the
nonnegative integer `repo_count`. The count reflects registered repository IDs. A
Repository contains string `repo_id`, `locator` and `git_common_dir`. Wrapper, locator
and Git common directory are registered absolute paths. Reading them performs no
existence check or new Git observation; missing checkouts remain in the results.

For example, a registered Project with two repositories produces:

```json
{
  "kind": "WorkspaceProjectListReadback@1",
  "schema_version": 1,
  "workspace": {"root": "/work"},
  "scope": {"project_id": null},
  "projects": [
    {
      "project_id": "alpha",
      "name": "Alpha – café / core",
      "wrapper": "/work/alpha wrapper",
      "repo_count": 2
    }
  ]
}
```

Its default text list is:

```text
Projects in Ply workspace /work:
  alpha: Alpha – café / core (2 repositories) /work/alpha wrapper
```

An empty filtered Task result still states the requested scope:

```json
{
  "kind": "WorkspaceTaskListReadback@1",
  "schema_version": 1,
  "workspace": {"root": "/work"},
  "scope": {"project_id": "alpha", "repo_id": null, "epic_id": null},
  "tasks": []
}
```

A Task has string `task_id`, `project_id`, `repo_id`, `parent_epic_id` and
`worktree_state`. The worktree state is the raw registration value, without an
additional aggregate status or a claim of fresh disk observation. Title availability
is explicit:

| Condition | `title` | `title_source` | `title_status` |
| --- | --- | --- | --- |
| Current problem revision is readable and verified | Current title string | `problem_revision` | `available` |
| Valid older registration has no problem head | Registered title string | `registration` | `available` |
| Problem head exists but its content or title cannot be read or verified | `null` | `problem_revision` | `unavailable` |

An unavailable current problem does not fall back to the registration title. A real
title equal to `Problem content unavailable` remains available. Consumers should read
the explicit fields instead of interpreting title text. The Task remains in the list
when its title is unavailable; this is successful readback with no separate stderr
warning.

Argument, scope and registry failures return a nonzero exit code, empty stdout and an
English error on stderr. Project operations retain their existing error classes;
invalid Project output formats use `workspace_project_invalid_arguments`. An
initialized workspace with no registrations succeeds with empty collections. A
corrupt or unreadable registry fails instead of appearing empty. Output I/O failure
can leave partial stdout and a nonzero exit; writes to a consumer's pipe are not
guaranteed atomic.

## Capabilities

`ply capabilities` works outside a workspace, without profiles, credentials, network
access, initialization or new files. It reports an explicit
`coverage: "workspace-core-read"` catalog, complete for the four operation IDs below.
Other commands, including Spec and Journal operations, are outside this catalog.

The JSON `build` object contains string `version`, nullable string `vcs_revision` and
nullable boolean `vcs_modified`. Metadata comes from the binary, never from a Git
command or the caller's current directory. Missing metadata is `null`; text output
shows `unknown` for missing values.

Each operation includes its `id`, `command` words after `ply`, `mode`, `selectors`,
`formats`, `filters`, `filter_policy`, `result_schemas`, and `effect: "read"`.
Operations are sorted by `(id, mode)`, so the `ready` entry comes before `registered`.
Every listed operation supports `formats: ["text", "json"]`. Selectors contain required
mode flags; filter names omit the `--` prefix.

| ID / command | Mode | Selectors | Filters / policy | Result schemas (`kind`, `schema_version`) |
| --- | --- | --- | --- | --- |
| `project.list` / `workspace project list` | `default` | `[]` | `[]` / `none` | `WorkspaceProjectListReadback@1`, `1` |
| `project.show` / `workspace project show` | `default` | `[]` | `[]` / `none` | `WorkspaceProjectReadback@1`, `1` |
| `task.list` / `workspace task list` | `ready` | `["--ready"]` | `["project","repo","epic"]` / `all_or_none` | `WorkspaceTaskQueueReadback@1`, `1` |
| `task.list` / `workspace task list` | `registered` | `[]` | `["project","repo","epic"]` / `independent_and` | `WorkspaceTaskListReadback@1`, `1` |
| `task.show` / `workspace task show` | `default` | `[]` | `[]` / `none` | `WorkspaceTaskReadback@1`, `null`; `WorkspaceTaskIntegrationReadback@1`, `1`; `WorkspaceTaskReadback@3`, `3` |

The `null` catalog version means that the legacy Task-show response **omits** its
`schema_version` field. These entries describe top-level response forms; the nested
integration version 2 object in Task-show version 3 has no separate entry. The catalog
does not add flags for choosing these versions.

Consumers should continue to construct argv from their own known command allowlist.
An operation being listed does not start an agent, grant write authority, guarantee
fresh local state or assert provider/runtime support. `read` means no persistent
workflow, Git or workspace mutation; older read modes may take advisory locks on
existing files. No general JSON error-body contract is advertised: ready mode already
has cases that return a JSON body with a nonzero exit.

Invalid capabilities arguments and formats fail with `capabilities_invalid_arguments`,
nonzero exit, empty stdout and English stderr. Text capabilities output reports the
build, catalog coverage, operations, modes and formats with stable English labels.

## Schemas and examples

The four self-contained JSON Schemas use Draft 2020-12 and permit additive fields at
every object level:

- [Project list schema](../schemas/read-contract/workspace-project-list-readback-v1.schema.json)
- [Project show schema](../schemas/read-contract/workspace-project-readback-v1.schema.json)
- [Task list schema](../schemas/read-contract/workspace-task-list-readback-v1.schema.json)
- [Capabilities schema](../schemas/read-contract/ply-capabilities-v1.schema.json)

[Examples](../schemas/read-contract/examples/) include populated and empty collections,
all three title states, a literal sentinel-like title, known and unknown build metadata,
and invalid documents for each response. Filenames end with `.valid.json` or
`.invalid.json`. Schemas validate document shape and the fixed capabilities catalog;
product tests also check ordering, registry facts, exact filtering and absence of
read-side mutation. Example paths describe registrations and need not exist on disk.

The schema test uses the test-only `jsonschema==4.26.0` distribution from
[PyPI](https://pypi.org/project/jsonschema/4.26.0/), pinned in
[test/read_contract_requirements.txt](../test/read_contract_requirements.txt). Ply's
runtime has no Python dependency. From the repository root, prepare a local test
environment and run:

```sh
python3 -m venv /private/tmp/ply-read-contract-190/schema-venv
/private/tmp/ply-read-contract-190/schema-venv/bin/python -m pip install -r test/read_contract_requirements.txt
/private/tmp/ply-read-contract-190/schema-venv/bin/python test/read_contract_schema_test.py --verbose
```

To validate actual serializer output, save successful JSON stdout for each of the four
new kinds as separate `.json` files under one directory, then pass that directory with
`--responses`. The test validates every `.json` file recursively, requires all four
kinds, checks the final newline and rejects duplicate object keys. It never runs Ply
or changes fixture state:

```sh
/private/tmp/ply-read-contract-190/schema-venv/bin/python test/read_contract_schema_test.py --responses /private/tmp/ply-read-contract-190/read-responses --verbose
```

## Compatibility for consumers

Existing argv and default text continue to work. Existing ready responses, Task-show
variants and other Task, Spec and Journal readbacks keep their contracts. A later
application migration can choose known advertised JSON forms while retaining its
existing help/text path for older binaries. Once JSON support is advertised, report
JSON read failures directly instead of silently retrying a text parser.
