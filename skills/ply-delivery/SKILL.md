---
name: ply-delivery
description: Register, inspect, execute or resume a native Ply Delivery for an exact verified Task candidate under its frozen PR or local integration agreement.
---

# Ply Delivery

Use Ply's native Delivery commands from the registered Task worktree. The CLI
owns Git, PR, evidence, locking and recovery; this skill supplies the workflow.

## Find the exact contract and candidate

Read `ply workspace task list --format json`, then
`ply workspace task show <task-id> --format json` and the bound
`ply workflow execute show <workflow-run-id> --format json`. Match the Task
worktree, result ID, commit/tree, exact Spec and frozen delivery agreement.
Read `ply workflow delivery list --task <task-id> --format json` for an existing
registration or partial attempt before creating another.

For a goal that has not started, propose the relevant mode and exact target:

- `pull_request`: publish the Task source and open/update its PR in the chosen
  GitHub repository and base. Completion stops before merge.
- `local_epic_integration`: return to the registered development Epic checkout.
- `local_branch_integration`: return to an explicitly chosen local product ref
  and checkout, including main/master. It works without a remote or GitHub tools.

Reuse an explicit human choice already given. Otherwise obtain the missing mode
or target choice before publishing that agreement. Repository visibility,
missing remote, cwd and branch names provide context, not integration authority.
Use the structured goal `delivery` field documented in
[the product guide](../../docs/workflow-delivery.md); execute-preview binds the
actual Task source branch. A started run retains its frozen agreement. A later
goal revision cannot change that run's target or permissions.

## Register and check

After meaningful verification and review qualify the exact candidate, write a
new private JSON file using the real IDs from readback:

```json
{
  "kind": "ply.delivery.registration",
  "schema_version": 1,
  "publication_key": "task/delivery/candidate-1",
  "task_id": "example-task",
  "task_result_id": "trs_<actual-result-id>",
  "workflow_run_id": "wfr_<actual-run-id>"
}
```

```shell
ply workflow delivery register --file /absolute/private/delivery.json --format json
ply workflow delivery check dlv_<returned-id> --format json
```

The native records supply the Spec, evidence, mode, target and actual authority.
Registration is allowed before human QA; it performs no delivery effect. Reuse
the same publication key and bytes after a lost response. Changed content needs
a new key; retain predecessor history when replacing a candidate.

All modes require the human's actual pass for the exact candidate. Use the
run's preserved QA callback and guide to record the exact answer; technical
success, silence and your own judgment cannot fill that gate. A fail keeps the
owner responsible for correction and qualification of a new candidate.

## Execute or resume

```shell
ply workflow delivery execute dlv_<returned-id> --format json
ply workflow delivery show dlv_<returned-id> --format json
```

Execute only within the already authorized delivery. It rechecks the contract,
candidate, evidence, runtime authority and human gate. Preserve the same ID after
an interruption and use execute to observe and resume the attempt. An unknown
outcome requires effect observation; do not replace the ID, push manually or
create another PR to bypass it.

For PR delivery, Ply reads local `~/.agents/local/pr-preferences.yaml` on the first
attempt and preserves the effective assignee/reviewer choices. Explicit choices
override repository preferences, which override global preferences. Keep personal
defaults out of versioned goals and examples. A metadata failure can leave a real
PR: preserve its URL and use `metadata --help` for an explicit correction, then
execute the same Delivery. Empty defaults do not remove existing people.

## Return the observed result

Report `delivered` only with the durable receipt: PR URL and observed exact head,
or local before/after refs plus registered base and current Task queue closure.
Preserve partial facts and the CLI's next action when blocked, failed or unknown.
Waiting is durable state; it starts no background worker. Distinguish technical
qualification, actual human judgment and final delivery. Task completion never
declares the whole Epic complete or authorizes its later publication.

Notification subscriptions and Slack delivery are separate from this workflow.
