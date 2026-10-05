# ply
A little "go help" for the Java/Kotlin developers using Maven.

Current main capability? 
Upgrade your pom.xml dependencies to latest and greatest! 

Why?
- No installs of maven-plugins required, so if you a working in a multi-repo developer environment with lots of 2party dependencies and repos, you can easily upgrade them with `ply upgrade 2party`. 
- Brings natural semantics and support for different types of dependencies to the table: Kotlin, 2party, spring-boot (curated dependencies), (other) 3party   
- Can be used as a library for other go-projects automating the upgrade process
- Easy and fast
- Brings feature to the table, not found anywhere else, stay tuned!

Heads up!
- ply rewrites your pom.xml, so make sure you have your pom.xml committed before testing out ply
- start with `ply format pom`, verify that the rewrite of the pom.xml is ok, commit, and from now on you will easily see the diff that ply introduces with ```ply upgrade <2party|3party|spring-boot|plugins|all>```
- or just use  `ply status` (no rewrite) and manually upgrade your pom.xml based on what is reported as outdated, current option if you need to keep your pom.xml formatting
  
Requirement: [Go 1.26.7](https://go.dev/doc/install)

```shell script
  _____  _       
 |  __ \| |      
 | |__) | |_   _ 
 |  ___/| | | | |
 | |    | | |_| |
 |_|    |_|\__, |
            __/ |
           |___/ 

Usage:
  ply [command]

Available Commands:
  build        Builds a ply project with ply files and formatting
  capabilities Show supported workspace core read contracts
  help         Help about any command
  plugin       Plugin functionality for plybuild
  profile      Manage profiles settings for ply
  status       Status functionality for a project
  tips         Use tips to learn information faster
  upgrade      Upgrade options
  workflow     Manage agent workflows
  workspace    Manage Ply workspaces

Flags:
      --debug   turn on debug output
      --doc     open documentation website
      --force   uses default for prompts
  -h, --help    help for ply
      --json    turn on json output logging

Additional help topics:
  ply about        About ply

Use "ply [command] --help" for more information about a command.
```

## Workspace
Initialize the current directory as an explicit Ply workspace:

```shell script
ply workspace init
```

The command creates `.ply/workspace.yaml` with format version 1 and the canonical physical
directory path. The directory does not need to be a Git repository. Re-running the command is
safe and leaves an existing compatible marker unchanged. It does not create a Git repository or
create workflows.

From an initialized workspace, or any directory below it, register a project with an explicit
wrapper and one or more Git worktree roots:

```shell script
ply workspace project add ply \
  --name Ply \
  --wrapper ../ply \
  --repo ply=../ply/main
```

Repeat `--repo` to register a multi-repository project:

```shell script
ply workspace project add trip \
  --name Trip \
  --wrapper ../trip \
  --repo trip-frontend=../trip/trip-frontend/main \
  --repo trip-openapi=../trip/trip-openapi/main \
  --repo trip-service=../trip/trip-service/main
```

The wrapper and repository members are explicit and may be outside the workspace. Ply validates
only the nominated worktree roots. Dirty repositories are accepted, and registration performs no
discovery or Git mutation. Read registrations with `ply workspace project show <id>` and
`ply workspace project list`. These commands also accept `--format json`.
See the [workspace read contract](docs/read-contract.md) for independent Task filters,
versioned JSON schemas, and `ply capabilities`.

Retrying the same registration is idempotent. Changing membership, relocating repositories, and
`project init` are not part of this command.

Adopt one existing clean worktree as the repository anchor for a workspace-owned Epic:

```shell script
ply workspace epic adopt ply-agentic-workflow-support \
  --title "Ply agentic workflow support" \
  --project ply \
  --repo ply \
  --worktree /Users/perottochristensen/github/ply/ply_agentic_workflow_support \
  --ref refs/heads/ply_agentic_workflow_support \
  --expected-oid 54f3631cbea789f25a4134945c7ca16d343139df
```

Adoption records the exact physical worktree, full local branch ref, commit, tree, and Git common
directory. It performs no Git change. Each Task belongs to an Epic and binds exactly one registered
repository and Git common directory:

```shell script
ply workspace task create workspace-work-item-bootstrap \
  --title "Workspace-owned Epic, Task, and worktree support" \
  --description "Add explicit workspace work items and prepare a Task worktree from the Epic base." \
  --epic ply-agentic-workflow-support \
  --project ply \
  --repo ply
```

Create the Task worktree from the Epic branch's exact expected commit:

```shell script
ply workspace task worktree create workspace-work-item-bootstrap \
  --branch ply_workspace_work_item_bootstrap \
  --path /Users/perottochristensen/github/ply/ply_workspace_work_item_bootstrap \
  --expected-parent-oid 54f3631cbea789f25a4134945c7ca16d343139df
```

Ply persists a durable create intent before the additive branch/worktree operation. Identical
retries recover safe no-effect, matching branch-only, or exact-effect outcomes with the same IDs.
Partial or unknown effects are preserved for explicit reconciliation and are never reset, removed,
or otherwise cleaned up automatically.

Use `ply workspace epic list`, `ply workspace epic show <id>`, `ply workspace task list`, and
`ply workspace task show <id>` for human readback. Both show commands accept `--format json` for
versioned deterministic machine readback. `worktree_ready` only records a clean local bootstrap
binding: it does not start an agent, select a workflow, or grant execution authority.

After a workflow handoff has an immutable terminal result, record its exact technical evidence
without creating a new Activity or Run identity:

```shell script
ply workspace task result record workspace-work-item-bootstrap \
  --file /absolute/task-result.json
ply workspace task qa record workspace-work-item-bootstrap \
  --file /absolute/human-qa.json
```

Human QA is a separate, locally claimed product judgment. Recording `pass` does not start
integration. First run a read-only check with the selected TaskResult, QA record, result commit,
and current Epic-parent commit. The check returns a canonical confirmation digest; a human starts
a separate `--apply --confirm <digest>` invocation. Apply permits exactly one local `ff-only`
move of the registered Epic parent. It never fetches, pushes, targets `main` or `master`, or cleans
worktrees.

The first accepted TaskResult upgrades a legacy format-1 work-item store atomically to format 2;
all existing facts remain immutable. A no-effect attempt can only be retried through a new check
that names `--retry-after`. Lost responses are reconciled from the durable authority, attempt,
Git state, and reflog; partial or unknown effects stop for read-only recovery control.

Format-2 `task show` presents the same integration state as check/apply in text or canonical JSON.
Treat `git_changed` and the single `next_action` as the authoritative outcome. Integration v1
requires Git 2.45.0 or newer, a `sha1` or `sha256` object format, and the `files` ref backend.

Sub-tasks, additional repository anchors on an Epic, editing, rebinding, cleanup, and
WorkflowRun creation are outside the resource-binding commands. Refresh is explicit only through
`workspace epic base update`; it appends a metadata base and never moves Git refs or old worktrees. Problem and solution revisions
are recorded separately, as described below.

### Prioritize Tasks and prepare a separate worktree

`task list` keeps its registered-Task meaning and existing `--epic` filter. The queue is a
separate, explicitly prioritized list of exact human selections. Run from the registered
Epic worktree (including a subdirectory), or supply all three target flags. A Task worktree
never implicitly selects its parent Epic.

```shell
ply workspace task queue list --project ply --repo ply --epic ply-agentic-workflow-support --format json
ply workspace task queue set --file /absolute/queue.json --format json
ply workspace task list --ready --project ply --repo ply --epic ply-agentic-workflow-support
ply workspace task prepare --next --check --format json
ply workspace task prepare --next --apply --confirm sha256:<preview-digest>
ply workspace task preparation show pre_<digest> --format json
```

`queue set` takes a strict `WorkspaceTaskQueueDraft@1` JSON object with `schema_version: 1`,
a unique `publication_key`, `project_id`, `repo_id`, `epic_id`, `expected_revision`, ordered
`entries`, `human_decision`, and `registry_upgrade`. Each entry is `{task_id, selection}`;
selection is null or the exact `{id, manifest_sha256}` of a Task-owned select event.
The human decision supplies `actor_claim`, UTC `decided_at_utc`, `source` (`human_cli` or
`explicit_human_instruction`), and `statement`. Copy `registry.required_upgrade` from
queue readback for the first format-4 mutation; use null thereafter. Migration preserves
an exact registry backup. Preview and list never migrate or create a queue.

A null, stale or blocked selection remains visible. `--ready` filters pending entries
without changing their ranks, and always shows the current preparation. A newer unselected
Spec draft does not replace the chosen solution. Set replaces the entire pending order;
it must omit the current reserved Task. Priority, selection, base or destination drift
invalidates a confirmation. Preview shows the exact separate branch/path and input hashes.
If preview includes a registry upgrade, apply additionally requires
`--upgrade-registry sha256:<registry-digest>`.

After a lost response, repeat the identical apply or inspect the preparation ID. One
reserved choice owns one durable intent. Normal later work does not erase `prepared`.
Partial or unknown effects retain that choice until explicitly recovered. Close it with
one metadata transition; neither command prepares the next Task:

```shell
ply workspace task queue advance --preparation pre_<digest> --expected-revision 2 --reason "Continue the queue"
ply workspace task queue release --preparation pre_<digest> --expected-revision 2 --reason "Stop this choice"
ply workspace task prepare --next --check
```

Advance requires a known prepared worktree. Release accepts known no-effect or prepared
state, keeps all refs/worktrees/evidence, and does not claim completion. Released Tasks
cannot be resumed or adopted in this version. Prepared and advanced do not imply technical
delivery, human QA, integration authority or agent start.

After separately integrating a delivered Task, explicitly register the clean descendant
Epic base before preparing the next Task:

```shell
ply workspace epic base show ply-agentic-workflow-support --repo ply --format json
ply workspace epic base update ply-agentic-workflow-support --repo ply --check --format json
ply workspace epic base update ply-agentic-workflow-support --repo ply --apply --confirm sha256:<preview-digest>
ply workspace epic base operation show ebu_<digest> --format json
```

The adopted base and old Task bindings stay immutable. Older selections become stale for
new effects. Record the next Task's new Spec revision against the current base, reuse
unchanged documents via snapshot references, assess it ready, explicitly select it, then
set the queue entry to that new selection. Historical Specs, results, QA and integration
retries retain their original digests. No command rebases worktrees, integrates code,
starts an agent, installs the CLI or selects a human QA outcome.

### Run a prepared Task in Herdr

`ply workflow run` previews and confirms one Herdr/Codex start from an existing
Task preparation, follows immutable round reports, and lets the named coordinator
accept, block or request bounded corrections in the same session. It returns a
reviewed report without publishing TaskResult or attesting provider inactivity.
Human QA and integration remain separate. See the [request, callbacks and recovery
guide](docs/workflow-run.md).

```shell
ply workflow run start --file /absolute/prepared-task-request.json --check
ply workflow run follow wfr_<digest> --timeout 60
ply workflow run review wfr_<digest> --file /absolute/review.json
```

### Run the local factory scenario with Codex exec

The separate `ply-factory-tests` repository owns the Python standard library rig,
scenario and independent oracle. Supply an explicitly built Ply binary; no installation
is required. `make build PLY_BUILD_OUTPUT=/absolute/private/ply` writes that binary.
The rig's `run --check` is read-only. Its later `run --human-start` is a separate human
authorization; development tests never prove native isolation or real provider behavior.

The same `workspace task run start --file` command supports request schema **2** with
`runtime.provider=codex`, `runtime.mode=exec`, `config_profile=null` and
`human_authority.start_surface=human_started_factory_test`. Its required `factory_test`
object contains `authorization_path`, `authorization_sha256`, `iteration` (1–2),
`task_slot` (1–2), `reasoning_effort=low` and `timeout_seconds` (1–240). The private
`PlyFactoryTestAuthorization@1` binds the output/work roots, actor, activity, scenario,
model, four executable hashes, effective policy hash and a deadline no more than 1200
seconds after first effect. Per-slot authorization files share one immutable root
budget ledger; changing a request key or authorization cannot reuse a spent slot.
Requests, preparations and independent fixture Git metadata must be inside the exact
iteration root. Symlinks, alternates, remotes and broader write policies are rejected.

```shell
ply workspace task run start --file /absolute/exec-request.json --check --format json
ply workspace task run start --file /absolute/exec-request.json --apply --confirm sha256:<preview-digest> --format json
```

Exec uses pipes, fresh `exec --json --ephemeral --ignore-user-config`, `--no-daemon`,
an explicit model, low reasoning and a private managed permission profile. Provider
JSONL and stderr are stored outside the recipient's write scope; each stream is limited
to 10 MiB and each JSONL line to 1 MiB. Stdout contains one Ply result. No resume,
interactive fallback, model upgrade or permission bypass is available.

Exec result and collect-preview use schema **3**, including a hash-bound reference to
runner-owned `ProviderCompletion@1` and runtime facts. `provider_terminal_event` can mark
the task inactive only after exactly one consistent native thread/turn completes with
all started items finished, complete valid JSONL, reaped exit 0 and an empty owned
process group. A completed command with status `failed` is closed and its failure stays
in the bound stream; a later correction and terminal turn can establish inactivity.
Verifier and report gates still decide whether the product result qualifies. Existing
negative completion projections remain unknown and byte-identical on readback.
Timeout, truncation, contradiction or unobservable lifetime stays
unknown. Timeout sends TERM to the still-owned group, waits five seconds, then sends
KILL only while ownership remains valid. Manual `--task-status` is rejected for exec.
An unknown run is preserved without restart or teardown.

Both modes require typed acceptance, a report with the mandatory `task-requirements`
artifact, real verifiers, correct clean Git state and agreement A. The explicit callback
`--context` works without inherited environment, and an unknown reported model is valid.
New reports undergo the existing workflow authority, verifier-binding and artifact
validation before the immutable report slot is reserved. Unknown budget or invalid
terminal input returns a concrete error that can be corrected in the same run. Accepted
reports retain their original bytes, including older unrepresentable reports. A terminal
rejection is reported with its workflow reason; clients must check `delivery.state` and
`terminal_sha256`, then the parent's collection result.
Exit codes remain 2/3/4/5. Existing schema 1 requests, interactive schema 2 results,
receipts, historical hashes and human task observations keep their original format.
Plan result control precedes separately authorized product QA or Ply integration.

### Start one interactive Codex run and preserve its return

After preparing a Task, use a private `ply.workspace.task-run-request` JSON document
(schema version 1) with the exact preparation ID and canonical SHA-256 of the whole
preserved preparation. Include its HandoffDraft@2, selected Spec and required inputs,
clean target, four total rounds, agreement A, human start claim and return policy.
The runtime binds physical Codex and Ply executables by SHA-256, the requested model,
an optional Codex **config** profile, and private evidence of the effective permission
policy. A permission profile name alone does not establish authority. Relative paths,
symlinks, unknown fields and unbound runtime overrides are rejected.

```shell
ply workspace task run start --file /absolute/request.json --check
ply workspace task run start --file /absolute/request.json --apply --confirm sha256:<preview-digest>
ply workspace task run show trn_<request-key-digest> --format json
ply workspace task run collect trn_<request-key-digest> --check --format json
```

Preview is read-only and prints a complete confirmation command. `--check` is the
default; `--apply` requires that exact confirmation. Interactive apply accepts text
output and a foreground macOS terminal. Codex inherits stdin, stdout and stderr;
Ply restores terminal ownership and state when possible. Other operating systems
support schema and readback, and reject interactive start before writing.

Ply reserves the request, run and target before its single launch attempt. Retrying
an already reserved request reads its preserved state, including after legitimate
Task changes. An unknown attempt is never automatically restarted. `show` and
`collect --check` create no files or locks. Recovery publication uses the confirmation
printed by `collect --check`; it cannot create a new launch attempt.

The recipient receives private instructions with the exact physical callback context:
`run accept RUN_ID --context /absolute/run/tmp/context.json --file /absolute/claim.json`
and `run report RUN_ID --context /absolute/run/tmp/context.json --file /absolute/report.json`.
These commands work in a clean tool shell. The legacy `PLY_TASK_RUN_CONTEXT` environment
value remains supported; when both locators are supplied they must be identical. The
flag does not expand authority: context bytes, run, request, session, cwd and the bound
Ply executable are still checked. No reply secret is placed in argv or stdout.

Acceptance schema version 2 separates the requested model from the recipient's claim.
A null reported model means unknown and does not alone block work. A known mismatch
cannot start. Runtime, profile, effective policy and scope must be known and match for
positive acceptance. Negative claims may honestly omit those facts and remain visible
even when the existing WF format cannot represent a StartReceipt. Acceptance@1 remains
readable; its literal model `unknown` also means unknown. Scope and effective OS policy
remain separate facts. No model is inferred from the requested `--model` value.
A received report does not mean that the interactive session has ended: the human
can exit Codex normally after reporting, and Ply collects automatically after Wait.
`Client process` describes the local client; `Task execution` describes the managed
task and names its evidence source. Without an explicit task observation, parent
collection preserves the return and exits 5 with `task_run_task_status_unknown`.
This is expected waiting for an observation; a return helper should handle that exact
code separately from other failures.

The return helper builds a private canonical `ply.workspace.task-run-task-status`
schema version 1 JSON document from the bound run and the human's explicit status
choice. Use `collect RUN_ID --task-status /absolute/task-status.json --check`, then
the complete apply/confirm command printed by that preview. The document binds the
run, request, Ply session, current journal tip, physical worktree, actor, UTC observation
time, visible group, task label, match count and optional native session ID. Its source
is `human_observation`, an attestation rather than provider authentication. An inactive
observation requires exactly one identified task in the matching group. A hidden native
ID is not required; a missing or ambiguous row remains unknown. The observation can be
collected with the human's final product assessment, but those are separate facts.
Check writes nothing. Apply rechecks the status digest and journal basis under lock;
changed evidence requires a new preview. Status history is immutable and exact retries
do not add duplicate observations.

Agreement A allows an initial execution, at most three correction rounds, 5400 active
seconds and two environment measures. Reported active time and observed elapsed time
remain separate. A complete return needs a valid started receipt, exact clean candidate,
required verification and review evidence, client exit 0 and proven process/group
quiescence, plus an inactive human task observation after the last received report.
A nonzero child exit, missing report, live group or unknown identity cannot qualify.
Ply records at most one qualifying TaskResult; raw blocked or incomplete reports remain
visible. Plan result control is the next gate, followed by separately authorized human
QA and integration. No queue advance or Epic base update happens automatically.
Active or unknown task execution retains the target reservation even after client exit.
Only a proven non-start or an inactive task with local quiescence can release it.
Historical qualified TaskResults remain readable and idempotent, with task execution
shown as unknown until observed. An explicit status can resolve that old reservation
without rewriting or requalifying the historical result. An unregistered old draft
still needs today's task status evidence before publication.

Run commands use exit 0 for preview/readback/idempotent return and qualified completion,
2 for invalid usage/schema, 3 for drift or policy conflicts, 4 for known start/I/O failure,
and 5 for missing or unknown process/delivery outcomes. The child exit is retained
separately. JSON output is one document on stdout; diagnostics go to stderr.

### Preserve a problem and choose a solution

A new Task atomically records its initial problem (P1) and a `spec_required` policy.
A problem is the need; a Spec is a proposed solution with abstract, functional and technical
parts. A readiness assessment and a human selection are separate immutable records.
Selection records the reported human claim. Worktree creation, agent start, product QA and
integration each remain separate actions.

```shell script
ply workspace task problem record explain-start-errors --file /absolute/problem.json
ply workspace task spec record explain-start-errors --file /absolute/solution-r1.json
ply workspace task spec assess explain-start-errors --file /absolute/ready-r1.json
ply workspace task spec select explain-start-errors --file /absolute/choose-r1.json
ply workspace task spec show explain-start-errors --spec explain-errors --revision 1 --format json
ply workspace task show explain-start-errors
```

The complete drafts below illustrate a Task that already has P1. Copy the exact predecessor
references from `problem show`, `spec show` and each mutation's JSON `outcome_ref`. Replace the
example hashes, WorktreeID, Git paths, OIDs and trees with those from your own readback. These
example claims are documentation data, not human approval. Write the following source files
as UTF-8 with LF line endings and a final newline; their exact hashes and sizes are included
in the corresponding drafts.

`/absolute/problem.md`:

```markdown
# Make start errors actionable

Explain the failed precondition and one safe next action.
```

`/absolute/solution.md`:

```markdown
# Explain start errors

## Purpose
Help the user understand a blocked start.

## Behavior
Show the failed precondition and one safe next action.

## Implementation
Change the error renderer only. Preserve exit codes and stored history.
No migration. Run start-error tests and inspect help. On failure, retain state.
```

`/absolute/problem.json` (replace all referenced identities with the exact prior readback):

```json
{
  "kind": "WorkspaceTaskProblemDraft@1",
  "schema_version": 1,
  "format": "json",
  "format_version": 1,
  "canonicalization": "RFC8785",
  "publication_key": "problem/p2",
  "task_id": "explain-start-errors",
  "registry_upgrade": null,
  "recorder": {
    "actor_claim": "example author",
    "control_surface": "local CLI",
    "recorded_at_utc": "2026-09-29T12:00:00Z"
  },
  "expected_previous": {
    "revision": 1,
    "manifest_sha256": "sha256:1111111111111111111111111111111111111111111111111111111111111111"
  },
  "origin": {
    "kind": "authored"
  },
  "title": "Make start errors actionable",
  "summary": "Explain the failed precondition and one safe next action.",
  "problem_document_id": "problem",
  "documents": [
    {
      "id": "problem",
      "source": {
        "kind": "file",
        "locator": "/absolute/problem.md",
        "sha256": "sha256:6216623ed48782578b432043428c91cbc98922da756fc23474331d45fa27c33a",
        "size_bytes": 90,
        "media_type": "text/markdown",
        "git_provenance": null
      }
    }
  ],
  "sources": [],
  "claims": [],
  "deadline": null,
  "change_reason": "Clarify the user-visible need."
}
```

`/absolute/solution-r1.json` (replace all referenced identities with the exact prior readback):

```json
{
  "kind": "WorkspaceTaskSpecDraft@1",
  "schema_version": 1,
  "format": "json",
  "format_version": 1,
  "canonicalization": "RFC8785",
  "publication_key": "solution/r1",
  "task_id": "explain-start-errors",
  "registry_upgrade": null,
  "recorder": {
    "actor_claim": "example author",
    "control_surface": "local CLI",
    "recorded_at_utc": "2026-09-29T12:00:00Z"
  },
  "spec_id": "explain-errors",
  "expected_previous": null,
  "problem": {
    "revision": 2,
    "manifest_sha256": "sha256:1111111111111111111111111111111111111111111111111111111111111111"
  },
  "title": "Explain start errors",
  "parts": {
    "abstract": {
      "state": "present",
      "reason": null,
      "documents": [
        {
          "document_id": "solution",
          "section": null
        }
      ]
    },
    "functional": {
      "state": "present",
      "reason": null,
      "documents": [
        {
          "document_id": "solution",
          "section": null
        }
      ]
    },
    "technical": {
      "state": "present",
      "reason": null,
      "documents": [
        {
          "document_id": "solution",
          "section": null
        }
      ]
    }
  },
  "documents": [
    {
      "id": "solution",
      "source": {
        "kind": "file",
        "locator": "/absolute/solution.md",
        "sha256": "sha256:3827cb23c3fe23f45e2c544a42061be969b31baf66e544ae5af7e5426a62c5aa",
        "size_bytes": 316,
        "media_type": "text/markdown",
        "git_provenance": null
      }
    }
  ],
  "supporting": [],
  "requirements": [
    {
      "id": "f-01",
      "functional_refs": [
        {
          "document_id": "solution",
          "section": null
        }
      ],
      "acceptance": "A blocked start names its failed precondition and one safe next action.",
      "verification_ids": [
        "start-errors"
      ],
      "technical_refs": [
        {
          "document_id": "solution",
          "section": null
        }
      ]
    }
  ],
  "removed_requirement_ids": [],
  "phases": [
    {
      "id": "implement",
      "purpose": "Change and verify error rendering.",
      "requirement_ids": [
        "f-01"
      ],
      "entry_criteria": [],
      "exit_criteria": [
        "Start-error tests pass."
      ],
      "verification_ids": [
        "start-errors"
      ]
    }
  ],
  "implementation_basis": {
    "project_id": "ply",
    "repo_id": "ply",
    "git_common_dir": "/absolute/ply/main/.git",
    "epic_id": "ply-agentic-workflow-support",
    "parent_worktree_id": "wt_66666666666666666666666666666666",
    "parent_ref": "refs/heads/ply_agentic_workflow_support",
    "parent_oid": "2222222222222222222222222222222222222222",
    "parent_tree": "3333333333333333333333333333333333333333",
    "start_oid": "2222222222222222222222222222222222222222",
    "start_tree": "3333333333333333333333333333333333333333"
  },
  "dependencies": [],
  "change_reason": "Record the proposed solution."
}
```

`/absolute/ready-r1.json` (replace all referenced identities with the exact prior readback):

```json
{
  "kind": "WorkspaceTaskSpecAssessmentDraft@1",
  "schema_version": 1,
  "format": "json",
  "format_version": 1,
  "canonicalization": "RFC8785",
  "publication_key": "solution/ready-r1",
  "task_id": "explain-start-errors",
  "registry_upgrade": null,
  "recorder": {
    "actor_claim": "example author",
    "control_surface": "local CLI",
    "recorded_at_utc": "2026-09-29T12:00:00Z"
  },
  "spec_id": "explain-errors",
  "spec": {
    "revision": 1,
    "manifest_sha256": "sha256:1111111111111111111111111111111111111111111111111111111111111111"
  },
  "expected_previous_assessment": null,
  "outcome": "ready",
  "reason": "The documented solution is ready for a separate human choice.",
  "open_questions": [],
  "checks": [
    {
      "id": "acceptance_coverage",
      "outcome": "pass",
      "reason": "Reviewed against the preserved problem, solution and local basis.",
      "evidence_document_ids": []
    },
    {
      "id": "implementation_basis",
      "outcome": "pass",
      "reason": "Reviewed against the preserved problem, solution and local basis.",
      "evidence_document_ids": []
    },
    {
      "id": "problem_coverage",
      "outcome": "pass",
      "reason": "Reviewed against the preserved problem, solution and local basis.",
      "evidence_document_ids": []
    },
    {
      "id": "recovery",
      "outcome": "pass",
      "reason": "Reviewed against the preserved problem, solution and local basis.",
      "evidence_document_ids": []
    },
    {
      "id": "scope_and_phases",
      "outcome": "pass",
      "reason": "Reviewed against the preserved problem, solution and local basis.",
      "evidence_document_ids": []
    },
    {
      "id": "three_parts",
      "outcome": "pass",
      "reason": "Reviewed against the preserved problem, solution and local basis.",
      "evidence_document_ids": []
    }
  ],
  "documents": []
}
```

`/absolute/choose-r1.json` (replace all referenced identities with the exact prior readback):

```json
{
  "kind": "WorkspaceTaskSolutionSelectionDraft@1",
  "schema_version": 1,
  "format": "json",
  "format_version": 1,
  "canonicalization": "RFC8785",
  "publication_key": "solution/choose-r1",
  "task_id": "explain-start-errors",
  "registry_upgrade": null,
  "recorder": {
    "actor_claim": "example author",
    "control_surface": "local CLI",
    "recorded_at_utc": "2026-09-29T12:00:00Z"
  },
  "expected_previous_selection": null,
  "action": "select",
  "solution": {
    "spec_id": "explain-errors",
    "spec": {
      "revision": 1,
      "manifest_sha256": "sha256:1111111111111111111111111111111111111111111111111111111111111111"
    },
    "problem": {
      "revision": 2,
      "manifest_sha256": "sha256:1111111111111111111111111111111111111111111111111111111111111111"
    },
    "assessment": {
      "id": "asm_44444444444444444444444444444444",
      "manifest_sha256": "sha256:1111111111111111111111111111111111111111111111111111111111111111"
    }
  },
  "reason": "Record the explicit human choice.",
  "human_decision": {
    "actor_claim": "example human claim; replace with the actual instruction",
    "decided_at_utc": "2026-09-29T12:00:00Z",
    "source": "explicit_human_instruction",
    "statement": "Select explain-errors revision 1 for the separately authorized delivery."
  }
}
```

All new leaves accept `--format text|json`. `problem show` optionally accepts `--revision`;
`spec show` requires `--spec` and `--revision`. `spec list` shows revisions and choices.
`spec withdraw --file /absolute/withdraw.json` uses the same selection draft: a new publication
key, `expected_previous_selection` from the latest event, `action="withdraw"`, `solution=null`,
and an explicit human withdrawal claim. It preserves all previous choices and starts or stops
no process.

To record r2, name r1 in `expected_previous`, keep the same `spec_id`, and use a new key.
Unchanged documents may use `source={"kind":"snapshot","manifest_sha256":"<exact prior
manifest digest>","document_id":"<prior document ID>"}` instead of a file source. A new draft
never changes the selected revision. A new problem, latest assessment, selection or withdrawal
can make the previous basis stale. Changing a problem requires a new solution revision and
human choice before another start; an old result remains historical evidence.

Imports preserve exact source bytes under private `.ply/task-content` paths. The source may
later move or disappear. Reads rehash snapshots and never repair or migrate them. A repeated
publication key with identical draft content finds the original outcome before predecessor
checks or source reads; different content conflicts. After an uncertain publication, inspect
`ply workspace task publication show explain-start-errors --key solution/r1 --format json`.
An identical retry may confirm durability without creating another revision. JSON mutations
also emit their observed outcome when an error follows a possible registry replacement;
inspect that outcome together with stderr and the nonzero exit.

New stores use format 3. Existing format-1/2 registries and legacy Task results remain readable.
Creating another Task or explicitly recording a legacy Task's first problem requires the exact
current registry SHA-256: `task create --upgrade-store sha256:<64 lowercase hex>`, or
`registry_upgrade={"from_version":1,"registry_sha256":"sha256:<64 lowercase hex>"}` in the
problem draft (use 2 for format 2). Ply preserves a byte-for-byte backup and upgrades atomically.
An identical legacy create retry never fabricates P1. A legacy summary import uses
`origin={"kind":"legacy_summary_import","legacy_summary_sha256":"<canonical summary digest>"}`,
its original title/description, empty documents/sources/claims, null deadline and
`problem_document_id="problem"`; it makes no historical approval claim.

Drafts and manifests are strict JSON, at most 256 KiB each. Documents are 1 byte–4 MiB of UTF-8
text/plain or text/markdown, without BOM or NUL. Limits are 32 documents, 128 requirements,
32 phases, 32 sources and 128 claims per applicable manifest. Titles allow 128 codepoints,
summaries 2048, and reasons/acceptance 2000. Redact known secrets before import; a withheld
source can preserve a claim without preserving its bytes. No remote URL is fetched.

Task show at format 3 includes the original resource view, integration state, current problem,
choice and every historical result. Its readiness is observation, never execution authority.
For a selected Task, prepare a schema-version-2 handoff with the exact `task_spec_binding` and
all `required_inputs` from Spec show. The closure includes all four manifests and every preserved
document, with at most 96 unique documents and 64 MiB of document bytes. Additional spec/design
inputs cannot replace the selected solution. Declare every requirement's verification ID and
explicitly require a managed application/json `task-requirements` artifact in the procedure.

The version-2 start repeats the exact basis; Ply adds and signs its own observation under the
existing Project/work-item locks. Version-1 handoffs remain supported for general or legacy
work, and version-2 null bindings only for those targets. Required Tasks cannot bypass the
choice through version 1 or null. A version-1 terminal works with either matching handoff/start
pair. Accepted start retries preserve the original slot after a later choice change.

A bound terminal's `task-requirements` artifact uses `WorkspaceTaskRequirementEvidence@1` with
`task_id`, `spec_id`, the exact `spec` revision reference, `result_oid`, `result_tree`, and one
`requirements` entry per selected requirement. Each entry has `id`, `outcome` (passed, failed,
not_run or unknown), sorted `verifier_ids`, sorted `artifact_ids`, and `reason`. References must
resolve in the same terminal and cannot refer to the requirement artifact itself. Missing or
non-passing functional coverage cannot produce a green Task result. Report failures honestly.

Result records bind the immutable accepted start, even after P2. Historical inspection does
not change merely because a later problem or choice exists. Spec-aware integration plans hash
the current relevance guard; a stale basis blocks a new attempt. Recovery of an already
attempted integration observes Git first and reports later relevance separately. Legacy plans
retain their original bytes and digests. Technical pass never means human QA or integration.

## Task process journal

An existing Task can retain process notes and steps across several native runs.
The journal combines validated native Task, problem, solution, preparation, run,
result, human QA and integration facts with immutable self-reported contributions.
It changes no native state and authorizes no next transition.

```sh
ply workspace task journal show my-task
ply workspace task journal show my-task --view timeline
ply workspace task journal show my-task --view details --order recorded
ply workspace task journal show my-task --run trn_<digest> --actor developer-session
ply workspace task journal show my-task --format json > /absolute/snapshot.json
ply workspace task journal show my-task --snapshot /absolute/snapshot.json --view timeline
```

`show` discovers the workspace from the current directory. It never creates a
journal directory, lock, cache, receipt or repair. `--snapshot` instead validates
and renders the saved snapshot without workspace discovery, source reads or a
new observation. JSON always contains the full evidence; filters change only
`selection`. Offline rendering preserves `snapshot_id` and `as_of`.

The banner separates Transport, Reported, Technical, Human QA and Integration.
Old judgments remain attached to their exact candidate. A dirty, missing or new
candidate does not inherit an earlier pass. Coverage describes the available
sources read, not a complete work or conversation history. Missing or changed
optional sources produce explicit partial coverage and a successful readback;
an unknown Task or invalid snapshot fails.

The timeline uses a shared UTC axis and one lane per actor/run. It preserves
occurrence and registration clocks, original offsets, source sequence and
uncertainty. `No end recorded`, `Unknown time` and `Time conflict` are explicit.
Only compatible known step endpoints have durations. Parent and child intervals
and parallel lanes are never added into an active-time or productivity total.
Waiting requires an explicit waiting step with reason, dependency and next actor.

To append a note, put the following in `/absolute/event.json`, replacing the Task,
source path and the source's 64 lowercase hexadecimal SHA-256. Every field is
required, including explicit `null` values and empty arrays. The source must be
an existing physical regular file whose bytes match the digest.

```json
{
  "kind": "ply.workspace.task-journal-event-input",
  "schema_version": 1,
  "publication_key": "delivery-1:verified",
  "task_id": "my-task",
  "actor": {"id": "developer-session", "role": "developer", "session_id": "external-session-1"},
  "activity_id": "delivery-1",
  "run_binding": null,
  "occurred_at": "2026-10-03T10:00:00+02:00",
  "time_basis": {"kind": "reported", "clock": "session-clock", "precision": "second", "uncertainty": null},
  "sources": [{"locator": "/absolute/test-report.txt", "sha256": "<64-lowercase-hex-digest>"}],
  "type": "note",
  "title": "Verification completed",
  "detail": "",
  "step": null,
  "outcome": "pass",
  "candidate": null,
  "relations": [],
  "waiting": null
}
```

```sh
ply workspace task journal record my-task --file /absolute/event.json --format json
```

`detail` is empty in v1. A native `run_binding`, when present, contains all four
exact fields: `run_id`, `request_sha256`, `preparation_id`, `preparation_sha256`.
An external session uses `null` and `actor.session_id`; a matching title is not a
run binding. Actor identities and roles are self-reported metadata.

Use `step_started` and `step_finished` with the same `{id, kind, parent_step_id}`
and actor/run for an interval. Step kinds are `clarification`, `planning`,
`implementation`, `verification`, `review`, `correction`, `coordination`, and
`waiting`. An original start and finish occur at most once; another correction
round gets a new step. Relations (`responds_to`, `corrects`, `verifies`,
`supersedes`) refer to existing events in this Task. A correcting note or decision
may supersede another contribution; native judgments cannot be superseded.
Conflicting successors are retained as a conflict, not resolved by timestamp.

An observer uses the same common fields, changes `kind` to
`ply.workspace.task-journal-observation-input`, and replaces the event-specific
fields (`type` through `waiting`) with this proposal body. Replace the target with
an `event_id` returned by `record` or `show`:

```json
{
  "type": "proposal",
  "proposal_event_id": null,
  "targets": {"event_ids": ["<existing-event-id>"], "step_ids": [], "interval": null},
  "finding": "The recorded review followed verification",
  "hypothesis": "An earlier review may shorten the next correction",
  "action": "Compare the next delivery with its review evidence",
  "owner": "planner",
  "next_signal": "The next comparable delivery",
  "decision": null,
  "assessment": null
}
```

```sh
ply workspace task journal observe my-task --file /absolute/observation.json --format json
```

A follow-up names its `proposal_event_id`. A `decision` uses `accepted` or
`rejected`; `applied` requires an unconflicted accepted decision; `assessment`
requires an applied record and uses `better`, `unchanged`, `worse`, or
`inconclusive`. Facts, hypotheses and proposals remain separate. Observer claims
of success never become native Technical, Human QA or Integration judgments.

Both append commands accept one UTF-8 JSON object of at most 64 KiB, reject
duplicates and unknown fields, and share a Task-local `publication_key` namespace.
An identical canonical retry returns the same event, registration time and digest
with `created=false`, even after a source disappears. Different input conflicts.
A post-publication error reports an unknown result and the key/readback path;
retry only the contribution, not the work it describes. Records are published
atomically under a local append lock in private `.ply/task-process/v1/` storage.
Sources are referenced and hashed, never executed or copied into the journal.
Do not submit transcripts, credentials, secrets or arbitrary payloads.

## Workflow handoffs

Workflow handoffs provide a local, immutable file protocol for giving one bounded task to a
human-started agent and receiving bound start and terminal reports. Initialize a workspace and
register its Project and Repo first, then create a handoff from an absolute draft path:

```shell script
ply workflow handoff create --file /absolute/handoff-draft.json
```

A successful create prints this five-line human start block:

```text
Created agent handoff hnd_<32-lowercase-hex>.
Purpose: <goal title>
Working directory: <exact target worktree>
Handoff: <absolute immutable handoff.json path>
Next action: Open a fresh recipient agent in the working directory and tell it: "Read and execute the handoff at <absolute immutable handoff.json path>."
```

The human—not Ply—opens the recipient in that exact working directory and supplies the one
handoff locator. The recipient submits `submit-start` before any target effect and
`submit-result` after completing or stopping. Use `show <handoff-id>` for a human-first summary
and `inspect --handoff <path> --format json` for a machine-first, redacted projection. Exact raw
handoff bytes expose the reply secret and therefore require both `--raw handoff` and
`--acknowledge-secret-exposure`.

Identical create/start/result retries return the existing immutable record; different bytes for
the same publication or slot are preserved as a conflict instead of replacing history. A reported
`complete` result never authorizes QA, integration, a pull request, or merge. The local file
capability binds replies to the handoff, but it is not strong principal attestation.

The following is a complete compact draft example; replace every absolute path, Git OID, and
identifier with observed values for the registered target. It is the full strict schema, not a
handwritten convenience format:

```json
{"activity_key":"example/readme","authority":{"allowed_effects":[{"id":"write-readme","max_occurrences":1,"scope":{"directory_prefixes":[],"kind":"filesystem","paths":["README.md"]},"sequence":1,"type":"filesystem_write"},{"id":"run-tests","max_occurrences":1,"scope":{"kind":"command","procedure_ids":["implement"],"verifier_ids":["tests"]},"sequence":2,"type":"command_execute"}],"forbidden_effects":[{"reason":"Network access is outside this local task.","type":"network"},{"reason":"Merge requires a later human gate.","type":"merge"}],"human_gates":["human_task_qa","local_integration"]},"binding_request":{"expected_oid":"0123456789abcdef0123456789abcdef01234567","project_id":"example","repo_id":"example","status_policy":{"mode":"clean"},"target_ref":"refs/heads/example_work","target_worktree":"/absolute/example/worktree"},"budget":{"max_rounds":2,"round_definition":{"command_retry_consumes_round":false,"retry_condition":"only_if_no_effect_started","unit":"implementation_or_review_fix_iteration"}},"canonicalization":"RFC8785","format":"json","format_version":1,"goal":{"done_when":"The bounded edit is complete and the declared verifier passes.","objective":"Update the documented file without effects outside the contract.","recipient_role":"Delivery agent","title":"Update the example documentation"},"inputs":[],"kind":"ply.workflow.handoff-draft","procedure":[{"id":"implement","instruction":"Make the bounded documentation change and review the diff.","required_before":[]}],"publication_key":"example/readme/run-1","recipient":{"principal_id":"codex-delivery-agent","principal_kind":"human_started_agent","runtime_constraints":["local"]},"reporting":{"meaning_max_codepoints":600,"required_start_fields":["acceptance","binding","contract_digests","issues","observed_inputs","observed_project","observed_target","observed_workspace","principal","receipt_id","sandbox"],"required_terminal_fields":["artifacts","binding","evidence_gaps","final_target","forbidden_effects_observed","meaning","observed_effects","principal","reported_outcome","result_id","review","rounds_used","start_binding","stop_reasons","summary","verifier_results"],"summary_max_codepoints":240},"schema_version":1,"stop_conditions":[{"description":"Stop when a product decision is required.","type":"product_decision_required"},{"description":"Stop before expanding scope or authority.","type":"scope_or_authority_expansion"},{"description":"Stop when target or input drift is observed.","type":"target_or_input_drift"},{"description":"Stop when an effect may be unknown or partial.","type":"unknown_or_partial_effect"},{"description":"Stop when unexpected sensitive data is observed.","type":"unexpected_sensitive_data"},{"description":"Stop when the round budget is exhausted.","type":"round_budget_exhausted"}],"verifiers":[{"argv":["go","test","./..."],"cwd":"/absolute/example/worktree","env":[],"evidence":{"binding":"target_oid","capture_stderr":false,"capture_stdout":false,"classification":"workspace_internal"},"expected_exit":0,"id":"tests","stop_on_failure":true}]}
```

## Install
```shell script
make install
```

## Build
```shell script
make build
```

## Help
```shell script
ply
```

## Report-ready Slack notifications

`ply workflow notification` previews and sends a report-ready human gate to one fixed
Slack channel. When separately authorized, create an Incoming Webhook for that channel,
put its URL in a named environment variable such as `PLY_SLACK_WEBHOOK` using your usual
secret practice, and create a route JSON containing the variable name and a private
state root, without the URL. No global Ply configuration or native workspace is needed.

See [the complete validated route/request example and check → apply → show procedure](docs/workflow-notification.md).
A preserved transport acknowledgement is not report control or human QA. Unknown
outcomes cannot be resent; local report and return instructions remain available.


### Agent completion and necessary feedback notifications

`ply workflow notification send` also accepts strict request schema version 2 for
`agent_finished`, `agent_stopped`, and `feedback_required`. These bind an external
activity/run and a preserved `PlyAgentNotificationEvent@1`; completion requires an
actual report, while a before-start stop needs no fabricated report or receipt.
The event is a sender claim. Transport acknowledgement does not attest product
correctness, result control, human QA, or that anyone read the message. Existing
v1 human-gate requests and stored state retain their contract.

The repository provides [ply-agent-notify](skills/ply-agent-notify/SKILL.md), with a
[task context and operational procedure](skills/ply-agent-notify/references/usage.md).
An authorized task adopts that procedure to preserve its return/question, preview,
apply at most once and read back through Ply. Optional questions, internal corrections
and an ending still waiting for the same answer send nothing. The same logical event
keeps its ID across restarts and port revisions. Unknown delivery never permits a new
ID, route or state root to resend. The skill provides no crash watcher.

During a separately authorized installation activity, run `make install-agent-notify`,
then run the installed skill's `scripts/bootstrap.py`. Bootstrap accepts the existing
webhook from `PLY_SLACK_WEBHOOK_URL` or hidden terminal input and stores it privately
at `/Users/perottochristensen/.config/ply/notifications/ply-log/webhook`. It creates
the adjacent non-secret `route.json` for `ply-log` and uses the fixed private state
root `/Users/perottochristensen/.local/state/ply/agent-notifications/ply-log`.
Normal helper invocations prefer the environment, otherwise read only that named
0600 file in its 0700 directory, and pass the credential only to Ply send/check.
Invalid environment values never fall back; conflicts and unsafe files are preserved.
No webhook belongs in a prompt, argv, repository, report or ordinary log.

Run `make test-agent-notify` for isolated installation, bootstrap and actual
helper/CLI interaction with controlled fake transport. Set `AGENT_NOTIFY_TEST_ROOT`
to a new private fixture path and `AGENT_NOTIFY_DEST` for fixture-only installation.
The technical tests use no Slack service. Personal installation and the useful,
human-observed agent journey remain a separate QA activity.
