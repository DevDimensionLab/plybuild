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
| `executor` | Explicit nullable `provider`, `model`, `effort` choices; optional boolean `allow_local_install` |
| `delivery` (optional) | Versioned PR, local Epic, or explicit local branch agreement; see [Delivery](workflow-delivery.md) |
| `change_reason` | Why this goal revision was published |

A goal has at least one design reference and requirement. It has no
`implementation_basis`, detailed technical parts, completed-test claims, or human
solution selection. The document store preserves actual bytes and digests using
the existing file/snapshot source rules. Revised goals get new immutable revisions.

New interactive executions grant the owner local installation only when the
selected Goal includes `"allow_local_install": true` in `executor`. Omission or
false reserves installation outside the owner mandate. For example:

```json
{"provider":"codex","model":null,"effort":null,"allow_local_install":false}
```

Choose this field from the human's actual scope before recording the Goal.
Delivery mode, provider choice and prose constraints do not grant installation.
The frozen selection supplies the generated mandate and initial owner prompt;
draft/request creation rejects an option that widens it. Headless workers still
cannot install, even when the Goal allows an interactive owner to do so. Separate
reviewers retain their review-only role. Existing frozen requests and contracts
keep their original bytes and authority; a revised Goal does not amend an old Run.

`workflow execute verify` snapshots and runs the acceptance shell script, captures
its output and qualifies candidate evidence. It does not itself perform ordinary
installation or integration. Interactive automatic acceptance supplies an isolated
HOME/TMPDIR and a restricted environment, but does not impose an OS sandbox: the
script can exercise the caller's filesystem and process rights. Review its actual
effects against the mandate. The installation field is authority, not a shell
containment mechanism. Worker verification retains its separate sandbox policy.

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

After the cause is clarified, the same owner names the changed action it will
take under the same Run mandate:

```shell
/absolute/ply-control workflow execute resolve-stop RUN_ID --context /absolute/context.json \
  --stop sha256:ACTUAL_STOP_DIGEST --resolution 'Actual clarified cause and who clarified it' \
  --action 'Changed action this owner takes now'
```

No permission confirmation is asked for: runtime still reviews the named action
when it happens. The refused action itself cannot be named again, and unknown or
partial effects cannot be cleared with text. Resolution requires fresh identity
checks and refuses an unresolved attempt. It preserves the denial, starts
nothing and sends no input. If acceptance is absent, the owner uses the startup
command selected by the frozen callback guide before Task writes. Older stops
resolved with the earlier `--permission-confirmed` flag remain readable history.
Neither operation releases a writer; selected replacement uses `execute retire`. These owner
commands must be present in the bound control; do not overwrite an older control
to acquire them. Existing accepted Runs can use the supported compatible-control
continuation when needed. Historical denied Runs remain preserved.

Native folder trust remains separate. If its existing configuration lock cannot
be acquired, the same Claude tab may require its native trust decision before
`execute resume`. Preserve that action as a startup limitation; do not retry the
launch or change permission mode. FL01 does not repair that configuration lock.

### Retire a selected owner for a fresh Task copy

For a misclassified **pre-execution Claude quota stop**, a coordinator can add
one source-checked retirement annotation. This applies only when the stop's
partial/unknown label described unfinished work and the named continuation
never started after quota. It cannot clear an earlier tool's uncertain effects.

```shell
ply workflow execute reconcile-quota-stop RUN_ID \
  --event-id quota-correction-01 --stop sha256:EXACT_STOP_DIGEST \
  --provider claude --native-session EXACT_SESSION \
  --action 'Exact original stopped continuation action' \
  --quota-event NATIVE_QUOTA_UUID --last-tool-result LAST_BASH_RESULT_UUID \
  --source-sha256 sha256:COMPLETE_NATIVE_TRANSCRIPT_DIGEST \
  --observer 'Actual coordinator' \
  --reason 'Why unfinished work was mistaken for an executed continuation'
```

Use `execute show RUN_ID --format json` to inspect the bound session and stop.
Inspect only necessary native transcript structure, then compute its full fresh
SHA-256 (for example `shasum -a 256 /absolute/native/session.jsonl`). Ply locates
that session through its existing provider-session locator; it accepts no
arbitrary evidence path. The native file must be a private physical regular
file, at most 32 MiB. It must bind the exact session, with the selected Bash
call/result and quota in the exact Task cwd, and contain the
selected `usage_limit_reached` API rejection (429), and end the preceding tool
activity at the selected synchronous Bash result with `interrupted=false`.
The selected result starts the precise stop window: no later tool call/result
or assistant activity is allowed apart from the exact quota rejection. At that
boundary every observed tool call must have a matching return, every async
Agent launch needs its matching native `session-task` completion notification,
and no interrupted tool or unsupported background effect may remain. Counts, peer
handbacks, tool output and free text cannot replace that notification. No later
tool use/result, task completion or assistant activity is permitted. Rule
refusals and unresolved native attempts remain blocked.

Matched returns before this boundary remain historical source only; the audit's
fixed scope explicitly says their effects are **not assessed**. A metadata gap
or error string in such an earlier return is not the named next continuation;
historical cwd changes likewise remain source rather than a new effect claim.
It cannot itself establish the selected boundary: that requires the exact
successful synchronous native Bash result with `interrupted=false`. Nor can
`is_error`, text or a later successful result hide a tool that actually started
inside the window. Selecting an earlier result fails if any tool ran afterwards.
Every other unresolved native unknown/partial stop still blocks retirement,
including one concerning an earlier action.

The immutable event binds Run, request, stop, action, provider, session, exact
native event IDs, original source path/hash and a private copied transcript.
Its scope attests only the unstarted continuation after the selected result;
it grants no historical effect attestation.
Original stop/evidence bytes remain unchanged. The same event/input retry is
idempotent, including after a lost response; different input cannot overwrite
it. A changed source requires a new event ID and fresh full validation. The
latest valid audit for the exact stop is considered; older snapshots stay
immutable. A new ID for an already audited identical source is rejected. The annotation keeps the
stop active and grants no acceptance, continuation, owner input or release.
`resolve-stop` retains its owner-only behavior and still refuses partial/unknown
effects. Historical contracts remain readable; older controls fail closed.

Only ordinary `execute retire` consumes this annotation. It rechecks the exact
unchanged native source before pane close and again after actual close and
writer absence, considers every other unresolved stop/attempt, and still
proves all writers absent before queue release. A
changed or missing native source remains blocked, even if a correction retry
returns its preserved observation. A close may already have happened: the
error preserves before/after source hashes, records no retirement and releases
no queue slot. Repeating retirement with the stale audit stops before another
close. If fresh inspection shows only benign final metadata, use a new audit
event ID and its fresh full source hash; ordinary retirement then observes the
already absent pane without replaying close. Actual new tool/assistant activity
fails that audit and requires its real effect to be established through the
existing stop workflow; this command cannot clear it. Preserve that sourced
observation with `observe-stop` when needed. There is no force or generic
partial-to-none operation.

When the coordinator selects replacement (for example standing cattle or
Codex backup), one operation from the containing workspace ends the old owner:

```shell
ply workflow execute retire RUN_ID --observer 'Actual coordinator' \
  --reason 'Selected replacement and its source'
```

Ply closes the exact bound Herdr pane itself, then requires the pane, its
managed agent and every known Task writer to be positively absent. It records
the retirement, closes every callback of the retired owner and releases only
this Task's queue slot. No last owner callback, private context, manual journal
edit, `queue release` or rebind is needed. The Run, its refusals, reviews and
Task worktree stay as history; useful code can be reused as reviewed input. A
fresh Task copy starts through ordinary goal publication and `execute` with its
own new mandate and inherits no acceptance.

An unresolved verification or integration attempt, a preserved stop with
unknown or partial effects, or a writer that is still observed stops retirement
with a concrete diagnosis before anything is recorded or released. A stopped
history whose effects are known does not block it. Repeating the same input
after a lost reply is safe; a different reason cannot overwrite the record.

## Follow the owner

```shell
ply workflow execute show wfr_<digest>
ply workflow execute follow wfr_<digest> --timeout 60
```

Show is read-only. Follow observes the same provider/session. It checks the
bound pane, terminal and native session, not the owner's current foreground
cwd: a transient read outside the Task does not hide that the same session is
working, idle or done. Owner callbacks still require the exact Task cwd. A timeout ends
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

### Independent review rounds

Review is an ordinary transition of the same Task and Run. When the owner has a
clean committed candidate, the coordinator runs, from the containing workspace:

```shell
ply workflow execute review RUN_ID --coordinator 'Actual coordinator'
```

Ply binds that exact candidate as the next numbered round and writes a reviewer
brief, a return template and a reviewer-only return context under the Run's
`review/rounds/NNN/`. No review Task, Spec, assessment, selection, queue move,
preparation or start request is written, and the owner keeps the queue slot.
The same command starts the separate reviewer through Ply: a new visible Herdr
tab in the owner's workspace (or `--herdr-workspace`), in the Task worktree,
with the Run's provider, model and permission mode. Ply binds the reviewer's
observed native session, which must differ from the owner's, and sends it the
brief. The owner's own pane cannot start its review. Each start effect is
reserved before it is sent and never replayed. Repeating the command for the
same candidate returns the existing round. An unknown start is shown as a stop;
after inspecting the tab, `--new-round` starts a fresh reviewer for the same
candidate. The replaced round is recorded as superseded: its history and
unknown start are kept, but its reviewer can no longer return and the fresh
round decides. When new evidence resolves a sealed `blocked` review with stated
evidence gaps and no findings, the coordinator explicitly uses `--new-round`
to reassess the same exact candidate. The original request, return, evidence
and reviewer session remain unchanged history. Ordinary repeats remain
idempotent, and the fresh round must return before verification can use it.
Known product findings, `changes_requested`, and blocks without stated evidence
gaps still govern the unchanged candidate and require a corrected commit.
An earlier sealed callback can only repeat its original bytes; it cannot
change the current judgment.

From its own pane, the reviewer returns an ordinary JSON file with `decision`
(`passed`, `changes_requested` or `blocked`), `summary`, `findings`
(`severity`, `summary`, optional `evidence`), `evidence_gaps` and an optional
`reviewer_claim`. Ply checks the calling pane and that pane's current session,
and records the session it observed:

```shell
ply workflow execute review-return RUN_ID --context /absolute/review/rounds/001/context.json --file /absolute/my-review.json
```

A complete negative review is a completed review: status shows the product
decision and its findings, and truthful evidence gaps are kept. Ply assigns
finding IDs and ordering and accepts ordinary file modes and letter case. A
sealed round cannot be overwritten; a changed judgment needs a new round.

Findings return to the owner, who corrects, commits and asks the coordinator
for a new round. That round binds the new candidate explicitly and keeps the
earlier rounds and their evidence. If the owner commits while a round is still
open, the return stays bound to the reviewed candidate and status shows the
drift. Product findings keep governing their exact candidate; an evidence-only
block governs until the coordinator explicitly supersedes it for reassessment. `verify`
without `--review` uses the completed passed round for exactly the current
candidate and carries earlier candidates' findings as fixes; any other
candidate has no pass to inherit. `execute show` lists each round, its findings
and gaps, and the next responsible actor.

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
