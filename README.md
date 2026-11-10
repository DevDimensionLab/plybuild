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
  
For a prebuilt macOS arm64 CLI, use the [consumer installation and recovery guide](docs/cli-installation.md).
It needs neither Go nor a source checkout and can replace an old CLI that cannot read the workspace.
See also [shared registry compatibility and controlled upgrades](docs/registry-compatibility.md).

Building from source requires [Go 1.26.10](https://go.dev/doc/install).

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
  agent        Show remaining Codex and Claude subscription quota
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

## Planning repositories

Create a standalone planning Git repository at an explicit root, without starting an agent:

```shell
ply workflow epic plan --repo /work/example/main --root /work/example --check
ply workflow epic plan --repo /work/example/main --root /work/example --goal "Plan the next release"
ply workflow epic plan --project ply --language nb
```

Project mode selects all registered members and defaults to the registered wrapper. Explicit
`--repo` paths require `--root` and work without a Ply workspace. The default document language
is English (`en`); Norwegian Bokmål is `nb`. CLI messages and help remain English.

Creation validates and renders first, then creates `<root>/planning` and one clean initial
commit on `main`. `--check` performs a read-only preview. Every existing target stops creation;
existing planning repositories receive a reminder to use the existing plan. Partial failures
are preserved for inspection. See the [command guide](docs/workflow-epic-plan.md) for all flags,
Git prerequisites, machine output and recovery behavior.

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

Retrying the same registration is idempotent. To append one member to an existing project:

```shell script
ply workspace project add-member plyapp --repo plyapp-pwa=../plyapp-pwa/main
ply workspace project show plyapp --format json
ply workspace project list --format json
```

`add-member` preserves the project name, wrapper, existing members and stored repository
locators. A known repository ID with the same physical Git common directory can be reused;
identical retries (including symlinks and linked worktrees) do not rewrite registration.
Conflicting IDs or Git identities and unavailable worktrees fail without partial membership.
This metadata operation accepts dirty repositories and paths outside the wrapper, requires no
Epic, and leaves Tasks, queues, Runs, history, Git refs and source unchanged. It does not relocate
or remove members, rebind execution, or grant cross-repository Epic support. `project init`
remains outside these commands.

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
directory. It performs no Git change and reserves no directory: several logical Epics may share the
same clean integration worktree and branch, for example a project's `main` checkout. Exclusivity
belongs to concrete work instead. Each Task gets its own branch and worktree, and an integration
into a shared branch blocks other Epics' base updates and preparations on that branch until it
completes. An older registration whose branch is no longer checked out remains visible as
`epic_ref_stale` history. It does not block a new adoption and gives no authority to continue old
work. When the current directory is the base of several Epics, Ply picks the one whose recorded
branch is checked out. If that is still ambiguous, it refuses and prints the exact
`--project <id> --repo <id> --epic <id>` choices.

Each Task belongs to an Epic and binds exactly one registered repository and Git common directory:

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

To continue an existing committed descendant, create an ordinary Task/worktree and use
human-owned `git merge --ff-only S` with the exact selected source commit. See the
[Task and Git recipe](docs/task-worktree-adoption.md) for the checks and limits.
Task worktree adoption and `release-adoption` are removed; Epic adoption remains.

Ply persists a durable create intent before the additive branch/worktree operation. Identical
retries recover safe no-effect, matching branch-only, or exact-effect outcomes with the same IDs.
Partial or unknown effects are preserved for explicit reconciliation and are never reset, removed,
or otherwise cleaned up automatically.

Use `ply workspace epic list`, `ply workspace epic show <id>`, `ply workspace task list`, and
`ply workspace task show <id>` for human readback. Both show commands accept `--format json` for
versioned deterministic machine readback. `worktree_ready` only records a clean local bootstrap
binding: it does not start an agent, select a workflow, or grant execution authority.

Every Task reaches its registered Epic through the same operation, with or without a
Goal. Commit the Task, run `ply workspace task integrate TASK --expected-result-oid
RESULT_OID --expected-parent-oid PARENT_OID --check`, then use the returned digest with
`--apply --confirm <digest>`. The human decides when to integrate; Ply records what happened
in Git. Deliver Tasks serially per Epic and update the Epic base after each delivery; see
[Explicit Goal agent start](docs/goal-agent-start.md). A Task adapted by hand to an
updated Epic base can be integrated from that base with `--refresh-parent`; see
[Integrer en tilpasset Task fra oppdatert Epic-basis](docs/task-integration-refresh-parent.md).

Apply permits one local `ff-only` move
of the exact registered parent, including an explicitly selected `main` or `master`.
It does not fetch, push or clean worktrees.

A newly created Epic may have only its creation reflog entry before integration.
A confirmed effect still requires the real result and preceding parent entries,
with the exact integration attempt marker in the reflog.

A format-1 work-item store is not integrated. A no-effect attempt can only be retried through a new check
that names `--retry-after`. Lost responses are reconciled from the durable authority, attempt,
Git state, and reflog; partial or unknown effects stop for read-only recovery control.

Format-2 `task show` presents the same integration state as check/apply in text or canonical JSON.
Top-level `next_action` er aktuell veiledning fra ferske observasjoner; `integration_result.next_action`
er historisk veiledning lagret med effekten. Lesende check eller show kan endre den første uten å endre
den siste, planen, digesten, confirmation, klassifikasjonen eller `git_changed`. Bruk `git_changed` og
top-level `next_action` for aktuell status. Tom `argv` betyr at det navngitte hinderet krever menneskelig
avklaring; det gir ikke adgang til apply. Integration v1
requires Git 2.45.0 or newer, a `sha1` or `sha256` object format, and the `files` ref backend.

Sub-tasks, additional repository anchors on an Epic, editing, rebinding, cleanup, and
WorkflowRun creation are outside the resource-binding commands. Refresh is explicit only through
`workspace epic base update`; it appends a metadata base and never moves Git refs or old worktrees. Problem and solution revisions
are recorded separately, as described below.

### Prioritize Goals

`task list` keeps its registered-Task meaning and existing `--epic` filter. The queue is a
separate, explicitly prioritized list of registered Goal references: `goal register` appends
one, `queue set` (draft `WorkspaceTaskQueueDraft@2`) orders them, and `queue list` and
`workflow status` show them. Run from the registered Epic worktree (including a
subdirectory), or supply all three target flags.

```shell
ply workspace task queue list --project ply --repo ply --epic ply-agentic-workflow-support --format json
ply workspace task queue set --file /absolute/queue.json --format json
```

### Human-owned process

Ply provides explicit workspace, Task, Git and agent-start operations. The compiled
`workflow execute`, `workflow worker`, `workflow week`, `workflow control`,
`workflow delivery`, `integration` and factory runners have been removed.
Choose and invoke the individual operations yourself. Task integration no longer
requires a selected Spec, a TaskResult, QA or an acceptance policy. Source, target,
expected commit, concurrency and uncertain-effect checks still protect Git changes.

Ply does not set Codex approval policy/reviewer or grant native project trust.
Process configuration at project/repository level is a future design direction;
no replacement process engine or configuration layer is included here.

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

An explicitly supplied version-2 document basis is checked against its exact preserved
inputs under the Project/work-item locks. A Task target can also be used without a Spec
binding; its physical path, ref and Git identity are still checked. A later selection change
does not replace a supplied historical basis. Accepted start retries preserve the original slot.

A bound terminal's `task-requirements` artifact uses `WorkspaceTaskRequirementEvidence@1` with
`task_id`, `spec_id`, the exact `spec` revision reference, `result_oid`, `result_tree`, and one
`requirements` entry per selected requirement. Each entry has `id`, `outcome` (passed, failed,
not_run or unknown), sorted `verifier_ids`, sorted `artifact_ids`, and `reason`. References must
resolve in the same terminal and cannot refer to the requirement artifact itself. Missing or
non-passing functional coverage cannot produce a green Task result. Report failures honestly.

Handoff-based results retain their immutable start and supplied input basis. Direct command
receipts retain the observed command, exit, source commit and output hashes without inventing
review or QA. A later Problem or selection does not veto an explicitly chosen integration.
Recovery observes the preserved Git attempt before deciding whether any effect remains unknown.

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
The JSON `duration_seconds` field preserves fractional seconds, including subsecond
intervals. Offline validation checks the duration against its original timestamp
evidence. Existing whole-second snapshot digests remain compatible; other journal
numeric fields retain their strict integer contract.

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

## Assignments

An Assignment is a short work description stored in a Project: `ASSIGNMENT-<n>`, a
title, text read byte for byte from stdin, one or more repo IDs of the Project and
optional topics. It is never changed after `add`; a correction is a new Assignment that
names the old one in its text.

```sh
printf '%s\n' 'Vis repo-tags i backlog list.' |
  ply assignment add --title 'Vis repo-tags i backlog list' --repo ply --topic backlog
ply assignment show ASSIGNMENT-7 --format json
ply assignment list
```

Without `--project` the single Project whose wrapper contains the current directory
is used. Files live in the Project's `data_paths.assignments` (default
`<project wrapper>/assignments/`). `list` reports unreadable files under `invalid` and
still exits 0. If `add` loses its answer, run `list` and look up the title before
trying again.

An agent takes an Assignment before working on it and marks it done when finished.
`done` itself moves it to the archive view; there is no archive command:

```sh
ply assignment show ASSIGNMENT-7 --format json       # lifecycle.revision, e.g. 0
ply assignment take ASSIGNMENT-7 --revision 0 --as claude:<session-id>
ply assignment done ASSIGNMENT-7 --claim <claim>     # or: release --claim <claim>
ply assignment list --view archive                   # available (default), taken, archive, all
```

`take` prints a claim; `release` and `done` need it. `--revision` is the revision
`show` just printed, so of two simultaneous takes exactly one succeeds and the other
gets `assignment_conflict`. `release --claim` makes it available again; `done` is
final, and a follow-up is a new Assignment. If the holder is gone (crash, lost
session), a human recovers it with
`ply assignment release ASSIGNMENT-7 --revoke --revision <r> --as <name> --reason <text>`;
the skill forbids agents to revoke. On `result unknown, run show`, run `show` before
retrying. The original Assignment file is never changed; the lifecycle is an
append-only record in `<assignments>/.lifecycle/ASSIGNMENT-<n>.json`. A record that
cannot be read is listed under `invalid` as `assignment_lifecycle_invalid`, its
Assignment appears in no view, and Ply refuses to change it; repair is a human task.

The repository-owned skill [ply-assignment](skills/ply-assignment/SKILL.md) lets a
human start an agent session on one Assignment: the agent reads it itself with
`ply assignment show`, takes it, and marks it done or releases it. Install it explicitly; there is no default destination and an
existing different `SKILL.md` is reported as `existing_install_conflict`, never
overwritten:

```sh
make install-assignment-skill ASSIGNMENT_SKILL_DEST=<repo>/.claude/skills/ply-assignment
```

Then start Claude in that repository with the Assignment as the initial prompt, either
interactively, `claude "/ply-assignment ASSIGNMENT-7"`, or non-interactively. This
form was qualified with Claude Code 2.1.309, as a project skill in an isolated
workspace; add the tools the work itself needs, such as `Write`:

```sh
claude -p "/ply-assignment ASSIGNMENT-7" --setting-sources project,local \
  --output-format stream-json --verbose --max-turns 20 \
  --allowedTools "Bash(ply assignment show:*)" "Bash(ply assignment take:*)" "Bash(ply assignment done:*)" "Bash(ply assignment release:*)" "Bash(ply workspace task:*)" "Read"
```

`--setting-sources project,local` keeps an older user-level
`~/.claude/skills/ply-assignment` from replacing the project skill. The agent labels
its take `claude:<session-id>`. A run that stops after `take` but before `done` or
`release`, for example at `--max-turns`, leaves the Assignment `taken`: resume that
session with `claude -p --resume <session-id> "/ply-assignment ASSIGNMENT-7"` (same
flags), which adopts the claim by its label, or, if the session is lost, a human
revokes the claim as above.

## Agent quota

`ply agent` shows the remaining subscription quota that Codex and Claude report
for your account on this machine, once, in the order Codex, Claude. `ply agent
--json` (or `ply --json agent`) prints one `schema_version: 1` document with a
`providers` list. It needs no Ply project, profile or workspace and changes no
settings.

```text
Codex   5-hour: 62.0% remaining (resets 2026-11-07T20:00:00Z)   Weekly: 18.0% remaining
        plan pro · observed 2026-11-07T17:59:55Z (5s ago)
Claude  not configured — configure claude-statusline-relay (see ply agent --help)
```

- Codex is read directly with a short-lived `codex app-server`
  (`account/rateLimits/read`, 15 s deadline) using the existing Codex login.
  No thread, turn or inference is started. Missing `codex` is `not_configured`.
- Claude reports quota only to its status line. The separate
  `claude-statusline-relay` stores only `five_hour`/`seven_day` in a private
  snapshot (`~/Library/Application Support/AI Usage/claude-rate-limits.json` on
  macOS, `$XDG_STATE_HOME/ai-usage/claude-rate-limits.json` on Linux; mode
  0600) and then runs your original status line command with the same input.
  The snapshot updates only while Claude Code runs and is stale after one hour.
- Statuses are `ok`, `not_configured`, `waiting`, `not_reported`, `unsupported`
  and `error`. The exit status is 1 only when a source has status `error`. A
  passed reset shows as `pending refresh`, never as full quota. Windows is
  reported as `unsupported` in this version.

One-time Claude relay setup (manual; the relay is not shipped with Ply releases
or Homebrew). From a Ply checkout, build it to a stable absolute path you choose:

```sh
go build -o /absolute/stable/path/claude-statusline-relay ./cmd/claude-statusline-relay
```

Then edit `statusLine` in `~/.claude/settings.json`. Keep any previous command
by passing it, POSIX single-quoted (each `'` written as `'\''`), to `--then`:

```json
{"statusLine": {"type": "command", "command": "'/absolute/stable/path/claude-statusline-relay' --then '~/.claude/statusline.sh'"}}
```

Without a previous command use only the quoted relay path. To undo, restore the
previous `statusLine` (or remove it) and delete the snapshot file. An already
configured compatible AI Usage relay with the same name is recognized as is.
`ply agent` inspects only `~/.claude/settings.json`; project-level status line
overrides and `CLAUDE_CONFIG_DIR` are not detected and can keep the snapshot
from updating.

## Install
```shell script
PLY_GO=/absolute/go1.26.10/bin/go make install
```

## Build
```shell script
PLY_GO=/absolute/go1.26.10/bin/go make build
```

Make's Go targets run through `scripts/with-go`, which requires `PLY_GO` to be
the absolute path of a Go binary reporting the exact go.mod `toolchain` under
`GOTOOLCHAIN=local`. It puts that binary first on `PATH`, so child `go` calls
use it too, and fails before building on a missing or different Go. Use
`PLY_GO=... scripts/with-go <command>` for direct `go test` or `go build` runs.
An explicit `GO=<path>` bypasses the selector and is a test seam only.

Local, reproducible build output lives under `.build/` in the worktree that
owns the Makefile, also when make runs with `-f` from another directory:

- `make build` writes `.build/bin/ply`. An explicit `PLY_BUILD_OUTPUT=<path>`
  keeps its meaning: relative paths resolve from the caller's working
  directory, and the binary is written exactly there.
- `make snapshot` runs GoReleaser from the repository root and writes its
  output to `.build/dist/`.

`make clean` removes only the physical `.build/` directory of that repository
root and succeeds when it is missing, so it can be repeated after an
interrupted run. It refuses a symlinked or non-directory `.build`, tracked
files under `.build/`, and a Makefile that is not at the root of the Ply Git
repository. It does not follow symlinks out of `.build/` and does not touch
explicit `PLY_BUILD_OUTPUT` paths elsewhere, source files, Go caches, installed
binaries or release receipts. Run `make clean` before Task closeout; closeout
still rejects ignored files that remain.

Older checkouts may still contain `./ply` or `dist/` from earlier builds;
those are not removed by `make clean`.

`make compat-cli`, and therefore also `make compatibility` and
`make preflight`, and the `test-workflow-handoff`, `test-workspace-work-items`
and `test-workspace-task-integration` targets require `PLY_ROUNDTRIP_ROOT=<dir>`
for the roundtrip scripts in `test/` (Task run, queue and removed source,
workflow handoff, work items, Task integration and Goal delivery), and so do
`make test-quality` and `make test`, whose `test/makefile_quality_test.sh` keeps
its fake tools, caches and quality output (which `make quality` requires outside
the repository) in one `ply-makefile-quality.XXXXXX` directory: an existing,
writable, absolute directory that is not a symlink and not inside a Git
repository or Ply workspace. There is no fallback to `TMPDIR` or `/tmp`. Each
run builds the CLI and runs its journey in its own `ply-*.XXXXXX` directory
under that root, with private HOME, fixtures, helpers and Go test and build
temp files there too. The scripts share this through the sourced
`test/roundtrip_root.sh`: `roundtrip_root <label> <dir-prefix>` exits 2 with
`<label>: unsuitable PLY_ROUNDTRIP_ROOT ...` before running anything, otherwise
creates `$run_root`, points `TMPDIR`/`GOTMPDIR` into it and installs the
cleanup trap, which restores write permission with `chmod -R` (not `find`)
before removing the directory. On HUP, INT or TERM the script waits for its foreground child and
exits 129, 130 or 143. The run reports its result and any failing command
output on stdout/stderr, and removes the directory on pass, failure or
interrupt. Shared Go caches are kept, and nothing else in the root is touched.
`test/workspace_task_roundtrip_root_test.sh` checks this with a fake toolchain.

`scripts/check-api-compat.sh`, `scripts/check-cli-compat.sh`, their
meta-tests `scripts/test-check-api-compat.sh` and
`scripts/test-check-cli-compat.sh`, and `scripts/test-search-replace.sh`
(`bash scripts/test-search-replace.sh`, no root variable needed) keep their work data under
`.build/compatibility/` of the physical worktree that holds the script, also
when run from another directory. They share the sourced
`scripts/compat_root.sh`: `compat_run_root <label> <dir-prefix>` exits 2 for a
symlinked or non-directory `.build` or `.build/compatibility`, otherwise
creates `$run_root` as one `<dir-prefix>.XXXXXX` directory there, points
`TMPDIR`/`GOTMPDIR` into it, exits 129, 130 or 143 on HUP, INT or TERM, and
removes that directory on pass, failure or interrupt; private Go build caches,
copied bases, fixtures and logs live in it as well. There is no fallback to
`TMPDIR` or `/tmp`, and other entries in `.build/compatibility/` are left
alone. Default reports are written to
`.build/compatibility/reports/api-report.json` and `cli-report.json`;
`.build/compatibility/reports` gets the same refusal; with
`API_COMPAT_REPORT_OUT` or `CLI_COMPAT_REPORT_OUT` set, the report is written
exactly there and that directory is neither checked nor created.
`test/cli_surface_contract_test.sh` runs `ply workspace init` and project
fixtures, which must live outside Git and Ply workspaces, so it uses
`roundtrip_root` like the roundtrip scripts: `make test-cli-surface` also
requires `PLY_ROUNDTRIP_ROOT`, and the binary, Go cache, HOME, logs and
fixtures of one run share one `ply-cli-surface.XXXXXX` directory there. A
failed assertion prints the first 80 lines of each file it checked before that
directory is removed. `scripts/compat_root_test.sh` holds `check_compat_root`,
used by both meta-tests to run copies of the checker in a fake worktree and
check the owned root, the tool temp environment, a kept neighbour, the original
exit on failure, with a failing `find` and on TERM, and the refusal of
symlinked roots and reports directories. `test/makefile_roots_test.sh` does the
same for the install, lint and preflight Makefile meta-tests and
`scripts/test-search-replace.sh`.

## Help
```shell script
ply
```

## Agent Slack notifications

`ply workflow notification` previews and sends an agent event (`agent_finished`,
`agent_stopped`, `feedback_required`) to one fixed Slack channel, as request v4–v6.
When separately authorized, create an Incoming Webhook for that channel, put its URL in a
named environment variable such as `PLY_SLACK_WEBHOOK` using your usual secret practice,
and create a route JSON containing the variable name and a private state root, without the
URL. No global Ply configuration or native workspace is needed.

See [the complete validated route/request example and check → apply → show procedure](docs/workflow-notification.md).
The event is a sender claim. Transport acknowledgement does not attest product
correctness, human acceptance, or that anyone read the message. Unknown outcomes cannot be
resent. Requests v1–v3 of the removed delivery process are refused; their stored records are
read as history and never sent again.

The repository provides [ply-agent-notify](skills/ply-agent-notify/SKILL.md), with an
[installation and diagnosis procedure](skills/ply-agent-notify/references/usage.md).
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

### Human input from an ordinary agent

`ply workflow input` lets an ordinary Claude or Codex agent publish one necessary
question (2-6 fixed choices or one exact short text) for a named activity, with an
optional existing Task, and read back one immutable answer. A configured app
discovers questions with `list`/`show` and records the answer with
`answer preview`/`apply`/`status`; the agent waits, acknowledges and reports its
use. Nothing starts or wakes an agent, and no Task is created. The versioned
contract is [schemas/workflow-human-input-v1](schemas/workflow-human-input-v1/README.md);
the agent procedure is the [ply-human-input skill](skills/ply-human-input/SKILL.md).
`make install-human-input-skill` installs it to `~/.agents/skills/ply-human-input`
(identical install is a no-op; a different skill is replaced only with
`HUMAN_INPUT_SKILL_BACKUP=/new/path`). `make test-human-input-skill` checks the
installer in an isolated `/private/tmp` root. A copied directory is not evidence that an agent
can use the skill: confirm that the Claude and Codex entry points actually
discover it, or that a readable mandate pointer such as the existing
`~/.agents/AGENTS.md` names it. That check is part of installation and needs no
broad provider configuration change.

Workspace Project/repository/Epic definitions can live outside disposable runtime data; see [durable workspace config](docs/durable-workspace-config.md).
