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
  bitbucket   Bitbucket functionality
  clean       Clean files and folder in a project
  completion  Generate the autocompletion script for the specified shell
  diagrams    Various tools for generating diagrams
  doc         Opens documentation in default browser
  examples    Examples found in cloud-config
  format      Format functionality for a project
  generate    Initializes a maven project with ply files and formatting
  git         Git commands
  help        Help about any command
  info        Prints info on spring version, dependencies etc
  init        Initializes a maven project with ply files and formatting
  install     Various install options for generating autocompletion etc
  lint        Linting commands
  maven       Run maven (mvn) commands
  merge       Merge functionalities for files to a project
  profiles    Manage profiles settings for ply
  query       Query dependencies in a project
  status      Status functionality for a project
  tips        Use tips to learn information faster
  upgrade     Upgrade options
  workflow    Manage agent workflows
  workspace   Manage Ply workspaces

Flags:
      --debug   turn on debug output
  -h, --help    help for ply
      --json    turn on json output logging

Additional help topics:
  ply about      About ply

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
`ply workspace project list`.

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

Sub-tasks, additional repository anchors on an Epic, editing, rebinding, refreshing, cleanup, and
WorkflowRun creation are outside this first workspace work-item version.

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
