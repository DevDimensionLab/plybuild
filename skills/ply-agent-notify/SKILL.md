---
name: ply-agent-notify
description: Notify an agreed Slack route through Ply when a human-requested task finishes, stops, or needs an answer. Bind ownership at task start and preserve delivery identity across continuations.
---

# Ply Agent Notify

## 1. Authority and owner

Use the user's explicit request or standing preference for the named route as outgoing authority.
Skill discovery alone grants no sending permission. Root owns the notification for a human-requested
task. Internal delegates return to root without separate sends unless the user agreed to a separate
delivery boundary.

An event that the removed bound-Task entrance already started is history; never send it again
through the short form (see [usage](references/usage.md)).

## 2. Send one short event

Run from the job's actual working directory:

```sh
python3 /Users/perottochristensen/.agents/skills/ply-agent-notify/scripts/notify.py \
  --activity stable-job-id --run stable-run-id --event-id finished --status done --provider claude \
  --title "ℹ️ Completed task" --summary "Concrete result" \
  --next-action "No action required." --ply-bin "$PLY_BIN"
```

`--provider` is `codex` or `claude` from the actual runtime, not a model name or executable lookup.
Pass the already-resolved absolute Ply executable used by the job as `PLY_BIN`
(here `/Users/perottochristensen/go/bin/ply`); no global PATH change. If it is unavailable, report
`not_attempted` with the missing executable. Credential setup requires a separate installation request.

Optional: `--next-actor coordinator` (default `user`),
`--phase before_start` (only for a stop), `--context-root`, `--timezone` (default Europe/Oslo) and
`--route` (default this user's ply-log).

The event identity is the whole `activity`/`run`/`event-id` tuple on one agreed route; it is immutable
across continuation. The event ID names one logical completion, stop or pending question. Reuse the
same command across restarts: a replay only reads the result and never sends again. A new event gets
a new ID. Changed text or origin on a used identity, or invalid input, reports `unknown` without a send.
Another run, event ID, route file directory or state directory is another namespace that the helper
cannot connect to this event, and may send again: changing any of them to retry or to bypass an
unknown outcome is forbidden. The first call freezes the event's instant; the helper keeps only its
private request and attempt marker in the route's `notification-events` directory. The message shows
time and provider, not the directory. Inside Herdr, a proven origin adds its workspace and tab
display labels (frozen at the first call); an unproven one is sent without them. Public text is
filtered only for mentions, webhooks, the credential and unsafe characters; you are responsible for keeping other private details out of it.

| Status | Use |
| --- | --- |
| `done` | The agreed task is complete. |
| `ready_for_review` | A return awaits coordinator review. |
| `ready_for_your_check` | Technical control is complete and a human must judge the product. |
| `stopped` | Preserve available work and state the blocker. |
| `needs_answer` | Necessary question, with `?` in the summary and `next_actor=user`. Preserve it before waiting. |

Optional clarification, internal corrections, and test failures that can be resolved within the
task do not send notifications. An ending still waiting for the same answer reuses that question
and uses `--trigger waiting_for_feedback`, not another stop/completion event.

For **new events**, prefix the title with ⚠️ when `next_actor=user` and
status is `needs_answer`, `ready_for_your_check`, `ready_for_review` or `stopped`.
Use ℹ️ for ordinary `done`/status notices and coordinator actions. The indicator
counts toward the title's 60-character limit; keep status and action readable in
plain text too. Never add an icon or change a title on continuation.

When `ply-human-input` has published a necessary Q, that skill owns the complete
question and answer. Reuse this job's authorized route and activity with one Q-derived
event ID `human-input/QUESTION_ID`. Immediately
before invocation, the owner reads the same Q; answered/withdrawn needs no alert.
For this branch only, the `needs_answer` summary is an alert containing
`?`, such as `A necessary question is waiting. Can you answer it in /plyapp?`;
the next action is `Open /plyapp. Select a question if needed, choose or write your answer, then confirm the exact answer with passkey.`
Keep the full question/reason/consequences in Ply/app rather than truncating or
copying them into Slack. Do not export private details or answer text. Missing route
authority is `not_attempted`, without setup or a new permission round. Other feedback
continues to preserve the actual necessary question in full.

Public text is single-line, with title/summary/action limits of 60/100/140 Unicode characters.
The notice remains three lines and at most
480 Unicode characters including indicators. Keep secrets, webhook URLs, mentions, raw reports,
and transcripts out of public fields. Preserve a necessary question in full except for the
Ply Q alert branch above. The next actor is `coordinator`
for coordinator continuation or `user` for a user return. When no follow-up is required, use
`No action required.` Necessary answers belong in the active conversation, or in `/plyapp`
with passkey for a published Q. Slack thread replies are not connected to either.

The helper uses the existing environment credential or named private file, checks the send,
applies at most once, and reads back the result. It does not automatically retry.

## 3. Report delivery separately

The notification owner invokes the helper before its final response or genuine wait.
Internal coding/review steps stay local while the coordinator can continue.
Report the helper's JSON `state`, including when its exit code is zero. An acknowledged transport
is not a read receipt, approval, or completed product QA. Preserve the task result and necessary
question even when sending fails. An unknown result stays unknown; reuse the same command for
readback rather than generating a new send. The helper cannot notify after hard process death
or a closed tab.
