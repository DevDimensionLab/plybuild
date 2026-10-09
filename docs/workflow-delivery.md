# Deliver an exact Task candidate

`ply workflow delivery` records and executes a durable Delivery under the Task's
frozen agreement. The package refers to the native TaskResult, exact Spec,
candidate commit/tree, verification, review and run authority. Human QA remains
a separate record bound to that candidate.

For explicit human ownership, schema 2 agreements add
`"integration_owner": "human"`. The developer stops after qualification and
handover (or exact PR publication). The separate [human integration workflow](integration.md)
confirms integration, optional installation and safe Task closeout. Schema 1
agreements retain the behavior described below.

| Mode | Exact target | Completion |
| --- | --- | --- |
| `pull_request` | GitHub repository and base ref | Source branch and PR observed with the exact candidate; metadata applied; current Task closed before merge |
| `local_epic_integration` | Registered development Epic ref and worktree | Local integration observed; registered base updated; current Task queue entry closed |
| `local_branch_integration` | Explicitly selected local ref and worktree | Same local evidence and closure, including main/master without any remote |

All three modes require meaningful tests, review, actual runtime authority and a
human pass for the exact candidate in historical v1/v2 agreements. New v3 local Epic
agreements may select automatic acceptance below. Registering a goal, queue entry or Delivery
does not supply these gates.

## Explicit automatic acceptance for a new local Epic task

Select agreement schema 3 when registering the Goal, before starting execution:

```json
{
  "schema_version": 3,
  "mode": "local_epic_integration",
  "project_id": "product",
  "repo_id": "product",
  "epic_id": "development",
  "target_ref": "refs/heads/development",
  "target_worktree": "/absolute/workspace/product/development",
  "acceptance": {
    "schema_version": 1,
    "mode": "automatic",
    "responsible_actor": "the selected orchestrator",
    "require_human_qa": false,
    "timeout_seconds": 300
  }
}
```

The bound source ref is added during Task preparation. This choice travels through
Goal, execution Spec, native mandate and Delivery; it cannot change an already
started run. PRs and local main/master delivery cannot select this policy. V1/v2
and absent policies retain their historical human gates. `require_human_qa: true`
adds a real human judgment after automatic verification. Optional
`integration_owner: human` preserves the separate human-started integration path.

From the accepted owner's Task checkout, use its preserved control executable:

```shell
/absolute/ply-control workflow execute verify wfr_<id> \
  --context /absolute/context.json --review /absolute/candidate-review.json \
  --candidate-binary /absolute/private/built-candidate
/absolute/ply-control workflow execute integrate wfr_<id> \
  --context /absolute/context.json
```

`verify` executes the declared acceptance script against the clean committed
candidate and applies the actual review. Automatic `PlyDeliveryVerification@2`
records preserve the contract argv/cwd, actual `automatic.executed_argv`, script
and executable snapshots, instruction digest, candidate OID/tree, responsible
actor, actual executor/session, environment, timestamps, exit, output and
`pass`, `fail` or `blocked`. `PLY_CANDIDATE_BINARY` names the preserved executable
snapshot; `PLY_CANDIDATE_OID` and `PLY_CANDIDATE_TREE` name its candidate. Tests
receive the recorded PATH, `LC_ALL=C`, and isolated HOME/TMPDIR. Provider secrets
and the rest of the caller's environment are not inherited. Acceptance scripts
must provision their local fixture within that declared environment. Candidate
executables are limited to 64 MiB to fit native evidence retention; larger inputs
are blocked before test execution.

A passing invocation with valid review satisfies automatic acceptance without
opening QA.txt, asking a human to start the script, or inventing a HumanQA record.
Later actual human `fail` or `blocked` still vetoes the same candidate bytes;
reverification alone cannot erase the answer. Changed candidates, scripts,
reviews or executables need fresh acceptance before any new integration effect.

A completed test result is preserved before candidate registration. Repeating
an interrupted invocation recovers that receipt exactly once. `verify --reuse
verify-<number>` explicitly requalifies a recorded successful attempt without
running it again. Unknown or still-active attempts block competing execution;
this version never automatically repeats an unknown test. Nonzero exit, timeout,
start failure, missing inputs and interruption cannot qualify a candidate.

Automatic local delivery reuses native Git, queue and Epic-base recovery. A
confirmed no-effect integration may be retried through the same Delivery; an
observed merge is reconciled without another merge. The owner is released only
after the exact integration is observed, then native Task closeout preserves the
worktree and branch (`keep=true`) and marks the lifecycle complete. Closeout
failure stays visible independently of successful integration; resuming the
same delivery completes closeout only. Removing retained resources remains a
separate explicit native cleanup action.

Text status and versioned JSON show acceptance separately from human judgment
and final delivery. `ply capabilities` advertises the new CLI through additive
`workflow_extensions`; historical read-catalog fields are unchanged.

## Record the agreement before starting

Add an optional `delivery` object to a Goal Spec@2. Keep the ordinary goal fields
and document-preservation rules in [the execution guide](workflow-execute.md).
For example, a PR agreement is:

```json
{
  "schema_version": 1,
  "mode": "pull_request",
  "project_id": "project",
  "repo_id": "product",
  "epic_id": "product-fixes",
  "target_ref": "refs/heads/main",
  "github_repository": "example/product",
  "remote": "origin"
}
```

A local branch agreement is:

```json
{
  "schema_version": 1,
  "mode": "local_branch_integration",
  "project_id": "project",
  "repo_id": "product",
  "epic_id": "product-fixes",
  "target_ref": "refs/heads/main",
  "target_worktree": "/absolute/workspace/product/main"
}
```

Use `local_epic_integration` with the exact development Epic ref and checkout for
an Epic return. That mode rejects main/master. In both local modes the selected
checkout is the Epic's registered return context; no extra development branch is
required for direct local product delivery. The actual ref matters, even when
the directory has a different name. Local contracts omit GitHub and remote fields.

Ply binds `source_ref` when execute creates the Task worktree. Preview and the
immutable execution Spec preserve the exact agreement. A pending goal revision
may select a different agreement for a future start; it cannot change a started
run. Historical goals without `delivery` retain their existing local behavior
and never acquire PR or main/master authority.

```shell
ply workspace task spec record example-task --file /absolute/goal.json --format json
ply workspace task spec show example-task --spec example-goal --revision 1 --format json
ply workspace task queue set --file /absolute/queue.json --format json
ply workflow execute --check --format json
```

The user chooses or confirms the proposed mode and target. An explicit choice
already supplied can be reused. Neither a private repository nor an absent remote
implicitly authorizes local integration. Actual startup and runtime acceptance
must cover the selected effects as well as the immutable agreement digest.

## Register the qualified candidate

Find the actual TaskResult and run using Task and workflow readback. Write a
private registration file (mode `0600`) with those exact IDs:

```json
{
  "kind": "ply.delivery.registration",
  "schema_version": 1,
  "publication_key": "example-task/delivery/candidate-1",
  "task_id": "example-task",
  "task_result_id": "trs_<actual-result-id>",
  "workflow_run_id": "wfr_<actual-run-id>"
}
```

```shell
ply workflow delivery register --file /absolute/private/delivery.json --format json
ply workflow delivery check dlv_<returned-id> --format json
ply workflow delivery show dlv_<returned-id>
ply workflow delivery list --task example-task --format json
ply workflow delivery list --epic product-fixes --format json
```

Registration validates native records and preserved bytes and causes no Git or PR
effect. Repeating the same publication key and content returns the same ID.
Different content under that key is a conflict. A new candidate uses a new key
and may name its preceding Delivery in `predecessor_id`; it needs its own
acceptance evidence under the frozen policy. Optional `title` and `body` provide PR text; optional `metadata` supplies
explicit people choices for this Delivery.

Check, show and list are read-only. Versioned JSON distinguishes registration,
waiting for a human answer, readiness, delivery in progress, completion,
blockage, failure and an unknown effect. Receipts include the frozen contract,
candidate and workflow links, actual gate observations, attempts, timestamps,
history and next action. Registration does not promise a background execution.

## Execute and recover

After the exact candidate satisfies its frozen automatic or human acceptance gate:

```shell
ply workflow delivery execute dlv_<returned-id> --format json
ply workflow delivery show dlv_<returned-id> --format json
```

Execution rechecks candidate/evidence integrity, the frozen mode and target,
native runtime acceptance and candidate-bound acceptance. Local execution also
requires the unchanged, clean registered parent. A moved parent or changed
candidate blocks delivery; Ply does not silently rebase or reuse old QA.

PR execution verifies the selected remote repository, source branch, candidate
head and base. It publishes only the source branch, reuses an exact existing PR,
or creates one and reads it back. Wrong head/target, ambiguous matches or a
required force-push block the effect. Ordinary advancement of the selected PR
base does not change its identity. PR delivery never integrates locally, updates
the registered Epic base, pushes the target branch, merges a PR or enables
auto-merge.

Local branch delivery uses no GitHub executable, credential, preference file or
network. When a repository has remotes, they remain untouched. Local delivery is
complete only after the observed integration, registered base and current Task
queue closure. Other Tasks and queue entries retain their history.

Repeat execute with the same ID after an interrupted command. The durable attempt
records effect intent; recovery observes an existing PR or local integration
before completing remaining metadata. An unavailable observation leaves the
effect unknown with a next action. Concurrent attempts, including separate
Delivery IDs for the same candidate and target, share the effect and its receipt.
Completed receipts retain delivery-time observations; later PR changes do not
rewrite that history.

## PR people preferences

On the first PR attempt Ply reads `~/.agents/local/pr-preferences.yaml` with
`schema_version: 1`. Optional `defaults` and case-insensitive `repositories`
entries contain independent `assignees` and `reviewers` fields. For each role,
explicit registration metadata overrides its repository value, which overrides
the global value. Missing files and fields invent no people; an empty list clears
the inherited default. Keep actual personal defaults in the local file.

`execute --preferences-file /absolute/private/preferences.yaml` selects an
explicit local file for the first attempt, including isolated product exercises.
The same validation and precedence rules apply; retries keep the preserved choice.

The effective choices and their sources are frozen in Delivery history. Retrying
does not pick new recipients because the preference file changed. Existing PR
people remain unless explicitly removed. Invalid configuration, unavailable
accounts, missing permission and rejected metadata remain visible errors. An
already observed PR survives a metadata failure and is reused on retry.

To revise the preserved choice, supply a private JSON object with optional
`assignees`, `reviewers`, `remove_assignees` and `remove_reviewers` string arrays.
For example, `{"assignees":[],"reviewers":[]}` clears desired additions while
preserving people already on the PR; only explicit removal arrays remove people.

```shell
ply workflow delivery metadata dlv_<returned-id> --file /absolute/private/metadata.json
ply workflow delivery execute dlv_<returned-id>
```

The receipt reports whether metadata was actually applied; a real PR URL alone is
insufficient for completed delivery.

## Events and agent use

Events have stable IDs, candidate/attempt correlation, actor/origin and UTC times.
An observed PR creation records `delivery.pr_created` with its URL; reusing an
older PR records reuse. Recovery retains the original logical event. These events
can support future consumers. This command starts no Slack sender, subscription,
notification queue or daemon.

The repository's [ply-delivery skill](../skills/ply-delivery/SKILL.md) describes
the agent journey: find the exact native Task, register, check, execute/resume and
report its observed receipt. Technical qualification, human judgment and final
delivery are distinct. Finishing one Task does not finish or publish its Epic.
