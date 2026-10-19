---
name: ply-agent-notify
description: Notify an agreed Slack route through Ply when a human-requested task finishes, stops, or needs an answer. Bind ownership at task start and preserve delivery identity across continuations.
---

# Ply Agent Notify

## 1. Bind the authorized task

Use the user's explicit request or standing preference for the named route as outgoing authority.
Reuse it across continuations; skill discovery alone grants no sending permission. Root owns the
notification for a human-requested task. Internal delegates return to root without separate sends
unless the user agreed to a separate delivery boundary.

For an assigned native delivery with an existing authorized route, use the bound
control's contract-selected startup command (`check-in` or legacy `begin`) with `--notification-route`,
`--notification-authority` and `--notification-timezone`. Ply creates the private
Task@2 binding with the actual provider, Task cwd and native activity/run identity.
Reuse the returned `delivery/notification-task.json`; do not create a second
binding. This prepares no transport and changes no global configuration.
If the frozen request already supplies a notification context, reuse that context.
Internal delegates whose goal gives notifications to a parent omit these flags.

For other new jobs, copy [the task context template](references/task-template.json) into the task's
return directory. Bind the activity/run, target worktree, mandate, route, fixed event root, actual
provider, physical origin/context root, and recipient timezone. Use `codex` or `claude` from the
actual agent runtime, not a model name or executable lookup. Record ownership and the authority
source in the mandate; the JSON schema accepts no extra fields. Europe/Oslo is this user's zone.

Check that the configured Ply executable is available (`command -v ply` for the default). If it
is unavailable, preserve the binding and local result and report `not_attempted` with the missing
executable. Use an explicitly located executable through `--ply-bin` when appropriate; preserve
the agreed route, state root, and runtime permissions. Credential setup requires a separate
installation request. For installation/upgrade or an existing Task@1 record, read the relevant
section of [installation and legacy usage](references/usage.md).

## 2. Preserve the result and event

A continuation or answer reuses the task binding. Preserve the same event file and ID for the
same logical completion or pending question, including across restarts. A new event gets a new
ID. Never alter an ID, route, or state directory to bypass an unknown outcome.

Use the actual local JSON report for a completion event, including a native delivery report.
Pass that file directly; do not write another notification report or add fields to the original.
The helper binds missing notification identities from the Task and rejects explicit mismatches,
including native `activity_id`/`run_id`. Include an actual start receipt when available.
Choose the event and public status from the observed delivery boundary:

| Event | Status and required local evidence |
| --- | --- |
| `agent_finished` | `done` for the agreed completed task; `ready_for_review` for coordinator review; `ready_for_your_check` for human product judgment. Requires the report. |
| `agent_stopped` | `stopped`. Preserve available work and state the blocker. `before_start` permits no start receipt and does not require a report. |
| `feedback_required` | `needs_answer`, `after_start`, `next_actor=user`. Preserve the actual necessary question, including `?`, before waiting. |

Optional clarification, internal corrections, and test failures that can be resolved within the
task do not send notifications. An ending still waiting for the same answer reuses that question
and uses `--trigger waiting_for_feedback`, not another stop/completion event.

For **new events**, prefix `public.task_title` with ⚠️ when `next_actor=user` and
status is `needs_answer`, `ready_for_your_check`, `ready_for_review` or `stopped`.
Use ℹ️ for ordinary `done`/status notices and coordinator actions. The indicator
counts toward the title's 60-character limit; keep status and action readable in
plain text too. Preserve older event/request bytes exactly: never add an icon or
change a title on continuation. The renderer and historical payloads stay unchanged.

When `ply-human-input` has published a necessary Q, that skill owns the complete
question and answer. Reuse this job's authorized binding and one Q-derived event
`human-input/QUESTION_ID`, with its preserved file and UTC instant. Immediately
before invocation, the owner reads the same Q; answered/withdrawn needs no alert.
For this branch only, the compact `feedback_required` summary is an alert containing
`?`, such as `A necessary question is waiting. Can you answer it in /plyapp?`;
`next_action` is `Open /plyapp. Select a question if needed, choose or write your answer, then confirm the exact answer with passkey.`
Keep the full question/reason/consequences in Ply/app rather than truncating or
copying them into Slack. Do not export private details or answer text. Missing route
authority is `not_attempted`, without setup or a new permission round. Other feedback
continues to preserve the actual necessary question in full.

For Task@2, preserve this event JSON, replacing every example value. The helper supplies omitted
`kind`, `activity`, and `run` from the bound Task; if supplied, they must match it.

```json
{
  "event_id": "stable-event-id",
  "event_type": "agent_finished",
  "phase": "after_start",
  "occurred_at": "actual UTC instant in YYYY-MM-DDTHH:MM:SSZ form",
  "public": {
    "status": "done",
    "task_title": "ℹ️ Completed task",
    "summary": "Concrete result",
    "next_action": "No action required.",
    "next_actor": "user"
  }
}
```

Public text is single-line, with title/summary/action limits of 60/100/140 Unicode characters.
Derived context is limited to 96 characters. The notice remains three lines and at most
480 Unicode characters including indicators. Keep secrets, webhook URLs, mentions, raw reports,
and transcripts out of public fields. Preserve a necessary question in full except for the
Ply Q alert branch above. Use the actual UTC
instant once; the helper freezes its local timestamp and offset. The next actor is `coordinator`
for coordinator continuation or `user` for a user return. When no follow-up is required, use
`No action required.` Necessary answers belong in the active conversation, or in `/plyapp`
with passkey for a published Q. Slack thread replies are not connected to either.

## 3. Invoke the helper

Run from the task's physical `origin_cwd`. The event root's parent must already exist.
For a completion:

```sh
python3 /Users/perottochristensen/.agents/skills/ply-agent-notify/scripts/notify.py \
  --task /absolute/task/notify-task.json --event /absolute/task/finished.json \
  --report /absolute/task/report.json --ply-bin "$PLY_BIN"
```

Pass the already-resolved absolute Ply executable as `PLY_BIN`, including when it
is absent from the helper's PATH (here `/Users/perottochristensen/go/bin/ply`).
Reuse the executable used for the job's Ply commands; no global PATH change.

Use the preserved stop/question event for those cases and supply report/start-receipt paths only
when they exist. `agent_finished` always needs its report. All paths are structured arguments.
The helper uses the existing environment credential or named private file, checks the send,
applies at most once, and reads back the result. It preserves original bytes/digests and any
needed binding envelopes in the existing private event directory. It does not automatically retry.

## 4. Report delivery separately

The notification owner invokes the helper before its final response or genuine wait.
Internal coding/review steps stay local while the coordinator can continue.
Report the helper's JSON `state`, including when its exit code is zero. An acknowledged transport
is not a read receipt, approval, or completed product QA. Preserve the task result and necessary
question even when sending fails. An unknown result stays unknown; reuse the original event for
readback rather than generating a new send. The helper cannot notify after hard process death
or a closed tab.
