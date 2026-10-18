# Deliver a queued goal

For one explicitly automatic local Epic Goal on macOS, use the
[durable headless worker](workflow-worker.md) with `--headless`. It has its own
process ownership and control commands. The interactive workflow follows below.

`ply workflow execute` gives one interactive implementor responsibility for a
feature from its goal through its frozen acceptance policy and local return.
The planner records the desired result and design; the implementor chooses the
detailed solution, tests, review and necessary corrections.

Run it from a Herdr terminal in the registered Epic worktree that should receive
the feature. A normal terminal can use `--check`, but starting requires the actual
Herdr caller context (`HERDR_ENV=1`); `--herdr-workspace` alone does not provide it:

```shell
ply workflow execute --check
ply workflow execute
# Or choose a unique pending goal Spec ID in this Epic:
ply workflow execute --spec explain-start-errors
```

`--check` reads the exact goal, current base and local runtime. For a new Claude
launch it also shows the planned persistent folder trust effect, including the
effective configuration file and exact project key. It creates no worktree,
launch file or agent and writes no provider configuration. Without it, the command
records the caller's exact choice, creates one sibling feature worktree, preserves
a control executable and starts the assigned Claude or Codex in Herdr. The caller must be inside the return
worktree, including when explicit `--project`, `--repo`, `--epic` flags are supplied.
Supply all three target flags or none.

Several Epics may share one clean return worktree and branch, such as a project's
`main` checkout. Adoption reserves no directory; each Task still gets its own
sibling branch and worktree, and integration stays bound to the exact candidate
and the freshly observed target. In a shared directory, Ply selects an Epic whose
recorded branch is checked out there, so a stale historical registration never wins. When several
current Epics remain, `--spec <id>` selects the one whose queue holds that goal. Otherwise
the command refuses and prints the exact `--project <id> --repo <id> --epic <id>`
choices. After another Epic integrates into the shared branch, this Epic's base is
stale until `ply workspace epic base` records the new commit; no merge, reset or
rerun happens implicitly.

The selected Epic is retained through Task preparation and retries, including
after the goal has been materialized into an execution Spec. Preparation rechecks
that exact target and base before creating the isolated Task copy; repeating the
same request recovers its history rather than selecting another shared Epic.

Each execution works for exactly one selected Epic in one repository. A planner may
continue recording goals and changing the pending queue while the owner works.
The current preparation is retained. Returning a second concurrent feature to a
parent that has moved needs a later refresh workflow; this version reports the
changed parent and preserves the work instead of silently merging it.

## What the planner records

Use the existing Task, Spec and queue commands:

```shell
ply workspace task spec record explain-errors --file /absolute/goal.json --format json
ply workspace task queue set --file /absolute/goal-queue.json --format json
ply workspace task queue list --ready --format json
```

The strict **WorkspaceTaskSpecDraft@2** goal variant uses `schema_version: 2` and
`contract_kind: "goal"`. It retains the existing publication key, Task and Spec
IDs, expected previous revision, exact Problem reference, recorder and document
snapshot fields. Its product fields are:

| Field | Meaning |
| --- | --- |
| `title`, `objective` | The user-visible outcome and why it matters |
| `documents`, `design` | Preserved design documents and `{document_id, section}` references |
| `requirements` | Sorted unique `{id, acceptance}` entries describing observable results |
| `constraints` | Sorted unique outcome and effect boundaries |
| `executor` | Explicit nullable `provider`, `model`, `effort` choices |
| `delivery` (optional) | Versioned PR, local Epic, or explicit local branch agreement; see [Delivery](workflow-delivery.md) |
| `change_reason` | Why this goal revision was published |

A goal has at least one design reference and requirement. It has no
`implementation_basis`, detailed technical parts, completed-test claims, or human
solution selection. The document store preserves actual bytes and digests using
the existing file/snapshot source rules. Revised goals get new immutable revisions.

**WorkspaceTaskQueueDraft@2** retains `publication_key`, target IDs,
`expected_revision`, `entries` and `registry_upgrade`. It adds `recorder` instead
of a fabricated `human_decision`. Each entry names exactly one Task and either:

```json
{"task_id":"explain-errors","goal":{"spec_id":"explain-start-errors","spec":{"revision":1,"manifest_sha256":"sha256:<actual manifest digest>"}}}
```

or an existing `selection` reference. Goal Spec IDs used with `--spec` must be
unique among pending goals in this target. Ambiguity, a foreign or missing goal,
an unknown provider and a missing executable are errors. There is no provider
fallback. The next selector chooses the first eligible pending goal; inspect
the queue for reasons why earlier entries are blocked.

At execution, Ply derives an **execution Spec@2** with the unchanged goal origin,
current parent/base, selected implementor and a real acceptance entrypoint. The
readiness assessment means the start contract is complete; it does not claim
that implementation, tests or product QA have already happened. The implementor
writes meaningful tests and the preserved `acceptance.sh` script. A missing script
cannot pass verification.

New structured delivery agreements travel through the goal, queue, preview and
frozen execution inputs. They bind the actual Task source, target and permitted
effects. `ply workflow delivery` registers and performs the resulting qualified
candidate after its exact automatic or human acceptance gate. New agreement v3
can explicitly select [automatic local Epic acceptance](workflow-delivery.md#explicit-automatic-acceptance-for-a-new-local-epic-task).
An existing schema 1 local Epic run can record a human's explicit
[acceptance policy choice](workflow-delivery.md#select-automatic-acceptance-for-an-existing-local-epic-run)
through `workflow execute acceptance RUN_ID --context PATH --file PATH`,
preserving the original frozen agreement and prior evidence. The choice is
separate from both test execution and human QA. Historical goals without this
field or an explicit supported amendment keep their original local return
contract.

## Runtime and permissions

When the goal omits its provider, the compatibility default is Codex. Its model,
effort and existing permission profile come from the local Codex settings. An
explicit model or effort in the goal wins. Claude defaults to `opus` with
`medium` effort. New Claude executions request its native auto permission mode
with the documented [`--permission-mode auto`](https://code.claude.com/docs/en/cli-reference)
argument. Omit `--permission-profile` for Claude, or specify `auto`; other values
are rejected. Model and effort remain separate from the provider identity checked
against the actual Herdr session.

The selected binaries are resolved and hash-bound before a worktree is created.
Herdr uses the current workspace, or `--herdr-workspace`. Codex uses its selected
managed permission profile with on-request automatic approval review; its
native project-trust and onboarding remain under Codex control.
Claude receives its own supported model, effort and permission arguments.
No Codex trust/config flags are forwarded to Claude.

New Claude launches preserve an explicit folder trust choice with the exact
physical Task worktree. Before its single provider start, Ply attempts to set
`projects[worktree_path].hasTrustDialogAccepted` to `true` in Claude's effective
configuration. The project key is the Task worktree's physical path; Ply does not
add a trust entry for the main checkout, parent Epic or another worktree. The
configuration path follows the actual launch environment, including supported
`CLAUDE_CONFIG_DIR` behavior. `--check` reports that path instead of assuming the
default personal configuration.

The update preserves unrelated configuration and creates a private backup before
changing an existing file. Already accepted trust is a file no-op. A missing
configuration file keeps normal onboarding. A failed or skipped trust update is
reported and the same ordinary Claude startup continues; its own trust dialog may
still need an answer. Ply does not repair malformed configuration or override
separate settings, MCP, authentication or tool-permission prompts.

Claude's [auto mode](https://code.claude.com/docs/en/permission-modes) reviews actions
with a classifier instead of asking about routine commands. Explicit ask rules,
denials, hooks and sandbox enforcement still apply. Auto mode can be unavailable
or fall back to prompts under the provider's policy. Handle remaining native
prompts in the same tab; folder trust does not bypass tool approvals. Preserved
launches retain their original trust choice and requested mode, including earlier
manual-mode requests. Older intents and requests without a trust choice never
gain one when read or resumed. Repeating execute or resuming startup does not
repeat the configuration update or switch an existing session's permission mode.

The launch-policy snapshot describes requested settings. New interactive Runs
bind `startup_protocol: "mandate_receipt_v1"` in both the request and native
mandate. The human mandate supplies scope and the exact delivery target;
runtime evaluates each action when it happens. Unknown future decisions do not
block all work. Known refusals and unknown previous effects remain binding.
Headless workers and frozen historical Runs retain their existing contracts.

### Receive the mandate with check-in

From the exact Task and bound immutable control, use the command in its callback
guide:

```shell
/absolute/ply-control workflow execute check-in RUN_ID --context /absolute/context.json \
  --provider claude --native-session ACTUAL_SESSION
```

The three selectors are the only required receipt arguments. Ply computes the
mandate and input digests, observes the live owner, checks cwd and records the
same immutable receipt on identical input. It does not ask for another human
approval, a handwritten acceptance file or blanket permission confirmation.
Optional `--model`, `--profile`, `--policy`, `--read-root` and `--write-root`
preserve actual sourced recipient observations. Omitted facts stay unknown.
Reported roots have `matches_contract: null`; profile names and root strings
never become sandbox attestations or future action approval.

Readback distinguishes the launch request, owner receipt, recipient observations,
unknown sandbox coverage, preserved human scope, per-action runtime decisions,
actual target changes, technical acceptance and final delivery. Native mutation
checks still bind the exact Run/owner/session, mandate, effects, candidate and
delivery target. A mismatched entrypoint cannot downgrade a frozen contract.

### Frozen legacy startup

A Run without the new protocol keeps its original callback guide and `workflow
execute begin` requirements, including its actual policy evidence and permission
confirmation. Its receipt and any stops remain unchanged. The new check-in
cannot convert old permission claims or clear either stopped CM01 Run.

When the owner has an existing notification authorization, add
`--notification-route PATH --notification-authority SOURCE --notification-timezone ZONE`.
Ply creates `delivery/notification-task.json` for the installed notification helper,
using the actual provider, Task cwd, native activity and workflow Run identity.
This binds an existing route; it performs no transport or global setup. Internal
delegates whose parent owns notifications omit these flags. Preserve this binding
and event identity across continuations.

### Preserve a provider stop without an owner callback

A coordinator in the containing workspace can preserve an actual denial with:

```shell
ply workflow execute observe-stop RUN_ID --event-id report-denied-01 \
  --provider claude --native-session ACTUAL_SESSION --stage report \
  --observer 'Actual coordinator' --reason 'Observed provider rejection' \
  --evidence /absolute/private/denial.txt
```

Stages are `startup`, `acceptance` and `report`. Omit the native session only for
startup before session binding. The command preserves a copy and hash of the
actual evidence. It records coordinator observation separately from recipient
acceptance and retains the writer reservation and any unresolved attempt.
Identical retries reuse the event; changed evidence under that ID is rejected.
A later owner report cannot hide the active stop. Readback separately shows the
terminal, Task prompt, observed owner receipt, runtime permission and target
changes. Target changes alone do not claim meaningful tests or product completion.

After the cause is actually repaired, the same owner uses the bound control:

```shell
/absolute/ply-control workflow execute resolve-stop RUN_ID --context /absolute/context.json \
  --stop sha256:ACTUAL_STOP_DIGEST --resolution 'Actual cause repair and current permission source' \
  --permission-confirmed
```

Resolution requires fresh identity checks and refuses an unknown effect attempt.
It preserves the denial, starts nothing and sends no input. If acceptance is
absent, the owner uses the startup command selected by the frozen callback guide before Task writes.
Neither operation releases a writer or authorizes a replacement Run. These owner
commands must be present in the bound control; do not overwrite an older control
to acquire them. Existing accepted Runs can use the supported compatible-control
continuation when needed. Historical denied Runs remain preserved.

Native folder trust remains separate. If its existing configuration lock cannot
be acquired, the same Claude tab may require its native trust decision before
`execute resume`. Preserve that action as a startup limitation; do not retry the
launch or change permission mode. FL01 does not repair that configuration lock.

## Follow the owner

```shell
ply workflow execute show wfr_<digest>
ply workflow execute follow wfr_<digest> --timeout 60
```

Show is read-only. Follow observes the same provider/session. A timeout ends
observation and leaves the provider running. Repeating execute reuses the
preserved launch; it never creates a second tab to conceal an uncertain start.
The interactive Herdr session is also where the human can steer the owner.

If a native trust dialog outlasts the initial startup window, resolve it in the
preserved tab, then continue that same startup:

```shell
ply workflow execute resume wfr_<digest> --timeout 60
```

Resume checks the existing provider, pane and terminal. It can complete the
readiness exchange and send the first Task prompt only when no Task prompt has
been attempted. It grants no trust or permissions, starts no additional process
and refuses to replay an uncertain prompt. `show` and `follow` remain observers.

Owner callbacks use the **absolute preserved control executable**, exact Task
cwd and private context printed in the start guide. Installing another Ply
candidate does not replace that control binary. Callback surfaces are `report`,
`acceptance`, `verify`, `qa` and `integrate` under `workflow execute`; runtime
acceptance uses `workflow run accept`. The `acceptance` callback records the
explicit policy choice for an eligible existing run. Consult each command's
help for its required evidence.

Reports preserve working, failed, unknown and not-run facts. Candidate
qualification executes the declared script and checks actual review evidence and
the clean exact candidate. It produces technical evidence and a TaskResult while
the same owner remains responsible for product QA and local integration.

Use the supplied review template and bind its exact candidate OID and tree.
`reviewer_claim` is a short actual identity claim (1–256 Unicode code points).
Finding and fix `evidence_ids` refer to managed candidate artifacts, such as
`candidate-review` for the review record itself or `verifier-stdout` and
`verifier-stderr` for actual command output. They are not arbitrary note labels.
The callback guide lists the available artifacts. Invalid review references are
rejected before a candidate handoff is created. Preserve rejected inputs; a
corrected review uses a new private file and a new verification after the prior
attempt has a known recorded outcome. Never replay an uncertain attempt.

After a later candidate qualifies, resolved qualification failures remain in
the saved event history; the current readback describes the qualified candidate.
Current drift or a newer failed attempt still remains visible.

When the agreement requires human QA, the owner presents one installed product
journey and preserves the actual human
answer as `DeliveryHumanAttestation@1`. That record binds `task_id`,
`task_result_id`, `result_oid`, `result_tree`, `outcome`, `actor_claim`,
`start_surface`, UTC start/completion timestamps, `answer` and `observation`.
It is a local human attestation, not cryptographic identity proof. Silence,
technical success or agent approval never supplies the answer. A fail permits
correction and a new immutable candidate generation in the same session.

Under historical human agreements, only a matching human pass permits the owner
to integrate the exact candidate into the
unchanged clean parent, advance that preparation and update the Epic base. This
does not push, merge master, deploy, close tabs or delete worktrees. A candidate
TaskResult alone does not release the delivery reservation.

An optional `--notification-context` binds an existing notification context to
the owner. Otherwise the owner follows applicable standing notification
preferences through the installed skill. Internal delegates return to the owner;
they do not duplicate completion notices. Slack delivery is separate from human
QA, and Slack thread replies are not an input channel in this version.

## Compatibility and recovery

Spec@1, QueueDraft@1, Handoff@1/@2 and workflow-run@1 retain their old meaning.
New goal/execution manifests, delivery Handoff@3 and workflow-run@2 are explicit
versions. Older binaries cannot silently reinterpret a new contract.

Requests, goal origins, launch intent, reports and candidate generations remain
preserved after failures. Follow or inspect the existing execution first. Missing
evidence, an unknown provider session, changed candidate, dirty/moved parent or
uncertain merge effect cannot be treated as success or authorization to replay.

### Restart an interrupted Codex startup

Use the same Spec selection when Codex exits during startup or its Task tab has
been closed before Task input:

```shell
ply workflow execute --spec explain-errors --restart
```

Ply finds the preserved attempt; no run ID is needed. Run from the registered
return worktree in a Herdr terminal. `--restart --check` previews the next attempt
without writing or starting anything. The ordinary command without `--restart`
continues to inspect an existing attempt. If this Spec has no attempt yet,
`--restart` performs the ordinary first start. The global `--force` flag retains
its existing prompt-default behavior and does not reset an execution.

Apply repeats the checks, preserves the installed Ply control executable, and
reserves a new startup generation. The run, Task, preparation, feature worktree
and queue ownership stay the same. Codex is resolved from the current `PATH`;
model, effort, permission profile, trust and goal retain their original choices.
A provider update is supported, but changing the provider binary is not required.

When the previous terminal still exists, fresh Herdr and process observations
must show an unclaimed idle shell with no pending managed agent, foreground
provider or descendants. Unavailable process inspection prevents reuse. When the
pane is missing, a successful global pane inventory must also show that its
terminal has not moved elsewhere. Ply then creates a new Task tab in the current
Herdr workspace (or `--herdr-workspace`). Missing terminal evidence does not prove
that every old OS process exited: replacement depends on the old attempt never
having been authorized to receive Task input.

Both paths require a clean Task on its original base, with no Task-prompt attempt,
acceptance, start receipt or delivery result. A readiness session alone is not
Task authority. Changed goal inputs, an active or unknown existing startup,
possible Task input, or unconfirmed creation of a replacement tab prevent another
start. Closing a tab after implementation started does not qualify it for this
startup recovery.

Each replacement is recorded under a new generation. Original requests, contexts,
control binaries, transport identities and failed attempts remain available.
A delayed controller or callback from an older generation cannot send the Task or
change the current attempt. A later interrupted startup can be recovered again
with the same Spec command when these checks establish that it is safe; there is
no single-use recovery limit. An uncertain effect is kept for inspection rather
than silently retried. `show --format json` exposes `startup_recovery`, including
the original and effective executable and transport bindings.

For callers that already have a run ID, `execute recover-start RUN_ID` provides
the same guarded startup recovery, with optional `--check` and `--timeout`.
Ordinary `resume` still uses the bound provider and never restarts a process.
Startup recovery does not release the queue, delete evidence, change provider
permissions or qualify the product. The actual human product pass and local
integration still belong to the owner's delivery journey.
