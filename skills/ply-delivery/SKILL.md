---
name: ply-delivery
description: Begin an assigned Ply Task using actual runtime observations, then verify and deliver its exact candidate under the frozen agreement. Also register, inspect or resume native Delivery.
---

# Ply Delivery

Use Ply's native Delivery commands from the registered Task worktree. The CLI
owns Git, PR, evidence, locking and recovery; this skill supplies the workflow.

## Begin an assigned Task

Read the frozen goal and mandate from the start message. Reuse their human
authorization. The bound callback guide selects the startup protocol and gives
the immutable control, context, provider and native session.

For `mandate_receipt_v1`, run that control's `workflow execute check-in RUN_ID
--context PATH --provider PROVIDER --native-session SESSION` before Task writes.
Ply checks the frozen mandate, exact owner and cwd and keeps an idempotent receipt.
Optional model, profile, policy and roots are actual sourced observations; omit
unknown facts. Reported roots do not establish sandbox coverage. The receipt
acknowledges the mandate; runtime evaluates actions when they happen. No blanket
permission confirmation, handwritten acceptance or repeated human yes is needed.
A known refusal or unresolved prior effect still stops dependent work.

For a frozen legacy contract, follow its `workflow execute begin --help` and
original actual-permission requirements. Do not move it to check-in, rewrite its
old claims or revive a stopped owner. Both protocols preserve identity and
scope, and reuse identical input after a lost response. Preserve provider controls.

Use the existing authorized notification context. For a new binding, supply
`--notification-route`, `--notification-authority` and `--notification-timezone`
to the contract-selected startup command; Ply returns the private `delivery/notification-task.json` path for
the installed ply-agent-notify helper. This creates no transport or global
configuration. If the goal gives notifications to a parent, omit these flags
and return to that parent without a separate notification.

If a provider rejects acceptance or report transport, stop the rejected action.
The coordinator can use `workflow execute observe-stop --help` with the actual
denial evidence. This preserves a sourced observation without owner acceptance
or writer release. After the cause is clarified, the same owner uses
`workflow execute resolve-stop --help` with the exact stop digest, the actual
clarified cause and the changed action it takes now; no permission confirmation
is needed. Then use the contract-selected startup command if receipt is still absent. The
resolution starts nothing and resends no rejected input. Never repeat the refused
action or change global permission rules to forget a denial. For a selected
replacement, the coordinator uses `workflow execute retire --help`; the retired
owner sends nothing further.

## Independent review rounds

Review is a round on the same Run, not a separate review Task. The coordinator
uses `workflow execute review --help` for the owner's clean committed candidate;
Ply starts the separate reviewer in its own Herdr tab and binds its session.
The reviewer follows the generated brief and returns from that pane with
`workflow execute review-return --help`. A negative review is a completed
review: the owner corrects, commits and asks for a new round. `verify` without
`--review` uses the round that passed exactly the current candidate.

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
Use the structured goal `delivery` field described by the native command help;
execute-preview binds the
actual Task source branch. A started run retains its frozen agreement. A later
goal revision cannot change that run's target or permissions.

## Existing Project membership

When the user authorizes adding a repository to an existing shared Project, use
`ply workspace project add-member PROJECT_ID --repo REPO_ID=PATH` from its containing
workspace. Verify the stored locator and Git common directory with
`ply workspace project show PROJECT_ID --format json` and the membership count with
`ply workspace project list --format json`. The operation appends one explicit member,
accepts dirty worktrees outside the wrapper, and requires no execution Epic. Matching
Git identity retries preserve the existing locator; identity conflicts fail without
partial membership. Preserve the chosen Project and its execution/history bindings.
Membership registration is separate from Task execution and delivery; it does not
extend an existing Epic to another repository. Use this only with membership authority,
not merely because an implementation Task mentions a future consumer registration.

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

Follow the frozen acceptance policy. An explicitly automatic agreement uses
the exact candidate's executed acceptance and independent review. A human gate
requires the human's actual answer through the preserved QA callback; technical
success and silence cannot fill it. A later human fail or blocked answer still
vetoes automatic delivery. Keep responsibility for correction and requalification.

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
For automatic local Epic delivery, `execute integrate` from the Task worktree
normally ends with exit 5 and `owner_handoff`: integration, queue and Epic base are
observed and the source cleanup waits for you to leave. That is delivered, not a
stop. Do not retry or report a stop; end the session and return to the
coordinator, whose `workflow execute follow` removes the worktree.
Preserve partial facts and the CLI's next action when blocked, failed or unknown.
Waiting is durable state; it starts no background worker. Distinguish technical
qualification, actual human judgment and final delivery. Task completion never
declares the whole Epic complete or authorizes its later publication.

Notification subscriptions and Slack delivery are separate from this workflow.
