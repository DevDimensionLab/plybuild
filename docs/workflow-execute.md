# Deliver a queued goal

`ply workflow execute` gives one interactive implementor responsibility for a
feature from its goal through an actual human product check and local return.
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

The first version supports one registered Epic in one repository. A planner may
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

The launch-policy snapshot describes requested settings. It is not an OS
permission grant. Before Task writes, the recipient reports its actual runtime,
native session and permission evidence. Unknown authority stops dependent writes.
Neither a goal nor the owner's broad delivery role changes sandbox enforcement.
Old Agreement A runs keep their original budget; a new delivery owner does not
inherit a hidden correction-count or time limit.

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
`verify`, `qa` and `integrate` under `workflow execute`; runtime acceptance uses
`workflow run accept`. Consult each command's help for its required evidence.

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

The owner presents one installed product journey and preserves the actual human
answer as `DeliveryHumanAttestation@1`. That record binds `task_id`,
`task_result_id`, `result_oid`, `result_tree`, `outcome`, `actor_claim`,
`start_surface`, UTC start/completion timestamps, `answer` and `observation`.
It is a local human attestation, not cryptographic identity proof. Silence,
technical success or agent approval never supplies the answer. A fail permits
correction and a new immutable candidate generation in the same session.

Only a matching pass permits the owner to integrate the exact candidate into the
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
